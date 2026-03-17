package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type assetGroupDataSource struct {
	terraformDatasource
}

func (t assetGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for DefectDojo Asset Group",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"asset": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Asset.",
				Computed:            true,
			},
			"group": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Group.",
				Computed:            true,
			},
			"role": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Role.",
				Computed:            true,
			},
		},
	}
}

func (d assetGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_asset_group"
}

var _ datasource.DataSource = &assetGroupDataSource{}

func NewAssetGroupDataSource() datasource.DataSource {
	return &assetGroupDataSource{
		terraformDatasource: terraformDatasource{
			dataProvider: assetGroupDataProvider{},
		},
	}
}
