package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
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
			// The live API defaults any channel field left unset in a
			// create/update request to ["alert"] rather than leaving it
			// empty, so these must be Computed (with UseStateForUnknown) to
			// avoid a "provider produced inconsistent result" error whenever
			// a config leaves one of these attributes unset.
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
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
				// The live API defaults this to false when unset in a
				// create/update request, so it must be Computed (with
				// UseStateForUnknown) to avoid a "provider produced
				// inconsistent result" error when left unset in config.
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"scan_added_empty": schema.StringAttribute{
				MarkdownDescription: "Triggered whenever an (re-)import has been done (even if that created/updated/closed no findings). Valid values: 'alert', 'mail', 'msteams', 'slack', 'webhooks'.",
				// The live API defaults this to "" when unset in a
				// create/update request, so it must be Computed (with
				// UseStateForUnknown) to avoid a "provider produced
				// inconsistent result" error when left unset in config.
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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

// notificationsScanAddedEmpty is the wire representation of scan_added_empty.
// The live 3.1.101 API is asymmetric for this one field: GET/LIST (and
// create/update responses when the field was left unset in the request)
// always return it as a JSON array of 0 or 1 strings (e.g. `[]` or
// `["alert"]`), while the write path (POST/PUT/PATCH) only accepts - and
// rejects an array for - a bare scalar string (e.g. `"alert"`; even `[]`
// or `["alert"]` fail with "... is not a valid choice"). The type's default
// MarshalJSON (inherited from the underlying string type) already produces
// the bare-scalar shape the write path requires; only UnmarshalJSON needs to
// be customized to also accept the array shape returned on reads.
type notificationsScanAddedEmpty string

func (v *notificationsScanAddedEmpty) UnmarshalJSON(data []byte) error {
	var asArray []string
	if err := json.Unmarshal(data, &asArray); err == nil {
		if len(asArray) > 0 {
			*v = notificationsScanAddedEmpty(asArray[0])
		} else {
			*v = ""
		}
		return nil
	}

	var asString string
	if err := json.Unmarshal(data, &asString); err != nil {
		return fmt.Errorf("scan_added_empty: expected a string or an array of strings: %w", err)
	}
	*v = notificationsScanAddedEmpty(asString)
	return nil
}

// notificationsModel mirrors the 3.1.101 /api/v2/notifications/ JSON. The
// generated dd.Notifications/dd.NotificationsRequest cannot be used here: the
// spec types scan_added_empty as a scalar enum, but the live API returns it as
// an array on every read path (see notificationsScanAddedEmpty above), so the
// generated response parsers fail to unmarshal every create/read/update.
type notificationsModel struct {
	Id       *int  `json:"id,omitempty"`
	Product  *int  `json:"product,omitempty"`
	User     *int  `json:"user,omitempty"`
	Template *bool `json:"template,omitempty"`

	ScanAddedEmpty *notificationsScanAddedEmpty `json:"scan_added_empty,omitempty"`

	AutoCloseEngagement      *[]string `json:"auto_close_engagement,omitempty"`
	CloseEngagement          *[]string `json:"close_engagement,omitempty"`
	CodeReview               *[]string `json:"code_review,omitempty"`
	EngagementAdded          *[]string `json:"engagement_added,omitempty"`
	JiraUpdate               *[]string `json:"jira_update,omitempty"`
	Other                    *[]string `json:"other,omitempty"`
	ProductAdded             *[]string `json:"product_added,omitempty"`
	ProductTypeAdded         *[]string `json:"product_type_added,omitempty"`
	ReviewRequested          *[]string `json:"review_requested,omitempty"`
	RiskAcceptanceExpiration *[]string `json:"risk_acceptance_expiration,omitempty"`
	ScanAdded                *[]string `json:"scan_added,omitempty"`
	SlaBreach                *[]string `json:"sla_breach,omitempty"`
	SlaBreachCombined        *[]string `json:"sla_breach_combined,omitempty"`
	StaleEngagement          *[]string `json:"stale_engagement,omitempty"`
	TestAdded                *[]string `json:"test_added,omitempty"`
	UpcomingEngagement       *[]string `json:"upcoming_engagement,omitempty"`
	UserMentioned            *[]string `json:"user_mentioned,omitempty"`
}

type notificationsDefectdojoResource struct {
	notificationsModel
}

func (ddr *notificationsDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	tflog.Info(ctx, "createApiCall")
	reqBody, err := json.Marshal(ddr.notificationsModel)
	if err != nil {
		return 0, nil, err
	}
	httpResp, err := client.NotificationsCreateWithBody(ctx, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return 0, nil, err
	}
	defer func() { _ = httpResp.Body.Close() }()
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %d: %s", httpResp.StatusCode, body))
	if httpResp.StatusCode == 201 {
		if err := json.Unmarshal(body, &ddr.notificationsModel); err != nil {
			return httpResp.StatusCode, body, err
		}
	}

	return httpResp.StatusCode, body, nil
}

func (ddr *notificationsDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "readApiCall")
	httpResp, err := client.NotificationsRetrieve(ctx, idNumber, &dd.NotificationsRetrieveParams{})
	if err != nil {
		return 0, nil, err
	}
	defer func() { _ = httpResp.Body.Close() }()
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %d: %s", httpResp.StatusCode, body))
	if httpResp.StatusCode == 200 {
		if err := json.Unmarshal(body, &ddr.notificationsModel); err != nil {
			return httpResp.StatusCode, body, err
		}
	}

	return httpResp.StatusCode, body, nil
}

func (ddr *notificationsDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "updateApiCall")
	reqBody, err := json.Marshal(ddr.notificationsModel)
	if err != nil {
		return 0, nil, err
	}
	httpResp, err := client.NotificationsUpdateWithBody(ctx, idNumber, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return 0, nil, err
	}
	defer func() { _ = httpResp.Body.Close() }()
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %d: %s", httpResp.StatusCode, body))
	if httpResp.StatusCode == 200 {
		if err := json.Unmarshal(body, &ddr.notificationsModel); err != nil {
			return httpResp.StatusCode, body, err
		}
	}
	return httpResp.StatusCode, body, nil
}

func (ddr *notificationsDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "deleteApiCall")
	httpResp, err := client.NotificationsDestroy(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	defer func() { _ = httpResp.Body.Close() }()
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %d: %s", httpResp.StatusCode, body))
	return httpResp.StatusCode, body, nil
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
		notificationsModel: notificationsModel{},
	}
}
