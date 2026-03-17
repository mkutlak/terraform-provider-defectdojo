package provider

import (
	"context"

	dd "github.com/doximity/terraform-provider-defectdojo/internal/ddclient"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (t announcementResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Announcement",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"message": schema.StringAttribute{
				MarkdownDescription: "The announcement message (can contain basic HTML tags).",
				Optional:            true,
				Computed:            true,
			},
			"style": schema.StringAttribute{
				MarkdownDescription: "The style of banner to display. Valid values: info, success, warning, danger.",
				Optional:            true,
				Computed:            true,
			},
			"dismissable": schema.BoolAttribute{
				MarkdownDescription: "Whether users can dismiss the announcement.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

type announcementResourceData struct {
	Id          types.String `tfsdk:"id" ddField:"Id"`
	Message     types.String `tfsdk:"message" ddField:"Message"`
	Style       types.String `tfsdk:"style" ddField:"Style"`
	Dismissable types.Bool   `tfsdk:"dismissable" ddField:"Dismissable"`
}

type announcementDefectdojoResource struct {
	dd.Announcement
}

func announcementToRequest(a dd.Announcement) dd.AnnouncementRequest {
	req := dd.AnnouncementRequest{
		Message:     a.Message,
		Dismissable: a.Dismissable,
	}
	if a.Style != nil {
		v := dd.AnnouncementRequestStyle(*a.Style)
		req.Style = &v
	}
	return req
}

func (ddr *announcementDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := announcementToRequest(ddr.Announcement)
	apiResp, err := client.AnnouncementsCreateWithResponse(ctx, reqBody)
	if apiResp.JSON201 != nil {
		ddr.Announcement = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *announcementDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.AnnouncementsRetrieveWithResponse(ctx, idNumber)
	if apiResp.JSON200 != nil {
		ddr.Announcement = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *announcementDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := announcementToRequest(ddr.Announcement)
	apiResp, err := client.AnnouncementsUpdateWithResponse(ctx, idNumber, reqBody)
	if apiResp.JSON200 != nil {
		ddr.Announcement = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *announcementDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.AnnouncementsDestroyWithResponse(ctx, idNumber)
	return apiResp.StatusCode(), apiResp.Body, err
}

func (d *announcementResourceData) id() types.String {
	return d.Id
}

func (d *announcementResourceData) defectdojoResource() defectdojoResource {
	return &announcementDefectdojoResource{Announcement: dd.Announcement{}}
}

type announcementResource struct {
	terraformResource
}

var _ resource.Resource = &announcementResource{}
var _ resource.ResourceWithImportState = &announcementResource{}

func NewAnnouncementResource() resource.Resource {
	return &announcementResource{
		terraformResource: terraformResource{
			dataProvider: announcementDataProvider{},
		},
	}
}

func (r announcementResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_announcement"
}

type announcementDataProvider struct{}

func (r announcementDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data announcementResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}
