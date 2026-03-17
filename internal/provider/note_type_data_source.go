package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type noteTypeDataSource struct {
	terraformDatasource
}

func (t noteTypeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for Defect Dojo Note Type",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Note Type",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the Note Type",
				Computed:            true,
			},
			"is_active": schema.BoolAttribute{
				MarkdownDescription: "Whether the Note Type is active",
				Computed:            true,
			},
			"is_mandatory": schema.BoolAttribute{
				MarkdownDescription: "Whether the Note Type is mandatory",
				Computed:            true,
			},
			"is_single": schema.BoolAttribute{
				MarkdownDescription: "Whether only a single note of this type is allowed",
				Computed:            true,
			},
		},
	}
}

func (d noteTypeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_note_type"
}

var _ datasource.DataSource = &noteTypeDataSource{}

func NewNoteTypeDataSource() datasource.DataSource {
	return &noteTypeDataSource{
		terraformDatasource: terraformDatasource{dataProvider: noteTypeDataProvider{}},
	}
}
