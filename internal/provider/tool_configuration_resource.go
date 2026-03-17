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

func (t toolConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Tool Configuration",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Tool Configuration",
				Required:            true,
			},
			"tool_type": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Tool Type",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the Tool Configuration",
				Optional:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "The URL of the tool",
				Optional:            true,
			},
			"authentication_type": schema.StringAttribute{
				MarkdownDescription: "Authentication type. Valid values: 'API', 'Password', 'SSH'",
				Optional:            true,
			},
			"auth_title": schema.StringAttribute{
				MarkdownDescription: "Title for authentication credentials",
				Optional:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username for authentication",
				Optional:            true,
			},
			"extras": schema.StringAttribute{
				MarkdownDescription: "Additional definitions that will be consumed by scanner",
				Optional:            true,
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

type toolConfigurationResourceData struct {
	Name               types.String `tfsdk:"name" ddField:"Name"`
	ToolType           types.Int64  `tfsdk:"tool_type" ddField:"ToolType"`
	Description        types.String `tfsdk:"description" ddField:"Description"`
	Url                types.String `tfsdk:"url" ddField:"Url"`
	AuthenticationType types.String `tfsdk:"authentication_type" ddField:"AuthenticationType"`
	AuthTitle          types.String `tfsdk:"auth_title" ddField:"AuthTitle"`
	Username           types.String `tfsdk:"username" ddField:"Username"`
	Extras             types.String `tfsdk:"extras" ddField:"Extras"`
	Id                 types.String `tfsdk:"id" ddField:"Id"`
}

type toolConfigurationDefectdojoResource struct {
	dd.ToolConfiguration
}

func toolConfigurationToRequest(obj dd.ToolConfiguration) dd.ToolConfigurationRequest {
	req := dd.ToolConfigurationRequest{
		Name:        obj.Name,
		ToolType:    obj.ToolType,
		Description: obj.Description,
		Url:         obj.Url,
		AuthTitle:   obj.AuthTitle,
		Username:    obj.Username,
		Extras:      obj.Extras,
	}
	if obj.AuthenticationType != nil {
		v := dd.ToolConfigurationRequestAuthenticationType(*obj.AuthenticationType)
		req.AuthenticationType = &v
	}
	return req
}

func (ddr *toolConfigurationDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := toolConfigurationToRequest(ddr.ToolConfiguration)
	apiResp, err := client.ToolConfigurationsCreateWithResponse(ctx, reqBody)
	if apiResp.JSON201 != nil {
		ddr.ToolConfiguration = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *toolConfigurationDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.ToolConfigurationsRetrieveWithResponse(ctx, idNumber, &dd.ToolConfigurationsRetrieveParams{})
	if apiResp.JSON200 != nil {
		ddr.ToolConfiguration = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *toolConfigurationDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := toolConfigurationToRequest(ddr.ToolConfiguration)
	apiResp, err := client.ToolConfigurationsUpdateWithResponse(ctx, idNumber, reqBody)
	if apiResp.JSON200 != nil {
		ddr.ToolConfiguration = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *toolConfigurationDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.ToolConfigurationsDestroyWithResponse(ctx, idNumber)
	return apiResp.StatusCode(), apiResp.Body, err
}

type toolConfigurationResource struct {
	terraformResource
}

var _ resource.Resource = &toolConfigurationResource{}
var _ resource.ResourceWithImportState = &toolConfigurationResource{}

func NewToolConfigurationResource() resource.Resource {
	return &toolConfigurationResource{
		terraformResource: terraformResource{dataProvider: toolConfigurationDataProvider{}},
	}
}

func (r toolConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tool_configuration"
}

type toolConfigurationDataProvider struct{}

func (r toolConfigurationDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data toolConfigurationResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *toolConfigurationResourceData) id() types.String { return d.Id }

func (d *toolConfigurationResourceData) defectdojoResource() defectdojoResource {
	return &toolConfigurationDefectdojoResource{ToolConfiguration: dd.ToolConfiguration{}}
}
