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

// tagValidationCases mirrors what DefectDojo 3.1.101 actually does. Each
// accepted value below was created through the API and read back verbatim; each
// rejected value either returns 400 or is silently split into several tags.
// See the table in tags.go.
var tagValidationCases = []struct {
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

// runTagSetValidatorsOnSet applies a Set attribute's validators to a whole tag
// set and reports whether they produced an error.
func runTagSetValidatorsOnSet(t *testing.T, tagsAttr schema.SetAttribute, tags ...string) bool {
	t.Helper()

	ctx := context.Background()
	elems := make([]attr.Value, 0, len(tags))
	for _, tag := range tags {
		elems = append(elems, types.StringValue(tag))
	}
	value := types.SetValueMust(types.StringType, elems)

	req := validator.SetRequest{
		Path:           path.Root("tags"),
		PathExpression: path.MatchRoot("tags"),
		ConfigValue:    value,
	}

	for _, v := range tagsAttr.Validators {
		resp := &validator.SetResponse{}
		v.ValidateSet(ctx, req, resp)
		if resp.Diagnostics.HasError() {
			return true
		}
	}
	return false
}

// runTagSetValidators applies a Set attribute's validators to one tag value and
// reports whether they produced an error.
func runTagSetValidators(t *testing.T, tagsAttr schema.SetAttribute, tag string) bool {
	t.Helper()

	return runTagSetValidatorsOnSet(t, tagsAttr, tag)
}

// TestTagsAttributesRejectNonCanonicalValues asserts that every resource
// exposing a "tags" attribute enforces DefectDojo's tag grammar.
//
// Three of the four tag-bearing resources (engagement, test, url) shipped with
// no validator at all, so `tags = ["Foo"]` was accepted at plan time and then
// failed at apply - after the object had already been created - because the
// server answers with the lower-cased form. product has had the validator since
// commit a54fb89; this test is what keeps the other three from drifting again.
//
// It checks behaviour rather than comparing validator slices: a
// stringvalidator.RegexMatches holds a *regexp.Regexp, which does not compare
// meaningfully with reflect.DeepEqual.
func TestTagsAttributesRejectNonCanonicalValues(t *testing.T) {
	t.Parallel()

	checked := 0
	for tfTypeName, resp := range providerResourceSchemas(t) {
		attr, ok := resp.Schema.Attributes["tags"]
		if !ok {
			continue
		}
		setAttr, ok := attr.(schema.SetAttribute)
		if !ok {
			t.Errorf("resource %s: \"tags\" is %T, expected schema.SetAttribute", tfTypeName, attr)
			continue
		}
		checked++

		for _, tc := range tagValidationCases {
			got := runTagSetValidators(t, setAttr, tc.tag)
			if got != tc.wantError {
				verb := "accepted"
				if got {
					verb = "rejected"
				}
				t.Errorf("resource %s: tags = [%q] was %s, but wantError=%v. "+
					"Every tag-bearing resource must use tagsSetAttribute() from tags.go.",
					tfTypeName, tc.tag, verb, tc.wantError)
			}
		}
	}

	if checked < 4 {
		t.Errorf("only %d resources expose a \"tags\" attribute; expected at least 4 "+
			"(product, engagement, test, url). Did a resource stop declaring tags, or is "+
			"the schema walk broken?", checked)
	}
}

// tagCollisionCases are whole tag SETS rather than single tags: DefectDojo keys
// its tag tables on a lower-cased slug, so two elements of one set that differ
// only in case are a single tag server-side and the set can never round-trip.
var tagCollisionCases = []struct {
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

// TestTagsAttributesRejectCaseCollidingElements asserts that every tag-bearing
// resource rejects a `tags` set holding two spellings of the same tag.
//
// DefectDojo matches tag names case-insensitively, so `tags = ["Foo","foo"]` is
// ONE tag once it reaches the server. The write is accepted, the create
// response echoes a single spelling twice, and the follow-up read returns one
// element - so the configuration can never be satisfied and the resource plans
// forever:
//
//	~ tags = [
//	      "TfBackend",
//	    + "tfbackend",
//	  ]
//
// The grammar in tags.go cannot catch this, because each element is a perfectly
// legal tag on its own; only the set as a whole is contradictory. Hence a
// set-level validator, wired through tagsSetAttribute so all four resources get
// it and cannot drift apart.
func TestTagsAttributesRejectCaseCollidingElements(t *testing.T) {
	t.Parallel()

	checked := 0
	for tfTypeName, resp := range providerResourceSchemas(t) {
		attr, ok := resp.Schema.Attributes["tags"]
		if !ok {
			continue
		}
		setAttr, ok := attr.(schema.SetAttribute)
		if !ok {
			t.Errorf("resource %s: \"tags\" is %T, expected schema.SetAttribute", tfTypeName, attr)
			continue
		}
		checked++

		for _, tc := range tagCollisionCases {
			got := runTagSetValidatorsOnSet(t, setAttr, tc.tags...)
			if got != tc.wantError {
				verb := "accepted"
				if got {
					verb = "rejected"
				}
				t.Errorf("resource %s: %s: tags = %q was %s, but wantError=%v. "+
					"Every tag-bearing resource must use tagsSetAttribute() from tags.go.",
					tfTypeName, tc.name, tc.tags, verb, tc.wantError)
			}
		}
	}

	if checked < 4 {
		t.Errorf("only %d resources expose a \"tags\" attribute; expected at least 4 "+
			"(product, engagement, test, url).", checked)
	}
}
