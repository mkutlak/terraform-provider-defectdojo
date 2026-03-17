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

func TestToolConfigurationResourcePopulate(t *testing.T) {
	expectedId := 99
	expectedName := "Test Tool Configuration"
	expectedToolType := 42
	expectedDescription := "A test tool configuration"
	expectedUrl := "https://example.com"

	ddObj := toolConfigurationDefectdojoResource{
		ToolConfiguration: dd.ToolConfiguration{
			Id:          &expectedId,
			Name:        expectedName,
			ToolType:    expectedToolType,
			Description: &expectedDescription,
			Url:         &expectedUrl,
		},
	}

	resourceData := toolConfigurationResourceData{}
	var tfResource terraformResourceData = &resourceData
	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddObj)

	assert.Equal(t, resourceData.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, resourceData.Name.ValueString(), expectedName)
	assert.Equal(t, resourceData.ToolType.ValueInt64(), int64(expectedToolType))
	assert.Equal(t, resourceData.Description.ValueString(), expectedDescription)
	assert.Equal(t, resourceData.Url.ValueString(), expectedUrl)
}

func TestToolConfigurationResource_defectdojoResource(t *testing.T) {
	expectedName := "Test Tool Configuration"
	expectedToolType := int64(42)
	expectedDescription := "A test tool configuration"

	resourceData := toolConfigurationResourceData{
		Name:        types.StringValue(expectedName),
		ToolType:    types.Int64Value(expectedToolType),
		Description: types.StringValue(expectedDescription),
	}

	ddResource := resourceData.defectdojoResource()
	var tfResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, tfResource, &ddResource)

	ddObj := ddResource.(*toolConfigurationDefectdojoResource)
	assert.Equal(t, ddObj.Name, expectedName)
	assert.Equal(t, int64(ddObj.ToolType), expectedToolType)
	assert.Equal(t, *ddObj.Description, expectedDescription)
}
