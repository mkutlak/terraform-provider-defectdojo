package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"gotest.tools/assert"
)

func TestUserResourcePopulate(t *testing.T) {
	expectedId := 99
	expectedUsername := "testuser"
	expectedEmail := openapi_types.Email("test@example.com")
	expectedFirstName := "Test"
	expectedLastName := "User"
	expectedIsActive := true
	expectedIsSuperuser := false

	ddUser := userDefectdojoResource{
		User: dd.User{
			Id:          &expectedId,
			Username:    expectedUsername,
			Email:       expectedEmail,
			FirstName:   &expectedFirstName,
			LastName:    &expectedLastName,
			IsActive:    &expectedIsActive,
			IsSuperuser: &expectedIsSuperuser,
		},
	}

	userResource := userResourceData{}
	var terraformResource terraformResourceData = &userResource

	populateResourceData(context.Background(), &diag.Diagnostics{}, &terraformResource, &ddUser)
	assert.Equal(t, userResource.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, userResource.Username.ValueString(), expectedUsername)
	assert.Equal(t, userResource.Email.ValueString(), string(expectedEmail))
	assert.Equal(t, userResource.FirstName.ValueString(), expectedFirstName)
	assert.Equal(t, userResource.LastName.ValueString(), expectedLastName)
	assert.Equal(t, userResource.IsActive.ValueBool(), expectedIsActive)
	assert.Equal(t, userResource.IsSuperuser.ValueBool(), expectedIsSuperuser)
}

func TestUserResourcePopulateNils(t *testing.T) {
	ddUser := userDefectdojoResource{
		User: dd.User{},
	}

	userResource := userResourceData{}
	var terraformResource terraformResourceData = &userResource

	populateResourceData(context.Background(), &diag.Diagnostics{}, &terraformResource, &ddUser)

	assert.Equal(t, userResource.Username.ValueString(), "")
	assert.Equal(t, userResource.FirstName.IsNull(), true)
	assert.Equal(t, userResource.LastName.IsNull(), true)
	assert.Equal(t, userResource.IsActive.IsNull(), true)
	assert.Equal(t, userResource.IsSuperuser.IsNull(), true)
}

func TestUserResource__defectdojoResource(t *testing.T) {
	expectedUsername := "testuser"
	expectedEmail := "test@example.com"
	expectedFirstName := "Test"
	expectedLastName := "User"
	expectedIsActive := true
	expectedIsSuperuser := false

	userResource := userResourceData{
		Username:    types.StringValue(expectedUsername),
		Email:       types.StringValue(expectedEmail),
		FirstName:   types.StringValue(expectedFirstName),
		LastName:    types.StringValue(expectedLastName),
		IsActive:    types.BoolValue(expectedIsActive),
		IsSuperuser: types.BoolValue(expectedIsSuperuser),
	}

	ddResource := userResource.defectdojoResource()
	ddUser := ddResource.(*userDefectdojoResource)
	var terraformResource terraformResourceData = &userResource
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, terraformResource, &ddResource)

	assert.Equal(t, ddUser.Username, expectedUsername)
	assert.Equal(t, string(ddUser.Email), expectedEmail)
	assert.Equal(t, *ddUser.FirstName, expectedFirstName)
	assert.Equal(t, *ddUser.LastName, expectedLastName)
	assert.Equal(t, *ddUser.IsActive, expectedIsActive)
	assert.Equal(t, *ddUser.IsSuperuser, expectedIsSuperuser)
}
