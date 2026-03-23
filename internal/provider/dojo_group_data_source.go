package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type dojoGroupDataSource struct {
	terraformDatasource
}

func (t dojoGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for Defect Dojo Dojo Group. You can specify either the `id` or the `name` to look up the Dojo Group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Dojo Group. Specify either id or name.",
				Optional:            true,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the Dojo Group",
				Computed:            true,
			},
		},
	}
}

func (d dojoGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dojo_group"
}

var _ datasource.DataSource = &dojoGroupDataSource{}

func NewDojoGroupDataSource() datasource.DataSource {
	return &dojoGroupDataSource{
		terraformDatasource: terraformDatasource{
			dataProvider: dojoGroupDataProvider{},
		},
	}
}
