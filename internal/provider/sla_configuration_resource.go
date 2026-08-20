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

// slaDayCountDescription builds the description shared by the four SLA day
// counts, which differ only in the severity they name. The warning is long
// enough that four copies would drift apart.
//
// It has to be spelled out because the obvious reading of an Optional
// attribute - delete the line, get the default back - is wrong here, and wrong
// in a way that leaves the server holding a value nothing in the configuration
// mentions any more.
func slaDayCountDescription(severity string) string {
	return "The number of days to remediate a " + severity + " finding. Omitting this attribute when the SLA " +
		"configuration is created lets DefectDojo apply its own default (which varies by version), and that value is " +
		"read back into state. Removing it from an existing configuration does not put the default back: the attribute " +
		"is computed, so Terraform reuses the value already in state, plans no change, and DefectDojo keeps what it " +
		"holds. DefectDojo offers no way back either - the column is not nullable, so a request sending null is refused. " +
		"To change the value, set it explicitly; to hand it back to DefectDojo, recreate the resource with " +
		"`terraform apply -replace=...`, which restores the version default but assigns a new `id`, so every " +
		"`defectdojo_product` whose `sla_configuration` points at it has to be updated to match. `terraform state rm` " +
		"does not help: re-importing reads the stored value straight back."
}

func (t slaConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo SLA Configuration",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "A unique name for the set of SLAs",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the SLA Configuration",
				Optional:            true,
			},
			// The SLA day counts are not nullable in DefectDojo: leaving one
			// unset in a create/update request makes the server apply its own
			// model default and answer with it. They must therefore be Computed
			// to avoid a "provider produced inconsistent result" error whenever
			// a config omits one. No Default is declared, because those model
			// defaults vary between DefectDojo versions - the server stays the
			// source of truth.
			//
			// Being Computed makes them write-once-then-sticky, which is what
			// slaDayCountDescription exists to say out loud.
			"critical": schema.Int64Attribute{
				MarkdownDescription: slaDayCountDescription("critical"),
				Optional:            true,
				Computed:            true,
			},
			"high": schema.Int64Attribute{
				MarkdownDescription: slaDayCountDescription("high"),
				Optional:            true,
				Computed:            true,
			},
			"medium": schema.Int64Attribute{
				MarkdownDescription: slaDayCountDescription("medium"),
				Optional:            true,
				Computed:            true,
			},
			"low": schema.Int64Attribute{
				MarkdownDescription: slaDayCountDescription("low"),
				Optional:            true,
				Computed:            true,
			},
			"enforce_critical": schema.BoolAttribute{
				MarkdownDescription: "When enabled, critical findings will be assigned an SLA expiration date",
				Optional:            true,
				Computed:            true,
			},
			"enforce_high": schema.BoolAttribute{
				MarkdownDescription: "When enabled, high findings will be assigned an SLA expiration date",
				Optional:            true,
				Computed:            true,
			},
			"enforce_medium": schema.BoolAttribute{
				MarkdownDescription: "When enabled, medium findings will be assigned an SLA expiration date",
				Optional:            true,
				Computed:            true,
			},
			"enforce_low": schema.BoolAttribute{
				MarkdownDescription: "When enabled, low findings will be assigned an SLA expiration date",
				Optional:            true,
				Computed:            true,
			},
			"restart_sla_on_reactivation": schema.BoolAttribute{
				MarkdownDescription: "When enabled, findings that were previously mitigated but are reactivated during reimport will have their SLA period restarted",
				Optional:            true,
				Computed:            true,
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

type slaConfigurationResourceData struct {
	Name                     types.String `tfsdk:"name" ddField:"Name"`
	Description              types.String `tfsdk:"description" ddField:"Description"`
	Critical                 types.Int64  `tfsdk:"critical" ddField:"Critical"`
	High                     types.Int64  `tfsdk:"high" ddField:"High"`
	Medium                   types.Int64  `tfsdk:"medium" ddField:"Medium"`
	Low                      types.Int64  `tfsdk:"low" ddField:"Low"`
	EnforceCritical          types.Bool   `tfsdk:"enforce_critical" ddField:"EnforceCritical"`
	EnforceHigh              types.Bool   `tfsdk:"enforce_high" ddField:"EnforceHigh"`
	EnforceMedium            types.Bool   `tfsdk:"enforce_medium" ddField:"EnforceMedium"`
	EnforceLow               types.Bool   `tfsdk:"enforce_low" ddField:"EnforceLow"`
	RestartSlaOnReactivation types.Bool   `tfsdk:"restart_sla_on_reactivation" ddField:"RestartSlaOnReactivation"`
	Id                       types.String `tfsdk:"id" ddField:"Id"`
}

type slaConfigurationDefectdojoResource struct {
	dd.SLAConfiguration
}

func slaConfigurationToRequest(obj dd.SLAConfiguration) dd.SLAConfigurationRequest {
	return dd.SLAConfigurationRequest{
		Name:                     obj.Name,
		Description:              obj.Description,
		Critical:                 obj.Critical,
		High:                     obj.High,
		Medium:                   obj.Medium,
		Low:                      obj.Low,
		EnforceCritical:          obj.EnforceCritical,
		EnforceHigh:              obj.EnforceHigh,
		EnforceMedium:            obj.EnforceMedium,
		EnforceLow:               obj.EnforceLow,
		RestartSlaOnReactivation: obj.RestartSlaOnReactivation,
	}
}

func (ddr *slaConfigurationDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := slaConfigurationToRequest(ddr.SLAConfiguration)
	apiResp, err := client.SlaConfigurationsCreateWithResponse(ctx, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON201 != nil {
		ddr.SLAConfiguration = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *slaConfigurationDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.SlaConfigurationsRetrieveWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.SLAConfiguration = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *slaConfigurationDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := slaConfigurationToRequest(ddr.SLAConfiguration)
	apiResp, err := client.SlaConfigurationsUpdateWithResponse(ctx, idNumber, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.SLAConfiguration = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *slaConfigurationDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.SlaConfigurationsDestroyWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

// partialUpdateWithBody names the PATCH that clears attributes removed from
// configuration. See clear.go.
func (ddr *slaConfigurationDefectdojoResource) partialUpdateWithBody(client *dd.ClientWithResponses) partialUpdateFunc {
	return client.SlaConfigurationsPartialUpdateWithBody
}

type slaConfigurationResource struct {
	terraformResource
}

var _ resource.Resource = &slaConfigurationResource{}
var _ resource.ResourceWithImportState = &slaConfigurationResource{}

func NewSlaConfigurationResource() resource.Resource {
	return &slaConfigurationResource{
		terraformResource: terraformResource{typeName: "defectdojo_sla_configuration", dataProvider: slaConfigurationDataProvider{}},
	}
}

func (r slaConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sla_configuration"
}

type slaConfigurationDataProvider struct{}

func (r slaConfigurationDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data slaConfigurationResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *slaConfigurationResourceData) id() types.String     { return d.Id }
func (d *slaConfigurationResourceData) setId(v types.String) { d.Id = v }

func (d *slaConfigurationResourceData) defectdojoResource() defectdojoResource {
	return &slaConfigurationDefectdojoResource{SLAConfiguration: dd.SLAConfiguration{}}
}
