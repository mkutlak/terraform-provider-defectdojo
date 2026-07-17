package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
	"gotest.tools/assert"
)

// TestNotificationsResourcePopulate exercises populateResourceData (DD ->
// Terraform) for the scalar fields plus a representative subset of the 17
// enum-slice Set(String) attributes, including scan_added, to verify the
// engine's element-wise conversion of defined string element types (e.g.
// dd.NotificationsScanAdded) into a types.Set(String).
func TestNotificationsResourcePopulate(t *testing.T) {
	expectedId := 99
	expectedProduct := 10
	expectedUser := 20
	expectedTemplate := true
	expectedScanAddedEmpty := "alert"

	expectedScanAdded := []dd.NotificationsScanAdded{"alert", "mail"}
	expectedScanAddedSet := types.SetValueMust(
		types.StringType,
		[]attr.Value{types.StringValue("alert"), types.StringValue("mail")},
	)

	expectedSlaBreach := []dd.NotificationsSlaBreach{"webhooks"}
	expectedSlaBreachSet := types.SetValueMust(
		types.StringType,
		[]attr.Value{types.StringValue("webhooks")},
	)

	expectedCloseEngagement := []dd.NotificationsCloseEngagement{"slack"}
	expectedCloseEngagementSet := types.SetValueMust(
		types.StringType,
		[]attr.Value{types.StringValue("slack")},
	)

	expectedJiraUpdate := []dd.NotificationsJiraUpdate{"msteams"}
	expectedJiraUpdateSet := types.SetValueMust(
		types.StringType,
		[]attr.Value{types.StringValue("msteams")},
	)

	expectedUserMentioned := []dd.NotificationsUserMentioned{"alert", "slack"}
	expectedUserMentionedSet := types.SetValueMust(
		types.StringType,
		[]attr.Value{types.StringValue("alert"), types.StringValue("slack")},
	)

	ddNotifications := notificationsDefectdojoResource{
		Notifications: dd.Notifications{
			Id:              &expectedId,
			Product:         &expectedProduct,
			User:            &expectedUser,
			Template:        &expectedTemplate,
			ScanAddedEmpty:  (*dd.NotificationsScanAddedEmpty)(&expectedScanAddedEmpty),
			ScanAdded:       &expectedScanAdded,
			SlaBreach:       &expectedSlaBreach,
			CloseEngagement: &expectedCloseEngagement,
			JiraUpdate:      &expectedJiraUpdate,
			UserMentioned:   &expectedUserMentioned,
		},
	}

	notificationsResource := notificationsResourceData{}
	var terraformResource terraformResourceData = &notificationsResource

	populateResourceData(context.Background(), &diag.Diagnostics{}, &terraformResource, &ddNotifications)
	assert.Equal(t, notificationsResource.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, notificationsResource.Product.ValueInt64(), (int64)(expectedProduct))
	assert.Equal(t, notificationsResource.User.ValueInt64(), (int64)(expectedUser))
	assert.Equal(t, notificationsResource.Template.ValueBool(), expectedTemplate)
	assert.Equal(t, notificationsResource.ScanAddedEmpty.ValueString(), expectedScanAddedEmpty)
	assert.DeepEqual(t, notificationsResource.ScanAdded, expectedScanAddedSet)
	assert.DeepEqual(t, notificationsResource.SlaBreach, expectedSlaBreachSet)
	assert.DeepEqual(t, notificationsResource.CloseEngagement, expectedCloseEngagementSet)
	assert.DeepEqual(t, notificationsResource.JiraUpdate, expectedJiraUpdateSet)
	assert.DeepEqual(t, notificationsResource.UserMentioned, expectedUserMentionedSet)

	// Now populate from an empty Notifications (all fields nil/zero) to verify
	// ALL 17 Set(String) attributes, plus the scalar fields, flip to null.
	ddNotifications = notificationsDefectdojoResource{
		Notifications: dd.Notifications{},
	}
	populateResourceData(context.Background(), &diag.Diagnostics{}, &terraformResource, &ddNotifications)

	nilStringSet := types.SetNull(types.StringType)

	assert.Equal(t, notificationsResource.Id.ValueString(), "")
	assert.Equal(t, notificationsResource.Product.IsNull(), true)
	assert.Equal(t, notificationsResource.User.IsNull(), true)
	assert.Equal(t, notificationsResource.Template.IsNull(), true)
	assert.Equal(t, notificationsResource.ScanAddedEmpty.IsNull(), true)

	assert.DeepEqual(t, notificationsResource.AutoCloseEngagement, nilStringSet)
	assert.DeepEqual(t, notificationsResource.CloseEngagement, nilStringSet)
	assert.DeepEqual(t, notificationsResource.CodeReview, nilStringSet)
	assert.DeepEqual(t, notificationsResource.EngagementAdded, nilStringSet)
	assert.DeepEqual(t, notificationsResource.JiraUpdate, nilStringSet)
	assert.DeepEqual(t, notificationsResource.Other, nilStringSet)
	assert.DeepEqual(t, notificationsResource.ProductAdded, nilStringSet)
	assert.DeepEqual(t, notificationsResource.ProductTypeAdded, nilStringSet)
	assert.DeepEqual(t, notificationsResource.ReviewRequested, nilStringSet)
	assert.DeepEqual(t, notificationsResource.RiskAcceptanceExpiration, nilStringSet)
	assert.DeepEqual(t, notificationsResource.ScanAdded, nilStringSet)
	assert.DeepEqual(t, notificationsResource.SlaBreach, nilStringSet)
	assert.DeepEqual(t, notificationsResource.SlaBreachCombined, nilStringSet)
	assert.DeepEqual(t, notificationsResource.StaleEngagement, nilStringSet)
	assert.DeepEqual(t, notificationsResource.TestAdded, nilStringSet)
	assert.DeepEqual(t, notificationsResource.UpcomingEngagement, nilStringSet)
	assert.DeepEqual(t, notificationsResource.UserMentioned, nilStringSet)
}

