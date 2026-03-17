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

func (t credentialResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Credential",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Credential",
				Required:            true,
			},
			"environment": schema.Int64Attribute{
				MarkdownDescription: "The ID of the environment for this credential",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The username for the credential",
				Required:            true,
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "The role for the credential",
				Required:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "The URL for the credential",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the credential",
				Optional:            true,
			},
			"authentication": schema.StringAttribute{
				MarkdownDescription: "Authentication type: 'Form' or 'SSO'",
				Optional:            true,
			},
			"http_authentication": schema.StringAttribute{
				MarkdownDescription: "HTTP authentication type: 'Basic' or 'NTLM'",
				Optional:            true,
			},
			"login_regex": schema.StringAttribute{
				MarkdownDescription: "Login regex pattern",
				Optional:            true,
			},
			"logout_regex": schema.StringAttribute{
				MarkdownDescription: "Logout regex pattern",
				Optional:            true,
			},
			"is_valid": schema.BoolAttribute{
				MarkdownDescription: "Whether the credential is valid",
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

type credentialResourceData struct {
	Name               types.String `tfsdk:"name" ddField:"Name"`
	Environment        types.Int64  `tfsdk:"environment" ddField:"Environment"`
	Username           types.String `tfsdk:"username" ddField:"Username"`
	Role               types.String `tfsdk:"role" ddField:"Role"`
	Url                types.String `tfsdk:"url" ddField:"Url"`
	Description        types.String `tfsdk:"description" ddField:"Description"`
	Authentication     types.String `tfsdk:"authentication" ddField:"Authentication"`
	HttpAuthentication types.String `tfsdk:"http_authentication" ddField:"HttpAuthentication"`
	LoginRegex         types.String `tfsdk:"login_regex" ddField:"LoginRegex"`
	LogoutRegex        types.String `tfsdk:"logout_regex" ddField:"LogoutRegex"`
	IsValid            types.Bool   `tfsdk:"is_valid" ddField:"IsValid"`
	Id                 types.String `tfsdk:"id" ddField:"Id"`
}

type credentialDefectdojoResource struct {
	dd.Credential
}

func credentialToRequest(c dd.Credential) dd.CredentialRequest {
	req := dd.CredentialRequest{
		Name:        c.Name,
		Environment: c.Environment,
		Username:    c.Username,
		Role:        c.Role,
		Url:         c.Url,
		Description: c.Description,
		LoginRegex:  c.LoginRegex,
		LogoutRegex: c.LogoutRegex,
		IsValid:     c.IsValid,
	}
	if c.Authentication != nil {
		v := dd.CredentialRequestAuthentication(*c.Authentication)
		req.Authentication = &v
	}
	if c.HttpAuthentication != nil {
		v := dd.CredentialRequestHttpAuthentication(*c.HttpAuthentication)
		req.HttpAuthentication = &v
	}
	return req
}

func (ddr *credentialDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := credentialToRequest(ddr.Credential)
	apiResp, err := client.CredentialsCreateWithResponse(ctx, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON201 != nil {
		ddr.Credential = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *credentialDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.CredentialsRetrieveWithResponse(ctx, idNumber, &dd.CredentialsRetrieveParams{})
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.Credential = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *credentialDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := credentialToRequest(ddr.Credential)
	apiResp, err := client.CredentialsUpdateWithResponse(ctx, idNumber, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.Credential = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *credentialDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.CredentialsDestroyWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

type credentialResource struct {
	terraformResource
}

var _ resource.Resource = &credentialResource{}
var _ resource.ResourceWithImportState = &credentialResource{}

func NewCredentialResource() resource.Resource {
	return &credentialResource{
		terraformResource: terraformResource{
			dataProvider: credentialDataProvider{},
		},
	}
}

func (r credentialResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential"
}

type credentialDataProvider struct{}

func (r credentialDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data credentialResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *credentialResourceData) id() types.String {
	return d.Id
}

func (d *credentialResourceData) defectdojoResource() defectdojoResource {
	return &credentialDefectdojoResource{Credential: dd.Credential{}}
}
