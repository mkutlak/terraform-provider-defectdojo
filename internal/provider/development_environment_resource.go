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

func (t developmentEnvironmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Development Environment",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Development Environment",
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

type developmentEnvironmentResourceData struct {
	Name types.String `tfsdk:"name" ddField:"Name"`
	Id   types.String `tfsdk:"id" ddField:"Id"`
}

type developmentEnvironmentDefectdojoResource struct {
	dd.DevelopmentEnvironment
}

func developmentEnvironmentToRequest(obj dd.DevelopmentEnvironment) dd.DevelopmentEnvironmentRequest {
	return dd.DevelopmentEnvironmentRequest{
		Name: obj.Name,
	}
}

func (ddr *developmentEnvironmentDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := developmentEnvironmentToRequest(ddr.DevelopmentEnvironment)
	apiResp, err := client.DevelopmentEnvironmentsCreateWithResponse(ctx, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON201 != nil {
		ddr.DevelopmentEnvironment = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *developmentEnvironmentDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.DevelopmentEnvironmentsRetrieveWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.DevelopmentEnvironment = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *developmentEnvironmentDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := developmentEnvironmentToRequest(ddr.DevelopmentEnvironment)
	apiResp, err := client.DevelopmentEnvironmentsUpdateWithResponse(ctx, idNumber, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.DevelopmentEnvironment = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *developmentEnvironmentDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.DevelopmentEnvironmentsDestroyWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

type developmentEnvironmentResource struct {
	terraformResource
}

var _ resource.Resource = &developmentEnvironmentResource{}
var _ resource.ResourceWithImportState = &developmentEnvironmentResource{}

func NewDevelopmentEnvironmentResource() resource.Resource {
	return &developmentEnvironmentResource{
		terraformResource: terraformResource{dataProvider: developmentEnvironmentDataProvider{}},
	}
}

func (r developmentEnvironmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_development_environment"
}

type developmentEnvironmentDataProvider struct{}

func (r developmentEnvironmentDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data developmentEnvironmentResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *developmentEnvironmentResourceData) id() types.String { return d.Id }

func (d *developmentEnvironmentResourceData) defectdojoResource() defectdojoResource {
	return &developmentEnvironmentDefectdojoResource{DevelopmentEnvironment: dd.DevelopmentEnvironment{}}
}
