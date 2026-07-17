package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
)

func (t announcementResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Announcement. DefectDojo allows a single global announcement; if one already exists, import it: `terraform import defectdojo_announcement.example <id>`.",
		Attributes: map[string]schema.Attribute{
			"message": schema.StringAttribute{
				MarkdownDescription: "This dismissable message will be displayed on all pages for authenticated users. It can contain basic html tags.",
				Required:            true,
			},
			"style": schema.StringAttribute{
				MarkdownDescription: "The style of banner to display. Valid values: 'info', 'success', 'warning', 'danger'.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("info", "success", "warning", "danger"),
				},
			},
			"dismissable": schema.BoolAttribute{
				MarkdownDescription: "Ticking this box allows users to dismiss the current announcement",
				Optional:            true,
				Computed:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

type announcementResourceData struct {
	Message     types.String `tfsdk:"message" ddField:"Message"`
	Style       types.String `tfsdk:"style" ddField:"Style"`
	Dismissable types.Bool   `tfsdk:"dismissable" ddField:"Dismissable"`
	Id          types.String `tfsdk:"id" ddField:"Id"`
}

type announcementDefectdojoResource struct {
	dd.Announcement
}

func announcementToRequest(obj dd.Announcement) dd.AnnouncementRequest {
	req := dd.AnnouncementRequest{
		Dismissable: obj.Dismissable,
		Message:     obj.Message,
	}
	if obj.Style != nil {
		style := dd.AnnouncementRequestStyle(string(*obj.Style))
		req.Style = &style
	}
	return req
}

func (ddr *announcementDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := announcementToRequest(ddr.Announcement)
	apiResp, err := client.AnnouncementsCreateWithResponse(ctx, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON201 != nil {
		ddr.Announcement = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *announcementDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.AnnouncementsRetrieveWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.Announcement = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *announcementDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := announcementToRequest(ddr.Announcement)
	apiResp, err := client.AnnouncementsUpdateWithResponse(ctx, idNumber, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.Announcement = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *announcementDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.AnnouncementsDestroyWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

type announcementResource struct {
	terraformResource
}

var _ resource.Resource = &announcementResource{}
var _ resource.ResourceWithImportState = &announcementResource{}

func NewAnnouncementResource() resource.Resource {
	return &announcementResource{
		terraformResource: terraformResource{typeName: "defectdojo_announcement", dataProvider: announcementDataProvider{}},
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

func (d *announcementResourceData) id() types.String     { return d.Id }
func (d *announcementResourceData) setId(v types.String) { d.Id = v }

func (d *announcementResourceData) defectdojoResource() defectdojoResource {
	return &announcementDefectdojoResource{Announcement: dd.Announcement{}}
}
