package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type notificationWebhookDataSource struct {
	terraformDatasource
}

func (t notificationWebhookDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for DefectDojo Notification Webhook",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the incoming webhook.",
				Computed:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "The full URL of the incoming webhook.",
				Computed:            true,
			},
			"header_name": schema.StringAttribute{
				MarkdownDescription: "Name of the header required for interacting with the webhook endpoint.",
				Computed:            true,
			},
			"header_value": schema.StringAttribute{
				MarkdownDescription: "Content of the header required for interacting with the webhook endpoint.",
				Computed:            true,
				Sensitive:           true,
			},
			"owner": schema.Int64Attribute{
				MarkdownDescription: "Owner/receiver of notification.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Status of the incoming webhook.",
				Computed:            true,
			},
		},
	}
}

func (d notificationWebhookDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_webhook"
}

var _ datasource.DataSource = &notificationWebhookDataSource{}

func NewNotificationWebhookDataSource() datasource.DataSource {
	return &notificationWebhookDataSource{
		terraformDatasource: terraformDatasource{
			dataProvider: notificationWebhookDataProvider{},
		},
	}
}
