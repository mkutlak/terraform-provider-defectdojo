package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
)

func (t productResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Product",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Product",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the Product",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`\A[^\s].*[^\s]\z`), "The description must not have leading or trailing whitespace"),
				},
			},
			"prod_numeric_grade": schema.Int64Attribute{
				MarkdownDescription: "The Numeric Grade of the Product",
				Optional:            true,
			},
			"business_criticality": schema.StringAttribute{
				MarkdownDescription: "The Business Criticality of the Product. Valid values are: 'very high', 'high', 'medium', 'low', 'very low', 'none'",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("very high", "high", "medium", "low", "very low", "none", ""),
				},
			},
			"platform": schema.StringAttribute{
				MarkdownDescription: "The Platform of the Product. Valid values are: 'web service', 'desktop', 'iot', 'mobile', 'web'",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("web service", "desktop", "iot", "mobile", "web", ""),
				},
			},
			"life_cycle": schema.StringAttribute{
				MarkdownDescription: "The Lifecycle state of the Product. Valid values are: 'construction', 'production', 'retirement'",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("construction", "production", "retirement", ""),
				},
			},
			"origin": schema.StringAttribute{
				MarkdownDescription: "The Origin of the Product. Valid values are: 'third party library', 'purchased', 'contractor', 'internal', 'open source', 'outsourced'",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("third party library", "purchased", "contractor", "internal", "open source", "outsourced", ""),
				},
			},
			"user_records": schema.Int64Attribute{
				MarkdownDescription: "Estimate the number of user records within the application.",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"revenue": schema.StringAttribute{
				MarkdownDescription: "Estimate the application's revenue.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`\A-?\d{0,13}(?:\.\d{0,2})?\z`), `Must be a decimal number format, i.e. /^-?\d{0,13}(?:\.\d{0,2})?$/`),
				},
			},
			"external_audience": schema.BoolAttribute{
				MarkdownDescription: "Specify if the application is used by people outside the organization.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"internet_accessible": schema.BoolAttribute{
				MarkdownDescription: "Specify if the application is accessible from the public internet.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"enable_skip_risk_acceptance": schema.BoolAttribute{
				MarkdownDescription: "Allows simple risk acceptance by checking/unchecking a checkbox.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"enable_full_risk_acceptance": schema.BoolAttribute{
				MarkdownDescription: "Allows full risk acceptance using a risk acceptance form, expiration date, uploaded proof, etc.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"product_manager_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the user who is the PM for this product.",
				Optional:            true,
			},
			"technical_contact_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the user who is the technical contact for this product.",
				Optional:            true,
			},
			"team_manager_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the user who is the manager for this product.",
				Optional:            true,
			},
			"regulation_ids": schema.SetAttribute{
				MarkdownDescription: "The IDs of the Regulations which apply to this product.",
				Optional:            true,
				ElementType:         types.Int64Type,
			},
			"product_type_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Product Type",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "Tags to apply to the product",
				Optional:            true,
				ElementType:         types.StringType,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(regexp.MustCompile(`\A[a-z0-9][a-z0-9_-]*\z`), "Tags must be lower case values (letters, digits, hyphens, underscores)"),
					),
				},
			},
			"disable_sla_breach_notifications": schema.BoolAttribute{
				MarkdownDescription: "Disable SLA breach notifications if configured in the global settings.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"enable_product_tag_inheritance": schema.BoolAttribute{
				MarkdownDescription: "Enables product tag inheritance. Any tags added on a product will automatically be added to all Engagements, Tests, and Findings.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"sla_configuration": schema.Int64Attribute{
				MarkdownDescription: "The ID of the SLA configuration to apply to this product.",
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

type productResourceData struct {
	Name                          types.String `tfsdk:"name" ddField:"Name"`
	Description                   types.String `tfsdk:"description" ddField:"Description"`
	ProductTypeId                 types.Int64  `tfsdk:"product_type_id" ddField:"ProdType"`
	Id                            types.String `tfsdk:"id" ddField:"Id"`
	BusinessCriticality           types.String `tfsdk:"business_criticality" ddField:"BusinessCriticality"`
	EnableFullRiskAcceptance      types.Bool   `tfsdk:"enable_full_risk_acceptance" ddField:"EnableFullRiskAcceptance"`
	EnableSimpleRiskAcceptance    types.Bool   `tfsdk:"enable_skip_risk_acceptance" ddField:"EnableSimpleRiskAcceptance"`
	ExternalAudience              types.Bool   `tfsdk:"external_audience" ddField:"ExternalAudience"`
	InternetAccessible            types.Bool   `tfsdk:"internet_accessible" ddField:"InternetAccessible"`
	Lifecycle                     types.String `tfsdk:"life_cycle" ddField:"Lifecycle"`
	Origin                        types.String `tfsdk:"origin" ddField:"Origin"`
	Platform                      types.String `tfsdk:"platform" ddField:"Platform"`
	ProdNumericGrade              types.Int64  `tfsdk:"prod_numeric_grade" ddField:"ProdNumericGrade"`
	ProductManagerId              types.Int64  `tfsdk:"product_manager_id" ddField:"ProductManager"`
	RegulationIds                 types.Set    `tfsdk:"regulation_ids" ddField:"Regulations"`
	Revenue                       types.String `tfsdk:"revenue" ddField:"Revenue"`
	Tags                          types.Set    `tfsdk:"tags" ddField:"Tags"`
	TeamManagerId                 types.Int64  `tfsdk:"team_manager_id" ddField:"TeamManager"`
	TechnicalContactId            types.Int64  `tfsdk:"technical_contact_id" ddField:"TechnicalContact"`
	UserRecords                   types.Int64  `tfsdk:"user_records" ddField:"UserRecords"`
	DisableSlaBreachNotifications types.Bool   `tfsdk:"disable_sla_breach_notifications" ddField:"DisableSlaBreachNotifications"`
	EnableProductTagInheritance   types.Bool   `tfsdk:"enable_product_tag_inheritance" ddField:"EnableProductTagInheritance"`
	SlaConfiguration              types.Int64  `tfsdk:"sla_configuration" ddField:"SlaConfiguration"`
}

type productDefectdojoResource struct {
	dd.Product
}

// productToRequest converts a Product (response model) to a ProductRequest (request model).
// The new API (v2.54.3) uses separate Request/Response schemas.
func productToRequest(p dd.Product) dd.ProductRequest {
	req := dd.ProductRequest{
		Name:                          p.Name,
		Description:                   p.Description,
		DisableSlaBreachNotifications: p.DisableSlaBreachNotifications,
		EnableFullRiskAcceptance:      p.EnableFullRiskAcceptance,
		EnableProductTagInheritance:   p.EnableProductTagInheritance,
		EnableSimpleRiskAcceptance:    p.EnableSimpleRiskAcceptance,
		ExternalAudience:              p.ExternalAudience,
		InternetAccessible:            p.InternetAccessible,
		ProdNumericGrade:              p.ProdNumericGrade,
		ProdType:                      p.ProdType,
		ProductManager:                p.ProductManager,
		Regulations:                   p.Regulations,
		Revenue:                       p.Revenue,
		Tags:                          p.Tags,
		TeamManager:                   p.TeamManager,
		TechnicalContact:              p.TechnicalContact,
		UserRecords:                   p.UserRecords,
	}
	// Only include SlaConfiguration when it has a meaningful value (non-nil and non-zero).
	// Sending pk=0 causes an API validation error: "Invalid pk \"0\" - object does not exist."
	if p.SlaConfiguration != nil && *p.SlaConfiguration != 0 {
		req.SlaConfiguration = p.SlaConfiguration
	}
	// Convert enum types: Product uses ProductXxx, ProductRequest uses ProductRequestXxx
	if p.BusinessCriticality != nil {
		v := dd.ProductRequestBusinessCriticality(*p.BusinessCriticality)
		req.BusinessCriticality = &v
	}
	if p.Lifecycle != nil {
		v := dd.ProductRequestLifecycle(*p.Lifecycle)
		req.Lifecycle = &v
	}
	if p.Origin != nil {
		v := dd.ProductRequestOrigin(*p.Origin)
		req.Origin = &v
	}
	if p.Platform != nil {
		v := dd.ProductRequestPlatform(*p.Platform)
		req.Platform = &v
	}
	return req
}

func (ddr *productDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	tflog.Info(ctx, "createApiCall")
	reqBody := productToRequest(ddr.Product)
	apiResp, err := client.ProductsCreateWithResponse(ctx, reqBody)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON201 != nil {
		ddr.Product = *apiResp.JSON201
	}

	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *productDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "readApiCall")
	apiResp, err := client.ProductsRetrieveWithResponse(ctx, idNumber, &dd.ProductsRetrieveParams{})
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON200 != nil {
		ddr.Product = *apiResp.JSON200
	}

	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *productDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "updateApiCall")
	reqBody := productToRequest(ddr.Product)
	apiResp, err := client.ProductsUpdateWithResponse(ctx, idNumber, reqBody)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON200 != nil {
		ddr.Product = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *productDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "deleteApiCall")
	apiResp, err := client.ProductsDestroyWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	return apiResp.StatusCode(), apiResp.Body, nil
}

type productResource struct {
	terraformResource
}

var _ resource.Resource = &productResource{}
var _ resource.ResourceWithImportState = &productResource{}

func NewProductResource() resource.Resource {
	return &productResource{
		terraformResource: terraformResource{
			typeName:     "defectdojo_product",
			dataProvider: productDataProvider{},
		},
	}
}

func (r productResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_product"
}

type productDataProvider struct{}

func (r productDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data productResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *productResourceData) id() types.String {
	return d.Id
}

func (d *productResourceData) defectdojoResource() defectdojoResource {
	return &productDefectdojoResource{
		Product: dd.Product{},
	}
}
