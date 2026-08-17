package provider

import (
	"strings"

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
