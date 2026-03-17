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

func TestTechnologyResourcePopulate(t *testing.T) {
	expectedId := 50
	expectedName := "React"
	expectedProduct := 3
	expectedUser := 1
	expectedVersion := "18.0"
	expectedWebsite := "https://reactjs.org"
	expectedConfidence := 90

	ddObj := technologyDefectdojoResource{
		AppAnalysis: dd.AppAnalysis{
			Id:         &expectedId,
			Name:       expectedName,
			Product:    expectedProduct,
			User:       expectedUser,
			Version:    &expectedVersion,
			Website:    &expectedWebsite,
			Confidence: &expectedConfidence,
		},
	}

	resourceData := technologyResourceData{}
	var tfResource terraformResourceData = &resourceData
	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddObj)

	assert.Equal(t, resourceData.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, resourceData.Name.ValueString(), expectedName)
	assert.Equal(t, resourceData.Product.ValueInt64(), int64(expectedProduct))
	assert.Equal(t, resourceData.User.ValueInt64(), int64(expectedUser))
	assert.Equal(t, resourceData.Version.ValueString(), expectedVersion)
	assert.Equal(t, resourceData.Website.ValueString(), expectedWebsite)
	assert.Equal(t, resourceData.Confidence.ValueInt64(), int64(expectedConfidence))
}

func TestTechnologyResource_defectdojoResource(t *testing.T) {
	expectedName := "React"
	expectedProduct := int64(3)
	expectedUser := int64(1)

	resourceData := technologyResourceData{
		Name:    types.StringValue(expectedName),
		Product: types.Int64Value(expectedProduct),
		User:    types.Int64Value(expectedUser),
	}

	ddResource := resourceData.defectdojoResource()
	var tfResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, tfResource, &ddResource)

	ddObj := ddResource.(*technologyDefectdojoResource)
	assert.Equal(t, ddObj.Name, expectedName)
	assert.Equal(t, int64(ddObj.Product), expectedProduct)
	assert.Equal(t, int64(ddObj.User), expectedUser)
}
