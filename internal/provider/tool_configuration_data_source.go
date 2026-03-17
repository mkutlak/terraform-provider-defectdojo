package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type toolConfigurationDataSource struct {
	terraformDatasource
}

func (t toolConfigurationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for Defect Dojo Tool Configuration",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Tool Configuration",
				Computed:            true,
			},
			"tool_type": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Tool Type",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the Tool Configuration",
				Computed:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "The URL of the tool",
				Computed:            true,
			},
			"authentication_type": schema.StringAttribute{
				MarkdownDescription: "Authentication type",
				Computed:            true,
			},
			"auth_title": schema.StringAttribute{
				MarkdownDescription: "Title for authentication credentials",
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username for authentication",
				Computed:            true,
			},
			"extras": schema.StringAttribute{
				MarkdownDescription: "Additional definitions that will be consumed by scanner",
				Computed:            true,
			},
		},
	}
}

func (d toolConfigurationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tool_configuration"
}

var _ datasource.DataSource = &toolConfigurationDataSource{}

func NewToolConfigurationDataSource() datasource.DataSource {
	return &toolConfigurationDataSource{
		terraformDatasource: terraformDatasource{dataProvider: toolConfigurationDataProvider{}},
	}
}
