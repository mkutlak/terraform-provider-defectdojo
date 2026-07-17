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

type locationDataSource struct {
	terraformDatasource
}

func (t locationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for a DefectDojo Location. Locations are the read-only polymorphic base of URL and network locations in DefectDojo 3.x: the `/api/v2/location/` API only supports GET, so this data source cannot create or manage locations. Manage `defectdojo_url` or `defectdojo_network_location` resources instead; a location's id matches the id of the concrete resource it backs (e.g. a `defectdojo_url` resource's id is also its location id).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier. This is the only supported lookup key: the locations API has no name-based filter.",
				Required:            true,
			},
			"location_type": schema.StringAttribute{
				MarkdownDescription: "The type of location that is stored, e.g. \"URL\" for HTTP(S) locations managed via `defectdojo_url`, or a network location type. This is a polymorphic discriminator that DefectDojo manages automatically.",
				Computed:            true,
			},
			"location_value": schema.StringAttribute{
				MarkdownDescription: "The string representation of the location, automatically managed by DefectDojo.",
				Computed:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "Tags applied to the location.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"inherited_tags": schema.SetAttribute{
				MarkdownDescription: "Tags inherited by the location from related objects.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"created": schema.StringAttribute{
				MarkdownDescription: "Time that the object was initially created, and saved to the database.",
				Computed:            true,
			},
			"updated": schema.StringAttribute{
				MarkdownDescription: "Time that the object was most recently saved to the database.",
				Computed:            true,
			},
		},
	}
}

func (d locationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_location"
}

var _ datasource.DataSource = &locationDataSource{}

func NewLocationDataSource() datasource.DataSource {
	return &locationDataSource{
		terraformDatasource: terraformDatasource{dataProvider: locationDataProvider{}},
	}
}

type locationResourceData struct {
	Id            types.String `tfsdk:"id" ddField:"Id"`
	LocationType  types.String `tfsdk:"location_type" ddField:"LocationType"`
	LocationValue types.String `tfsdk:"location_value" ddField:"LocationValue"`
	Tags          types.Set    `tfsdk:"tags" ddField:"Tags"`
	InheritedTags types.Set    `tfsdk:"inherited_tags" ddField:"InheritedTags"`
	Created       types.String `tfsdk:"created" ddField:"Created"`
	Updated       types.String `tfsdk:"updated" ddField:"Updated"`
}

// locationDefectdojoResource wraps the Location model, the polymorphic
// read-only base of URL and network locations served by the DefectDojo 3.x
// locations API.
type locationDefectdojoResource struct {
	dd.Location
}

// errLocationReadOnly is returned by the write stubs below. The locations API
// only supports GET; the methods only exist to satisfy the defectdojoResource
// interface for the data source read path.
var errLocationReadOnly = errors.New("locations are read-only in DefectDojo 3.x; manage defectdojo_url or defectdojo_network_location resources instead")

func (ddr *locationDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	return 0, nil, errLocationReadOnly
}

func (ddr *locationDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.LocationRetrieveWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.Location = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *locationDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	return 0, nil, errLocationReadOnly
}

func (ddr *locationDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	return 0, nil, errLocationReadOnly
}

type locationDataProvider struct{}

func (r locationDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data locationResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *locationResourceData) id() types.String { return d.Id }

func (d *locationResourceData) setId(v types.String) { d.Id = v }

func (d *locationResourceData) defectdojoResource() defectdojoResource {
	return &locationDefectdojoResource{Location: dd.Location{}}
}
