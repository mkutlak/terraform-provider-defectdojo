package provider

// This file contains a reflection audit that protects the ddField struct tag
// contract between the Terraform resource data structs and the generated
// ddclient structs (internal/ddclient).
//
// Why this exists: populateDefectdojoResource() and populateResourceData()
// (resource.go) look up DD struct fields by the `ddField` tag at RUNTIME via
// FieldByName(). If the generated client is regenerated from a newer
// DefectDojo OpenAPI spec and a field is renamed, removed, or changes type,
// nothing fails at compile time; users only see "No such field" diagnostics
// (or silently dropped values on a type mismatch) when they run terraform.
// Ordinary populate unit tests do not catch this either, because null TF
// fields are skipped before the FieldByName lookup.
//
// This audit makes any such drift a plain `go test` failure:
//
//  1. Every `ddField` tag on every terraformResourceData implementation must
//     resolve to a field on the struct returned by defectdojoResource().
//  2. The (terraform type, dd type) pairing must be one the reflection engine
//     in resource.go actually handles (see ddFieldPairingSupported below,
//     which mirrors the switch statements in populateDefectdojoResource and
//     populateResourceData).
//  3. The table below must stay complete: a companion check parses the
//     package sources (go/ast) and fails if a defectdojoResource() method
//     exists on a type that is not listed here.
//
// HOW TO ADD A NEW RESOURCE: add one line to ddFieldAuditTable with the new
// *ResourceData zero value. The completeness check (TestDdFieldAuditTableIsComplete)
// will remind you if you forget.

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

// ddFieldAuditTable lists every terraformResourceData implementation in this
// package. Keys are the Go type names of the implementations (used by the
// completeness check); values are zero-value instances whose
// defectdojoResource() result is audited against their ddField tags.
var ddFieldAuditTable = map[string]terraformResourceData{
	"announcementResourceData":                &announcementResourceData{},
	"configurationPermissionResourceData":     &configurationPermissionResourceData{},
	"ddTestResourceData":                      &ddTestResourceData{},
	"developmentEnvironmentResourceData":      &developmentEnvironmentResourceData{},
	"endpointResourceData":                    &endpointResourceData{},
	"engagementPresetResourceData":            &engagementPresetResourceData{},
	"engagementResourceData":                  &engagementResourceData{},
	"findingTemplateResourceData":             &findingTemplateResourceData{},
	"jiraInstanceResourceData":                &jiraInstanceResourceData{},
	"jiraProductConfigurationResourceData":    &jiraProductConfigurationResourceData{},
	"languageTypeResourceData":                &languageTypeResourceData{},
	"locationProductResourceData":             &locationProductResourceData{},
	"locationResourceData":                    &locationResourceData{},
	"metadataResourceData":                    &metadataResourceData{},
	"networkLocationResourceData":             &networkLocationResourceData{},
	"noteTypeResourceData":                    &noteTypeResourceData{},
	"notificationWebhookResourceData":         &notificationWebhookResourceData{},
	"notificationsResourceData":               &notificationsResourceData{},
	"productAPIScanConfigurationResourceData": &productAPIScanConfigurationResourceData{},
	"productResourceData":                     &productResourceData{},
	"productTypeResourceData":                 &productTypeResourceData{},
	"regulationResourceData":                  &regulationResourceData{},
	"riskAcceptanceResourceData":              &riskAcceptanceResourceData{},
	"slaConfigurationResourceData":            &slaConfigurationResourceData{},
	"systemSettingsResourceData":              &systemSettingsResourceData{},
	"testTypeResourceData":                    &testTypeResourceData{},
	"toolConfigurationResourceData":           &toolConfigurationResourceData{},
	"toolProductSettingsResourceData":         &toolProductSettingsResourceData{},
	"toolTypeResourceData":                    &toolTypeResourceData{},
	"urlResourceData":                         &urlResourceData{},
	"userContactInfoResourceData":             &userContactInfoResourceData{},
	"userProfileResourceData":                 &userProfileResourceData{},
	"userResourceData":                        &userResourceData{},
}

var (
	auditTypeInt     = reflect.TypeFor[int]()
	auditTypeInt32   = reflect.TypeFor[int32]()
	auditTypeBool    = reflect.TypeFor[bool]()
	auditTypeFloat32 = reflect.TypeFor[float32]()
	auditTypeFloat64 = reflect.TypeFor[float64]()
	auditTypeTime    = reflect.TypeFor[time.Time]()
	auditTypeDate    = reflect.TypeFor[openapi_types.Date]()
)

