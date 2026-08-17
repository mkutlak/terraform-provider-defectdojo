package provider

// Round-trip guardrails. These sit at the intersection of the two audits that
// already exist: schema_validation_unit_test.go owns schema SHAPE, and
// ddfield_audit_unit_test.go owns the ddField TAG CONTRACT. Neither can see the
// bug class that produced issue #23, because that bug lives in the join between
// them - a schema flag that is only correct if the ddclient field behind it can
// represent the same set of values.
//
// The invariant both tests defend:
//
//	An attribute that is Optional, not Computed, and has no Default plans as
//	NULL when the practitioner omits it. If the value that comes back on the
//	read path can never be null, Terraform core rejects the apply with
//	"Provider produced inconsistent result after apply: .attr: was null, but
//	now ...". The practitioner cannot work around it, and the server has
//	already been mutated by the time it fires.
//
// TestResourceOptionalAttributesCanBeNull checks that against the Go type
// (always runnable). TestResourceOptionalAttributesAreNullableInSpec checks it
// against DefectDojo's own OpenAPI document (runnable whenever a spec has been
// collected). They are complementary, not redundant: the first catches a
// non-pointer target, the second catches a pointer target that the server
// nevertheless always populates - which is what every real defect in this class
// has turned out to be.

import (
	"cmp"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// ---------------------------------------------------------------------------
// Joining a resource schema to the ddclient struct behind it
// ---------------------------------------------------------------------------

// ddTarget describes where a single Terraform attribute lands on the generated
// ddclient struct.
type ddTarget struct {
	tfsdkName  string
	tfType     reflect.Type // types.String, types.Int64, ...
	ddFieldTag string       // the ddField tag value
	ddType     reflect.Type // the resolved ddclient field type
	jsonName   string       // the json tag on that field, first comma-separated token
	ddFormat   string       // the ddFormat tag value, "" when absent
}

// resourceDdTargets joins a registered resource's Terraform model to the
// ddclient struct its defectdojoResource() wraps, keyed by tfsdk attribute name.
//
// The join goes through testAccResourceFactories (acctest_registry_test.go)
// rather than ddFieldAuditTable: the factory map is keyed by the same Terraform
// type name as providerResourceSchemas and is derived from the provider's own
// registry, so it cannot go stale. ddFieldAuditTable is keyed by Go type name
// and also lists data-source-only models that have no resource schema at all.
func resourceDdTargets(t *testing.T, tfTypeName string) (map[string]ddTarget, reflect.Type, bool) {
	t.Helper()

	model, ok := testAccResourceModel(tfTypeName)
	if !ok {
		return nil, nil, false // reported by TestAccResourceRegistryIsDerivable
	}

	modelType := reflect.TypeOf(model)
	if modelType.Kind() != reflect.Ptr || modelType.Elem().Kind() != reflect.Struct {
		t.Errorf("resource %s: model must be a pointer to a struct, got %s", tfTypeName, modelType)
		return nil, nil, false
	}
	ddStruct := ddStructType(t, tfTypeName, model.defectdojoResource())

	out := map[string]ddTarget{}
	for field := range modelType.Elem().Fields() {
		ddFieldTag := field.Tag.Get("ddField")
		if ddFieldTag == "" {
			continue // the engine never writes this attribute
		}
		ddField, ok := ddStruct.FieldByName(ddFieldTag)
		if !ok {
			continue // reported by TestDdFieldAuditTagsResolve
		}
		// Promoted fields keep the tags they were declared with, so this
		// reaches through the embedded dd.* struct to the generated json tag.
		jsonName, _, _ := strings.Cut(ddField.Tag.Get("json"), ",")

		out[field.Tag.Get("tfsdk")] = ddTarget{
			tfsdkName:  field.Tag.Get("tfsdk"),
			tfType:     field.Type,
			ddFieldTag: ddFieldTag,
			ddType:     ddField.Type,
			jsonName:   jsonName,
			ddFormat:   field.Tag.Get("ddFormat"),
		}
	}
	return out, ddStruct, true
}

// schemaAttributeHasDefault reports whether a schema attribute declares a
// Default. The framework's own AttributeWithXDefaultValue interfaces live under
// internal/, but the accessor methods and their return types are public, so
// these local interfaces reach them the same way schemaAttribute does in
// schema_validation_unit_test.go.
//
// An attribute kind that is not listed here reads as "no default", which makes
// the guardrails STRICTER rather than weaker - the safe direction to fail if a
// future schema starts using a kind this list has not caught up with.
func schemaAttributeHasDefault(attr any) bool {
	switch a := attr.(type) {
	case interface{ StringDefaultValue() defaults.String }:
		return a.StringDefaultValue() != nil
	case interface{ BoolDefaultValue() defaults.Bool }:
		return a.BoolDefaultValue() != nil
	case interface{ Int64DefaultValue() defaults.Int64 }:
		return a.Int64DefaultValue() != nil
	case interface{ Int32DefaultValue() defaults.Int32 }:
		return a.Int32DefaultValue() != nil
	case interface{ Float64DefaultValue() defaults.Float64 }:
		return a.Float64DefaultValue() != nil
	case interface{ Float32DefaultValue() defaults.Float32 }:
		return a.Float32DefaultValue() != nil
	case interface{ SetDefaultValue() defaults.Set }:
		return a.SetDefaultValue() != nil
	case interface{ ListDefaultValue() defaults.List }:
		return a.ListDefaultValue() != nil
	case interface{ MapDefaultValue() defaults.Map }:
		return a.MapDefaultValue() != nil
	case interface{ NumberDefaultValue() defaults.Number }:
		return a.NumberDefaultValue() != nil
	case interface{ ObjectDefaultValue() defaults.Object }:
		return a.ObjectDefaultValue() != nil
	}
	return false
}

// planningAsNull reports whether omitting this attribute from config makes
// Terraform plan a null value for it - the precondition for the whole bug
// class. Optional+Computed plans as unknown (the server may fill it) and a
// Default plans as the default, so neither can produce the inconsistency.
func planningAsNull(attr any) bool {
	a, ok := attr.(schemaAttribute)
	return ok && a.IsOptional() && !a.IsComputed() && !schemaAttributeHasDefault(attr)
}

// isSliceTarget reports whether a ddclient field is a collection (or a pointer
// to one).
func isSliceTarget(ddType reflect.Type) bool {
	if ddType.Kind() == reflect.Ptr {
		ddType = ddType.Elem()
	}
	return ddType.Kind() == reflect.Slice
}

// ---------------------------------------------------------------------------
// Guardrail A: the ddclient field must be able to represent null
// ---------------------------------------------------------------------------

// ddTargetCanBeNull reports whether populateResourceData can ever write a null
// Terraform value for a ddclient field of the given type.
//
// Pointers get a nil check in every read branch. time.Time and
// openapi_types.Date get an IsZero() check (resource.go:640, :658). Slices and
// pointers-to-slices get the empty-collection branch (resource.go:739, :751,
// :768, :783). Everything else - a bare string/bool/int/int32/int64/float32/
// float64 - is read unconditionally, so the engine ALWAYS writes a known value
// and an Optional-only attribute mapped to one can never survive an apply that
// omits it.
func ddTargetCanBeNull(ddType reflect.Type) bool {
	switch {
	case ddType.Kind() == reflect.Ptr,
		ddType.Kind() == reflect.Slice,
		ddType == reflect.TypeFor[time.Time](),
		ddType == reflect.TypeFor[openapi_types.Date]():
		return true
	}
	return false
}

// TestResourceOptionalAttributesCanBeNull fails when an attribute that plans as
// null is backed by a ddclient field the engine always reads a concrete value
// out of. This is the issue #23 bug class reduced to a property of the Go
// types, so it needs no DefectDojo instance and no collected spec.
//
// It passes today. Its value is that it cannot stop passing quietly: a
// `make regen-client` that turns a *string into a string, or a new resource
// written against a non-pointer field, becomes a plain `go test` failure.
func TestResourceOptionalAttributesCanBeNull(t *testing.T) {
	t.Parallel()

	checked := 0
	for tfTypeName, resp := range providerResourceSchemas(t) {
		if resp.Diagnostics.HasError() {
			continue // reported by TestSchemasValidateImplementation
		}
		targets, ddStruct, ok := resourceDdTargets(t, tfTypeName)
		if !ok {
			continue
		}

		for attrName, attr := range resp.Schema.Attributes {
			if !planningAsNull(attr) {
				continue
			}
			target, ok := targets[attrName]
			if !ok {
				continue // no ddField tag: the engine never writes it
			}
			checked++
			if ddTargetCanBeNull(target.ddType) {
				continue
			}
			t.Errorf("resource %s: attribute %q is Optional without Computed and without a "+
				"Default, but its ddField target %s.%s has type %s, for which populateResourceData "+
				"always writes a known value. Every apply that omits %q will fail with "+
				"\"Provider produced inconsistent result after apply: .%s: was null, but now ...\".\n"+
				"Fix by one of:\n"+
				"  - Computed: true                     (DefectDojo picks the value)\n"+
				"  - Computed: true plus a Default:     (the provider picks the value)\n"+
				"  - make the ddclient field a pointer  (only if the API can really omit it)",
				tfTypeName, attrName, ddStruct, target.ddFieldTag, target.ddType, attrName, attrName)
		}
	}

	if checked == 0 {
		t.Fatal("no Optional-without-Computed attributes were checked; the join is broken and " +
			"this test would pass vacuously")
	}
	t.Logf("checked %d Optional-without-Computed attributes against their ddclient field types", checked)
}

// TestDdTargetCanBeNull pins the predicate the guardrail above depends on.
//
// It exists because the repo currently has zero non-pointer Optional targets,
// so the usual sanity check - break the schema on purpose, watch the guardrail
// go red - cannot be performed against real code. This table is the executable
// proof that ddTargetCanBeNull means what its doc comment claims.
func TestDdTargetCanBeNull(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		typ  reflect.Type
		want bool
	}{
		{reflect.TypeFor[string](), false},
		{reflect.TypeFor[bool](), false},
		{reflect.TypeFor[int](), false},
		{reflect.TypeFor[int32](), false},
		{reflect.TypeFor[int64](), false},
		{reflect.TypeFor[float32](), false},
		{reflect.TypeFor[float64](), false},
		{reflect.TypeFor[*string](), true},
		{reflect.TypeFor[*bool](), true},
		{reflect.TypeFor[*int](), true},
		{reflect.TypeFor[*int32](), true},
		{reflect.TypeFor[*float64](), true},
		{reflect.TypeFor[time.Time](), true},
		{reflect.TypeFor[*time.Time](), true},
		{reflect.TypeFor[openapi_types.Date](), true},
		{reflect.TypeFor[*openapi_types.Date](), true},
		{reflect.TypeFor[[]string](), true},
		{reflect.TypeFor[[]int](), true},
		{reflect.TypeFor[*[]string](), true},
		{reflect.TypeFor[*[]int](), true},
	} {
		if got := ddTargetCanBeNull(tc.typ); got != tc.want {
			t.Errorf("ddTargetCanBeNull(%s) = %v, want %v", tc.typ, got, tc.want)
		}
	}
}

