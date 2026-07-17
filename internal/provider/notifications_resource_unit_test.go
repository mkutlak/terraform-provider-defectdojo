package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gotest.tools/assert"
)

// TestNotificationsResourcePopulate exercises populateResourceData (DD ->
// Terraform) for the scalar fields plus a representative subset of the 17
// []string Set(String) attributes, including scan_added, to verify the
// engine's element-wise conversion of a *[]string into a types.Set(String).
func TestNotificationsResourcePopulate(t *testing.T) {
	expectedId := 99
	expectedProduct := 10
	expectedUser := 20
	expectedTemplate := true
	expectedScanAddedEmpty := notificationsScanAddedEmpty("alert")

	expectedScanAdded := []string{"alert", "mail"}
	expectedScanAddedSet := types.SetValueMust(
		types.StringType,
		[]attr.Value{types.StringValue("alert"), types.StringValue("mail")},
	)

	expectedSlaBreach := []string{"webhooks"}
	expectedSlaBreachSet := types.SetValueMust(
		types.StringType,
		[]attr.Value{types.StringValue("webhooks")},
	)

	expectedCloseEngagement := []string{"slack"}
	expectedCloseEngagementSet := types.SetValueMust(
		types.StringType,
		[]attr.Value{types.StringValue("slack")},
	)

	expectedJiraUpdate := []string{"msteams"}
	expectedJiraUpdateSet := types.SetValueMust(
		types.StringType,
		[]attr.Value{types.StringValue("msteams")},
	)

	expectedUserMentioned := []string{"alert", "slack"}
	expectedUserMentionedSet := types.SetValueMust(
		types.StringType,
		[]attr.Value{types.StringValue("alert"), types.StringValue("slack")},
	)

	ddNotifications := notificationsDefectdojoResource{
		notificationsModel: notificationsModel{
			Id:              &expectedId,
			Product:         &expectedProduct,
			User:            &expectedUser,
			Template:        &expectedTemplate,
			ScanAddedEmpty:  &expectedScanAddedEmpty,
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
	assert.Equal(t, notificationsResource.ScanAddedEmpty.ValueString(), string(expectedScanAddedEmpty))
	assert.DeepEqual(t, notificationsResource.ScanAdded, expectedScanAddedSet)
	assert.DeepEqual(t, notificationsResource.SlaBreach, expectedSlaBreachSet)
	assert.DeepEqual(t, notificationsResource.CloseEngagement, expectedCloseEngagementSet)
	assert.DeepEqual(t, notificationsResource.JiraUpdate, expectedJiraUpdateSet)
	assert.DeepEqual(t, notificationsResource.UserMentioned, expectedUserMentionedSet)

	// Now populate from an empty notificationsModel (all fields nil/zero) to
	// verify ALL 17 Set(String) attributes, plus the scalar fields, flip to
	// null.
	ddNotifications = notificationsDefectdojoResource{
		notificationsModel: notificationsModel{},
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
// from an empty notificationsModel.
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
		notificationsModel: notificationsModel{},
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
// populateDefectdojoResource (Terraform -> DD) for the enum-slice Set(String)
// attributes, now backed by plain *[]string fields, and then verifies the
// JSON wire format built for an outgoing request: every channel field (all
// genuinely multi-select) is serialized as an array, while scan_added_empty
// is serialized as a bare scalar string -- the only shape the live API's
// write path accepts (a JSON array is rejected with "... is not a valid
// choice", see notificationsScanAddedEmpty).
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
	assert.DeepEqual(t, *ddNotifications.ScanAdded, []string{"alert", "mail"})
	assert.DeepEqual(t, *ddNotifications.SlaBreach, []string{"webhooks"})
	assert.DeepEqual(t, *ddNotifications.CloseEngagement, []string{"slack"})
	assert.DeepEqual(t, *ddNotifications.JiraUpdate, []string{"msteams"})
	assert.DeepEqual(t, *ddNotifications.UserMentioned, []string{"alert", "slack"})

	wire, err := json.Marshal(ddNotifications.notificationsModel)
	assert.NilError(t, err)

	var decoded map[string]any
	assert.NilError(t, json.Unmarshal(wire, &decoded))
	assert.Equal(t, decoded["scan_added_empty"], expectedScanAddedEmpty)
	scanAdded, ok := decoded["scan_added"].([]any)
	assert.Equal(t, ok, true)
	assert.DeepEqual(t, scanAdded, []any{"alert", "mail"})
}

// TestNotificationsResource__defectdojoResource_Nulls verifies that null
// Terraform values are skipped by the engine (populateDefectdojoResource),
// leaving all corresponding notificationsModel pointer fields nil, for every
// field including all 17 []string Set(String) attributes, and that the
// resulting JSON wire format (thanks to the omitempty tags) omits every one
// of them.
func TestNotificationsResource__defectdojoResource_Nulls(t *testing.T) {
	var nilInt *int
	var nilBool *bool
	var nilScanAddedEmpty *notificationsScanAddedEmpty

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

	wire, err := json.Marshal(ddNotifications.notificationsModel)
	assert.NilError(t, err)
	assert.Equal(t, string(wire), "{}")
}

// TestNotificationsScanAddedEmptyUnmarshalJSON verifies the custom
// UnmarshalJSON handles both shapes the live API is observed to return for
// scan_added_empty: a JSON array of 0 or 1 strings (GET/LIST, and
// create/update responses when the field was left unset in the request), and
// a bare scalar string (create/update responses that echo back an
// explicitly-set value).
func TestNotificationsScanAddedEmptyUnmarshalJSON(t *testing.T) {
	var v notificationsScanAddedEmpty

	assert.NilError(t, json.Unmarshal([]byte(`[]`), &v))
	assert.Equal(t, v, notificationsScanAddedEmpty(""))

	assert.NilError(t, json.Unmarshal([]byte(`["alert"]`), &v))
	assert.Equal(t, v, notificationsScanAddedEmpty("alert"))

	assert.NilError(t, json.Unmarshal([]byte(`"mail"`), &v))
	assert.Equal(t, v, notificationsScanAddedEmpty("mail"))

	assert.NilError(t, json.Unmarshal([]byte(`""`), &v))
	assert.Equal(t, v, notificationsScanAddedEmpty(""))

	assert.ErrorContains(t, json.Unmarshal([]byte(`42`), &v), "scan_added_empty")
}

// TestNotificationsModelUnmarshalFullResponse round-trips a full
// /api/v2/notifications/ response body (as observed live, with
// scan_added_empty as an array) through notificationsModel, mirroring what
// readApiCall/createApiCall/updateApiCall do with the raw HTTP response body.
func TestNotificationsModelUnmarshalFullResponse(t *testing.T) {
	body := []byte(`{
		"id": 5,
		"product": 10,
		"user": null,
		"scan_added": ["alert"],
		"scan_added_empty": ["alert"],
		"template": false
	}`)

	var m notificationsModel
	assert.NilError(t, json.Unmarshal(body, &m))
	assert.Equal(t, *m.Id, 5)
	assert.Equal(t, *m.Product, 10)
	assert.Equal(t, *m.ScanAddedEmpty, notificationsScanAddedEmpty("alert"))
	assert.DeepEqual(t, *m.ScanAdded, []string{"alert"})
}
