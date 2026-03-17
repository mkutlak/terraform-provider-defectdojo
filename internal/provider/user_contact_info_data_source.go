package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type userContactInfoDataSource struct {
	terraformDatasource
}

func (t userContactInfoDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for Defect Dojo User Contact Info",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"user": schema.Int64Attribute{
				MarkdownDescription: "The ID of the User this contact info belongs to",
				Computed:            true,
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "Title of the user",
				Computed:            true,
			},
			"phone_number": schema.StringAttribute{
				MarkdownDescription: "Phone number in format '+999999999'. Up to 15 digits allowed.",
				Computed:            true,
			},
			"cell_number": schema.StringAttribute{
				MarkdownDescription: "Cell number in format '+999999999'. Up to 15 digits allowed.",
				Computed:            true,
			},
			"twitter_username": schema.StringAttribute{
				MarkdownDescription: "Twitter username",
				Computed:            true,
			},
			"github_username": schema.StringAttribute{
				MarkdownDescription: "GitHub username",
				Computed:            true,
			},
			"slack_username": schema.StringAttribute{
				MarkdownDescription: "Email address associated with your Slack account",
				Computed:            true,
			},
			"slack_user_id": schema.StringAttribute{
				MarkdownDescription: "Slack user ID",
				Computed:            true,
			},
			"block_execution": schema.BoolAttribute{
				MarkdownDescription: "Instead of async deduping a finding the findings will be deduped synchronously",
				Computed:            true,
			},
			"force_password_reset": schema.BoolAttribute{
				MarkdownDescription: "Forces this user to reset their password on next login",
				Computed:            true,
			},
		},
	}
}

func (d userContactInfoDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_contact_info"
}

var _ datasource.DataSource = &userContactInfoDataSource{}

func NewUserContactInfoDataSource() datasource.DataSource {
	return &userContactInfoDataSource{
		terraformDatasource: terraformDatasource{
			dataProvider: userContactInfoDataProvider{},
		},
	}
}
