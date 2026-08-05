package provider

import (
	"context"
	"slices"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

// The generic acceptance checks in acctest_checks_test.go need, for any
// Terraform type name found in state, a zero-valued model whose
// defectdojoResource() exposes readApiCall/deleteApiCall. Rather than
// hand-maintaining a typeName -> model table (which would silently rot the
// moment a resource is added), the table is DERIVED from the provider itself.
//
// This works because every resource embeds terraformResource by value, and
// terraformResource embeds the dataProvider interface, so getData is promoted
// onto the resource type. Calling it with a no-op dataGetter yields the zero
// value of the resource's model - exactly what readApiCall needs, since only
// the id (passed separately) matters.
//
// TestAccResourceRegistryIsDerivable is the guard: if a future resource stops
// embedding terraformResource, derivation would silently skip it and its
// destroy check would vanish without any test failing. That test makes the
// omission loud.

// dataProviderResource is the promoted-method view of a resource that is backed
// by the reflection CRUD engine.
type dataProviderResource interface {
	getData(context.Context, dataGetter) (terraformResourceData, diag.Diagnostics)
}

// nullDataGetter satisfies dataGetter without touching the target, so getData
// returns the model's zero value.
type nullDataGetter struct{}

func (nullDataGetter) Get(context.Context, any) diag.Diagnostics { return nil }

// testAccResourceFactories maps Terraform type name -> zero-value model
// factory, derived once from the provider's own resource list.
var testAccResourceFactories = sync.OnceValue(func() map[string]func() terraformResourceData {
	ctx := context.Background()
	out := map[string]func() terraformResourceData{}

	for _, newResource := range New("test")().Resources(ctx) {
		r := newResource()

		var md fwresource.MetadataResponse
		r.Metadata(ctx, fwresource.MetadataRequest{ProviderTypeName: "defectdojo"}, &md)

		dp, ok := r.(dataProviderResource)
		if !ok {
			continue // reported by TestAccResourceRegistryIsDerivable
		}

		out[md.TypeName] = func() terraformResourceData {
			data, diags := dp.getData(ctx, nullDataGetter{})
			if diags.HasError() {
				return nil
			}
			return data
		}
	}
	return out
})

// testAccResourceModel returns a freshly zero-valued model for the given
// Terraform type name (e.g. "defectdojo_product"), plus whether the type is
// known to the provider at all.
func testAccResourceModel(typeName string) (terraformResourceData, bool) {
	factory, ok := testAccResourceFactories()[typeName]
	if !ok {
		return nil, false
	}
	data := factory()
	return data, data != nil
}

// testAccIsUndestroyable reports whether the resource cannot be deleted
// server-side, so a destroy check must not expect a 404. This is detected from
// the singletonAdopter interface rather than a maintained list: the engine
// itself uses that same interface to turn Delete into a state-removal
// (resource.go:339-344), so the two can never disagree.
func testAccIsUndestroyable(data terraformResourceData) bool {
	_, ok := data.defectdojoResource().(singletonAdopter)
	return ok
}

// TestAccResourceRegistryIsDerivable asserts that every resource the provider
// registers can be resolved to a model by testAccResourceModel. A resource that
// fails here would silently lose its CheckDestroy coverage.
//
// This is a unit test: it needs no backend and runs in the fast CI job.
func TestAccResourceRegistryIsDerivable(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	resources := New("test")().Resources(ctx)

	if len(resources) == 0 {
		t.Fatal("provider registered no resources; the derivation below would vacuously pass")
	}

	seen := map[string]bool{}

	for _, newResource := range resources {
		r := newResource()

		var md fwresource.MetadataResponse
		r.Metadata(ctx, fwresource.MetadataRequest{ProviderTypeName: "defectdojo"}, &md)

		if md.TypeName == "" {
			t.Errorf("%T: Metadata produced an empty TypeName", r)
			continue
		}
		if seen[md.TypeName] {
			t.Errorf("%s: registered more than once", md.TypeName)
			continue
		}
		seen[md.TypeName] = true

		data, ok := testAccResourceModel(md.TypeName)
		if !ok {
			t.Errorf("%s (%T): could not derive a model. It likely does not embed "+
				"terraformResource, so the generic CheckDestroy would silently skip it.",
				md.TypeName, r)
			continue
		}

		// defectdojoResource() must not panic and must be non-nil, since the
		// generic checks call readApiCall/deleteApiCall on it.
		ddResource := data.defectdojoResource()
		if ddResource == nil {
			t.Errorf("%s: defectdojoResource() returned nil", md.TypeName)
		}
	}

	t.Logf("derived models for %d resources", len(seen))
}

// TestAccResourceRegistryReportsSingletons documents which resources are exempt
// from the 404-after-destroy expectation, so the exemption set is visible in
// test output rather than buried in an interface assertion.
func TestAccResourceRegistryReportsSingletons(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	var singletons []string

	for _, newResource := range New("test")().Resources(ctx) {
		r := newResource()

		var md fwresource.MetadataResponse
		r.Metadata(ctx, fwresource.MetadataRequest{ProviderTypeName: "defectdojo"}, &md)

		data, ok := testAccResourceModel(md.TypeName)
		if !ok {
			continue // reported by TestAccResourceRegistryIsDerivable
		}
		if testAccIsUndestroyable(data) {
			singletons = append(singletons, md.TypeName)
		}
	}

	t.Logf("resources exempt from destroy verification (singletonAdopter): %v", singletons)

	// system_settings is the known singleton. Assert it explicitly: if it ever
	// stopped implementing singletonAdopter, the engine would start issuing a
	// real DELETE and the generic destroy check would start expecting a 404,
	// producing a confusing failure far from the cause.
	if !slices.Contains(singletons, "defectdojo_system_settings") {
		t.Error("defectdojo_system_settings no longer implements singletonAdopter; " +
			"destroy semantics changed (see resource.go:339-344)")
	}
}
