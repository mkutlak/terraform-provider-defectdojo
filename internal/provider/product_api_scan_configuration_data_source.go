package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type productAPIScanConfigurationDataSource struct {
	terraformDatasource
}

func (t productAPIScanConfigurationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for DefectDojo Product API Scan Configuration. **Deprecated:** DefectDojo removes the `product_api_scan_configuration` API in version 3.5.0, on 2026-11-01. The server adds `X-Deprecated` and `X-End-Of-Life-Date` headers to each response.",
		DeprecationMessage:  "DefectDojo deprecated the product_api_scan_configuration API in version 3.2.0. DefectDojo removes the API in version 3.5.0, on 2026-11-01. This data source then fails.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"product": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Product.",
				Computed:            true,
			},
			"tool_configuration": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Tool Configuration.",
				Computed:            true,
			},
			"service_key_1": schema.StringAttribute{
				MarkdownDescription: "Service key 1 for the API scan configuration.",
				Computed:            true,
			},
			"service_key_2": schema.StringAttribute{
				MarkdownDescription: "Service key 2 for the API scan configuration.",
				Computed:            true,
			},
			"service_key_3": schema.StringAttribute{
				MarkdownDescription: "Service key 3 for the API scan configuration.",
				Computed:            true,
			},
		},
	}
}

func (d productAPIScanConfigurationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_product_api_scan_configuration"
}

var _ datasource.DataSource = &productAPIScanConfigurationDataSource{}

func NewProductAPIScanConfigurationDataSource() datasource.DataSource {
	return &productAPIScanConfigurationDataSource{
		terraformDatasource: terraformDatasource{
			dataProvider: productAPIScanConfigurationDataProvider{},
		},
	}
}