// TestNotificationsResourcePopulateNils verifies a zero-value
// notificationsResourceData starts empty, and stays empty after populating
// from an empty dd.Notifications.
func TestNotificationsResourcePopulateNils(t *testing.T) {
	notificationsResource := notificationsResourceData{}
	var terraformResource terraformResourceData = &notificationsResource

	assert.Equal(t, notificationsResource.Id.ValueString(), "")
	assert.Equal(t, notificationsResource.Product.ValueInt64(), (int64)(0))
	assert.Equal(t, notificationsResource.User.ValueInt64(), (int64)(0))
	assert.Equal(t, notificationsResource.Template.ValueBool(), false)
	assert.Equal(t, notificationsResource.ScanAddedEmpty.ValueString(), "")

	assert.DeepEqual(t, notificationsResource.AutoCloseEngagement.Elements(), []attr.Value{})
	assert.DeepEqual(t, notificationsResource.ScanAdded.Elements(), []attr.Value{})

	ddNotifications := notificationsDefectdojoResource{
		Notifications: dd.Notifications{},
	}
	populateResourceData(context.Background(), &diag.Diagnostics{}, &terraformResource, &ddNotifications)

	nilStringSet := types.SetNull(types.StringType)

	// still all empty/null values after running populate
	assert.Equal(t, notificationsResource.Id.ValueString(), "")
	assert.Equal(t, notificationsResource.Product.IsNull(), true)
	assert.Equal(t, notificationsResource.User.IsNull(), true)
	assert.Equal(t, notificationsResource.Template.IsNull(), true)
	assert.Equal(t, notificationsResource.ScanAddedEmpty.IsNull(), true)
	assert.DeepEqual(t, notificationsResource.AutoCloseEngagement, nilStringSet)
	assert.DeepEqual(t, notificationsResource.CloseEngagement, nilStringSet)
	assert.DeepEqual(t, notificationsResource.CodeReview, nilStringSet)
	assert.DeepEqual(t, notificationsResource.EngagementAdded, nilStringSet)
	assert.DeepEqual(t, notificationsResource.JiraUpdate, nilStringSet)
	assert.DeepEqual(t, notificationsResource.Other, nilStringSet)
	assert.DeepEqual(t, notificationsResource.ProductAdded, nilStringSet)
	assert.DeepEqual(t, notificationsResource.ProductTypeAdded, nilStringSet)
	assert.DeepEqual(t, notificationsResource.ReviewRequested, nilStringSet)
	assert.DeepEqual(t, notificationsResource.RiskAcceptanceExpiration, nilStringSet)
	assert.DeepEqual(t, notificationsResource.ScanAdded, nilStringSet)
	assert.DeepEqual(t, notificationsResource.SlaBreach, nilStringSet)
	assert.DeepEqual(t, notificationsResource.SlaBreachCombined, nilStringSet)
	assert.DeepEqual(t, notificationsResource.StaleEngagement, nilStringSet)
	assert.DeepEqual(t, notificationsResource.TestAdded, nilStringSet)
	assert.DeepEqual(t, notificationsResource.UpcomingEngagement, nilStringSet)
	assert.DeepEqual(t, notificationsResource.UserMentioned, nilStringSet)
}

