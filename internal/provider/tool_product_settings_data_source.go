package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type toolProductSettingsDataSource struct {
	terraformDatasource
}

func (t toolProductSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for DefectDojo Tool Product Settings",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the tool product settings.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the tool product settings.",
				Computed:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "URL for the tool product settings.",
				Computed:            true,
			},
			"setting_url": schema.StringAttribute{
				MarkdownDescription: "The settings URL for the tool.",
				Computed:            true,
			},
			"product": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Product.",
				Computed:            true,
			},
			"tool_configuration": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Tool Configuration.",
				Computed:            true,
			},
			"tool_project_id": schema.StringAttribute{
				MarkdownDescription: "The project ID in the tool.",
				Computed:            true,
			},
		},
	}
}

func (d toolProductSettingsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tool_product_settings"
}

var _ datasource.DataSource = &toolProductSettingsDataSource{}

func NewToolProductSettingsDataSource() datasource.DataSource {
	return &toolProductSettingsDataSource{
		terraformDatasource: terraformDatasource{
			dataProvider: toolProductSettingsDataProvider{},
		},
	}
}
