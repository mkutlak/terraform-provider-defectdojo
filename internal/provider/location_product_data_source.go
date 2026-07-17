package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type locationProductDataSource struct {
	terraformDatasource
}

func (t locationProductDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for DefectDojo Location Product. Links a Location (e.g. a URL or Network Location) to a Product. Looked up by id only.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Required:            true,
			},
			"location": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Location.",
				Computed:            true,
			},
			"product": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Product.",
				Computed:            true,
			},
			"relationship": schema.StringAttribute{
				MarkdownDescription: "The relationship between the location and the product. Valid values: 'owned_by', 'used_by'.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The status of the given Location. Valid values: 'Active', 'Mitigated'.",
				Computed:            true,
			},
			"location_type": schema.StringAttribute{
				MarkdownDescription: "The type of the linked Location.",
				Computed:            true,
			},
			"location_value": schema.StringAttribute{
				MarkdownDescription: "The value representation of the linked Location.",
				Computed:            true,
			},
		},
	}
}

func (d locationProductDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_location_product"
}

var _ datasource.DataSource = &locationProductDataSource{}

func NewLocationProductDataSource() datasource.DataSource {
	return &locationProductDataSource{
		terraformDatasource: terraformDatasource{
			dataProvider: locationProductDataProvider{},
		},
	}
}
