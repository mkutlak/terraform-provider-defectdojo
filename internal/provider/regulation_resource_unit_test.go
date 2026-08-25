package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
	"gotest.tools/assert"
)

func TestRegulationResourcePopulate(t *testing.T) {
	expectedId := 99
	expectedName := "Test Regulation"
	expectedAcronym := "TST"
	expectedCategory := "other"
	expectedJurisdiction := "US"
	expectedDescription := "A test regulation"
	expectedReference := "https://example.com"

	// Reference is an oapi-codegen oneOf-string wrapper as of DefectDojo 3.2
	// (see isOapiUnionStringType in resource.go); it cannot be built with a
	// plain &string literal.
	reference := &dd.Regulation_Reference{}
	assert.NilError(t, reference.FromRegulationReference0(expectedReference))

	ddObj := regulationDefectdojoResource{
		Regulation: dd.Regulation{
			Id:           &expectedId,
			Name:         expectedName,
			Acronym:      expectedAcronym,
			Category:     dd.RegulationCategory(expectedCategory),
			Jurisdiction: expectedJurisdiction,
			Description:  &expectedDescription,
			Reference:    reference,
		},
	}

	resourceData := regulationResourceData{}
	var tfResource terraformResourceData = &resourceData
	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddObj)

	assert.Equal(t, resourceData.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, resourceData.Name.ValueString(), expectedName)
	assert.Equal(t, resourceData.Acronym.ValueString(), expectedAcronym)
	assert.Equal(t, resourceData.Category.ValueString(), expectedCategory)
	assert.Equal(t, resourceData.Jurisdiction.ValueString(), expectedJurisdiction)
	assert.Equal(t, resourceData.Description.ValueString(), expectedDescription)
	assert.Equal(t, resourceData.Reference.ValueString(), expectedReference)
}

// TestRegulationResourcePopulateNils covers the read path when Reference is
// absent server-side, so its wrapper pointer is nil. Reference is an
// oapi-codegen oneOf-string wrapper (see isOapiUnionStringType in
// resource.go).
func TestRegulationResourcePopulateNils(t *testing.T) {
	ddObj := regulationDefectdojoResource{
		Regulation: dd.Regulation{
			Category:     dd.RegulationCategory("other"),
			Jurisdiction: "US",
		},
	}

	resourceData := regulationResourceData{}
	var tfResource terraformResourceData = &resourceData
	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddObj)

	assert.Equal(t, resourceData.Reference.IsNull(), true)
}

// TestRegulationResourcePopulate_ReferenceEmpty exercises the read path for
// the oneOf's second arm (maxLength: 0), which DefectDojo sends when
// Reference is blank. Empty-string coverage was flagged as the top gap in
// the isOapiUnionStringType engine (resource.go).
func TestRegulationResourcePopulate_ReferenceEmpty(t *testing.T) {
	reference := &dd.Regulation_Reference{}
	assert.NilError(t, reference.FromRegulationReference1(""))

	ddObj := regulationDefectdojoResource{
		Regulation: dd.Regulation{
			Category:     dd.RegulationCategory("other"),
			Jurisdiction: "US",
			Reference:    reference,
		},
	}

	resourceData := regulationResourceData{}
	var tfResource terraformResourceData = &resourceData
	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddObj)

	assert.Equal(t, resourceData.Reference.IsNull(), false)
	assert.Equal(t, resourceData.Reference.ValueString(), "")
}

// TestRegulationResource_defectdojoResource_ReferenceEmpty exercises the
// write path when Reference is configured as an explicit empty string, the
// other arm of the union a null round-trip does not cover.
func TestRegulationResource_defectdojoResource_ReferenceEmpty(t *testing.T) {
	resourceData := regulationResourceData{
		Category:     types.StringValue("other"),
		Jurisdiction: types.StringValue("US"),
		Reference:    types.StringValue(""),
	}

	ddResource := resourceData.defectdojoResource()
	var tfResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, tfResource, &ddResource)

	ddObj := ddResource.(*regulationDefectdojoResource)
	assert.Assert(t, ddObj.Reference != nil)

	raw, err := ddObj.Reference.MarshalJSON()
	assert.NilError(t, err)
	assert.Equal(t, string(raw), `""`)
}

// TestRegulationResource_defectdojoResource_ReferenceNull exercises the write
// path when Reference is left null in configuration: the wrapper pointer
// must stay nil, or the request would send a value the practitioner never
// configured.
func TestRegulationResource_defectdojoResource_ReferenceNull(t *testing.T) {
	resourceData := regulationResourceData{
		Category:     types.StringValue("other"),
		Jurisdiction: types.StringValue("US"),
		Reference:    types.StringNull(),
	}

	ddResource := resourceData.defectdojoResource()
	var tfResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, tfResource, &ddResource)

	ddObj := ddResource.(*regulationDefectdojoResource)
	assert.Assert(t, ddObj.Reference == nil)
}

func TestRegulationResource_defectdojoResource(t *testing.T) {
	expectedName := "Test Regulation"
	expectedAcronym := "TST"
	expectedCategory := "other"
	expectedJurisdiction := "US"
	expectedReference := "https://example.com"

	resourceData := regulationResourceData{
		Name:         types.StringValue(expectedName),
		Acronym:      types.StringValue(expectedAcronym),
		Category:     types.StringValue(expectedCategory),
		Jurisdiction: types.StringValue(expectedJurisdiction),
		Reference:    types.StringValue(expectedReference),
	}

	ddResource := resourceData.defectdojoResource()
	var tfResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, tfResource, &ddResource)

	ddObj := ddResource.(*regulationDefectdojoResource)
	assert.Equal(t, ddObj.Name, expectedName)
	assert.Equal(t, ddObj.Acronym, expectedAcronym)
	assert.Equal(t, string(ddObj.Category), expectedCategory)
	assert.Equal(t, ddObj.Jurisdiction, expectedJurisdiction)

	// Reference is an oapi-codegen oneOf-string wrapper as of DefectDojo 3.2
	// (see isOapiUnionStringType in resource.go); confirm the write path sets
	// it correctly, then confirm regulationToRequest can carry it across to
	// the distinct RegulationRequest_Reference wrapper type.
	gotReference, err := ddObj.Reference.AsRegulationReference0()
	assert.NilError(t, err)
	assert.Equal(t, gotReference, expectedReference)

	req := regulationToRequest(ddObj.Regulation)
	gotReqReference, err := req.Reference.AsRegulationRequestReference0()
	assert.NilError(t, err)
	assert.Equal(t, gotReqReference, expectedReference)
}
