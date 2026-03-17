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

func TestSlaConfigurationResourcePopulate(t *testing.T) {
	expectedId := 99
	expectedName := "Test SLA Configuration"
	expectedCritical := 7
	expectedHigh := 30
	expectedMedium := 90
	expectedLow := 180
	expectedEnforceCritical := true
	expectedEnforceHigh := true

	ddObj := slaConfigurationDefectdojoResource{
		SLAConfiguration: dd.SLAConfiguration{
			Id:              &expectedId,
			Name:            expectedName,
			Critical:        &expectedCritical,
			High:            &expectedHigh,
			Medium:          &expectedMedium,
			Low:             &expectedLow,
			EnforceCritical: &expectedEnforceCritical,
			EnforceHigh:     &expectedEnforceHigh,
		},
	}

	resourceData := slaConfigurationResourceData{}
	var tfResource terraformResourceData = &resourceData
	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddObj)

	assert.Equal(t, resourceData.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, resourceData.Name.ValueString(), expectedName)
	assert.Equal(t, resourceData.Critical.ValueInt64(), int64(expectedCritical))
	assert.Equal(t, resourceData.High.ValueInt64(), int64(expectedHigh))
	assert.Equal(t, resourceData.Medium.ValueInt64(), int64(expectedMedium))
	assert.Equal(t, resourceData.Low.ValueInt64(), int64(expectedLow))
	assert.Equal(t, resourceData.EnforceCritical.ValueBool(), expectedEnforceCritical)
	assert.Equal(t, resourceData.EnforceHigh.ValueBool(), expectedEnforceHigh)
}

func TestSlaConfigurationResource_defectdojoResource(t *testing.T) {
	expectedName := "Test SLA Configuration"
	expectedCritical := int64(7)
	expectedHigh := int64(30)

	resourceData := slaConfigurationResourceData{
		Name:     types.StringValue(expectedName),
		Critical: types.Int64Value(expectedCritical),
		High:     types.Int64Value(expectedHigh),
	}

	ddResource := resourceData.defectdojoResource()
	var tfResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, tfResource, &ddResource)

	ddObj := ddResource.(*slaConfigurationDefectdojoResource)
	assert.Equal(t, ddObj.Name, expectedName)
	assert.Equal(t, int64(*ddObj.Critical), expectedCritical)
	assert.Equal(t, int64(*ddObj.High), expectedHigh)
}
