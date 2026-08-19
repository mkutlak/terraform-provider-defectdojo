package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestNormalizeDecimal(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		in   string
		want string
	}{
		{"100", "100"},
		{"100.00", "100"}, // the DRF rendering of a configured "100"
		{"100.0", "100"},
		{"100.50", "100.5"},
		{"0100.50", "100.5"},
		{"+7", "7"},
		{"-1.50", "-1.5"},
		{"-1.05", "-1.05"},
		{"0", "0"},
		{"0.00", "0"},
		{"-0", "0"},
		{"-0.00", "0"},
		{"+0.000", "0"},
		{".5", "0.5"},
		{"5.", "5"},
		{"0000", "0"},
	} {
		got, ok := normalizeDecimal(tc.in)
		if !ok {
			t.Errorf("normalizeDecimal(%q) reported not-a-decimal, want %q", tc.in, tc.want)
			continue
		}
		if got != tc.want {
			t.Errorf("normalizeDecimal(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestNormalizeDecimalRejects(t *testing.T) {
	t.Parallel()

	for _, in := range []string{
		"",
		".",
		"+",
		"-",
		"1.2.3",
		"abc",
		"1e3",
		" 100",
		"100 ",
		"1,000.00",
		"a lot of money",
		"NaN",
	} {
		if got, ok := normalizeDecimal(in); ok {
			t.Errorf("normalizeDecimal(%q) = %q, true; want not-a-decimal", in, got)
		}
	}
}

// TestNormalizeDecimalIsExactBeyondFloat64 is an executable record of why
// normalizeDecimal compares strings rather than parsing to a float. DefectDojo's
// revenue column is DecimalField(max_digits=15, decimal_places=2), which needs
// 15 significant digits - past the point where float64 can tell neighbouring
// values apart.
func TestNormalizeDecimalIsExactBeyondFloat64(t *testing.T) {
	t.Parallel()

	a, aok := normalizeDecimal("1234567890123.45")
	b, bok := normalizeDecimal("1234567890123.46")
	if !aok || !bok {
		t.Fatalf("both inputs should normalise: %q(%v), %q(%v)", a, aok, b, bok)
	}
	if a == b {
		t.Errorf("normalizeDecimal collapsed two distinct 15-digit decimals onto %q", a)
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
	// the configured literal by preserveDecimalLiteral.
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

	// Not decimal literals at all. normalizeDecimal cannot key them, so
	// preserveDecimalLiteral could not protect them even where the server
	// takes the write.
	{"-", true},
	{".", true},
	{"-.", true},
	{"+", true},
	{"abc", true},
	{"1e3", true},
	{"1,000.00", true},
	{" 100", true},
	{"100 ", true},
	{"1.2.3", true},
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

// TestProductRevenueValidatorAcceptsOnlyNormalisableLiterals pins the property
// the schema and the read path have to share: every literal the validator lets
// through, normalizeDecimal must be able to key.
//
// Where they disagree the validator promises a value preserveDecimalLiteral
// then declines to protect, and the apply fails after the object has already
// been created. The containment is deliberately one-way - normalizeDecimal
// keys "1.000", which the column's two decimal places refuse - so only this
// direction is a defect.
//
// TestDecimalSpecPropertiesCarryDdFormat pins the analogous property for the
// ddFormat tag.
func TestProductRevenueValidatorAcceptsOnlyNormalisableLiterals(t *testing.T) {
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
				if _, ok := normalizeDecimal(literal); !ok {
					t.Errorf("revenue = %q passes the schema validators, but normalizeDecimal "+
						"reports it is not a decimal. preserveDecimalLiteral therefore stores "+
						"whatever the server answered instead, and the apply fails with "+
						"\"Provider produced inconsistent result after apply\".", literal)
				}
			}
		}
	}

	if checked == 0 {
		t.Fatal("the validators rejected every generated literal; this test would pass vacuously")
	}
	t.Logf("checked %d accepted literals against normalizeDecimal", checked)
}

func TestPreserveDecimalLiteral(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name    string
		current types.String
		server  string
		want    types.String
	}{
		{
			// The issue this whole mechanism exists for.
			name:    "configured short form survives DRF canonicalisation",
			current: types.StringValue("100"),
			server:  "100.00",
			want:    types.StringValue("100"),
		},
		{
			name:    "one decimal place survives",
			current: types.StringValue("1.5"),
			server:  "1.50",
			want:    types.StringValue("1.5"),
		},
		{
			name:    "already canonical is kept as-is",
			current: types.StringValue("100.00"),
			server:  "100.00",
			want:    types.StringValue("100.00"),
		},
		{
			name:    "negative short form survives",
			current: types.StringValue("-2.5"),
			server:  "-2.50",
			want:    types.StringValue("-2.5"),
		},
		{
			// Real drift must still be reported, or the mechanism would hide
			// out-of-band changes.
			name:    "different amount takes the server value",
			current: types.StringValue("99"),
			server:  "100.00",
			want:    types.StringValue("100.00"),
		},
		{
			name:    "null current takes the server value",
			current: types.StringNull(),
			server:  "100.00",
			want:    types.StringValue("100.00"),
		},
		{
			name:    "unknown current takes the server value",
			current: types.StringUnknown(),
			server:  "100.00",
			want:    types.StringValue("100.00"),
		},
		{
			name:    "unparseable current takes the server value",
			current: types.StringValue("a lot of money"),
			server:  "100.00",
			want:    types.StringValue("100.00"),
		},
		{
			name:    "unparseable server value is passed through verbatim",
			current: types.StringValue("100"),
			server:  "not a number",
			want:    types.StringValue("not a number"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := preserveDecimalLiteral(tc.current, tc.server); !got.Equal(tc.want) {
				t.Errorf("preserveDecimalLiteral(%v, %q) = %v, want %v", tc.current, tc.server, got, tc.want)
			}
		})
	}
}
