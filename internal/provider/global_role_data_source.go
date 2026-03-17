package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type globalRoleDataSource struct {
	terraformDatasource
}

func (t globalRoleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for Defect Dojo Global Role",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"user": schema.Int64Attribute{
				MarkdownDescription: "The ID of the User",
				Computed:            true,
			},
			"group": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Dojo Group",
				Computed:            true,
			},
			"role": schema.Int64Attribute{
				MarkdownDescription: "The global role ID",
				Computed:            true,
			},
		},
	}
}

func (d globalRoleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_global_role"
}

var _ datasource.DataSource = &globalRoleDataSource{}

func NewGlobalRoleDataSource() datasource.DataSource {
	return &globalRoleDataSource{
		terraformDatasource: terraformDatasource{
			dataProvider: globalRoleDataProvider{},
		},
	}
}
