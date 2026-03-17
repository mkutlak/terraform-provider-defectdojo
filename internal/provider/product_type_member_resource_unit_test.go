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

func TestProductTypeMemberResourcePopulate(t *testing.T) {
	expectedId := 9
	expectedProductType := 2
	expectedUser := 42
	expectedRole := 3

	ddResource := productTypeMemberDefectdojoResource{
		ProductTypeMember: dd.ProductTypeMember{
			Id:          &expectedId,
			ProductType: expectedProductType,
			User:        expectedUser,
			Role:        expectedRole,
		},
	}

	resourceData := productTypeMemberResourceData{}
	var terraformResource terraformResourceData = &resourceData

	populateResourceData(context.Background(), &diag.Diagnostics{}, &terraformResource, &ddResource)
	assert.Equal(t, resourceData.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, resourceData.ProductType.ValueInt64(), int64(expectedProductType))
	assert.Equal(t, resourceData.User.ValueInt64(), int64(expectedUser))
	assert.Equal(t, resourceData.Role.ValueInt64(), int64(expectedRole))
}

func TestProductTypeMemberResource__defectdojoResource(t *testing.T) {
	expectedProductType := 2
	expectedUser := 42
	expectedRole := 3

	resourceData := productTypeMemberResourceData{
		ProductType: types.Int64Value(int64(expectedProductType)),
		User:        types.Int64Value(int64(expectedUser)),
		Role:        types.Int64Value(int64(expectedRole)),
	}

	ddRes := resourceData.defectdojoResource()
	ddMember := ddRes.(*productTypeMemberDefectdojoResource)
	var terraformResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, terraformResource, &ddRes)

	assert.Equal(t, ddMember.ProductType, expectedProductType)
	assert.Equal(t, ddMember.User, expectedUser)
	assert.Equal(t, ddMember.Role, expectedRole)
}
