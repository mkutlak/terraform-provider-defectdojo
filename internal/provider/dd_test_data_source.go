package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ddTestDataSource struct {
	terraformDatasource
}

var _ datasource.DataSource = &ddTestDataSource{}
var _ datasource.DataSourceWithConfigure = &ddTestDataSource{}

func NewDDTestDataSource() datasource.DataSource {
	return &ddTestDataSource{
		terraformDatasource: terraformDatasource{
			dataProvider: ddTestDataProvider{},
		},
	}
}

func (d ddTestDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_test"
}

func (d ddTestDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Test data source",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Identifier",
			},
			"test_type": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "ID of the Test Type",
			},
			"engagement": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "ID of the Engagement",
			},
			"target_start": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Start datetime of the Test (RFC3339 format)",
			},
			"target_end": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "End datetime of the Test (RFC3339 format)",
			},
			"title": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Title of the Test",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Description of the Test",
			},
			"version": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Version tested",
			},
			"branch_tag": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Tag or branch that was tested",
			},
			"commit_hash": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Commit hash tested",
			},
			"build_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Build ID that was tested",
			},
			"percent_complete": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Percentage of test completion",
			},
			"environment": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "ID of the environment",
			},
			"lead": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "ID of the lead user",
			},
			"scan_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Type of scan",
			},
			"tags": schema.SetAttribute{
				Computed:            true,
				MarkdownDescription: "Tags for this Test",
				ElementType:         types.StringType,
			},
		},
	}
}
