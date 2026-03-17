package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type regulationDataSource struct {
	terraformDatasource
}

func (t regulationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for Defect Dojo Regulation",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Regulation",
				Computed:            true,
			},
			"acronym": schema.StringAttribute{
				MarkdownDescription: "A shortened representation of the name",
				Computed:            true,
			},
			"category": schema.StringAttribute{
				MarkdownDescription: "The subject of the regulation",
				Computed:            true,
			},
			"jurisdiction": schema.StringAttribute{
				MarkdownDescription: "The territory over which the regulation applies",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Information about the regulation's purpose",
				Computed:            true,
			},
			"reference": schema.StringAttribute{
				MarkdownDescription: "An external URL for more information",
				Computed:            true,
			},
		},
	}
}

func (d regulationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_regulation"
}

var _ datasource.DataSource = &regulationDataSource{}

func NewRegulationDataSource() datasource.DataSource {
	return &regulationDataSource{
		terraformDatasource: terraformDatasource{dataProvider: regulationDataProvider{}},
	}
}
