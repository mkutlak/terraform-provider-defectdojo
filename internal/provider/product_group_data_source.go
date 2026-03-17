package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type productGroupDataSource struct {
	terraformDatasource
}

func (t productGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for Defect Dojo Product Group",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"product": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Product",
				Computed:            true,
			},
			"group": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Dojo Group",
				Computed:            true,
			},
			"role": schema.Int64Attribute{
				MarkdownDescription: "The role ID for this product group membership",
				Computed:            true,
			},
		},
	}
}

func (d productGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_product_group"
}

var _ datasource.DataSource = &productGroupDataSource{}

func NewProductGroupDataSource() datasource.DataSource {
	return &productGroupDataSource{
		terraformDatasource: terraformDatasource{
			dataProvider: productGroupDataProvider{},
		},
	}
}
