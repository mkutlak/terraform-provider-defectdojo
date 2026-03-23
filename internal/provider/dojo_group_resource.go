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

func (t dojoGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo Dojo Group",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Dojo Group",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the Dojo Group",
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

type dojoGroupResourceData struct {
	Name        types.String `tfsdk:"name" ddField:"Name"`
	Description types.String `tfsdk:"description" ddField:"Description"`
	Id          types.String `tfsdk:"id" ddField:"Id"`
}

type dojoGroupDefectdojoResource struct {
	dd.DojoGroup
}

func dojoGroupToRequest(g dd.DojoGroup) dd.DojoGroupRequest {
	req := dd.DojoGroupRequest{
		Name:        g.Name,
		Description: g.Description,
	}
	if g.SocialProvider != nil {
		v := dd.DojoGroupRequestSocialProvider(*g.SocialProvider)
		req.SocialProvider = &v
	}
	return req
}

func (ddr *dojoGroupDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	reqBody := dojoGroupToRequest(ddr.DojoGroup)
	apiResp, err := client.DojoGroupsCreateWithResponse(ctx, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON201 != nil {
		ddr.DojoGroup = *apiResp.JSON201
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *dojoGroupDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.DojoGroupsRetrieveWithResponse(ctx, idNumber, &dd.DojoGroupsRetrieveParams{})
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.DojoGroup = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *dojoGroupDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	reqBody := dojoGroupToRequest(ddr.DojoGroup)
	apiResp, err := client.DojoGroupsUpdateWithResponse(ctx, idNumber, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.DojoGroup = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *dojoGroupDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.DojoGroupsDestroyWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

type dojoGroupResource struct {
	terraformResource
}

var _ resource.Resource = &dojoGroupResource{}
var _ resource.ResourceWithImportState = &dojoGroupResource{}

func NewDojoGroupResource() resource.Resource {
	return &dojoGroupResource{
		terraformResource: terraformResource{
			typeName:     "defectdojo_dojo_group",
			dataProvider: dojoGroupDataProvider{},
		},
	}
}

func (r dojoGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dojo_group"
}

type dojoGroupDataProvider struct{}

func (r dojoGroupDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data dojoGroupResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *dojoGroupResourceData) id() types.String {
	return d.Id
}

func (d *dojoGroupResourceData) defectdojoResource() defectdojoResource {
	return &dojoGroupDefectdojoResource{DojoGroup: dd.DojoGroup{}}
}
