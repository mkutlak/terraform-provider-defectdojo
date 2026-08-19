package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DefectDojo stores money-ish fields as Django DecimalField columns, and DRF
// serialises them in exactly one canonical form: decimal_places digits after
// the point, always. So a configured revenue of "100" comes back as "100.00",
// and "1.5" comes back as "1.50".
//
// That is issue #23 again in a different costume - a permissive input grammar
// with a single canonical output form. Without preservation, `revenue = "100"`
// fails every apply with "Provider produced inconsistent result after apply",
// and the practitioner cannot fix it except by guessing the canonical spelling.

// normalizeDecimal reduces a decimal literal to a canonical comparison key:
// sign, integer digits with leading zeros stripped, fractional digits with
// trailing zeros stripped, and any flavour of zero collapsed onto "0". It
// reports false for anything that is not a plain decimal literal, so callers
// fall back to the server's own rendering rather than guessing.
//
// String normalisation rather than arithmetic is deliberate. The widest such
// column in the API is DecimalField(max_digits=15, decimal_places=2) -
// product.revenue, whose spec pattern is ^-?\d{0,13}(?:\.\d{0,2})?$ - which
// needs 15 significant digits, past the range where float64 is exact.
// TestNormalizeDecimalIsExactBeyondFloat64 pins that.
func normalizeDecimal(s string) (string, bool) {
	if s == "" {
		return "", false
	}

	sign := ""
	switch s[0] {
	case '+':
		s = s[1:]
	case '-':
		sign, s = "-", s[1:]
	}

	intPart, fracPart, _ := strings.Cut(s, ".")
	if strings.Contains(fracPart, ".") {
		return "", false // more than one decimal point
	}
	if !isASCIIDigits(intPart) || !isASCIIDigits(fracPart) {
		return "", false
	}
	if intPart == "" && fracPart == "" {
		return "", false // "", "+", "-" and "." all land here
	}

	if intPart = strings.TrimLeft(intPart, "0"); intPart == "" {
		intPart = "0"
	}
	fracPart = strings.TrimRight(fracPart, "0")

	switch {
	case intPart == "0" && fracPart == "":
		return "0", true // collapses "-0", "+0.00" and "0.000"
	case fracPart == "":
		return sign + intPart, true
	default:
		return sign + intPart + "." + fracPart, true
	}
}

// isASCIIDigits reports whether s consists solely of ASCII digits. The empty
// string qualifies on purpose, so ".5" and "5." both normalise.
func isASCIIDigits(s string) bool {
	for i := range len(s) {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// preserveDecimalLiteral picks the string to store in state for a
// server-provided decimal. When the attribute already holds a literal denoting
// the same amount - the practitioner's configured value on create/update, or
// the prior state value on refresh - that literal is kept verbatim. Otherwise
// the server's rendering is used.
//
// A genuinely different amount is never preserved, so real drift is still
// reported. This is the decimal counterpart of preserveDateTimeLiteral.
func preserveDecimalLiteral(current types.String, server string) types.String {
	if !current.IsNull() && !current.IsUnknown() {
		if want, ok := normalizeDecimal(current.ValueString()); ok {
			if got, ok := normalizeDecimal(server); ok && want == got {
				return current
			}
		}
	}

	return types.StringValue(server)
}

// decimalValidator refuses at plan time every literal preserveDecimalLiteral
// could not protect, plus every literal outside the column's own limits.
//
// It reads the accepted grammar off normalizeDecimal rather than restating it
// as a pattern, because the two are one rule seen from two ends: whatever the
// schema admits, the read path has to be able to key, or the practitioner gets
// a value that cannot survive its own apply. Restating it invites exactly the
// drift this replaces - a pattern whose quantifiers were all {0,n} and so
// admitted "", "-", "." and "-.", none of which normalizeDecimal parses.
//
// The containment is one-way on purpose. normalizeDecimal keys "1.000", which
// the column's two decimal places refuse with a 400, so the digit limits below
// are genuinely narrower than the grammar.
// TestProductRevenueValidatorAcceptsOnlyNormalisableLiterals pins the direction
// that matters.
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
	normalised, ok := normalizeDecimal(literal)
	if !ok {
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

	// Count whole digits on the normalised form. DRF counts the digits of the
	// parsed Decimal, so leading zeros do not count against the limit and
	// "000000000000005" is a perfectly good 1-digit amount.
	whole, _, _ := strings.Cut(strings.TrimPrefix(normalised, "-"), ".")
	if len(whole) > v.maxWholeDigits {
		resp.Diagnostics.AddAttributeError(req.Path, "Decimal Value Out Of Range", fmt.Sprintf(
			"%s is %q, which carries %d digits before the decimal point; DefectDojo accepts at "+
				"most %d and answers the write with 400 \"Ensure that there are no more than %d "+
				"digits before the decimal point.\".",
			req.Path, literal, len(whole), v.maxWholeDigits, v.maxWholeDigits))
		return
	}

	// Decimal places are counted on the literal instead, because the server
	// counts them on the literal too: normalizeDecimal keys "1.000" as "1",
	// but DRF reads it as a three-place Decimal and refuses it.
	if _, fraction, hasPoint := strings.Cut(literal, "."); hasPoint && len(fraction) > v.maxDecimalPlaces {
		resp.Diagnostics.AddAttributeError(req.Path, "Decimal Value Out Of Range", fmt.Sprintf(
			"%s is %q, which carries %d digits after the decimal point; DefectDojo accepts at "+
				"most %d and answers the write with 400 \"Ensure that there are no more than %d "+
				"decimal places.\".",
			req.Path, literal, len(fraction), v.maxDecimalPlaces, v.maxDecimalPlaces))
	}
}
