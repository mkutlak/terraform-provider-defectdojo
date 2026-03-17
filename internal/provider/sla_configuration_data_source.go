package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type slaConfigurationDataSource struct {
	terraformDatasource
}

func (t slaConfigurationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for Defect Dojo SLA Configuration",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "A unique name for the set of SLAs",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the SLA Configuration",
				Computed:            true,
			},
			"critical": schema.Int64Attribute{
				MarkdownDescription: "The number of days to remediate a critical finding",
				Computed:            true,
			},
			"high": schema.Int64Attribute{
				MarkdownDescription: "The number of days to remediate a high finding",
				Computed:            true,
			},
			"medium": schema.Int64Attribute{
				MarkdownDescription: "The number of days to remediate a medium finding",
				Computed:            true,
			},
			"low": schema.Int64Attribute{
				MarkdownDescription: "The number of days to remediate a low finding",
				Computed:            true,
			},
			"enforce_critical": schema.BoolAttribute{
				MarkdownDescription: "When enabled, critical findings will be assigned an SLA expiration date",
				Computed:            true,
			},
			"enforce_high": schema.BoolAttribute{
				MarkdownDescription: "When enabled, high findings will be assigned an SLA expiration date",
				Computed:            true,
			},
			"enforce_medium": schema.BoolAttribute{
				MarkdownDescription: "When enabled, medium findings will be assigned an SLA expiration date",
				Computed:            true,
			},
			"enforce_low": schema.BoolAttribute{
				MarkdownDescription: "When enabled, low findings will be assigned an SLA expiration date",
				Computed:            true,
			},
			"restart_sla_on_reactivation": schema.BoolAttribute{
				MarkdownDescription: "When enabled, findings that were previously mitigated but are reactivated during reimport will have their SLA period restarted",
				Computed:            true,
			},
		},
	}
}

func (d slaConfigurationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sla_configuration"
}

var _ datasource.DataSource = &slaConfigurationDataSource{}

func NewSlaConfigurationDataSource() datasource.DataSource {
	return &slaConfigurationDataSource{
		terraformDatasource: terraformDatasource{dataProvider: slaConfigurationDataProvider{}},
	}
}
