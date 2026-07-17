package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type metadataDataSource struct {
	terraformDatasource
}

func (t metadataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for DefectDojo Metadata. You can specify either the `id` or the `name` to look up the Metadata. Note that names are only unique per parent object, so looking up by name when more than one parent object has metadata with that name will fail; use `id` instead in that case.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the metadata field. Specify either id or name.",
				Optional:            true,
				Computed:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "The value of the metadata field.",
				Computed:            true,
			},
			"product": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Product this metadata is attached to.",
				Computed:            true,
			},
			"finding": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Finding this metadata is attached to.",
				Computed:            true,
			},
		},
	}
}

func (d metadataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_metadata"
}

var _ datasource.DataSource = &metadataDataSource{}

func NewMetadataDataSource() datasource.DataSource {
	return &metadataDataSource{
		terraformDatasource: terraformDatasource{
			dataProvider: metadataDataProvider{},
		},
	}
}
