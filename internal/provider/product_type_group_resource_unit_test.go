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

func TestProductTypeGroupResourcePopulate(t *testing.T) {
	expectedId := 11
	expectedProductType := 2
	expectedGroup := 10
	expectedRole := 3

	ddResource := productTypeGroupDefectdojoResource{
		ProductTypeGroup: dd.ProductTypeGroup{
			Id:          &expectedId,
			ProductType: expectedProductType,
			Group:       expectedGroup,
			Role:        expectedRole,
		},
	}

	resourceData := productTypeGroupResourceData{}
	var terraformResource terraformResourceData = &resourceData

	populateResourceData(context.Background(), &diag.Diagnostics{}, &terraformResource, &ddResource)
	assert.Equal(t, resourceData.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, resourceData.ProductType.ValueInt64(), int64(expectedProductType))
	assert.Equal(t, resourceData.Group.ValueInt64(), int64(expectedGroup))
	assert.Equal(t, resourceData.Role.ValueInt64(), int64(expectedRole))
}

func TestProductTypeGroupResource__defectdojoResource(t *testing.T) {
	expectedProductType := 2
	expectedGroup := 10
	expectedRole := 3

	resourceData := productTypeGroupResourceData{
		ProductType: types.Int64Value(int64(expectedProductType)),
		Group:       types.Int64Value(int64(expectedGroup)),
		Role:        types.Int64Value(int64(expectedRole)),
	}

	ddRes := resourceData.defectdojoResource()
	ddGroup := ddRes.(*productTypeGroupDefectdojoResource)
	var terraformResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, terraformResource, &ddRes)

	assert.Equal(t, ddGroup.ProductType, expectedProductType)
	assert.Equal(t, ddGroup.Group, expectedGroup)
	assert.Equal(t, ddGroup.Role, expectedRole)
}
