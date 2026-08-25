package provider

import (
	"context"
	"fmt"
	"strings"

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

func (r cicdInfrastructureResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo CI/CD Infrastructure. DefectDojo 3.2 added this resource to replace Tool Configuration for CI/CD purposes. DefectDojo hides `description` and `url` from a token without the `view_cicdinfrastructure` or `change_cicdinfrastructure` permission; both attributes then read as empty.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the CI/CD infrastructure. A name is unique together with `infrastructure_type`, not on its own.",
				Required:            true,
			},
			// infrastructure_type carries RequiresReplace() because DefectDojo
			// rejects a change to this field once set. Verified on 3.2.300: a PUT
			// with a different infrastructure_type returns 200 but the server
			// silently keeps the original value, so an in-place update through
			// this engine would leave Terraform state permanently diverged from
			// the server without RequiresReplace forcing a destroy/create
			// instead.
			"infrastructure_type": schema.StringAttribute{
				MarkdownDescription: "The kind of CI/CD infrastructure. Valid values: `scm_server`, `build_server`, `orchestration`. Changing this value replaces the resource: DefectDojo does not allow this field to change once set.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("scm_server", "build_server", "orchestration"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// description and url are Optional+Computed rather than bare Optional.
			// DefectDojo declares both blank=True, default="", so the server
			// always answers with a concrete string ("" when unset); a bare
			// Optional attribute would plan as null and every apply that omits it
			// would fail with "Provider produced inconsistent result after
			// apply". See regulation_resource.go:35-44 for the identical shape.
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the CI/CD infrastructure. DefectDojo hides this from a token without the `view_cicdinfrastructure` or `change_cicdinfrastructure` permission.",
				Optional:            true,
				Computed:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "The public URL of the CI/CD infrastructure, for example `https://jenkins.company.com`. DefectDojo hides this from a token without the `view_cicdinfrastructure` or `change_cicdinfrastructure` permission.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

type cicdInfrastructureResourceData struct {
	Id                 types.String `tfsdk:"id" ddField:"Id"`
	Name               types.String `tfsdk:"name" ddField:"Name"`
	InfrastructureType types.String `tfsdk:"infrastructure_type" ddField:"InfrastructureType"`
	Description        types.String `tfsdk:"description" ddField:"Description"`
	Url                types.String `tfsdk:"url" ddField:"Url"`
}

type cicdInfrastructureDefectdojoResource struct {
	dd.CICDInfrastructure
}

func cicdInfrastructureToRequest(obj dd.CICDInfrastructure) dd.CICDInfrastructureRequest {
	return dd.CICDInfrastructureRequest{
		Name:               obj.Name,
		InfrastructureType: dd.CICDInfrastructureRequestInfrastructureType(obj.InfrastructureType),
		Description:        obj.Description,
		// Url is an oapi-codegen oneOf-string wrapper (see isOapiUnionStringType
		// in resource.go). CICDInfrastructure and CICDInfrastructureRequest use
		// distinct named wrapper types with identical underlying types, so a
		// direct field copy does not compile; convert explicitly.
		Url: (*dd.CICDInfrastructureRequest_Url)(obj.Url),
	}
}

func (ddr *cicdInfrastructureDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := cicdInfrastructureToRequest(ddr.CICDInfrastructure)
	apiResp, err := client.CicdInfrastructureCreateWithResponse(ctx, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON201 != nil {
		ddr.CICDInfrastructure = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *cicdInfrastructureDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.CicdInfrastructureRetrieveWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.CICDInfrastructure = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *cicdInfrastructureDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := cicdInfrastructureToRequest(ddr.CICDInfrastructure)
	apiResp, err := client.CicdInfrastructureUpdateWithResponse(ctx, idNumber, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.CICDInfrastructure = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *cicdInfrastructureDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.CicdInfrastructureDestroyWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

type cicdInfrastructureResource struct {
	terraformResource
}

var _ resource.Resource = &cicdInfrastructureResource{}
var _ resource.ResourceWithImportState = &cicdInfrastructureResource{}

func NewCicdInfrastructureResource() resource.Resource {
	return &cicdInfrastructureResource{
		terraformResource: terraformResource{typeName: "defectdojo_cicd_infrastructure", dataProvider: cicdInfrastructureDataProvider{}},
	}
}

func (r cicdInfrastructureResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cicd_infrastructure"
}

type cicdInfrastructureDataProvider struct{}

func (r cicdInfrastructureDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data cicdInfrastructureResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *cicdInfrastructureResourceData) id() types.String     { return d.Id }
func (d *cicdInfrastructureResourceData) setId(v types.String) { d.Id = v }

func (r cicdInfrastructureDataProvider) nameFromData(data terraformResourceData) (string, bool) {
	d := data.(*cicdInfrastructureResourceData)
	if !d.Name.IsNull() && !d.Name.IsUnknown() {
		return d.Name.ValueString(), true
	}
	return "", false
}

// listByName looks up a CI/CD infrastructure by name. Because uniqueness is on
// (name, infrastructure_type), a name alone can match up to three rows - one
// per infrastructure_type. This does not add an infrastructure_type filter
// argument: nameFilterable (datasource.go) is a single-field interface, and
// widening it would ripple through every other data source that implements it.
func (r cicdInfrastructureDataProvider) listByName(ctx context.Context, client *dd.ClientWithResponses, name string, data terraformResourceData) error {
	apiResp, err := client.CicdInfrastructureListWithResponse(ctx, &dd.CicdInfrastructureListParams{
		Name: &name,
	})
	if err != nil {
		return fmt.Errorf("error listing cicd infrastructure: %w", err)
	}
	if apiResp.StatusCode() != 200 || apiResp.JSON200 == nil {
		return fmt.Errorf("unexpected API response: status %d, body: %s", apiResp.StatusCode(), string(apiResp.Body))
	}
	var matched []dd.CICDInfrastructure
	for _, ci := range apiResp.JSON200.Results {
		if strings.EqualFold(ci.Name, name) {
			matched = append(matched, ci)
		}
	}
	if len(matched) == 0 {
		return fmt.Errorf("no cicd infrastructure found with name %q", name)
	}
	if len(matched) > 1 {
		return fmt.Errorf("%d cicd infrastructure entries matched name %q: names are only unique per infrastructure_type; look up by id instead", len(matched), name)
	}
	if matched[0].Id != nil {
		data.setId(types.StringValue(fmt.Sprintf("%d", *matched[0].Id)))
	}
	return nil
}

func (d *cicdInfrastructureResourceData) defectdojoResource() defectdojoResource {
	return &cicdInfrastructureDefectdojoResource{CICDInfrastructure: dd.CICDInfrastructure{}}
}
