package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
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

// TestJiraInstanceResourcePreservesPasswordOnRefresh asserts that a refresh
// does not clobber the configured secret.
//
// DefectDojo never echoes the password back - it exists on JIRAInstanceRequest
// but not on the JIRAInstance response model - so while the Terraform field
// carried a ddField tag, every refresh rebuilt the wrapper with a nil password
// and populateResourceData wrote StringNull() over it. The result was a
// permanent, retry-proof diff on a credential.
//
// user_resource.go has always used the shape this now copies.
func TestJiraInstanceResourcePreservesPasswordOnRefresh(t *testing.T) {
	resourceData := jiraInstanceResourceData{
		Url:      types.StringValue("https://jira.example.com"),
		Username: types.StringValue("jirauser"),
		Password: types.StringValue("s3cret"),
	}

	// The write path must still carry the secret into the request body.
	ddResource := resourceData.defectdojoResource().(*jiraInstanceDefectdojoResource)
	req := jiraInstanceToRequest(ddResource)
	assert.Assert(t, req.Password != nil)
	assert.Equal(t, *req.Password, "s3cret")

	// Simulate a read: the response model has no password field at all.
	refreshed := &jiraInstanceDefectdojoResource{
		JIRAInstance: dd.JIRAInstance{Url: "https://jira.example.com", Username: "jirauser"},
	}
	var terraformResource terraformResourceData = &resourceData
	diags := diag.Diagnostics{}
	populateResourceData(context.Background(), &diags, &terraformResource, refreshed)

	assert.Equal(t, diags.HasError(), false)
	assert.Equal(t, resourceData.Password.ValueString(), "s3cret")

	// An unset password goes out as nil rather than as an empty string, so
	// DefectDojo leaves the stored credential untouched.
	resourceData.Password = types.StringNull()
	ddResource = resourceData.defectdojoResource().(*jiraInstanceDefectdojoResource)
	assert.Assert(t, jiraInstanceToRequest(ddResource).Password == nil)
}
