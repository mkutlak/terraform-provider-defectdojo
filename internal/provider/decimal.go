package provider

import (
	"context"
	"fmt"
	"math/big"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// decimalPattern is the literal grammar the schema accepts: an optional sign,
// then digits with an optional fractional part.
//
// It gates big.Rat.SetString, which is more permissive than DefectDojo: it also
// reads "1e3", "1/2" and "0x1p-2".
var decimalPattern = regexp.MustCompile(`\A[+-]?(?:\d+(?:\.\d*)?|\.\d+)\z`)

// parseDecimal reads a plain decimal literal exactly.
//
// big.Rat is arbitrary precision. The widest such column in the API is
// DecimalField(max_digits=15, decimal_places=2) - product.revenue - so a value
// can need 15 significant digits, past the range where float64 is exact.
func parseDecimal(s string) (*big.Rat, bool) {
	if !decimalPattern.MatchString(s) {
		return nil, false
	}
	return new(big.Rat).SetString(s)
}

// decimalsEqual reports whether two literals denote the same amount.
//
// DRF renders a DecimalField with decimal_places digits after the point,
// always, so a configured "100" comes back as "100.00". The two are the same
// amount and the configured spelling is kept. See preserveLiteral.
func decimalsEqual(current, server string) bool {
	a, aok := parseDecimal(current)
	b, bok := parseDecimal(server)
	return aok && bok && a.Cmp(b) == 0
}

// decimalValidator refuses at plan time every literal outside the grammar
// above, plus every literal outside the column's own limits.
type decimalValidator struct {
	// maxWholeDigits and maxDecimalPlaces mirror the Django
	// DecimalField(max_digits, decimal_places) behind the column. DRF checks
	// both separately, so both are reported separately.
	maxWholeDigits   int
	maxDecimalPlaces int
}

func (v decimalValidator) Description(_ context.Context) string {
	return fmt.Sprintf("must be a decimal number with at most %d digits before the decimal point "+
		"and %d after it", v.maxWholeDigits, v.maxDecimalPlaces)
}

func (v decimalValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v decimalValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	literal := req.ConfigValue.ValueString()
	if _, ok := parseDecimal(literal); !ok {
		// The empty string is worth its own message. It is what
		// `revenue = var.x` produces for an unset variable, which is the
		// ordinary Terraform spelling of "leave this alone", and the server
		// takes the write, so the practitioner has no way to tell from the
		// failure what to do instead.
		if literal == "" {
			resp.Diagnostics.AddAttributeError(req.Path, "Invalid Decimal Value", fmt.Sprintf(
				"%s is an empty string, which is not a decimal number. DefectDojo answers a "+
					"write of \"\" with a 201 and stores the column as null, so the apply then "+
					"fails with \"Provider produced inconsistent result after apply: .%s: was "+
					"cty.StringVal(\"\"), but now null\" - and fails the same way on every "+
					"retry. Omit the attribute to leave the value unset.",
				req.Path, req.Path))
			return
		}
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid Decimal Value", fmt.Sprintf(
			"%s is %q, which is not a decimal number. Write a plain decimal literal such as "+
				"\"100\", \"1.5\" or \"-2.50\"; exponents, thousands separators and surrounding "+
				"whitespace are not accepted.",
			req.Path, literal))
		return
	}

	// DRF counts whole digits on the parsed Decimal, so leading zeros do not
	// count against the limit and "000000000000005" is a good 1-digit amount.
	whole, fraction, hasPoint := strings.Cut(strings.TrimPrefix(strings.TrimPrefix(literal, "+"), "-"), ".")
	if whole = strings.TrimLeft(whole, "0"); len(whole) > v.maxWholeDigits {
		resp.Diagnostics.AddAttributeError(req.Path, "Decimal Value Out Of Range", fmt.Sprintf(
			"%s is %q, which carries %d digits before the decimal point; DefectDojo accepts at "+
				"most %d and answers the write with 400 \"Ensure that there are no more than %d "+
				"digits before the decimal point.\".",
			req.Path, literal, len(whole), v.maxWholeDigits, v.maxWholeDigits))
		return
	}

	// Decimal places are counted on the literal instead, because the server
	// counts them on the literal too: "1.000" is one amount to big.Rat, but DRF
	// reads it as a three-place Decimal and refuses it.
	if hasPoint && len(fraction) > v.maxDecimalPlaces {
		resp.Diagnostics.AddAttributeError(req.Path, "Decimal Value Out Of Range", fmt.Sprintf(
			"%s is %q, which carries %d digits after the decimal point; DefectDojo accepts at "+
				"most %d and answers the write with 400 \"Ensure that there are no more than %d "+
				"decimal places.\".",
			req.Path, literal, len(fraction), v.maxDecimalPlaces, v.maxDecimalPlaces))
	}
}
