package provider

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"gotest.tools/assert"
)

func TestUserProfileDataSourcePopulate(t *testing.T) {
	expectedId := 7
	expectedUsername := "admin"
	expectedEmail := "admin@example.com"
	expectedFirstName := "Ada"
	expectedLastName := "Min"
	expectedIsActive := true
	expectedIsStaff := true
	expectedIsSuperuser := true
	expectedDateJoined := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	expectedLastLogin := time.Date(2024, 6, 7, 8, 9, 10, 0, time.UTC)

	ddUserProfile := userProfileDefectdojoResource{
		User: dd.User{
			Id:          &expectedId,
			Username:    expectedUsername,
			Email:       openapi_types.Email(expectedEmail),
			FirstName:   &expectedFirstName,
			LastName:    &expectedLastName,
			IsActive:    &expectedIsActive,
			IsStaff:     &expectedIsStaff,
			IsSuperuser: &expectedIsSuperuser,
			DateJoined:  &expectedDateJoined,
			LastLogin:   &expectedLastLogin,
		},
	}

	profileResource := userProfileResourceData{}
	var tfResource terraformResourceData = &profileResource
	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddUserProfile)

	assert.Equal(t, profileResource.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, profileResource.Username.ValueString(), expectedUsername)
	assert.Equal(t, profileResource.Email.ValueString(), expectedEmail)
	assert.Equal(t, profileResource.FirstName.ValueString(), expectedFirstName)
	assert.Equal(t, profileResource.LastName.ValueString(), expectedLastName)
	assert.Equal(t, profileResource.IsActive.ValueBool(), expectedIsActive)
	assert.Equal(t, profileResource.IsStaff.ValueBool(), expectedIsStaff)
	assert.Equal(t, profileResource.IsSuperuser.ValueBool(), expectedIsSuperuser)
	assert.Equal(t, profileResource.DateJoined.ValueString(), expectedDateJoined.Format(time.RFC3339))
	assert.Equal(t, profileResource.LastLogin.ValueString(), expectedLastLogin.Format(time.RFC3339))
}

func TestUserProfileDataSourcePopulateNils(t *testing.T) {
	profileResource := userProfileResourceData{}
	var tfResource terraformResourceData = &profileResource

	ddUserProfile := userProfileDefectdojoResource{
		User: dd.User{
			Username: "admin",
			Email:    openapi_types.Email("admin@example.com"),
		},
	}
	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddUserProfile)

	assert.Equal(t, profileResource.Id.IsNull(), true)
	assert.Equal(t, profileResource.Username.ValueString(), "admin")
	assert.Equal(t, profileResource.Email.ValueString(), "admin@example.com")
	assert.Equal(t, profileResource.FirstName.IsNull(), true)
	assert.Equal(t, profileResource.LastName.IsNull(), true)
	assert.Equal(t, profileResource.IsActive.IsNull(), true)
	assert.Equal(t, profileResource.IsStaff.IsNull(), true)
	assert.Equal(t, profileResource.IsSuperuser.IsNull(), true)
	assert.Equal(t, profileResource.DateJoined.IsNull(), true)
	assert.Equal(t, profileResource.LastLogin.IsNull(), true)
}

func TestUserProfileDataSource_defectdojoResource(t *testing.T) {
	profileResource := userProfileResourceData{}
	ddResource := profileResource.defectdojoResource()
	ddUserProfile, ok := ddResource.(*userProfileDefectdojoResource)
	assert.Assert(t, ok)
	assert.Equal(t, ddUserProfile.Username, "")
}

func TestUserProfileDataSource_writeStubsAreReadOnly(t *testing.T) {
	ddObj := &userProfileDefectdojoResource{}

	status, body, err := ddObj.createApiCall(context.Background(), nil)
	assert.Equal(t, status, 0)
	assert.Assert(t, body == nil)
	assert.ErrorContains(t, err, "read-only data source describing the authenticated user")

	status, body, err = ddObj.updateApiCall(context.Background(), nil, 1)
	assert.Equal(t, status, 0)
	assert.Assert(t, body == nil)
	assert.ErrorContains(t, err, "read-only data source describing the authenticated user")

	status, body, err = ddObj.deleteApiCall(context.Background(), nil, 1)
	assert.Equal(t, status, 0)
	assert.Assert(t, body == nil)
	assert.ErrorContains(t, err, "read-only data source describing the authenticated user")
}