// TestNotificationsResource__defectdojoResource exercises
// populateDefectdojoResource (Terraform -> DD) plus notificationsToRequest,
// verifying the enum-slice round trip through the reflection engine (TF Set
// -> *[]dd.NotificationsXxx) and then through the generic
// notificationsConvertEnumSlice helper (*[]dd.NotificationsXxx ->
// *[]dd.NotificationsRequestXxx).
func TestNotificationsResource__defectdojoResource(t *testing.T) {
	expectedProduct := 10
	expectedUser := 20
	expectedTemplate := true
	expectedScanAddedEmpty := "alert"

	expectedScanAddedSet := types.SetValueMust(
		types.StringType,
		[]attr.Value{types.StringValue("alert"), types.StringValue("mail")},
	)
	expectedSlaBreachSet := types.SetValueMust(
		types.StringType,
		[]attr.Value{types.StringValue("webhooks")},
	)
	expectedCloseEngagementSet := types.SetValueMust(
		types.StringType,
		[]attr.Value{types.StringValue("slack")},
	)
	expectedJiraUpdateSet := types.SetValueMust(
		types.StringType,
		[]attr.Value{types.StringValue("msteams")},
	)
	expectedUserMentionedSet := types.SetValueMust(
		types.StringType,
		[]attr.Value{types.StringValue("alert"), types.StringValue("slack")},
	)

	notificationsResource := notificationsResourceData{
		Product:         types.Int64Value(int64(expectedProduct)),
		User:            types.Int64Value(int64(expectedUser)),
		Template:        types.BoolValue(expectedTemplate),
		ScanAddedEmpty:  types.StringValue(expectedScanAddedEmpty),
		ScanAdded:       expectedScanAddedSet,
		SlaBreach:       expectedSlaBreachSet,
		CloseEngagement: expectedCloseEngagementSet,
		JiraUpdate:      expectedJiraUpdateSet,
		UserMentioned:   expectedUserMentionedSet,
	}

	ddResource := notificationsResource.defectdojoResource()
	ddNotifications := ddResource.(*notificationsDefectdojoResource)
	var terraformResource terraformResourceData = &notificationsResource
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, terraformResource, &ddResource)

	assert.Equal(t, *ddNotifications.Product, expectedProduct)
	assert.Equal(t, *ddNotifications.User, expectedUser)
	assert.Equal(t, *ddNotifications.Template, expectedTemplate)
	assert.Equal(t, string(*ddNotifications.ScanAddedEmpty), expectedScanAddedEmpty)
	assert.DeepEqual(t, *ddNotifications.ScanAdded, []dd.NotificationsScanAdded{"alert", "mail"})
	assert.DeepEqual(t, *ddNotifications.SlaBreach, []dd.NotificationsSlaBreach{"webhooks"})
	assert.DeepEqual(t, *ddNotifications.CloseEngagement, []dd.NotificationsCloseEngagement{"slack"})
	assert.DeepEqual(t, *ddNotifications.JiraUpdate, []dd.NotificationsJiraUpdate{"msteams"})
	assert.DeepEqual(t, *ddNotifications.UserMentioned, []dd.NotificationsUserMentioned{"alert", "slack"})

	req := notificationsToRequest(ddNotifications.Notifications)
	assert.Equal(t, *req.Product, expectedProduct)
	assert.Equal(t, *req.User, expectedUser)
	assert.Equal(t, *req.Template, expectedTemplate)
	assert.Equal(t, string(*req.ScanAddedEmpty), expectedScanAddedEmpty)
	assert.DeepEqual(t, *req.ScanAdded, []dd.NotificationsRequestScanAdded{"alert", "mail"})
	assert.DeepEqual(t, *req.SlaBreach, []dd.NotificationsRequestSlaBreach{"webhooks"})
	assert.DeepEqual(t, *req.CloseEngagement, []dd.NotificationsRequestCloseEngagement{"slack"})
	assert.DeepEqual(t, *req.JiraUpdate, []dd.NotificationsRequestJiraUpdate{"msteams"})
	assert.DeepEqual(t, *req.UserMentioned, []dd.NotificationsRequestUserMentioned{"alert", "slack"})
}