// ---------------------------------------------------------------------------
// Guardrail B: DefectDojo must actually be willing to answer null
// ---------------------------------------------------------------------------

// Minimal projections of the OpenAPI document. The full file is ~756 KB across
// 262 components; only these four fields matter here.
type specDocument struct {
	Components struct {
		Schemas map[string]specSchema `json:"schemas"`
	} `json:"components"`
}

type specSchema struct {
	Properties map[string]specProperty `json:"properties"`
	Required   []string                `json:"required"`
}

type specProperty struct {
	Type     string `json:"type"`
	Format   string `json:"format"`
	Nullable bool   `json:"nullable"`
	ReadOnly bool   `json:"readOnly"`
}

// defaultDdSpecVersion mirrors DD_VERSION in GNUmakefile:3. The Makefile
// exports it (GNUmakefile:4), so every `make test-unit` / `make testacc-local`
// run supplies the real value; this constant only matters for a bare `go test`.
const defaultDdSpecVersion = "3.1.101"

// ddClientPkgPath is the import path of the generated client, used to pick the
// embedded ddclient struct out of a defectdojoResource() wrapper.
const ddClientPkgPath = "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"

// loadOpenAPISpec reads the collected DefectDojo spec, or skips the test when
// none has been collected.
//
// Tests run with the package directory as cwd (ddfield_audit_unit_test.go reads
// "." the same way), so the spec is two levels up.
func loadOpenAPISpec(t *testing.T) (specDocument, string) {
	t.Helper()

	version := cmp.Or(os.Getenv("DD_VERSION"), defaultDdSpecVersion)
	path := filepath.Join("..", "..", "openapi-specs", version, "defect_dojo.json")

	raw, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		// The collected specs are local artifacts (.gitignore:44), so this
		// cannot be a hard requirement - a clean checkout has none. Skipping is
		// honest; failing would make every fresh clone red.
		t.Skipf("OpenAPI spec %s is not present, skipping the spec-backed round-trip audit. "+
			"It is a local artifact: run `make dd-up && make dd-spec` to collect it "+
			"(or set DD_VERSION to a version you have already collected).", path)
	}
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}

	var doc specDocument
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("parsing %s: %v", path, err)
	}
	if len(doc.Components.Schemas) == 0 {
		t.Fatalf("%s parsed but contains no components.schemas; is it really an OpenAPI 3 document?", path)
	}
	return doc, path
}

