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

func (t productMemberResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Product Member",
		Attributes: map[string]schema.Attribute{
			"product": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Product",
				Required:            true,
			},
			"user": schema.Int64Attribute{
				MarkdownDescription: "The ID of the User",
				Required:            true,
			},
			"role": schema.Int64Attribute{
				MarkdownDescription: "The role ID for this product membership",
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

type productMemberResourceData struct {
	Product types.Int64  `tfsdk:"product" ddField:"Product"`
	User    types.Int64  `tfsdk:"user" ddField:"User"`
	Role    types.Int64  `tfsdk:"role" ddField:"Role"`
	Id      types.String `tfsdk:"id" ddField:"Id"`
}

type productMemberDefectdojoResource struct {
	dd.ProductMember
}

func productMemberToRequest(m dd.ProductMember) dd.ProductMemberRequest {
	return dd.ProductMemberRequest{
		Product: m.Product,
		User:    m.User,
		Role:    m.Role,
	}
}

func (ddr *productMemberDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := productMemberToRequest(ddr.ProductMember)
	apiResp, err := client.ProductMembersCreateWithResponse(ctx, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON201 != nil {
		ddr.ProductMember = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *productMemberDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.ProductMembersRetrieveWithResponse(ctx, idNumber, &dd.ProductMembersRetrieveParams{})
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.ProductMember = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *productMemberDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := productMemberToRequest(ddr.ProductMember)
	apiResp, err := client.ProductMembersUpdateWithResponse(ctx, idNumber, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.ProductMember = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *productMemberDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.ProductMembersDestroyWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

type productMemberResource struct {
	terraformResource
}

var _ resource.Resource = &productMemberResource{}
var _ resource.ResourceWithImportState = &productMemberResource{}

func NewProductMemberResource() resource.Resource {
	return &productMemberResource{
		terraformResource: terraformResource{
			typeName:     "defectdojo_product_member",
			dataProvider: productMemberDataProvider{},
		},
	}
}

func (r productMemberResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_product_member"
}

type productMemberDataProvider struct{}

func (r productMemberDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data productMemberResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *productMemberResourceData) id() types.String {
	return d.Id
}

func (d *productMemberResourceData) setId(v types.String) { d.Id = v }

func (d *productMemberResourceData) defectdojoResource() defectdojoResource {
	return &productMemberDefectdojoResource{ProductMember: dd.ProductMember{}}
}
