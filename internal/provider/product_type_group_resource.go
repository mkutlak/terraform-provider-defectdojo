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

func (t productTypeGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Product Type Group",
		Attributes: map[string]schema.Attribute{
			"product_type": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Product Type",
				Required:            true,
			},
			"group": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Dojo Group",
				Required:            true,
			},
			"role": schema.Int64Attribute{
				MarkdownDescription: "The role ID for this product type group membership",
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

type productTypeGroupResourceData struct {
	ProductType types.Int64  `tfsdk:"product_type" ddField:"ProductType"`
	Group       types.Int64  `tfsdk:"group" ddField:"Group"`
	Role        types.Int64  `tfsdk:"role" ddField:"Role"`
	Id          types.String `tfsdk:"id" ddField:"Id"`
}

type productTypeGroupDefectdojoResource struct {
	dd.ProductTypeGroup
}

func productTypeGroupToRequest(g dd.ProductTypeGroup) dd.ProductTypeGroupRequest {
	return dd.ProductTypeGroupRequest{
		ProductType: g.ProductType,
		Group:       g.Group,
		Role:        g.Role,
	}
}

func (ddr *productTypeGroupDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := productTypeGroupToRequest(ddr.ProductTypeGroup)
	apiResp, err := client.ProductTypeGroupsCreateWithResponse(ctx, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON201 != nil {
		ddr.ProductTypeGroup = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *productTypeGroupDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.ProductTypeGroupsRetrieveWithResponse(ctx, idNumber, &dd.ProductTypeGroupsRetrieveParams{})
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.ProductTypeGroup = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *productTypeGroupDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := productTypeGroupToRequest(ddr.ProductTypeGroup)
	apiResp, err := client.ProductTypeGroupsUpdateWithResponse(ctx, idNumber, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.ProductTypeGroup = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *productTypeGroupDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.ProductTypeGroupsDestroyWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

type productTypeGroupResource struct {
	terraformResource
}

var _ resource.Resource = &productTypeGroupResource{}
var _ resource.ResourceWithImportState = &productTypeGroupResource{}

func NewProductTypeGroupResource() resource.Resource {
	return &productTypeGroupResource{
		terraformResource: terraformResource{
			typeName:     "defectdojo_product_type_group",
			dataProvider: productTypeGroupDataProvider{},
		},
	}
}

func (r productTypeGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_product_type_group"
}

type productTypeGroupDataProvider struct{}

func (r productTypeGroupDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data productTypeGroupResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *productTypeGroupResourceData) id() types.String {
	return d.Id
}

func (d *productTypeGroupResourceData) setId(v types.String) { d.Id = v }

func (d *productTypeGroupResourceData) defectdojoResource() defectdojoResource {
	return &productTypeGroupDefectdojoResource{ProductTypeGroup: dd.ProductTypeGroup{}}
}
