package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
)

// notificationChannels is the set of valid values for every notification
// channel attribute (both the Set(String) attributes and the single-value
// scan_added_empty attribute).
var notificationChannels = []string{"alert", "mail", "msteams", "slack", "webhooks"}

func (t notificationsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	channelSetAttribute := func(description string) schema.SetAttribute {
		return schema.SetAttribute{
			MarkdownDescription: description,
			Optional:            true,
			ElementType:         types.StringType,
			Validators: []validator.Set{
				setvalidator.ValueStringsAre(
					stringvalidator.OneOf(notificationChannels...),
				),
			},
		}
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Notifications: per-channel notification preferences for a scope (global default, per-user, or per-product; at most one of `product` or `user` may be set). " +
			"A fresh DefectDojo instance pre-creates the global default row and a per-user row for each existing user (e.g. admin); creating a duplicate row for a scope that already exists returns an API error " +
			"(\"Notification for user and product already exists\"). To manage a pre-existing row (the global default, or a user's default row), import it instead of trying to create it, e.g. " +
			"`terraform import defectdojo_notifications.global <id>`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"product": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Product this notification preference row is scoped to. Mutually exclusive with `user`.",
				Optional:            true,
			},
			"user": schema.Int64Attribute{
				MarkdownDescription: "The ID of the User this notification preference row is scoped to. Mutually exclusive with `product`.",
				Optional:            true,
			},
			"template": schema.BoolAttribute{
				MarkdownDescription: "Whether this row is the notification template applied to new users/products.",
				Optional:            true,
			},
			"scan_added_empty": schema.StringAttribute{
				MarkdownDescription: "Triggered whenever an (re-)import has been done (even if that created/updated/closed no findings). Valid values: 'alert', 'mail', 'msteams', 'slack', 'webhooks'.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(notificationChannels...),
				},
			},
			"auto_close_engagement":      channelSetAttribute("Notification channels for when an engagement is auto-closed."),
			"close_engagement":           channelSetAttribute("Notification channels for when an engagement is closed."),
			"code_review":                channelSetAttribute("Notification channels for when a code review is requested."),
			"engagement_added":           channelSetAttribute("Notification channels for when an engagement is added."),
			"jira_update":                channelSetAttribute("Notification channels for when a JIRA issue is updated."),
			"other":                      channelSetAttribute("Notification channels for other events."),
			"product_added":              channelSetAttribute("Notification channels for when a product is added."),
			"product_type_added":         channelSetAttribute("Notification channels for when a product type is added."),
			"review_requested":           channelSetAttribute("Notification channels for when a review is requested."),
			"risk_acceptance_expiration": channelSetAttribute("Notification channels for when a risk acceptance is about to expire."),
			"scan_added":                 channelSetAttribute("Notification channels for when a scan is added/imported."),
			"sla_breach":                 channelSetAttribute("Notification channels for when an SLA is breached."),
			"sla_breach_combined":        channelSetAttribute("Notification channels for combined SLA breach notifications."),
			"stale_engagement":           channelSetAttribute("Notification channels for stale engagements."),
			"test_added":                 channelSetAttribute("Notification channels for when a test is added."),
			"upcoming_engagement":        channelSetAttribute("Notification channels for upcoming engagements."),
			"user_mentioned":             channelSetAttribute("Notification channels for when a user is mentioned."),
		},
	}
}

