package provider

import (
	"bytes"
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
)

func (t metadataResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Metadata: a custom key/value field attached to exactly one parent object. Exactly one of `product` or `finding` must be set. DefectDojo 3.1.101 does not support location- or endpoint-attached metadata via the API (the location parent is silently ignored and the endpoint parent is rejected), so only product and finding are exposed.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the metadata field.",
				Required:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "The value of the metadata field.",
				Required:            true,
			},
			"product": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Product this metadata is attached to. This is the recommended parent object for metadata. (Location- and endpoint-attached metadata are not supported: the DefectDojo 3.1.101 API ignores or rejects those parents despite advertising them.)",
				Optional:            true,
			},
			"finding": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Finding this metadata is attached to. Findings are import-managed; attaching metadata couples state to scan artifacts.",
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

type metadataResourceData struct {
	Id      types.String `tfsdk:"id" ddField:"Id"`
	Name    types.String `tfsdk:"name" ddField:"Name"`
	Value   types.String `tfsdk:"value" ddField:"Value"`
	Product types.Int64  `tfsdk:"product" ddField:"Product"`
	Finding types.Int64  `tfsdk:"finding" ddField:"Finding"`
}

type metadataDefectdojoResource struct {
	dd.Meta
}

// metadataToRequest converts a Meta (response model) to a MetaRequest (request model).
func metadataToRequest(m dd.Meta) dd.MetaRequest {
	return dd.MetaRequest{
		Finding: m.Finding,
		Name:    m.Name,
		Product: m.Product,
		Value:   m.Value,
	}
}

func (ddr *metadataDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	tflog.Info(ctx, "createApiCall")
	reqBody := metadataToRequest(ddr.Meta)
	apiResp, err := client.MetadataCreateWithResponse(ctx, reqBody)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON201 != nil {
		ddr.Meta = *apiResp.JSON201
	}

	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *metadataDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "readApiCall")
	apiResp, err := client.MetadataRetrieveWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON200 != nil {
		ddr.Meta = *apiResp.JSON200
	}

	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *metadataDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "updateApiCall")
	reqBody := metadataToRequest(ddr.Meta)
	apiResp, err := client.MetadataUpdateWithResponse(ctx, idNumber, reqBody)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON200 != nil {
		ddr.Meta = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *metadataDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "deleteApiCall")
	apiResp, err := client.MetadataDestroyWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	return apiResp.StatusCode(), apiResp.Body, nil
}

// clearFieldsApiCall sends the explicit-null PATCH that clears attributes
// removed from configuration. See clear.go: omitting a field from an update
// request leaves it unchanged, so clearing needs its own request.
func (ddr *metadataDefectdojoResource) clearFieldsApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int, body []byte) (int, []byte, error) {
	tflog.Info(ctx, "metadataDefectdojoResource clearFieldsApiCall")
	apiResp, err := client.MetadataPartialUpdateWithBodyWithResponse(ctx, idNumber, "application/json", bytes.NewReader(body))
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	return apiResp.StatusCode(), apiResp.Body, nil
}

type metadataResource struct {
	terraformResource
}

var _ resource.Resource = &metadataResource{}
var _ resource.ResourceWithImportState = &metadataResource{}
var _ resource.ResourceWithConfigValidators = &metadataResource{}

func NewMetadataResource() resource.Resource {
	return &metadataResource{
		terraformResource: terraformResource{
			typeName:     "defectdojo_metadata",
			dataProvider: metadataDataProvider{},
		},
	}
}

func (r metadataResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_metadata"
}

// ConfigValidators enforces that exactly one parent object (product or
// finding) is set, since DefectDojo metadata must attach to
// exactly one parent.
func (r metadataResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("product"), path.MatchRoot("finding"),
		),
	}
}

type metadataDataProvider struct{}

func (r metadataDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data metadataResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *metadataResourceData) id() types.String {
	return d.Id
}

func (d *metadataResourceData) setId(v types.String) { d.Id = v }

func (r metadataDataProvider) nameFromData(data terraformResourceData) (string, bool) {
	d := data.(*metadataResourceData)
	if !d.Name.IsNull() && !d.Name.IsUnknown() {
		return d.Name.ValueString(), true
	}
	return "", false
}

func (r metadataDataProvider) listByName(ctx context.Context, client *dd.ClientWithResponses, name string, data terraformResourceData) error {
	apiResp, err := client.MetadataListWithResponse(ctx, &dd.MetadataListParams{
		Name: &name,
	})
	if err != nil {
		return fmt.Errorf("error listing metadata: %w", err)
	}
	if apiResp.StatusCode() != 200 || apiResp.JSON200 == nil {
		return fmt.Errorf("unexpected API response: status %d, body: %s", apiResp.StatusCode(), string(apiResp.Body))
	}
	if apiResp.JSON200.Count == 0 {
		return fmt.Errorf("no metadata found with name %q", name)
	}
	if apiResp.JSON200.Count > 1 {
		return fmt.Errorf("more than one metadata found with name %q: metadata names are only unique per parent object; look up by id instead", name)
	}
	result := apiResp.JSON200.Results[0]
	if result.Id != nil {
		data.setId(types.StringValue(fmt.Sprintf("%d", *result.Id)))
	}
	return nil
}

func (d *metadataResourceData) defectdojoResource() defectdojoResource {
	return &metadataDefectdojoResource{
		Meta: dd.Meta{},
	}
}
