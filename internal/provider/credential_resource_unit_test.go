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

func TestCredentialResourcePopulate(t *testing.T) {
	expectedId := 12
	expectedName := "Test Credential"
	expectedEnvironment := 1
	expectedUsername := "creduser"
	expectedRole := "viewer"
	expectedUrl := "https://example.com"
	expectedDescription := "A test credential"
	expectedIsValid := true

	ddResource := credentialDefectdojoResource{
		Credential: dd.Credential{
			Id:          &expectedId,
			Name:        expectedName,
			Environment: expectedEnvironment,
			Username:    expectedUsername,
			Role:        expectedRole,
			Url:         expectedUrl,
			Description: &expectedDescription,
			IsValid:     &expectedIsValid,
		},
	}

	resourceData := credentialResourceData{}
	var terraformResource terraformResourceData = &resourceData

	populateResourceData(context.Background(), &diag.Diagnostics{}, &terraformResource, &ddResource)
	assert.Equal(t, resourceData.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, resourceData.Name.ValueString(), expectedName)
	assert.Equal(t, resourceData.Environment.ValueInt64(), int64(expectedEnvironment))
	assert.Equal(t, resourceData.Username.ValueString(), expectedUsername)
	assert.Equal(t, resourceData.Role.ValueString(), expectedRole)
	assert.Equal(t, resourceData.Url.ValueString(), expectedUrl)
	assert.Equal(t, resourceData.Description.ValueString(), expectedDescription)
	assert.Equal(t, resourceData.IsValid.ValueBool(), expectedIsValid)
	assert.Equal(t, resourceData.Authentication.IsNull(), true)
	assert.Equal(t, resourceData.HttpAuthentication.IsNull(), true)
}

func TestCredentialResourcePopulateNils(t *testing.T) {
	ddResource := credentialDefectdojoResource{
		Credential: dd.Credential{},
	}

	resourceData := credentialResourceData{}
	var terraformResource terraformResourceData = &resourceData

	populateResourceData(context.Background(), &diag.Diagnostics{}, &terraformResource, &ddResource)
	assert.Equal(t, resourceData.Name.ValueString(), "")
	assert.Equal(t, resourceData.Description.IsNull(), true)
	assert.Equal(t, resourceData.Authentication.IsNull(), true)
	assert.Equal(t, resourceData.IsValid.IsNull(), true)
}

func TestCredentialResource__defectdojoResource(t *testing.T) {
	expectedName := "Test Credential"
	expectedEnvironment := 1
	expectedUsername := "creduser"
	expectedRole := "viewer"
	expectedUrl := "https://example.com"

	resourceData := credentialResourceData{
		Name:        types.StringValue(expectedName),
		Environment: types.Int64Value(int64(expectedEnvironment)),
		Username:    types.StringValue(expectedUsername),
		Role:        types.StringValue(expectedRole),
		Url:         types.StringValue(expectedUrl),
	}

	ddRes := resourceData.defectdojoResource()
	ddCred := ddRes.(*credentialDefectdojoResource)
	var terraformResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, terraformResource, &ddRes)

	assert.Equal(t, ddCred.Name, expectedName)
	assert.Equal(t, ddCred.Environment, expectedEnvironment)
	assert.Equal(t, ddCred.Username, expectedUsername)
	assert.Equal(t, ddCred.Role, expectedRole)
	assert.Equal(t, ddCred.Url, expectedUrl)
}
