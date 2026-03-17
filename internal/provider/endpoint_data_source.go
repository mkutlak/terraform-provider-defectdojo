package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type endpointDataSource struct {
	terraformDatasource
}

func (t endpointDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for DefectDojo Endpoint",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "The communication protocol/scheme such as 'http', 'ftp', 'dns', etc.",
				Computed:            true,
			},
			"userinfo": schema.StringAttribute{
				MarkdownDescription: "User info as 'alice', 'bob', etc.",
				Computed:            true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "The host name or IP address.",
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "The network port associated with the endpoint.",
				Computed:            true,
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "The location of the resource.",
				Computed:            true,
			},
			"query": schema.StringAttribute{
				MarkdownDescription: "The query string.",
				Computed:            true,
			},
			"fragment": schema.StringAttribute{
				MarkdownDescription: "The fragment identifier.",
				Computed:            true,
			},
			"product": schema.Int64Attribute{
				MarkdownDescription: "The product this endpoint belongs to.",
				Computed:            true,
			},
		},
	}
}

func (d endpointDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_endpoint"
}

var _ datasource.DataSource = &endpointDataSource{}

func NewEndpointDataSource() datasource.DataSource {
	return &endpointDataSource{
		terraformDatasource: terraformDatasource{dataProvider: endpointDataProvider{}},
	}
}
