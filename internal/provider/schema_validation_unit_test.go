package provider

import (
	"context"
	"sort"
	"testing"

	fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

// Schema meta-tests. These walk every resource and data source the provider
// registers and assert invariants that would otherwise only surface at
// `terraform plan` time or as a silently undocumented attribute in docs/.
//
// They need no backend and run in the fast CI job.

// schemaAttribute is the subset of the framework's attribute interface these
// tests rely on. Declaring it locally avoids importing the framework's internal
// fwschema package.
type schemaAttribute interface {
	GetDescription() string
	GetMarkdownDescription() string
	IsRequired() bool
	IsOptional() bool
	IsComputed() bool
}

func providerResourceSchemas(t *testing.T) map[string]fwresource.SchemaResponse {
	t.Helper()

	ctx := context.Background()
	out := map[string]fwresource.SchemaResponse{}

	for _, newResource := range New("test")().Resources(ctx) {
		r := newResource()

		var md fwresource.MetadataResponse
		r.Metadata(ctx, fwresource.MetadataRequest{ProviderTypeName: "defectdojo"}, &md)

		var resp fwresource.SchemaResponse
		r.Schema(ctx, fwresource.SchemaRequest{}, &resp)
		out[md.TypeName] = resp
	}
	return out
}

func providerDataSourceSchemas(t *testing.T) map[string]fwdatasource.SchemaResponse {
	t.Helper()

	ctx := context.Background()
	out := map[string]fwdatasource.SchemaResponse{}

	for _, newDataSource := range New("test")().DataSources(ctx) {
		d := newDataSource()

		var md fwdatasource.MetadataResponse
		d.Metadata(ctx, fwdatasource.MetadataRequest{ProviderTypeName: "defectdojo"}, &md)

		var resp fwdatasource.SchemaResponse
		d.Schema(ctx, fwdatasource.SchemaRequest{}, &resp)
		out[md.TypeName] = resp
	}
	return out
}

// TestSchemasValidateImplementation runs the framework's own structural
// validation over every schema. It catches invalid attribute names, attributes
// that are neither Required nor Optional nor Computed, and defaults declared on
// non-computed attributes - all of which are runtime errors otherwise.
func TestSchemasValidateImplementation(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	resources := providerResourceSchemas(t)
	if len(resources) == 0 {
		t.Fatal("no resources found; this test would vacuously pass")
	}
	for name, resp := range resources {
		if resp.Diagnostics.HasError() {
			t.Errorf("resource %s: Schema() diagnostics: %v", name, resp.Diagnostics)
			continue
		}
		if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
			t.Errorf("resource %s: ValidateImplementation: %v", name, diags)
		}
	}

	dataSources := providerDataSourceSchemas(t)
	if len(dataSources) == 0 {
		t.Fatal("no data sources found; this test would vacuously pass")
	}
	for name, resp := range dataSources {
		if resp.Diagnostics.HasError() {
			t.Errorf("data source %s: Schema() diagnostics: %v", name, resp.Diagnostics)
			continue
		}
		if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
			t.Errorf("data source %s: ValidateImplementation: %v", name, diags)
		}
	}

	t.Logf("validated %d resource and %d data source schemas", len(resources), len(dataSources))
}

// TestSchemaAttributesAreDocumented asserts every attribute carries a
// description. tfplugindocs renders MarkdownDescription when present and falls
// back to Description, so either satisfies it; an attribute with neither ships
// as a blank row in docs/.
func TestSchemaAttributesAreDocumented(t *testing.T) {
	t.Parallel()

	var undocumented []string

	check := func(kind, owner, attrName string, attr any) {
		a, ok := attr.(schemaAttribute)
		if !ok {
			t.Errorf("%s %s: attribute %q does not implement the expected interface (%T)",
				kind, owner, attrName, attr)
			return
		}
		if a.GetDescription() == "" && a.GetMarkdownDescription() == "" {
			undocumented = append(undocumented, kind+" "+owner+"."+attrName)
		}
	}

	for name, resp := range providerResourceSchemas(t) {
		for attrName, attr := range resp.Schema.Attributes {
			check("resource", name, attrName, attr)
		}
	}
	for name, resp := range providerDataSourceSchemas(t) {
		for attrName, attr := range resp.Schema.Attributes {
			check("data source", name, attrName, attr)
		}
	}

	if len(undocumented) > 0 {
		sort.Strings(undocumented)
		t.Errorf("%d attribute(s) have neither Description nor MarkdownDescription "+
			"and will render blank in docs/:\n  %v", len(undocumented), undocumented)
	}
}

// TestResourceSchemasHaveComputedID asserts every resource exposes a Computed
// "id". The CRUD engine reads and writes state through terraformResourceData.id()
// and ImportState is a passthrough on path.Root("id") (resource.go:366-368), so a
// resource without one cannot be imported or refreshed.
func TestResourceSchemasHaveComputedID(t *testing.T) {
	t.Parallel()

	for name, resp := range providerResourceSchemas(t) {
		attr, ok := resp.Schema.Attributes["id"]
		if !ok {
			t.Errorf("resource %s: no \"id\" attribute; ImportState and Read both require one", name)
			continue
		}
		a, ok := attr.(schemaAttribute)
		if !ok {
			t.Errorf("resource %s: \"id\" does not implement the expected interface (%T)", name, attr)
			continue
		}
		if !a.IsComputed() {
			t.Errorf("resource %s: \"id\" must be Computed", name)
		}
	}
}

// TestDataSourceSchemasHaveID asserts every data source exposes an "id". The
// generic datasource Read resolves the object by numeric id, either supplied
// directly or filled in by a name-based lookup (datasource.go:65-96).
func TestDataSourceSchemasHaveID(t *testing.T) {
	t.Parallel()

	for name, resp := range providerDataSourceSchemas(t) {
		if _, ok := resp.Schema.Attributes["id"]; !ok {
			t.Errorf("data source %s: no \"id\" attribute; the generic Read requires one", name)
		}
	}
}
