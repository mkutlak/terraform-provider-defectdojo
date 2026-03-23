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

func (t jiraInstanceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Jira Instance",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "The URL of the Jira instance.",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username or Email Address for Jira authentication.",
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password or API Token for Jira authentication.",
				Optional:            true,
				Sensitive:           true,
			},
			"configuration_name": schema.StringAttribute{
				MarkdownDescription: "A name for this Jira configuration.",
				Optional:            true,
				Computed:            true,
			},
			"epic_name_id": schema.Int64Attribute{
				MarkdownDescription: "The Epic Name field ID in Jira (from cf[number] in /rest/api/2/field).",
				Required:            true,
			},
			"open_status_key": schema.Int64Attribute{
				MarkdownDescription: "Transition ID to Re-Open JIRA issues.",
				Required:            true,
			},
			"close_status_key": schema.Int64Attribute{
				MarkdownDescription: "Transition ID to Close JIRA issues.",
				Required:            true,
			},
			"info_mapping_severity": schema.StringAttribute{
				MarkdownDescription: "Maps to the Priority field in Jira for Info severity.",
				Required:            true,
			},
			"low_mapping_severity": schema.StringAttribute{
				MarkdownDescription: "Maps to the Priority field in Jira for Low severity.",
				Required:            true,
			},
			"medium_mapping_severity": schema.StringAttribute{
				MarkdownDescription: "Maps to the Priority field in Jira for Medium severity.",
				Required:            true,
			},
			"high_mapping_severity": schema.StringAttribute{
				MarkdownDescription: "Maps to the Priority field in Jira for High severity.",
				Required:            true,
			},
			"critical_mapping_severity": schema.StringAttribute{
				MarkdownDescription: "Maps to the Priority field in Jira for Critical severity.",
				Required:            true,
			},
			"finding_text": schema.StringAttribute{
				MarkdownDescription: "Additional text added to findings in Jira.",
				Optional:            true,
				Computed:            true,
			},
			"accepted_mapping_resolution": schema.StringAttribute{
				MarkdownDescription: "JIRA resolution that maps to Risk Accepted in DefectDojo.",
				Optional:            true,
				Computed:            true,
			},
			"false_positive_mapping_resolution": schema.StringAttribute{
				MarkdownDescription: "JIRA resolution that maps to False Positive in DefectDojo.",
				Optional:            true,
				Computed:            true,
			},
			"global_jira_sla_notification": schema.BoolAttribute{
				MarkdownDescription: "Send SLA notifications via Jira.",
				Optional:            true,
				Computed:            true,
			},
			"finding_jira_sync": schema.BoolAttribute{
				MarkdownDescription: "If enabled, sync finding changes automatically to JIRA.",
				Optional:            true,
				Computed:            true,
			},
			"issue_template_dir": schema.StringAttribute{
				MarkdownDescription: "Folder containing Django templates for JIRA issue descriptions.",
				Optional:            true,
				Computed:            true,
			},
			"default_issue_type": schema.StringAttribute{
				MarkdownDescription: "Default issue type. Valid values: Task, Story, Epic, Spike, Bug, Security.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

type jiraInstanceResourceData struct {
	Id                             types.String `tfsdk:"id" ddField:"Id"`
	Url                            types.String `tfsdk:"url" ddField:"Url"`
	Username                       types.String `tfsdk:"username" ddField:"Username"`
	Password                       types.String `tfsdk:"password" ddField:"Password"`
	ConfigurationName              types.String `tfsdk:"configuration_name" ddField:"ConfigurationName"`
	EpicNameId                     types.Int64  `tfsdk:"epic_name_id" ddField:"EpicNameId"`
	OpenStatusKey                  types.Int64  `tfsdk:"open_status_key" ddField:"OpenStatusKey"`
	CloseStatusKey                 types.Int64  `tfsdk:"close_status_key" ddField:"CloseStatusKey"`
	InfoMappingSeverity            types.String `tfsdk:"info_mapping_severity" ddField:"InfoMappingSeverity"`
	LowMappingSeverity             types.String `tfsdk:"low_mapping_severity" ddField:"LowMappingSeverity"`
	MediumMappingSeverity          types.String `tfsdk:"medium_mapping_severity" ddField:"MediumMappingSeverity"`
	HighMappingSeverity            types.String `tfsdk:"high_mapping_severity" ddField:"HighMappingSeverity"`
	CriticalMappingSeverity        types.String `tfsdk:"critical_mapping_severity" ddField:"CriticalMappingSeverity"`
	FindingText                    types.String `tfsdk:"finding_text" ddField:"FindingText"`
	AcceptedMappingResolution      types.String `tfsdk:"accepted_mapping_resolution" ddField:"AcceptedMappingResolution"`
	FalsePositiveMappingResolution types.String `tfsdk:"false_positive_mapping_resolution" ddField:"FalsePositiveMappingResolution"`
	GlobalJiraSlaNotification      types.Bool   `tfsdk:"global_jira_sla_notification" ddField:"GlobalJiraSlaNotification"`
	FindingJiraSync                types.Bool   `tfsdk:"finding_jira_sync" ddField:"FindingJiraSync"`
	IssueTemplateDir               types.String `tfsdk:"issue_template_dir" ddField:"IssueTemplateDir"`
	DefaultIssueType               types.String `tfsdk:"default_issue_type" ddField:"DefaultIssueType"`
}

type jiraInstanceDefectdojoResource struct {
	dd.JIRAInstance
}

func jiraInstanceToRequest(j dd.JIRAInstance) dd.JIRAInstanceRequest {
	req := dd.JIRAInstanceRequest{
		Url:                            j.Url,
		Username:                       j.Username,
		ConfigurationName:              j.ConfigurationName,
		EpicNameId:                     j.EpicNameId,
		OpenStatusKey:                  j.OpenStatusKey,
		CloseStatusKey:                 j.CloseStatusKey,
		InfoMappingSeverity:            j.InfoMappingSeverity,
		LowMappingSeverity:             j.LowMappingSeverity,
		MediumMappingSeverity:          j.MediumMappingSeverity,
		HighMappingSeverity:            j.HighMappingSeverity,
		CriticalMappingSeverity:        j.CriticalMappingSeverity,
		FindingText:                    j.FindingText,
		AcceptedMappingResolution:      j.AcceptedMappingResolution,
		FalsePositiveMappingResolution: j.FalsePositiveMappingResolution,
		GlobalJiraSlaNotification:      j.GlobalJiraSlaNotification,
		FindingJiraSync:                j.FindingJiraSync,
		IssueTemplateDir:               j.IssueTemplateDir,
	}
	if j.DefaultIssueType != nil {
		v := dd.JIRAInstanceRequestDefaultIssueType(*j.DefaultIssueType)
		req.DefaultIssueType = &v
	}
	return req
}

func (ddr *jiraInstanceDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := jiraInstanceToRequest(ddr.JIRAInstance)
	apiResp, err := client.JiraInstancesCreateWithResponse(ctx, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON201 != nil {
		ddr.JIRAInstance = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *jiraInstanceDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.JiraInstancesRetrieveWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.JIRAInstance = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *jiraInstanceDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := jiraInstanceToRequest(ddr.JIRAInstance)
	apiResp, err := client.JiraInstancesUpdateWithResponse(ctx, idNumber, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.JIRAInstance = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *jiraInstanceDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.JiraInstancesDestroyWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (d *jiraInstanceResourceData) id() types.String {
	return d.Id
}

func (d *jiraInstanceResourceData) defectdojoResource() defectdojoResource {
	return &jiraInstanceDefectdojoResource{JIRAInstance: dd.JIRAInstance{}}
}

type jiraInstanceResource struct {
	terraformResource
}

var _ resource.Resource = &jiraInstanceResource{}
var _ resource.ResourceWithImportState = &jiraInstanceResource{}

func NewJiraInstanceResource() resource.Resource {
	return &jiraInstanceResource{
		terraformResource: terraformResource{
			typeName:     "defectdojo_jira_instance",
			dataProvider: jiraInstanceDataProvider{},
		},
	}
}

func (r jiraInstanceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_instance"
}

type jiraInstanceDataProvider struct{}

func (r jiraInstanceDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data jiraInstanceResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}
