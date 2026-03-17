package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type findingTemplateDataSource struct {
	terraformDatasource
}

func (t findingTemplateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for DefectDojo Finding Template",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "The title of the Finding Template",
				Computed:            true,
			},
			"severity": schema.StringAttribute{
				MarkdownDescription: "The severity of the finding",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Detailed description of the vulnerability",
				Computed:            true,
			},
			"mitigation": schema.StringAttribute{
				MarkdownDescription: "Steps to mitigate the vulnerability",
				Computed:            true,
			},
			"impact": schema.StringAttribute{
				MarkdownDescription: "Impact of the vulnerability",
				Computed:            true,
			},
			"references": schema.StringAttribute{
				MarkdownDescription: "References for the vulnerability",
				Computed:            true,
			},
			"cwe": schema.Int64Attribute{
				MarkdownDescription: "CWE number",
				Computed:            true,
			},
			"cvssv3": schema.StringAttribute{
				MarkdownDescription: "CVSSv3 score string",
				Computed:            true,
			},
			"cvssv3_score": schema.Float64Attribute{
				MarkdownDescription: "CVSSv3 numeric score",
				Computed:            true,
			},
			"cvssv4": schema.StringAttribute{
				MarkdownDescription: "CVSSv4 score string",
				Computed:            true,
			},
			"cvssv4_score": schema.Float64Attribute{
				MarkdownDescription: "CVSSv4 numeric score",
				Computed:            true,
			},
			"component_name": schema.StringAttribute{
				MarkdownDescription: "Affected component name",
				Computed:            true,
			},
			"component_version": schema.StringAttribute{
				MarkdownDescription: "Affected component version",
				Computed:            true,
			},
			"fix_available": schema.BoolAttribute{
				MarkdownDescription: "Indicates if a fix is available",
				Computed:            true,
			},
			"fix_version": schema.StringAttribute{
				MarkdownDescription: "Version where fix is available",
				Computed:            true,
			},
			"planned_remediation_version": schema.StringAttribute{
				MarkdownDescription: "Target version for remediation",
				Computed:            true,
			},
			"effort_for_fixing": schema.StringAttribute{
				MarkdownDescription: "Effort estimate for fixing",
				Computed:            true,
			},
			"severity_justification": schema.StringAttribute{
				MarkdownDescription: "Explanation of why this severity level is appropriate",
				Computed:            true,
			},
			"steps_to_reproduce": schema.StringAttribute{
				MarkdownDescription: "Standard reproduction steps for this vulnerability type",
				Computed:            true,
			},
			"endpoints_text": schema.StringAttribute{
				MarkdownDescription: "Endpoint URLs (one per line)",
				Computed:            true,
			},
			"notes": schema.StringAttribute{
				MarkdownDescription: "Note content to add when applying this template",
				Computed:            true,
			},
		},
	}
}

func (d findingTemplateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_finding_template"
}

var _ datasource.DataSource = &findingTemplateDataSource{}

func NewFindingTemplateDataSource() datasource.DataSource {
	return &findingTemplateDataSource{
		terraformDatasource: terraformDatasource{dataProvider: findingTemplateDataProvider{}},
	}
}
