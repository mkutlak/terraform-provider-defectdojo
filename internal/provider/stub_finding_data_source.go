package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type stubFindingDataSource struct {
	terraformDatasource
}

func (t stubFindingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for DefectDojo Stub Finding",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "The title of the Stub Finding",
				Computed:            true,
			},
			"test": schema.Int64Attribute{
				MarkdownDescription: "The test this stub finding belongs to.",
				Computed:            true,
			},
			"date": schema.StringAttribute{
				MarkdownDescription: "Date of the stub finding (format: 2006-01-02).",
				Computed:            true,
			},
			"severity": schema.StringAttribute{
				MarkdownDescription: "The severity of the stub finding.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the stub finding.",
				Computed:            true,
			},
			"reporter": schema.Int64Attribute{
				MarkdownDescription: "The user who reported this stub finding.",
				Computed:            true,
			},
		},
	}
}

func (d stubFindingDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stub_finding"
}

var _ datasource.DataSource = &stubFindingDataSource{}

func NewStubFindingDataSource() datasource.DataSource {
	return &stubFindingDataSource{
		terraformDatasource: terraformDatasource{dataProvider: stubFindingDataProvider{}},
	}
}
