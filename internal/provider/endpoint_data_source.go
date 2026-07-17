package provider

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
)

type endpointDataSource struct {
	terraformDatasource
}

func (t endpointDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for DefectDojo Endpoint. Since DefectDojo 3.x endpoints are a read-only projection of Locations: this data source reads the legacy endpoint representation, but endpoints can no longer be created or modified via the API (manage Locations instead). **Deprecated:** use the `defectdojo_url` resource/data source and the `defectdojo_location` / `defectdojo_location_product` data sources instead; this data source will be removed in the next major version.",
		DeprecationMessage:  "Endpoints are a read-only projection in DefectDojo 3.x. Use the defectdojo_url resource/data source and the defectdojo_location / defectdojo_location_product data sources instead. This data source will be removed in the next major version.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "The communication protocol/scheme such as 'http', 'ftp', 'dns', etc.",
				Computed:            true,
			},
			"userinfo": schema.StringAttribute{
				MarkdownDescription: "User info as 'alice', 'bob', etc.",
				Computed:            true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "The host name or IP address.",
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "The network port associated with the endpoint.",
				Computed:            true,
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "The location of the resource.",
				Computed:            true,
			},
			"query": schema.StringAttribute{
				MarkdownDescription: "The query string.",
				Computed:            true,
			},
			"fragment": schema.StringAttribute{
				MarkdownDescription: "The fragment identifier.",
				Computed:            true,
			},
			"product": schema.Int64Attribute{
				MarkdownDescription: "The product this endpoint belongs to.",
				Computed:            true,
			},
		},
	}
}

func (d endpointDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_endpoint"
}

var _ datasource.DataSource = &endpointDataSource{}

func NewEndpointDataSource() datasource.DataSource {
	return &endpointDataSource{
		terraformDatasource: terraformDatasource{dataProvider: endpointDataProvider{}},
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

// endpointDefectdojoResource wraps the V3EndpointCompatible model, which is the
// read-only endpoint projection served by the DefectDojo 3.x endpoints API.
type endpointDefectdojoResource struct {
	dd.V3EndpointCompatible
}

// errEndpointReadOnly is returned by the write stubs below. The endpoints API
// lost its write operations in DefectDojo 3.0; the methods only exist to
// satisfy the defectdojoResource interface for the data source read path.
var errEndpointReadOnly = errors.New("endpoints are a read-only projection in DefectDojo 3.x; manage Locations instead")

func (ddr *endpointDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	return 0, nil, errEndpointReadOnly
}

func (ddr *endpointDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.EndpointsRetrieveWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.V3EndpointCompatible = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *endpointDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	return 0, nil, errEndpointReadOnly
}

func (ddr *endpointDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	return 0, nil, errEndpointReadOnly
}

type endpointDataProvider struct{}

func (r endpointDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data endpointResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *endpointResourceData) id() types.String { return d.Id }

func (d *endpointResourceData) setId(v types.String) { d.Id = v }

func (d *endpointResourceData) defectdojoResource() defectdojoResource {
	return &endpointDefectdojoResource{V3EndpointCompatible: dd.V3EndpointCompatible{}}
}
