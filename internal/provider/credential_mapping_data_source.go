package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type credentialMappingDataSource struct {
	terraformDatasource
}

func (t credentialMappingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for DefectDojo Credential Mapping",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"cred_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Credential.",
				Computed:            true,
			},
			"product": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Product.",
				Computed:            true,
			},
			"engagement": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Engagement.",
				Computed:            true,
			},
			"finding": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Finding.",
				Computed:            true,
			},
			"test": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Test.",
				Computed:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "URL associated with the credential mapping.",
				Computed:            true,
			},
			"is_authn_provider": schema.BoolAttribute{
				MarkdownDescription: "Whether this is an authentication provider.",
				Computed:            true,
			},
		},
	}
}

func (d credentialMappingDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_mapping"
}

var _ datasource.DataSource = &credentialMappingDataSource{}

func NewCredentialMappingDataSource() datasource.DataSource {
	return &credentialMappingDataSource{
		terraformDatasource: terraformDatasource{
			dataProvider: credentialMappingDataProvider{},
		},
	}
}