// ddFieldPairingSupported reports whether the reflection engine in resource.go
// can convert between the given Terraform framework type and the given
// ddclient field type in BOTH directions (populateDefectdojoResource and
// populateResourceData). It mirrors the switch statements in resource.go,
// including the places where the engine uses reflect.Value.Set with an
// exactly-typed value (which panics or silently warns for defined types)
// versus reflect.Value.Convert (which tolerates defined types with a matching
// underlying kind).
func ddFieldPairingSupported(tfType, ddType reflect.Type) (bool, string) {
	isPtrTo := func(t reflect.Type, elem func(reflect.Type) bool) bool {
		return t.Kind() == reflect.Ptr && elem(t.Elem())
	}

	switch tfType {
	case typeOfTypesString:
		switch {
		// string / defined string types: engine converts via Convert() and
		// reads via .String(), so named enum string types are fine.
		case ddType.Kind() == reflect.String:
			return true, ""
		case isPtrTo(ddType, func(e reflect.Type) bool { return e.Kind() == reflect.String }):
			return true, ""
		// int paths assign a raw int value (reflect.ValueOf(intVal)), so the
		// DD type must be exactly int / *int.
		case ddType == auditTypeInt:
			return true, ""
		case isPtrTo(ddType, func(e reflect.Type) bool { return e == auditTypeInt }):
			return true, ""
		// datetime / date paths use exact types.
		case ddType == auditTypeTime:
			return true, ""
		case isPtrTo(ddType, func(e reflect.Type) bool { return e == auditTypeTime }):
			return true, ""
		case ddType == auditTypeDate:
			return true, ""
		case isPtrTo(ddType, func(e reflect.Type) bool { return e == auditTypeDate }):
			return true, ""
		}
		return false, "types.String maps only to string, *string (defined string types ok), int, *int, time.Time, *time.Time, openapi_types.Date, *openapi_types.Date"

	case typeOfTypesBool:
		switch {
		// bare bool is assigned with an exactly-typed value.
		case ddType == auditTypeBool:
			return true, ""
		// *bool uses Convert() on the element, so defined bool types are ok.
		case isPtrTo(ddType, func(e reflect.Type) bool { return e.Kind() == reflect.Bool }):
			return true, ""
		}
		return false, "types.Bool maps only to bool or *bool"

	case typeOfTypesInt64:
		switch {
		// int / *int use Convert(), so defined int types are ok.
		case ddType.Kind() == reflect.Int:
			return true, ""
		case isPtrTo(ddType, func(e reflect.Type) bool { return e.Kind() == reflect.Int }):
			return true, ""
		// int32 / *int32 assign exactly-typed values.
		case ddType == auditTypeInt32:
			return true, ""
		case isPtrTo(ddType, func(e reflect.Type) bool { return e == auditTypeInt32 }):
			return true, ""
		}
		return false, "types.Int64 maps only to int, *int (defined int types ok), int32, *int32"

	case typeOfTypesFloat64:
		switch {
		// all float paths assign exactly-typed values (or type-assert on read).
		case ddType == auditTypeFloat64:
			return true, ""
		case isPtrTo(ddType, func(e reflect.Type) bool { return e == auditTypeFloat64 }):
			return true, ""
		case ddType == auditTypeFloat32:
			return true, ""
		case isPtrTo(ddType, func(e reflect.Type) bool { return e == auditTypeFloat32 }):
			return true, ""
		}
		return false, "types.Float64 maps only to float64, *float64, float32 or *float32"

	case typeOfTypesSet:
		switch {
		// direct slices: the engine converts element-wise via reflect.Convert,
		// so any int-kind or string-kind element type works (incl. defined
		// types like enum strings).
		case ddType.Kind() == reflect.Slice && (ddType.Elem().Kind() == reflect.Int || ddType.Elem().Kind() == reflect.String):
			return true, ""
		// pointer to slice: same element-wise conversion in both directions
		// (populateResourceData iterates reflectively, no type assertion).
		case isPtrTo(ddType, func(e reflect.Type) bool {
			return e.Kind() == reflect.Slice && (e.Elem().Kind() == reflect.Int || e.Elem().Kind() == reflect.String)
		}):
			return true, ""
		}
		return false, "types.Set maps only to slices (or pointers to slices) whose element kind is int or string"
	}

	return false, fmt.Sprintf("terraform type %s is not handled by the reflection engine at all (only types.String, types.Bool, types.Int64, types.Float64, types.Set are)", tfType)
}