// embeddedDdClientTypeName finds the generated ddclient struct a
// defectdojoResource() wrapper embeds. Its Go name is also its OpenAPI
// component name, because oapi-codegen names generated types after the
// components they come from (dd.SLAConfiguration -> SLAConfiguration,
// dd.EngagementPresets -> EngagementPresets, dd.LocationProductReference ->
// LocationProductReference, ...).
//
// Scanning for the embed rather than matching on the wrapper's own name is
// deliberate: jira_instance and user wrappers carry an extra write-only shadow
// field alongside the embedded type.
func embeddedDdClientTypeName(wrapper reflect.Type) (string, bool) {
	for f := range wrapper.Fields() {
		if !f.Anonymous {
			continue
		}
		ft := f.Type
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}
		if ft.PkgPath() == ddClientPkgPath {
			return ft.Name(), true
		}
	}
	return "", false
}

// specComponentUnavailable excuses resources whose wrapper embeds no generated
// ddclient struct, so there is no component to join against. Anything not
// listed here must resolve: silently skipping a whole resource is exactly the
// failure mode this audit exists to prevent.
var specComponentUnavailable = map[string]string{
	"defectdojo_notifications": "wraps the hand-written notificationsModel rather than a dd.* type, " +
		"because the live API returns an array for scan_added_empty where the spec declares a " +
		"scalar enum (see the custom UnmarshalJSON in notifications_resource.go)",
}

