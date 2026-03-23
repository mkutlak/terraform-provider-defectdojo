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

func (t languageTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Language Type",
		Attributes: map[string]schema.Attribute{
			"language": schema.StringAttribute{
				MarkdownDescription: "The name of the language",
				Required:            true,
			},
			"color": schema.StringAttribute{
				MarkdownDescription: "Color associated with the language",
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

type languageTypeResourceData struct {
	Language types.String `tfsdk:"language" ddField:"Language"`
	Color    types.String `tfsdk:"color" ddField:"Color"`
	Id       types.String `tfsdk:"id" ddField:"Id"`
}

type languageTypeDefectdojoResource struct {
	dd.LanguageType
}

func languageTypeToRequest(obj dd.LanguageType) dd.LanguageTypeRequest {
	return dd.LanguageTypeRequest{
		Language: obj.Language,
		Color:    obj.Color,
	}
}

func (ddr *languageTypeDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := languageTypeToRequest(ddr.LanguageType)
	apiResp, err := client.LanguageTypesCreateWithResponse(ctx, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON201 != nil {
		ddr.LanguageType = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *languageTypeDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.LanguageTypesRetrieveWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.LanguageType = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *languageTypeDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := languageTypeToRequest(ddr.LanguageType)
	apiResp, err := client.LanguageTypesUpdateWithResponse(ctx, idNumber, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.LanguageType = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *languageTypeDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.LanguageTypesDestroyWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

type languageTypeResource struct {
	terraformResource
}

var _ resource.Resource = &languageTypeResource{}
var _ resource.ResourceWithImportState = &languageTypeResource{}

func NewLanguageTypeResource() resource.Resource {
	return &languageTypeResource{
		terraformResource: terraformResource{typeName: "defectdojo_language_type", dataProvider: languageTypeDataProvider{}},
	}
}

func (r languageTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_language_type"
}

type languageTypeDataProvider struct{}

func (r languageTypeDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data languageTypeResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *languageTypeResourceData) id() types.String { return d.Id }

func (d *languageTypeResourceData) defectdojoResource() defectdojoResource {
	return &languageTypeDefectdojoResource{LanguageType: dd.LanguageType{}}
}
