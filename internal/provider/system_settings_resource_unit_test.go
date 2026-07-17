package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
	"gotest.tools/assert"
)

// systemSettingsResourceData / systemSettingsDefectdojoResource together map
// ~72 fields (46 bool, 13 string, 12 int, 1 enum) plus id. These tests cover
// every mapped field in both reflection directions (Populate and
// defectdojoResource), plus the JiraMinimumSeverity enum cast in both
// directions and the systemSettingsToRequest converter.

// compile-time assertion: systemSettingsDefectdojoResource must implement
// singletonAdopter for the engine's Create/Delete adoption path to activate.
var _ singletonAdopter = &systemSettingsDefectdojoResource{}

func TestSystemSettingsResourcePopulate(t *testing.T) {
	expectedId := 501

	// bools
	bAddVulnerabilityIdToJiraLabel := true
	bAllowAnonymousSurveyRepsonse := false
	bApiExposeErrorDetails := true
	bDeleteDuplicates := false
	bDisableJiraWebhookSecret := true
	bDisclaimerReportsForced := false
	bEnableBenchmark := true
	bEnableCalendar := false
	bEnableChecklists := true
	bEnableCvss3Display := false
	bEnableCvss4Display := true
	bEnableDeduplication := false
	bEnableEndpointMetadataImport := true
	bEnableFindingGroups := false
	bEnableFindingSla := true
	bEnableGithub := false
	bEnableJira := true
	bEnableJiraWebHook := false
	bEnableMailNotifications := true
	bEnableMsteamsNotifications := false
	bEnableNotifySlaActive := true
	bEnableNotifySlaActiveVerified := false
	bEnableNotifySlaExponentialBackoff := true
	bEnableNotifySlaJiraOnly := false
	bEnableProductGrade := true
	bEnableProductTagInheritance := false
	bEnableProductTrackingFiles := true
	bEnableQuestionnaires := false
	bEnableSimilarFindings := true
	bEnableSlackNotifications := false
	bEnableUiTableBasedSearching := true
	bEnableUserProfileEditable := false
	bEnableWebhooksNotifications := true
	bEnforceVerifiedStatus := false
	bEnforceVerifiedStatusJira := true
	bEnforceVerifiedStatusMetrics := false
	bEnforceVerifiedStatusProductGrading := true
	bEngagementAutoClose := false
	bFalsePositiveHistory := true
	bFilterStringMatching := false
	bLowercaseCharacterRequired := true
	bNonCommonPasswordRequired := false
	bNumberCharacterRequired := true
	bRetroactiveFalsePositiveHistory := false
	bSpecialCharacterRequired := true
	bUppercaseCharacterRequired := false

	// strings
	sDisclaimerNotes := "notes disclaimer"
	sDisclaimerNotifications := "notifications disclaimer"
	sDisclaimerReports := "reports disclaimer"
	sEmailFrom := "noreply@example.com"
	sJiraLabels := "label1 label2"
	sJiraWebhookSecret := "shh-secret"
	sMailNotificationsTo := "admin@example.com"
	sMsteamsUrl := "https://example.com/msteams"
	sSlackChannel := "#security"
	sSlackToken := "xoxb-token"
	sSlackUsername := "defectdojo-bot"
	sTeamName := "Acme Security"
	sUrlPrefix := "/dojo"

	// ints
	iEngagementAutoCloseDays := 3
	iMaxDupes := 5
	iMaximumPasswordLength := 64
	iMinimumPasswordLength := 8
	iProductGradeA := 90
	iProductGradeB := 80
	iProductGradeC := 70
	iProductGradeD := 60
	iProductGradeF := 59
	iRiskAcceptanceFormDefaultDays := 180
	iRiskAcceptanceNotifyBeforeExpiration := 10
	iWebhooksNotificationsTimeout := 30

	eJiraMinimumSeverity := dd.SystemSettingsJiraMinimumSeverityCritical

	ddSettings := systemSettingsDefectdojoResource{
		SystemSettings: dd.SystemSettings{
			Id: &expectedId,

			AddVulnerabilityIdToJiraLabel:       &bAddVulnerabilityIdToJiraLabel,
			AllowAnonymousSurveyRepsonse:        &bAllowAnonymousSurveyRepsonse,
			ApiExposeErrorDetails:               &bApiExposeErrorDetails,
			DeleteDuplicates:                    &bDeleteDuplicates,
			DisableJiraWebhookSecret:            &bDisableJiraWebhookSecret,
			DisclaimerReportsForced:             &bDisclaimerReportsForced,
			EnableBenchmark:                     &bEnableBenchmark,
			EnableCalendar:                      &bEnableCalendar,
			EnableChecklists:                    &bEnableChecklists,
			EnableCvss3Display:                  &bEnableCvss3Display,
			EnableCvss4Display:                  &bEnableCvss4Display,
			EnableDeduplication:                 &bEnableDeduplication,
			EnableEndpointMetadataImport:        &bEnableEndpointMetadataImport,
			EnableFindingGroups:                 &bEnableFindingGroups,
			EnableFindingSla:                    &bEnableFindingSla,
			EnableGithub:                        &bEnableGithub,
			EnableJira:                          &bEnableJira,
			EnableJiraWebHook:                   &bEnableJiraWebHook,
			EnableMailNotifications:             &bEnableMailNotifications,
			EnableMsteamsNotifications:          &bEnableMsteamsNotifications,
			EnableNotifySlaActive:               &bEnableNotifySlaActive,
			EnableNotifySlaActiveVerified:       &bEnableNotifySlaActiveVerified,
			EnableNotifySlaExponentialBackoff:   &bEnableNotifySlaExponentialBackoff,
			EnableNotifySlaJiraOnly:             &bEnableNotifySlaJiraOnly,
			EnableProductGrade:                  &bEnableProductGrade,
			EnableProductTagInheritance:         &bEnableProductTagInheritance,
			EnableProductTrackingFiles:          &bEnableProductTrackingFiles,
			EnableQuestionnaires:                &bEnableQuestionnaires,
			EnableSimilarFindings:               &bEnableSimilarFindings,
			EnableSlackNotifications:            &bEnableSlackNotifications,
			EnableUiTableBasedSearching:         &bEnableUiTableBasedSearching,
			EnableUserProfileEditable:           &bEnableUserProfileEditable,
			EnableWebhooksNotifications:         &bEnableWebhooksNotifications,
			EnforceVerifiedStatus:               &bEnforceVerifiedStatus,
			EnforceVerifiedStatusJira:           &bEnforceVerifiedStatusJira,
			EnforceVerifiedStatusMetrics:        &bEnforceVerifiedStatusMetrics,
			EnforceVerifiedStatusProductGrading: &bEnforceVerifiedStatusProductGrading,
			EngagementAutoClose:                 &bEngagementAutoClose,
			FalsePositiveHistory:                &bFalsePositiveHistory,
			FilterStringMatching:                &bFilterStringMatching,
			LowercaseCharacterRequired:          &bLowercaseCharacterRequired,
			NonCommonPasswordRequired:           &bNonCommonPasswordRequired,
			NumberCharacterRequired:             &bNumberCharacterRequired,
			RetroactiveFalsePositiveHistory:     &bRetroactiveFalsePositiveHistory,
			SpecialCharacterRequired:            &bSpecialCharacterRequired,
			UppercaseCharacterRequired:          &bUppercaseCharacterRequired,

			DisclaimerNotes:         &sDisclaimerNotes,
			DisclaimerNotifications: &sDisclaimerNotifications,
			DisclaimerReports:       &sDisclaimerReports,
			EmailFrom:               &sEmailFrom,
			JiraLabels:              &sJiraLabels,
			JiraWebhookSecret:       &sJiraWebhookSecret,
			MailNotificationsTo:     &sMailNotificationsTo,
			MsteamsUrl:              &sMsteamsUrl,
			SlackChannel:            &sSlackChannel,
			SlackToken:              &sSlackToken,
			SlackUsername:           &sSlackUsername,
			TeamName:                &sTeamName,
			UrlPrefix:               &sUrlPrefix,

			EngagementAutoCloseDays:              &iEngagementAutoCloseDays,
			MaxDupes:                             &iMaxDupes,
			MaximumPasswordLength:                &iMaximumPasswordLength,
			MinimumPasswordLength:                &iMinimumPasswordLength,
			ProductGradeA:                        &iProductGradeA,
			ProductGradeB:                        &iProductGradeB,
			ProductGradeC:                        &iProductGradeC,
			ProductGradeD:                        &iProductGradeD,
			ProductGradeF:                        &iProductGradeF,
			RiskAcceptanceFormDefaultDays:        &iRiskAcceptanceFormDefaultDays,
			RiskAcceptanceNotifyBeforeExpiration: &iRiskAcceptanceNotifyBeforeExpiration,
			WebhooksNotificationsTimeout:         &iWebhooksNotificationsTimeout,

			JiraMinimumSeverity: &eJiraMinimumSeverity,
		},
	}

	data := systemSettingsResourceData{}
	var terraformResource terraformResourceData = &data
	populateResourceData(context.Background(), &diag.Diagnostics{}, &terraformResource, &ddSettings)

	assert.Equal(t, data.Id.ValueString(), "501")

	assert.Equal(t, data.AddVulnerabilityIdToJiraLabel.ValueBool(), bAddVulnerabilityIdToJiraLabel)
	assert.Equal(t, data.AllowAnonymousSurveyRepsonse.ValueBool(), bAllowAnonymousSurveyRepsonse)
	assert.Equal(t, data.ApiExposeErrorDetails.ValueBool(), bApiExposeErrorDetails)
	assert.Equal(t, data.DeleteDuplicates.ValueBool(), bDeleteDuplicates)
	assert.Equal(t, data.DisableJiraWebhookSecret.ValueBool(), bDisableJiraWebhookSecret)
	assert.Equal(t, data.DisclaimerReportsForced.ValueBool(), bDisclaimerReportsForced)
	assert.Equal(t, data.EnableBenchmark.ValueBool(), bEnableBenchmark)
	assert.Equal(t, data.EnableCalendar.ValueBool(), bEnableCalendar)
	assert.Equal(t, data.EnableChecklists.ValueBool(), bEnableChecklists)
	assert.Equal(t, data.EnableCvss3Display.ValueBool(), bEnableCvss3Display)
	assert.Equal(t, data.EnableCvss4Display.ValueBool(), bEnableCvss4Display)
	assert.Equal(t, data.EnableDeduplication.ValueBool(), bEnableDeduplication)
	assert.Equal(t, data.EnableEndpointMetadataImport.ValueBool(), bEnableEndpointMetadataImport)
	assert.Equal(t, data.EnableFindingGroups.ValueBool(), bEnableFindingGroups)
	assert.Equal(t, data.EnableFindingSla.ValueBool(), bEnableFindingSla)
	assert.Equal(t, data.EnableGithub.ValueBool(), bEnableGithub)
	assert.Equal(t, data.EnableJira.ValueBool(), bEnableJira)
	assert.Equal(t, data.EnableJiraWebHook.ValueBool(), bEnableJiraWebHook)
	assert.Equal(t, data.EnableMailNotifications.ValueBool(), bEnableMailNotifications)
	assert.Equal(t, data.EnableMsteamsNotifications.ValueBool(), bEnableMsteamsNotifications)
	assert.Equal(t, data.EnableNotifySlaActive.ValueBool(), bEnableNotifySlaActive)
	assert.Equal(t, data.EnableNotifySlaActiveVerified.ValueBool(), bEnableNotifySlaActiveVerified)
	assert.Equal(t, data.EnableNotifySlaExponentialBackoff.ValueBool(), bEnableNotifySlaExponentialBackoff)
	assert.Equal(t, data.EnableNotifySlaJiraOnly.ValueBool(), bEnableNotifySlaJiraOnly)
	assert.Equal(t, data.EnableProductGrade.ValueBool(), bEnableProductGrade)
	assert.Equal(t, data.EnableProductTagInheritance.ValueBool(), bEnableProductTagInheritance)
	assert.Equal(t, data.EnableProductTrackingFiles.ValueBool(), bEnableProductTrackingFiles)
	assert.Equal(t, data.EnableQuestionnaires.ValueBool(), bEnableQuestionnaires)
	assert.Equal(t, data.EnableSimilarFindings.ValueBool(), bEnableSimilarFindings)
	assert.Equal(t, data.EnableSlackNotifications.ValueBool(), bEnableSlackNotifications)
	assert.Equal(t, data.EnableUiTableBasedSearching.ValueBool(), bEnableUiTableBasedSearching)
	assert.Equal(t, data.EnableUserProfileEditable.ValueBool(), bEnableUserProfileEditable)
	assert.Equal(t, data.EnableWebhooksNotifications.ValueBool(), bEnableWebhooksNotifications)
	assert.Equal(t, data.EnforceVerifiedStatus.ValueBool(), bEnforceVerifiedStatus)
	assert.Equal(t, data.EnforceVerifiedStatusJira.ValueBool(), bEnforceVerifiedStatusJira)
	assert.Equal(t, data.EnforceVerifiedStatusMetrics.ValueBool(), bEnforceVerifiedStatusMetrics)
	assert.Equal(t, data.EnforceVerifiedStatusProductGrading.ValueBool(), bEnforceVerifiedStatusProductGrading)
	assert.Equal(t, data.EngagementAutoClose.ValueBool(), bEngagementAutoClose)
	assert.Equal(t, data.FalsePositiveHistory.ValueBool(), bFalsePositiveHistory)
	assert.Equal(t, data.FilterStringMatching.ValueBool(), bFilterStringMatching)
	assert.Equal(t, data.LowercaseCharacterRequired.ValueBool(), bLowercaseCharacterRequired)
	assert.Equal(t, data.NonCommonPasswordRequired.ValueBool(), bNonCommonPasswordRequired)
	assert.Equal(t, data.NumberCharacterRequired.ValueBool(), bNumberCharacterRequired)
	assert.Equal(t, data.RetroactiveFalsePositiveHistory.ValueBool(), bRetroactiveFalsePositiveHistory)
	assert.Equal(t, data.SpecialCharacterRequired.ValueBool(), bSpecialCharacterRequired)
	assert.Equal(t, data.UppercaseCharacterRequired.ValueBool(), bUppercaseCharacterRequired)

	assert.Equal(t, data.DisclaimerNotes.ValueString(), sDisclaimerNotes)
	assert.Equal(t, data.DisclaimerNotifications.ValueString(), sDisclaimerNotifications)
	assert.Equal(t, data.DisclaimerReports.ValueString(), sDisclaimerReports)
	assert.Equal(t, data.EmailFrom.ValueString(), sEmailFrom)
	assert.Equal(t, data.JiraLabels.ValueString(), sJiraLabels)
	assert.Equal(t, data.JiraWebhookSecret.ValueString(), sJiraWebhookSecret)
	assert.Equal(t, data.MailNotificationsTo.ValueString(), sMailNotificationsTo)
	assert.Equal(t, data.MsteamsUrl.ValueString(), sMsteamsUrl)
	assert.Equal(t, data.SlackChannel.ValueString(), sSlackChannel)
	assert.Equal(t, data.SlackToken.ValueString(), sSlackToken)
	assert.Equal(t, data.SlackUsername.ValueString(), sSlackUsername)
	assert.Equal(t, data.TeamName.ValueString(), sTeamName)
	assert.Equal(t, data.UrlPrefix.ValueString(), sUrlPrefix)

	assert.Equal(t, data.EngagementAutoCloseDays.ValueInt64(), int64(iEngagementAutoCloseDays))
	assert.Equal(t, data.MaxDupes.ValueInt64(), int64(iMaxDupes))
	assert.Equal(t, data.MaximumPasswordLength.ValueInt64(), int64(iMaximumPasswordLength))
	assert.Equal(t, data.MinimumPasswordLength.ValueInt64(), int64(iMinimumPasswordLength))
	assert.Equal(t, data.ProductGradeA.ValueInt64(), int64(iProductGradeA))
	assert.Equal(t, data.ProductGradeB.ValueInt64(), int64(iProductGradeB))
	assert.Equal(t, data.ProductGradeC.ValueInt64(), int64(iProductGradeC))
	assert.Equal(t, data.ProductGradeD.ValueInt64(), int64(iProductGradeD))
	assert.Equal(t, data.ProductGradeF.ValueInt64(), int64(iProductGradeF))
	assert.Equal(t, data.RiskAcceptanceFormDefaultDays.ValueInt64(), int64(iRiskAcceptanceFormDefaultDays))
	assert.Equal(t, data.RiskAcceptanceNotifyBeforeExpiration.ValueInt64(), int64(iRiskAcceptanceNotifyBeforeExpiration))
	assert.Equal(t, data.WebhooksNotificationsTimeout.ValueInt64(), int64(iWebhooksNotificationsTimeout))

	assert.Equal(t, data.JiraMinimumSeverity.ValueString(), string(eJiraMinimumSeverity))

	// Nil case: every pointer field nil -> every Terraform value null.
	ddSettingsNil := systemSettingsDefectdojoResource{SystemSettings: dd.SystemSettings{}}
	populateResourceData(context.Background(), &diag.Diagnostics{}, &terraformResource, &ddSettingsNil)

	assert.Equal(t, data.Id.IsNull(), true)

	assert.Equal(t, data.AddVulnerabilityIdToJiraLabel.IsNull(), true)
	assert.Equal(t, data.AllowAnonymousSurveyRepsonse.IsNull(), true)
	assert.Equal(t, data.ApiExposeErrorDetails.IsNull(), true)
	assert.Equal(t, data.DeleteDuplicates.IsNull(), true)
	assert.Equal(t, data.DisableJiraWebhookSecret.IsNull(), true)
	assert.Equal(t, data.DisclaimerReportsForced.IsNull(), true)
	assert.Equal(t, data.EnableBenchmark.IsNull(), true)
	assert.Equal(t, data.EnableCalendar.IsNull(), true)
	assert.Equal(t, data.EnableChecklists.IsNull(), true)
	assert.Equal(t, data.EnableCvss3Display.IsNull(), true)
	assert.Equal(t, data.EnableCvss4Display.IsNull(), true)
	assert.Equal(t, data.EnableDeduplication.IsNull(), true)
	assert.Equal(t, data.EnableEndpointMetadataImport.IsNull(), true)
	assert.Equal(t, data.EnableFindingGroups.IsNull(), true)
	assert.Equal(t, data.EnableFindingSla.IsNull(), true)
	assert.Equal(t, data.EnableGithub.IsNull(), true)
	assert.Equal(t, data.EnableJira.IsNull(), true)
	assert.Equal(t, data.EnableJiraWebHook.IsNull(), true)
	assert.Equal(t, data.EnableMailNotifications.IsNull(), true)
	assert.Equal(t, data.EnableMsteamsNotifications.IsNull(), true)
	assert.Equal(t, data.EnableNotifySlaActive.IsNull(), true)
	assert.Equal(t, data.EnableNotifySlaActiveVerified.IsNull(), true)
	assert.Equal(t, data.EnableNotifySlaExponentialBackoff.IsNull(), true)
	assert.Equal(t, data.EnableNotifySlaJiraOnly.IsNull(), true)
	assert.Equal(t, data.EnableProductGrade.IsNull(), true)
	assert.Equal(t, data.EnableProductTagInheritance.IsNull(), true)
	assert.Equal(t, data.EnableProductTrackingFiles.IsNull(), true)
	assert.Equal(t, data.EnableQuestionnaires.IsNull(), true)
	assert.Equal(t, data.EnableSimilarFindings.IsNull(), true)
	assert.Equal(t, data.EnableSlackNotifications.IsNull(), true)
	assert.Equal(t, data.EnableUiTableBasedSearching.IsNull(), true)
	assert.Equal(t, data.EnableUserProfileEditable.IsNull(), true)
	assert.Equal(t, data.EnableWebhooksNotifications.IsNull(), true)
	assert.Equal(t, data.EnforceVerifiedStatus.IsNull(), true)
	assert.Equal(t, data.EnforceVerifiedStatusJira.IsNull(), true)
	assert.Equal(t, data.EnforceVerifiedStatusMetrics.IsNull(), true)
	assert.Equal(t, data.EnforceVerifiedStatusProductGrading.IsNull(), true)
	assert.Equal(t, data.EngagementAutoClose.IsNull(), true)
	assert.Equal(t, data.FalsePositiveHistory.IsNull(), true)
	assert.Equal(t, data.FilterStringMatching.IsNull(), true)
	assert.Equal(t, data.LowercaseCharacterRequired.IsNull(), true)
	assert.Equal(t, data.NonCommonPasswordRequired.IsNull(), true)
	assert.Equal(t, data.NumberCharacterRequired.IsNull(), true)
	assert.Equal(t, data.RetroactiveFalsePositiveHistory.IsNull(), true)
	assert.Equal(t, data.SpecialCharacterRequired.IsNull(), true)
	assert.Equal(t, data.UppercaseCharacterRequired.IsNull(), true)

	assert.Equal(t, data.DisclaimerNotes.IsNull(), true)
	assert.Equal(t, data.DisclaimerNotifications.IsNull(), true)
	assert.Equal(t, data.DisclaimerReports.IsNull(), true)
	assert.Equal(t, data.EmailFrom.IsNull(), true)
	assert.Equal(t, data.JiraLabels.IsNull(), true)
	assert.Equal(t, data.JiraWebhookSecret.IsNull(), true)
	assert.Equal(t, data.MailNotificationsTo.IsNull(), true)
	assert.Equal(t, data.MsteamsUrl.IsNull(), true)
	assert.Equal(t, data.SlackChannel.IsNull(), true)
	assert.Equal(t, data.SlackToken.IsNull(), true)
	assert.Equal(t, data.SlackUsername.IsNull(), true)
	assert.Equal(t, data.TeamName.IsNull(), true)
	assert.Equal(t, data.UrlPrefix.IsNull(), true)

	assert.Equal(t, data.EngagementAutoCloseDays.IsNull(), true)
	assert.Equal(t, data.MaxDupes.IsNull(), true)
	assert.Equal(t, data.MaximumPasswordLength.IsNull(), true)
	assert.Equal(t, data.MinimumPasswordLength.IsNull(), true)
	assert.Equal(t, data.ProductGradeA.IsNull(), true)
	assert.Equal(t, data.ProductGradeB.IsNull(), true)
	assert.Equal(t, data.ProductGradeC.IsNull(), true)
	assert.Equal(t, data.ProductGradeD.IsNull(), true)
	assert.Equal(t, data.ProductGradeF.IsNull(), true)
	assert.Equal(t, data.RiskAcceptanceFormDefaultDays.IsNull(), true)
	assert.Equal(t, data.RiskAcceptanceNotifyBeforeExpiration.IsNull(), true)
	assert.Equal(t, data.WebhooksNotificationsTimeout.IsNull(), true)

	assert.Equal(t, data.JiraMinimumSeverity.IsNull(), true)
}

