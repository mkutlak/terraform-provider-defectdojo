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
