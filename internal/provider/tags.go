package provider

import (
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DefectDojo accepts almost anything as a tag. Verified against 3.1.101 by
// creating a product per input and reading it back:
//
//	"env:prod", "v1.2.3", "team/security", "tag#1", "_internal", "-legacy",
//	"café", "Foo", "UPPER"      -> stored and returned verbatim
//	"a b"                       -> 400 "Invalid tag: 'a b'. Tags should not
//	                               contain spaces, commas, or quotes."
//	"a'b"                       -> 400, same message
//	"a,b"  and  "a\"b"          -> 201, but silently SPLIT into ["a", "b"]
//
// Only the last two rows are the issue #23 shape: the server accepts the write
// and answers with something the practitioner did not ask for, so state
// disagrees with config and the apply fails with "Provider produced
// inconsistent result after apply". The 400s are worth catching at plan time
// too, since the message is identical and the round trip never happens.
//
// So the grammar below rejects exactly the four characters DefectDojo's own
// error message names, and nothing else.
//
// It deliberately does NOT reject uppercase. product_resource.go carried a
// `\A[a-z0-9][a-z0-9_-]*\z` pattern with the message "Tags must be lower case
// values" from commit a54fb89 (2022), but that is not true of DefectDojo
// 3.1.101: "Foo" round-trips as "Foo". Enforcing it rejected working
// configurations - `env:prod`, `owasp:a01`, `v1.2.3` and `team/security` are
// ordinary tag conventions that the server stores untouched.
var tagPattern = regexp.MustCompile(`\A[^ ,'"]+\z`)

const tagPatternMessage = `Tags must not be empty or contain spaces, commas, single quotes or double quotes. ` +
	`DefectDojo rejects spaces and single quotes outright, and silently splits a tag containing ` +
	`a comma or a double quote into several tags, which would make state disagree with configuration.`

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

// preserveTagCase keeps the configured spelling of a tag set when the server's
// answer differs only by letter case.
//
// DefectDojo stores tags in a single instance-wide table, matched
// case-insensitively, and answers with whichever spelling was registered FIRST
// anywhere on that instance. Verified on 3.1.101: with "Foo" already present,
// creating a product with tags ["foo","bar"] returns ["Foo","bar"]; with
// "ZZTOP" present, submitting "zztop" returns "ZZTOP".
//
// That makes the returned spelling a property of global server state rather
// than of the practitioner's configuration, so no schema validator can prevent
// the mismatch - a perfectly lower-case config still breaks the moment someone
// registers the capitalised spelling through the UI or another tool. Without
// this, Terraform fails the apply with:
//
//	.tags: planned set element cty.StringVal("foo") does not correlate with
//	any element in actual
//
// Since the two spellings denote the SAME tag server-side, the configured
// literal is kept. A genuine difference in the set of tags is never preserved,
// so real drift is still reported.
func preserveTagCase(current types.Set, server types.Set) types.Set {
	if current.IsNull() || current.IsUnknown() || server.IsNull() || server.IsUnknown() {
		return server
	}

	currentElems := current.Elements()
	serverElems := server.Elements()
	if len(currentElems) != len(serverElems) {
		return server
	}

	fold := func(elems []attr.Value) (map[string]bool, bool) {
		out := make(map[string]bool, len(elems))
		for _, e := range elems {
			s, ok := e.(types.String)
			if !ok || s.IsNull() || s.IsUnknown() {
				return nil, false
			}
			out[strings.ToLower(s.ValueString())] = true
		}
		return out, true
	}

	currentFolded, ok := fold(currentElems)
	if !ok {
		return server
	}
	serverFolded, ok := fold(serverElems)
	if !ok {
		return server
	}
	if len(currentFolded) != len(serverFolded) {
		return server
	}
	for k := range currentFolded {
		if !serverFolded[k] {
			return server
		}
	}

	return current
}
