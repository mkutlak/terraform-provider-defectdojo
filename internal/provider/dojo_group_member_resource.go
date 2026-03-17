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

func (t dojoGroupMemberResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Dojo Group Member",
		Attributes: map[string]schema.Attribute{
			"group": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Dojo Group",
				Required:            true,
			},
			"user": schema.Int64Attribute{
				MarkdownDescription: "The ID of the User",
				Required:            true,
			},
			"role": schema.Int64Attribute{
				MarkdownDescription: "The role determining the permissions of the user to manage the group",
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

type dojoGroupMemberResourceData struct {
	Group types.Int64  `tfsdk:"group" ddField:"Group"`
	User  types.Int64  `tfsdk:"user" ddField:"User"`
	Role  types.Int64  `tfsdk:"role" ddField:"Role"`
	Id    types.String `tfsdk:"id" ddField:"Id"`
}

type dojoGroupMemberDefectdojoResource struct {
	dd.DojoGroupMember
}

func dojoGroupMemberToRequest(m dd.DojoGroupMember) dd.DojoGroupMemberRequest {
	return dd.DojoGroupMemberRequest{
		Group: m.Group,
		User:  m.User,
		Role:  m.Role,
	}
}

func (ddr *dojoGroupMemberDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := dojoGroupMemberToRequest(ddr.DojoGroupMember)
	apiResp, err := client.DojoGroupMembersCreateWithResponse(ctx, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON201 != nil {
		ddr.DojoGroupMember = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *dojoGroupMemberDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.DojoGroupMembersRetrieveWithResponse(ctx, idNumber, &dd.DojoGroupMembersRetrieveParams{})
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.DojoGroupMember = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *dojoGroupMemberDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := dojoGroupMemberToRequest(ddr.DojoGroupMember)
	apiResp, err := client.DojoGroupMembersUpdateWithResponse(ctx, idNumber, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.DojoGroupMember = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *dojoGroupMemberDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.DojoGroupMembersDestroyWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

type dojoGroupMemberResource struct {
	terraformResource
}

var _ resource.Resource = &dojoGroupMemberResource{}
var _ resource.ResourceWithImportState = &dojoGroupMemberResource{}

func NewDojoGroupMemberResource() resource.Resource {
	return &dojoGroupMemberResource{
		terraformResource: terraformResource{
			dataProvider: dojoGroupMemberDataProvider{},
		},
	}
}

func (r dojoGroupMemberResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dojo_group_member"
}

type dojoGroupMemberDataProvider struct{}

func (r dojoGroupMemberDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data dojoGroupMemberResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *dojoGroupMemberResourceData) id() types.String {
	return d.Id
}

func (d *dojoGroupMemberResourceData) defectdojoResource() defectdojoResource {
	return &dojoGroupMemberDefectdojoResource{DojoGroupMember: dd.DojoGroupMember{}}
}
