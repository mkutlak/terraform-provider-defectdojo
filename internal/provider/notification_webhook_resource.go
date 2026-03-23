package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
)

func (t notificationWebhookResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Notification Webhook",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the incoming webhook.",
				Optional:            true,
				Computed:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "The full URL of the incoming webhook.",
				Optional:            true,
				Computed:            true,
			},
			"header_name": schema.StringAttribute{
				MarkdownDescription: "Name of the header required for interacting with the webhook endpoint.",
				Optional:            true,
				Computed:            true,
			},
			"header_value": schema.StringAttribute{
				MarkdownDescription: "Content of the header required for interacting with the webhook endpoint.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"owner": schema.Int64Attribute{
				MarkdownDescription: "Owner/receiver of notification. If empty, processed as system notification.",
				Optional:            true,
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Status of the incoming webhook (read-only).",
				Computed:            true,
			},
		},
	}
}

type notificationWebhookResourceData struct {
	Id          types.String `tfsdk:"id" ddField:"Id"`
	Name        types.String `tfsdk:"name" ddField:"Name"`
	Url         types.String `tfsdk:"url" ddField:"Url"`
	HeaderName  types.String `tfsdk:"header_name" ddField:"HeaderName"`
	HeaderValue types.String `tfsdk:"header_value" ddField:"HeaderValue"`
	Owner       types.Int64  `tfsdk:"owner" ddField:"Owner"`
	Status      types.String `tfsdk:"status" ddField:"Status"`
}

type notificationWebhookDefectdojoResource struct {
	dd.NotificationWebhooks
}

func notificationWebhookToRequest(n dd.NotificationWebhooks) dd.NotificationWebhooksRequest {
	return dd.NotificationWebhooksRequest{
		Name:        n.Name,
		Url:         n.Url,
		HeaderName:  n.HeaderName,
		HeaderValue: n.HeaderValue,
		Owner:       n.Owner,
	}
}

func (ddr *notificationWebhookDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := notificationWebhookToRequest(ddr.NotificationWebhooks)
	apiResp, err := client.NotificationWebhooksCreateWithResponse(ctx, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON201 != nil {
		ddr.NotificationWebhooks = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *notificationWebhookDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.NotificationWebhooksRetrieveWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.NotificationWebhooks = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *notificationWebhookDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := notificationWebhookToRequest(ddr.NotificationWebhooks)
	apiResp, err := client.NotificationWebhooksUpdateWithResponse(ctx, idNumber, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.NotificationWebhooks = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *notificationWebhookDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.NotificationWebhooksDestroyWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (d *notificationWebhookResourceData) id() types.String {
	return d.Id
}

func (d *notificationWebhookResourceData) defectdojoResource() defectdojoResource {
	return &notificationWebhookDefectdojoResource{NotificationWebhooks: dd.NotificationWebhooks{}}
}

type notificationWebhookResource struct {
	terraformResource
}

var _ resource.Resource = &notificationWebhookResource{}
var _ resource.ResourceWithImportState = &notificationWebhookResource{}

func NewNotificationWebhookResource() resource.Resource {
	return &notificationWebhookResource{
		terraformResource: terraformResource{
			typeName:     "defectdojo_notification_webhook",
			dataProvider: notificationWebhookDataProvider{},
		},
	}
}

func (r notificationWebhookResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_webhook"
}

type notificationWebhookDataProvider struct{}

func (r notificationWebhookDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data notificationWebhookResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}
