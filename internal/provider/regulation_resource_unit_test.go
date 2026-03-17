package provider

import (
	"context"
	"fmt"
	"testing"

	dd "github.com/doximity/terraform-provider-defectdojo/internal/ddclient"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

	ddObj := regulationDefectdojoResource{
		Regulation: dd.Regulation{
			Id:           &expectedId,
			Name:         expectedName,
			Acronym:      expectedAcronym,
			Category:     dd.RegulationCategory(expectedCategory),
			Jurisdiction: expectedJurisdiction,
			Description:  &expectedDescription,
			Reference:    &expectedReference,
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

func TestRegulationResource_defectdojoResource(t *testing.T) {
	expectedName := "Test Regulation"
	expectedAcronym := "TST"
	expectedCategory := "other"
	expectedJurisdiction := "US"

	resourceData := regulationResourceData{
		Name:         types.StringValue(expectedName),
		Acronym:      types.StringValue(expectedAcronym),
		Category:     types.StringValue(expectedCategory),
		Jurisdiction: types.StringValue(expectedJurisdiction),
	}

	ddResource := resourceData.defectdojoResource()
	var tfResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, tfResource, &ddResource)

	ddObj := ddResource.(*regulationDefectdojoResource)
	assert.Equal(t, ddObj.Name, expectedName)
	assert.Equal(t, ddObj.Acronym, expectedAcronym)
	assert.Equal(t, string(ddObj.Category), expectedCategory)
	assert.Equal(t, ddObj.Jurisdiction, expectedJurisdiction)
}
