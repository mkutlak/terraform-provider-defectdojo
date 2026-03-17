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

func (t userContactInfoResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo User Contact Info",
		Attributes: map[string]schema.Attribute{
			"user": schema.Int64Attribute{
				MarkdownDescription: "The ID of the User this contact info belongs to",
				Required:            true,
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "Title of the user",
				Optional:            true,
			},
			"phone_number": schema.StringAttribute{
				MarkdownDescription: "Phone number in format '+999999999'. Up to 15 digits allowed.",
				Optional:            true,
			},
			"cell_number": schema.StringAttribute{
				MarkdownDescription: "Cell number in format '+999999999'. Up to 15 digits allowed.",
				Optional:            true,
			},
			"twitter_username": schema.StringAttribute{
				MarkdownDescription: "Twitter username",
				Optional:            true,
			},
			"github_username": schema.StringAttribute{
				MarkdownDescription: "GitHub username",
				Optional:            true,
			},
			"slack_username": schema.StringAttribute{
				MarkdownDescription: "Email address associated with your Slack account",
				Optional:            true,
			},
			"slack_user_id": schema.StringAttribute{
				MarkdownDescription: "Slack user ID",
				Optional:            true,
			},
			"block_execution": schema.BoolAttribute{
				MarkdownDescription: "Instead of async deduping a finding the findings will be deduped synchronously",
				Optional:            true,
				Computed:            true,
			},
			"force_password_reset": schema.BoolAttribute{
				MarkdownDescription: "Forces this user to reset their password on next login",
				Optional:            true,
				Computed:            true,
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

type userContactInfoResourceData struct {
	User               types.Int64  `tfsdk:"user" ddField:"User"`
	Title              types.String `tfsdk:"title" ddField:"Title"`
	PhoneNumber        types.String `tfsdk:"phone_number" ddField:"PhoneNumber"`
	CellNumber         types.String `tfsdk:"cell_number" ddField:"CellNumber"`
	TwitterUsername    types.String `tfsdk:"twitter_username" ddField:"TwitterUsername"`
	GithubUsername     types.String `tfsdk:"github_username" ddField:"GithubUsername"`
	SlackUsername      types.String `tfsdk:"slack_username" ddField:"SlackUsername"`
	SlackUserId        types.String `tfsdk:"slack_user_id" ddField:"SlackUserId"`
	BlockExecution     types.Bool   `tfsdk:"block_execution" ddField:"BlockExecution"`
	ForcePasswordReset types.Bool   `tfsdk:"force_password_reset" ddField:"ForcePasswordReset"`
	Id                 types.String `tfsdk:"id" ddField:"Id"`
}

type userContactInfoDefectdojoResource struct {
	dd.UserContactInfo
}

func userContactInfoToRequest(u dd.UserContactInfo) dd.UserContactInfoRequest {
	return dd.UserContactInfoRequest{
		User:               u.User,
		Title:              u.Title,
		PhoneNumber:        u.PhoneNumber,
		CellNumber:         u.CellNumber,
		TwitterUsername:    u.TwitterUsername,
		GithubUsername:     u.GithubUsername,
		SlackUsername:      u.SlackUsername,
		SlackUserId:        u.SlackUserId,
		BlockExecution:     u.BlockExecution,
		ForcePasswordReset: u.ForcePasswordReset,
	}
}

func (ddr *userContactInfoDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := userContactInfoToRequest(ddr.UserContactInfo)
	apiResp, err := client.UserContactInfosCreateWithResponse(ctx, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON201 != nil {
		ddr.UserContactInfo = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *userContactInfoDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.UserContactInfosRetrieveWithResponse(ctx, idNumber, &dd.UserContactInfosRetrieveParams{})
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.UserContactInfo = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *userContactInfoDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := userContactInfoToRequest(ddr.UserContactInfo)
	apiResp, err := client.UserContactInfosUpdateWithResponse(ctx, idNumber, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.UserContactInfo = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *userContactInfoDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.UserContactInfosDestroyWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

type userContactInfoResource struct {
	terraformResource
}

var _ resource.Resource = &userContactInfoResource{}
var _ resource.ResourceWithImportState = &userContactInfoResource{}

func NewUserContactInfoResource() resource.Resource {
	return &userContactInfoResource{
		terraformResource: terraformResource{
			dataProvider: userContactInfoDataProvider{},
		},
	}
}

func (r userContactInfoResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_contact_info"
}

type userContactInfoDataProvider struct{}

func (r userContactInfoDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data userContactInfoResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *userContactInfoResourceData) id() types.String {
	return d.Id
}

func (d *userContactInfoResourceData) defectdojoResource() defectdojoResource {
	return &userContactInfoDefectdojoResource{UserContactInfo: dd.UserContactInfo{}}
}
