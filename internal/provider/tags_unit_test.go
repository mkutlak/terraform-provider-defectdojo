package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// wantTagBearingResources is the number of resources that must declare a `tags`
// attribute: product, engagement, test and url.
const wantTagBearingResources = 4

func tagSet(vals ...string) types.Set {
	elems := make([]attr.Value, 0, len(vals))
	for _, v := range vals {
		elems = append(elems, types.StringValue(v))
	}
	return types.SetValueMust(types.StringType, elems)
}

// runTagSetValidators applies a Set attribute's validators to a whole tag set
// and reports whether they produced an error.
func runTagSetValidators(t *testing.T, tagsAttr schema.SetAttribute, tags ...string) bool {
	t.Helper()

	elems := make([]attr.Value, 0, len(tags))
	for _, tag := range tags {
		elems = append(elems, types.StringValue(tag))
	}

	req := validator.SetRequest{
		Path:           path.Root("tags"),
		PathExpression: path.MatchRoot("tags"),
		ConfigValue:    types.SetValueMust(types.StringType, elems),
	}

	for _, v := range tagsAttr.Validators {
		resp := &validator.SetResponse{}
		v.ValidateSet(context.Background(), req, resp)
		if resp.Diagnostics.HasError() {
			return true
		}
	}
	return false
}

// forEachTagsAttribute runs fn against the `tags` attribute of every resource
// that declares one.
//
// It checks behaviour rather than comparing validator slices: a
// stringvalidator.RegexMatches holds a *regexp.Regexp, which does not compare
// meaningfully with reflect.DeepEqual.
func forEachTagsAttribute(t *testing.T, fn func(tfTypeName string, tagsAttr schema.SetAttribute)) {
	t.Helper()

	checked := 0
	for tfTypeName, resp := range providerResourceSchemas(t) {
		attribute, ok := resp.Schema.Attributes["tags"]
		if !ok {
			continue
		}
		setAttr, ok := attribute.(schema.SetAttribute)
		if !ok {
			t.Errorf("resource %s: \"tags\" is %T, expected schema.SetAttribute", tfTypeName, attribute)
			continue
		}
		checked++
		fn(tfTypeName, setAttr)
	}

	if checked < wantTagBearingResources {
		t.Errorf("only %d resources expose a \"tags\" attribute; expected at least %d "+
			"(product, engagement, test, url). Did a resource stop declaring tags, or is "+
			"the schema walk broken?", checked, wantTagBearingResources)
	}
}

// TestTagsAttributesRejectNonCanonicalValues asserts that every resource
// exposing a `tags` attribute enforces DefectDojo's tag grammar.
//
// Three of the four tag-bearing resources (engagement, test, url) shipped with
// no validator at all, so `tags = ["a,b"]` was accepted at plan time and then
// failed at apply - after the object had already been created - because
// DefectDojo silently splits that into two tags. Letter case is not the problem:
// force_lowercase reaches only the slug, so the name keeps the spelling it was
// submitted with and "Foo" round-trips unchanged.
//
// Each accepted value below was created through the API on 3.1.101 and read back
// verbatim; each rejected value either returns 400 or is silently split. See the
// table in tags.go.
func TestTagsAttributesRejectNonCanonicalValues(t *testing.T) {
	t.Parallel()

	cases := []struct {
		tag       string
		wantError bool
	}{
		// Stored verbatim - the validator must not stand in the way of these.
		{"ok", false},
		{"ok-tag_1", false},
		{"0-starts-with-a-digit", false},
		{"env:prod", false},
		{"owasp:a01", false},
		{"v1.2.3", false},
		{"team/security", false},
		{"tag#1", false},
		{"tag+1", false},
		{"_internal", false},
		{"-legacy", false},
		{"café", false},
		// Uppercase round-trips unchanged on 3.1.101; rejecting it broke working
		// configurations. See the note in tags.go.
		{"Foo", false},
		{"MixedCase", false},

		// 400 "Invalid tag: ... Tags should not contain spaces, commas, or quotes."
		{"needs review", true},
		{"a'b", true},
		// Accepted with a 201, but silently split into two tags - the issue #23
		// shape, and the reason these are rejected at plan time.
		{"a,b", true},
		{`a"b`, true},
		// The serializer drops an empty tag, so it can never round-trip.
		{"", true},
	}

	forEachTagsAttribute(t, func(tfTypeName string, tagsAttr schema.SetAttribute) {
		for _, tc := range cases {
			got := runTagSetValidators(t, tagsAttr, tc.tag)
			if got == tc.wantError {
				continue
			}
			verb := "accepted"
			if got {
				verb = "rejected"
			}
			t.Errorf("resource %s: tags = [%q] was %s, but wantError=%v. "+
				"Every tag-bearing resource must use tagsSetAttribute() from tags.go.",
				tfTypeName, tc.tag, verb, tc.wantError)
		}
	})
}

// TestTagsAttributesRejectCaseCollidingElements asserts that every tag-bearing
// resource rejects a `tags` set holding two spellings of the same tag.
//
// DefectDojo matches tag names case-insensitively, so `tags = ["Foo","foo"]` is
// ONE tag once it reaches the server. The write is accepted, the create response
// echoes a single spelling twice, and the follow-up read returns one element - so
// the configuration can never be satisfied and the resource plans forever:
//
//	~ tags = [
//	      "TfBackend",
//	    + "tfbackend",
//	  ]
func TestTagsAttributesRejectCaseCollidingElements(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		tags      []string
		wantError bool
	}{
		// Nothing here folds together, so none of these may be rejected.
		{"an empty set", nil, false},
		{"a single tag", []string{"Foo"}, false},
		{"distinct tags", []string{"foo", "bar"}, false},
		{"distinct tags in different cases", []string{"Foo", "bar", "BAZ"}, false},
		{"a shared prefix is not a collision", []string{"foo", "foobar"}, false},

		// Each of these is one tag server-side.
		{"two spellings of the same tag", []string{"Foo", "foo"}, true},
		{"a collision the caller did not put next to each other", []string{"a", "Backend", "b", "backend"}, true},
		{"three spellings of the same tag", []string{"Foo", "FOO", "foo"}, true},
		{"a collision that differs in an inner letter only", []string{"envProd", "envprod"}, true},
	}

	forEachTagsAttribute(t, func(tfTypeName string, tagsAttr schema.SetAttribute) {
		for _, tc := range cases {
			got := runTagSetValidators(t, tagsAttr, tc.tags...)
			if got == tc.wantError {
				continue
			}
			verb := "accepted"
			if got {
				verb = "rejected"
			}
			t.Errorf("resource %s: %s: tags = %q was %s, but wantError=%v. "+
				"Every tag-bearing resource must use tagsSetAttribute() from tags.go.",
				tfTypeName, tc.name, tc.tags, verb, tc.wantError)
		}
	})
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
			name:    "a differing element count takes the server value",
			current: tagSet("foo", "bar"),
			server:  tagSet("foo"),
			want:    tagSet("foo"),
		},
		{
			// State written before the case-collision validator existed can
			// still hold two spellings of one tag. Folding it down would make it
			// look like a match for a server set of the same folded size, so the
			// folded counts are compared as well as the raw ones. Here the raw
			// counts agree (2 and 2) and only the folded ones differ.
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
			t.Parallel()

			if got := preserveTagCase(tc.current, tc.server); !got.Equal(tc.want) {
				t.Errorf("preserveTagCase(%v, %v) = %v, want %v", tc.current, tc.server, got, tc.want)
			}
		})
	}
}
