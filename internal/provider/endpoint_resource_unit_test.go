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

func TestEndpointResourcePopulate(t *testing.T) {
	expectedId := 10
	expectedHost := "example.com"
	expectedProtocol := "https"
	expectedPort := 443
	expectedPath := "api/v1"
	expectedProduct := 5

	ddObj := endpointDefectdojoResource{
		Endpoint: dd.Endpoint{
			Id:       &expectedId,
			Host:     &expectedHost,
			Protocol: &expectedProtocol,
			Port:     &expectedPort,
			Path:     &expectedPath,
			Product:  &expectedProduct,
		},
	}

	resourceData := endpointResourceData{}
	var tfResource terraformResourceData = &resourceData
	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddObj)

	assert.Equal(t, resourceData.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, resourceData.Host.ValueString(), expectedHost)
	assert.Equal(t, resourceData.Protocol.ValueString(), expectedProtocol)
	assert.Equal(t, resourceData.Port.ValueInt64(), int64(expectedPort))
	assert.Equal(t, resourceData.Path.ValueString(), expectedPath)
	assert.Equal(t, resourceData.Product.ValueInt64(), int64(expectedProduct))
}

func TestEndpointResource_defectdojoResource(t *testing.T) {
	expectedHost := "example.com"
	expectedProtocol := "https"

	resourceData := endpointResourceData{
		Host:     types.StringValue(expectedHost),
		Protocol: types.StringValue(expectedProtocol),
	}

	ddResource := resourceData.defectdojoResource()
	var tfResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, tfResource, &ddResource)

	ddObj := ddResource.(*endpointDefectdojoResource)
	assert.Equal(t, *ddObj.Host, expectedHost)
	assert.Equal(t, *ddObj.Protocol, expectedProtocol)
}
