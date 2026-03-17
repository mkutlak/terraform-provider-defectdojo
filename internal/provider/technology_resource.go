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

func (t technologyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Technology (App Analysis)",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the technology",
				Required:            true,
			},
			"product": schema.Int64Attribute{
				MarkdownDescription: "The product this technology is associated with.",
				Required:            true,
			},
			"user": schema.Int64Attribute{
				MarkdownDescription: "The user who added this technology.",
				Required:            true,
			},
			"confidence": schema.Int64Attribute{
				MarkdownDescription: "Confidence level of the detection.",
				Optional:            true,
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "Version of the technology.",
				Optional:            true,
			},
			"icon": schema.StringAttribute{
				MarkdownDescription: "Icon for the technology.",
				Optional:            true,
			},
			"website": schema.StringAttribute{
				MarkdownDescription: "Website of the technology.",
				Optional:            true,
			},
			"website_found": schema.StringAttribute{
				MarkdownDescription: "Website where the technology was found.",
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

type technologyResourceData struct {
	Name         types.String `tfsdk:"name" ddField:"Name"`
	Product      types.Int64  `tfsdk:"product" ddField:"Product"`
	User         types.Int64  `tfsdk:"user" ddField:"User"`
	Confidence   types.Int64  `tfsdk:"confidence" ddField:"Confidence"`
	Version      types.String `tfsdk:"version" ddField:"Version"`
	Icon         types.String `tfsdk:"icon" ddField:"Icon"`
	Website      types.String `tfsdk:"website" ddField:"Website"`
	WebsiteFound types.String `tfsdk:"website_found" ddField:"WebsiteFound"`
	Id           types.String `tfsdk:"id" ddField:"Id"`
}

type technologyDefectdojoResource struct {
	dd.AppAnalysis
}

func technologyToRequest(obj dd.AppAnalysis) dd.AppAnalysisRequest {
	return dd.AppAnalysisRequest{
		Name:         obj.Name,
		Product:      obj.Product,
		User:         obj.User,
		Confidence:   obj.Confidence,
		Version:      obj.Version,
		Icon:         obj.Icon,
		Website:      obj.Website,
		WebsiteFound: obj.WebsiteFound,
	}
}

func (ddr *technologyDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := technologyToRequest(ddr.AppAnalysis)
	apiResp, err := client.TechnologiesCreateWithResponse(ctx, reqBody)
	if apiResp.JSON201 != nil {
		ddr.AppAnalysis = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *technologyDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.TechnologiesRetrieveWithResponse(ctx, idNumber, &dd.TechnologiesRetrieveParams{})
	if apiResp.JSON200 != nil {
		ddr.AppAnalysis = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *technologyDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := technologyToRequest(ddr.AppAnalysis)
	apiResp, err := client.TechnologiesUpdateWithResponse(ctx, idNumber, reqBody)
	if apiResp.JSON200 != nil {
		ddr.AppAnalysis = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *technologyDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.TechnologiesDestroyWithResponse(ctx, idNumber)
	return apiResp.StatusCode(), apiResp.Body, err
}

type technologyResource struct {
	terraformResource
}

var _ resource.Resource = &technologyResource{}
var _ resource.ResourceWithImportState = &technologyResource{}

func NewTechnologyResource() resource.Resource {
	return &technologyResource{
		terraformResource: terraformResource{dataProvider: technologyDataProvider{}},
	}
}

func (r technologyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_technology"
}

type technologyDataProvider struct{}

func (r technologyDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data technologyResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *technologyResourceData) id() types.String { return d.Id }

func (d *technologyResourceData) defectdojoResource() defectdojoResource {
	return &technologyDefectdojoResource{AppAnalysis: dd.AppAnalysis{}}
}
