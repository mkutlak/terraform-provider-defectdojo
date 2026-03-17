package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type technologyDataSource struct {
	terraformDatasource
}

func (t technologyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for DefectDojo Technology (App Analysis)",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the technology",
				Computed:            true,
			},
			"product": schema.Int64Attribute{
				MarkdownDescription: "The product this technology is associated with.",
				Computed:            true,
			},
			"user": schema.Int64Attribute{
				MarkdownDescription: "The user who added this technology.",
				Computed:            true,
			},
			"confidence": schema.Int64Attribute{
				MarkdownDescription: "Confidence level of the detection.",
				Computed:            true,
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "Version of the technology.",
				Computed:            true,
			},
			"icon": schema.StringAttribute{
				MarkdownDescription: "Icon for the technology.",
				Computed:            true,
			},
			"website": schema.StringAttribute{
				MarkdownDescription: "Website of the technology.",
				Computed:            true,
			},
			"website_found": schema.StringAttribute{
				MarkdownDescription: "Website where the technology was found.",
				Computed:            true,
			},
		},
	}
}

func (d technologyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_technology"
}

var _ datasource.DataSource = &technologyDataSource{}

func NewTechnologyDataSource() datasource.DataSource {
	return &technologyDataSource{
		terraformDatasource: terraformDatasource{dataProvider: technologyDataProvider{}},
	}
}
