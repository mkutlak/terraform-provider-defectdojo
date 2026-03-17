package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type engagementDataSource struct {
	terraformDatasource
}

var _ datasource.DataSource = &engagementDataSource{}
var _ datasource.DataSourceWithConfigure = &engagementDataSource{}

func NewEngagementDataSource() datasource.DataSource {
	return &engagementDataSource{
		terraformDatasource: terraformDatasource{
			dataProvider: engagementDataProvider{},
		},
	}
}

func (d engagementDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_engagement"
}

func (d engagementDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Engagement data source",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Identifier",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Name of the Engagement",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Description of the Engagement",
			},
			"product": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "ID of the Product this Engagement belongs to",
			},
			"target_start": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Start date of the Engagement (format: 2006-01-02)",
			},
			"target_end": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "End date of the Engagement (format: 2006-01-02)",
			},
			"engagement_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Type of Engagement",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Status of the Engagement",
			},
			"lead": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "ID of the lead user",
			},
			"reason": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Reason for the Engagement",
			},
			"version": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Version of the product being tested",
			},
			"branch_tag": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Tag or branch of the product the engagement tested",
			},
			"commit_hash": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Commit hash from repo",
			},
			"build_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Build ID of the product the engagement tested",
			},
			"tracker": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Link to epic or ticket system",
			},
			"test_strategy": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Test strategy for the engagement",
			},
			"threat_model": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether a threat model was performed",
			},
			"api_test": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether an API test was performed",
			},
			"pen_test": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether a pen test was performed",
			},
			"check_list": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether a check list was used",
			},
			"deduplication_on_engagement": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether deduplication is on engagement level",
			},
			"first_contacted": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date first contacted (format: 2006-01-02)",
			},
			"source_code_management_uri": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource link to source code",
			},
			"preset": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "ID of the preset for this Engagement",
			},
			"report_type": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "ID of the report type",
			},
			"requester": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "ID of the requester",
			},
			"tags": schema.SetAttribute{
				Computed:            true,
				MarkdownDescription: "Tags for this Engagement",
				ElementType:         types.StringType,
			},
		},
	}
}
