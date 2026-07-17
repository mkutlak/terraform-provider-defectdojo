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

func TestConfigurationPermissionDataSourcePopulate(t *testing.T) {
	expectedId := 10
	expectedCodename := "add_product"
	expectedName := "Can add product"

	ddObj := configurationPermissionDefectdojoResource{
		ConfigurationPermission: dd.ConfigurationPermission{
			Id:       &expectedId,
			Codename: expectedCodename,
			Name:     expectedName,
		},
	}

	resourceData := configurationPermissionResourceData{}
	var tfResource terraformResourceData = &resourceData
	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddObj)

	assert.Equal(t, resourceData.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, resourceData.Codename.ValueString(), expectedCodename)
	assert.Equal(t, resourceData.Name.ValueString(), expectedName)
}

func TestConfigurationPermissionDataSource_defectdojoResource(t *testing.T) {
	expectedCodename := "add_product"

	resourceData := configurationPermissionResourceData{
		Codename: types.StringValue(expectedCodename),
	}

	ddResource := resourceData.defectdojoResource()
	var tfResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, tfResource, &ddResource)

	ddObj := ddResource.(*configurationPermissionDefectdojoResource)
	assert.Equal(t, ddObj.Codename, expectedCodename)
}

func TestConfigurationPermissionDataSource_writeStubsAreReadOnly(t *testing.T) {
	ddObj := &configurationPermissionDefectdojoResource{}

	status, body, err := ddObj.createApiCall(context.Background(), nil)
	assert.Equal(t, status, 0)
	assert.Assert(t, body == nil)
	assert.ErrorContains(t, err, "read-only")

	status, body, err = ddObj.updateApiCall(context.Background(), nil, 1)
	assert.Equal(t, status, 0)
	assert.Assert(t, body == nil)
	assert.ErrorContains(t, err, "read-only")

	status, body, err = ddObj.deleteApiCall(context.Background(), nil, 1)
	assert.Equal(t, status, 0)
	assert.Assert(t, body == nil)
	assert.ErrorContains(t, err, "read-only")
}

func TestConfigurationPermissionDataSource_nameFromData(t *testing.T) {
	provider := configurationPermissionDataProvider{}

	withCodename := &configurationPermissionResourceData{
		Codename: types.StringValue("add_product"),
	}
	name, ok := provider.nameFromData(withCodename)
	assert.Assert(t, ok)
	assert.Equal(t, name, "add_product")

	withoutCodename := &configurationPermissionResourceData{}
	_, ok = provider.nameFromData(withoutCodename)
	assert.Assert(t, !ok)
}
