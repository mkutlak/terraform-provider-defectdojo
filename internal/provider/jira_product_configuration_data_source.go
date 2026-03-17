package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type jiraProductConfigurationDataSource struct {
	terraformDatasource
}

func (t jiraProductConfigurationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for a DefectDojo Jira Product Configuration",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"project_key": schema.StringAttribute{
				MarkdownDescription: "The Jira Project Key",
				Computed:            true,
			},
			"issue_template_dir": schema.StringAttribute{
				MarkdownDescription: "The folder containing Django templates used to render the JIRA issue description.",
				Computed:            true,
			},
			"push_all_issues": schema.BoolAttribute{
				MarkdownDescription: "Automatically maintain parity with JIRA.",
				Computed:            true,
			},
			"enable_engagement_epic_mapping": schema.BoolAttribute{
				MarkdownDescription: "Whether to map engagements to epics in Jira",
				Computed:            true,
			},
			"push_notes": schema.BoolAttribute{
				MarkdownDescription: "Whether to push notes to Jira",
				Computed:            true,
			},
			"product_jira_sla_notification": schema.BoolAttribute{
				MarkdownDescription: "Send SLA notifications as comments",
				Computed:            true,
			},
			"risk_acceptance_expiration_notification": schema.BoolAttribute{
				MarkdownDescription: "Send Risk Acceptance expiration notifications as comments",
				Computed:            true,
			},
			"jira_instance_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the Jira Instance",
				Computed:            true,
			},
			"product_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the Product",
				Computed:            true,
			},
			"engagement_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the Engagement",
				Computed:            true,
			},
			"add_vulnerability_id_to_jira_label": schema.BoolAttribute{
				MarkdownDescription: "Whether to add the vulnerability ID to the Jira label.",
				Computed:            true,
			},
			"component": schema.StringAttribute{
				MarkdownDescription: "The Jira component to use for issues.",
				Computed:            true,
			},
			"default_assignee": schema.StringAttribute{
				MarkdownDescription: "JIRA default assignee (name).",
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the Jira integration is enabled.",
				Computed:            true,
			},
			"epic_issue_type_name": schema.StringAttribute{
				MarkdownDescription: "The name of the structure that represents an Epic.",
				Computed:            true,
			},
			"jira_labels": schema.StringAttribute{
				MarkdownDescription: "JIRA issue labels space separated.",
				Computed:            true,
			},
		},
	}
}

func (d jiraProductConfigurationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_product_configuration"
}

var _ datasource.DataSource = &jiraProductConfigurationDataSource{}

func NewJiraProductConfigurationDataSource() datasource.DataSource {
	return &jiraProductConfigurationDataSource{
		terraformDatasource: terraformDatasource{
			dataProvider: jiraProductConfigurationDataProvider{},
		},
	}
}
