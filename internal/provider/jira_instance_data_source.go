package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type jiraInstanceDataSource struct {
	terraformDatasource
}

func (t jiraInstanceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for DefectDojo Jira Instance",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "The URL of the Jira instance.",
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username or Email Address for Jira authentication.",
				Computed:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password or API Token for Jira authentication.",
				Computed:            true,
				Sensitive:           true,
			},
			"configuration_name": schema.StringAttribute{
				MarkdownDescription: "A name for this Jira configuration.",
				Computed:            true,
			},
			"epic_name_id": schema.Int64Attribute{
				MarkdownDescription: "The Epic Name field ID in Jira.",
				Computed:            true,
			},
			"open_status_key": schema.Int64Attribute{
				MarkdownDescription: "Transition ID to Re-Open JIRA issues.",
				Computed:            true,
			},
			"close_status_key": schema.Int64Attribute{
				MarkdownDescription: "Transition ID to Close JIRA issues.",
				Computed:            true,
			},
			"info_mapping_severity": schema.StringAttribute{
				MarkdownDescription: "Maps to the Priority field in Jira for Info severity.",
				Computed:            true,
			},
			"low_mapping_severity": schema.StringAttribute{
				MarkdownDescription: "Maps to the Priority field in Jira for Low severity.",
				Computed:            true,
			},
			"medium_mapping_severity": schema.StringAttribute{
				MarkdownDescription: "Maps to the Priority field in Jira for Medium severity.",
				Computed:            true,
			},
			"high_mapping_severity": schema.StringAttribute{
				MarkdownDescription: "Maps to the Priority field in Jira for High severity.",
				Computed:            true,
			},
			"critical_mapping_severity": schema.StringAttribute{
				MarkdownDescription: "Maps to the Priority field in Jira for Critical severity.",
				Computed:            true,
			},
			"finding_text": schema.StringAttribute{
				MarkdownDescription: "Additional text added to findings in Jira.",
				Computed:            true,
			},
			"accepted_mapping_resolution": schema.StringAttribute{
				MarkdownDescription: "JIRA resolution that maps to Risk Accepted in DefectDojo.",
				Computed:            true,
			},
			"false_positive_mapping_resolution": schema.StringAttribute{
				MarkdownDescription: "JIRA resolution that maps to False Positive in DefectDojo.",
				Computed:            true,
			},
			"global_jira_sla_notification": schema.BoolAttribute{
				MarkdownDescription: "Send SLA notifications via Jira.",
				Computed:            true,
			},
			"finding_jira_sync": schema.BoolAttribute{
				MarkdownDescription: "If enabled, sync finding changes automatically to JIRA.",
				Computed:            true,
			},
			"issue_template_dir": schema.StringAttribute{
				MarkdownDescription: "Folder containing Django templates for JIRA issue descriptions.",
				Computed:            true,
			},
			"default_issue_type": schema.StringAttribute{
				MarkdownDescription: "Default issue type.",
				Computed:            true,
			},
		},
	}
}

func (d jiraInstanceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_instance"
}

var _ datasource.DataSource = &jiraInstanceDataSource{}

func NewJiraInstanceDataSource() datasource.DataSource {
	return &jiraInstanceDataSource{
		terraformDatasource: terraformDatasource{
			dataProvider: jiraInstanceDataProvider{},
		},
	}
}
