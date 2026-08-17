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

// tagValidationCases are the inputs every tag-bearing resource must agree on.
// The rejected ones are all values DefectDojo either lower-cases or refuses, so
// accepting them in the plan only defers the failure to apply time.
var tagValidationCases = []struct {
	tag       string
	wantError bool
}{
	{"ok", false},
	{"ok-tag_1", false},
	{"0-starts-with-a-digit", false},
	{"Foo", true},           // server lower-cases
	{"UPPER", true},         // server lower-cases
	{"team.security", true}, // dots are outside the grammar
	{"needs review", true},  // spaces are outside the grammar
	{"-leading-hyphen", true},
	{"_leading-underscore", true},
	{"", true},
}

// runTagSetValidators applies a Set attribute's validators to one tag value and
// reports whether they produced an error.
func runTagSetValidators(t *testing.T, tagsAttr schema.SetAttribute, tag string) bool {
	t.Helper()

	ctx := context.Background()
	value := types.SetValueMust(types.StringType, []attr.Value{types.StringValue(tag)})

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
