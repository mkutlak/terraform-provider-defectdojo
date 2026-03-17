package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type languageDataSource struct {
	terraformDatasource
}

func (t languageDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for DefectDojo Language",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"language_type": schema.Int64Attribute{
				MarkdownDescription: "The language type ID.",
				Computed:            true,
			},
			"product": schema.Int64Attribute{
				MarkdownDescription: "The product this language is associated with.",
				Computed:            true,
			},
			"user": schema.Int64Attribute{
				MarkdownDescription: "The user who added this language.",
				Computed:            true,
			},
			"files": schema.Int64Attribute{
				MarkdownDescription: "Number of files.",
				Computed:            true,
			},
			"code": schema.Int64Attribute{
				MarkdownDescription: "Number of lines of code.",
				Computed:            true,
			},
			"blank": schema.Int64Attribute{
				MarkdownDescription: "Number of blank lines.",
				Computed:            true,
			},
			"comment": schema.Int64Attribute{
				MarkdownDescription: "Number of comment lines.",
				Computed:            true,
			},
		},
	}
}

func (d languageDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_language"
}

var _ datasource.DataSource = &languageDataSource{}

func NewLanguageDataSource() datasource.DataSource {
	return &languageDataSource{
		terraformDatasource: terraformDatasource{dataProvider: languageDataProvider{}},
	}
}
