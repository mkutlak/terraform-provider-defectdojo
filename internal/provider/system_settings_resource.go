package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
)

// errSystemSettingsSingleton is returned by createApiCall/deleteApiCall,
// which the engine never actually invokes for a singletonAdopter resource,
// but which must exist to satisfy the defectdojoResource interface.
var errSystemSettingsSingleton = errors.New("system settings is a singleton managed by DefectDojo; it is adopted on create and cannot be deleted")

func (t systemSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo System Settings.\n\n" +
			"**Singleton**: `terraform apply` adopts the existing system settings row and updates it. " +
			"`terraform destroy` only removes it from Terraform state — settings keep their last-applied values. " +
			"Only attributes set in configuration are managed; unset attributes retain their server-side values.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"add_vulnerability_id_to_jira_label": schema.BoolAttribute{
				MarkdownDescription: "Add the vulnerability ID as a label when pushing findings to Jira.",
				Optional:            true,
				Computed:            true,
			},
			"allow_anonymous_survey_repsonse": schema.BoolAttribute{
				MarkdownDescription: "Enable anyone with a link to the survey to answer a survey.",
				Optional:            true,
				Computed:            true,
			},
			"api_expose_error_details": schema.BoolAttribute{
				MarkdownDescription: "When turned on, the API will expose error details in the response.",
				Optional:            true,
				Computed:            true,
			},
			"delete_duplicates": schema.BoolAttribute{
				MarkdownDescription: "Requires next setting: maximum number of duplicates to retain.",
				Optional:            true,
				Computed:            true,
			},
			"disable_jira_webhook_secret": schema.BoolAttribute{
				MarkdownDescription: "Allows incoming requests without a secret (discouraged legacy behaviour).",
				Optional:            true,
				Computed:            true,
			},
			"disclaimer_notes": schema.StringAttribute{
				MarkdownDescription: "Include this custom disclaimer next to input form for notes.",
				Optional:            true,
				Computed:            true,
			},
			"disclaimer_notifications": schema.StringAttribute{
				MarkdownDescription: "Include this custom disclaimer on all notifications.",
				Optional:            true,
				Computed:            true,
			},
			"disclaimer_reports": schema.StringAttribute{
				MarkdownDescription: "Include this custom disclaimer on generated reports.",
				Optional:            true,
				Computed:            true,
			},
			"disclaimer_reports_forced": schema.BoolAttribute{
				MarkdownDescription: "Disclaimer will be added to all reports even if user didn't select 'Include disclaimer'.",
				Optional:            true,
				Computed:            true,
			},
			"email_from": schema.StringAttribute{
				MarkdownDescription: "The From address used for outgoing DefectDojo email notifications.",
				Optional:            true,
				Computed:            true,
			},
			"enable_benchmark": schema.BoolAttribute{
				MarkdownDescription: "Enables Benchmarks such as the OWASP ASVS (Application Security Verification Standard).",
				Optional:            true,
				Computed:            true,
			},
			"enable_calendar": schema.BoolAttribute{
				MarkdownDescription: "With this setting turned off, the Calendar will be disabled in the user interface.",
				Optional:            true,
				Computed:            true,
			},
			"enable_checklists": schema.BoolAttribute{
				MarkdownDescription: "With this setting turned off, checklists will be disabled in the user interface.",
				Optional:            true,
				Computed:            true,
			},
			"enable_cvss3_display": schema.BoolAttribute{
				MarkdownDescription: "With this setting turned off, CVSS3 fields will be hidden in the user interface.",
				Optional:            true,
				Computed:            true,
			},
			"enable_cvss4_display": schema.BoolAttribute{
				MarkdownDescription: "With this setting turned off, CVSS4 fields will be hidden in the user interface.",
				Optional:            true,
				Computed:            true,
			},
			"enable_deduplication": schema.BoolAttribute{
				MarkdownDescription: "With this setting turned on, DefectDojo deduplicates findings by comparing endpoints, cwe fields, and titles. If two findings share a URL and have the same CWE or title, DefectDojo marks the recent finding as a duplicate. When deduplication is enabled, a list of deduplicated findings is added to the engagement view.",
				Optional:            true,
				Computed:            true,
			},
			"enable_endpoint_metadata_import": schema.BoolAttribute{
				MarkdownDescription: "With this setting turned off, endpoint metadata import will be disabled in the user interface.",
				Optional:            true,
				Computed:            true,
			},
			"enable_finding_groups": schema.BoolAttribute{
				MarkdownDescription: "With this setting turned off, the Finding Groups will be disabled.",
				Optional:            true,
				Computed:            true,
			},
			"enable_finding_sla": schema.BoolAttribute{
				MarkdownDescription: "Enables Finding SLA's for time to remediate.",
				Optional:            true,
				Computed:            true,
			},
			"enable_github": schema.BoolAttribute{
				MarkdownDescription: "Enables linking DefectDojo findings to Github issues.",
				Optional:            true,
				Computed:            true,
			},
			"enable_jira": schema.BoolAttribute{
				MarkdownDescription: "Enables JIRA integration.",
				Optional:            true,
				Computed:            true,
			},
			"enable_jira_web_hook": schema.BoolAttribute{
				MarkdownDescription: "Please note: It is strongly recommended to use a secret and / or IP whitelist the JIRA server using a proxy such as Nginx.",
				Optional:            true,
				Computed:            true,
			},
			"enable_mail_notifications": schema.BoolAttribute{
				MarkdownDescription: "Enables email notifications.",
				Optional:            true,
				Computed:            true,
			},
			"enable_msteams_notifications": schema.BoolAttribute{
				MarkdownDescription: "Enables Microsoft Teams notifications.",
				Optional:            true,
				Computed:            true,
			},
			"enable_notify_sla_active": schema.BoolAttribute{
				MarkdownDescription: "Enables Notify when time to remediate according to Finding SLA's is breached for active Findings.",
				Optional:            true,
				Computed:            true,
			},
			"enable_notify_sla_active_verified": schema.BoolAttribute{
				MarkdownDescription: "Enables Notify when time to remediate according to Finding SLA's is breached for active, verified Findings.",
				Optional:            true,
				Computed:            true,
			},
			"enable_notify_sla_exponential_backoff": schema.BoolAttribute{
				MarkdownDescription: "Enable an exponential backoff strategy for SLA breach notifications, e.g. 1, 2, 4, 8, etc. Otherwise it alerts every day.",
				Optional:            true,
				Computed:            true,
			},
			"enable_notify_sla_jira_only": schema.BoolAttribute{
				MarkdownDescription: "Enables Notify when time to remediate according to Finding SLA's is breached for Findings that are linked to JIRA issues. Notification is disabled for Findings not linked to JIRA issues.",
				Optional:            true,
				Computed:            true,
			},
			"enable_product_grade": schema.BoolAttribute{
				MarkdownDescription: "Displays a grade letter next to a product to show the overall health.",
				Optional:            true,
				Computed:            true,
			},
			"enable_product_tag_inheritance": schema.BoolAttribute{
				MarkdownDescription: "Enables product tag inheritance globally for all products. Any tags added on a product will automatically be added to all Engagements, Tests, and Findings.",
				Optional:            true,
				Computed:            true,
			},
			"enable_product_tracking_files": schema.BoolAttribute{
				MarkdownDescription: "With this setting turned off, the product tracking files will be disabled in the user interface.",
				Optional:            true,
				Computed:            true,
			},
			"enable_questionnaires": schema.BoolAttribute{
				MarkdownDescription: "With this setting turned off, questionnaires will be disabled in the user interface.",
				Optional:            true,
				Computed:            true,
			},
			"enable_similar_findings": schema.BoolAttribute{
				MarkdownDescription: "Enable the query of similar findings on the view finding page. This feature can involve potentially large queries and negatively impact performance.",
				Optional:            true,
				Computed:            true,
			},
			"enable_slack_notifications": schema.BoolAttribute{
				MarkdownDescription: "Enables Slack notifications.",
				Optional:            true,
				Computed:            true,
			},
			"enable_ui_table_based_searching": schema.BoolAttribute{
				MarkdownDescription: "With this setting enabled, table headings will contain sort buttons for the current page of data in addition to sorting buttons that consider data from all pages.",
				Optional:            true,
				Computed:            true,
			},
			"enable_user_profile_editable": schema.BoolAttribute{
				MarkdownDescription: "When turned on users can edit their profiles.",
				Optional:            true,
				Computed:            true,
			},
			"enable_webhooks_notifications": schema.BoolAttribute{
				MarkdownDescription: "Enables outgoing webhook notifications.",
				Optional:            true,
				Computed:            true,
			},
			"enforce_verified_status": schema.BoolAttribute{
				MarkdownDescription: "When enabled, features such as product grading, jira integration, metrics, and reports will only interact with verified findings. This setting will override individually scoped verified toggles.",
				Optional:            true,
				Computed:            true,
			},
			"enforce_verified_status_jira": schema.BoolAttribute{
				MarkdownDescription: "When enabled, findings must have a verified status to be pushed to jira.",
				Optional:            true,
				Computed:            true,
			},
			"enforce_verified_status_metrics": schema.BoolAttribute{
				MarkdownDescription: "When enabled, findings must have a verified status to be counted in metric calculations, be included in reports, and filters.",
				Optional:            true,
				Computed:            true,
			},
			"enforce_verified_status_product_grading": schema.BoolAttribute{
				MarkdownDescription: "When enabled, findings must have a verified status to be considered as part of a product's grading.",
				Optional:            true,
				Computed:            true,
			},
			"engagement_auto_close": schema.BoolAttribute{
				MarkdownDescription: "Closes an engagement after 3 days (default) past due date including last update.",
				Optional:            true,
				Computed:            true,
			},
			"engagement_auto_close_days": schema.Int64Attribute{
				MarkdownDescription: "Closes an engagement after the specified number of days past due date including last update.",
				Optional:            true,
				Computed:            true,
			},
			"false_positive_history": schema.BoolAttribute{
				MarkdownDescription: "(EXPERIMENTAL) DefectDojo will automatically mark the finding as a false positive if an equal finding (according to its dedupe algorithm) has been previously marked as a false positive on the same product. ATTENTION: Although the deduplication algorithm is used to determine if a finding should be marked as a false positive, this feature will not work if deduplication is enabled since it doesn't make sense to use both.",
				Optional:            true,
				Computed:            true,
			},
			"filter_string_matching": schema.BoolAttribute{
				MarkdownDescription: "When turned on, all filter operations in the UI will require string matches rather than ID. This is a performance enhancement to avoid fetching objects unnecessarily.",
				Optional:            true,
				Computed:            true,
			},
			"jira_labels": schema.StringAttribute{
				MarkdownDescription: "JIRA issue labels space seperated.",
				Optional:            true,
				Computed:            true,
			},
			"jira_minimum_severity": schema.StringAttribute{
				MarkdownDescription: "Minimum severity level for pushing findings to JIRA. Valid values are: 'Critical', 'High', 'Medium', 'Low', 'Info'.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("Critical", "High", "Medium", "Low", "Info"),
				},
			},
			"jira_webhook_secret": schema.StringAttribute{
				MarkdownDescription: "Secret needed in URL for incoming JIRA Webhook.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"lowercase_character_required": schema.BoolAttribute{
				MarkdownDescription: "Requires user passwords to contain at least one lowercase letter (a-z).",
				Optional:            true,
				Computed:            true,
			},
			"mail_notifications_to": schema.StringAttribute{
				MarkdownDescription: "Recipient address for global mail notifications.",
				Optional:            true,
				Computed:            true,
			},
			"max_dupes": schema.Int64Attribute{
				MarkdownDescription: "When enabled, if a single issue reaches the maximum number of duplicates, the oldest will be deleted. Duplicate will not be deleted when left empty. A value of 0 will remove all duplicates.",
				Optional:            true,
				Computed:            true,
			},
			"maximum_password_length": schema.Int64Attribute{
				MarkdownDescription: "Requires user to set passwords less than maximum length.",
				Optional:            true,
				Computed:            true,
			},
			"minimum_password_length": schema.Int64Attribute{
				MarkdownDescription: "Requires user to set passwords greater than minimum length.",
				Optional:            true,
				Computed:            true,
			},
			"msteams_url": schema.StringAttribute{
				MarkdownDescription: "The full URL of the incoming webhook.",
				Optional:            true,
				Computed:            true,
			},
			"non_common_password_required": schema.BoolAttribute{
				MarkdownDescription: "Requires user passwords to not be part of list of common passwords.",
				Optional:            true,
				Computed:            true,
			},
			"number_character_required": schema.BoolAttribute{
				MarkdownDescription: "Requires user passwords to contain at least one digit (0-9).",
				Optional:            true,
				Computed:            true,
			},
			"product_grade_a": schema.Int64Attribute{
				MarkdownDescription: "Percentage score for an 'A' >=.",
				Optional:            true,
				Computed:            true,
			},
			"product_grade_b": schema.Int64Attribute{
				MarkdownDescription: "Percentage score for a 'B' >=.",
				Optional:            true,
				Computed:            true,
			},
			"product_grade_c": schema.Int64Attribute{
				MarkdownDescription: "Percentage score for a 'C' >=.",
				Optional:            true,
				Computed:            true,
			},
			"product_grade_d": schema.Int64Attribute{
				MarkdownDescription: "Percentage score for a 'D' >=.",
				Optional:            true,
				Computed:            true,
			},
			"product_grade_f": schema.Int64Attribute{
				MarkdownDescription: "Percentage score for an 'F' <=.",
				Optional:            true,
				Computed:            true,
			},
			"retroactive_false_positive_history": schema.BoolAttribute{
				MarkdownDescription: "(EXPERIMENTAL) FP History will also retroactively mark/unmark all existing equal findings in the same product as a false positives. Only works if the False Positive History feature is also enabled.",
				Optional:            true,
				Computed:            true,
			},
			"risk_acceptance_form_default_days": schema.Int64Attribute{
				MarkdownDescription: "Default expiry period for risk acceptance form.",
				Optional:            true,
				Computed:            true,
			},
			"risk_acceptance_notify_before_expiration": schema.Int64Attribute{
				MarkdownDescription: "Notify X days before risk acceptance expires. Leave empty to disable.",
				Optional:            true,
				Computed:            true,
			},
			"slack_channel": schema.StringAttribute{
				MarkdownDescription: "Optional. Needed if you want to send global notifications.",
				Optional:            true,
				Computed:            true,
			},
			"slack_token": schema.StringAttribute{
				MarkdownDescription: "Token required for interacting with Slack. Get one at https://api.slack.com/tokens.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"slack_username": schema.StringAttribute{
				MarkdownDescription: "Optional. Will take your bot name otherwise.",
				Optional:            true,
				Computed:            true,
			},
			"special_character_required": schema.BoolAttribute{
				MarkdownDescription: "Requires user passwords to contain at least one special character (e.g. !@#$%^&*).",
				Optional:            true,
				Computed:            true,
			},
			"team_name": schema.StringAttribute{
				MarkdownDescription: "The name of the team/organization, displayed in the DefectDojo UI.",
				Optional:            true,
				Computed:            true,
			},
			"uppercase_character_required": schema.BoolAttribute{
				MarkdownDescription: "Requires user passwords to contain at least one uppercase letter (A-Z).",
				Optional:            true,
				Computed:            true,
			},
			"url_prefix": schema.StringAttribute{
				MarkdownDescription: "URL prefix if DefectDojo is installed in it's own virtual subdirectory.",
				Optional:            true,
				Computed:            true,
			},
			"webhooks_notifications_timeout": schema.Int64Attribute{
				MarkdownDescription: "How many seconds will DefectDojo waits for response from webhook endpoint.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

type systemSettingsResourceData struct {
	Id types.String `tfsdk:"id" ddField:"Id"`

	AddVulnerabilityIdToJiraLabel       types.Bool  `tfsdk:"add_vulnerability_id_to_jira_label" ddField:"AddVulnerabilityIdToJiraLabel"`
	AllowAnonymousSurveyRepsonse        types.Bool  `tfsdk:"allow_anonymous_survey_repsonse" ddField:"AllowAnonymousSurveyRepsonse"`
	ApiExposeErrorDetails               types.Bool  `tfsdk:"api_expose_error_details" ddField:"ApiExposeErrorDetails"`
	DeleteDuplicates                    types.Bool  `tfsdk:"delete_duplicates" ddField:"DeleteDuplicates"`
	DisableJiraWebhookSecret            types.Bool  `tfsdk:"disable_jira_webhook_secret" ddField:"DisableJiraWebhookSecret"`
	DisclaimerReportsForced             types.Bool  `tfsdk:"disclaimer_reports_forced" ddField:"DisclaimerReportsForced"`
	EnableBenchmark                     types.Bool  `tfsdk:"enable_benchmark" ddField:"EnableBenchmark"`
	EnableCalendar                      types.Bool  `tfsdk:"enable_calendar" ddField:"EnableCalendar"`
	EnableChecklists                    types.Bool  `tfsdk:"enable_checklists" ddField:"EnableChecklists"`
	EnableCvss3Display                  types.Bool  `tfsdk:"enable_cvss3_display" ddField:"EnableCvss3Display"`
	EnableCvss4Display                  types.Bool  `tfsdk:"enable_cvss4_display" ddField:"EnableCvss4Display"`
	EnableDeduplication                 types.Bool  `tfsdk:"enable_deduplication" ddField:"EnableDeduplication"`
	EnableEndpointMetadataImport        types.Bool  `tfsdk:"enable_endpoint_metadata_import" ddField:"EnableEndpointMetadataImport"`
	EnableFindingGroups                 types.Bool  `tfsdk:"enable_finding_groups" ddField:"EnableFindingGroups"`
	EnableFindingSla                    types.Bool  `tfsdk:"enable_finding_sla" ddField:"EnableFindingSla"`
	EnableGithub                        types.Bool  `tfsdk:"enable_github" ddField:"EnableGithub"`
	EnableJira                          types.Bool  `tfsdk:"enable_jira" ddField:"EnableJira"`
	EnableJiraWebHook                   types.Bool  `tfsdk:"enable_jira_web_hook" ddField:"EnableJiraWebHook"`
	EnableMailNotifications             types.Bool  `tfsdk:"enable_mail_notifications" ddField:"EnableMailNotifications"`
	EnableMsteamsNotifications          types.Bool  `tfsdk:"enable_msteams_notifications" ddField:"EnableMsteamsNotifications"`
	EnableNotifySlaActive               types.Bool  `tfsdk:"enable_notify_sla_active" ddField:"EnableNotifySlaActive"`
	EnableNotifySlaActiveVerified       types.Bool  `tfsdk:"enable_notify_sla_active_verified" ddField:"EnableNotifySlaActiveVerified"`
	EnableNotifySlaExponentialBackoff   types.Bool  `tfsdk:"enable_notify_sla_exponential_backoff" ddField:"EnableNotifySlaExponentialBackoff"`
	EnableNotifySlaJiraOnly             types.Bool  `tfsdk:"enable_notify_sla_jira_only" ddField:"EnableNotifySlaJiraOnly"`
	EnableProductGrade                  types.Bool  `tfsdk:"enable_product_grade" ddField:"EnableProductGrade"`
	EnableProductTagInheritance         types.Bool  `tfsdk:"enable_product_tag_inheritance" ddField:"EnableProductTagInheritance"`
	EnableProductTrackingFiles          types.Bool  `tfsdk:"enable_product_tracking_files" ddField:"EnableProductTrackingFiles"`
	EnableQuestionnaires                types.Bool  `tfsdk:"enable_questionnaires" ddField:"EnableQuestionnaires"`
	EnableSimilarFindings               types.Bool  `tfsdk:"enable_similar_findings" ddField:"EnableSimilarFindings"`
	EnableSlackNotifications            types.Bool  `tfsdk:"enable_slack_notifications" ddField:"EnableSlackNotifications"`
	EnableUiTableBasedSearching         types.Bool  `tfsdk:"enable_ui_table_based_searching" ddField:"EnableUiTableBasedSearching"`
	EnableUserProfileEditable           types.Bool  `tfsdk:"enable_user_profile_editable" ddField:"EnableUserProfileEditable"`
	EnableWebhooksNotifications         types.Bool  `tfsdk:"enable_webhooks_notifications" ddField:"EnableWebhooksNotifications"`
	EnforceVerifiedStatus               types.Bool  `tfsdk:"enforce_verified_status" ddField:"EnforceVerifiedStatus"`
	EnforceVerifiedStatusJira           types.Bool  `tfsdk:"enforce_verified_status_jira" ddField:"EnforceVerifiedStatusJira"`
	EnforceVerifiedStatusMetrics        types.Bool  `tfsdk:"enforce_verified_status_metrics" ddField:"EnforceVerifiedStatusMetrics"`
	EnforceVerifiedStatusProductGrading types.Bool  `tfsdk:"enforce_verified_status_product_grading" ddField:"EnforceVerifiedStatusProductGrading"`
	EngagementAutoClose                 types.Bool  `tfsdk:"engagement_auto_close" ddField:"EngagementAutoClose"`
	EngagementAutoCloseDays             types.Int64 `tfsdk:"engagement_auto_close_days" ddField:"EngagementAutoCloseDays"`
	FalsePositiveHistory                types.Bool  `tfsdk:"false_positive_history" ddField:"FalsePositiveHistory"`
	FilterStringMatching                types.Bool  `tfsdk:"filter_string_matching" ddField:"FilterStringMatching"`

	DisclaimerNotes         types.String `tfsdk:"disclaimer_notes" ddField:"DisclaimerNotes"`
	DisclaimerNotifications types.String `tfsdk:"disclaimer_notifications" ddField:"DisclaimerNotifications"`
	DisclaimerReports       types.String `tfsdk:"disclaimer_reports" ddField:"DisclaimerReports"`
	EmailFrom               types.String `tfsdk:"email_from" ddField:"EmailFrom"`
	JiraLabels              types.String `tfsdk:"jira_labels" ddField:"JiraLabels"`
	JiraMinimumSeverity     types.String `tfsdk:"jira_minimum_severity" ddField:"JiraMinimumSeverity"`
	JiraWebhookSecret       types.String `tfsdk:"jira_webhook_secret" ddField:"JiraWebhookSecret"`
	MailNotificationsTo     types.String `tfsdk:"mail_notifications_to" ddField:"MailNotificationsTo"`
	MsteamsUrl              types.String `tfsdk:"msteams_url" ddField:"MsteamsUrl"`
	SlackChannel            types.String `tfsdk:"slack_channel" ddField:"SlackChannel"`
	SlackToken              types.String `tfsdk:"slack_token" ddField:"SlackToken"`
	SlackUsername           types.String `tfsdk:"slack_username" ddField:"SlackUsername"`
	TeamName                types.String `tfsdk:"team_name" ddField:"TeamName"`
	UrlPrefix               types.String `tfsdk:"url_prefix" ddField:"UrlPrefix"`

	LowercaseCharacterRequired      types.Bool `tfsdk:"lowercase_character_required" ddField:"LowercaseCharacterRequired"`
	NonCommonPasswordRequired       types.Bool `tfsdk:"non_common_password_required" ddField:"NonCommonPasswordRequired"`
	NumberCharacterRequired         types.Bool `tfsdk:"number_character_required" ddField:"NumberCharacterRequired"`
	RetroactiveFalsePositiveHistory types.Bool `tfsdk:"retroactive_false_positive_history" ddField:"RetroactiveFalsePositiveHistory"`
	SpecialCharacterRequired        types.Bool `tfsdk:"special_character_required" ddField:"SpecialCharacterRequired"`
	UppercaseCharacterRequired      types.Bool `tfsdk:"uppercase_character_required" ddField:"UppercaseCharacterRequired"`

	MaxDupes                             types.Int64 `tfsdk:"max_dupes" ddField:"MaxDupes"`
	MaximumPasswordLength                types.Int64 `tfsdk:"maximum_password_length" ddField:"MaximumPasswordLength"`
	MinimumPasswordLength                types.Int64 `tfsdk:"minimum_password_length" ddField:"MinimumPasswordLength"`
	ProductGradeA                        types.Int64 `tfsdk:"product_grade_a" ddField:"ProductGradeA"`
	ProductGradeB                        types.Int64 `tfsdk:"product_grade_b" ddField:"ProductGradeB"`
	ProductGradeC                        types.Int64 `tfsdk:"product_grade_c" ddField:"ProductGradeC"`
	ProductGradeD                        types.Int64 `tfsdk:"product_grade_d" ddField:"ProductGradeD"`
	ProductGradeF                        types.Int64 `tfsdk:"product_grade_f" ddField:"ProductGradeF"`
	RiskAcceptanceFormDefaultDays        types.Int64 `tfsdk:"risk_acceptance_form_default_days" ddField:"RiskAcceptanceFormDefaultDays"`
	RiskAcceptanceNotifyBeforeExpiration types.Int64 `tfsdk:"risk_acceptance_notify_before_expiration" ddField:"RiskAcceptanceNotifyBeforeExpiration"`
	WebhooksNotificationsTimeout         types.Int64 `tfsdk:"webhooks_notifications_timeout" ddField:"WebhooksNotificationsTimeout"`
}

type systemSettingsDefectdojoResource struct {
	dd.SystemSettings
}

var _ singletonAdopter = &systemSettingsDefectdojoResource{}

// systemSettingsToRequest converts a SystemSettings (response model) into the
// request body expected by SystemSettingsPartialUpdateWithResponse
// (PatchedSystemSettingsRequest). Only the pointer fields actually populated
// by populateDefectdojoResource (i.e. attributes present in configuration)
// carry a non-nil value here, and thanks to `omitempty` those are the only
// fields serialized into the PATCH body.
func systemSettingsToRequest(s dd.SystemSettings) dd.PatchedSystemSettingsRequest {
	req := dd.PatchedSystemSettingsRequest{
		AddVulnerabilityIdToJiraLabel:        s.AddVulnerabilityIdToJiraLabel,
		AllowAnonymousSurveyRepsonse:         s.AllowAnonymousSurveyRepsonse,
		ApiExposeErrorDetails:                s.ApiExposeErrorDetails,
		DeleteDuplicates:                     s.DeleteDuplicates,
		DisableJiraWebhookSecret:             s.DisableJiraWebhookSecret,
		DisclaimerNotes:                      s.DisclaimerNotes,
		DisclaimerNotifications:              s.DisclaimerNotifications,
		DisclaimerReports:                    s.DisclaimerReports,
		DisclaimerReportsForced:              s.DisclaimerReportsForced,
		EmailFrom:                            s.EmailFrom,
		EnableBenchmark:                      s.EnableBenchmark,
		EnableCalendar:                       s.EnableCalendar,
		EnableChecklists:                     s.EnableChecklists,
		EnableCvss3Display:                   s.EnableCvss3Display,
		EnableCvss4Display:                   s.EnableCvss4Display,
		EnableDeduplication:                  s.EnableDeduplication,
		EnableEndpointMetadataImport:         s.EnableEndpointMetadataImport,
		EnableFindingGroups:                  s.EnableFindingGroups,
		EnableFindingSla:                     s.EnableFindingSla,
		EnableGithub:                         s.EnableGithub,
		EnableJira:                           s.EnableJira,
		EnableJiraWebHook:                    s.EnableJiraWebHook,
		EnableMailNotifications:              s.EnableMailNotifications,
		EnableMsteamsNotifications:           s.EnableMsteamsNotifications,
		EnableNotifySlaActive:                s.EnableNotifySlaActive,
		EnableNotifySlaActiveVerified:        s.EnableNotifySlaActiveVerified,
		EnableNotifySlaExponentialBackoff:    s.EnableNotifySlaExponentialBackoff,
		EnableNotifySlaJiraOnly:              s.EnableNotifySlaJiraOnly,
		EnableProductGrade:                   s.EnableProductGrade,
		EnableProductTagInheritance:          s.EnableProductTagInheritance,
		EnableProductTrackingFiles:           s.EnableProductTrackingFiles,
		EnableQuestionnaires:                 s.EnableQuestionnaires,
		EnableSimilarFindings:                s.EnableSimilarFindings,
		EnableSlackNotifications:             s.EnableSlackNotifications,
		EnableUiTableBasedSearching:          s.EnableUiTableBasedSearching,
		EnableUserProfileEditable:            s.EnableUserProfileEditable,
		EnableWebhooksNotifications:          s.EnableWebhooksNotifications,
		EnforceVerifiedStatus:                s.EnforceVerifiedStatus,
		EnforceVerifiedStatusJira:            s.EnforceVerifiedStatusJira,
		EnforceVerifiedStatusMetrics:         s.EnforceVerifiedStatusMetrics,
		EnforceVerifiedStatusProductGrading:  s.EnforceVerifiedStatusProductGrading,
		EngagementAutoClose:                  s.EngagementAutoClose,
		EngagementAutoCloseDays:              s.EngagementAutoCloseDays,
		FalsePositiveHistory:                 s.FalsePositiveHistory,
		FilterStringMatching:                 s.FilterStringMatching,
		JiraLabels:                           s.JiraLabels,
		JiraWebhookSecret:                    s.JiraWebhookSecret,
		LowercaseCharacterRequired:           s.LowercaseCharacterRequired,
		MailNotificationsTo:                  s.MailNotificationsTo,
		MaxDupes:                             s.MaxDupes,
		MaximumPasswordLength:                s.MaximumPasswordLength,
		MinimumPasswordLength:                s.MinimumPasswordLength,
		MsteamsUrl:                           s.MsteamsUrl,
		NonCommonPasswordRequired:            s.NonCommonPasswordRequired,
		NumberCharacterRequired:              s.NumberCharacterRequired,
		ProductGradeA:                        s.ProductGradeA,
		ProductGradeB:                        s.ProductGradeB,
		ProductGradeC:                        s.ProductGradeC,
		ProductGradeD:                        s.ProductGradeD,
		ProductGradeF:                        s.ProductGradeF,
		RetroactiveFalsePositiveHistory:      s.RetroactiveFalsePositiveHistory,
		RiskAcceptanceFormDefaultDays:        s.RiskAcceptanceFormDefaultDays,
		RiskAcceptanceNotifyBeforeExpiration: s.RiskAcceptanceNotifyBeforeExpiration,
		SlackChannel:                         s.SlackChannel,
		SlackToken:                           s.SlackToken,
		SlackUsername:                        s.SlackUsername,
		SpecialCharacterRequired:             s.SpecialCharacterRequired,
		TeamName:                             s.TeamName,
		UppercaseCharacterRequired:           s.UppercaseCharacterRequired,
		UrlPrefix:                            s.UrlPrefix,
		WebhooksNotificationsTimeout:         s.WebhooksNotificationsTimeout,
	}
	// JiraMinimumSeverity is an enum: SystemSettings uses
	// SystemSettingsJiraMinimumSeverity, PatchedSystemSettingsRequest uses
	// PatchedSystemSettingsRequestJiraMinimumSeverity. Both are named string
	// types with identical values, so a direct cast suffices.
	if s.JiraMinimumSeverity != nil {
		v := dd.PatchedSystemSettingsRequestJiraMinimumSeverity(*s.JiraMinimumSeverity)
		req.JiraMinimumSeverity = &v
	}
	return req
}

// createApiCall/deleteApiCall exist only to satisfy the defectdojoResource
// interface. system_settings is a singletonAdopter: the engine calls
// adoptApiCall + updateApiCall instead of createApiCall on Create, and only
// removes the resource from state (without calling deleteApiCall) on Delete.
func (ddr *systemSettingsDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	return 0, nil, errSystemSettingsSingleton
}

func (ddr *systemSettingsDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	return 0, nil, errSystemSettingsSingleton
}

// adoptApiCall resolves the id of the single system_settings row via the
// list endpoint (there is no dedicated "retrieve" or "create" endpoint).
// It deliberately does NOT copy the row into ddr.SystemSettings: the engine
// has already staged the Terraform plan's values there via
// populateDefectdojoResource, and immediately follows adoptApiCall with
// updateApiCall using that same struct - overwriting it here would discard
// the plan.
func (ddr *systemSettingsDefectdojoResource) adoptApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, int, []byte, error) {
	tflog.Info(ctx, "adoptApiCall")
	apiResp, err := client.SystemSettingsListWithResponse(ctx, &dd.SystemSettingsListParams{})
	if err != nil {
		return 0, 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))

	if apiResp.JSON200 == nil || len(apiResp.JSON200.Results) == 0 || apiResp.JSON200.Results[0].Id == nil {
		// Unexpected shape: the request succeeded but did not return the
		// single system_settings row we require. Force a non-200 status
		// (even if the real status was 200) so the engine's diagnostic
		// fires and surfaces the raw body for debugging.
		statusCode := apiResp.StatusCode()
		if statusCode == 200 {
			statusCode = 0
		}
		return 0, statusCode, apiResp.Body, nil
	}

	return *apiResp.JSON200.Results[0].Id, apiResp.StatusCode(), apiResp.Body, nil
}

