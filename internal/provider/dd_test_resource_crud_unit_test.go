package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
	"gotest.tools/assert"
)

// The tests in this file drive terraformResource.Create/Read directly against
// an httptest-backed client, rather than exercising populateDefectdojoResource
// and populateResourceData in isolation like the other unit tests do.
//
// That distinction is the whole point. The engine used to collect conversion
// diagnostics into a local `diags` that Create then overwrote with
// `diags = resp.State.Set(...)`, so the errors never reached the practitioner
// and the zero value was transmitted anyway. Because that defect lived in
// Create/Read/Update and not in the populate functions, no unit test that calls
// the populate functions directly can catch a regression of it - only a test
// that goes through the framework entry points can. See issue #23.

// newDDTestResourceForTest wires a ddTestResource to client and returns it with
// the null-valued Plan/State skeleton the framework itself builds before
// invoking a resource (compare fwserver.Server.CreateResource, which does
// tftypes.NewValue(req.ResourceSchema.Type().TerraformType(ctx), nil)).
//
// The resource is returned as a pointer because Create/Update/Delete have value
// receivers but Read has a pointer receiver; only *ddTestResource has all four
// in its method set.
func newDDTestResourceForTest(t *testing.T, ctx context.Context, client *dd.ClientWithResponses) (*ddTestResource, tfsdk.Plan, tfsdk.State) {
	t.Helper()

	r := &ddTestResource{
		terraformResource: terraformResource{
			typeName:     "defectdojo_test",
			dataProvider: ddTestDataProvider{},
			client:       client,
		},
	}

	schemaResp := &resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, schemaResp)
	assert.Equal(t, schemaResp.Diagnostics.HasError(), false)

	nullRaw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), nil)

	return r,
		tfsdk.Plan{Schema: schemaResp.Schema, Raw: nullRaw},
		tfsdk.State{Schema: schemaResp.Schema, Raw: nullRaw}
}

// ddTestPlanData returns a fully-populated ddTestResourceData for Plan.Set.
// Every attribute needs a typed null - in particular Tags, whose zero value
// carries no element type and would be rejected against the schema's
// SetType{ElemType: StringType}.
func ddTestPlanData(id types.String, targetStart, targetEnd string) *ddTestResourceData {
	return &ddTestResourceData{
		Id:              id,
		TestType:        types.Int64Value(1),
		Engagement:      types.Int64Value(8),
		TargetStart:     types.StringValue(targetStart),
		TargetEnd:       types.StringValue(targetEnd),
		Title:           types.StringValue("SAST Scan"),
		Description:     types.StringNull(),
		Version:         types.StringNull(),
		BranchTag:       types.StringNull(),
		CommitHash:      types.StringNull(),
		BuildId:         types.StringNull(),
		PercentComplete: types.Int64Null(),
		Environment:     types.Int64Null(),
		Lead:            types.Int64Null(),
		ScanType:        types.StringNull(),
		Tags:            types.SetNull(types.StringType),
	}
}

// TestDDTestResourceCreateAcceptsDateOnlyTargets is the end-to-end regression
// test for issue #23: a date-only target_start must reach the API as a real
// instant, and must come back out of state as the literal the practitioner
// wrote. Before the fix the request carried 0001-01-01T00:00:00Z and state held
// null, which Terraform core rejected with "Provider produced inconsistent
// result after apply".
func TestDDTestResourceCreateAcceptsDateOnlyTargets(t *testing.T) {
	ctx := context.Background()

	var (
		mu      sync.Mutex
		gotBody dd.TestCreateRequest
		gotPath string
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		gotPath = req.URL.Path
		if err := json.NewDecoder(req.Body).Decode(&gotBody); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id := 25
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(dd.TestCreate{
			Id:          &id,
			TestType:    gotBody.TestType,
			Engagement:  gotBody.Engagement,
			TargetStart: gotBody.TargetStart,
			TargetEnd:   gotBody.TargetEnd,
			Title:       gotBody.Title,
		})
	}))
	defer srv.Close()

	client, err := dd.NewClientWithResponses(srv.URL)
	assert.NilError(t, err)

	r, plan, state := newDDTestResourceForTest(t, ctx, client)
	assert.Equal(t, plan.Set(ctx, ddTestPlanData(types.StringUnknown(), "2026-07-28", "2031-07-28")).HasError(), false)

	resp := &resource.CreateResponse{State: state}
	r.Create(ctx, resource.CreateRequest{Plan: plan}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics from Create: %v", resp.Diagnostics.Errors())
	}

	mu.Lock()
	defer mu.Unlock()

	assert.Equal(t, gotPath, "/api/v2/tests/")
	assert.Equal(t, gotBody.TargetStart.Format(time.RFC3339), "2026-07-28T00:00:00Z")
	assert.Equal(t, gotBody.TargetEnd.Format(time.RFC3339), "2031-07-28T00:00:00Z")

	var out ddTestResourceData
	assert.Equal(t, resp.State.Get(ctx, &out).HasError(), false)
	assert.Equal(t, out.Id.ValueString(), "25")
	assert.Equal(t, out.TargetStart.ValueString(), "2026-07-28")
	assert.Equal(t, out.TargetEnd.ValueString(), "2031-07-28")
}

