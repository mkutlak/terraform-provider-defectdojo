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

func TestDojoGroupMemberResourcePopulate(t *testing.T) {
	expectedId := 5
	expectedGroup := 10
	expectedUser := 20
	expectedRole := 3

	ddResource := dojoGroupMemberDefectdojoResource{
		DojoGroupMember: dd.DojoGroupMember{
			Id:    &expectedId,
			Group: expectedGroup,
			User:  expectedUser,
			Role:  expectedRole,
		},
	}

	resourceData := dojoGroupMemberResourceData{}
	var terraformResource terraformResourceData = &resourceData

	populateResourceData(context.Background(), &diag.Diagnostics{}, &terraformResource, &ddResource)
	assert.Equal(t, resourceData.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, resourceData.Group.ValueInt64(), int64(expectedGroup))
	assert.Equal(t, resourceData.User.ValueInt64(), int64(expectedUser))
	assert.Equal(t, resourceData.Role.ValueInt64(), int64(expectedRole))
}

func TestDojoGroupMemberResource__defectdojoResource(t *testing.T) {
	expectedGroup := 10
	expectedUser := 20
	expectedRole := 3

	resourceData := dojoGroupMemberResourceData{
		Group: types.Int64Value(int64(expectedGroup)),
		User:  types.Int64Value(int64(expectedUser)),
		Role:  types.Int64Value(int64(expectedRole)),
	}

	ddRes := resourceData.defectdojoResource()
	ddMember := ddRes.(*dojoGroupMemberDefectdojoResource)
	var terraformResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, terraformResource, &ddRes)

	assert.Equal(t, ddMember.Group, expectedGroup)
	assert.Equal(t, ddMember.User, expectedUser)
	assert.Equal(t, ddMember.Role, expectedRole)
}
