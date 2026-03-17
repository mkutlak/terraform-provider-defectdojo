package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type credentialDataSource struct {
	terraformDatasource
}

func (t credentialDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for Defect Dojo Credential",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Credential",
				Computed:            true,
			},
			"environment": schema.Int64Attribute{
				MarkdownDescription: "The ID of the environment for this credential",
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The username for the credential",
				Computed:            true,
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "The role for the credential",
				Computed:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "The URL for the credential",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the credential",
				Computed:            true,
			},
			"authentication": schema.StringAttribute{
				MarkdownDescription: "Authentication type: 'Form' or 'SSO'",
				Computed:            true,
			},
			"http_authentication": schema.StringAttribute{
				MarkdownDescription: "HTTP authentication type: 'Basic' or 'NTLM'",
				Computed:            true,
			},
			"login_regex": schema.StringAttribute{
				MarkdownDescription: "Login regex pattern",
				Computed:            true,
			},
			"logout_regex": schema.StringAttribute{
				MarkdownDescription: "Logout regex pattern",
				Computed:            true,
			},
			"is_valid": schema.BoolAttribute{
				MarkdownDescription: "Whether the credential is valid",
				Computed:            true,
			},
		},
	}
}

func (d credentialDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential"
}

var _ datasource.DataSource = &credentialDataSource{}

func NewCredentialDataSource() datasource.DataSource {
	return &credentialDataSource{
		terraformDatasource: terraformDatasource{
			dataProvider: credentialDataProvider{},
		},
	}
}
