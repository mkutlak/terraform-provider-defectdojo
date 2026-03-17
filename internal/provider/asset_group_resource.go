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

func (t assetGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Asset Group",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"asset": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Asset.",
				Required:            true,
			},
			"group": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Group.",
				Required:            true,
			},
			"role": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Role.",
				Required:            true,
			},
		},
	}
}

type assetGroupResourceData struct {
	Id    types.String `tfsdk:"id" ddField:"Id"`
	Asset types.Int64  `tfsdk:"asset" ddField:"Asset"`
	Group types.Int64  `tfsdk:"group" ddField:"Group"`
	Role  types.Int64  `tfsdk:"role" ddField:"Role"`
}

type assetGroupDefectdojoResource struct {
	dd.AssetGroup
}

func assetGroupToRequest(a dd.AssetGroup) dd.AssetGroupRequest {
	return dd.AssetGroupRequest{
		Asset: a.Asset,
		Group: a.Group,
		Role:  a.Role,
	}
}

func (ddr *assetGroupDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := assetGroupToRequest(ddr.AssetGroup)
	apiResp, err := client.AssetGroupsCreateWithResponse(ctx, reqBody)
	if apiResp.JSON201 != nil {
		ddr.AssetGroup = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *assetGroupDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.AssetGroupsRetrieveWithResponse(ctx, idNumber, &dd.AssetGroupsRetrieveParams{})
	if apiResp.JSON200 != nil {
		ddr.AssetGroup = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *assetGroupDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := assetGroupToRequest(ddr.AssetGroup)
	apiResp, err := client.AssetGroupsUpdateWithResponse(ctx, idNumber, reqBody)
	if apiResp.JSON200 != nil {
		ddr.AssetGroup = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *assetGroupDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.AssetGroupsDestroyWithResponse(ctx, idNumber)
	return apiResp.StatusCode(), apiResp.Body, err
}

func (d *assetGroupResourceData) id() types.String {
	return d.Id
}

func (d *assetGroupResourceData) defectdojoResource() defectdojoResource {
	return &assetGroupDefectdojoResource{AssetGroup: dd.AssetGroup{}}
}

type assetGroupResource struct {
	terraformResource
}

var _ resource.Resource = &assetGroupResource{}
var _ resource.ResourceWithImportState = &assetGroupResource{}

func NewAssetGroupResource() resource.Resource {
	return &assetGroupResource{
		terraformResource: terraformResource{
			dataProvider: assetGroupDataProvider{},
		},
	}
}

func (r assetGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_asset_group"
}

type assetGroupDataProvider struct{}

func (r assetGroupDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data assetGroupResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}
