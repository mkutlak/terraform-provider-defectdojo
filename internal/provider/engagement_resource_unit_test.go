package provider

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"gotest.tools/assert"
)

func TestEngagementResourcePopulate(t *testing.T) {
	expectedId := 42
	expectedName := "Test Engagement"
	expectedDescription := "An engagement description"
	expectedProduct := 5
	expectedLead := 3
	expectedReason := "Test reason"
	expectedVersion := "v1.0"
	expectedBranchTag := "main"
	expectedCommitHash := "abc123"
	expectedBuildId := "build-456"
	expectedTracker := "https://jira.example.com/browse/PROJ-1"
	expectedTestStrategy := "https://example.com/strategy"
	expectedThreatModel := true
	expectedApiTest := false
	expectedPenTest := true
	expectedCheckList := false
	expectedDeduplicationOnEngagement := true
	expectedSourceCodeManagementUri := "https://github.com/example/repo"
	expectedPreset := 2
	expectedReportType := 4
	expectedRequester := 6
	expectedEngagementType := dd.EngagementEngagementType("Interactive")
	expectedStatus := dd.EngagementStatus("In Progress")

	targetStartTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	targetEndTime := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
	firstContactedTime := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)

	ddResource := engagementDefectdojoResource{
		Engagement: dd.Engagement{
			Id:                        &expectedId,
			Name:                      &expectedName,
			Description:               &expectedDescription,
			Product:                   expectedProduct,
			TargetStart:               openapi_types.Date{Time: targetStartTime},
			TargetEnd:                 openapi_types.Date{Time: targetEndTime},
			EngagementType:            &expectedEngagementType,
			Status:                    &expectedStatus,
			Lead:                      &expectedLead,
			Reason:                    &expectedReason,
			Version:                   &expectedVersion,
			BranchTag:                 &expectedBranchTag,
			CommitHash:                &expectedCommitHash,
			BuildId:                   &expectedBuildId,
			Tracker:                   &expectedTracker,
			TestStrategy:              &expectedTestStrategy,
			ThreatModel:               &expectedThreatModel,
			ApiTest:                   &expectedApiTest,
			PenTest:                   &expectedPenTest,
			CheckList:                 &expectedCheckList,
			DeduplicationOnEngagement: &expectedDeduplicationOnEngagement,
			FirstContacted:            &openapi_types.Date{Time: firstContactedTime},
			SourceCodeManagementUri:   &expectedSourceCodeManagementUri,
			Preset:                    &expectedPreset,
			ReportType:                &expectedReportType,
			Requester:                 &expectedRequester,
		},
	}

	resourceData := engagementResourceData{}
	var tfResource terraformResourceData = &resourceData

	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddResource)

	assert.Equal(t, resourceData.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, resourceData.Name.ValueString(), expectedName)
	assert.Equal(t, resourceData.Description.ValueString(), expectedDescription)
	assert.Equal(t, resourceData.Product.ValueInt64(), int64(expectedProduct))
	assert.Equal(t, resourceData.TargetStart.ValueString(), "2025-01-01")
	assert.Equal(t, resourceData.TargetEnd.ValueString(), "2025-12-31")
	assert.Equal(t, resourceData.EngagementType.ValueString(), "Interactive")
	assert.Equal(t, resourceData.Status.ValueString(), "In Progress")
	assert.Equal(t, resourceData.Lead.ValueInt64(), int64(expectedLead))
	assert.Equal(t, resourceData.Reason.ValueString(), expectedReason)
	assert.Equal(t, resourceData.Version.ValueString(), expectedVersion)
	assert.Equal(t, resourceData.BranchTag.ValueString(), expectedBranchTag)
	assert.Equal(t, resourceData.CommitHash.ValueString(), expectedCommitHash)
	assert.Equal(t, resourceData.BuildId.ValueString(), expectedBuildId)
	assert.Equal(t, resourceData.Tracker.ValueString(), expectedTracker)
	assert.Equal(t, resourceData.TestStrategy.ValueString(), expectedTestStrategy)
	assert.Equal(t, resourceData.ThreatModel.ValueBool(), expectedThreatModel)
	assert.Equal(t, resourceData.ApiTest.ValueBool(), expectedApiTest)
	assert.Equal(t, resourceData.PenTest.ValueBool(), expectedPenTest)
	assert.Equal(t, resourceData.CheckList.ValueBool(), expectedCheckList)
	assert.Equal(t, resourceData.DeduplicationOnEngagement.ValueBool(), expectedDeduplicationOnEngagement)
	assert.Equal(t, resourceData.FirstContacted.ValueString(), "2024-06-15")
	assert.Equal(t, resourceData.SourceCodeManagementUri.ValueString(), expectedSourceCodeManagementUri)
	assert.Equal(t, resourceData.Preset.ValueInt64(), int64(expectedPreset))
	assert.Equal(t, resourceData.ReportType.ValueInt64(), int64(expectedReportType))
	assert.Equal(t, resourceData.Requester.ValueInt64(), int64(expectedRequester))
}

