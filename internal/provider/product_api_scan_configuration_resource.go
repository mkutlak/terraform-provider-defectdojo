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

func (t productAPIScanConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Product API Scan Configuration",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"product": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Product.",
				Required:            true,
			},
			"tool_configuration": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Tool Configuration.",
				Required:            true,
			},
			"service_key_1": schema.StringAttribute{
				MarkdownDescription: "Service key 1 for the API scan configuration.",
				Optional:            true,
				Computed:            true,
			},
			"service_key_2": schema.StringAttribute{
				MarkdownDescription: "Service key 2 for the API scan configuration.",
				Optional:            true,
				Computed:            true,
			},
			"service_key_3": schema.StringAttribute{
				MarkdownDescription: "Service key 3 for the API scan configuration.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

type productAPIScanConfigurationResourceData struct {
	Id                types.String `tfsdk:"id" ddField:"Id"`
	Product           types.Int64  `tfsdk:"product" ddField:"Product"`
	ToolConfiguration types.Int64  `tfsdk:"tool_configuration" ddField:"ToolConfiguration"`
	ServiceKey1       types.String `tfsdk:"service_key_1" ddField:"ServiceKey1"`
	ServiceKey2       types.String `tfsdk:"service_key_2" ddField:"ServiceKey2"`
	ServiceKey3       types.String `tfsdk:"service_key_3" ddField:"ServiceKey3"`
}

type productAPIScanConfigurationDefectdojoResource struct {
	dd.ProductAPIScanConfiguration
}

func productAPIScanConfigurationToRequest(p dd.ProductAPIScanConfiguration) dd.ProductAPIScanConfigurationRequest {
	return dd.ProductAPIScanConfigurationRequest{
		Product:           p.Product,
		ToolConfiguration: p.ToolConfiguration,
		ServiceKey1:       p.ServiceKey1,
		ServiceKey2:       p.ServiceKey2,
		ServiceKey3:       p.ServiceKey3,
	}
}

func (ddr *productAPIScanConfigurationDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := productAPIScanConfigurationToRequest(ddr.ProductAPIScanConfiguration)
	apiResp, err := client.ProductApiScanConfigurationsCreateWithResponse(ctx, reqBody)
	if apiResp.JSON201 != nil {
		ddr.ProductAPIScanConfiguration = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *productAPIScanConfigurationDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.ProductApiScanConfigurationsRetrieveWithResponse(ctx, idNumber, &dd.ProductApiScanConfigurationsRetrieveParams{})
	if apiResp.JSON200 != nil {
		ddr.ProductAPIScanConfiguration = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *productAPIScanConfigurationDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := productAPIScanConfigurationToRequest(ddr.ProductAPIScanConfiguration)
	apiResp, err := client.ProductApiScanConfigurationsUpdateWithResponse(ctx, idNumber, reqBody)
	if apiResp.JSON200 != nil {
		ddr.ProductAPIScanConfiguration = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *productAPIScanConfigurationDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.ProductApiScanConfigurationsDestroyWithResponse(ctx, idNumber)
	return apiResp.StatusCode(), apiResp.Body, err
}

func (d *productAPIScanConfigurationResourceData) id() types.String {
	return d.Id
}

func (d *productAPIScanConfigurationResourceData) defectdojoResource() defectdojoResource {
	return &productAPIScanConfigurationDefectdojoResource{ProductAPIScanConfiguration: dd.ProductAPIScanConfiguration{}}
}

type productAPIScanConfigurationResource struct {
	terraformResource
}

var _ resource.Resource = &productAPIScanConfigurationResource{}
var _ resource.ResourceWithImportState = &productAPIScanConfigurationResource{}

func NewProductAPIScanConfigurationResource() resource.Resource {
	return &productAPIScanConfigurationResource{
		terraformResource: terraformResource{
			dataProvider: productAPIScanConfigurationDataProvider{},
		},
	}
}

func (r productAPIScanConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_product_api_scan_configuration"
}

type productAPIScanConfigurationDataProvider struct{}

func (r productAPIScanConfigurationDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data productAPIScanConfigurationResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}
