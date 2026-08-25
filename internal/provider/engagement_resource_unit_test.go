package provider

import (
	"context"
	"fmt"
	"strings"
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

	// Tracker, TestStrategy and SourceCodeManagementUri are oapi-codegen
	// oneOf-string wrappers as of DefectDojo 3.2 (see isOapiUnionStringType
	// in resource.go); they cannot be built with a plain &string literal.
	tracker := &dd.Engagement_Tracker{}
	assert.NilError(t, tracker.FromEngagementTracker0(expectedTracker))
	testStrategy := &dd.Engagement_TestStrategy{}
	assert.NilError(t, testStrategy.FromEngagementTestStrategy0(expectedTestStrategy))
	sourceCodeManagementUri := &dd.Engagement_SourceCodeManagementUri{}
	assert.NilError(t, sourceCodeManagementUri.FromEngagementSourceCodeManagementUri0(expectedSourceCodeManagementUri))

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
			Tracker:                   tracker,
			TestStrategy:              testStrategy,
			ThreatModel:               &expectedThreatModel,
			ApiTest:                   &expectedApiTest,
			PenTest:                   &expectedPenTest,
			CheckList:                 &expectedCheckList,
			DeduplicationOnEngagement: &expectedDeduplicationOnEngagement,
			FirstContacted:            &openapi_types.Date{Time: firstContactedTime},
			SourceCodeManagementUri:   sourceCodeManagementUri,
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
	// Tracker, TestStrategy and SourceCodeManagementUri are oapi-codegen
	// oneOf-string wrappers (see isOapiUnionStringType in resource.go). A
	// zero-value dd.Engagement leaves their pointers nil, so state must be
	// null, not a wrapper holding an empty string.
	assert.Equal(t, resourceData.Tracker.IsNull(), true)
	assert.Equal(t, resourceData.TestStrategy.IsNull(), true)
	assert.Equal(t, resourceData.SourceCodeManagementUri.IsNull(), true)
}

// TestEngagementResourcePopulate_UnionFieldsEmpty exercises the read path for
// the oneOf's second arm (maxLength: 0), which DefectDojo sends when Tracker,
// TestStrategy or SourceCodeManagementUri is blank. Empty-string coverage was
// flagged as the top gap in the isOapiUnionStringType engine (resource.go).
func TestEngagementResourcePopulate_UnionFieldsEmpty(t *testing.T) {
	tracker := &dd.Engagement_Tracker{}
	assert.NilError(t, tracker.FromEngagementTracker1(""))
	testStrategy := &dd.Engagement_TestStrategy{}
	assert.NilError(t, testStrategy.FromEngagementTestStrategy1(""))
	sourceCodeManagementUri := &dd.Engagement_SourceCodeManagementUri{}
	assert.NilError(t, sourceCodeManagementUri.FromEngagementSourceCodeManagementUri1(""))

	ddResource := engagementDefectdojoResource{
		Engagement: dd.Engagement{
			Product:                 5,
			Tracker:                 tracker,
			TestStrategy:            testStrategy,
			SourceCodeManagementUri: sourceCodeManagementUri,
		},
	}

	resourceData := engagementResourceData{}
	var tfResource terraformResourceData = &resourceData
	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddResource)

	assert.Equal(t, resourceData.Tracker.IsNull(), false)
	assert.Equal(t, resourceData.Tracker.ValueString(), "")
	assert.Equal(t, resourceData.TestStrategy.IsNull(), false)
	assert.Equal(t, resourceData.TestStrategy.ValueString(), "")
	assert.Equal(t, resourceData.SourceCodeManagementUri.IsNull(), false)
	assert.Equal(t, resourceData.SourceCodeManagementUri.ValueString(), "")
}

// TestEngagementResource_defectdojoResource_UnionFieldsEmpty exercises the
// write path when Tracker, TestStrategy or SourceCodeManagementUri is
// configured as an explicit empty string, the other arm of the union a null
// round-trip does not cover.
func TestEngagementResource_defectdojoResource_UnionFieldsEmpty(t *testing.T) {
	resourceData := engagementResourceData{
		Product:                 types.Int64Value(5),
		TargetStart:             types.StringValue("2025-01-01"),
		TargetEnd:               types.StringValue("2025-12-31"),
		Tracker:                 types.StringValue(""),
		TestStrategy:            types.StringValue(""),
		SourceCodeManagementUri: types.StringValue(""),
	}

	ddRes := resourceData.defectdojoResource()
	var tfResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, tfResource, &ddRes)

	ddEng := ddRes.(*engagementDefectdojoResource)
	assert.Assert(t, ddEng.Tracker != nil)
	assert.Assert(t, ddEng.TestStrategy != nil)
	assert.Assert(t, ddEng.SourceCodeManagementUri != nil)

	trackerRaw, err := ddEng.Tracker.MarshalJSON()
	assert.NilError(t, err)
	assert.Equal(t, string(trackerRaw), `""`)

	testStrategyRaw, err := ddEng.TestStrategy.MarshalJSON()
	assert.NilError(t, err)
	assert.Equal(t, string(testStrategyRaw), `""`)

	sourceCodeManagementUriRaw, err := ddEng.SourceCodeManagementUri.MarshalJSON()
	assert.NilError(t, err)
	assert.Equal(t, string(sourceCodeManagementUriRaw), `""`)
}

