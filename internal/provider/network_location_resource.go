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

func (t networkLocationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Network Location",
		Attributes: map[string]schema.Attribute{
			"location": schema.StringAttribute{
				MarkdownDescription: "Location of network testing: Examples: VPN, Internet or Internal",
				Required:            true,
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

type networkLocationResourceData struct {
	Location types.String `tfsdk:"location" ddField:"Location"`
	Id       types.String `tfsdk:"id" ddField:"Id"`
}

type networkLocationDefectdojoResource struct {
	dd.NetworkLocations
}

func networkLocationToRequest(obj dd.NetworkLocations) dd.NetworkLocationsRequest {
	return dd.NetworkLocationsRequest{
		Location: obj.Location,
	}
}

func (ddr *networkLocationDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := networkLocationToRequest(ddr.NetworkLocations)
	apiResp, err := client.NetworkLocationsCreateWithResponse(ctx, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON201 != nil {
		ddr.NetworkLocations = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *networkLocationDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.NetworkLocationsRetrieveWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.NetworkLocations = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *networkLocationDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := networkLocationToRequest(ddr.NetworkLocations)
	apiResp, err := client.NetworkLocationsUpdateWithResponse(ctx, idNumber, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.NetworkLocations = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *networkLocationDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.NetworkLocationsDestroyWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

type networkLocationResource struct {
	terraformResource
}

var _ resource.Resource = &networkLocationResource{}
var _ resource.ResourceWithImportState = &networkLocationResource{}

func NewNetworkLocationResource() resource.Resource {
	return &networkLocationResource{
		terraformResource: terraformResource{dataProvider: networkLocationDataProvider{}},
	}
}

func (r networkLocationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_location"
}

type networkLocationDataProvider struct{}

func (r networkLocationDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data networkLocationResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *networkLocationResourceData) id() types.String { return d.Id }

func (d *networkLocationResourceData) defectdojoResource() defectdojoResource {
	return &networkLocationDefectdojoResource{NetworkLocations: dd.NetworkLocations{}}
}
