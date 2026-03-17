package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type riskAcceptanceDataSource struct {
	terraformDatasource
}

func (t riskAcceptanceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for DefectDojo Risk Acceptance",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Descriptive name for the risk acceptance.",
				Computed:            true,
			},
			"owner": schema.Int64Attribute{
				MarkdownDescription: "User ID of the owner of this risk acceptance.",
				Computed:            true,
			},
			"accepted_findings": schema.SetAttribute{
				MarkdownDescription: "IDs of findings included in this risk acceptance.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"accepted_by": schema.StringAttribute{
				MarkdownDescription: "The person that accepts the risk.",
				Computed:            true,
			},
			"expiration_date": schema.StringAttribute{
				MarkdownDescription: "When the risk acceptance expires.",
				Computed:            true,
			},
			"decision": schema.StringAttribute{
				MarkdownDescription: "Risk treatment decision.",
				Computed:            true,
			},
			"decision_details": schema.StringAttribute{
				MarkdownDescription: "Details about the risk treatment decision.",
				Computed:            true,
			},
			"recommendation": schema.StringAttribute{
				MarkdownDescription: "Security team recommendation.",
				Computed:            true,
			},
			"recommendation_details": schema.StringAttribute{
				MarkdownDescription: "Details explaining the security recommendation.",
				Computed:            true,
			},
			"reactivate_expired": schema.BoolAttribute{
				MarkdownDescription: "Reactivate findings when risk acceptance expires.",
				Computed:            true,
			},
			"restart_sla_expired": schema.BoolAttribute{
				MarkdownDescription: "Restart SLA for findings when risk acceptance expires.",
				Computed:            true,
			},
		},
	}
}

func (d riskAcceptanceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_risk_acceptance"
}

var _ datasource.DataSource = &riskAcceptanceDataSource{}

func NewRiskAcceptanceDataSource() datasource.DataSource {
	return &riskAcceptanceDataSource{
		terraformDatasource: terraformDatasource{
			dataProvider: riskAcceptanceDataProvider{},
		},
	}
}
