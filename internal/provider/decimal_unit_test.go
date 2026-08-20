package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestParseDecimal(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		in     string
		wantOK bool
	}{
		{"100", true},
		{"100.00", true},
		{"100.50", true},
		{"0100.50", true},
		{"+7", true},
		{"-1.50", true},
		{"0", true},
		{"-0", true},
		{"+0.000", true},
		{".5", true},
		{"5.", true},
		{"0000", true},

		// Rejected. big.Rat.SetString would read the exponent, ratio and hex
		// forms, which DefectDojo does not accept; decimalPattern gates them.
		{"", false},
		{".", false},
		{"+", false},
		{"-", false},
		{"1.2.3", false},
		{"abc", false},
		{"1e3", false},
		{"1/2", false},
		{"0x1p-2", false},
		{" 100", false},
		{"100 ", false},
		{"1,000.00", false},
		{"NaN", false},
	} {
		if _, ok := parseDecimal(tc.in); ok != tc.wantOK {
			t.Errorf("parseDecimal(%q) ok = %v, want %v", tc.in, ok, tc.wantOK)
		}
	}
}

func TestDecimalsEqual(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		a, b string
		want bool
	}{
		{"100", "100.00", true}, // the DRF rendering of a configured "100"
		{"1.5", "1.50", true},
		{"-2.5", "-2.50", true},
		{"0", "-0.00", true},
		{".5", "0.50", true},
		{"000000000000005", "5.00", true},
		{"99", "100.00", false},
		{"100", "not a number", false},
		{"not a number", "100.00", false},

		// DecimalField(max_digits=15, decimal_places=2) needs 15 significant
		// digits, past the point where float64 can tell neighbours apart. big.Rat
		// is exact, so these stay distinct.
		{"1234567890123.45", "1234567890123.46", false},
	} {
		if got := decimalsEqual(tc.a, tc.b); got != tc.want {
			t.Errorf("decimalsEqual(%q, %q) = %v, want %v", tc.a, tc.b, got, tc.want)
		}
	}
}

// revenueValidationCases mirrors what DefectDojo 3.1.101 actually does with
// product.revenue. Every accepted value below was POSTed to /api/v2/products/
// and read back; every rejected one either fails with a 400 or is stored as
// something the configuration can never match.
var revenueValidationCases = []struct {
	revenue   string
	wantError bool
}{
	// Accepted, stored in DRF's canonical two-place form, and folded back onto
	// the configured literal on read.
	{"100", false},              // -> "100.00"
	{"100.00", false},           // -> "100.00"
	{"1.5", false},              // -> "1.50"
	{"-2.5", false},             // -> "-2.50"
	{"0", false},                // -> "0.00"
	{"-0", false},               // -> "-0.00"
	{"+5", false},               // -> "5.00"
	{".5", false},               // -> "0.50"
	{"5.", false},               // -> "5.00"
	{"000000000000005", false},  // -> "5.00"; DRF counts significant digits
	{"1234567890123.45", false}, // the widest amount the column holds

	// 201, with the column stored as null - so every apply fails with
	// "Provider produced inconsistent result after apply: .revenue: was
	// cty.StringVal(""), but now null". `revenue = var.x` with an unset
	// variable is the ordinary way to reach it.
	{"", true},

	// 400 "Ensure that there are no more than 2 decimal places."
	{"1.000", true},
	{"12.3456", true},
	// 400 "Ensure that there are no more than 13 digits before the decimal point."
	{"12345678901234", true},

	// Not decimal literals at all; TestParseDecimal covers the grammar in full.
	{"abc", true},
	{"1e3", true},
	{"1,000.00", true},
}

func TestProductRevenueValidator(t *testing.T) {
	t.Parallel()

	attr := resourceStringAttribute(t, "defectdojo_product", "revenue")
	for _, tc := range revenueValidationCases {
		got := runStringValidators(t, "revenue", attr, tc.revenue)
		if got == tc.wantError {
			continue
		}
		verb := "accepted"
		if got {
			verb = "rejected"
		}
		t.Errorf("revenue = %q was %s, but wantError=%v", tc.revenue, verb, tc.wantError)
	}
}

// TestProductRevenueValidatorAcceptsOnlyParseableLiterals pins the property the
// schema and the read path have to share: every literal the validator lets
// through, parseDecimal must be able to read.
//
// Where they disagree the validator promises a value the read path then
// declines to protect, and the apply fails after the object has already been
// created. The containment is deliberately one-way - parseDecimal reads
// "1.000", which the column's two decimal places refuse - so only this
// direction is a defect.
func TestProductRevenueValidatorAcceptsOnlyParseableLiterals(t *testing.T) {
	t.Parallel()

	attr := resourceStringAttribute(t, "defectdojo_product", "revenue")
	checked := 0

	// A systematic sweep of sign/integer/fraction shapes rather than a list,
	// so a validator that admits some other degenerate spelling is caught too.
	for _, sign := range []string{"", "-", "+", "--", "-+"} {
		for _, intPart := range []string{"", "0", "5", "13", "0000000000000000"} {
			for _, frac := range []string{"", ".", ".0", ".55", ".555"} {
				literal := sign + intPart + frac
				if runStringValidators(t, "revenue", attr, literal) {
					continue
				}
				checked++
				if _, ok := parseDecimal(literal); !ok {
					t.Errorf("revenue = %q passes the schema validators, but parseDecimal "+
						"cannot read it. The read path therefore stores whatever the server "+
						"answered instead, and the apply fails with \"Provider produced "+
						"inconsistent result after apply\".", literal)
				}
			}
		}
	}

	if checked == 0 {
		t.Fatal("the validators rejected every generated literal; this test would pass vacuously")
	}
	t.Logf("checked %d accepted literals against parseDecimal", checked)
}

func TestPreserveDecimalLiteral(t *testing.T) {
	t.Parallel()

	runPreserveCases(t, func(current types.String, server string) types.String {
		return preserveLiteral(current, server, decimalsEqual)
	}, []preserveCase{
		// The issue this whole mechanism exists for.
		{"configured short form survives DRF canonicalisation", types.StringValue("100"), "100.00", "100"},
		{"already canonical is kept as-is", types.StringValue("100.00"), "100.00", "100.00"},
		// Real drift must still be reported, or the mechanism would hide
		// out-of-band changes.
		{"different amount takes the server value", types.StringValue("99"), "100.00", "100.00"},
		{"unparseable current takes the server value", types.StringValue("a lot of money"), "100.00", "100.00"},
		{"unparseable server value is passed through verbatim", types.StringValue("100"), "not a number", "not a number"},
	})
}
