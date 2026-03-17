package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type engagementPresetDataSource struct {
	terraformDatasource
}

var _ datasource.DataSource = &engagementPresetDataSource{}
var _ datasource.DataSourceWithConfigure = &engagementPresetDataSource{}

func NewEngagementPresetDataSource() datasource.DataSource {
	return &engagementPresetDataSource{
		terraformDatasource: terraformDatasource{
			dataProvider: engagementPresetDataProvider{},
		},
	}
}

func (d engagementPresetDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_engagement_preset"
}

func (d engagementPresetDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Engagement Preset data source",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Identifier",
			},
			"title": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Brief description of preset",
			},
			"product": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "ID of the Product this Preset belongs to",
			},
			"notes": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Description of what needs to be tested",
			},
			"scope": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Scope of Engagement testing",
			},
			"network_locations": schema.SetAttribute{
				Computed:            true,
				MarkdownDescription: "IDs of network locations",
				ElementType:         types.Int64Type,
			},
			"test_type": schema.SetAttribute{
				Computed:            true,
				MarkdownDescription: "IDs of test types",
				ElementType:         types.Int64Type,
			},
		},
	}
}
