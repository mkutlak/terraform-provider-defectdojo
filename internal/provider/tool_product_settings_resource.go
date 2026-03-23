package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
)

func (t toolProductSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Tool Product Settings",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the tool product settings.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the tool product settings.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "URL for the tool product settings. Automatically set from setting_url by the API.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"setting_url": schema.StringAttribute{
				MarkdownDescription: "The settings URL for the tool.",
				Required:            true,
			},
			"product": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Product.",
				Required:            true,
			},
			"tool_configuration": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Tool Configuration.",
				Required:            true,
			},
			"tool_project_id": schema.StringAttribute{
				MarkdownDescription: "The project ID in the tool.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
		},
	}
}

type toolProductSettingsResourceData struct {
	Id                types.String `tfsdk:"id" ddField:"Id"`
	Name              types.String `tfsdk:"name" ddField:"Name"`
	Description       types.String `tfsdk:"description" ddField:"Description"`
	Url               types.String `tfsdk:"url" ddField:"Url"`
	SettingUrl        types.String `tfsdk:"setting_url" ddField:"SettingUrl"`
	Product           types.Int64  `tfsdk:"product" ddField:"Product"`
	ToolConfiguration types.Int64  `tfsdk:"tool_configuration" ddField:"ToolConfiguration"`
	ToolProjectId     types.String `tfsdk:"tool_project_id" ddField:"ToolProjectId"`
}

type toolProductSettingsDefectdojoResource struct {
	dd.ToolProductSettings
}

func toolProductSettingsToRequest(t dd.ToolProductSettings) dd.ToolProductSettingsRequest {
	return dd.ToolProductSettingsRequest{
		Name:              t.Name,
		Description:       t.Description,
		SettingUrl:        t.SettingUrl,
		Product:           t.Product,
		ToolConfiguration: t.ToolConfiguration,
		ToolProjectId:     t.ToolProjectId,
		// Url intentionally omitted — the DD API auto-copies setting_url to url.
		// Sending url (even empty) causes the API to overwrite setting_url.
	}
}

func (ddr *toolProductSettingsDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := toolProductSettingsToRequest(ddr.ToolProductSettings)
	apiResp, err := client.ToolProductSettingsCreateWithResponse(ctx, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON201 != nil {
		ddr.ToolProductSettings = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *toolProductSettingsDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.ToolProductSettingsRetrieveWithResponse(ctx, idNumber, &dd.ToolProductSettingsRetrieveParams{})
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.ToolProductSettings = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *toolProductSettingsDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := toolProductSettingsToRequest(ddr.ToolProductSettings)
	apiResp, err := client.ToolProductSettingsUpdateWithResponse(ctx, idNumber, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.ToolProductSettings = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *toolProductSettingsDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.ToolProductSettingsDestroyWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (d *toolProductSettingsResourceData) id() types.String {
	return d.Id
}

func (d *toolProductSettingsResourceData) setId(v types.String) { d.Id = v }

func (d *toolProductSettingsResourceData) defectdojoResource() defectdojoResource {
	return &toolProductSettingsDefectdojoResource{ToolProductSettings: dd.ToolProductSettings{}}
}

type toolProductSettingsResource struct {
	terraformResource
}

var _ resource.Resource = &toolProductSettingsResource{}
var _ resource.ResourceWithImportState = &toolProductSettingsResource{}

func NewToolProductSettingsResource() resource.Resource {
	return &toolProductSettingsResource{
		terraformResource: terraformResource{
			typeName:     "defectdojo_tool_product_settings",
			dataProvider: toolProductSettingsDataProvider{},
		},
	}
}

func (r toolProductSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tool_product_settings"
}

type toolProductSettingsDataProvider struct{}

func (r toolProductSettingsDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data toolProductSettingsResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}