// TestResourceOptionalAttributesAreNullableInSpec is the audit that finds the
// defects Guardrail A structurally cannot: an attribute whose ddclient field IS
// a pointer - so the Go types look fine - but whose DefectDojo column is not
// nullable, so the server fills it in and answers with a concrete value anyway.
//
// The join is exact rather than name-based:
//
//	tfsdk attribute -> ddField Go field name -> that field's `json:` tag
//	                -> the property on the embedded type's OpenAPI component
//
// A property that is not nullable, not readOnly and not required is one the
// server always answers with. Paired with an attribute that plans as null, that
// is a guaranteed "was null, but now ..." on every apply that omits it.
//
// Skips when no spec has been collected. It runs for real under
// scripts/dd-version-compat.sh, which collects a spec per supported DefectDojo
// version before running the suite - precisely when upstream drift arrives.
func TestResourceOptionalAttributesAreNullableInSpec(t *testing.T) {
	t.Parallel()

	spec, specPath := loadOpenAPISpec(t)
	checked := 0

	for tfTypeName, resp := range providerResourceSchemas(t) {
		if resp.Diagnostics.HasError() {
			continue // reported by TestSchemasValidateImplementation
		}
		targets, ddStruct, ok := resourceDdTargets(t, tfTypeName)
		if !ok {
			continue
		}

		component, ok := embeddedDdClientTypeName(ddStruct)
		if !ok {
			if specComponentUnavailable[tfTypeName] == "" {
				t.Errorf("resource %s: defectdojoResource() embeds no %s struct, so its attributes "+
					"cannot be audited against the OpenAPI spec. Either embed the generated type, "+
					"or add an entry to specComponentUnavailable explaining why not.",
					tfTypeName, ddClientPkgPath)
			}
			continue
		}
		schemaDef, ok := spec.Components.Schemas[component]
		if !ok {
			t.Errorf("resource %s: OpenAPI component %q is missing from %s (spec drift after a "+
				"DefectDojo upgrade?)", tfTypeName, component, specPath)
			continue
		}

		for attrName, attr := range resp.Schema.Attributes {
			if !planningAsNull(attr) {
				continue
			}
			target, ok := targets[attrName]
			if !ok {
				continue // no ddField tag: the engine never writes it
			}
			// Collections are exempt. The read path falls back to SetNull when
			// the server sends an empty collection and the current value is
			// null (resource.go:739, :751, :768, :783), so a server-returned []
			// against a null plan stays null. Scalars have no such branch.
			//
			// time.Time / openapi_types.Date are deliberately NOT exempt here,
			// unlike in Guardrail A: their IsZero() fallback never fires for a
			// non-nullable column the server always populates.
			if isSliceTarget(target.ddType) {
				continue
			}
			if target.jsonName == "" {
				t.Errorf("resource %s: attribute %q resolves to %s.%s, which carries no json tag, "+
					"so it cannot be joined to the spec. A write-only shadow field must not carry "+
					"a ddField tag - see user_resource.go for the pattern that avoids this.",
					tfTypeName, attrName, ddStruct, target.ddFieldTag)
				continue
			}
			prop, ok := schemaDef.Properties[target.jsonName]
			if !ok {
				t.Errorf("resource %s: %s.%s (json:%q) is not a property of OpenAPI component %s",
					tfTypeName, ddStruct, target.ddFieldTag, target.jsonName, component)
				continue
			}
			checked++

			if prop.Nullable || prop.ReadOnly || slices.Contains(schemaDef.Required, target.jsonName) {
				continue
			}
			t.Errorf("resource %s: attribute %q is Optional without Computed and without a "+
				"Default, but OpenAPI component %s.%s (reached via ddField:%q -> json:%q) is "+
				"neither nullable, nor readOnly, nor required. DefectDojo therefore always "+
				"answers with a concrete value, and every apply that omits %q fails with "+
				"\"Provider produced inconsistent result after apply: .%s: was null, but now ...\".\n"+
				"Fix by adding `Computed: true` (DefectDojo picks the value), or `Computed: true` "+
				"plus a `Default:` (the provider picks it).\nSpec: %s",
				tfTypeName, attrName, component, target.jsonName, target.ddFieldTag,
				target.jsonName, attrName, attrName, specPath)
		}
	}

	if checked == 0 {
		t.Fatal("no attributes were checked against the spec; the join is broken and this test " +
			"would pass vacuously")
	}
	t.Logf("audited %d Optional-without-Computed attributes against %s", checked, specPath)
}

