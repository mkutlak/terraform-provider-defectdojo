package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type urlDataSource struct {
	terraformDatasource
}

func (t urlDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for Defect Dojo URL. You can specify either the `id` or the `host` to look up the URL. " +
			"Note that if multiple URL records share the same host, lookup by `host` will fail; use `id` instead.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "The host of the URL, which can be a domain name or an IP address. Specify either id or host.",
				Optional:            true,
				Computed:            true,
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "The protocol of the URL (e.g., http, https, ftp, etc.)",
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "The port number of the URL (optional)",
				Computed:            true,
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "The path of the URL (optional)",
				Computed:            true,
			},
			"query": schema.StringAttribute{
				MarkdownDescription: "The query string of the URL (optional)",
				Computed:            true,
			},
			"fragment": schema.StringAttribute{
				MarkdownDescription: "The fragment identifier of the URL (optional)",
				Computed:            true,
			},
			"user_info": schema.StringAttribute{
				MarkdownDescription: "Connection details for a given user",
				Computed:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "Tags applied to the URL",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"host_validation_failure": schema.BoolAttribute{
				MarkdownDescription: "Dictates whether the endpoint was found to have host validation issues during creation",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of the URL",
				Computed:            true,
			},
		},
	}
}

func (d urlDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_url"
}

var _ datasource.DataSource = &urlDataSource{}

func NewUrlDataSource() datasource.DataSource {
	return &urlDataSource{
		terraformDatasource: terraformDatasource{
			dataProvider: urlDataProvider{},
		},
	}
}
