package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
)

func (r engagementResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Engagement",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the Engagement",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the Engagement",
				Optional:            true,
			},
			"product": schema.Int64Attribute{
				MarkdownDescription: "ID of the Product this Engagement belongs to",
				Required:            true,
			},
			"target_start": schema.StringAttribute{
				MarkdownDescription: "Start date of the Engagement (format: 2006-01-02)",
				Required:            true,
			},
			"target_end": schema.StringAttribute{
				MarkdownDescription: "End date of the Engagement (format: 2006-01-02)",
				Required:            true,
			},
			"engagement_type": schema.StringAttribute{
				MarkdownDescription: "Type of Engagement: 'Interactive' or 'CI/CD'",
				Optional:            true,
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Status of the Engagement (Not Started, Blocked, Cancelled, Completed, In Progress, On Hold, Waiting for Resource)",
				Optional:            true,
				Computed:            true,
			},
			"lead": schema.Int64Attribute{
				MarkdownDescription: "ID of the lead user for this Engagement",
				Optional:            true,
			},
			"reason": schema.StringAttribute{
				MarkdownDescription: "Reason for the Engagement",
				Optional:            true,
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "Version of the product being tested",
				Optional:            true,
			},
			"branch_tag": schema.StringAttribute{
				MarkdownDescription: "Tag or branch of the product the engagement tested",
				Optional:            true,
			},
			"commit_hash": schema.StringAttribute{
				MarkdownDescription: "Commit hash from repo",
				Optional:            true,
			},
			"build_id": schema.StringAttribute{
				MarkdownDescription: "Build ID of the product the engagement tested",
				Optional:            true,
			},
			"tracker": schema.StringAttribute{
				MarkdownDescription: "Link to epic or ticket system with changes to version",
				Optional:            true,
			},
			"test_strategy": schema.StringAttribute{
				MarkdownDescription: "Test strategy for the engagement",
				Optional:            true,
			},
			"threat_model": schema.BoolAttribute{
				MarkdownDescription: "Whether a threat model was performed",
				Optional:            true,
				Computed:            true,
			},
			"api_test": schema.BoolAttribute{
				MarkdownDescription: "Whether an API test was performed",
				Optional:            true,
				Computed:            true,
			},
			"pen_test": schema.BoolAttribute{
				MarkdownDescription: "Whether a pen test was performed",
				Optional:            true,
				Computed:            true,
			},
			"check_list": schema.BoolAttribute{
				MarkdownDescription: "Whether a check list was used",
				Optional:            true,
				Computed:            true,
			},
			"deduplication_on_engagement": schema.BoolAttribute{
				MarkdownDescription: "If enabled deduplication will only mark a finding in this engagement as duplicate of another finding if both findings are in this engagement",
				Optional:            true,
				Computed:            true,
			},
			"first_contacted": schema.StringAttribute{
				MarkdownDescription: "Date first contacted (format: 2006-01-02)",
				Optional:            true,
			},
			"source_code_management_uri": schema.StringAttribute{
				MarkdownDescription: "Resource link to source code",
				Optional:            true,
			},
			"preset": schema.Int64Attribute{
				MarkdownDescription: "ID of the preset for this Engagement",
				Optional:            true,
			},
			"report_type": schema.Int64Attribute{
				MarkdownDescription: "ID of the report type",
				Optional:            true,
			},
			"requester": schema.Int64Attribute{
				MarkdownDescription: "ID of the requester",
				Optional:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "Tags for this Engagement",
				Optional:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

type engagementResourceData struct {
	Id                        types.String `tfsdk:"id" ddField:"Id"`
	Name                      types.String `tfsdk:"name" ddField:"Name"`
	Description               types.String `tfsdk:"description" ddField:"Description"`
	Product                   types.Int64  `tfsdk:"product" ddField:"Product"`
	TargetStart               types.String `tfsdk:"target_start" ddField:"TargetStart"`
	TargetEnd                 types.String `tfsdk:"target_end" ddField:"TargetEnd"`
	EngagementType            types.String `tfsdk:"engagement_type" ddField:"EngagementType"`
	Status                    types.String `tfsdk:"status" ddField:"Status"`
	Lead                      types.Int64  `tfsdk:"lead" ddField:"Lead"`
	Reason                    types.String `tfsdk:"reason" ddField:"Reason"`
	Version                   types.String `tfsdk:"version" ddField:"Version"`
	BranchTag                 types.String `tfsdk:"branch_tag" ddField:"BranchTag"`
	CommitHash                types.String `tfsdk:"commit_hash" ddField:"CommitHash"`
	BuildId                   types.String `tfsdk:"build_id" ddField:"BuildId"`
	Tracker                   types.String `tfsdk:"tracker" ddField:"Tracker"`
	TestStrategy              types.String `tfsdk:"test_strategy" ddField:"TestStrategy"`
	ThreatModel               types.Bool   `tfsdk:"threat_model" ddField:"ThreatModel"`
	ApiTest                   types.Bool   `tfsdk:"api_test" ddField:"ApiTest"`
	PenTest                   types.Bool   `tfsdk:"pen_test" ddField:"PenTest"`
	CheckList                 types.Bool   `tfsdk:"check_list" ddField:"CheckList"`
	DeduplicationOnEngagement types.Bool   `tfsdk:"deduplication_on_engagement" ddField:"DeduplicationOnEngagement"`
	FirstContacted            types.String `tfsdk:"first_contacted" ddField:"FirstContacted"`
	SourceCodeManagementUri   types.String `tfsdk:"source_code_management_uri" ddField:"SourceCodeManagementUri"`
	Preset                    types.Int64  `tfsdk:"preset" ddField:"Preset"`
	ReportType                types.Int64  `tfsdk:"report_type" ddField:"ReportType"`
	Requester                 types.Int64  `tfsdk:"requester" ddField:"Requester"`
	Tags                      types.Set    `tfsdk:"tags" ddField:"Tags"`
}

type engagementDefectdojoResource struct {
	dd.Engagement
}

func engagementToRequest(e dd.Engagement) dd.EngagementRequest {
	req := dd.EngagementRequest{
		Name:                      e.Name,
		Description:               e.Description,
		Product:                   e.Product,
		TargetStart:               e.TargetStart,
		TargetEnd:                 e.TargetEnd,
		Lead:                      e.Lead,
		Reason:                    e.Reason,
		Version:                   e.Version,
		BranchTag:                 e.BranchTag,
		CommitHash:                e.CommitHash,
		BuildId:                   e.BuildId,
		Tracker:                   e.Tracker,
		TestStrategy:              e.TestStrategy,
		ThreatModel:               e.ThreatModel,
		ApiTest:                   e.ApiTest,
		PenTest:                   e.PenTest,
		CheckList:                 e.CheckList,
		DeduplicationOnEngagement: e.DeduplicationOnEngagement,
		FirstContacted:            e.FirstContacted,
		SourceCodeManagementUri:   e.SourceCodeManagementUri,
		Preset:                    e.Preset,
		ReportType:                e.ReportType,
		Requester:                 e.Requester,
		Tags:                      e.Tags,
	}
	if e.EngagementType != nil && *e.EngagementType != "" {
		v := dd.EngagementRequestEngagementType(*e.EngagementType)
		req.EngagementType = &v
	}
	if e.Status != nil && *e.Status != "" {
		v := dd.EngagementRequestStatus(*e.Status)
		req.Status = &v
	}
	return req
}

func (ddr *engagementDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	tflog.Info(ctx, "engagementDefectdojoResource createApiCall")
	reqBody := engagementToRequest(ddr.Engagement)
	apiResp, err := client.EngagementsCreateWithResponse(ctx, reqBody)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON201 != nil {
		ddr.Engagement = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *engagementDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "engagementDefectdojoResource readApiCall")
	apiResp, err := client.EngagementsRetrieveWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON200 != nil {
		ddr.Engagement = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *engagementDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "engagementDefectdojoResource updateApiCall")
	reqBody := engagementToRequest(ddr.Engagement)
	apiResp, err := client.EngagementsUpdateWithResponse(ctx, idNumber, reqBody)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON200 != nil {
		ddr.Engagement = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *engagementDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "engagementDefectdojoResource deleteApiCall")
	apiResp, err := client.EngagementsDestroyWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	return apiResp.StatusCode(), apiResp.Body, nil
}

type engagementResource struct {
	terraformResource
}

var _ resource.Resource = &engagementResource{}
var _ resource.ResourceWithImportState = &engagementResource{}

func NewEngagementResource() resource.Resource {
	return &engagementResource{
		terraformResource: terraformResource{
			dataProvider: engagementDataProvider{},
		},
	}
}

func (r engagementResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_engagement"
}

type engagementDataProvider struct{}

func (r engagementDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data engagementResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *engagementResourceData) id() types.String {
	return d.Id
}

func (d *engagementResourceData) defectdojoResource() defectdojoResource {
	return &engagementDefectdojoResource{
		Engagement: dd.Engagement{},
	}
}
