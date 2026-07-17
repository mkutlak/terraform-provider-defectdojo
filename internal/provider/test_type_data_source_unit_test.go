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

func TestTestTypeDataSourcePopulate(t *testing.T) {
	expectedId := 42
	expectedName := "Burp Scan"
	expectedActive := true
	expectedStaticTool := false
	expectedDynamicTool := true

	ddObj := testTypeDefectdojoResource{
		TestType: dd.TestType{
			Id:          &expectedId,
			Name:        &expectedName,
			Active:      &expectedActive,
			StaticTool:  &expectedStaticTool,
			DynamicTool: &expectedDynamicTool,
		},
	}

	resourceData := testTypeResourceData{}
	var tfResource terraformResourceData = &resourceData
	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddObj)

	assert.Equal(t, resourceData.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, resourceData.Name.ValueString(), expectedName)
	assert.Equal(t, resourceData.Active.ValueBool(), expectedActive)
	assert.Equal(t, resourceData.StaticTool.ValueBool(), expectedStaticTool)
	assert.Equal(t, resourceData.DynamicTool.ValueBool(), expectedDynamicTool)
}

func TestTestTypeDataSource_defectdojoResource(t *testing.T) {
	expectedName := "Burp Scan"
	expectedActive := true

	resourceData := testTypeResourceData{
		Name:   types.StringValue(expectedName),
		Active: types.BoolValue(expectedActive),
	}

	ddResource := resourceData.defectdojoResource()
	var tfResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, tfResource, &ddResource)

	ddObj := ddResource.(*testTypeDefectdojoResource)
	assert.Equal(t, *ddObj.Name, expectedName)
	assert.Equal(t, *ddObj.Active, expectedActive)
}

func TestTestTypeDataSource_writeStubsAreReadOnly(t *testing.T) {
	ddObj := &testTypeDefectdojoResource{}

	status, body, err := ddObj.createApiCall(context.Background(), nil)
	assert.Equal(t, status, 0)
	assert.Assert(t, body == nil)
	assert.ErrorContains(t, err, "read-only in this provider")

	status, body, err = ddObj.updateApiCall(context.Background(), nil, 1)
	assert.Equal(t, status, 0)
	assert.Assert(t, body == nil)
	assert.ErrorContains(t, err, "read-only in this provider")

	status, body, err = ddObj.deleteApiCall(context.Background(), nil, 1)
	assert.Equal(t, status, 0)
	assert.Assert(t, body == nil)
	assert.ErrorContains(t, err, "read-only in this provider")
}

func TestTestTypeDataSource_nameFromData(t *testing.T) {
	dp := testTypeDataProvider{}

	withName := &testTypeResourceData{Name: types.StringValue("Burp Scan")}
	name, ok := dp.nameFromData(withName)
	assert.Assert(t, ok)
	assert.Equal(t, name, "Burp Scan")

	withoutName := &testTypeResourceData{Name: types.StringNull()}
	_, ok = dp.nameFromData(withoutName)
	assert.Assert(t, !ok)

	unknownName := &testTypeResourceData{Name: types.StringUnknown()}
	_, ok = dp.nameFromData(unknownName)
	assert.Assert(t, !ok)
}
