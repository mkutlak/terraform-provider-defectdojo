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

func TestGlobalRoleResourcePopulate(t *testing.T) {
	expectedId := 7
	expectedUser := 42
	expectedRole := 5

	ddResource := globalRoleDefectdojoResource{
		GlobalRole: dd.GlobalRole{
			Id:   &expectedId,
			User: &expectedUser,
			Role: &expectedRole,
		},
	}

	resourceData := globalRoleResourceData{}
	var terraformResource terraformResourceData = &resourceData

	populateResourceData(context.Background(), &diag.Diagnostics{}, &terraformResource, &ddResource)
	assert.Equal(t, resourceData.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, resourceData.User.ValueInt64(), int64(expectedUser))
	assert.Equal(t, resourceData.Role.ValueInt64(), int64(expectedRole))
	assert.Equal(t, resourceData.Group.IsNull(), true)
}

func TestGlobalRoleResource__defectdojoResource(t *testing.T) {
	expectedUser := 42
	expectedRole := 5

	resourceData := globalRoleResourceData{
		User: types.Int64Value(int64(expectedUser)),
		Role: types.Int64Value(int64(expectedRole)),
	}

	ddRes := resourceData.defectdojoResource()
	ddRole := ddRes.(*globalRoleDefectdojoResource)
	var terraformResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, terraformResource, &ddRes)

	assert.Equal(t, *ddRole.User, expectedUser)
	assert.Equal(t, *ddRole.Role, expectedRole)
}