func TestSystemSettingsResource__defectdojoResource(t *testing.T) {
	data := systemSettingsResourceData{
		AddVulnerabilityIdToJiraLabel:       types.BoolValue(true),
		AllowAnonymousSurveyRepsonse:        types.BoolValue(false),
		ApiExposeErrorDetails:               types.BoolValue(true),
		DeleteDuplicates:                    types.BoolValue(false),
		DisableJiraWebhookSecret:            types.BoolValue(true),
		DisclaimerReportsForced:             types.BoolValue(false),
		EnableBenchmark:                     types.BoolValue(true),
		EnableCalendar:                      types.BoolValue(false),
		EnableChecklists:                    types.BoolValue(true),
		EnableCvss3Display:                  types.BoolValue(false),
		EnableCvss4Display:                  types.BoolValue(true),
		EnableDeduplication:                 types.BoolValue(false),
		EnableEndpointMetadataImport:        types.BoolValue(true),
		EnableFindingGroups:                 types.BoolValue(false),
		EnableFindingSla:                    types.BoolValue(true),
		EnableGithub:                        types.BoolValue(false),
		EnableJira:                          types.BoolValue(true),
		EnableJiraWebHook:                   types.BoolValue(false),
		EnableMailNotifications:             types.BoolValue(true),
		EnableMsteamsNotifications:          types.BoolValue(false),
		EnableNotifySlaActive:               types.BoolValue(true),
		EnableNotifySlaActiveVerified:       types.BoolValue(false),
		EnableNotifySlaExponentialBackoff:   types.BoolValue(true),
		EnableNotifySlaJiraOnly:             types.BoolValue(false),
		EnableProductGrade:                  types.BoolValue(true),
		EnableProductTagInheritance:         types.BoolValue(false),
		EnableProductTrackingFiles:          types.BoolValue(true),
		EnableQuestionnaires:                types.BoolValue(false),
		EnableSimilarFindings:               types.BoolValue(true),
		EnableSlackNotifications:            types.BoolValue(false),
		EnableUiTableBasedSearching:         types.BoolValue(true),
		EnableUserProfileEditable:           types.BoolValue(false),
		EnableWebhooksNotifications:         types.BoolValue(true),
		EnforceVerifiedStatus:               types.BoolValue(false),
		EnforceVerifiedStatusJira:           types.BoolValue(true),
		EnforceVerifiedStatusMetrics:        types.BoolValue(false),
		EnforceVerifiedStatusProductGrading: types.BoolValue(true),
		EngagementAutoClose:                 types.BoolValue(false),
		FalsePositiveHistory:                types.BoolValue(true),
		FilterStringMatching:                types.BoolValue(false),
		LowercaseCharacterRequired:          types.BoolValue(true),
		NonCommonPasswordRequired:           types.BoolValue(false),
		NumberCharacterRequired:             types.BoolValue(true),
		RetroactiveFalsePositiveHistory:     types.BoolValue(false),
		SpecialCharacterRequired:            types.BoolValue(true),
		UppercaseCharacterRequired:          types.BoolValue(false),

		DisclaimerNotes:         types.StringValue("notes disclaimer"),
		DisclaimerNotifications: types.StringValue("notifications disclaimer"),
		DisclaimerReports:       types.StringValue("reports disclaimer"),
		EmailFrom:               types.StringValue("noreply@example.com"),
		JiraLabels:              types.StringValue("label1 label2"),
		JiraWebhookSecret:       types.StringValue("shh-secret"),
		MailNotificationsTo:     types.StringValue("admin@example.com"),
		MsteamsUrl:              types.StringValue("https://example.com/msteams"),
		SlackChannel:            types.StringValue("#security"),
		SlackToken:              types.StringValue("xoxb-token"),
		SlackUsername:           types.StringValue("defectdojo-bot"),
		TeamName:                types.StringValue("Acme Security"),
		UrlPrefix:               types.StringValue("/dojo"),

		EngagementAutoCloseDays:              types.Int64Value(3),
		MaxDupes:                             types.Int64Value(5),
		MaximumPasswordLength:                types.Int64Value(64),
		MinimumPasswordLength:                types.Int64Value(8),
		ProductGradeA:                        types.Int64Value(90),
		ProductGradeB:                        types.Int64Value(80),
		ProductGradeC:                        types.Int64Value(70),
		ProductGradeD:                        types.Int64Value(60),
		ProductGradeF:                        types.Int64Value(59),
		RiskAcceptanceFormDefaultDays:        types.Int64Value(180),
		RiskAcceptanceNotifyBeforeExpiration: types.Int64Value(10),
		WebhooksNotificationsTimeout:         types.Int64Value(30),

		JiraMinimumSeverity: types.StringValue("Critical"),
	}

	ddResource := data.defectdojoResource()
	ddSettings := ddResource.(*systemSettingsDefectdojoResource)
	var terraformResource terraformResourceData = &data
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, terraformResource, &ddResource)

	assert.Equal(t, *ddSettings.AddVulnerabilityIdToJiraLabel, true)
	assert.Equal(t, *ddSettings.AllowAnonymousSurveyRepsonse, false)
	assert.Equal(t, *ddSettings.ApiExposeErrorDetails, true)
	assert.Equal(t, *ddSettings.DeleteDuplicates, false)
	assert.Equal(t, *ddSettings.DisableJiraWebhookSecret, true)
	assert.Equal(t, *ddSettings.DisclaimerReportsForced, false)
	assert.Equal(t, *ddSettings.EnableBenchmark, true)
	assert.Equal(t, *ddSettings.EnableCalendar, false)
	assert.Equal(t, *ddSettings.EnableChecklists, true)
	assert.Equal(t, *ddSettings.EnableCvss3Display, false)
	assert.Equal(t, *ddSettings.EnableCvss4Display, true)
	assert.Equal(t, *ddSettings.EnableDeduplication, false)
	assert.Equal(t, *ddSettings.EnableEndpointMetadataImport, true)
	assert.Equal(t, *ddSettings.EnableFindingGroups, false)
	assert.Equal(t, *ddSettings.EnableFindingSla, true)
	assert.Equal(t, *ddSettings.EnableGithub, false)
	assert.Equal(t, *ddSettings.EnableJira, true)
	assert.Equal(t, *ddSettings.EnableJiraWebHook, false)
	assert.Equal(t, *ddSettings.EnableMailNotifications, true)
	assert.Equal(t, *ddSettings.EnableMsteamsNotifications, false)
	assert.Equal(t, *ddSettings.EnableNotifySlaActive, true)
	assert.Equal(t, *ddSettings.EnableNotifySlaActiveVerified, false)
	assert.Equal(t, *ddSettings.EnableNotifySlaExponentialBackoff, true)
	assert.Equal(t, *ddSettings.EnableNotifySlaJiraOnly, false)
	assert.Equal(t, *ddSettings.EnableProductGrade, true)
	assert.Equal(t, *ddSettings.EnableProductTagInheritance, false)
	assert.Equal(t, *ddSettings.EnableProductTrackingFiles, true)
	assert.Equal(t, *ddSettings.EnableQuestionnaires, false)
	assert.Equal(t, *ddSettings.EnableSimilarFindings, true)
	assert.Equal(t, *ddSettings.EnableSlackNotifications, false)
	assert.Equal(t, *ddSettings.EnableUiTableBasedSearching, true)
	assert.Equal(t, *ddSettings.EnableUserProfileEditable, false)
	assert.Equal(t, *ddSettings.EnableWebhooksNotifications, true)
	assert.Equal(t, *ddSettings.EnforceVerifiedStatus, false)
	assert.Equal(t, *ddSettings.EnforceVerifiedStatusJira, true)
	assert.Equal(t, *ddSettings.EnforceVerifiedStatusMetrics, false)
	assert.Equal(t, *ddSettings.EnforceVerifiedStatusProductGrading, true)
	assert.Equal(t, *ddSettings.EngagementAutoClose, false)
	assert.Equal(t, *ddSettings.FalsePositiveHistory, true)
	assert.Equal(t, *ddSettings.FilterStringMatching, false)
	assert.Equal(t, *ddSettings.LowercaseCharacterRequired, true)
	assert.Equal(t, *ddSettings.NonCommonPasswordRequired, false)
	assert.Equal(t, *ddSettings.NumberCharacterRequired, true)
	assert.Equal(t, *ddSettings.RetroactiveFalsePositiveHistory, false)
	assert.Equal(t, *ddSettings.SpecialCharacterRequired, true)
	assert.Equal(t, *ddSettings.UppercaseCharacterRequired, false)

	assert.Equal(t, *ddSettings.DisclaimerNotes, "notes disclaimer")
	assert.Equal(t, *ddSettings.DisclaimerNotifications, "notifications disclaimer")
	assert.Equal(t, *ddSettings.DisclaimerReports, "reports disclaimer")
	assert.Equal(t, *ddSettings.EmailFrom, "noreply@example.com")
	assert.Equal(t, *ddSettings.JiraLabels, "label1 label2")
	assert.Equal(t, *ddSettings.JiraWebhookSecret, "shh-secret")
	assert.Equal(t, *ddSettings.MailNotificationsTo, "admin@example.com")
	assert.Equal(t, *ddSettings.MsteamsUrl, "https://example.com/msteams")
	assert.Equal(t, *ddSettings.SlackChannel, "#security")
	assert.Equal(t, *ddSettings.SlackToken, "xoxb-token")
	assert.Equal(t, *ddSettings.SlackUsername, "defectdojo-bot")
	assert.Equal(t, *ddSettings.TeamName, "Acme Security")
	assert.Equal(t, *ddSettings.UrlPrefix, "/dojo")

	assert.Equal(t, *ddSettings.EngagementAutoCloseDays, 3)
	assert.Equal(t, *ddSettings.MaxDupes, 5)
	assert.Equal(t, *ddSettings.MaximumPasswordLength, 64)
	assert.Equal(t, *ddSettings.MinimumPasswordLength, 8)
	assert.Equal(t, *ddSettings.ProductGradeA, 90)
	assert.Equal(t, *ddSettings.ProductGradeB, 80)
	assert.Equal(t, *ddSettings.ProductGradeC, 70)
	assert.Equal(t, *ddSettings.ProductGradeD, 60)
	assert.Equal(t, *ddSettings.ProductGradeF, 59)
	assert.Equal(t, *ddSettings.RiskAcceptanceFormDefaultDays, 180)
	assert.Equal(t, *ddSettings.RiskAcceptanceNotifyBeforeExpiration, 10)
	assert.Equal(t, *ddSettings.WebhooksNotificationsTimeout, 30)

	// enum cast: Terraform string -> dd.SystemSettingsJiraMinimumSeverity
	assert.Equal(t, string(*ddSettings.JiraMinimumSeverity), "Critical")
	assert.Equal(t, *ddSettings.JiraMinimumSeverity, dd.SystemSettingsJiraMinimumSeverityCritical)

	// systemSettingsToRequest: verify the PATCH request body carries the
	// same values, including the second enum cast
	// (SystemSettingsJiraMinimumSeverity -> PatchedSystemSettingsRequestJiraMinimumSeverity).
	req := systemSettingsToRequest(ddSettings.SystemSettings)
	assert.Equal(t, *req.AddVulnerabilityIdToJiraLabel, true)
	assert.Equal(t, *req.TeamName, "Acme Security")
	assert.Equal(t, *req.MaxDupes, 5)
	assert.Equal(t, *req.JiraMinimumSeverity, dd.PatchedSystemSettingsRequestJiraMinimumSeverityCritical)
}

