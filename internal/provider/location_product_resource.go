package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
)

func (t locationProductResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Location Product. Links a Location (e.g. a URL or Network Location) to a Product.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"location": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Location.",
				Required:            true,
			},
			"product": schema.Int64Attribute{
				MarkdownDescription: "The ID of the Product.",
				Required:            true,
			},
			// Not nullable in DefectDojo - the API enum carries "" as its blank
			// marker, so omitting this stores and returns an empty string. It
			// must be Computed to avoid a "provider produced inconsistent
			// result" error when left unset, matching `status` below.
			//
			// "" is listed in the OneOf because it is a real member of the enum
			// rather than a spelling of null. Verified on 3.1.101: a create
			// that omits the field answers 201 carrying "", POST and PUT with
			// an explicit "" answer 201 and 200 and a GET confirms the blank
			// stuck, while "bogus" is refused with `"bogus" is not a valid
			// choice.`. Leaving "" out made the one value a create produces the
			// one value a configuration could not contain, so a practitioner
			// matching config to state was rejected at plan time and, since
			// Computed attributes cannot be cleared by omission either, had no
			// way back to the blank the provider itself had created.
			//
			// `status` keeps a two-value OneOf because its enum really has no
			// blank member: the same request with "" is refused with `"" is not
			// a valid choice.`, and DefectDojo fills the field with "Mitigated"
			// when it is omitted.
			"relationship": schema.StringAttribute{
				MarkdownDescription: "The relationship between the location and the product. Valid values: 'owned_by', 'used_by', and '' - " +
					"the empty string is DefectDojo's own blank marker, and is what it stores when this attribute is omitted, so it is " +
					"accepted here to let a configuration match the state a create leaves behind. Deleting the attribute from an existing " +
					"configuration does not blank it: the value is computed, so Terraform reuses the one already in state and plans no " +
					"change. Set `relationship = \"\"` to clear it.",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("", "owned_by", "used_by"),
				},
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The status of the given Location. Valid values: 'Active', 'Mitigated'. DefectDojo stores 'Mitigated' when this " +
					"attribute is omitted. Deleting it from an existing configuration does not restore that default: the value is computed, so " +
					"Terraform reuses the one already in state and plans no change. Set the value explicitly to change it.",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("Active", "Mitigated"),
				},
			},
			"location_type": schema.StringAttribute{
				MarkdownDescription: "The type of the linked Location.",
				Computed:            true,
			},
			"location_value": schema.StringAttribute{
				MarkdownDescription: "The value representation of the linked Location.",
				Computed:            true,
			},
		},
	}
}

type locationProductResourceData struct {
	Id            types.String `tfsdk:"id" ddField:"Id"`
	Location      types.Int64  `tfsdk:"location" ddField:"Location"`
	Product       types.Int64  `tfsdk:"product" ddField:"Product"`
	Relationship  types.String `tfsdk:"relationship" ddField:"Relationship"`
	Status        types.String `tfsdk:"status" ddField:"Status"`
	LocationType  types.String `tfsdk:"location_type" ddField:"LocationType"`
	LocationValue types.String `tfsdk:"location_value" ddField:"LocationValue"`
	// Created/Updated (server-managed timestamps) and RelationshipData (an
	// interface{} field on dd.LocationProductReference with no fixed shape)
	// are intentionally not mapped here - see locationProductToRequest.
}

type locationProductDefectdojoResource struct {
	dd.LocationProductReference
}

func locationProductToRequest(l dd.LocationProductReference) dd.LocationProductReferenceRequest {
	req := dd.LocationProductReferenceRequest{
		Location: l.Location,
		Product:  l.Product,
	}
	if l.Relationship != nil {
		v := dd.LocationProductReferenceRequestRelationship(*l.Relationship)
		req.Relationship = &v
	}
	if l.Status != nil {
		v := dd.LocationProductReferenceRequestStatus(*l.Status)
		req.Status = &v
	}
	// RelationshipData is intentionally omitted: it is typed as interface{} on
	// dd.LocationProductReference(Request), and the reflection-based CRUD engine
	// (populateDefectdojoResource/populateResourceData in resource.go) has no
	// case for mapping an arbitrary interface{} to/from a Terraform type. There
	// is no corresponding schema attribute for it.
	return req
}

func (ddr *locationProductDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := locationProductToRequest(ddr.LocationProductReference)
	apiResp, err := client.LocationProductsCreateWithResponse(ctx, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON201 != nil {
		ddr.LocationProductReference = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *locationProductDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.LocationProductsRetrieveWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.LocationProductReference = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *locationProductDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := locationProductToRequest(ddr.LocationProductReference)
	apiResp, err := client.LocationProductsUpdateWithResponse(ctx, idNumber, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.LocationProductReference = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *locationProductDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.LocationProductsDestroyWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (d *locationProductResourceData) id() types.String {
	return d.Id
}

func (d *locationProductResourceData) setId(v types.String) { d.Id = v }

func (d *locationProductResourceData) defectdojoResource() defectdojoResource {
	return &locationProductDefectdojoResource{LocationProductReference: dd.LocationProductReference{}}
}

type locationProductResource struct {
	terraformResource
}

var _ resource.Resource = &locationProductResource{}
var _ resource.ResourceWithImportState = &locationProductResource{}

func NewLocationProductResource() resource.Resource {
	return &locationProductResource{
		terraformResource: terraformResource{
			typeName:     "defectdojo_location_product",
			dataProvider: locationProductDataProvider{},
		},
	}
}

func (r locationProductResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_location_product"
}

type locationProductDataProvider struct{}

func (r locationProductDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data locationProductResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}