func TestEngagementResourcePopulateNils(t *testing.T) {
	ddResource := engagementDefectdojoResource{
		Engagement: dd.Engagement{
			// Product is non-pointer required int, TargetStart/TargetEnd are non-pointer Date
			Product: 0,
		},
	}

	resourceData := engagementResourceData{}
	var tfResource terraformResourceData = &resourceData

	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddResource)

	assert.Equal(t, resourceData.Id.IsNull(), true)
	assert.Equal(t, resourceData.Name.IsNull(), true)
	assert.Equal(t, resourceData.Description.IsNull(), true)
	assert.Equal(t, resourceData.Product.ValueInt64(), int64(0))
	// TargetStart/TargetEnd are zero time -> should be null
	assert.Equal(t, resourceData.TargetStart.IsNull(), true)
	assert.Equal(t, resourceData.TargetEnd.IsNull(), true)
	assert.Equal(t, resourceData.EngagementType.IsNull(), true)
	assert.Equal(t, resourceData.Status.IsNull(), true)
	assert.Equal(t, resourceData.Lead.IsNull(), true)
	assert.Equal(t, resourceData.Reason.IsNull(), true)
	assert.Equal(t, resourceData.Version.IsNull(), true)
	assert.Equal(t, resourceData.BranchTag.IsNull(), true)
	assert.Equal(t, resourceData.ThreatModel.IsNull(), true)
	assert.Equal(t, resourceData.Preset.IsNull(), true)
	assert.Equal(t, resourceData.FirstContacted.IsNull(), true)
}

func TestEngagementResource__defectdojoResource(t *testing.T) {
	expectedProduct := 5
	expectedName := "Test Engagement"
	expectedVersion := "v1.0"
	expectedThreatModel := true

	resourceData := engagementResourceData{
		Product:     types.Int64Value(int64(expectedProduct)),
		Name:        types.StringValue(expectedName),
		TargetStart: types.StringValue("2025-01-01"),
		TargetEnd:   types.StringValue("2025-12-31"),
		Version:     types.StringValue(expectedVersion),
		ThreatModel: types.BoolValue(expectedThreatModel),
	}

	ddRes := resourceData.defectdojoResource()
	ddEng := ddRes.(*engagementDefectdojoResource)
	var tfResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, tfResource, &ddRes)

	assert.Equal(t, ddEng.Product, expectedProduct)
	assert.Equal(t, *ddEng.Name, expectedName)
	assert.Equal(t, ddEng.TargetStart.Format("2006-01-02"), "2025-01-01")
	assert.Equal(t, ddEng.TargetEnd.Format("2006-01-02"), "2025-12-31")
	assert.Equal(t, *ddEng.Version, expectedVersion)
	assert.Equal(t, *ddEng.ThreatModel, expectedThreatModel)
}
