package provider

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DefectDojo normalises tags server-side: a tag is lower-cased and rejected if
// it contains characters outside the grammar below. A config carrying "Foo"
// therefore comes back as "foo", which is the issue #23 shape - the write path
// accepts a form the read path cannot reproduce, and the apply dies with
// "Provider produced inconsistent result after apply" only after the server has
// been mutated. Rejecting it during plan turns that into a pointed config error
// before anything is written.
//
// The pattern and message are the ones product_resource.go has carried (with a
// negative acceptance test) since commit a54fb89, "Add validation to product
// fields to match API". engagement, test and url declared the identical
// attribute with no validation at all; this helper exists so the four cannot
// drift apart again.
var tagPattern = regexp.MustCompile(`\A[a-z0-9][a-z0-9_-]*\z`)

const tagPatternMessage = "Tags must be lower case values (letters, digits, hyphens, underscores)"

// tagsSetAttribute builds the `tags` attribute shared by every tag-bearing
// resource. Callers pass their own description so the generated docs keep the
// wording each resource already used.
func tagsSetAttribute(markdownDescription string) schema.SetAttribute {
	return schema.SetAttribute{
		MarkdownDescription: markdownDescription,
		Optional:            true,
		ElementType:         types.StringType,
		Validators: []validator.Set{
			setvalidator.ValueStringsAre(
				stringvalidator.RegexMatches(tagPattern, tagPatternMessage),
			),
		},
	}
}
