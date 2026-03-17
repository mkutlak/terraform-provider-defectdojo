package provider

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
	"gotest.tools/assert"
)

func TestDDTestResourcePopulate(t *testing.T) {
	expectedId := 25
	expectedTestType := 2
	expectedEngagement := 8
	expectedTitle := "SAST Scan"
	expectedDescription := "Static analysis run"
	expectedVersion := "v2.0"
	expectedBranchTag := "feature/new-feature"
	expectedCommitHash := "def456"
	expectedBuildId := "ci-789"
	expectedPercentComplete := 100
	expectedEnvironment := 3
	expectedLead := 7
	expectedScanType := "Bandit Scan"

	targetStart := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	targetEnd := time.Date(2025, 1, 1, 18, 0, 0, 0, time.UTC)

	ddResource := ddTestDefectdojoResource{
		TestCreate: dd.TestCreate{
			Id:              &expectedId,
			TestType:        expectedTestType,
			Engagement:      expectedEngagement,
			TargetStart:     targetStart,
			TargetEnd:       targetEnd,
			Title:           &expectedTitle,
			Description:     &expectedDescription,
			Version:         &expectedVersion,
			BranchTag:       &expectedBranchTag,
			CommitHash:      &expectedCommitHash,
			BuildId:         &expectedBuildId,
			PercentComplete: &expectedPercentComplete,
			Environment:     &expectedEnvironment,
			Lead:            &expectedLead,
			ScanType:        &expectedScanType,
		},
	}

	resourceData := ddTestResourceData{}
	var tfResource terraformResourceData = &resourceData

	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddResource)

	assert.Equal(t, resourceData.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, resourceData.TestType.ValueInt64(), int64(expectedTestType))
	assert.Equal(t, resourceData.Engagement.ValueInt64(), int64(expectedEngagement))
	assert.Equal(t, resourceData.TargetStart.ValueString(), "2025-01-01T10:00:00Z")
	assert.Equal(t, resourceData.TargetEnd.ValueString(), "2025-01-01T18:00:00Z")
	assert.Equal(t, resourceData.Title.ValueString(), expectedTitle)
	assert.Equal(t, resourceData.Description.ValueString(), expectedDescription)
	assert.Equal(t, resourceData.Version.ValueString(), expectedVersion)
	assert.Equal(t, resourceData.BranchTag.ValueString(), expectedBranchTag)
	assert.Equal(t, resourceData.CommitHash.ValueString(), expectedCommitHash)
	assert.Equal(t, resourceData.BuildId.ValueString(), expectedBuildId)
	assert.Equal(t, resourceData.PercentComplete.ValueInt64(), int64(expectedPercentComplete))
	assert.Equal(t, resourceData.Environment.ValueInt64(), int64(expectedEnvironment))
	assert.Equal(t, resourceData.Lead.ValueInt64(), int64(expectedLead))
	assert.Equal(t, resourceData.ScanType.ValueString(), expectedScanType)
}

func TestDDTestResourcePopulateNils(t *testing.T) {
	ddResource := ddTestDefectdojoResource{
		TestCreate: dd.TestCreate{
			// TestType and Engagement are non-pointer required ints
			TestType:   0,
			Engagement: 0,
			// TargetStart and TargetEnd are non-pointer time.Time (zero value)
		},
	}

	resourceData := ddTestResourceData{}
	var tfResource terraformResourceData = &resourceData

	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddResource)

	assert.Equal(t, resourceData.Id.IsNull(), true)
	assert.Equal(t, resourceData.TestType.ValueInt64(), int64(0))
	assert.Equal(t, resourceData.Engagement.ValueInt64(), int64(0))
	// Zero time.Time -> null
	assert.Equal(t, resourceData.TargetStart.IsNull(), true)
	assert.Equal(t, resourceData.TargetEnd.IsNull(), true)
	assert.Equal(t, resourceData.Title.IsNull(), true)
	assert.Equal(t, resourceData.Description.IsNull(), true)
	assert.Equal(t, resourceData.Version.IsNull(), true)
	assert.Equal(t, resourceData.PercentComplete.IsNull(), true)
	assert.Equal(t, resourceData.Environment.IsNull(), true)
	assert.Equal(t, resourceData.Lead.IsNull(), true)
	assert.Equal(t, resourceData.ScanType.IsNull(), true)
}

func TestDDTestResource__defectdojoResource(t *testing.T) {
	expectedTestType := 2
	expectedEngagement := 8
	expectedTitle := "SAST Scan"
	expectedDescription := "Static analysis run"

	resourceData := ddTestResourceData{
		TestType:    types.Int64Value(int64(expectedTestType)),
		Engagement:  types.Int64Value(int64(expectedEngagement)),
		TargetStart: types.StringValue("2025-01-01T10:00:00Z"),
		TargetEnd:   types.StringValue("2025-01-01T18:00:00Z"),
		Title:       types.StringValue(expectedTitle),
		Description: types.StringValue(expectedDescription),
	}

	ddRes := resourceData.defectdojoResource()
	ddTest := ddRes.(*ddTestDefectdojoResource)
	var tfResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, tfResource, &ddRes)

	assert.Equal(t, ddTest.TestType, expectedTestType)
	assert.Equal(t, ddTest.Engagement, expectedEngagement)
	assert.Equal(t, ddTest.TargetStart.Format(time.RFC3339), "2025-01-01T10:00:00Z")
	assert.Equal(t, ddTest.TargetEnd.Format(time.RFC3339), "2025-01-01T18:00:00Z")
	assert.Equal(t, *ddTest.Title, expectedTitle)
	assert.Equal(t, *ddTest.Description, expectedDescription)
}
