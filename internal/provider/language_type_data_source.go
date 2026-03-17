package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type languageTypeDataSource struct {
	terraformDatasource
}

func (t languageTypeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for Defect Dojo Language Type",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"language": schema.StringAttribute{
				MarkdownDescription: "The name of the language",
				Computed:            true,
			},
			"color": schema.StringAttribute{
				MarkdownDescription: "Color associated with the language",
				Computed:            true,
			},
		},
	}
}

func (d languageTypeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_language_type"
}

var _ datasource.DataSource = &languageTypeDataSource{}

func NewLanguageTypeDataSource() datasource.DataSource {
	return &languageTypeDataSource{
		terraformDatasource: terraformDatasource{dataProvider: languageTypeDataProvider{}},
	}
}
