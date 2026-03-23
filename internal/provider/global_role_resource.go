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

func (t globalRoleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Global Role",
		Attributes: map[string]schema.Attribute{
			"user": schema.Int64Attribute{
				MarkdownDescription: "The ID of the User to assign the global role to",
				Optional:            true,
			},
			"group": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Dojo Group to assign the global role to",
				Optional:            true,
			},
			"role": schema.Int64Attribute{
				MarkdownDescription: "The global role ID (applied to all product types and products)",
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

type globalRoleResourceData struct {
	User  types.Int64  `tfsdk:"user" ddField:"User"`
	Group types.Int64  `tfsdk:"group" ddField:"Group"`
	Role  types.Int64  `tfsdk:"role" ddField:"Role"`
	Id    types.String `tfsdk:"id" ddField:"Id"`
}

type globalRoleDefectdojoResource struct {
	dd.GlobalRole
}

func globalRoleToRequest(g dd.GlobalRole) dd.GlobalRoleRequest {
	return dd.GlobalRoleRequest{
		User:  g.User,
		Group: g.Group,
		Role:  g.Role,
	}
}

func (ddr *globalRoleDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := globalRoleToRequest(ddr.GlobalRole)
	apiResp, err := client.GlobalRolesCreateWithResponse(ctx, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON201 != nil {
		ddr.GlobalRole = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *globalRoleDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.GlobalRolesRetrieveWithResponse(ctx, idNumber, &dd.GlobalRolesRetrieveParams{})
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.GlobalRole = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *globalRoleDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := globalRoleToRequest(ddr.GlobalRole)
	apiResp, err := client.GlobalRolesUpdateWithResponse(ctx, idNumber, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.GlobalRole = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *globalRoleDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.GlobalRolesDestroyWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

type globalRoleResource struct {
	terraformResource
}

var _ resource.Resource = &globalRoleResource{}
var _ resource.ResourceWithImportState = &globalRoleResource{}

func NewGlobalRoleResource() resource.Resource {
	return &globalRoleResource{
		terraformResource: terraformResource{
			typeName:     "defectdojo_global_role",
			dataProvider: globalRoleDataProvider{},
		},
	}
}

func (r globalRoleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_global_role"
}

type globalRoleDataProvider struct{}

func (r globalRoleDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data globalRoleResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *globalRoleResourceData) id() types.String {
	return d.Id
}

func (d *globalRoleResourceData) defectdojoResource() defectdojoResource {
	return &globalRoleDefectdojoResource{GlobalRole: dd.GlobalRole{}}
}
