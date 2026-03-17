package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type networkLocationDataSource struct {
	terraformDatasource
}

func (t networkLocationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for Defect Dojo Network Location",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"location": schema.StringAttribute{
				MarkdownDescription: "Location of network testing: Examples: VPN, Internet or Internal",
				Computed:            true,
			},
		},
	}
}

func (d networkLocationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_location"
}

var _ datasource.DataSource = &networkLocationDataSource{}

func NewNetworkLocationDataSource() datasource.DataSource {
	return &networkLocationDataSource{
		terraformDatasource: terraformDatasource{dataProvider: networkLocationDataProvider{}},
	}
}
