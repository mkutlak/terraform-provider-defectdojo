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

func TestDevelopmentEnvironmentResourcePopulate(t *testing.T) {
	expectedId := 99
	expectedName := "Test Development Environment"

	ddObj := developmentEnvironmentDefectdojoResource{
		DevelopmentEnvironment: dd.DevelopmentEnvironment{
			Id:   &expectedId,
			Name: expectedName,
		},
	}

	resourceData := developmentEnvironmentResourceData{}
	var tfResource terraformResourceData = &resourceData
	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddObj)

	assert.Equal(t, resourceData.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, resourceData.Name.ValueString(), expectedName)
}

func TestDevelopmentEnvironmentResource_defectdojoResource(t *testing.T) {
	expectedName := "Test Development Environment"

	resourceData := developmentEnvironmentResourceData{
		Name: types.StringValue(expectedName),
	}

	ddResource := resourceData.defectdojoResource()
	var tfResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, tfResource, &ddResource)

	ddObj := ddResource.(*developmentEnvironmentDefectdojoResource)
	assert.Equal(t, ddObj.Name, expectedName)
}