// TestDecimalSpecPropertiesCarryDdFormat pins the ddFormat:"decimal" tag to
// exactly the set of string attributes whose OpenAPI property is
// format: decimal.
//
// Today that set has one member, product.revenue, which makes the ddFormat
// mechanism a one-user abstraction. The point of this test is that it cannot
// stay that way silently: a `make regen-client` against a newer DefectDojo that
// adds a second decimal column would otherwise reintroduce the exact bug
// ddFormat was built to fix, with nothing failing.
//
// Both directions are checked, so a tag on a non-decimal property is also an
// error - a tag that does not correspond to real server behaviour is worse than
// no tag at all.
func TestDecimalSpecPropertiesCarryDdFormat(t *testing.T) {
	t.Parallel()

	spec, specPath := loadOpenAPISpec(t)
	tagged, untagged := 0, 0

	for tfTypeName := range providerResourceSchemas(t) {
		targets, ddStruct, ok := resourceDdTargets(t, tfTypeName)
		if !ok {
			continue
		}
		component, ok := embeddedDdClientTypeName(ddStruct)
		if !ok {
			continue // excused by specComponentUnavailable, reported above
		}
		schemaDef, ok := spec.Components.Schemas[component]
		if !ok {
			continue // reported above
		}

		for _, target := range targets {
			if target.tfType != typeOfTypesString || target.jsonName == "" {
				continue
			}
			prop, ok := schemaDef.Properties[target.jsonName]
			if !ok {
				continue
			}

			isDecimal := prop.Format == "decimal"
			hasTag := target.ddFormat == ddFormatDecimal

			switch {
			case isDecimal && !hasTag:
				untagged++
				t.Errorf("resource %s: attribute %q maps to %s.%s, which the spec declares as "+
					"format: decimal, but it carries no ddFormat:%q tag. DefectDojo renders "+
					"decimals in one canonical form (\"100\" comes back as \"100.00\"), so "+
					"without the tag any non-canonical configured value fails the apply with "+
					"\"Provider produced inconsistent result after apply\".\nSpec: %s",
					tfTypeName, target.tfsdkName, component, target.jsonName, ddFormatDecimal, specPath)
			case !isDecimal && hasTag:
				t.Errorf("resource %s: attribute %q carries ddFormat:%q, but %s.%s is not "+
					"format: decimal in the spec (it is %q). Remove the tag, or confirm the "+
					"server really does canonicalise this field.\nSpec: %s",
					tfTypeName, target.tfsdkName, ddFormatDecimal, component, target.jsonName,
					prop.Format, specPath)
			case isDecimal:
				tagged++
			}
		}
	}

	if tagged == 0 && untagged == 0 {
		t.Errorf("no format: decimal properties were found in %s at all. product.revenue is "+
			"expected to be one, so either the join is broken or the spec changed shape.", specPath)
	}
	t.Logf("matched %d decimal properties to ddFormat tags", tagged)
}
