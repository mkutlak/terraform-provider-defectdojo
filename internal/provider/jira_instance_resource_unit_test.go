package provider

import (
	"context"
	"testing"

	dd "github.com/doximity/terraform-provider-defectdojo/internal/ddclient"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gotest.tools/assert"
)

func TestJiraInstanceResourcePopulate(t *testing.T) {
	expectedId := 10
	expectedUrl := "https://jira.example.com"
	expectedUsername := "jirauser"
	expectedConfigName := "My Jira"
	expectedEpicNameId := 10001
	expectedOpenStatusKey := 11
	expectedCloseStatusKey := 21
	expectedInfo := "Lowest"
	expectedLow := "Low"
	expectedMedium := "Medium"
	expectedHigh := "High"
	expectedCritical := "Critical"

	ddResource := jiraInstanceDefectdojoResource{
		JIRAInstance: dd.JIRAInstance{
			Id:                      &expectedId,
			Url:                     expectedUrl,
			Username:                expectedUsername,
			ConfigurationName:       &expectedConfigName,
			EpicNameId:              expectedEpicNameId,
			OpenStatusKey:           expectedOpenStatusKey,
			CloseStatusKey:          expectedCloseStatusKey,
			InfoMappingSeverity:     expectedInfo,
			LowMappingSeverity:      expectedLow,
			MediumMappingSeverity:   expectedMedium,
			HighMappingSeverity:     expectedHigh,
			CriticalMappingSeverity: expectedCritical,
		},
	}

	resourceData := jiraInstanceResourceData{}
	var trd terraformResourceData = &resourceData

	populateResourceData(context.Background(), &diag.Diagnostics{}, &trd, &ddResource)

	assert.Equal(t, resourceData.Id.ValueString(), "10")
	assert.Equal(t, resourceData.Url.ValueString(), expectedUrl)
	assert.Equal(t, resourceData.Username.ValueString(), expectedUsername)
	assert.Equal(t, resourceData.ConfigurationName.ValueString(), expectedConfigName)
	assert.Equal(t, resourceData.EpicNameId.ValueInt64(), int64(expectedEpicNameId))
	assert.Equal(t, resourceData.OpenStatusKey.ValueInt64(), int64(expectedOpenStatusKey))
	assert.Equal(t, resourceData.CloseStatusKey.ValueInt64(), int64(expectedCloseStatusKey))
	assert.Equal(t, resourceData.InfoMappingSeverity.ValueString(), expectedInfo)
	assert.Equal(t, resourceData.LowMappingSeverity.ValueString(), expectedLow)
	assert.Equal(t, resourceData.MediumMappingSeverity.ValueString(), expectedMedium)
	assert.Equal(t, resourceData.HighMappingSeverity.ValueString(), expectedHigh)
	assert.Equal(t, resourceData.CriticalMappingSeverity.ValueString(), expectedCritical)
}

func TestJiraInstanceResourcePopulateDefectdojo(t *testing.T) {
	resourceData := jiraInstanceResourceData{
		Url:                     types.StringValue("https://jira.example.com"),
		Username:                types.StringValue("user"),
		EpicNameId:              types.Int64Value(10001),
		OpenStatusKey:           types.Int64Value(11),
		CloseStatusKey:          types.Int64Value(21),
		InfoMappingSeverity:     types.StringValue("Lowest"),
		LowMappingSeverity:      types.StringValue("Low"),
		MediumMappingSeverity:   types.StringValue("Medium"),
		HighMappingSeverity:     types.StringValue("High"),
		CriticalMappingSeverity: types.StringValue("Critical"),
	}

	ddRes := resourceData.defectdojoResource()
	ddInstance := ddRes.(*jiraInstanceDefectdojoResource)
	var trd terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, trd, &ddRes)

	assert.Equal(t, ddInstance.Url, "https://jira.example.com")
	assert.Equal(t, ddInstance.Username, "user")
	assert.Equal(t, ddInstance.EpicNameId, 10001)
	assert.Equal(t, ddInstance.OpenStatusKey, 11)
	assert.Equal(t, ddInstance.CloseStatusKey, 21)
	assert.Equal(t, ddInstance.InfoMappingSeverity, "Lowest")
}
