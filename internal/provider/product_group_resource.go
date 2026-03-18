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

func (t productGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Product Group",
		Attributes: map[string]schema.Attribute{
			"product": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Product",
				Required:            true,
			},
			"group": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Dojo Group",
				Required:            true,
			},
			"role": schema.Int64Attribute{
				MarkdownDescription: "The role ID for this product group membership",
				Required:            true,
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

type productGroupResourceData struct {
	Product types.Int64  `tfsdk:"product" ddField:"Product"`
	Group   types.Int64  `tfsdk:"group" ddField:"Group"`
	Role    types.Int64  `tfsdk:"role" ddField:"Role"`
	Id      types.String `tfsdk:"id" ddField:"Id"`
}

type productGroupDefectdojoResource struct {
	dd.ProductGroup
}

func productGroupToRequest(g dd.ProductGroup) dd.ProductGroupRequest {
	return dd.ProductGroupRequest{
		Product: g.Product,
		Group:   g.Group,
		Role:    g.Role,
	}
}

func (ddr *productGroupDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := productGroupToRequest(ddr.ProductGroup)
	apiResp, err := client.ProductGroupsCreateWithResponse(ctx, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON201 != nil {
		ddr.ProductGroup = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *productGroupDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.ProductGroupsRetrieveWithResponse(ctx, idNumber, &dd.ProductGroupsRetrieveParams{})
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.ProductGroup = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *productGroupDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := productGroupToRequest(ddr.ProductGroup)
	apiResp, err := client.ProductGroupsUpdateWithResponse(ctx, idNumber, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.ProductGroup = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *productGroupDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.ProductGroupsDestroyWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

type productGroupResource struct {
	terraformResource
}

var _ resource.Resource = &productGroupResource{}
var _ resource.ResourceWithImportState = &productGroupResource{}

func NewProductGroupResource() resource.Resource {
	return &productGroupResource{
		terraformResource: terraformResource{
			dataProvider: productGroupDataProvider{},
		},
	}
}

func (r productGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_product_group"
}

type productGroupDataProvider struct{}

func (r productGroupDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data productGroupResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *productGroupResourceData) id() types.String {
	return d.Id
}

func (d *productGroupResourceData) defectdojoResource() defectdojoResource {
	return &productGroupDefectdojoResource{ProductGroup: dd.ProductGroup{}}
}
