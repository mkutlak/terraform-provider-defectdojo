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

func (t languageResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Language",
		Attributes: map[string]schema.Attribute{
			"language_type": schema.Int64Attribute{
				MarkdownDescription: "The language type ID.",
				Required:            true,
			},
			"product": schema.Int64Attribute{
				MarkdownDescription: "The product this language is associated with.",
				Required:            true,
			},
			"user": schema.Int64Attribute{
				MarkdownDescription: "The user who added this language.",
				Optional:            true,
			},
			"files": schema.Int64Attribute{
				MarkdownDescription: "Number of files.",
				Optional:            true,
			},
			"code": schema.Int64Attribute{
				MarkdownDescription: "Number of lines of code.",
				Optional:            true,
			},
			"blank": schema.Int64Attribute{
				MarkdownDescription: "Number of blank lines.",
				Optional:            true,
			},
			"comment": schema.Int64Attribute{
				MarkdownDescription: "Number of comment lines.",
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

type languageResourceData struct {
	Language types.Int64  `tfsdk:"language_type" ddField:"LanguageTypeId"`
	Product  types.Int64  `tfsdk:"product" ddField:"Product"`
	User     types.Int64  `tfsdk:"user" ddField:"User"`
	Files    types.Int64  `tfsdk:"files" ddField:"Files"`
	Code     types.Int64  `tfsdk:"code" ddField:"Code"`
	Blank    types.Int64  `tfsdk:"blank" ddField:"Blank"`
	Comment  types.Int64  `tfsdk:"comment" ddField:"Comment"`
	Id       types.String `tfsdk:"id" ddField:"Id"`
}

// ddLanguageWrapper is a thin struct wrapper that exposes dd.Language fields
// without the name collision that would occur from directly embedding dd.Language
// (which has a field named "Language" same as the embedding type name).
type ddLanguageWrapper struct {
	Id       *int       `json:"id,omitempty"`
	LanguageTypeId int  `json:"language"`
	Product  int        `json:"product"`
	User     *int       `json:"user,omitempty"`
	Files    *int       `json:"files,omitempty"`
	Code     *int       `json:"code,omitempty"`
	Blank    *int       `json:"blank,omitempty"`
	Comment  *int       `json:"comment,omitempty"`
}

func ddLanguageToWrapper(obj dd.Language) ddLanguageWrapper {
	return ddLanguageWrapper{
		Id:             obj.Id,
		LanguageTypeId: obj.Language,
		Product:        obj.Product,
		User:           obj.User,
		Files:          obj.Files,
		Code:           obj.Code,
		Blank:          obj.Blank,
		Comment:        obj.Comment,
	}
}

func wrapperToLanguage(w ddLanguageWrapper) dd.Language {
	return dd.Language{
		Id:       w.Id,
		Language: w.LanguageTypeId,
		Product:  w.Product,
		User:     w.User,
		Files:    w.Files,
		Code:     w.Code,
		Blank:    w.Blank,
		Comment:  w.Comment,
	}
}

type languageDefectdojoResource struct {
	ddLanguageWrapper
}

func languageToRequest(w ddLanguageWrapper) dd.LanguageRequest {
	return dd.LanguageRequest{
		Language: w.LanguageTypeId,
		Product:  w.Product,
		User:     w.User,
		Files:    w.Files,
		Code:     w.Code,
		Blank:    w.Blank,
		Comment:  w.Comment,
	}
}

func (ddr *languageDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := languageToRequest(ddr.ddLanguageWrapper)
	apiResp, err := client.LanguagesCreateWithResponse(ctx, reqBody)
	if apiResp.JSON201 != nil {
		ddr.ddLanguageWrapper = ddLanguageToWrapper(*apiResp.JSON201)
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *languageDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.LanguagesRetrieveWithResponse(ctx, idNumber, &dd.LanguagesRetrieveParams{})
	if apiResp.JSON200 != nil {
		ddr.ddLanguageWrapper = ddLanguageToWrapper(*apiResp.JSON200)
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *languageDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := languageToRequest(ddr.ddLanguageWrapper)
	apiResp, err := client.LanguagesUpdateWithResponse(ctx, idNumber, reqBody)
	if apiResp.JSON200 != nil {
		ddr.ddLanguageWrapper = ddLanguageToWrapper(*apiResp.JSON200)
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *languageDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.LanguagesDestroyWithResponse(ctx, idNumber)
	return apiResp.StatusCode(), apiResp.Body, err
}

type languageResource struct {
	terraformResource
}

var _ resource.Resource = &languageResource{}
var _ resource.ResourceWithImportState = &languageResource{}

func NewLanguageResource() resource.Resource {
	return &languageResource{
		terraformResource: terraformResource{dataProvider: languageDataProvider{}},
	}
}

func (r languageResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_language"
}

type languageDataProvider struct{}

func (r languageDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data languageResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *languageResourceData) id() types.String { return d.Id }

func (d *languageResourceData) defectdojoResource() defectdojoResource {
	return &languageDefectdojoResource{ddLanguageWrapper: ddLanguageWrapper{}}
}