// TestNotificationsResource__defectdojoResource_Nulls verifies that null
// Terraform values are skipped by the engine (populateDefectdojoResource),
// leaving all corresponding dd.Notifications pointer fields nil, for every
// field including all 17 enum-slice Set(String) attributes.
func TestNotificationsResource__defectdojoResource_Nulls(t *testing.T) {
	var nilInt *int
	var nilBool *bool
	var nilScanAddedEmpty *dd.NotificationsScanAddedEmpty

	notificationsResource := notificationsResourceData{
		Id:             types.StringNull(),
		Product:        types.Int64Null(),
		User:           types.Int64Null(),
		Template:       types.BoolNull(),
		ScanAddedEmpty: types.StringNull(),

		AutoCloseEngagement:      types.SetNull(types.StringType),
		CloseEngagement:          types.SetNull(types.StringType),
		CodeReview:               types.SetNull(types.StringType),
		EngagementAdded:          types.SetNull(types.StringType),
		JiraUpdate:               types.SetNull(types.StringType),
		Other:                    types.SetNull(types.StringType),
		ProductAdded:             types.SetNull(types.StringType),
		ProductTypeAdded:         types.SetNull(types.StringType),
		ReviewRequested:          types.SetNull(types.StringType),
		RiskAcceptanceExpiration: types.SetNull(types.StringType),
		ScanAdded:                types.SetNull(types.StringType),
		SlaBreach:                types.SetNull(types.StringType),
		SlaBreachCombined:        types.SetNull(types.StringType),
		StaleEngagement:          types.SetNull(types.StringType),
		TestAdded:                types.SetNull(types.StringType),
		UpcomingEngagement:       types.SetNull(types.StringType),
		UserMentioned:            types.SetNull(types.StringType),
	}

	ddResource := notificationsResource.defectdojoResource()
	ddNotifications := ddResource.(*notificationsDefectdojoResource)
	var terraformResource terraformResourceData = &notificationsResource
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, terraformResource, &ddResource)

	assert.Equal(t, ddNotifications.Id, nilInt)
	assert.Equal(t, ddNotifications.Product, nilInt)
	assert.Equal(t, ddNotifications.User, nilInt)
	assert.Equal(t, ddNotifications.Template, nilBool)
	assert.Equal(t, ddNotifications.ScanAddedEmpty, nilScanAddedEmpty)

	// Null TF values are skipped, so pointer-to-slice fields remain nil.
	assert.Equal(t, ddNotifications.AutoCloseEngagement == nil, true)
	assert.Equal(t, ddNotifications.CloseEngagement == nil, true)
	assert.Equal(t, ddNotifications.CodeReview == nil, true)
	assert.Equal(t, ddNotifications.EngagementAdded == nil, true)
	assert.Equal(t, ddNotifications.JiraUpdate == nil, true)
	assert.Equal(t, ddNotifications.Other == nil, true)
	assert.Equal(t, ddNotifications.ProductAdded == nil, true)
	assert.Equal(t, ddNotifications.ProductTypeAdded == nil, true)
	assert.Equal(t, ddNotifications.ReviewRequested == nil, true)
	assert.Equal(t, ddNotifications.RiskAcceptanceExpiration == nil, true)
	assert.Equal(t, ddNotifications.ScanAdded == nil, true)
	assert.Equal(t, ddNotifications.SlaBreach == nil, true)
	assert.Equal(t, ddNotifications.SlaBreachCombined == nil, true)
	assert.Equal(t, ddNotifications.StaleEngagement == nil, true)
	assert.Equal(t, ddNotifications.TestAdded == nil, true)
	assert.Equal(t, ddNotifications.UpcomingEngagement == nil, true)
	assert.Equal(t, ddNotifications.UserMentioned == nil, true)

	req := notificationsToRequest(ddNotifications.Notifications)
	assert.Equal(t, req.Product, nilInt)
	assert.Equal(t, req.User, nilInt)
	assert.Equal(t, req.Template, nilBool)
	assert.Equal(t, req.ScanAdded == nil, true)
	assert.Equal(t, req.UserMentioned == nil, true)
}
