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

func TestFindingTemplateResourcePopulate(t *testing.T) {
	expectedId := 42
	expectedTitle := "Test Finding Template"
	expectedSeverity := "High"
	expectedDescription := "A test description"
	expectedMitigation := "Apply patch"
	expectedImpact := "Data exposure"
	expectedReferences := "https://example.com/cve"
	expectedCwe := 79
	expectedCvssv3 := "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H"
	expectedCvssv3Score := 9.8
	expectedComponentName := "openssl"
	expectedComponentVersion := "1.0.2"
	expectedFixAvailable := true
	expectedFixVersion := "1.0.3"

	ddObj := findingTemplateDefectdojoResource{
		FindingTemplate: dd.FindingTemplate{
			Id:               &expectedId,
			Title:            expectedTitle,
			Severity:         &expectedSeverity,
			Description:      &expectedDescription,
			Mitigation:       &expectedMitigation,
			Impact:           &expectedImpact,
			References:       &expectedReferences,
			Cwe:              &expectedCwe,
			Cvssv3:           &expectedCvssv3,
			Cvssv3Score:      &expectedCvssv3Score,
			ComponentName:    &expectedComponentName,
			ComponentVersion: &expectedComponentVersion,
			FixAvailable:     &expectedFixAvailable,
			FixVersion:       &expectedFixVersion,
		},
	}

	resourceData := findingTemplateResourceData{}
	var tfResource terraformResourceData = &resourceData
	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddObj)

	assert.Equal(t, resourceData.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, resourceData.Title.ValueString(), expectedTitle)
	assert.Equal(t, resourceData.Severity.ValueString(), expectedSeverity)
	assert.Equal(t, resourceData.Description.ValueString(), expectedDescription)
	assert.Equal(t, resourceData.Mitigation.ValueString(), expectedMitigation)
	assert.Equal(t, resourceData.Impact.ValueString(), expectedImpact)
	assert.Equal(t, resourceData.References.ValueString(), expectedReferences)
	assert.Equal(t, resourceData.Cwe.ValueInt64(), int64(expectedCwe))
	assert.Equal(t, resourceData.Cvssv3.ValueString(), expectedCvssv3)
	assert.Equal(t, resourceData.ComponentName.ValueString(), expectedComponentName)
	assert.Equal(t, resourceData.ComponentVersion.ValueString(), expectedComponentVersion)
	assert.Equal(t, resourceData.FixAvailable.ValueBool(), expectedFixAvailable)
	assert.Equal(t, resourceData.FixVersion.ValueString(), expectedFixVersion)
}

func TestFindingTemplateResource_defectdojoResource(t *testing.T) {
	expectedTitle := "Test Finding Template"
	expectedSeverity := "Medium"
	expectedDescription := "A test description"

	resourceData := findingTemplateResourceData{
		Title:       types.StringValue(expectedTitle),
		Severity:    types.StringValue(expectedSeverity),
		Description: types.StringValue(expectedDescription),
	}

	ddResource := resourceData.defectdojoResource()
	var tfResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, tfResource, &ddResource)

	ddObj := ddResource.(*findingTemplateDefectdojoResource)
	assert.Equal(t, ddObj.Title, expectedTitle)
	assert.Equal(t, *ddObj.Severity, expectedSeverity)
	assert.Equal(t, *ddObj.Description, expectedDescription)
}
