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
		MarkdownDescription: "Data source for Defect Dojo Announcement. DefectDojo allows a single global announcement; if one already exists, import it: `terraform import defectdojo_announcement.example <id>`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Required:            true,
			},
			"message": schema.StringAttribute{
				MarkdownDescription: "This dismissable message will be displayed on all pages for authenticated users. It can contain basic html tags.",
				Computed:            true,
			},
			"style": schema.StringAttribute{
				MarkdownDescription: "The style of banner to display. (info, success, warning, danger)",
				Computed:            true,
			},
			"dismissable": schema.BoolAttribute{
				MarkdownDescription: "Ticking this box allows users to dismiss the current announcement",
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
		terraformDatasource: terraformDatasource{dataProvider: announcementDataProvider{}},
	}
}
