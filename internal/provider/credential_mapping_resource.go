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

func (t credentialMappingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Credential Mapping",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cred_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Credential.",
				Required:            true,
			},
			"product": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Product to map the credential to.",
				Optional:            true,
				Computed:            true,
			},
			"engagement": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Engagement to map the credential to.",
				Optional:            true,
				Computed:            true,
			},
			"finding": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Finding to map the credential to.",
				Optional:            true,
				Computed:            true,
			},
			"test": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Test to map the credential to.",
				Optional:            true,
				Computed:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "URL associated with the credential mapping.",
				Optional:            true,
				Computed:            true,
			},
			"is_authn_provider": schema.BoolAttribute{
				MarkdownDescription: "Whether this is an authentication provider.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

type credentialMappingResourceData struct {
	Id              types.String `tfsdk:"id" ddField:"Id"`
	CredId          types.Int64  `tfsdk:"cred_id" ddField:"CredId"`
	Product         types.Int64  `tfsdk:"product" ddField:"Product"`
	Engagement      types.Int64  `tfsdk:"engagement" ddField:"Engagement"`
	Finding         types.Int64  `tfsdk:"finding" ddField:"Finding"`
	Test            types.Int64  `tfsdk:"test" ddField:"Test"`
	Url             types.String `tfsdk:"url" ddField:"Url"`
	IsAuthnProvider types.Bool   `tfsdk:"is_authn_provider" ddField:"IsAuthnProvider"`
}

type credentialMappingDefectdojoResource struct {
	dd.CredentialMapping
}

func credentialMappingToRequest(c dd.CredentialMapping) dd.CredentialMappingRequest {
	return dd.CredentialMappingRequest{
		CredId:          c.CredId,
		Product:         c.Product,
		Engagement:      c.Engagement,
		Finding:         c.Finding,
		Test:            c.Test,
		Url:             c.Url,
		IsAuthnProvider: c.IsAuthnProvider,
	}
}

func (ddr *credentialMappingDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := credentialMappingToRequest(ddr.CredentialMapping)
	apiResp, err := client.CredentialMappingsCreateWithResponse(ctx, reqBody)
	if apiResp.JSON201 != nil {
		ddr.CredentialMapping = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *credentialMappingDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.CredentialMappingsRetrieveWithResponse(ctx, idNumber)
	if apiResp.JSON200 != nil {
		ddr.CredentialMapping = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *credentialMappingDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := credentialMappingToRequest(ddr.CredentialMapping)
	apiResp, err := client.CredentialMappingsUpdateWithResponse(ctx, idNumber, reqBody)
	if apiResp.JSON200 != nil {
		ddr.CredentialMapping = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *credentialMappingDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.CredentialMappingsDestroyWithResponse(ctx, idNumber)
	return apiResp.StatusCode(), apiResp.Body, err
}

func (d *credentialMappingResourceData) id() types.String {
	return d.Id
}

func (d *credentialMappingResourceData) defectdojoResource() defectdojoResource {
	return &credentialMappingDefectdojoResource{CredentialMapping: dd.CredentialMapping{}}
}

type credentialMappingResource struct {
	terraformResource
}

var _ resource.Resource = &credentialMappingResource{}
var _ resource.ResourceWithImportState = &credentialMappingResource{}

func NewCredentialMappingResource() resource.Resource {
	return &credentialMappingResource{
		terraformResource: terraformResource{
			dataProvider: credentialMappingDataProvider{},
		},
	}
}

func (r credentialMappingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_mapping"
}

type credentialMappingDataProvider struct{}

func (r credentialMappingDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data credentialMappingResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}
