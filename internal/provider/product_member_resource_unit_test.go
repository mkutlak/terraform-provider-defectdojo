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

func TestProductMemberResourcePopulate(t *testing.T) {
	expectedId := 15
	expectedProduct := 1
	expectedUser := 42
	expectedRole := 3

	ddResource := productMemberDefectdojoResource{
		ProductMember: dd.ProductMember{
			Id:      &expectedId,
			Product: expectedProduct,
			User:    expectedUser,
			Role:    expectedRole,
		},
	}

	resourceData := productMemberResourceData{}
	var terraformResource terraformResourceData = &resourceData

	populateResourceData(context.Background(), &diag.Diagnostics{}, &terraformResource, &ddResource)
	assert.Equal(t, resourceData.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, resourceData.Product.ValueInt64(), int64(expectedProduct))
	assert.Equal(t, resourceData.User.ValueInt64(), int64(expectedUser))
	assert.Equal(t, resourceData.Role.ValueInt64(), int64(expectedRole))
}

func TestProductMemberResource__defectdojoResource(t *testing.T) {
	expectedProduct := 1
	expectedUser := 42
	expectedRole := 3

	resourceData := productMemberResourceData{
		Product: types.Int64Value(int64(expectedProduct)),
		User:    types.Int64Value(int64(expectedUser)),
		Role:    types.Int64Value(int64(expectedRole)),
	}

	ddRes := resourceData.defectdojoResource()
	ddMember := ddRes.(*productMemberDefectdojoResource)
	var terraformResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, terraformResource, &ddRes)

	assert.Equal(t, ddMember.Product, expectedProduct)
	assert.Equal(t, ddMember.User, expectedUser)
	assert.Equal(t, ddMember.Role, expectedRole)
}
