package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
)

// userProfileDataSource exposes GET /api/v2/user_profile/, which always
// describes the authenticated user (the caller identified by the provider's
// credentials). It takes no lookup parameters, so it overrides the generic
// terraformDatasource.Read instead of relying on the id/name resolution
// logic in datasource.go.
type userProfileDataSource struct {
	terraformDatasource
}

func (t userProfileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for the authenticated DefectDojo user's profile (`/api/v2/user_profile/`). This always describes the user identified by the provider's credentials and takes no lookup parameters. Contact information is available via the `defectdojo_user_contact_info` data source.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the authenticated user",
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The username of the authenticated user",
				Computed:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "The email address of the authenticated user",
				Computed:            true,
			},
			"first_name": schema.StringAttribute{
				MarkdownDescription: "The first name of the authenticated user",
				Computed:            true,
			},
			"last_name": schema.StringAttribute{
				MarkdownDescription: "The last name of the authenticated user",
				Computed:            true,
			},
			"is_active": schema.BoolAttribute{
				MarkdownDescription: "Whether this user account is active",
				Computed:            true,
			},
			"is_staff": schema.BoolAttribute{
				MarkdownDescription: "Whether the user can log into the admin site",
				Computed:            true,
			},
			"is_superuser": schema.BoolAttribute{
				MarkdownDescription: "Whether this user has all permissions without explicitly assigning them",
				Computed:            true,
			},
			"date_joined": schema.StringAttribute{
				MarkdownDescription: "Timestamp (RFC3339) of when the user account was created",
				Computed:            true,
			},
			"last_login": schema.StringAttribute{
				MarkdownDescription: "Timestamp (RFC3339) of the user's last login",
				Computed:            true,
			},
		},
	}
}

func (d userProfileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_profile"
}

// Read overrides the generic terraformDatasource.Read: that generic
// implementation resolves an id (or a name-based lookup) before reading,
// but /api/v2/user_profile/ takes no identifier at all - it always returns
// the authenticated user. Configure is still inherited from the embedded
// terraformDatasource, so d.client is populated normally.
func (d userProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data, diags := d.getData(ctx, req.Config)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	ddResource := data.defectdojoResource()

	statusCode, body, err := ddResource.readApiCall(ctx, d.client, 0)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Resource",
			err.Error())
		return
	}

	if statusCode != 200 {
		resp.Diagnostics.AddError(
			"API Error Retrieving Resource",
			fmt.Sprintf("Unexpected response code from API: %d", statusCode)+
				fmt.Sprintf("\n\nbody:\n\n%+v", string(body)),
		)
		return
	}

	populateResourceData(ctx, &diags, &data, ddResource)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

var _ datasource.DataSource = &userProfileDataSource{}

func NewUserProfileDataSource() datasource.DataSource {
	return &userProfileDataSource{
		terraformDatasource: terraformDatasource{dataProvider: userProfileDataProvider{}},
	}
}

// userProfileResourceData maps the fields of dd.User (the "user" object
// flattened out of dd.UserProfile), since the reflection-based CRUD engine
// cannot map the nested UserProfile struct directly. ConfigurationPermissions
// (*[]*int) is unsupported by the reflection engine and is intentionally
// omitted, as are password-like fields (none are returned by this endpoint).
type userProfileResourceData struct {
	Id          types.String `tfsdk:"id" ddField:"Id"`
	Username    types.String `tfsdk:"username" ddField:"Username"`
	Email       types.String `tfsdk:"email" ddField:"Email"`
	FirstName   types.String `tfsdk:"first_name" ddField:"FirstName"`
	LastName    types.String `tfsdk:"last_name" ddField:"LastName"`
	IsActive    types.Bool   `tfsdk:"is_active" ddField:"IsActive"`
	IsStaff     types.Bool   `tfsdk:"is_staff" ddField:"IsStaff"`
	IsSuperuser types.Bool   `tfsdk:"is_superuser" ddField:"IsSuperuser"`
	DateJoined  types.String `tfsdk:"date_joined" ddField:"DateJoined"`
	LastLogin   types.String `tfsdk:"last_login" ddField:"LastLogin"`
}

// userProfileDefectdojoResource wraps dd.User, the "user" object nested
// inside dd.UserProfile's response envelope (which also carries an optional
// UserContactInfo, exposed separately via defectdojo_user_contact_info).
type userProfileDefectdojoResource struct {
	dd.User
}

// errUserProfileReadOnly is returned by the write stubs below. The
// user_profile endpoint only supports GET: it always describes the
// authenticated user and cannot be created, updated, or deleted.
var errUserProfileReadOnly = errors.New("user_profile is a read-only data source describing the authenticated user")

func (ddr *userProfileDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	return 0, nil, errUserProfileReadOnly
}

// readApiCall ignores the id parameter: GET /api/v2/user_profile/ takes no
// identifier and always returns the authenticated user's profile.
func (ddr *userProfileDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, _ int) (int, []byte, error) {
	apiResp, err := client.UserProfileRetrieveWithResponse(ctx)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.User = apiResp.JSON200.User
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *userProfileDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	return 0, nil, errUserProfileReadOnly
}

func (ddr *userProfileDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	return 0, nil, errUserProfileReadOnly
}

type userProfileDataProvider struct{}

func (r userProfileDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data userProfileResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *userProfileResourceData) id() types.String { return d.Id }

func (d *userProfileResourceData) setId(v types.String) { d.Id = v }

func (d *userProfileResourceData) defectdojoResource() defectdojoResource {
	return &userProfileDefectdojoResource{User: dd.User{}}
}
