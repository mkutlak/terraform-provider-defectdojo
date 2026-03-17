package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
	"gotest.tools/assert"
)

func TestProductAPIScanConfigurationResourcePopulate(t *testing.T) {
	expectedId := 7
	expectedProduct := 1
	expectedToolConfig := 3
	expectedKey1 := "key-one"

	ddResource := productAPIScanConfigurationDefectdojoResource{
		ProductAPIScanConfiguration: dd.ProductAPIScanConfiguration{
			Id:                &expectedId,
			Product:           expectedProduct,
			ToolConfiguration: expectedToolConfig,
			ServiceKey1:       &expectedKey1,
		},
	}

	resourceData := productAPIScanConfigurationResourceData{}
	var trd terraformResourceData = &resourceData

	populateResourceData(context.Background(), &diag.Diagnostics{}, &trd, &ddResource)

	assert.Equal(t, resourceData.Id.ValueString(), "7")
	assert.Equal(t, resourceData.Product.ValueInt64(), int64(expectedProduct))
	assert.Equal(t, resourceData.ToolConfiguration.ValueInt64(), int64(expectedToolConfig))
	assert.Equal(t, resourceData.ServiceKey1.ValueString(), expectedKey1)
}

func TestProductAPIScanConfigurationResourcePopulateDefectdojo(t *testing.T) {
	resourceData := productAPIScanConfigurationResourceData{
		Product:           types.Int64Value(1),
		ToolConfiguration: types.Int64Value(3),
		ServiceKey1:       types.StringValue("key-one"),
	}

	ddRes := resourceData.defectdojoResource()
	ddScan := ddRes.(*productAPIScanConfigurationDefectdojoResource)
	var trd terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, trd, &ddRes)

	assert.Equal(t, ddScan.Product, 1)
	assert.Equal(t, ddScan.ToolConfiguration, 3)
	assert.Equal(t, *ddScan.ServiceKey1, "key-one")
}