// ddStructType unwraps the defectdojoResource() result down to the struct
// type the reflection engine operates on: resource.go does
// reflect.ValueOf(*ddResource).Elem(), i.e. it dereferences the pointer held
// in the interface and does FieldByName on the wrapper struct (which promotes
// the embedded ddclient struct's fields).
func ddStructType(t *testing.T, name string, ddRes defectdojoResource) reflect.Type {
	t.Helper()
	ddType := reflect.TypeOf(ddRes)
	if ddType.Kind() != reflect.Ptr || ddType.Elem().Kind() != reflect.Struct {
		t.Fatalf("%s: defectdojoResource() must return a pointer to a struct, got %s", name, ddType)
	}
	return ddType.Elem()
}

// TestDdFieldAuditTagsResolve asserts that every non-empty ddField tag on
// every terraformResourceData implementation resolves to an existing field of
// the corresponding ddclient struct, and that the type pairing is one the
// reflection engine supports. This is the test that fails when the generated
// client drifts (field renamed/removed/retyped in a new DefectDojo spec).
func TestDdFieldAuditTagsResolve(t *testing.T) {
	totalTags := 0

	for name, data := range ddFieldAuditTable {
		t.Run(name, func(t *testing.T) {
			tfType := reflect.TypeOf(data)
			if tfType.Kind() != reflect.Ptr || tfType.Elem().Kind() != reflect.Struct {
				t.Fatalf("%s: table entry must be a pointer to a struct, got %s", name, tfType)
			}
			tfStruct := tfType.Elem()

			ddStruct := ddStructType(t, name, data.defectdojoResource())

			tagged := 0
			for tfField := range tfStruct.Fields() {
				tag := tfField.Tag.Get("ddField")
				if tag == "" {
					continue
				}
				tagged++
				totalTags++

				ddField, ok := ddStruct.FieldByName(tag)
				if !ok {
					t.Errorf("%s: field %q has ddField:%q but %s has no such field (generated client drift: the runtime engine would emit a \"No such field\" diagnostic)",
						name, tfField.Name, tag, ddStruct)
					continue
				}

				if ok, reason := ddFieldPairingSupported(tfField.Type, ddField.Type); !ok {
					t.Errorf("%s: field %q (ddField:%q) pairs %s with %s.%s of type %s, which the reflection engine in resource.go cannot convert: %s",
						name, tfField.Name, tag, tfField.Type, ddStruct, tag, ddField.Type, reason)
				}
			}

			if tagged == 0 {
				t.Errorf("%s: no ddField tags found; either the struct lost its tags or the table entry is wrong", name)
			}
		})
	}

	t.Logf("audited %d resource models, %d ddField tags total", len(ddFieldAuditTable), totalTags)
}

// TestDdFieldAuditTableIsComplete parses the package sources and fails if a
// defectdojoResource() method exists on a receiver type that is missing from
// ddFieldAuditTable (i.e. someone added a new resource but forgot to add it
// to the audit above).
func TestDdFieldAuditTableIsComplete(t *testing.T) {
	fset := token.NewFileSet()
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("failed to read package directory: %v", err)
	}

	implementations := map[string]bool{}
	for _, entry := range entries {
		fileName := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(fileName, ".go") || strings.HasSuffix(fileName, "_test.go") {
			continue
		}
		file, err := parser.ParseFile(fset, fileName, nil, 0)
		if err != nil {
			t.Fatalf("failed to parse %s: %v", fileName, err)
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Name.Name != "defectdojoResource" || fn.Recv == nil || len(fn.Recv.List) != 1 {
				continue
			}
			recv := receiverTypeName(fn.Recv.List[0].Type)
			if recv == "" {
				t.Errorf("%s: could not determine receiver type of defectdojoResource()", fileName)
				continue
			}
			implementations[recv] = true
		}
	}

	if len(implementations) == 0 {
		t.Fatal("found no defectdojoResource() implementations in package sources; the parser filter is probably broken")
	}

	for recv := range implementations {
		if _, ok := ddFieldAuditTable[recv]; !ok {
			t.Errorf("type %s implements defectdojoResource() but is missing from ddFieldAuditTable in ddfield_audit_unit_test.go — add it so its ddField tags are audited", recv)
		}
	}
	for name := range ddFieldAuditTable {
		if !implementations[name] {
			t.Errorf("ddFieldAuditTable lists %s but no defectdojoResource() method was found for it in the package sources — remove the stale entry", name)
		}
	}

	t.Logf("found %d defectdojoResource() implementations in sources, %d entries in audit table", len(implementations), len(ddFieldAuditTable))
}

func receiverTypeName(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.StarExpr:
		return receiverTypeName(e.X)
	case *ast.Ident:
		return e.Name
	case *ast.IndexExpr: // generic receiver, not expected here
		return receiverTypeName(e.X)
	}
	return ""
}