// TestEngagementResource_defectdojoResource_UnionFieldsNull exercises the
// write path when Tracker, TestStrategy and SourceCodeManagementUri are left
// null in configuration: the wrapper pointers must stay nil, or the request
// would send a value the practitioner never configured.
func TestEngagementResource_defectdojoResource_UnionFieldsNull(t *testing.T) {
	resourceData := engagementResourceData{
		Product:                 types.Int64Value(5),
		TargetStart:             types.StringValue("2025-01-01"),
		TargetEnd:               types.StringValue("2025-12-31"),
		Tracker:                 types.StringNull(),
		TestStrategy:            types.StringNull(),
		SourceCodeManagementUri: types.StringNull(),
	}

	ddRes := resourceData.defectdojoResource()
	var tfResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, tfResource, &ddRes)

	ddEng := ddRes.(*engagementDefectdojoResource)
	assert.Assert(t, ddEng.Tracker == nil)
	assert.Assert(t, ddEng.TestStrategy == nil)
	assert.Assert(t, ddEng.SourceCodeManagementUri == nil)
}

func TestEngagementResource__defectdojoResource(t *testing.T) {
	expectedProduct := 5
	expectedName := "Test Engagement"
	expectedVersion := "v1.0"
	expectedThreatModel := true
	expectedTracker := "https://jira.example.com/browse/PROJ-1"

	resourceData := engagementResourceData{
		Product:     types.Int64Value(int64(expectedProduct)),
		Name:        types.StringValue(expectedName),
		TargetStart: types.StringValue("2025-01-01"),
		TargetEnd:   types.StringValue("2025-12-31"),
		Version:     types.StringValue(expectedVersion),
		ThreatModel: types.BoolValue(expectedThreatModel),
		Tracker:     types.StringValue(expectedTracker),
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

	// Tracker is an oapi-codegen oneOf-string wrapper as of DefectDojo 3.2
	// (see isOapiUnionStringType in resource.go); confirm the write path sets
	// it correctly, then confirm engagementToRequest can carry it across to
	// the distinct EngagementRequest_Tracker wrapper type.
	gotTracker, err := ddEng.Tracker.AsEngagementTracker0()
	assert.NilError(t, err)
	assert.Equal(t, gotTracker, expectedTracker)

	req := engagementToRequest(ddEng.Engagement)
	gotReqTracker, err := req.Tracker.AsEngagementRequestTracker0()
	assert.NilError(t, err)
	assert.Equal(t, gotReqTracker, expectedTracker)
}

// TestEngagementResourcePopulateDedupesServerTags covers the create response
// DefectDojo sends for a child of a product with tag inheritance enabled.
//
// The inheritance signal merges the product's tags into the child and the
// serializer then renders the child's OWN tag twice. Verified on 3.2.300 with a
// product carrying tags ['team-a'] and enable_product_tag_inheritance = true:
//
//	POST /api/v2/engagements/ {"tags":["sprint-1"]}
//	  -> 201, create response tags ['sprint-1', 'team-a', 'sprint-1']
//	  -> GET returns ['sprint-1', 'team-a']
//
// A Terraform set cannot hold that, so the server's list is deduplicated before
// the set is built. See dedupeTagElements in tags.go.
func TestEngagementResourcePopulateDedupesServerTags(t *testing.T) {
	serverTags := []string{"sprint-1", "team-a", "sprint-1"}
	ddEngagement := engagementDefectdojoResource{
		Engagement: dd.Engagement{Tags: &serverTags},
	}

	// The practitioner configured only the engagement's own tag.
	engagementResource := engagementResourceData{Tags: tagSet("sprint-1")}
	var terraformResource terraformResourceData = &engagementResource

	diags := diag.Diagnostics{}
	populateResourceData(context.Background(), &diags, &terraformResource, &ddEngagement)

	assert.Equal(t, diags.HasError(), false)
	assert.Equal(t, len(engagementResource.Tags.Elements()), 2)
	assert.Assert(t, engagementResource.Tags.Equal(tagSet("sprint-1", "team-a")))
}

// TestEngagementResource__defectdojoResourceRejectsDatetime is deliberately
// asymmetric with defectdojo_test: engagement.target_start is an
// openapi_types.Date with no time component, and datetime -> date is
// narrowing and ambiguous, so datetime literals are rejected outright rather
// than silently truncated (see parseDate).
func TestEngagementResource__defectdojoResourceRejectsDatetime(t *testing.T) {
	expectedProduct := 5

	resourceData := engagementResourceData{
		Product:     types.Int64Value(int64(expectedProduct)),
		TargetStart: types.StringValue("2025-01-01T00:00:00Z"),
	}

	ddRes := resourceData.defectdojoResource()
	var tfResource terraformResourceData = &resourceData
	var diags diag.Diagnostics
	populateDefectdojoResource(context.Background(), &diags, tfResource, &ddRes)

	assert.Equal(t, diags.HasError(), true)
	assert.Assert(t, strings.Contains(diags.Errors()[0].Detail(), "target_start"))
	assert.Assert(t, strings.Contains(diags.Errors()[0].Detail(), "2006-01-02"))
}
