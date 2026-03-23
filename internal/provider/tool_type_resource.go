package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
)

func (t toolTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Tool Type",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Tool Type",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the Tool Type",
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

type toolTypeResourceData struct {
	Name        types.String `tfsdk:"name" ddField:"Name"`
	Description types.String `tfsdk:"description" ddField:"Description"`
	Id          types.String `tfsdk:"id" ddField:"Id"`
}

type toolTypeDefectdojoResource struct {
	dd.ToolType
}

func toolTypeToRequest(obj dd.ToolType) dd.ToolTypeRequest {
	return dd.ToolTypeRequest{
		Name:        obj.Name,
		Description: obj.Description,
	}
}

func (ddr *toolTypeDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := toolTypeToRequest(ddr.ToolType)
	apiResp, err := client.ToolTypesCreateWithResponse(ctx, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON201 != nil {
		ddr.ToolType = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *toolTypeDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.ToolTypesRetrieveWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.ToolType = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *toolTypeDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := toolTypeToRequest(ddr.ToolType)
	apiResp, err := client.ToolTypesUpdateWithResponse(ctx, idNumber, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.ToolType = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *toolTypeDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.ToolTypesDestroyWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

type toolTypeResource struct {
	terraformResource
}

var _ resource.Resource = &toolTypeResource{}
var _ resource.ResourceWithImportState = &toolTypeResource{}

func NewToolTypeResource() resource.Resource {
	return &toolTypeResource{
		terraformResource: terraformResource{typeName: "defectdojo_tool_type", dataProvider: toolTypeDataProvider{}},
	}
}

func (r toolTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tool_type"
}

type toolTypeDataProvider struct{}

func (r toolTypeDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data toolTypeResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *toolTypeResourceData) id() types.String { return d.Id }
func (d *toolTypeResourceData) setId(v types.String) { d.Id = v }

func (r toolTypeDataProvider) nameFromData(data terraformResourceData) (string, bool) {
	d := data.(*toolTypeResourceData)
	if !d.Name.IsNull() && !d.Name.IsUnknown() {
		return d.Name.ValueString(), true
	}
	return "", false
}

func (r toolTypeDataProvider) listByName(ctx context.Context, client *dd.ClientWithResponses, name string, data terraformResourceData) error {
	apiResp, err := client.ToolTypesListWithResponse(ctx, &dd.ToolTypesListParams{
		Name: &name,
	})
	if err != nil {
		return fmt.Errorf("error listing tool types: %w", err)
	}
	if apiResp.StatusCode() != 200 || apiResp.JSON200 == nil {
		return fmt.Errorf("unexpected API response: status %d, body: %s", apiResp.StatusCode(), string(apiResp.Body))
	}
	var matched []dd.ToolType
	for _, tt := range apiResp.JSON200.Results {
		if strings.EqualFold(tt.Name, name) {
			matched = append(matched, tt)
		}
	}
	if len(matched) == 0 {
		return fmt.Errorf("no tool type found with name %q", name)
	}
	if len(matched) > 1 {
		return fmt.Errorf("%d tool types matched name %q, expected exactly 1", len(matched), name)
	}
	if matched[0].Id != nil {
		data.setId(types.StringValue(fmt.Sprintf("%d", *matched[0].Id)))
	}
	return nil
}

func (d *toolTypeResourceData) defectdojoResource() defectdojoResource {
	return &toolTypeDefectdojoResource{ToolType: dd.ToolType{}}
}
