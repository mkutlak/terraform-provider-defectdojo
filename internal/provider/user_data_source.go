package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type userDataSource struct {
	terraformDatasource
}

func (t userDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for Defect Dojo User",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The username of the User",
				Computed:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "The email address of the User",
				Computed:            true,
			},
			"first_name": schema.StringAttribute{
				MarkdownDescription: "The first name of the User",
				Computed:            true,
			},
			"last_name": schema.StringAttribute{
				MarkdownDescription: "The last name of the User",
				Computed:            true,
			},
			"is_active": schema.BoolAttribute{
				MarkdownDescription: "Whether this user account is active",
				Computed:            true,
			},
			"is_superuser": schema.BoolAttribute{
				MarkdownDescription: "Whether this user has all permissions without explicitly assigning them",
				Computed:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password for the User (not returned by API)",
				Computed:            true,
				Sensitive:           true,
			},
		},
	}
}

func (d userDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

var _ datasource.DataSource = &userDataSource{}

func NewUserDataSource() datasource.DataSource {
	return &userDataSource{
		terraformDatasource: terraformDatasource{
			dataProvider: userDataProvider{},
		},
	}
}
