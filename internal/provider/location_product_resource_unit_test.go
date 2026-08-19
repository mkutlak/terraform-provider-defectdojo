package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
	"gotest.tools/assert"
)

func TestLocationProductResourcePopulate(t *testing.T) {
	expectedId := 5
	expectedLocation := 10
	expectedProduct := 20
	expectedRelationship := dd.LocationProductReferenceRelationship("owned_by")
	expectedStatus := dd.LocationProductReferenceStatus("Active")
	expectedLocationType := "url"
	expectedLocationValue := "https://example.com"

	ddResource := locationProductDefectdojoResource{
		LocationProductReference: dd.LocationProductReference{
			Id:            &expectedId,
			Location:      expectedLocation,
			Product:       expectedProduct,
			Relationship:  &expectedRelationship,
			Status:        &expectedStatus,
			LocationType:  &expectedLocationType,
			LocationValue: &expectedLocationValue,
		},
	}

	resourceData := locationProductResourceData{}
	var trd terraformResourceData = &resourceData

	populateResourceData(context.Background(), &diag.Diagnostics{}, &trd, &ddResource)

	assert.Equal(t, resourceData.Id.ValueString(), "5")
	assert.Equal(t, resourceData.Location.ValueInt64(), int64(expectedLocation))
	assert.Equal(t, resourceData.Product.ValueInt64(), int64(expectedProduct))
	assert.Equal(t, resourceData.Relationship.ValueString(), string(expectedRelationship))
	assert.Equal(t, resourceData.Status.ValueString(), string(expectedStatus))
	assert.Equal(t, resourceData.LocationType.ValueString(), expectedLocationType)
	assert.Equal(t, resourceData.LocationValue.ValueString(), expectedLocationValue)

	ddResource = locationProductDefectdojoResource{
		LocationProductReference: dd.LocationProductReference{},
	}
	populateResourceData(context.Background(), &diag.Diagnostics{}, &trd, &ddResource)

	assert.Equal(t, resourceData.Id.IsNull(), true)
	// Location and Product are bare (non-pointer) ints on the generated
	// struct, so a zero-value response leaves them at their int zero value
	// rather than null.
	assert.Equal(t, resourceData.Location.ValueInt64(), int64(0))
	assert.Equal(t, resourceData.Product.ValueInt64(), int64(0))
	assert.Equal(t, resourceData.Relationship.IsNull(), true)
	assert.Equal(t, resourceData.Status.IsNull(), true)
	assert.Equal(t, resourceData.LocationType.IsNull(), true)
	assert.Equal(t, resourceData.LocationValue.IsNull(), true)
}

func TestLocationProductResourcePopulateNils(t *testing.T) {
	resourceData := locationProductResourceData{}
	var trd terraformResourceData = &resourceData

	assert.Equal(t, resourceData.Id.ValueString(), "")
	assert.Equal(t, resourceData.Location.ValueInt64(), int64(0))
	assert.Equal(t, resourceData.Product.ValueInt64(), int64(0))
	assert.Equal(t, resourceData.Relationship.ValueString(), "")
	assert.Equal(t, resourceData.Status.ValueString(), "")
	assert.Equal(t, resourceData.LocationType.ValueString(), "")
	assert.Equal(t, resourceData.LocationValue.ValueString(), "")

	ddResource := locationProductDefectdojoResource{
		LocationProductReference: dd.LocationProductReference{},
	}
	populateResourceData(context.Background(), &diag.Diagnostics{}, &trd, &ddResource)

	// still all empty/null values after running populate against a zero-value response
	assert.Equal(t, resourceData.Id.ValueString(), "")
	assert.Equal(t, resourceData.Location.ValueInt64(), int64(0))
	assert.Equal(t, resourceData.Product.ValueInt64(), int64(0))
	assert.Equal(t, resourceData.Relationship.ValueString(), "")
	assert.Equal(t, resourceData.Status.ValueString(), "")
	assert.Equal(t, resourceData.LocationType.ValueString(), "")
	assert.Equal(t, resourceData.LocationValue.ValueString(), "")
}

func TestLocationProductResource__defectdojoResource(t *testing.T) {
	expectedLocation := 10
	expectedProduct := 20
	expectedRelationship := "owned_by"
	expectedStatus := "Active"

	resourceData := locationProductResourceData{
		Location:     types.Int64Value(int64(expectedLocation)),
		Product:      types.Int64Value(int64(expectedProduct)),
		Relationship: types.StringValue(expectedRelationship),
		Status:       types.StringValue(expectedStatus),
	}

	ddRes := resourceData.defectdojoResource()
	ddLocationProduct := ddRes.(*locationProductDefectdojoResource)
	var trd terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, trd, &ddRes)

	assert.Equal(t, ddLocationProduct.Location, expectedLocation)
	assert.Equal(t, ddLocationProduct.Product, expectedProduct)
	assert.Equal(t, string(*ddLocationProduct.Relationship), expectedRelationship)
	assert.Equal(t, string(*ddLocationProduct.Status), expectedStatus)

	// locationProductToRequest must cast from the response enum types
	// (LocationProductReferenceRelationship/Status) to the distinct request
	// enum types (LocationProductReferenceRequestRelationship/Status).
	req := locationProductToRequest(ddLocationProduct.LocationProductReference)
	assert.Equal(t, req.Location, expectedLocation)
	assert.Equal(t, req.Product, expectedProduct)
	assert.Equal(t, string(*req.Relationship), expectedRelationship)
	assert.Equal(t, string(*req.Status), expectedStatus)
}

