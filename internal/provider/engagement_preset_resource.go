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

func (r engagementPresetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Engagement Preset",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "Brief description of preset",
				Optional:            true,
			},
			"product": schema.Int64Attribute{
				MarkdownDescription: "ID of the Product this Preset belongs to",
				Required:            true,
			},
			"notes": schema.StringAttribute{
				MarkdownDescription: "Description of what needs to be tested or setting up environment for testing",
				Optional:            true,
			},
			"scope": schema.StringAttribute{
				MarkdownDescription: "Scope of Engagement testing, IP's/Resources/URL's",
				Optional:            true,
			},
			"network_locations": schema.SetAttribute{
				MarkdownDescription: "IDs of network locations",
				Optional:            true,
				ElementType:         types.Int64Type,
			},
			"test_type": schema.SetAttribute{
				MarkdownDescription: "IDs of test types",
				Optional:            true,
				ElementType:         types.Int64Type,
			},
		},
	}
}

type engagementPresetResourceData struct {
	Id               types.String `tfsdk:"id" ddField:"Id"`
	Title            types.String `tfsdk:"title" ddField:"Title"`
	Product          types.Int64  `tfsdk:"product" ddField:"Product"`
	Notes            types.String `tfsdk:"notes" ddField:"Notes"`
	Scope            types.String `tfsdk:"scope" ddField:"Scope"`
	NetworkLocations types.Set    `tfsdk:"network_locations" ddField:"NetworkLocations"`
	TestType         types.Set    `tfsdk:"test_type" ddField:"TestType"`
}

type engagementPresetDefectdojoResource struct {
	dd.EngagementPresets
}

func engagementPresetToRequest(e dd.EngagementPresets) dd.EngagementPresetsRequest {
	return dd.EngagementPresetsRequest{
		Title:            e.Title,
		Product:          e.Product,
		Notes:            e.Notes,
		Scope:            e.Scope,
		NetworkLocations: e.NetworkLocations,
		TestType:         e.TestType,
	}
}

func (ddr *engagementPresetDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	tflog.Info(ctx, "engagementPresetDefectdojoResource createApiCall")
	reqBody := engagementPresetToRequest(ddr.EngagementPresets)
	apiResp, err := client.EngagementPresetsCreateWithResponse(ctx, reqBody)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON201 != nil {
		ddr.EngagementPresets = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *engagementPresetDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "engagementPresetDefectdojoResource readApiCall")
	apiResp, err := client.EngagementPresetsRetrieveWithResponse(ctx, idNumber, &dd.EngagementPresetsRetrieveParams{})
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON200 != nil {
		ddr.EngagementPresets = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *engagementPresetDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "engagementPresetDefectdojoResource updateApiCall")
	reqBody := engagementPresetToRequest(ddr.EngagementPresets)
	apiResp, err := client.EngagementPresetsUpdateWithResponse(ctx, idNumber, reqBody)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON200 != nil {
		ddr.EngagementPresets = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *engagementPresetDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "engagementPresetDefectdojoResource deleteApiCall")
	apiResp, err := client.EngagementPresetsDestroyWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	return apiResp.StatusCode(), apiResp.Body, nil
}

type engagementPresetResource struct {
	terraformResource
}

var _ resource.Resource = &engagementPresetResource{}
var _ resource.ResourceWithImportState = &engagementPresetResource{}

func NewEngagementPresetResource() resource.Resource {
	return &engagementPresetResource{
		terraformResource: terraformResource{
			typeName:     "defectdojo_engagement_preset",
			dataProvider: engagementPresetDataProvider{},
		},
	}
}

func (r engagementPresetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_engagement_preset"
}

type engagementPresetDataProvider struct{}

func (r engagementPresetDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data engagementPresetResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *engagementPresetResourceData) id() types.String {
	return d.Id
}

func (d *engagementPresetResourceData) setId(v types.String) { d.Id = v }

func (d *engagementPresetResourceData) defectdojoResource() defectdojoResource {
	return &engagementPresetDefectdojoResource{
		EngagementPresets: dd.EngagementPresets{},
	}
}
