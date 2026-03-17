package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type productTypeGroupDataSource struct {
	terraformDatasource
}

func (t productTypeGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for Defect Dojo Product Type Group",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"product_type": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Product Type",
				Computed:            true,
			},
			"group": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Dojo Group",
				Computed:            true,
			},
			"role": schema.Int64Attribute{
				MarkdownDescription: "The role ID for this product type group membership",
				Computed:            true,
			},
		},
	}
}

func (d productTypeGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_product_type_group"
}

var _ datasource.DataSource = &productTypeGroupDataSource{}

func NewProductTypeGroupDataSource() datasource.DataSource {
	return &productTypeGroupDataSource{
		terraformDatasource: terraformDatasource{
			dataProvider: productTypeGroupDataProvider{},
		},
	}
}
