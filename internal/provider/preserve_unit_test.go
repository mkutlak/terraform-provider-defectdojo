package provider

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// preserveCase is one row of a preserveLiteral table: the value already in
// state, the value the server sent, and the value that should reach state.
type preserveCase struct {
	name    string
	current types.String
	server  string
	want    string
}

// runPreserveCases drives a table against one of the preserveLiteral wrappers.
// The three of them - decimal, host and datetime - share a shape, so they share
// this runner rather than each restating it.
func runPreserveCases(t *testing.T, fn func(types.String, string) types.String, cases []preserveCase) {
	t.Helper()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := fn(tc.current, tc.server); got.ValueString() != tc.want {
				t.Errorf("current=%v, server=%q -> %q, want %q", tc.current, tc.server, got.ValueString(), tc.want)
			}
		})
	}
}

// TestPreserveLiteral covers the parts of preserveLiteral that do not depend on
// which equivalence relation is passed in. The per-format tests then only have
// to exercise their own relation.
func TestPreserveLiteral(t *testing.T) {
	t.Parallel()

	runPreserveCases(t, func(current types.String, server string) types.String {
		return preserveLiteral(current, server, strings.EqualFold)
	}, []preserveCase{
		{"equivalent current is kept verbatim", types.StringValue("VALUE"), "value", "VALUE"},
		{"non-equivalent current takes the server value", types.StringValue("other"), "value", "value"},
		{"null current takes the server value", types.StringNull(), "value", "value"},
		{"unknown current takes the server value", types.StringUnknown(), "value", "value"},
	})
}
