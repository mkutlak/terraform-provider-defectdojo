package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type productMemberDataSource struct {
	terraformDatasource
}

func (t productMemberDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for Defect Dojo Product Member",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"product": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Product",
				Computed:            true,
			},
			"user": schema.Int64Attribute{
				MarkdownDescription: "The ID of the User",
				Computed:            true,
			},
			"role": schema.Int64Attribute{
				MarkdownDescription: "The role ID for this product membership",
				Computed:            true,
			},
		},
	}
}

func (d productMemberDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_product_member"
}

var _ datasource.DataSource = &productMemberDataSource{}

func NewProductMemberDataSource() datasource.DataSource {
	return &productMemberDataSource{
		terraformDatasource: terraformDatasource{
			dataProvider: productMemberDataProvider{},
		},
	}
}
