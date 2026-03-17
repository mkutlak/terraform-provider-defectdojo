package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
)

func (t findingTemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Finding Template",
		Attributes: map[string]schema.Attribute{
			"title": schema.StringAttribute{
				MarkdownDescription: "The title of the Finding Template",
				Required:            true,
			},
			"severity": schema.StringAttribute{
				MarkdownDescription: "The severity of the finding (Critical, High, Medium, Low, Info)",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Detailed description of the vulnerability",
				Optional:            true,
			},
			"mitigation": schema.StringAttribute{
				MarkdownDescription: "Steps to mitigate the vulnerability",
				Optional:            true,
			},
			"impact": schema.StringAttribute{
				MarkdownDescription: "Impact of the vulnerability",
				Optional:            true,
			},
			"references": schema.StringAttribute{
				MarkdownDescription: "References for the vulnerability",
				Optional:            true,
			},
			"cwe": schema.Int64Attribute{
				MarkdownDescription: "CWE number",
				Optional:            true,
			},
			"cvssv3": schema.StringAttribute{
				MarkdownDescription: "CVSSv3 score string",
				Optional:            true,
			},
			"cvssv3_score": schema.Float64Attribute{
				MarkdownDescription: "CVSSv3 numeric score",
				Optional:            true,
			},
			"cvssv4": schema.StringAttribute{
				MarkdownDescription: "CVSSv4 score string",
				Optional:            true,
			},
			"cvssv4_score": schema.Float64Attribute{
				MarkdownDescription: "CVSSv4 numeric score",
				Optional:            true,
			},
			"component_name": schema.StringAttribute{
				MarkdownDescription: "Affected component name",
				Optional:            true,
			},
			"component_version": schema.StringAttribute{
				MarkdownDescription: "Affected component version",
				Optional:            true,
			},
			"fix_available": schema.BoolAttribute{
				MarkdownDescription: "Indicates if a fix is available",
				Optional:            true,
			},
			"fix_version": schema.StringAttribute{
				MarkdownDescription: "Version where fix is available",
				Optional:            true,
			},
			"planned_remediation_version": schema.StringAttribute{
				MarkdownDescription: "Target version for remediation",
				Optional:            true,
			},
			"effort_for_fixing": schema.StringAttribute{
				MarkdownDescription: "Effort estimate for fixing (e.g., Low/Medium/High)",
				Optional:            true,
			},
			"severity_justification": schema.StringAttribute{
				MarkdownDescription: "Explanation of why this severity level is appropriate",
				Optional:            true,
			},
			"steps_to_reproduce": schema.StringAttribute{
				MarkdownDescription: "Standard reproduction steps for this vulnerability type",
				Optional:            true,
			},
			"endpoints_text": schema.StringAttribute{
				MarkdownDescription: "Endpoint URLs (one per line)",
				Optional:            true,
			},
			"notes": schema.StringAttribute{
				MarkdownDescription: "Note content to add when applying this template",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

type findingTemplateResourceData struct {
	Title                     types.String  `tfsdk:"title" ddField:"Title"`
	Severity                  types.String  `tfsdk:"severity" ddField:"Severity"`
	Description               types.String  `tfsdk:"description" ddField:"Description"`
	Mitigation                types.String  `tfsdk:"mitigation" ddField:"Mitigation"`
	Impact                    types.String  `tfsdk:"impact" ddField:"Impact"`
	References                types.String  `tfsdk:"references" ddField:"References"`
	Cwe                       types.Int64   `tfsdk:"cwe" ddField:"Cwe"`
	Cvssv3                    types.String  `tfsdk:"cvssv3" ddField:"Cvssv3"`
	Cvssv3Score               types.Float64 `tfsdk:"cvssv3_score" ddField:"Cvssv3Score"`
	Cvssv4                    types.String  `tfsdk:"cvssv4" ddField:"Cvssv4"`
	Cvssv4Score               types.Float64 `tfsdk:"cvssv4_score" ddField:"Cvssv4Score"`
	ComponentName             types.String  `tfsdk:"component_name" ddField:"ComponentName"`
	ComponentVersion          types.String  `tfsdk:"component_version" ddField:"ComponentVersion"`
	FixAvailable              types.Bool    `tfsdk:"fix_available" ddField:"FixAvailable"`
	FixVersion                types.String  `tfsdk:"fix_version" ddField:"FixVersion"`
	PlannedRemediationVersion types.String  `tfsdk:"planned_remediation_version" ddField:"PlannedRemediationVersion"`
	EffortForFixing           types.String  `tfsdk:"effort_for_fixing" ddField:"EffortForFixing"`
	SeverityJustification     types.String  `tfsdk:"severity_justification" ddField:"SeverityJustification"`
	StepsToReproduce          types.String  `tfsdk:"steps_to_reproduce" ddField:"StepsToReproduce"`
	EndpointsText             types.String  `tfsdk:"endpoints_text" ddField:"EndpointsText"`
	Notes                     types.String  `tfsdk:"notes" ddField:"Notes"`
	Id                        types.String  `tfsdk:"id" ddField:"Id"`
}

type findingTemplateDefectdojoResource struct {
	dd.FindingTemplate
}

func findingTemplateToRequest(obj dd.FindingTemplate) dd.FindingTemplateRequest {
	return dd.FindingTemplateRequest{
		Title:                     obj.Title,
		Severity:                  obj.Severity,
		Description:               obj.Description,
		Mitigation:                obj.Mitigation,
		Impact:                    obj.Impact,
		References:                obj.References,
		Cwe:                       obj.Cwe,
		Cvssv3:                    obj.Cvssv3,
		Cvssv3Score:               obj.Cvssv3Score,
		Cvssv4:                    obj.Cvssv4,
		Cvssv4Score:               obj.Cvssv4Score,
		ComponentName:             obj.ComponentName,
		ComponentVersion:          obj.ComponentVersion,
		FixAvailable:              obj.FixAvailable,
		FixVersion:                obj.FixVersion,
		PlannedRemediationVersion: obj.PlannedRemediationVersion,
		EffortForFixing:           obj.EffortForFixing,
		SeverityJustification:     obj.SeverityJustification,
		StepsToReproduce:          obj.StepsToReproduce,
		EndpointsText:             obj.EndpointsText,
		Notes:                     obj.Notes,
	}
}

func (ddr *findingTemplateDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := findingTemplateToRequest(ddr.FindingTemplate)
	apiResp, err := client.FindingTemplatesCreateWithResponse(ctx, reqBody)
	if apiResp.JSON201 != nil {
		ddr.FindingTemplate = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *findingTemplateDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.FindingTemplatesRetrieveWithResponse(ctx, idNumber)
	if apiResp.JSON200 != nil {
		ddr.FindingTemplate = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *findingTemplateDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := findingTemplateToRequest(ddr.FindingTemplate)
	apiResp, err := client.FindingTemplatesUpdateWithResponse(ctx, idNumber, reqBody)
	if apiResp.JSON200 != nil {
		ddr.FindingTemplate = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *findingTemplateDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.FindingTemplatesDestroyWithResponse(ctx, idNumber)
	return apiResp.StatusCode(), apiResp.Body, err
}

type findingTemplateResource struct {
	terraformResource
}

var _ resource.Resource = &findingTemplateResource{}
var _ resource.ResourceWithImportState = &findingTemplateResource{}

func NewFindingTemplateResource() resource.Resource {
	return &findingTemplateResource{
		terraformResource: terraformResource{dataProvider: findingTemplateDataProvider{}},
	}
}

func (r findingTemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_finding_template"
}

type findingTemplateDataProvider struct{}

func (r findingTemplateDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data findingTemplateResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *findingTemplateResourceData) id() types.String { return d.Id }

func (d *findingTemplateResourceData) defectdojoResource() defectdojoResource {
	return &findingTemplateDefectdojoResource{FindingTemplate: dd.FindingTemplate{}}
}
