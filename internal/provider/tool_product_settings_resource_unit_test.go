package provider

import (
	"context"
	"testing"

	dd "github.com/doximity/terraform-provider-defectdojo/internal/ddclient"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gotest.tools/assert"
)

func TestToolProductSettingsResourcePopulate(t *testing.T) {
	expectedId := 5
	expectedName := "My Tool Settings"
	expectedSettingUrl := "https://tool.example.com/settings"
	expectedProduct := 1
	expectedToolConfig := 2

	ddResource := toolProductSettingsDefectdojoResource{
		ToolProductSettings: dd.ToolProductSettings{
			Id:                &expectedId,
			Name:              expectedName,
			SettingUrl:        expectedSettingUrl,
			Product:           expectedProduct,
			ToolConfiguration: expectedToolConfig,
		},
	}

	resourceData := toolProductSettingsResourceData{}
	var trd terraformResourceData = &resourceData

	populateResourceData(context.Background(), &diag.Diagnostics{}, &trd, &ddResource)

	assert.Equal(t, resourceData.Id.ValueString(), "5")
	assert.Equal(t, resourceData.Name.ValueString(), expectedName)
	assert.Equal(t, resourceData.SettingUrl.ValueString(), expectedSettingUrl)
	assert.Equal(t, resourceData.Product.ValueInt64(), int64(expectedProduct))
	assert.Equal(t, resourceData.ToolConfiguration.ValueInt64(), int64(expectedToolConfig))
}

func TestToolProductSettingsResourcePopulateDefectdojo(t *testing.T) {
	resourceData := toolProductSettingsResourceData{
		Name:              types.StringValue("My Tool Settings"),
		SettingUrl:        types.StringValue("https://tool.example.com/settings"),
		Product:           types.Int64Value(1),
		ToolConfiguration: types.Int64Value(2),
	}

	ddRes := resourceData.defectdojoResource()
	ddTool := ddRes.(*toolProductSettingsDefectdojoResource)
	var trd terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, trd, &ddRes)

	assert.Equal(t, ddTool.Name, "My Tool Settings")
	assert.Equal(t, ddTool.SettingUrl, "https://tool.example.com/settings")
	assert.Equal(t, ddTool.Product, 1)
	assert.Equal(t, ddTool.ToolConfiguration, 2)
}
