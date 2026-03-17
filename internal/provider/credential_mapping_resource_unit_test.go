package provider

import (
	"context"
	"testing"

	dd "github.com/doximity/terraform-provider-defectdojo/internal/ddclient"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gotest.tools/assert"
)

func TestCredentialMappingResourcePopulate(t *testing.T) {
	expectedId := 8
	expectedCredId := 42
	expectedProduct := 1
	expectedIsAuthn := true

	ddResource := credentialMappingDefectdojoResource{
		CredentialMapping: dd.CredentialMapping{
			Id:              &expectedId,
			CredId:          expectedCredId,
			Product:         &expectedProduct,
			IsAuthnProvider: &expectedIsAuthn,
		},
	}

	resourceData := credentialMappingResourceData{}
	var trd terraformResourceData = &resourceData

	populateResourceData(context.Background(), &diag.Diagnostics{}, &trd, &ddResource)

	assert.Equal(t, resourceData.Id.ValueString(), "8")
	assert.Equal(t, resourceData.CredId.ValueInt64(), int64(expectedCredId))
	assert.Equal(t, resourceData.Product.ValueInt64(), int64(expectedProduct))
	assert.Equal(t, resourceData.IsAuthnProvider.ValueBool(), expectedIsAuthn)
}

func TestCredentialMappingResourcePopulateDefectdojo(t *testing.T) {
	resourceData := credentialMappingResourceData{
		CredId:          types.Int64Value(42),
		Product:         types.Int64Value(1),
		IsAuthnProvider: types.BoolValue(true),
	}

	ddRes := resourceData.defectdojoResource()
	ddCred := ddRes.(*credentialMappingDefectdojoResource)
	var trd terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, trd, &ddRes)

	assert.Equal(t, ddCred.CredId, 42)
	assert.Equal(t, *ddCred.Product, 1)
	assert.Equal(t, *ddCred.IsAuthnProvider, true)
}
