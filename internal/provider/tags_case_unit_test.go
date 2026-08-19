package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func tagSet(vals ...string) types.Set {
	elems := make([]attr.Value, 0, len(vals))
	for _, v := range vals {
		elems = append(elems, types.StringValue(v))
	}
	return types.SetValueMust(types.StringType, elems)
}

// TestPreserveTagCase pins the helper that keeps a configured tag spelling when
// DefectDojo answers with a different case for the same tag. See tags.go for
// why the server does that.
func TestPreserveTagCase(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name    string
		current types.Set
		server  types.Set
		want    types.Set
	}{
		{
			// The bug: config says "foo", the instance already knows "Foo".
			name:    "case-only difference keeps the configured spelling",
			current: tagSet("foo", "bar"),
			server:  tagSet("Foo", "bar"),
			want:    tagSet("foo", "bar"),
		},
		{
			name:    "identical sets are unchanged",
			current: tagSet("foo", "bar"),
			server:  tagSet("foo", "bar"),
			want:    tagSet("foo", "bar"),
		},
		{
			name:    "set order does not matter",
			current: tagSet("foo", "bar"),
			server:  tagSet("BAR", "FOO"),
			want:    tagSet("foo", "bar"),
		},
		{
			// Real drift must still be reported.
			name:    "a genuinely different tag takes the server value",
			current: tagSet("foo", "bar"),
			server:  tagSet("foo", "baz"),
			want:    tagSet("foo", "baz"),
		},
		{
			name:    "a removed tag takes the server value",
			current: tagSet("foo", "bar"),
			server:  tagSet("foo"),
			want:    tagSet("foo"),
		},
		{
			name:    "an added tag takes the server value",
			current: tagSet("foo"),
			server:  tagSet("foo", "bar"),
			want:    tagSet("foo", "bar"),
		},
		{
			// State written before the case-collision validator existed can
			// still hold two spellings of one tag. Folding it down would make
			// it look like a match for a server set of the same folded size,
			// so the folded counts are compared as well as the raw ones. Here
			// the raw counts agree (2 and 2) and only the folded ones differ.
			name:    "a case-colliding current does not match a larger server set",
			current: tagSet("Foo", "foo"),
			server:  tagSet("Foo", "bar"),
			want:    tagSet("Foo", "bar"),
		},
		{
			// On import there is nothing to preserve.
			name:    "null current takes the server value",
			current: types.SetNull(types.StringType),
			server:  tagSet("Foo"),
			want:    tagSet("Foo"),
		},
		{
			name:    "unknown current takes the server value",
			current: types.SetUnknown(types.StringType),
			server:  tagSet("Foo"),
			want:    tagSet("Foo"),
		},
		{
			name:    "null server is passed through",
			current: tagSet("foo"),
			server:  types.SetNull(types.StringType),
			want:    types.SetNull(types.StringType),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := preserveTagCase(tc.current, tc.server); !got.Equal(tc.want) {
				t.Errorf("preserveTagCase(%v, %v) = %v, want %v", tc.current, tc.server, got, tc.want)
			}
		})
	}
}
