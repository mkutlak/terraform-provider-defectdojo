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

func TestEndpointDataSourcePopulate(t *testing.T) {
	expectedId := 10
	expectedHost := "example.com"
	expectedProtocol := "https"
	expectedUserinfo := "alice"
	expectedPort := 443
	expectedPath := "api/v1"
	expectedQuery := "limit=10"
	expectedFragment := "section"
	expectedProduct := 5

	ddObj := endpointDefectdojoResource{
		V3EndpointCompatible: dd.V3EndpointCompatible{
			Id:       &expectedId,
			Host:     expectedHost,
			Protocol: expectedProtocol,
			Userinfo: expectedUserinfo,
			Port:     expectedPort,
			Path:     expectedPath,
			Query:    expectedQuery,
			Fragment: expectedFragment,
			Product:  expectedProduct,
		},
	}

	resourceData := endpointResourceData{}
	var tfResource terraformResourceData = &resourceData
	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddObj)

	assert.Equal(t, resourceData.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, resourceData.Host.ValueString(), expectedHost)
	assert.Equal(t, resourceData.Protocol.ValueString(), expectedProtocol)
	assert.Equal(t, resourceData.Userinfo.ValueString(), expectedUserinfo)
	assert.Equal(t, resourceData.Port.ValueInt64(), int64(expectedPort))
	assert.Equal(t, resourceData.Path.ValueString(), expectedPath)
	assert.Equal(t, resourceData.Query.ValueString(), expectedQuery)
	assert.Equal(t, resourceData.Fragment.ValueString(), expectedFragment)
	assert.Equal(t, resourceData.Product.ValueInt64(), int64(expectedProduct))
}

func TestEndpointDataSource_defectdojoResource(t *testing.T) {
	expectedHost := "example.com"
	expectedProtocol := "https"
	expectedPort := 8443

	resourceData := endpointResourceData{
		Host:     types.StringValue(expectedHost),
		Protocol: types.StringValue(expectedProtocol),
		Port:     types.Int64Value(int64(expectedPort)),
	}

	ddResource := resourceData.defectdojoResource()
	var tfResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, tfResource, &ddResource)

	ddObj := ddResource.(*endpointDefectdojoResource)
	assert.Equal(t, ddObj.Host, expectedHost)
	assert.Equal(t, ddObj.Protocol, expectedProtocol)
	assert.Equal(t, ddObj.Port, expectedPort)
}

func TestEndpointDataSource_writeStubsAreReadOnly(t *testing.T) {
	ddObj := &endpointDefectdojoResource{}

	status, body, err := ddObj.createApiCall(context.Background(), nil)
	assert.Equal(t, status, 0)
	assert.Assert(t, body == nil)
	assert.ErrorContains(t, err, "read-only projection in DefectDojo 3.x")

	status, body, err = ddObj.updateApiCall(context.Background(), nil, 1)
	assert.Equal(t, status, 0)
	assert.Assert(t, body == nil)
	assert.ErrorContains(t, err, "read-only projection in DefectDojo 3.x")

	status, body, err = ddObj.deleteApiCall(context.Background(), nil, 1)
	assert.Equal(t, status, 0)
	assert.Assert(t, body == nil)
	assert.ErrorContains(t, err, "read-only projection in DefectDojo 3.x")
}
