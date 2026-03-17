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

func (t stubFindingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Stub Finding",
		Attributes: map[string]schema.Attribute{
			"title": schema.StringAttribute{
				MarkdownDescription: "The title of the Stub Finding",
				Required:            true,
			},
			"test": schema.Int64Attribute{
				MarkdownDescription: "The test this stub finding belongs to.",
				Required:            true,
			},
			"date": schema.StringAttribute{
				MarkdownDescription: "Date of the stub finding (format: 2006-01-02).",
				Optional:            true,
			},
			"severity": schema.StringAttribute{
				MarkdownDescription: "The severity of the stub finding.",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the stub finding.",
				Optional:            true,
			},
			"reporter": schema.Int64Attribute{
				MarkdownDescription: "The user who reported this stub finding.",
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

type stubFindingResourceData struct {
	Title       types.String `tfsdk:"title" ddField:"Title"`
	Test        types.Int64  `tfsdk:"test" ddField:"Test"`
	Date        types.String `tfsdk:"date" ddField:"Date"`
	Severity    types.String `tfsdk:"severity" ddField:"Severity"`
	Description types.String `tfsdk:"description" ddField:"Description"`
	Reporter    types.Int64  `tfsdk:"reporter" ddField:"Reporter"`
	Id          types.String `tfsdk:"id" ddField:"Id"`
}

// stubFindingDefectdojoResource wraps StubFindingCreate which is what the API returns on create/read.
// Create uses StubFindingCreateRequest, update uses StubFindingRequest.
type stubFindingDefectdojoResource struct {
	dd.StubFindingCreate
}

func stubFindingToCreateRequest(obj dd.StubFindingCreate) dd.StubFindingCreateRequest {
	return dd.StubFindingCreateRequest{
		Title:       obj.Title,
		Test:        obj.Test,
		Date:        obj.Date,
		Severity:    obj.Severity,
		Description: obj.Description,
	}
}

func stubFindingToUpdateRequest(obj dd.StubFindingCreate) dd.StubFindingRequest {
	return dd.StubFindingRequest{
		Title:       obj.Title,
		Date:        obj.Date,
		Severity:    obj.Severity,
		Description: obj.Description,
	}
}

func (ddr *stubFindingDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := stubFindingToCreateRequest(ddr.StubFindingCreate)
	apiResp, err := client.StubFindingsCreateWithResponse(ctx, reqBody)
	if apiResp.JSON201 != nil {
		ddr.StubFindingCreate = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *stubFindingDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.StubFindingsRetrieveWithResponse(ctx, idNumber)
	if apiResp.JSON200 != nil {
		// StubFinding (read) has Test as *int, StubFindingCreate has Test as int
		sf := apiResp.JSON200
		ddr.StubFindingCreate = dd.StubFindingCreate{
			Id:          sf.Id,
			Title:       sf.Title,
			Date:        sf.Date,
			Severity:    sf.Severity,
			Description: sf.Description,
			Reporter:    sf.Reporter,
		}
		if sf.Test != nil {
			ddr.StubFindingCreate.Test = *sf.Test
		}
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *stubFindingDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := stubFindingToUpdateRequest(ddr.StubFindingCreate)
	apiResp, err := client.StubFindingsUpdateWithResponse(ctx, idNumber, reqBody)
	if apiResp.JSON200 != nil {
		sf := apiResp.JSON200
		ddr.StubFindingCreate = dd.StubFindingCreate{
			Id:          sf.Id,
			Title:       sf.Title,
			Date:        sf.Date,
			Severity:    sf.Severity,
			Description: sf.Description,
			Reporter:    sf.Reporter,
		}
		if sf.Test != nil {
			ddr.StubFindingCreate.Test = *sf.Test
		}
	}
	return apiResp.StatusCode(), apiResp.Body, err
}

func (ddr *stubFindingDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.StubFindingsDestroyWithResponse(ctx, idNumber)
	return apiResp.StatusCode(), apiResp.Body, err
}

type stubFindingResource struct {
	terraformResource
}

var _ resource.Resource = &stubFindingResource{}
var _ resource.ResourceWithImportState = &stubFindingResource{}

func NewStubFindingResource() resource.Resource {
	return &stubFindingResource{
		terraformResource: terraformResource{dataProvider: stubFindingDataProvider{}},
	}
}

func (r stubFindingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stub_finding"
}

type stubFindingDataProvider struct{}

func (r stubFindingDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data stubFindingResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *stubFindingResourceData) id() types.String { return d.Id }

func (d *stubFindingResourceData) defectdojoResource() defectdojoResource {
	return &stubFindingDefectdojoResource{StubFindingCreate: dd.StubFindingCreate{}}
}
