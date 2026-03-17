package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
)

func (t noteTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Note Type",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Note Type",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the Note Type",
				Required:            true,
			},
			"is_active": schema.BoolAttribute{
				MarkdownDescription: "Whether the Note Type is active",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"is_mandatory": schema.BoolAttribute{
				MarkdownDescription: "Whether the Note Type is mandatory",
				Optional:            true,
				Computed:            true,
			},
			"is_single": schema.BoolAttribute{
				MarkdownDescription: "Whether only a single note of this type is allowed",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
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

type noteTypeResourceData struct {
	Name        types.String `tfsdk:"name" ddField:"Name"`
	Description types.String `tfsdk:"description" ddField:"Description"`
	IsActive    types.Bool   `tfsdk:"is_active" ddField:"IsActive"`
	IsMandatory types.Bool   `tfsdk:"is_mandatory" ddField:"IsMandatory"`
	IsSingle    types.Bool   `tfsdk:"is_single" ddField:"IsSingle"`
	Id          types.String `tfsdk:"id" ddField:"Id"`
}

type noteTypeDefectdojoResource struct {
	dd.NoteType
}

func noteTypeToRequest(obj dd.NoteType) dd.NoteTypeRequest {
	return dd.NoteTypeRequest{
		Name:        obj.Name,
		Description: obj.Description,
		IsActive:    obj.IsActive,
		IsMandatory: obj.IsMandatory,
		IsSingle:    obj.IsSingle,
	}
}

func (ddr *noteTypeDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := noteTypeToRequest(ddr.NoteType)
	apiResp, err := client.NoteTypeCreateWithResponse(ctx, reqBody)
	if apiResp.JSON201 != nil {
		ddr.NoteType = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *noteTypeDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.NoteTypeRetrieveWithResponse(ctx, idNumber)
	if apiResp.JSON200 != nil {
		ddr.NoteType = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *noteTypeDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := noteTypeToRequest(ddr.NoteType)
	apiResp, err := client.NoteTypeUpdateWithResponse(ctx, idNumber, reqBody)
	if apiResp.JSON200 != nil {
		ddr.NoteType = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *noteTypeDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.NoteTypeDestroyWithResponse(ctx, idNumber)
	return apiResp.StatusCode(), apiResp.Body, err
}

type noteTypeResource struct {
	terraformResource
}

var _ resource.Resource = &noteTypeResource{}
var _ resource.ResourceWithImportState = &noteTypeResource{}

func NewNoteTypeResource() resource.Resource {
	return &noteTypeResource{
		terraformResource: terraformResource{dataProvider: noteTypeDataProvider{}},
	}
}

func (r noteTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_note_type"
}

type noteTypeDataProvider struct{}

func (r noteTypeDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data noteTypeResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *noteTypeResourceData) id() types.String { return d.Id }

func (d *noteTypeResourceData) defectdojoResource() defectdojoResource {
	return &noteTypeDefectdojoResource{NoteType: dd.NoteType{}}
}
