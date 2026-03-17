package provider

import (
	"context"
	"testing"

	dd "github.com/doximity/terraform-provider-defectdojo/internal/ddclient"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gotest.tools/assert"
)

func TestAssetGroupResourcePopulate(t *testing.T) {
	expectedId := 30
	expectedAsset := 10
	expectedGroup := 20
	expectedRole := 5

	ddResource := assetGroupDefectdojoResource{
		AssetGroup: dd.AssetGroup{
			Id:    &expectedId,
			Asset: expectedAsset,
			Group: expectedGroup,
			Role:  expectedRole,
		},
	}

	resourceData := assetGroupResourceData{}
	var trd terraformResourceData = &resourceData

	populateResourceData(context.Background(), &diag.Diagnostics{}, &trd, &ddResource)

	assert.Equal(t, resourceData.Id.ValueString(), "30")
	assert.Equal(t, resourceData.Asset.ValueInt64(), int64(expectedAsset))
	assert.Equal(t, resourceData.Group.ValueInt64(), int64(expectedGroup))
	assert.Equal(t, resourceData.Role.ValueInt64(), int64(expectedRole))
}

func TestAssetGroupResourcePopulateDefectdojo(t *testing.T) {
	resourceData := assetGroupResourceData{
		Asset: types.Int64Value(10),
		Group: types.Int64Value(20),
		Role:  types.Int64Value(5),
	}

	ddRes := resourceData.defectdojoResource()
	ddGroup := ddRes.(*assetGroupDefectdojoResource)
	var trd terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, trd, &ddRes)

	assert.Equal(t, ddGroup.Asset, 10)
	assert.Equal(t, ddGroup.Group, 20)
	assert.Equal(t, ddGroup.Role, 5)
}
