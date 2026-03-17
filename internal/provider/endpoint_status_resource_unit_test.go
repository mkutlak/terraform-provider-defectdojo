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

func TestEndpointStatusResourcePopulate(t *testing.T) {
	expectedId := 20
	expectedEndpoint := 5
	expectedFinding := 10
	expectedFalsePositive := false
	expectedMitigated := true
	expectedOutOfScope := false
	expectedRiskAccepted := false

	ddObj := endpointStatusDefectdojoResource{
		EndpointStatus: dd.EndpointStatus{
			Id:            &expectedId,
			Endpoint:      expectedEndpoint,
			Finding:       expectedFinding,
			FalsePositive: &expectedFalsePositive,
			Mitigated:     &expectedMitigated,
			OutOfScope:    &expectedOutOfScope,
			RiskAccepted:  &expectedRiskAccepted,
		},
	}

	resourceData := endpointStatusResourceData{}
	var tfResource terraformResourceData = &resourceData
	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddObj)

	assert.Equal(t, resourceData.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, resourceData.Endpoint.ValueInt64(), int64(expectedEndpoint))
	assert.Equal(t, resourceData.Finding.ValueInt64(), int64(expectedFinding))
	assert.Equal(t, resourceData.FalsePositive.ValueBool(), expectedFalsePositive)
	assert.Equal(t, resourceData.Mitigated.ValueBool(), expectedMitigated)
	assert.Equal(t, resourceData.OutOfScope.ValueBool(), expectedOutOfScope)
	assert.Equal(t, resourceData.RiskAccepted.ValueBool(), expectedRiskAccepted)
}

func TestEndpointStatusResource_defectdojoResource(t *testing.T) {
	expectedEndpoint := int64(5)
	expectedFinding := int64(10)

	resourceData := endpointStatusResourceData{
		Endpoint: types.Int64Value(expectedEndpoint),
		Finding:  types.Int64Value(expectedFinding),
	}

	ddResource := resourceData.defectdojoResource()
	var tfResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, tfResource, &ddResource)

	ddObj := ddResource.(*endpointStatusDefectdojoResource)
	assert.Equal(t, int64(ddObj.Endpoint), expectedEndpoint)
	assert.Equal(t, int64(ddObj.Finding), expectedFinding)
}
