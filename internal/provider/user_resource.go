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

func (t userResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo User",
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				MarkdownDescription: "The username of the User",
				Required:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "The email address of the User",
				Required:            true,
			},
			"first_name": schema.StringAttribute{
				MarkdownDescription: "The first name of the User",
				Optional:            true,
			},
			"last_name": schema.StringAttribute{
				MarkdownDescription: "The last name of the User",
				Optional:            true,
			},
			"is_active": schema.BoolAttribute{
				MarkdownDescription: "Whether this user account is active",
				Optional:            true,
				Computed:            true,
			},
			"is_superuser": schema.BoolAttribute{
				MarkdownDescription: "Whether this user has all permissions without explicitly assigning them",
				Optional:            true,
				Computed:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password for the User",
				Optional:            true,
				Sensitive:           true,
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

type userResourceData struct {
	Username    types.String `tfsdk:"username" ddField:"Username"`
	Email       types.String `tfsdk:"email" ddField:"Email"`
	FirstName   types.String `tfsdk:"first_name" ddField:"FirstName"`
	LastName    types.String `tfsdk:"last_name" ddField:"LastName"`
	IsActive    types.Bool   `tfsdk:"is_active" ddField:"IsActive"`
	IsSuperuser types.Bool   `tfsdk:"is_superuser" ddField:"IsSuperuser"`
	Password    types.String `tfsdk:"password" ddField:"Password"`
	Id          types.String `tfsdk:"id" ddField:"Id"`
}

type userDefectdojoResource struct {
	dd.User
	// Password is stored separately since it's not returned in API responses
	password *string
}

func userToRequest(u dd.User, password *string) dd.UserRequest {
	req := dd.UserRequest{
		Username:    u.Username,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		IsActive:    u.IsActive,
		IsSuperuser: u.IsSuperuser,
		Password:    password,
	}
	// Email field: UserRequest uses openapi_types.Email same as User
	req.Email = u.Email
	return req
}

func (ddr *userDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := userToRequest(ddr.User, ddr.password)
	apiResp, err := client.UsersCreateWithResponse(ctx, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON201 != nil {
		ddr.User = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *userDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.UsersRetrieveWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.User = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *userDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := userToRequest(ddr.User, ddr.password)
	apiResp, err := client.UsersUpdateWithResponse(ctx, idNumber, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.User = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *userDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.UsersDestroyWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

type userResource struct {
	terraformResource
}

var _ resource.Resource = &userResource{}
var _ resource.ResourceWithImportState = &userResource{}

func NewUserResource() resource.Resource {
	return &userResource{
		terraformResource: terraformResource{
			dataProvider: userDataProvider{},
		},
	}
}

func (r userResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

type userDataProvider struct{}

func (r userDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data userResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *userResourceData) id() types.String {
	return d.Id
}

func (d *userResourceData) defectdojoResource() defectdojoResource {
	r := &userDefectdojoResource{User: dd.User{}}
	// Password is write-only; store separately so it can be passed to requests
	// but is not overwritten by API responses
	if !d.Password.IsNull() && !d.Password.IsUnknown() {
		p := d.Password.ValueString()
		r.password = &p
	}
	return r
}
