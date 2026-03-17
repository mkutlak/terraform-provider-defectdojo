package provider

import (
	"context"

	dd "github.com/doximity/terraform-provider-defectdojo/internal/ddclient"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (t riskAcceptanceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Risk Acceptance",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Descriptive name for the risk acceptance.",
				Required:            true,
			},
			"owner": schema.Int64Attribute{
				MarkdownDescription: "User ID of the owner of this risk acceptance.",
				Required:            true,
			},
			"accepted_findings": schema.SetAttribute{
				MarkdownDescription: "IDs of findings included in this risk acceptance.",
				Required:            true,
				ElementType:         types.Int64Type,
			},
			"accepted_by": schema.StringAttribute{
				MarkdownDescription: "The person that accepts the risk (can be outside DefectDojo).",
				Optional:            true,
				Computed:            true,
			},
			"expiration_date": schema.StringAttribute{
				MarkdownDescription: "When the risk acceptance expires (RFC3339 format).",
				Optional:            true,
				Computed:            true,
			},
			"decision": schema.StringAttribute{
				MarkdownDescription: "Risk treatment decision. Valid values: A (Accept), V (Avoid), M (Mitigate), F (Fix), T (Transfer).",
				Optional:            true,
				Computed:            true,
			},
			"decision_details": schema.StringAttribute{
				MarkdownDescription: "Details about the risk treatment decision.",
				Optional:            true,
				Computed:            true,
			},
			"recommendation": schema.StringAttribute{
				MarkdownDescription: "Security team recommendation. Valid values: A, V, M, F, T.",
				Optional:            true,
				Computed:            true,
			},
			"recommendation_details": schema.StringAttribute{
				MarkdownDescription: "Details explaining the security recommendation.",
				Optional:            true,
				Computed:            true,
			},
			"reactivate_expired": schema.BoolAttribute{
				MarkdownDescription: "Reactivate findings when risk acceptance expires.",
				Optional:            true,
				Computed:            true,
			},
			"restart_sla_expired": schema.BoolAttribute{
				MarkdownDescription: "Restart SLA for findings when risk acceptance expires.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

type riskAcceptanceResourceData struct {
	Id                    types.String `tfsdk:"id" ddField:"Id"`
	Name                  types.String `tfsdk:"name" ddField:"Name"`
	Owner                 types.Int64  `tfsdk:"owner" ddField:"Owner"`
	AcceptedFindings      types.Set    `tfsdk:"accepted_findings" ddField:"AcceptedFindings"`
	AcceptedBy            types.String `tfsdk:"accepted_by" ddField:"AcceptedBy"`
	ExpirationDate        types.String `tfsdk:"expiration_date" ddField:"ExpirationDate"`
	Decision              types.String `tfsdk:"decision" ddField:"Decision"`
	DecisionDetails       types.String `tfsdk:"decision_details" ddField:"DecisionDetails"`
	Recommendation        types.String `tfsdk:"recommendation" ddField:"Recommendation"`
	RecommendationDetails types.String `tfsdk:"recommendation_details" ddField:"RecommendationDetails"`
	ReactivateExpired     types.Bool   `tfsdk:"reactivate_expired" ddField:"ReactivateExpired"`
	RestartSlaExpired     types.Bool   `tfsdk:"restart_sla_expired" ddField:"RestartSlaExpired"`
}

type riskAcceptanceDefectdojoResource struct {
	dd.RiskAcceptance
}

func riskAcceptanceToRequest(r dd.RiskAcceptance) dd.RiskAcceptanceRequest {
	req := dd.RiskAcceptanceRequest{
		Name:                  r.Name,
		Owner:                 r.Owner,
		AcceptedFindings:      r.AcceptedFindings,
		AcceptedBy:            r.AcceptedBy,
		ExpirationDate:        r.ExpirationDate,
		DecisionDetails:       r.DecisionDetails,
		RecommendationDetails: r.RecommendationDetails,
		ReactivateExpired:     r.ReactivateExpired,
		RestartSlaExpired:     r.RestartSlaExpired,
	}
	if r.Decision != nil {
		v := dd.RiskAcceptanceRequestDecision(*r.Decision)
		req.Decision = &v
	}
	if r.Recommendation != nil {
		v := dd.RiskAcceptanceRequestRecommendation(*r.Recommendation)
		req.Recommendation = &v
	}
	return req
}

func (ddr *riskAcceptanceDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := riskAcceptanceToRequest(ddr.RiskAcceptance)
	apiResp, err := client.RiskAcceptanceCreateWithResponse(ctx, reqBody)
	if apiResp.JSON201 != nil {
		ddr.RiskAcceptance = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *riskAcceptanceDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.RiskAcceptanceRetrieveWithResponse(ctx, idNumber)
	if apiResp.JSON200 != nil {
		ddr.RiskAcceptance = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *riskAcceptanceDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := riskAcceptanceToRequest(ddr.RiskAcceptance)
	apiResp, err := client.RiskAcceptanceUpdateWithResponse(ctx, idNumber, reqBody)
	if apiResp.JSON200 != nil {
		ddr.RiskAcceptance = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *riskAcceptanceDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.RiskAcceptanceDestroyWithResponse(ctx, idNumber)
	return apiResp.StatusCode(), apiResp.Body, err
}

func (d *riskAcceptanceResourceData) id() types.String {
	return d.Id
}

func (d *riskAcceptanceResourceData) defectdojoResource() defectdojoResource {
	return &riskAcceptanceDefectdojoResource{RiskAcceptance: dd.RiskAcceptance{}}
}

type riskAcceptanceResource struct {
	terraformResource
}

var _ resource.Resource = &riskAcceptanceResource{}
var _ resource.ResourceWithImportState = &riskAcceptanceResource{}

func NewRiskAcceptanceResource() resource.Resource {
	return &riskAcceptanceResource{
		terraformResource: terraformResource{
			dataProvider: riskAcceptanceDataProvider{},
		},
	}
}

func (r riskAcceptanceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_risk_acceptance"
}

type riskAcceptanceDataProvider struct{}

func (r riskAcceptanceDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data riskAcceptanceResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}