type notificationsResourceData struct {
	Id       types.String `tfsdk:"id" ddField:"Id"`
	Product  types.Int64  `tfsdk:"product" ddField:"Product"`
	User     types.Int64  `tfsdk:"user" ddField:"User"`
	Template types.Bool   `tfsdk:"template" ddField:"Template"`

	ScanAddedEmpty types.String `tfsdk:"scan_added_empty" ddField:"ScanAddedEmpty"`

	AutoCloseEngagement      types.Set `tfsdk:"auto_close_engagement" ddField:"AutoCloseEngagement"`
	CloseEngagement          types.Set `tfsdk:"close_engagement" ddField:"CloseEngagement"`
	CodeReview               types.Set `tfsdk:"code_review" ddField:"CodeReview"`
	EngagementAdded          types.Set `tfsdk:"engagement_added" ddField:"EngagementAdded"`
	JiraUpdate               types.Set `tfsdk:"jira_update" ddField:"JiraUpdate"`
	Other                    types.Set `tfsdk:"other" ddField:"Other"`
	ProductAdded             types.Set `tfsdk:"product_added" ddField:"ProductAdded"`
	ProductTypeAdded         types.Set `tfsdk:"product_type_added" ddField:"ProductTypeAdded"`
	ReviewRequested          types.Set `tfsdk:"review_requested" ddField:"ReviewRequested"`
	RiskAcceptanceExpiration types.Set `tfsdk:"risk_acceptance_expiration" ddField:"RiskAcceptanceExpiration"`
	ScanAdded                types.Set `tfsdk:"scan_added" ddField:"ScanAdded"`
	SlaBreach                types.Set `tfsdk:"sla_breach" ddField:"SlaBreach"`
	SlaBreachCombined        types.Set `tfsdk:"sla_breach_combined" ddField:"SlaBreachCombined"`
	StaleEngagement          types.Set `tfsdk:"stale_engagement" ddField:"StaleEngagement"`
	TestAdded                types.Set `tfsdk:"test_added" ddField:"TestAdded"`
	UpcomingEngagement       types.Set `tfsdk:"upcoming_engagement" ddField:"UpcomingEngagement"`
	UserMentioned            types.Set `tfsdk:"user_mentioned" ddField:"UserMentioned"`

	// Prefetch (a nested struct of *map[string]Product / *map[string]UserStub)
	// is intentionally not mapped here: the reflection-based CRUD engine
	// (populateDefectdojoResource/populateResourceData in resource.go) has no
	// case for mapping an arbitrary nested struct to/from a Terraform type,
	// and there is no corresponding schema attribute for it.
}

type notificationsDefectdojoResource struct {
	dd.Notifications
}

// notificationsConvertEnumSlice converts a *[]S (e.g. *[]dd.NotificationsScanAdded)
// to a *[]D (e.g. *[]dd.NotificationsRequestScanAdded). S and D are both defined
// string types that share the same underlying values, but are distinct Go types
// generated for the Notifications (response) and NotificationsRequest (request)
// schemas respectively, so a direct assignment does not compile.
func notificationsConvertEnumSlice[S, D ~string](in *[]S) *[]D {
	if in == nil {
		return nil
	}
	out := make([]D, len(*in))
	for i, v := range *in {
		out[i] = D(v)
	}
	return &out
}

