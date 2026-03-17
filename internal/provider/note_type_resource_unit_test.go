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

func TestNoteTypeResourcePopulate(t *testing.T) {
	expectedId := 99
	expectedName := "Test Note Type"
	expectedDescription := "A test note type"
	expectedIsActive := true
	expectedIsMandatory := false
	expectedIsSingle := false

	ddObj := noteTypeDefectdojoResource{
		NoteType: dd.NoteType{
			Id:          &expectedId,
			Name:        expectedName,
			Description: expectedDescription,
			IsActive:    &expectedIsActive,
			IsMandatory: &expectedIsMandatory,
			IsSingle:    &expectedIsSingle,
		},
	}

	resourceData := noteTypeResourceData{}
	var tfResource terraformResourceData = &resourceData
	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddObj)

	assert.Equal(t, resourceData.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, resourceData.Name.ValueString(), expectedName)
	assert.Equal(t, resourceData.Description.ValueString(), expectedDescription)
	assert.Equal(t, resourceData.IsActive.ValueBool(), expectedIsActive)
	assert.Equal(t, resourceData.IsMandatory.ValueBool(), expectedIsMandatory)
	assert.Equal(t, resourceData.IsSingle.ValueBool(), expectedIsSingle)
}

func TestNoteTypeResource_defectdojoResource(t *testing.T) {
	expectedName := "Test Note Type"
	expectedDescription := "A test note type"
	expectedIsActive := true

	resourceData := noteTypeResourceData{
		Name:        types.StringValue(expectedName),
		Description: types.StringValue(expectedDescription),
		IsActive:    types.BoolValue(expectedIsActive),
		IsMandatory: types.BoolValue(false),
		IsSingle:    types.BoolValue(false),
	}

	ddResource := resourceData.defectdojoResource()
	var tfResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, tfResource, &ddResource)

	ddObj := ddResource.(*noteTypeDefectdojoResource)
	assert.Equal(t, ddObj.Name, expectedName)
	assert.Equal(t, ddObj.Description, expectedDescription)
	assert.Equal(t, *ddObj.IsActive, expectedIsActive)
}
