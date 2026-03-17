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

func (t productTypeMemberResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Product Type Member",
		Attributes: map[string]schema.Attribute{
			"product_type": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Product Type",
				Required:            true,
			},
			"user": schema.Int64Attribute{
				MarkdownDescription: "The ID of the User",
				Required:            true,
			},
			"role": schema.Int64Attribute{
				MarkdownDescription: "The role ID for this product type membership",
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

type productTypeMemberResourceData struct {
	ProductType types.Int64  `tfsdk:"product_type" ddField:"ProductType"`
	User        types.Int64  `tfsdk:"user" ddField:"User"`
	Role        types.Int64  `tfsdk:"role" ddField:"Role"`
	Id          types.String `tfsdk:"id" ddField:"Id"`
}

type productTypeMemberDefectdojoResource struct {
	dd.ProductTypeMember
}

func productTypeMemberToRequest(m dd.ProductTypeMember) dd.ProductTypeMemberRequest {
	return dd.ProductTypeMemberRequest{
		ProductType: m.ProductType,
		User:        m.User,
		Role:        m.Role,
	}
}

func (ddr *productTypeMemberDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := productTypeMemberToRequest(ddr.ProductTypeMember)
	apiResp, err := client.ProductTypeMembersCreateWithResponse(ctx, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON201 != nil {
		ddr.ProductTypeMember = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *productTypeMemberDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.ProductTypeMembersRetrieveWithResponse(ctx, idNumber, &dd.ProductTypeMembersRetrieveParams{})
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.ProductTypeMember = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *productTypeMemberDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := productTypeMemberToRequest(ddr.ProductTypeMember)
	apiResp, err := client.ProductTypeMembersUpdateWithResponse(ctx, idNumber, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.ProductTypeMember = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *productTypeMemberDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.ProductTypeMembersDestroyWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

type productTypeMemberResource struct {
	terraformResource
}

var _ resource.Resource = &productTypeMemberResource{}
var _ resource.ResourceWithImportState = &productTypeMemberResource{}

func NewProductTypeMemberResource() resource.Resource {
	return &productTypeMemberResource{
		terraformResource: terraformResource{
			dataProvider: productTypeMemberDataProvider{},
		},
	}
}

func (r productTypeMemberResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_product_type_member"
}

type productTypeMemberDataProvider struct{}

func (r productTypeMemberDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data productTypeMemberResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *productTypeMemberResourceData) id() types.String {
	return d.Id
}

func (d *productTypeMemberResourceData) defectdojoResource() defectdojoResource {
	return &productTypeMemberDefectdojoResource{ProductTypeMember: dd.ProductTypeMember{}}
}
