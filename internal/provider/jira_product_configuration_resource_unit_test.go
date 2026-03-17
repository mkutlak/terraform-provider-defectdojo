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

func TestJiraProductConfigurationResourcePopulate(t *testing.T) {
	expectedId := 3
	expectedProjectKey := "PROJ"
	expectedIssueTemplateDir := "custom_templates"
	expectedPushAllIssues := true
	expectedEnableEngagementEpicMapping := false
	expectedPushNotes := true
	expectedProductJiraSlaNotification := false
	expectedRiskAcceptanceExpirationNotification := true
	expectedJiraInstance := 5
	expectedProduct := 10
	expectedEngagement := 15
	expectedAddVulnerabilityIdToJiraLabel := true
	expectedComponent := "backend"
	expectedDefaultAssignee := "john.doe"
	expectedEnabled := true
	expectedEpicIssueTypeName := "Epic"
	expectedJiraLabels := "security defectdojo"

	ddResource := jiraProductConfigurationDefectdojoResource{
		JIRAProject: dd.JIRAProject{
			Id:                                   &expectedId,
			ProjectKey:                           &expectedProjectKey,
			IssueTemplateDir:                     &expectedIssueTemplateDir,
			PushAllIssues:                        &expectedPushAllIssues,
			EnableEngagementEpicMapping:          &expectedEnableEngagementEpicMapping,
			PushNotes:                            &expectedPushNotes,
			ProductJiraSlaNotification:           &expectedProductJiraSlaNotification,
			RiskAcceptanceExpirationNotification: &expectedRiskAcceptanceExpirationNotification,
			JiraInstance:                         &expectedJiraInstance,
			Product:                              &expectedProduct,
			Engagement:                           &expectedEngagement,
			AddVulnerabilityIdToJiraLabel:        &expectedAddVulnerabilityIdToJiraLabel,
			Component:                            &expectedComponent,
			DefaultAssignee:                      &expectedDefaultAssignee,
			Enabled:                              &expectedEnabled,
			EpicIssueTypeName:                    &expectedEpicIssueTypeName,
			JiraLabels:                           &expectedJiraLabels,
		},
	}

	resourceData := jiraProductConfigurationResourceData{}
	var tfResource terraformResourceData = &resourceData

	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddResource)

	assert.Equal(t, resourceData.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, resourceData.ProjectKey.ValueString(), expectedProjectKey)
	assert.Equal(t, resourceData.IssueTemplateDir.ValueString(), expectedIssueTemplateDir)
	assert.Equal(t, resourceData.PushAllIssues.ValueBool(), expectedPushAllIssues)
	assert.Equal(t, resourceData.EnableEngagementEpicMapping.ValueBool(), expectedEnableEngagementEpicMapping)
	assert.Equal(t, resourceData.PushNotes.ValueBool(), expectedPushNotes)
	assert.Equal(t, resourceData.ProductJiraSlaNotification.ValueBool(), expectedProductJiraSlaNotification)
	assert.Equal(t, resourceData.RiskAcceptanceExpirationNotification.ValueBool(), expectedRiskAcceptanceExpirationNotification)
	assert.Equal(t, resourceData.JiraInstance.ValueString(), fmt.Sprint(expectedJiraInstance))
	assert.Equal(t, resourceData.Product.ValueString(), fmt.Sprint(expectedProduct))
	assert.Equal(t, resourceData.Engagement.ValueString(), fmt.Sprint(expectedEngagement))
	assert.Equal(t, resourceData.AddVulnerabilityIdToJiraLabel.ValueBool(), expectedAddVulnerabilityIdToJiraLabel)
	assert.Equal(t, resourceData.Component.ValueString(), expectedComponent)
	assert.Equal(t, resourceData.DefaultAssignee.ValueString(), expectedDefaultAssignee)
	assert.Equal(t, resourceData.Enabled.ValueBool(), expectedEnabled)
	assert.Equal(t, resourceData.EpicIssueTypeName.ValueString(), expectedEpicIssueTypeName)
	assert.Equal(t, resourceData.JiraLabels.ValueString(), expectedJiraLabels)
}

func TestJiraProductConfigurationResourcePopulateNils(t *testing.T) {
	ddResource := jiraProductConfigurationDefectdojoResource{
		JIRAProject: dd.JIRAProject{},
	}

	resourceData := jiraProductConfigurationResourceData{}
	var tfResource terraformResourceData = &resourceData

	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddResource)

	assert.Equal(t, resourceData.Id.IsNull(), true)
	assert.Equal(t, resourceData.ProjectKey.IsNull(), true)
	assert.Equal(t, resourceData.IssueTemplateDir.IsNull(), true)
	assert.Equal(t, resourceData.PushAllIssues.IsNull(), true)
	assert.Equal(t, resourceData.EnableEngagementEpicMapping.IsNull(), true)
	assert.Equal(t, resourceData.PushNotes.IsNull(), true)
	assert.Equal(t, resourceData.ProductJiraSlaNotification.IsNull(), true)
	assert.Equal(t, resourceData.RiskAcceptanceExpirationNotification.IsNull(), true)
	assert.Equal(t, resourceData.JiraInstance.IsNull(), true)
	assert.Equal(t, resourceData.Product.IsNull(), true)
	assert.Equal(t, resourceData.Engagement.IsNull(), true)
	assert.Equal(t, resourceData.AddVulnerabilityIdToJiraLabel.IsNull(), true)
	assert.Equal(t, resourceData.Component.IsNull(), true)
	assert.Equal(t, resourceData.DefaultAssignee.IsNull(), true)
	assert.Equal(t, resourceData.Enabled.IsNull(), true)
	assert.Equal(t, resourceData.EpicIssueTypeName.IsNull(), true)
	assert.Equal(t, resourceData.JiraLabels.IsNull(), true)
}

func TestJiraProductConfigurationResource__defectdojoResource(t *testing.T) {
	expectedProjectKey := "PROJ"
	expectedPushAllIssues := true
	expectedJiraInstance := 5
	expectedProduct := 10
	expectedEnabled := true
	expectedComponent := "backend"

	resourceData := jiraProductConfigurationResourceData{
		ProjectKey:    types.StringValue(expectedProjectKey),
		PushAllIssues: types.BoolValue(expectedPushAllIssues),
		JiraInstance:  types.StringValue(fmt.Sprint(expectedJiraInstance)),
		Product:       types.StringValue(fmt.Sprint(expectedProduct)),
		Enabled:       types.BoolValue(expectedEnabled),
		Component:     types.StringValue(expectedComponent),
	}

	ddRes := resourceData.defectdojoResource()
	ddJIRA := ddRes.(*jiraProductConfigurationDefectdojoResource)
	var tfResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, tfResource, &ddRes)

	assert.Equal(t, *ddJIRA.ProjectKey, expectedProjectKey)
	assert.Equal(t, *ddJIRA.PushAllIssues, expectedPushAllIssues)
	assert.Equal(t, *ddJIRA.JiraInstance, expectedJiraInstance)
	assert.Equal(t, *ddJIRA.Product, expectedProduct)
	assert.Equal(t, *ddJIRA.Enabled, expectedEnabled)
	assert.Equal(t, *ddJIRA.Component, expectedComponent)
}
