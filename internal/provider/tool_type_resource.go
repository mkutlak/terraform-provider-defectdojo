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

func (t toolTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Tool Type",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Tool Type",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the Tool Type",
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

type toolTypeResourceData struct {
	Name        types.String `tfsdk:"name" ddField:"Name"`
	Description types.String `tfsdk:"description" ddField:"Description"`
	Id          types.String `tfsdk:"id" ddField:"Id"`
}

type toolTypeDefectdojoResource struct {
	dd.ToolType
}

func toolTypeToRequest(obj dd.ToolType) dd.ToolTypeRequest {
	return dd.ToolTypeRequest{
		Name:        obj.Name,
		Description: obj.Description,
	}
}

func (ddr *toolTypeDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := toolTypeToRequest(ddr.ToolType)
	apiResp, err := client.ToolTypesCreateWithResponse(ctx, reqBody)
	if apiResp.JSON201 != nil {
		ddr.ToolType = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *toolTypeDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.ToolTypesRetrieveWithResponse(ctx, idNumber)
	if apiResp.JSON200 != nil {
		ddr.ToolType = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *toolTypeDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := toolTypeToRequest(ddr.ToolType)
	apiResp, err := client.ToolTypesUpdateWithResponse(ctx, idNumber, reqBody)
	if apiResp.JSON200 != nil {
		ddr.ToolType = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *toolTypeDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.ToolTypesDestroyWithResponse(ctx, idNumber)
	return apiResp.StatusCode(), apiResp.Body, err
}

type toolTypeResource struct {
	terraformResource
}

var _ resource.Resource = &toolTypeResource{}
var _ resource.ResourceWithImportState = &toolTypeResource{}

func NewToolTypeResource() resource.Resource {
	return &toolTypeResource{
		terraformResource: terraformResource{dataProvider: toolTypeDataProvider{}},
	}
}

func (r toolTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tool_type"
}

type toolTypeDataProvider struct{}

func (r toolTypeDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data toolTypeResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *toolTypeResourceData) id() types.String { return d.Id }

func (d *toolTypeResourceData) defectdojoResource() defectdojoResource {
	return &toolTypeDefectdojoResource{ToolType: dd.ToolType{}}
}
