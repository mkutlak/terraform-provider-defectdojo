package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type toolTypeDataSource struct {
	terraformDatasource
}

func (t toolTypeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for Defect Dojo Tool Type. You can specify either the `id` or the `name` to look up the Tool Type. **Deprecated:** DefectDojo removes the `tool_type` API in version 3.5.0, on 2026-11-01. The server adds `X-Deprecated` and `X-End-Of-Life-Date` headers to each response.",
		DeprecationMessage:  "DefectDojo deprecated the tool_type API in version 3.2.0. DefectDojo removes the API in version 3.5.0, on 2026-11-01. This data source then fails.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Tool Type. Specify either id or name.",
				Optional:            true,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the Tool Type",
				Computed:            true,
			},
		},
	}
}

func (d toolTypeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tool_type"
}

var _ datasource.DataSource = &toolTypeDataSource{}

func NewToolTypeDataSource() datasource.DataSource {
	return &toolTypeDataSource{
		terraformDatasource: terraformDatasource{dataProvider: toolTypeDataProvider{}},
	}
}
