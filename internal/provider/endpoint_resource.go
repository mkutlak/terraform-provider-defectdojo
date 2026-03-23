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

func (t endpointResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Endpoint",
		Attributes: map[string]schema.Attribute{
			"protocol": schema.StringAttribute{
				MarkdownDescription: "The communication protocol/scheme such as 'http', 'ftp', 'dns', etc.",
				Optional:            true,
			},
			"userinfo": schema.StringAttribute{
				MarkdownDescription: "User info as 'alice', 'bob', etc.",
				Optional:            true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "The host name or IP address. It must not include the port number.",
				Optional:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "The network port associated with the endpoint.",
				Optional:            true,
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "The location of the resource, it must not start with a '/'.",
				Optional:            true,
			},
			"query": schema.StringAttribute{
				MarkdownDescription: "The query string, the question mark should be omitted.",
				Optional:            true,
			},
			"fragment": schema.StringAttribute{
				MarkdownDescription: "The fragment identifier which follows the hash mark. The hash mark should be omitted.",
				Optional:            true,
			},
			"product": schema.Int64Attribute{
				MarkdownDescription: "The product this endpoint belongs to.",
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

type endpointResourceData struct {
	Protocol types.String `tfsdk:"protocol" ddField:"Protocol"`
	Userinfo types.String `tfsdk:"userinfo" ddField:"Userinfo"`
	Host     types.String `tfsdk:"host" ddField:"Host"`
	Port     types.Int64  `tfsdk:"port" ddField:"Port"`
	Path     types.String `tfsdk:"path" ddField:"Path"`
	Query    types.String `tfsdk:"query" ddField:"Query"`
	Fragment types.String `tfsdk:"fragment" ddField:"Fragment"`
	Product  types.Int64  `tfsdk:"product" ddField:"Product"`
	Id       types.String `tfsdk:"id" ddField:"Id"`
}

type endpointDefectdojoResource struct {
	dd.Endpoint
}

func endpointToRequest(obj dd.Endpoint) dd.EndpointRequest {
	return dd.EndpointRequest{
		Protocol: obj.Protocol,
		Userinfo: obj.Userinfo,
		Host:     obj.Host,
		Port:     obj.Port,
		Path:     obj.Path,
		Query:    obj.Query,
		Fragment: obj.Fragment,
		Product:  obj.Product,
	}
}

func (ddr *endpointDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := endpointToRequest(ddr.Endpoint)
	apiResp, err := client.EndpointsCreateWithResponse(ctx, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON201 != nil {
		ddr.Endpoint = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *endpointDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.EndpointsRetrieveWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.Endpoint = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *endpointDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := endpointToRequest(ddr.Endpoint)
	apiResp, err := client.EndpointsUpdateWithResponse(ctx, idNumber, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.Endpoint = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *endpointDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.EndpointsDestroyWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

type endpointResource struct {
	terraformResource
}

var _ resource.Resource = &endpointResource{}
var _ resource.ResourceWithImportState = &endpointResource{}

func NewEndpointResource() resource.Resource {
	return &endpointResource{
		terraformResource: terraformResource{typeName: "defectdojo_endpoint", dataProvider: endpointDataProvider{}},
	}
}

func (r endpointResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_endpoint"
}

type endpointDataProvider struct{}

func (r endpointDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data endpointResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *endpointResourceData) id() types.String { return d.Id }

func (d *endpointResourceData) defectdojoResource() defectdojoResource {
	return &endpointDefectdojoResource{Endpoint: dd.Endpoint{}}
}
