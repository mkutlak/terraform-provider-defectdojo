package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type notificationsDataSource struct {
	terraformDatasource
}

func (t notificationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	channelSetAttribute := func(description string) schema.SetAttribute {
		return schema.SetAttribute{
			MarkdownDescription: description,
			Computed:            true,
			ElementType:         types.StringType,
		}
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for DefectDojo Notifications: per-channel notification preferences for a scope (global default, per-user, or per-product). Looked up by id only.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Required:            true,
			},
			"product": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Product this notification preference row is scoped to.",
				Computed:            true,
			},
			"user": schema.Int64Attribute{
				MarkdownDescription: "The ID of the User this notification preference row is scoped to.",
				Computed:            true,
			},
			"template": schema.BoolAttribute{
				MarkdownDescription: "Whether this row is the notification template applied to new users/products.",
				Computed:            true,
			},
			"scan_added_empty": schema.StringAttribute{
				MarkdownDescription: "Triggered whenever an (re-)import has been done (even if that created/updated/closed no findings).",
				Computed:            true,
			},
			"auto_close_engagement":      channelSetAttribute("Notification channels for when an engagement is auto-closed."),
			"close_engagement":           channelSetAttribute("Notification channels for when an engagement is closed."),
			"code_review":                channelSetAttribute("Notification channels for when a code review is requested."),
			"engagement_added":           channelSetAttribute("Notification channels for when an engagement is added."),
			"jira_update":                channelSetAttribute("Notification channels for when a JIRA issue is updated."),
			"other":                      channelSetAttribute("Notification channels for other events."),
			"product_added":              channelSetAttribute("Notification channels for when a product is added."),
			"product_type_added":         channelSetAttribute("Notification channels for when a product type is added."),
			"review_requested":           channelSetAttribute("Notification channels for when a review is requested."),
			"risk_acceptance_expiration": channelSetAttribute("Notification channels for when a risk acceptance is about to expire."),
			"scan_added":                 channelSetAttribute("Notification channels for when a scan is added/imported."),
			"sla_breach":                 channelSetAttribute("Notification channels for when an SLA is breached."),
			"sla_breach_combined":        channelSetAttribute("Notification channels for combined SLA breach notifications."),
			"stale_engagement":           channelSetAttribute("Notification channels for stale engagements."),
			"test_added":                 channelSetAttribute("Notification channels for when a test is added."),
			"upcoming_engagement":        channelSetAttribute("Notification channels for upcoming engagements."),
			"user_mentioned":             channelSetAttribute("Notification channels for when a user is mentioned."),
		},
	}
}

func (d notificationsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notifications"
}

var _ datasource.DataSource = &notificationsDataSource{}

func NewNotificationsDataSource() datasource.DataSource {
	return &notificationsDataSource{
		terraformDatasource: terraformDatasource{
			dataProvider: notificationsDataProvider{},
		},
	}
}
