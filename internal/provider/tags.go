package provider

import (
	"context"
	"fmt"
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

// tagCaseCollisionValidator rejects a `tags` set that holds two spellings of
// what DefectDojo stores as a single tag.
//
// Every tag table keys its rows on a lower-cased slug, so names are matched
// case-insensitively and `tags = ["Foo","foo"]` is one tag by the time it
// reaches the database. Verified on 3.1.101:
//
//	POST /api/v2/products/ {"tags":["Backend","backend"]}
//	  -> 201, create response ['Backend', 'Backend']   (one spelling, twice)
//	  -> GET returns ['Backend']                       (one tag)
//
// The apply itself succeeds, which is what makes this worth catching early: the
// configuration asks for two elements, the server only ever reports one, and
// every later plan proposes adding the missing spelling again, forever.
//
// This cannot live in tagPattern. Each element is a perfectly legal tag on its
// own - "Foo" and "foo" both round-trip when they are the only spelling in
// play - so only the set as a whole is contradictory, and only a set-level
// validator can see that.
type tagCaseCollisionValidator struct{}

func (v tagCaseCollisionValidator) Description(_ context.Context) string {
	return "tag names must not differ only by letter case, because DefectDojo matches them case-insensitively"
}

func (v tagCaseCollisionValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v tagCaseCollisionValidator) ValidateSet(_ context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	elems := req.ConfigValue.Elements()
	firstSpelling := make(map[string]string, len(elems))
	for _, e := range elems {
		s, ok := e.(types.String)
		if !ok || s.IsNull() || s.IsUnknown() {
			continue
		}
		tag := s.ValueString()
		folded := strings.ToLower(tag)
		if first, collides := firstSpelling[folded]; collides {
			resp.Diagnostics.AddAttributeError(req.Path, "Case-Colliding Tags", fmt.Sprintf(
				"tags contains both %q and %q, which DefectDojo stores as a single tag: tag "+
					"names are keyed on a lower-cased slug, so they are matched "+
					"case-insensitively. Writing both stores one tag and the read returns one "+
					"element, so this configuration can never be satisfied - the apply succeeds "+
					"and every later plan proposes adding the missing spelling again. Keep one "+
					"spelling of the tag.",
				first, tag))
			return
		}
		firstSpelling[folded] = tag
	}
}

// dedupeTagElements drops repeated elements from a tag list the server sent,
// keeping the first occurrence of each.
//
// A Terraform set cannot hold a duplicate, and DefectDojo can answer with one:
// a create under product tag inheritance echoes the child's own tag twice, and
// a case-colliding write echoes one spelling twice. types.SetValue does not
// police duplicates, so the value used to reach state intact and the framework
// rejected it afterwards with a diagnostic that named neither DefectDojo nor
// the cause.
//
// The comparison is exact rather than case-folded. Two spellings of one tag
// never come back from a single object - the tag table resolves them to one row
// before the response is rendered - so folding here would only ever discard a
// genuine second tag, and it would have to invent a rule for which spelling
// survives.
func dedupeTagElements(elems []attr.Value) []attr.Value {
	out := make([]attr.Value, 0, len(elems))
	seen := make(map[string]bool, len(elems))
	for _, e := range elems {
		s, ok := e.(types.String)
		if !ok || s.IsNull() || s.IsUnknown() {
			out = append(out, e)
			continue
		}
		if seen[s.ValueString()] {
			continue
		}
		seen[s.ValueString()] = true
		out = append(out, e)
	}
	return out
}

// tagsSharedDescription documents what this file wires in, as opposed to what
// any one resource is for. It lives beside the code that enforces it, and is
// appended to all four descriptions from here, because four copies of it in
// four schemas could drift apart from each other and from the behaviour - the
// same reason the validators and the case preservation are shared rather than
// repeated per resource.
const tagsSharedDescription = "Tags must not contain spaces, commas or quotes, and the configured " +
	"spelling is kept when DefectDojo answers with a different letter case."

// tagsSetAttribute builds the `tags` attribute shared by every tag-bearing
// resource. Callers pass their own description so the generated docs keep the
// wording each resource already used; it is a bare phrase, and the sentence
// break before tagsSharedDescription is added here.
func tagsSetAttribute(markdownDescription string) schema.SetAttribute {
	return schema.SetAttribute{
		MarkdownDescription: markdownDescription + ". " + tagsSharedDescription,
		Optional:            true,
		ElementType:         types.StringType,
		Validators: []validator.Set{
			setvalidator.ValueStringsAre(
				stringvalidator.RegexMatches(tagPattern, tagPatternMessage),
			),
			tagCaseCollisionValidator{},
		},
	}
}

// preserveTagCase keeps the configured spelling of a tag set when the server's
// answer differs only by letter case.
//
// DefectDojo gives each object type its own tag table - Tagulous_Product_tags,
// Tagulous_Engagement_tags, Tagulous_Test_tags and so on - and keys every row
// on a lower-cased slug beside the name it was created with. Names are
// therefore matched case-insensitively, and the server answers with whichever
// spelling was registered FIRST in that type's table, by any object of that
// type on the instance. Verified on 3.1.101: with "Foo" already present,
// creating a product with tags ["foo","bar"] returns ["Foo","bar"]; with
// "ZZTOP" present, submitting "zztop" returns "ZZTOP".
//
// The scope is the object type, not the instance: a product storing
// "MixedCase" leaves an engagement free to store "mixedcase", because the two
// rows live in different tables. Every TagField is declared force_lowercase, but
// that only reaches the slug - the name keeps the submitted spelling, which is
// why uppercase round-trips at all.
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
	// The raw counts already agree, so this only differs from that check when
	// one side collapses under folding - which means `current` holds two
	// spellings of one tag. tagCaseCollisionValidator keeps that out of new
	// configurations, but state written before it existed can still carry it,
	// and folding such a set down would make it look like a match for a server
	// answer it does not describe.
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
