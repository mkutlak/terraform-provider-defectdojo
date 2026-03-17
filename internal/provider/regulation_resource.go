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

func (t regulationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Regulation",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Regulation",
				Required:            true,
			},
			"acronym": schema.StringAttribute{
				MarkdownDescription: "A shortened representation of the name",
				Required:            true,
			},
			"category": schema.StringAttribute{
				MarkdownDescription: "The subject of the regulation. Valid values: 'privacy', 'finance', 'education', 'medical', 'corporate', 'security', 'government', 'other'",
				Required:            true,
			},
			"jurisdiction": schema.StringAttribute{
				MarkdownDescription: "The territory over which the regulation applies",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Information about the regulation's purpose",
				Optional:            true,
				Computed:            true,
			},
			"reference": schema.StringAttribute{
				MarkdownDescription: "An external URL for more information",
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

type regulationResourceData struct {
	Name         types.String `tfsdk:"name" ddField:"Name"`
	Acronym      types.String `tfsdk:"acronym" ddField:"Acronym"`
	Category     types.String `tfsdk:"category" ddField:"Category"`
	Jurisdiction types.String `tfsdk:"jurisdiction" ddField:"Jurisdiction"`
	Description  types.String `tfsdk:"description" ddField:"Description"`
	Reference    types.String `tfsdk:"reference" ddField:"Reference"`
	Id           types.String `tfsdk:"id" ddField:"Id"`
}

type regulationDefectdojoResource struct {
	dd.Regulation
}

func regulationToRequest(obj dd.Regulation) dd.RegulationRequest {
	return dd.RegulationRequest{
		Name:         obj.Name,
		Acronym:      obj.Acronym,
		Category:     dd.RegulationRequestCategory(obj.Category),
		Jurisdiction: obj.Jurisdiction,
		Description:  obj.Description,
		Reference:    obj.Reference,
	}
}

func (ddr *regulationDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := regulationToRequest(ddr.Regulation)
	apiResp, err := client.RegulationsCreateWithResponse(ctx, reqBody)
	if apiResp.JSON201 != nil {
		ddr.Regulation = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *regulationDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.RegulationsRetrieveWithResponse(ctx, idNumber)
	if apiResp.JSON200 != nil {
		ddr.Regulation = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *regulationDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := regulationToRequest(ddr.Regulation)
	apiResp, err := client.RegulationsUpdateWithResponse(ctx, idNumber, reqBody)
	if apiResp.JSON200 != nil {
		ddr.Regulation = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *regulationDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.RegulationsDestroyWithResponse(ctx, idNumber)
	return apiResp.StatusCode(), apiResp.Body, err
}

type regulationResource struct {
	terraformResource
}

var _ resource.Resource = &regulationResource{}
var _ resource.ResourceWithImportState = &regulationResource{}

func NewRegulationResource() resource.Resource {
	return &regulationResource{
		terraformResource: terraformResource{dataProvider: regulationDataProvider{}},
	}
}

func (r regulationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_regulation"
}

type regulationDataProvider struct{}

func (r regulationDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data regulationResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *regulationResourceData) id() types.String { return d.Id }

func (d *regulationResourceData) defectdojoResource() defectdojoResource {
	return &regulationDefectdojoResource{Regulation: dd.Regulation{}}
}
