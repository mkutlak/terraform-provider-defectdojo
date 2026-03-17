package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type dojoGroupMemberDataSource struct {
	terraformDatasource
}

func (t dojoGroupMemberDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for Defect Dojo Dojo Group Member",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"group": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Dojo Group",
				Computed:            true,
			},
			"user": schema.Int64Attribute{
				MarkdownDescription: "The ID of the User",
				Computed:            true,
			},
			"role": schema.Int64Attribute{
				MarkdownDescription: "The role determining the permissions of the user to manage the group",
				Computed:            true,
			},
		},
	}
}

func (d dojoGroupMemberDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dojo_group_member"
}

var _ datasource.DataSource = &dojoGroupMemberDataSource{}

func NewDojoGroupMemberDataSource() datasource.DataSource {
	return &dojoGroupMemberDataSource{
		terraformDatasource: terraformDatasource{
			dataProvider: dojoGroupMemberDataProvider{},
		},
	}
}
