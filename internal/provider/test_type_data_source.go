package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
)

type testTypeDataSource struct {
	terraformDatasource
}

func (t testTypeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for DefectDojo Test Type. Test Types are registered by DefectDojo's scan-type integrations (e.g. Burp Suite, Trivy) and cannot be created or modified via this provider. Use this data source to reference a Test Type by name (e.g. in `defectdojo_test`) without hardcoding its numeric id. You can specify either `id` or `name` to look it up.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Test Type, e.g. \"Burp Scan\". Specify either id or name.",
				Optional:            true,
				Computed:            true,
			},
			"active": schema.BoolAttribute{
				MarkdownDescription: "Whether the Test Type is active.",
				Computed:            true,
			},
			"static_tool": schema.BoolAttribute{
				MarkdownDescription: "Whether the Test Type represents a static analysis tool.",
				Computed:            true,
			},
			"dynamic_tool": schema.BoolAttribute{
				MarkdownDescription: "Whether the Test Type represents a dynamic analysis tool.",
				Computed:            true,
			},
		},
	}
}

func (d testTypeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_test_type"
}

var _ datasource.DataSource = &testTypeDataSource{}

func NewTestTypeDataSource() datasource.DataSource {
	return &testTypeDataSource{
		terraformDatasource: terraformDatasource{dataProvider: testTypeDataProvider{}},
	}
}

type testTypeResourceData struct {
	Name        types.String `tfsdk:"name" ddField:"Name"`
	Active      types.Bool   `tfsdk:"active" ddField:"Active"`
	StaticTool  types.Bool   `tfsdk:"static_tool" ddField:"StaticTool"`
	DynamicTool types.Bool   `tfsdk:"dynamic_tool" ddField:"DynamicTool"`
	Id          types.String `tfsdk:"id" ddField:"Id"`
}

// testTypeDefectdojoResource wraps dd.TestType. Test Types are registered by
// DefectDojo's scan-type integrations; the write stubs below only exist to
// satisfy the defectdojoResource interface for the data source read path.
type testTypeDefectdojoResource struct {
	dd.TestType
}

// errTestTypeReadOnly is returned by the write stubs below. The test types
// API has no DELETE operation and its update body cannot change the name, so
// Terraform lifecycle semantics would be broken; this provider only exposes
// Test Type as a read-only data source.
var errTestTypeReadOnly = errors.New("test types are registered by DefectDojo scan-type integrations; read-only in this provider")

func (ddr *testTypeDefectdojoResource) createApiCall(ctx context.Context, client *dd.ClientWithResponses) (int, []byte, error) {
	return 0, nil, errTestTypeReadOnly
}

func (ddr *testTypeDefectdojoResource) readApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	apiResp, err := client.TestTypesRetrieveWithResponse(ctx, idNumber)
	if err != nil {
		return 0, nil, err
	}
	if apiResp.JSON200 != nil {
		ddr.TestType = *apiResp.JSON200
	}
	return apiResp.StatusCode(), apiResp.Body, nil
}

func (ddr *testTypeDefectdojoResource) updateApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	return 0, nil, errTestTypeReadOnly
}

func (ddr *testTypeDefectdojoResource) deleteApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int) (int, []byte, error) {
	return 0, nil, errTestTypeReadOnly
}

type testTypeDataProvider struct{}

func (r testTypeDataProvider) getData(ctx context.Context, getter dataGetter) (terraformResourceData, diag.Diagnostics) {
	var data testTypeResourceData
	diags := getter.Get(ctx, &data)
	return &data, diags
}

func (d *testTypeResourceData) id() types.String { return d.Id }

func (d *testTypeResourceData) setId(v types.String) { d.Id = v }

func (r testTypeDataProvider) nameFromData(data terraformResourceData) (string, bool) {
	d := data.(*testTypeResourceData)
	if !d.Name.IsNull() && !d.Name.IsUnknown() {
		return d.Name.ValueString(), true
	}
	return "", false
}

func (r testTypeDataProvider) listByName(ctx context.Context, client *dd.ClientWithResponses, name string, data terraformResourceData) error {
	apiResp, err := client.TestTypesListWithResponse(ctx, &dd.TestTypesListParams{
		Name: &name,
	})
	if err != nil {
		return fmt.Errorf("error listing test types: %w", err)
	}
	if apiResp.StatusCode() != 200 || apiResp.JSON200 == nil {
		return fmt.Errorf("unexpected API response: status %d, body: %s", apiResp.StatusCode(), string(apiResp.Body))
	}
	var matched []dd.TestType
	for _, tt := range apiResp.JSON200.Results {
		if tt.Name != nil && strings.EqualFold(*tt.Name, name) {
			matched = append(matched, tt)
		}
	}
	if len(matched) == 0 {
		return fmt.Errorf("no test type found with name %q", name)
	}
	if len(matched) > 1 {
		return fmt.Errorf("%d test types matched name %q, expected exactly 1", len(matched), name)
	}
	if matched[0].Id != nil {
		data.setId(types.StringValue(fmt.Sprintf("%d", *matched[0].Id)))
	}
	return nil
}

func (d *testTypeResourceData) defectdojoResource() defectdojoResource {
	return &testTypeDefectdojoResource{TestType: dd.TestType{}}
}