// readApiCall has no dedicated retrieve endpoint to call: system_settings
// only exposes List and Update/PartialUpdate. It lists and finds the row
// matching id, falling back to the sole row when there is exactly one
// (covering the case where the row's id wasn't known ahead of time).
func (ddr *systemSettingsDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "readApiCall")
	apiResp, err := client.SystemSettingsListWithResponse(ctx, &dd.SystemSettingsListParams{})
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))

	if apiResp.StatusCode() != 200 || apiResp.JSON200 == nil {
		return apiResp.StatusCode(), apiResp.Body, nil
	}

	var found *dd.SystemSettings
	for i := range apiResp.JSON200.Results {
		row := apiResp.JSON200.Results[i]
		if row.Id != nil && *row.Id == idNumber {
			found = &row
			break
		}
	}
	if found == nil && len(apiResp.JSON200.Results) == 1 {
		found = &apiResp.JSON200.Results[0]
	}
	if found == nil {
		return 404, apiResp.Body, nil
	}

	ddr.SystemSettings = *found
	return 200, apiResp.Body, nil
}

// updateApiCall uses PATCH (SystemSettingsPartialUpdateWithResponse) rather
// than PUT (SystemSettingsUpdateWithResponse). THIS IS A DELIBERATE
// DEVIATION from every other resource in this provider: system_settings is a
// singleton with ~70 fields, and Terraform only manages the subset actually
// present in configuration - the engine skips null/unknown fields when
// staging the plan into ddr.SystemSettings via populateDefectdojoResource,
// so most fields on ddr.SystemSettings are nil for any given apply. A full
// PUT would serialize every omitted field as its Go zero value and the
// server would persist it verbatim, silently resetting every unconfigured
// setting (JIRA config, SLA notification toggles, password policy, etc.) to
// its zero/default value and wiping out whatever an administrator had
// already configured directly in DefectDojo. PATCH only sends the fields we
// actually populated (thanks to `omitempty` on every pointer field), leaving
// every other server-side setting untouched.
func (ddr *systemSettingsDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "updateApiCall")
	reqBody := systemSettingsToRequest(ddr.SystemSettings)
	apiResp, err := client.SystemSettingsPartialUpdateWithResponse(ctx, idNumber, reqBody)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON200 != nil {
		ddr.SystemSettings = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

type systemSettingsResource struct {
	terraformResource
}

var _ resource.Resource = &systemSettingsResource{}
var _ resource.ResourceWithImportState = &systemSettingsResource{}

func NewSystemSettingsResource() resource.Resource {
	return &systemSettingsResource{
		terraformResource: terraformResource{
			typeName:     "defectdojo_system_settings",
			dataProvider: systemSettingsDataProvider{},
		},
	}
}

func (r systemSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_settings"
}

type systemSettingsDataProvider struct{}

func (r systemSettingsDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data systemSettingsResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *systemSettingsResourceData) id() types.String {
	return d.Id
}

func (d *systemSettingsResourceData) setId(v types.String) { d.Id = v }

func (d *systemSettingsResourceData) defectdojoResource() defectdojoResource {
	return &systemSettingsDefectdojoResource{
		SystemSettings: dd.SystemSettings{},
	}
}
