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

func TestProductGroupResourcePopulate(t *testing.T) {
	expectedId := 8
	expectedProduct := 1
	expectedGroup := 10
	expectedRole := 3

	ddResource := productGroupDefectdojoResource{
		ProductGroup: dd.ProductGroup{
			Id:      &expectedId,
			Product: expectedProduct,
			Group:   expectedGroup,
			Role:    expectedRole,
		},
	}

	resourceData := productGroupResourceData{}
	var terraformResource terraformResourceData = &resourceData

	populateResourceData(context.Background(), &diag.Diagnostics{}, &terraformResource, &ddResource)
	assert.Equal(t, resourceData.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, resourceData.Product.ValueInt64(), int64(expectedProduct))
	assert.Equal(t, resourceData.Group.ValueInt64(), int64(expectedGroup))
	assert.Equal(t, resourceData.Role.ValueInt64(), int64(expectedRole))
}

func TestProductGroupResource__defectdojoResource(t *testing.T) {
	expectedProduct := 1
	expectedGroup := 10
	expectedRole := 3

	resourceData := productGroupResourceData{
		Product: types.Int64Value(int64(expectedProduct)),
		Group:   types.Int64Value(int64(expectedGroup)),
		Role:    types.Int64Value(int64(expectedRole)),
	}

	ddRes := resourceData.defectdojoResource()
	ddGroup := ddRes.(*productGroupDefectdojoResource)
	var terraformResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, terraformResource, &ddRes)

	assert.Equal(t, ddGroup.Product, expectedProduct)
	assert.Equal(t, ddGroup.Group, expectedGroup)
	assert.Equal(t, ddGroup.Role, expectedRole)
}
