package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type endpointStatusDataSource struct {
	terraformDatasource
}

func (t endpointStatusDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for DefectDojo Endpoint Status",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"endpoint": schema.Int64Attribute{
				MarkdownDescription: "The endpoint this status is associated with.",
				Computed:            true,
			},
			"finding": schema.Int64Attribute{
				MarkdownDescription: "The finding this status is associated with.",
				Computed:            true,
			},
			"date": schema.StringAttribute{
				MarkdownDescription: "Date of the endpoint status (format: 2006-01-02).",
				Computed:            true,
			},
			"false_positive": schema.BoolAttribute{
				MarkdownDescription: "Whether the finding is a false positive.",
				Computed:            true,
			},
			"mitigated": schema.BoolAttribute{
				MarkdownDescription: "Whether the finding has been mitigated.",
				Computed:            true,
			},
			"mitigated_by": schema.Int64Attribute{
				MarkdownDescription: "The user who mitigated the finding.",
				Computed:            true,
			},
			"out_of_scope": schema.BoolAttribute{
				MarkdownDescription: "Whether the finding is out of scope.",
				Computed:            true,
			},
			"risk_accepted": schema.BoolAttribute{
				MarkdownDescription: "Whether the risk has been accepted.",
				Computed:            true,
			},
		},
	}
}

func (d endpointStatusDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_endpoint_status"
}

var _ datasource.DataSource = &endpointStatusDataSource{}

func NewEndpointStatusDataSource() datasource.DataSource {
	return &endpointStatusDataSource{
		terraformDatasource: terraformDatasource{dataProvider: endpointStatusDataProvider{}},
	}
}
