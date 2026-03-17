package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type developmentEnvironmentDataSource struct {
	terraformDatasource
}

func (t developmentEnvironmentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for Defect Dojo Development Environment",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Development Environment",
				Computed:            true,
			},
		},
	}
}

func (d developmentEnvironmentDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_development_environment"
}

var _ datasource.DataSource = &developmentEnvironmentDataSource{}

func NewDevelopmentEnvironmentDataSource() datasource.DataSource {
	return &developmentEnvironmentDataSource{
		terraformDatasource: terraformDatasource{dataProvider: developmentEnvironmentDataProvider{}},
	}
}
