package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestPreserveHostCase pins the helper that keeps a configured host spelling
// when DefectDojo answers with the case-folded one. See host.go for why the
// server does that.
func TestPreserveHostCase(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name    string
		current types.String
		server  string
		want    types.String
	}{
		{
			// The bug: config says "API.Example.COM", DefectDojo stores
			// "api.example.com" and the apply fails on the difference.
			name:    "case-only difference keeps the configured spelling",
			current: types.StringValue("API.Example.COM"),
			server:  "api.example.com",
			want:    types.StringValue("API.Example.COM"),
		},
		{
			name:    "mixed case in one label only",
			current: types.StringValue("api.Example.com"),
			server:  "api.example.com",
			want:    types.StringValue("api.Example.com"),
		},
		{
			name:    "already lower case is kept as-is",
			current: types.StringValue("api.example.com"),
			server:  "api.example.com",
			want:    types.StringValue("api.example.com"),
		},
		{
			// clean_host() lower-cases the hex digits of an IPv6 literal too.
			name:    "upper-case IPv6 hex digits keep the configured spelling",
			current: types.StringValue("2001:DB8::1"),
			server:  "2001:db8::1",
			want:    types.StringValue("2001:DB8::1"),
		},
		{
			// Real drift must still be reported, or the mechanism would hide
			// out-of-band changes.
			name:    "a different host takes the server value",
			current: types.StringValue("API.Example.COM"),
			server:  "other.example.com",
			want:    types.StringValue("other.example.com"),
		},
		{
			// A near miss is still a miss: preservation is not a fuzzy match.
			name:    "a differing label takes the server value",
			current: types.StringValue("API.Example.COM"),
			server:  "api.example.net",
			want:    types.StringValue("api.example.net"),
		},
		{
			// clean_host() punycodes a non-ASCII name, which is a different
			// string rather than a different case, so state takes the server's.
			name:    "IDNA punycoding takes the server value",
			current: types.StringValue("Bücher.example"),
			server:  "xn--bcher-kva.example",
			want:    types.StringValue("xn--bcher-kva.example"),
		},
		{
			// Likewise IPv6 compression, which changes more than the case.
			name:    "IPv6 compression takes the server value",
			current: types.StringValue("2001:db8:0:0:0:0:0:1"),
			server:  "2001:db8::1",
			want:    types.StringValue("2001:db8::1"),
		},
		{
			// On import there is nothing configured to preserve.
			name:    "null current takes the server value",
			current: types.StringNull(),
			server:  "api.example.com",
			want:    types.StringValue("api.example.com"),
		},
		{
			name:    "unknown current takes the server value",
			current: types.StringUnknown(),
			server:  "api.example.com",
			want:    types.StringValue("api.example.com"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := preserveHostCase(tc.current, tc.server); !got.Equal(tc.want) {
				t.Errorf("preserveHostCase(%v, %q) = %v, want %v", tc.current, tc.server, got, tc.want)
			}
		})
	}
}
