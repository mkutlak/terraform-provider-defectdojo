package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type productTypeMemberDataSource struct {
	terraformDatasource
}

func (t productTypeMemberDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for Defect Dojo Product Type Member",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"product_type": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Product Type",
				Computed:            true,
			},
			"user": schema.Int64Attribute{
				MarkdownDescription: "The ID of the User",
				Computed:            true,
			},
			"role": schema.Int64Attribute{
				MarkdownDescription: "The role ID for this product type membership",
				Computed:            true,
			},
		},
	}
}

func (d productTypeMemberDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_product_type_member"
}

var _ datasource.DataSource = &productTypeMemberDataSource{}

func NewProductTypeMemberDataSource() datasource.DataSource {
	return &productTypeMemberDataSource{
		terraformDatasource: terraformDatasource{
			dataProvider: productTypeMemberDataProvider{},
		},
	}
}
