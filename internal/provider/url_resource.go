package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
)

func (t urlResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DefectDojo URL",

		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "The host of the URL, which can be a domain name or an IP address",
				Required:            true,
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "The protocol of the URL (e.g., http, https, ftp, etc.)",
				Optional:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "The port number of the URL (optional)",
				Optional:            true,
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "The path of the URL (optional)",
				Optional:            true,
			},
			"query": schema.StringAttribute{
				MarkdownDescription: "The query string of the URL (optional)",
				Optional:            true,
			},
			"fragment": schema.StringAttribute{
				MarkdownDescription: "The fragment identifier of the URL (optional)",
				Optional:            true,
			},
			"user_info": schema.StringAttribute{
				MarkdownDescription: "Connection details for a given user",
				Optional:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "Tags to apply to the URL",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"host_validation_failure": schema.BoolAttribute{
				MarkdownDescription: "Dictates whether the endpoint was found to have host validation issues during creation",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of the URL",
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

type urlResourceData struct {
	Id                    types.String `tfsdk:"id" ddField:"Id"`
	Host                  types.String `tfsdk:"host" ddField:"Host"`
	Protocol              types.String `tfsdk:"protocol" ddField:"Protocol"`
	Port                  types.Int64  `tfsdk:"port" ddField:"Port"`
	Path                  types.String `tfsdk:"path" ddField:"Path"`
	Query                 types.String `tfsdk:"query" ddField:"Query"`
	Fragment              types.String `tfsdk:"fragment" ddField:"Fragment"`
	UserInfo              types.String `tfsdk:"user_info" ddField:"UserInfo"`
	Tags                  types.Set    `tfsdk:"tags" ddField:"Tags"`
	HostValidationFailure types.Bool   `tfsdk:"host_validation_failure" ddField:"HostValidationFailure"`
	Type                  types.String `tfsdk:"type" ddField:"Type"`
}

type urlDefectdojoResource struct {
	dd.URL
}

// urlToRequest converts a URL (response model) to a URLRequest (request model).
// HostValidationFailure is intentionally omitted: it is server-managed and not writable.
func urlToRequest(u dd.URL) dd.URLRequest {
	return dd.URLRequest{
		Fragment: u.Fragment,
		Host:     u.Host,
		Path:     u.Path,
		Port:     u.Port,
		Protocol: u.Protocol,
		Query:    u.Query,
		Tags:     u.Tags,
		UserInfo: u.UserInfo,
	}
}

func (ddr *urlDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	tflog.Info(ctx, "createApiCall")
	reqBody := urlToRequest(ddr.URL)
	apiResp, err := client.UrlCreateWithResponse(ctx, reqBody)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON201 != nil {
		ddr.URL = *apiResp.JSON201
	}

	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *urlDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "readApiCall")
	apiResp, err := client.UrlRetrieveWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON200 != nil {
		ddr.URL = *apiResp.JSON200
	}

	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *urlDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "updateApiCall")
	reqBody := urlToRequest(ddr.URL)
	apiResp, err := client.UrlUpdateWithResponse(ctx, idNumber, reqBody)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	if apiResp.JSON200 != nil {
		ddr.URL = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *urlDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	tflog.Info(ctx, "deleteApiCall")
	apiResp, err := client.UrlDestroyWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", apiResp.Status(), apiResp.Body))
	return apiResp.StatusCode(), apiResp.Body, nil
}

type urlResource struct {
	terraformResource
}

var _ resource.Resource = &urlResource{}
var _ resource.ResourceWithImportState = &urlResource{}

func NewUrlResource() resource.Resource {
	return &urlResource{
		terraformResource: terraformResource{
			typeName:     "defectdojo_url",
			dataProvider: urlDataProvider{},
		},
	}
}

func (r urlResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_url"
}

type urlDataProvider struct{}

func (r urlDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data urlResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *urlResourceData) id() types.String { return d.Id }

func (d *urlResourceData) setId(v types.String) { d.Id = v }

func (r urlDataProvider) nameFromData(data terraformResourceData) (string, bool) {
	d := data.(*urlResourceData)
	if !d.Host.IsNull() && !d.Host.IsUnknown() {
		return d.Host.ValueString(), true
	}
	return "", false
}

func (r urlDataProvider) listByName(ctx context.Context, client *dd.ClientWithResponses, name string, data terraformResourceData) error {
	apiResp, err := client.UrlListWithResponse(ctx, &dd.UrlListParams{
		Host: &name,
	})
	if err != nil {
		return fmt.Errorf("error listing urls: %w", err)
	}
	if apiResp.StatusCode() != 200 || apiResp.JSON200 == nil {
		return fmt.Errorf("unexpected API response: status %d, body: %s", apiResp.StatusCode(), string(apiResp.Body))
	}
	if apiResp.JSON200.Count == 0 {
		return fmt.Errorf("no url found with host %q", name)
	}
	if apiResp.JSON200.Count > 1 {
		return fmt.Errorf("more than one url found with host %q, expected exactly 1", name)
	}
	result := apiResp.JSON200.Results[0]
	if result.Id != nil {
		data.setId(types.StringValue(fmt.Sprintf("%d", *result.Id)))
	}
	return nil
}

func (d *urlResourceData) defectdojoResource() defectdojoResource {
	return &urlDefectdojoResource{
		URL: dd.URL{},
	}
}
