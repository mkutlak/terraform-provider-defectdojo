package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type cicdInfrastructureDataSource struct {
	terraformDatasource
}

func (t cicdInfrastructureDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for DefectDojo CI/CD Infrastructure. You can specify either the `id` or the `name` to look up the CI/CD infrastructure. A name is unique together with `infrastructure_type`, so looking up by name when more than one infrastructure_type shares that name will fail; use `id` instead in that case. DefectDojo hides `description` and `url` from a token without the `view_cicdinfrastructure` or `change_cicdinfrastructure` permission; both attributes then read as empty.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the CI/CD infrastructure. Specify either id or name.",
				Optional:            true,
				Computed:            true,
			},
			"infrastructure_type": schema.StringAttribute{
				MarkdownDescription: "The kind of CI/CD infrastructure. Valid values: `scm_server`, `build_server`, `orchestration`.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the CI/CD infrastructure. DefectDojo hides this from a token without the `view_cicdinfrastructure` or `change_cicdinfrastructure` permission.",
				Computed:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "The public URL of the CI/CD infrastructure. DefectDojo hides this from a token without the `view_cicdinfrastructure` or `change_cicdinfrastructure` permission.",
				Computed:            true,
			},
		},
	}
}

func (d cicdInfrastructureDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cicd_infrastructure"
}

var _ datasource.DataSource = &cicdInfrastructureDataSource{}

func NewCicdInfrastructureDataSource() datasource.DataSource {
	return &cicdInfrastructureDataSource{
		terraformDatasource: terraformDatasource{dataProvider: cicdInfrastructureDataProvider{}},
	}
}
