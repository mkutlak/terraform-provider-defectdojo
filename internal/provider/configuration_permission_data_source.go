package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
)

type configurationPermissionDataSource struct {
	terraformDatasource
}

func (t configurationPermissionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for a DefectDojo Configuration Permission. Configuration Permissions are DefectDojo's built-in permission registry, useful for RBAC wiring (e.g. authorized_users on products) without hardcoding ids. You can specify either the `id` or the `codename` to look up the Configuration Permission.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"codename": schema.StringAttribute{
				MarkdownDescription: "The machine-stable codename of the Configuration Permission (e.g. `add_product`). Specify either id or codename.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The human-readable name of the Configuration Permission.",
				Computed:            true,
			},
		},
	}
}

func (d configurationPermissionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_configuration_permission"
}

var _ datasource.DataSource = &configurationPermissionDataSource{}

func NewConfigurationPermissionDataSource() datasource.DataSource {
	return &configurationPermissionDataSource{
		terraformDatasource: terraformDatasource{dataProvider: configurationPermissionDataProvider{}},
	}
}

type configurationPermissionResourceData struct {
	Codename types.String `tfsdk:"codename" ddField:"Codename"`
	Name     types.String `tfsdk:"name" ddField:"Name"`
	Id       types.String `tfsdk:"id" ddField:"Id"`
}

// configurationPermissionDefectdojoResource wraps dd.ConfigurationPermission,
// DefectDojo's built-in permission registry.
type configurationPermissionDefectdojoResource struct {
	dd.ConfigurationPermission
}

// errConfigurationPermissionReadOnly is returned by the write stubs below.
// Configuration permissions are defined by DefectDojo itself; the methods
// only exist to satisfy the defectdojoResource interface for the data
// source read path.
var errConfigurationPermissionReadOnly = errors.New("configuration permissions are defined by DefectDojo itself; read-only")

func (ddr *configurationPermissionDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	return 0, nil, errConfigurationPermissionReadOnly
}

func (ddr *configurationPermissionDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.ConfigurationPermissionsRetrieveWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.ConfigurationPermission = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *configurationPermissionDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	return 0, nil, errConfigurationPermissionReadOnly
}

func (ddr *configurationPermissionDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	return 0, nil, errConfigurationPermissionReadOnly
}

type configurationPermissionDataProvider struct{}

func (r configurationPermissionDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data configurationPermissionResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *configurationPermissionResourceData) id() types.String { return d.Id }

func (d *configurationPermissionResourceData) setId(v types.String) { d.Id = v }

func (d *configurationPermissionResourceData) defectdojoResource() defectdojoResource {
	return &configurationPermissionDefectdojoResource{ConfigurationPermission: dd.ConfigurationPermission{}}
}

func (r configurationPermissionDataProvider) nameFromData(data terraformResourceData) (string, bool) {
	d := data.(*configurationPermissionResourceData)
	if !d.Codename.IsNull() && !d.Codename.IsUnknown() {
		return d.Codename.ValueString(), true
	}
	return "", false
}

func (r configurationPermissionDataProvider) listByName(ctx context.Context, client *dd.ClientWithResponses, name string, data terraformResourceData) error {
	apiResp, err := client.ConfigurationPermissionsListWithResponse(ctx, &dd.ConfigurationPermissionsListParams{
		Codename: &name,
	})
	if err != nil {
		return fmt.Errorf("error listing configuration permissions: %w", err)
	}
	if apiResp.StatusCode() != 200 || apiResp.JSON200 == nil {
		return fmt.Errorf("unexpected API response: status %d, body: %s", apiResp.StatusCode(), string(apiResp.Body))
	}
	var matched []dd.ConfigurationPermission
	for _, cp := range apiResp.JSON200.Results {
		if cp.Codename == name {
			matched = append(matched, cp)
		}
	}
	if len(matched) == 0 {
		return fmt.Errorf("no configuration permission found with codename %q", name)
	}
	if len(matched) > 1 {
		return fmt.Errorf("%d configuration permissions matched codename %q, expected exactly 1", len(matched), name)
	}
	if matched[0].Id != nil {
		data.setId(types.StringValue(fmt.Sprintf("%d", *matched[0].Id)))
	}
	return nil
}