// TestDDTestResourceCreateSurfacesConversionError covers the other half of the
// bug: an unparseable value must produce a diagnostic and must abort before any
// API call is made. Asserting the request count is the load-bearing part -
// previously the zero time was transmitted and the error was dropped on the
// floor.
func TestDDTestResourceCreateSurfacesConversionError(t *testing.T) {
	ctx := context.Background()

	var (
		mu       sync.Mutex
		requests int
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		mu.Lock()
		requests++
		mu.Unlock()
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	client, err := dd.NewClientWithResponses(srv.URL)
	assert.NilError(t, err)

	r, plan, state := newDDTestResourceForTest(t, ctx, client)
	assert.Equal(t, plan.Set(ctx, ddTestPlanData(types.StringUnknown(), "yesterday", "2031-07-28")).HasError(), false)

	resp := &resource.CreateResponse{State: state}
	r.Create(ctx, resource.CreateRequest{Plan: plan}, resp)

	assert.Equal(t, resp.Diagnostics.HasError(), true)

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, requests, 0)

	detail := resp.Diagnostics.Errors()[0].Detail()
	assert.Assert(t, strings.Contains(detail, "target_start"), "diagnostic should name the attribute, got: %s", detail)
	assert.Assert(t, strings.Contains(detail, "yesterday"), "diagnostic should quote the bad value, got: %s", detail)
}

// TestDDTestResourceReadPreservesDateOnlyLiteral proves the refresh path does
// not reintroduce drift: prior state holding "2026-07-28" must survive a Read
// whose API response renders the same instant as "2026-07-28T00:00:00Z". Without
// preservation this would produce a perpetual diff.
func TestDDTestResourceReadPreservesDateOnlyLiteral(t *testing.T) {
	ctx := context.Background()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.Path, "/api/v2/tests/25/")

		id, engagement, title := 25, 8, "SAST Scan"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(dd.Test{
			Id:          &id,
			TestType:    1,
			Engagement:  &engagement,
			Title:       &title,
			TargetStart: time.Date(2026, 7, 28, 0, 0, 0, 0, time.UTC),
			TargetEnd:   time.Date(2031, 7, 28, 0, 0, 0, 0, time.UTC),
		})
	}))
	defer srv.Close()

	client, err := dd.NewClientWithResponses(srv.URL)
	assert.NilError(t, err)

	r, _, state := newDDTestResourceForTest(t, ctx, client)
	assert.Equal(t, state.Set(ctx, ddTestPlanData(types.StringValue("25"), "2026-07-28", "2031-07-28")).HasError(), false)

	resp := &resource.ReadResponse{State: state}
	r.Read(ctx, resource.ReadRequest{State: state}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics from Read: %v", resp.Diagnostics.Errors())
	}

	var out ddTestResourceData
	assert.Equal(t, resp.State.Get(ctx, &out).HasError(), false)
	assert.Equal(t, out.TargetStart.ValueString(), "2026-07-28")
	assert.Equal(t, out.TargetEnd.ValueString(), "2031-07-28")
}