func TestLocationProductResource__defectdojoResource_Nulls(t *testing.T) {
	var nilRelationship *dd.LocationProductReferenceRelationship
	var nilStatus *dd.LocationProductReferenceStatus
	var nilString *string
	var nilInt *int

	resourceData := locationProductResourceData{
		Id:            types.StringNull(),
		Location:      types.Int64Null(),
		Product:       types.Int64Null(),
		Relationship:  types.StringNull(),
		Status:        types.StringNull(),
		LocationType:  types.StringNull(),
		LocationValue: types.StringNull(),
	}

	ddRes := resourceData.defectdojoResource()
	ddLocationProduct := ddRes.(*locationProductDefectdojoResource)
	var trd terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, trd, &ddRes)

	assert.Equal(t, ddLocationProduct.Id, nilInt)
	// Null TF values are skipped by populateDefectdojoResource, so the bare
	// (non-pointer) Location/Product fields remain at their int zero value.
	assert.Equal(t, ddLocationProduct.Location, 0)
	assert.Equal(t, ddLocationProduct.Product, 0)
	assert.Equal(t, ddLocationProduct.Relationship, nilRelationship)
	assert.Equal(t, ddLocationProduct.Status, nilStatus)
	assert.Equal(t, ddLocationProduct.LocationType, nilString)
	assert.Equal(t, ddLocationProduct.LocationValue, nilString)

	var nilReqRelationship *dd.LocationProductReferenceRequestRelationship
	var nilReqStatus *dd.LocationProductReferenceRequestStatus
	req := locationProductToRequest(ddLocationProduct.LocationProductReference)
	assert.Equal(t, req.Location, 0)
	assert.Equal(t, req.Product, 0)
	assert.Equal(t, req.Relationship, nilReqRelationship)
	assert.Equal(t, req.Status, nilReqStatus)
}

// TestLocationProductRelationshipAcceptsBlank pins the empty string as a value
// a configuration may carry.
//
// A create that omits relationship stores and reads back "", so state holds a
// value the OneOf used to reject at plan time:
//
//	Attribute relationship value must be one of: ["owned_by" "used_by"], got: ""
//
// That made the blank the provider itself produces unwritable, and since a
// Computed attribute cannot be cleared by deleting it from configuration
// either, there was no way back to it. "" is a genuine member of the enum -
// dd.LocationProductReferenceRequestRelationshipEmpty - and DefectDojo 3.1.101
// accepts it on create and update.
func TestLocationProductRelationshipAcceptsBlank(t *testing.T) {
	t.Parallel()

	attr := resourceStringAttribute(t, "defectdojo_location_product", "relationship")
	for _, tc := range []struct {
		value     string
		wantError bool
	}{
		{"", false},
		{"owned_by", false},
		{"used_by", false},
		{"bogus", true},
		{"Owned_By", true},
	} {
		got := runStringValidators(t, "relationship", attr, tc.value)
		if got != tc.wantError {
			verb := "accepted"
			if got {
				verb = "rejected"
			}
			t.Errorf("relationship = %q was %s, but wantError=%v", tc.value, verb, tc.wantError)
		}
	}
}

// TestLocationProductStatusRejectsBlank is the other half of the pair above.
//
// status looks like the same shape - Optional, Computed, a OneOf - but its enum
// has no blank member: DefectDojo 3.1.101 answers a status of "" with
// `"" is not a valid choice.` and fills the field with "Mitigated" when it is
// omitted, so there is no state value the validator has to be able to express.
// Without this, "fix the relationship validator" reads like an argument for
// loosening every enum on the resource.
func TestLocationProductStatusRejectsBlank(t *testing.T) {
	t.Parallel()

	attr := resourceStringAttribute(t, "defectdojo_location_product", "status")
	if !runStringValidators(t, "status", attr, "") {
		t.Error(`status = "" was accepted, but DefectDojo rejects it as an invalid choice`)
	}
	for _, value := range []string{"Active", "Mitigated"} {
		if runStringValidators(t, "status", attr, value) {
			t.Errorf("status = %q was rejected, but it is a valid choice", value)
		}
	}
}