// notificationsToRequest converts a Notifications (response model) to a
// NotificationsRequest (request model).
func notificationsToRequest(n dd.Notifications) dd.NotificationsRequest {
	req := dd.NotificationsRequest{
		Product:  n.Product,
		User:     n.User,
		Template: n.Template,

		AutoCloseEngagement:      notificationsConvertEnumSlice[dd.NotificationsAutoCloseEngagement, dd.NotificationsRequestAutoCloseEngagement](n.AutoCloseEngagement),
		CloseEngagement:          notificationsConvertEnumSlice[dd.NotificationsCloseEngagement, dd.NotificationsRequestCloseEngagement](n.CloseEngagement),
		CodeReview:               notificationsConvertEnumSlice[dd.NotificationsCodeReview, dd.NotificationsRequestCodeReview](n.CodeReview),
		EngagementAdded:          notificationsConvertEnumSlice[dd.NotificationsEngagementAdded, dd.NotificationsRequestEngagementAdded](n.EngagementAdded),
		JiraUpdate:               notificationsConvertEnumSlice[dd.NotificationsJiraUpdate, dd.NotificationsRequestJiraUpdate](n.JiraUpdate),
		Other:                    notificationsConvertEnumSlice[dd.NotificationsOther, dd.NotificationsRequestOther](n.Other),
		ProductAdded:             notificationsConvertEnumSlice[dd.NotificationsProductAdded, dd.NotificationsRequestProductAdded](n.ProductAdded),
		ProductTypeAdded:         notificationsConvertEnumSlice[dd.NotificationsProductTypeAdded, dd.NotificationsRequestProductTypeAdded](n.ProductTypeAdded),
		ReviewRequested:          notificationsConvertEnumSlice[dd.NotificationsReviewRequested, dd.NotificationsRequestReviewRequested](n.ReviewRequested),
		RiskAcceptanceExpiration: notificationsConvertEnumSlice[dd.NotificationsRiskAcceptanceExpiration, dd.NotificationsRequestRiskAcceptanceExpiration](n.RiskAcceptanceExpiration),
		ScanAdded:                notificationsConvertEnumSlice[dd.NotificationsScanAdded, dd.NotificationsRequestScanAdded](n.ScanAdded),
		SlaBreach:                notificationsConvertEnumSlice[dd.NotificationsSlaBreach, dd.NotificationsRequestSlaBreach](n.SlaBreach),
		SlaBreachCombined:        notificationsConvertEnumSlice[dd.NotificationsSlaBreachCombined, dd.NotificationsRequestSlaBreachCombined](n.SlaBreachCombined),
		StaleEngagement:          notificationsConvertEnumSlice[dd.NotificationsStaleEngagement, dd.NotificationsRequestStaleEngagement](n.StaleEngagement),
		TestAdded:                notificationsConvertEnumSlice[dd.NotificationsTestAdded, dd.NotificationsRequestTestAdded](n.TestAdded),
		UpcomingEngagement:       notificationsConvertEnumSlice[dd.NotificationsUpcomingEngagement, dd.NotificationsRequestUpcomingEngagement](n.UpcomingEngagement),
		UserMentioned:            notificationsConvertEnumSlice[dd.NotificationsUserMentioned, dd.NotificationsRequestUserMentioned](n.UserMentioned),
	}
	if n.ScanAddedEmpty != nil {
		v := dd.NotificationsRequestScanAddedEmpty(*n.ScanAddedEmpty)
		req.ScanAddedEmpty = &v
	}
	return req
}

func (ddr *notificationsDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	tflog.Info(ctx, "createApiCall")
	reqBody := notificationsToRequest(ddr.Notifications)
	apiResp, err := client.NotificationsCreateWithResponse(ctx, reqBody)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON201 != nil {
		ddr.Notifications = *apiResp.JSON201
	}

	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *notificationsDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "readApiCall")
	apiResp, err := client.NotificationsRetrieveWithResponse(ctx, idNumber, &dd.NotificationsRetrieveParams{})
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON200 != nil {
		ddr.Notifications = *apiResp.JSON200
	}

	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *notificationsDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "updateApiCall")
	reqBody := notificationsToRequest(ddr.Notifications)
	apiResp, err := client.NotificationsUpdateWithResponse(ctx, idNumber, reqBody)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON200 != nil {
		ddr.Notifications = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *notificationsDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "deleteApiCall")
	apiResp, err := client.NotificationsDestroyWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	return apiResp.StatusCode(), apiResp.Body, nil
}

type notificationsResource struct {
	terraformResource
}

var _ resource.Resource = &notificationsResource{}
var _ resource.ResourceWithImportState = &notificationsResource{}
var _ resource.ResourceWithConfigValidators = &notificationsResource{}

func NewNotificationsResource() resource.Resource {
	return &notificationsResource{
		terraformResource: terraformResource{
			typeName:     "defectdojo_notifications",
			dataProvider: notificationsDataProvider{},
		},
	}
}

func (r notificationsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notifications"
}

// ConfigValidators enforces that a notifications row is scoped to at most one
// of product or user (a row with neither set is the global default row).
func (r notificationsResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("product"), path.MatchRoot("user"),
		),
	}
}

type notificationsDataProvider struct{}

func (r notificationsDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data notificationsResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *notificationsResourceData) id() types.String {
	return d.Id
}

func (d *notificationsResourceData) setId(v types.String) { d.Id = v }

func (d *notificationsResourceData) defectdojoResource() defectdojoResource {
	return &notificationsDefectdojoResource{
		Notifications: dd.Notifications{},
	}
}
