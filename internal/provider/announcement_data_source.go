package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type announcementDataSource struct {
	terraformDatasource
}

func (t announcementDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for DefectDojo Announcement",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"message": schema.StringAttribute{
				MarkdownDescription: "The announcement message.",
				Computed:            true,
			},
			"style": schema.StringAttribute{
				MarkdownDescription: "The style of banner to display.",
				Computed:            true,
			},
			"dismissable": schema.BoolAttribute{
				MarkdownDescription: "Whether users can dismiss the announcement.",
				Computed:            true,
			},
		},
	}
}

func (d announcementDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_announcement"
}

var _ datasource.DataSource = &announcementDataSource{}

func NewAnnouncementDataSource() datasource.DataSource {
	return &announcementDataSource{
		terraformDatasource: terraformDatasource{
			dataProvider: announcementDataProvider{},
		},
	}
}
