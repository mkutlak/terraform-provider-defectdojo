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

func TestToolTypeResourcePopulate(t *testing.T) {
	expectedId := 99
	expectedName := "Test Tool Type"
	expectedDescription := "A test tool type"

	ddObj := toolTypeDefectdojoResource{
		ToolType: dd.ToolType{
			Id:          &expectedId,
			Name:        expectedName,
			Description: &expectedDescription,
		},
	}

	resourceData := toolTypeResourceData{}
	var tfResource terraformResourceData = &resourceData
	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddObj)

	assert.Equal(t, resourceData.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, resourceData.Name.ValueString(), expectedName)
	assert.Equal(t, resourceData.Description.ValueString(), expectedDescription)
}

func TestToolTypeResource_defectdojoResource(t *testing.T) {
	expectedName := "Test Tool Type"
	expectedDescription := "A test tool type"

	resourceData := toolTypeResourceData{
		Name:        types.StringValue(expectedName),
		Description: types.StringValue(expectedDescription),
	}

	ddResource := resourceData.defectdojoResource()
	var tfResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, tfResource, &ddResource)

	ddObj := ddResource.(*toolTypeDefectdojoResource)
	assert.Equal(t, ddObj.Name, expectedName)
	assert.Equal(t, *ddObj.Description, expectedDescription)
}
