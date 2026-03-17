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

func (t endpointStatusResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Endpoint Status",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.Int64Attribute{
				MarkdownDescription: "The endpoint this status is associated with.",
				Required:            true,
			},
			"finding": schema.Int64Attribute{
				MarkdownDescription: "The finding this status is associated with.",
				Required:            true,
			},
			"date": schema.StringAttribute{
				MarkdownDescription: "Date of the endpoint status (format: 2006-01-02).",
				Optional:            true,
			},
			"false_positive": schema.BoolAttribute{
				MarkdownDescription: "Whether the finding is a false positive.",
				Optional:            true,
			},
			"mitigated": schema.BoolAttribute{
				MarkdownDescription: "Whether the finding has been mitigated.",
				Optional:            true,
			},
			"mitigated_by": schema.Int64Attribute{
				MarkdownDescription: "The user who mitigated the finding.",
				Optional:            true,
			},
			"out_of_scope": schema.BoolAttribute{
				MarkdownDescription: "Whether the finding is out of scope.",
				Optional:            true,
			},
			"risk_accepted": schema.BoolAttribute{
				MarkdownDescription: "Whether the risk has been accepted.",
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

type endpointStatusResourceData struct {
	Endpoint      types.Int64  `tfsdk:"endpoint" ddField:"Endpoint"`
	Finding       types.Int64  `tfsdk:"finding" ddField:"Finding"`
	Date          types.String `tfsdk:"date" ddField:"Date"`
	FalsePositive types.Bool   `tfsdk:"false_positive" ddField:"FalsePositive"`
	Mitigated     types.Bool   `tfsdk:"mitigated" ddField:"Mitigated"`
	MitigatedBy   types.Int64  `tfsdk:"mitigated_by" ddField:"MitigatedBy"`
	OutOfScope    types.Bool   `tfsdk:"out_of_scope" ddField:"OutOfScope"`
	RiskAccepted  types.Bool   `tfsdk:"risk_accepted" ddField:"RiskAccepted"`
	Id            types.String `tfsdk:"id" ddField:"Id"`
}

type endpointStatusDefectdojoResource struct {
	dd.EndpointStatus
}

func endpointStatusToRequest(obj dd.EndpointStatus) dd.EndpointStatusRequest {
	return dd.EndpointStatusRequest{
		Endpoint:      obj.Endpoint,
		Finding:       obj.Finding,
		Date:          obj.Date,
		FalsePositive: obj.FalsePositive,
		Mitigated:     obj.Mitigated,
		MitigatedBy:   obj.MitigatedBy,
		OutOfScope:    obj.OutOfScope,
		RiskAccepted:  obj.RiskAccepted,
	}
}

func (ddr *endpointStatusDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := endpointStatusToRequest(ddr.EndpointStatus)
	apiResp, err := client.EndpointStatusCreateWithResponse(ctx, reqBody)
	if apiResp.JSON201 != nil {
		ddr.EndpointStatus = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *endpointStatusDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.EndpointStatusRetrieveWithResponse(ctx, idNumber)
	if apiResp.JSON200 != nil {
		ddr.EndpointStatus = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *endpointStatusDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := endpointStatusToRequest(ddr.EndpointStatus)
	apiResp, err := client.EndpointStatusUpdateWithResponse(ctx, idNumber, reqBody)
	if apiResp.JSON200 != nil {
		ddr.EndpointStatus = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *endpointStatusDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.EndpointStatusDestroyWithResponse(ctx, idNumber)
	return apiResp.StatusCode(), apiResp.Body, err
}

type endpointStatusResource struct {
	terraformResource
}

var _ resource.Resource = &endpointStatusResource{}
var _ resource.ResourceWithImportState = &endpointStatusResource{}

func NewEndpointStatusResource() resource.Resource {
	return &endpointStatusResource{
		terraformResource: terraformResource{dataProvider: endpointStatusDataProvider{}},
	}
}

func (r endpointStatusResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_endpoint_status"
}

type endpointStatusDataProvider struct{}

func (r endpointStatusDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data endpointStatusResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *endpointStatusResourceData) id() types.String { return d.Id }

func (d *endpointStatusResourceData) defectdojoResource() defectdojoResource {
	return &endpointStatusDefectdojoResource{EndpointStatus: dd.EndpointStatus{}}
}