func TestSystemSettingsResource__defectdojoResource_Nulls(t *testing.T) {
	data := systemSettingsResourceData{
		Id: types.StringNull(),

		AddVulnerabilityIdToJiraLabel:       types.BoolNull(),
		AllowAnonymousSurveyRepsonse:        types.BoolNull(),
		ApiExposeErrorDetails:               types.BoolNull(),
		DeleteDuplicates:                    types.BoolNull(),
		DisableJiraWebhookSecret:            types.BoolNull(),
		DisclaimerReportsForced:             types.BoolNull(),
		EnableBenchmark:                     types.BoolNull(),
		EnableCalendar:                      types.BoolNull(),
		EnableChecklists:                    types.BoolNull(),
		EnableCvss3Display:                  types.BoolNull(),
		EnableCvss4Display:                  types.BoolNull(),
		EnableDeduplication:                 types.BoolNull(),
		EnableEndpointMetadataImport:        types.BoolNull(),
		EnableFindingGroups:                 types.BoolNull(),
		EnableFindingSla:                    types.BoolNull(),
		EnableGithub:                        types.BoolNull(),
		EnableJira:                          types.BoolNull(),
		EnableJiraWebHook:                   types.BoolNull(),
		EnableMailNotifications:             types.BoolNull(),
		EnableMsteamsNotifications:          types.BoolNull(),
		EnableNotifySlaActive:               types.BoolNull(),
		EnableNotifySlaActiveVerified:       types.BoolNull(),
		EnableNotifySlaExponentialBackoff:   types.BoolNull(),
		EnableNotifySlaJiraOnly:             types.BoolNull(),
		EnableProductGrade:                  types.BoolNull(),
		EnableProductTagInheritance:         types.BoolNull(),
		EnableProductTrackingFiles:          types.BoolNull(),
		EnableQuestionnaires:                types.BoolNull(),
		EnableSimilarFindings:               types.BoolNull(),
		EnableSlackNotifications:            types.BoolNull(),
		EnableUiTableBasedSearching:         types.BoolNull(),
		EnableUserProfileEditable:           types.BoolNull(),
		EnableWebhooksNotifications:         types.BoolNull(),
		EnforceVerifiedStatus:               types.BoolNull(),
		EnforceVerifiedStatusJira:           types.BoolNull(),
		EnforceVerifiedStatusMetrics:        types.BoolNull(),
		EnforceVerifiedStatusProductGrading: types.BoolNull(),
		EngagementAutoClose:                 types.BoolNull(),
		FalsePositiveHistory:                types.BoolNull(),
		FilterStringMatching:                types.BoolNull(),
		LowercaseCharacterRequired:          types.BoolNull(),
		NonCommonPasswordRequired:           types.BoolNull(),
		NumberCharacterRequired:             types.BoolNull(),
		RetroactiveFalsePositiveHistory:     types.BoolNull(),
		SpecialCharacterRequired:            types.BoolNull(),
		UppercaseCharacterRequired:          types.BoolNull(),

		DisclaimerNotes:         types.StringNull(),
		DisclaimerNotifications: types.StringNull(),
		DisclaimerReports:       types.StringNull(),
		EmailFrom:               types.StringNull(),
		JiraLabels:              types.StringNull(),
		JiraWebhookSecret:       types.StringNull(),
		MailNotificationsTo:     types.StringNull(),
		MsteamsUrl:              types.StringNull(),
		SlackChannel:            types.StringNull(),
		SlackToken:              types.StringNull(),
		SlackUsername:           types.StringNull(),
		TeamName:                types.StringNull(),
		UrlPrefix:               types.StringNull(),

		EngagementAutoCloseDays:              types.Int64Null(),
		MaxDupes:                             types.Int64Null(),
		MaximumPasswordLength:                types.Int64Null(),
		MinimumPasswordLength:                types.Int64Null(),
		ProductGradeA:                        types.Int64Null(),
		ProductGradeB:                        types.Int64Null(),
		ProductGradeC:                        types.Int64Null(),
		ProductGradeD:                        types.Int64Null(),
		ProductGradeF:                        types.Int64Null(),
		RiskAcceptanceFormDefaultDays:        types.Int64Null(),
		RiskAcceptanceNotifyBeforeExpiration: types.Int64Null(),
		WebhooksNotificationsTimeout:         types.Int64Null(),

		JiraMinimumSeverity: types.StringNull(),
	}

	ddResource := data.defectdojoResource()
	ddSettings := ddResource.(*systemSettingsDefectdojoResource)
	var terraformResource terraformResourceData = &data
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, terraformResource, &ddResource)

	var nilBool *bool
	var nilString *string
	var nilInt *int
	var nilEnum *dd.SystemSettingsJiraMinimumSeverity

	assert.Equal(t, ddSettings.AddVulnerabilityIdToJiraLabel, nilBool)
	assert.Equal(t, ddSettings.AllowAnonymousSurveyRepsonse, nilBool)
	assert.Equal(t, ddSettings.ApiExposeErrorDetails, nilBool)
	assert.Equal(t, ddSettings.DeleteDuplicates, nilBool)
	assert.Equal(t, ddSettings.DisableJiraWebhookSecret, nilBool)
	assert.Equal(t, ddSettings.DisclaimerReportsForced, nilBool)
	assert.Equal(t, ddSettings.EnableBenchmark, nilBool)
	assert.Equal(t, ddSettings.EnableCalendar, nilBool)
	assert.Equal(t, ddSettings.EnableChecklists, nilBool)
	assert.Equal(t, ddSettings.EnableCvss3Display, nilBool)
	assert.Equal(t, ddSettings.EnableCvss4Display, nilBool)
	assert.Equal(t, ddSettings.EnableDeduplication, nilBool)
	assert.Equal(t, ddSettings.EnableEndpointMetadataImport, nilBool)
	assert.Equal(t, ddSettings.EnableFindingGroups, nilBool)
	assert.Equal(t, ddSettings.EnableFindingSla, nilBool)
	assert.Equal(t, ddSettings.EnableGithub, nilBool)
	assert.Equal(t, ddSettings.EnableJira, nilBool)
	assert.Equal(t, ddSettings.EnableJiraWebHook, nilBool)
	assert.Equal(t, ddSettings.EnableMailNotifications, nilBool)
	assert.Equal(t, ddSettings.EnableMsteamsNotifications, nilBool)
	assert.Equal(t, ddSettings.EnableNotifySlaActive, nilBool)
	assert.Equal(t, ddSettings.EnableNotifySlaActiveVerified, nilBool)
	assert.Equal(t, ddSettings.EnableNotifySlaExponentialBackoff, nilBool)
	assert.Equal(t, ddSettings.EnableNotifySlaJiraOnly, nilBool)
	assert.Equal(t, ddSettings.EnableProductGrade, nilBool)
	assert.Equal(t, ddSettings.EnableProductTagInheritance, nilBool)
	assert.Equal(t, ddSettings.EnableProductTrackingFiles, nilBool)
	assert.Equal(t, ddSettings.EnableQuestionnaires, nilBool)
	assert.Equal(t, ddSettings.EnableSimilarFindings, nilBool)
	assert.Equal(t, ddSettings.EnableSlackNotifications, nilBool)
	assert.Equal(t, ddSettings.EnableUiTableBasedSearching, nilBool)
	assert.Equal(t, ddSettings.EnableUserProfileEditable, nilBool)
	assert.Equal(t, ddSettings.EnableWebhooksNotifications, nilBool)
	assert.Equal(t, ddSettings.EnforceVerifiedStatus, nilBool)
	assert.Equal(t, ddSettings.EnforceVerifiedStatusJira, nilBool)
	assert.Equal(t, ddSettings.EnforceVerifiedStatusMetrics, nilBool)
	assert.Equal(t, ddSettings.EnforceVerifiedStatusProductGrading, nilBool)
	assert.Equal(t, ddSettings.EngagementAutoClose, nilBool)
	assert.Equal(t, ddSettings.FalsePositiveHistory, nilBool)
	assert.Equal(t, ddSettings.FilterStringMatching, nilBool)
	assert.Equal(t, ddSettings.LowercaseCharacterRequired, nilBool)
	assert.Equal(t, ddSettings.NonCommonPasswordRequired, nilBool)
	assert.Equal(t, ddSettings.NumberCharacterRequired, nilBool)
	assert.Equal(t, ddSettings.RetroactiveFalsePositiveHistory, nilBool)
	assert.Equal(t, ddSettings.SpecialCharacterRequired, nilBool)
	assert.Equal(t, ddSettings.UppercaseCharacterRequired, nilBool)

	assert.Equal(t, ddSettings.DisclaimerNotes, nilString)
	assert.Equal(t, ddSettings.DisclaimerNotifications, nilString)
	assert.Equal(t, ddSettings.DisclaimerReports, nilString)
	assert.Equal(t, ddSettings.EmailFrom, nilString)
	assert.Equal(t, ddSettings.JiraLabels, nilString)
	assert.Equal(t, ddSettings.JiraWebhookSecret, nilString)
	assert.Equal(t, ddSettings.MailNotificationsTo, nilString)
	assert.Equal(t, ddSettings.MsteamsUrl, nilString)
	assert.Equal(t, ddSettings.SlackChannel, nilString)
	assert.Equal(t, ddSettings.SlackToken, nilString)
	assert.Equal(t, ddSettings.SlackUsername, nilString)
	assert.Equal(t, ddSettings.TeamName, nilString)
	assert.Equal(t, ddSettings.UrlPrefix, nilString)

	assert.Equal(t, ddSettings.EngagementAutoCloseDays, nilInt)
	assert.Equal(t, ddSettings.MaxDupes, nilInt)
	assert.Equal(t, ddSettings.MaximumPasswordLength, nilInt)
	assert.Equal(t, ddSettings.MinimumPasswordLength, nilInt)
	assert.Equal(t, ddSettings.ProductGradeA, nilInt)
	assert.Equal(t, ddSettings.ProductGradeB, nilInt)
	assert.Equal(t, ddSettings.ProductGradeC, nilInt)
	assert.Equal(t, ddSettings.ProductGradeD, nilInt)
	assert.Equal(t, ddSettings.ProductGradeF, nilInt)
	assert.Equal(t, ddSettings.RiskAcceptanceFormDefaultDays, nilInt)
	assert.Equal(t, ddSettings.RiskAcceptanceNotifyBeforeExpiration, nilInt)
	assert.Equal(t, ddSettings.WebhooksNotificationsTimeout, nilInt)

	assert.Equal(t, ddSettings.JiraMinimumSeverity, nilEnum)

	// createApiCall/deleteApiCall must return the sentinel singleton error.
	_, _, err := ddSettings.createApiCall(context.Background(), nil)
	assert.ErrorContains(t, err, "singleton")
	_, _, err = ddSettings.deleteApiCall(context.Background(), nil, 1)
	assert.ErrorContains(t, err, "singleton")
}
