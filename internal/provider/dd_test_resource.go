package provider

import (
	"context"
	"fmt"

	dd "github.com/doximity/terraform-provider-defectdojo/internal/ddclient"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (r ddTestResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Test",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"test_type": schema.Int64Attribute{
				MarkdownDescription: "ID of the Test Type",
				Required:            true,
			},
			"engagement": schema.Int64Attribute{
				MarkdownDescription: "ID of the Engagement this Test belongs to",
				Required:            true,
			},
			"target_start": schema.StringAttribute{
				MarkdownDescription: "Start datetime of the Test (RFC3339 format, e.g. 2006-01-02T15:04:05Z)",
				Required:            true,
			},
			"target_end": schema.StringAttribute{
				MarkdownDescription: "End datetime of the Test (RFC3339 format, e.g. 2006-01-02T15:04:05Z)",
				Required:            true,
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "Title of the Test",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the Test",
				Optional:            true,
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "Version tested",
				Optional:            true,
			},
			"branch_tag": schema.StringAttribute{
				MarkdownDescription: "Tag or branch that was tested",
				Optional:            true,
			},
			"commit_hash": schema.StringAttribute{
				MarkdownDescription: "Commit hash tested",
				Optional:            true,
			},
			"build_id": schema.StringAttribute{
				MarkdownDescription: "Build ID that was tested",
				Optional:            true,
			},
			"percent_complete": schema.Int64Attribute{
				MarkdownDescription: "Percentage of test completion",
				Optional:            true,
			},
			"environment": schema.Int64Attribute{
				MarkdownDescription: "ID of the environment",
				Optional:            true,
			},
			"lead": schema.Int64Attribute{
				MarkdownDescription: "ID of the lead user",
				Optional:            true,
			},
			"scan_type": schema.StringAttribute{
				MarkdownDescription: "Type of scan",
				Optional:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "Tags for this Test",
				Optional:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

type ddTestResourceData struct {
	Id              types.String `tfsdk:"id" ddField:"Id"`
	TestType        types.Int64  `tfsdk:"test_type" ddField:"TestType"`
	Engagement      types.Int64  `tfsdk:"engagement" ddField:"Engagement"`
	TargetStart     types.String `tfsdk:"target_start" ddField:"TargetStart"`
	TargetEnd       types.String `tfsdk:"target_end" ddField:"TargetEnd"`
	Title           types.String `tfsdk:"title" ddField:"Title"`
	Description     types.String `tfsdk:"description" ddField:"Description"`
	Version         types.String `tfsdk:"version" ddField:"Version"`
	BranchTag       types.String `tfsdk:"branch_tag" ddField:"BranchTag"`
	CommitHash      types.String `tfsdk:"commit_hash" ddField:"CommitHash"`
	BuildId         types.String `tfsdk:"build_id" ddField:"BuildId"`
	PercentComplete types.Int64  `tfsdk:"percent_complete" ddField:"PercentComplete"`
	Environment     types.Int64  `tfsdk:"environment" ddField:"Environment"`
	Lead            types.Int64  `tfsdk:"lead" ddField:"Lead"`
	ScanType        types.String `tfsdk:"scan_type" ddField:"ScanType"`
	Tags            types.Set    `tfsdk:"tags" ddField:"Tags"`
}

// ddTestDefectdojoResource wraps the dd.TestCreate struct for create, and dd.Test for read/update.
// We use TestCreate for the embedded struct because it has Engagement as a required int (not *int).
type ddTestDefectdojoResource struct {
	dd.TestCreate
}

func ddTestToCreateRequest(t dd.TestCreate) dd.TestCreateRequest {
	return dd.TestCreateRequest{
		TestType:        t.TestType,
		Engagement:      t.Engagement,
		TargetStart:     t.TargetStart,
		TargetEnd:       t.TargetEnd,
		Title:           t.Title,
		Description:     t.Description,
		Version:         t.Version,
		BranchTag:       t.BranchTag,
		CommitHash:      t.CommitHash,
		BuildId:         t.BuildId,
		PercentComplete: t.PercentComplete,
		Environment:     t.Environment,
		Lead:            t.Lead,
		ScanType:        t.ScanType,
		Tags:            t.Tags,
	}
}

func ddTestToRequest(t dd.TestCreate) dd.TestRequest {
	return dd.TestRequest{
		TestType:        t.TestType,
		TargetStart:     t.TargetStart,
		TargetEnd:       t.TargetEnd,
		Title:           t.Title,
		Description:     t.Description,
		Version:         t.Version,
		BranchTag:       t.BranchTag,
		CommitHash:      t.CommitHash,
		BuildId:         t.BuildId,
		PercentComplete: t.PercentComplete,
		Environment:     t.Environment,
		Lead:            t.Lead,
		ScanType:        t.ScanType,
		Tags:            t.Tags,
	}
}

func (ddr *ddTestDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	tflog.Info(ctx, "ddTestDefectdojoResource createApiCall")
	reqBody := ddTestToCreateRequest(ddr.TestCreate)
	apiResp, err := client.TestsCreateWithResponse(ctx, reqBody)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON201 != nil {
		// Copy fields from Test back into TestCreate
		t := apiResp.JSON201
		ddr.TestCreate.Id = t.Id
		ddr.TestCreate.TestType = t.TestType
		ddr.TestCreate.Engagement = t.Engagement
		ddr.TestCreate.TargetStart = t.TargetStart
		ddr.TestCreate.TargetEnd = t.TargetEnd
		ddr.TestCreate.Title = t.Title
		ddr.TestCreate.Description = t.Description
		ddr.TestCreate.Version = t.Version
		ddr.TestCreate.BranchTag = t.BranchTag
		ddr.TestCreate.CommitHash = t.CommitHash
		ddr.TestCreate.BuildId = t.BuildId
		ddr.TestCreate.PercentComplete = t.PercentComplete
		ddr.TestCreate.Environment = t.Environment
		ddr.TestCreate.Lead = t.Lead
		ddr.TestCreate.ScanType = t.ScanType
		ddr.TestCreate.Tags = t.Tags
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *ddTestDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "ddTestDefectdojoResource readApiCall")
	apiResp, err := client.TestsRetrieveWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON200 != nil {
		t := apiResp.JSON200
		ddr.TestCreate.Id = t.Id
		ddr.TestCreate.TestType = t.TestType
		if t.Engagement != nil {
			ddr.TestCreate.Engagement = *t.Engagement
		}
		ddr.TestCreate.TargetStart = t.TargetStart
		ddr.TestCreate.TargetEnd = t.TargetEnd
		ddr.TestCreate.Title = t.Title
		ddr.TestCreate.Description = t.Description
		ddr.TestCreate.Version = t.Version
		ddr.TestCreate.BranchTag = t.BranchTag
		ddr.TestCreate.CommitHash = t.CommitHash
		ddr.TestCreate.BuildId = t.BuildId
		ddr.TestCreate.PercentComplete = t.PercentComplete
		ddr.TestCreate.Environment = t.Environment
		ddr.TestCreate.Lead = t.Lead
		ddr.TestCreate.ScanType = t.ScanType
		ddr.TestCreate.Tags = t.Tags
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *ddTestDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "ddTestDefectdojoResource updateApiCall")
	reqBody := ddTestToRequest(ddr.TestCreate)
	apiResp, err := client.TestsUpdateWithResponse(ctx, idNumber, reqBody)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON200 != nil {
		t := apiResp.JSON200
		ddr.TestCreate.Id = t.Id
		ddr.TestCreate.TestType = t.TestType
		if t.Engagement != nil {
			ddr.TestCreate.Engagement = *t.Engagement
		}
		ddr.TestCreate.TargetStart = t.TargetStart
		ddr.TestCreate.TargetEnd = t.TargetEnd
		ddr.TestCreate.Title = t.Title
		ddr.TestCreate.Description = t.Description
		ddr.TestCreate.Version = t.Version
		ddr.TestCreate.BranchTag = t.BranchTag
		ddr.TestCreate.CommitHash = t.CommitHash
		ddr.TestCreate.BuildId = t.BuildId
		ddr.TestCreate.PercentComplete = t.PercentComplete
		ddr.TestCreate.Environment = t.Environment
		ddr.TestCreate.Lead = t.Lead
		ddr.TestCreate.ScanType = t.ScanType
		ddr.TestCreate.Tags = t.Tags
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *ddTestDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "ddTestDefectdojoResource deleteApiCall")
	apiResp, err := client.TestsDestroyWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	return apiResp.StatusCode(), apiResp.Body, nil
}

type ddTestResource struct {
	terraformResource
}

var _ resource.Resource = &ddTestResource{}
var _ resource.ResourceWithImportState = &ddTestResource{}
var _ resource.ResourceWithConfigure = &ddTestResource{}

func NewDDTestResource() resource.Resource {
	return &ddTestResource{
		terraformResource: terraformResource{
			dataProvider: ddTestDataProvider{},
		},
	}
}

func (r ddTestResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_test"
}

type ddTestDataProvider struct{}

func (r ddTestDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data ddTestResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *ddTestResourceData) id() types.String {
	return d.Id
}

func (d *ddTestResourceData) defectdojoResource() defectdojoResource {
	return &ddTestDefectdojoResource{
		TestCreate: dd.TestCreate{},
	}
}
