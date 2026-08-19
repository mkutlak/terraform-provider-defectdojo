package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
)

// TestResourcesWithClearableAttributesImplementNullClearer is the completeness
// gate for the clearing mechanism.
//
// Any attribute that is Optional without Computed and without a Default can be
// removed from configuration after having held a value. When that happens the
// engine must PATCH an explicit null, because omitting the field from an update
// request leaves it unchanged (GitHub issue #30). A resource with such an
// attribute but no clearFieldsApiCall would fail at apply time with a provider
// bug, so the omission is a build failure instead.
func TestResourcesWithClearableAttributesImplementNullClearer(t *testing.T) {
	t.Parallel()

	needed := 0
	for tfTypeName, resp := range providerResourceSchemas(t) {
		if resp.Diagnostics.HasError() {
			continue // reported by TestSchemasValidateImplementation
		}
		targets, _, ok := resourceDdTargets(t, tfTypeName)
		if !ok {
			continue
		}

		var clearable []string
		for attrName, attribute := range resp.Schema.Attributes {
			if !planningAsNull(attribute) {
				continue
			}
			if _, ok := targets[attrName]; !ok {
				continue // no ddField tag: the engine never writes it
			}
			clearable = append(clearable, attrName)
		}
		if len(clearable) == 0 {
			continue
		}
		needed++

		model, ok := testAccResourceModel(tfTypeName)
		if !ok {
			continue
		}
		if _, ok := model.defectdojoResource().(nullClearer); !ok {
			sort.Strings(clearable)
			t.Errorf("resource %s has %d attribute(s) that can be cleared (%v) but its "+
				"defectdojoResource() does not implement clearFieldsApiCall. Removing one of "+
				"them from configuration would fail the apply, because omitting a field from an "+
				"update request leaves it unchanged. Add the method - see clear.go for the "+
				"seven-line shape.", tfTypeName, len(clearable), clearable)
		}
	}

	if needed == 0 {
		t.Fatal("no resource was found to have clearable attributes; the join is broken and this " +
			"test would pass vacuously")
	}
	t.Logf("%d resources have clearable attributes and implement clearFieldsApiCall", needed)
}

// TestClearedDdFieldsDetectsRemovedAttributes drives the plan-vs-state diff over
// a real model. dd_test is the resource from issue #30.
func TestClearedDdFieldsDetectsRemovedAttributes(t *testing.T) {
	t.Parallel()

	plan := &ddTestResourceData{
		Id:              types.StringValue("1"),
		TestType:        types.Int64Value(1),
		Engagement:      types.Int64Value(2),
		TargetStart:     types.StringValue("2026-01-01T00:00:00Z"),
		TargetEnd:       types.StringValue("2026-01-02T00:00:00Z"),
		BranchTag:       types.StringNull(), // removed from config
		BuildId:         types.StringNull(), // removed from config
		CommitHash:      types.StringNull(), // removed from config
		PercentComplete: types.Int64Null(),  // removed from config
		Version:         types.StringValue("v2"),
		Title:           types.StringNull(), // never set: nothing to clear
		Tags:            types.SetNull(types.StringType),
	}
	state := &ddTestResourceData{
		Id:              types.StringValue("1"),
		TestType:        types.Int64Value(1),
		Engagement:      types.Int64Value(2),
		TargetStart:     types.StringValue("2026-01-01T00:00:00Z"),
		TargetEnd:       types.StringValue("2026-01-02T00:00:00Z"),
		BranchTag:       types.StringValue("main"),
		BuildId:         types.StringValue("b-1"),
		CommitHash:      types.StringValue("abc123"),
		PercentComplete: types.Int64Value(50),
		Version:         types.StringValue("v1"),
		Title:           types.StringNull(),
		Tags:            types.SetValueMust(types.StringType, []attr.Value{types.StringValue("a")}),
	}

	targets := clearedDdFields(plan, state, plan.defectdojoResource())
	var nulls, empties []string
	for _, t := range targets {
		if t.isSlice {
			empties = append(empties, t.jsonName)
		} else {
			nulls = append(nulls, t.jsonName)
		}
	}
	sort.Strings(nulls)
	sort.Strings(empties)

	wantNulls := []string{"branch_tag", "build_id", "commit_hash", "percent_complete"}
	if len(nulls) != len(wantNulls) {
		t.Fatalf("clearedDdFields nulls = %v, want %v", nulls, wantNulls)
	}
	for i := range wantNulls {
		if nulls[i] != wantNulls[i] {
			t.Errorf("clearedDdFields nulls = %v, want %v", nulls, wantNulls)
			break
		}
	}

	// A Set that held elements and is now null clears to [], not null: the spec
	// types tags as a non-nullable array.
	if len(empties) != 1 || empties[0] != "tags" {
		t.Errorf("clearedDdFields empties = %v, want [tags]", empties)
	}
}

// TestClearedDdFieldsIgnoresUnchangedAndUnknown guards the two ways this diff
// could do damage: clearing something the practitioner never removed, and
// clearing something the server is about to fill in.
func TestClearedDdFieldsIgnoresUnchangedAndUnknown(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name      string
		planValue types.String
	}{
		{"value unchanged", types.StringValue("main")},
		{"value changed", types.StringValue("release")},
		{"plan unknown (server will fill it)", types.StringUnknown()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			plan := &ddTestResourceData{BranchTag: tc.planValue}
			state := &ddTestResourceData{BranchTag: types.StringValue("main")}
			if got := clearedDdFields(plan, state, plan.defectdojoResource()); len(got) != 0 {
				t.Errorf("clearedDdFields = %v, want empty: nothing was removed from config", got)
			}
		})
	}

	t.Run("state already null", func(t *testing.T) {
		plan := &ddTestResourceData{BranchTag: types.StringNull()}
		state := &ddTestResourceData{BranchTag: types.StringNull()}
		if got := clearedDdFields(plan, state, plan.defectdojoResource()); len(got) != 0 {
			t.Errorf("clearedDdFields = %v, want empty: there was nothing to clear", got)
		}
	})
}

// TestClearedDdFieldsDetectsRemovedZeroValues is the state-side half of the
// same distinction stillSetAfterUpdate has to make.
//
// This diff reads Terraform values, not Go ones, and there null is its own
// sentinel: types.Int64Value(0), types.BoolValue(false) and
// types.StringValue("") are all non-null, so an attribute that held a zero and
// was then removed from configuration is a clear target like any other. It has
// to be, or removing branch_tag = "" would leave the empty string on the server
// and fail the apply.
func TestClearedDdFieldsDetectsRemovedZeroValues(t *testing.T) {
	t.Parallel()

	plan := &ddTestResourceData{
		BranchTag:       types.StringNull(),
		PercentComplete: types.Int64Null(),
		Tags:            types.SetNull(types.StringType),
	}
	state := &ddTestResourceData{
		BranchTag:       types.StringValue(""),
		PercentComplete: types.Int64Value(0),
		Tags:            types.SetValueMust(types.StringType, []attr.Value{}),
	}

	var got []string
	for _, target := range clearedDdFields(plan, state, plan.defectdojoResource()) {
		got = append(got, target.jsonName)
	}
	sort.Strings(got)

	want := []string{"branch_tag", "percent_complete", "tags"}
	if !slices.Equal(got, want) {
		t.Errorf("clearedDdFields = %v, want %v: a zero is a value, so removing it from "+
			"configuration still has to be cleared", got, want)
	}
}

// TestApplyClearedFieldsNamesConfigurationAttributes pins the vocabulary of the
// clearing diagnostics.
//
// The practitioner writes `product_manager_id` and `regulation_ids`; DefectDojo
// calls the same two fields `product_manager` and `regulations`. Reporting the
// wire names sent the reader looking for arguments this provider does not have:
//
//	Removing product_manager from the configuration requires DefectDojo to
//	accept an explicit null for those fields, but it answered 400.
//
// So the attribute list is spelled the way the configuration spells it. The
// echoed request body keeps the wire names, which is the point of dumping it.
func TestApplyClearedFieldsNamesConfigurationAttributes(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name         string
		plan, state  *productResourceData
		wantAttrList string
		wantWireName string
	}{
		{
			name:         "a scalar attribute",
			plan:         &productResourceData{ProductManagerId: types.Int64Null()},
			state:        &productResourceData{ProductManagerId: types.Int64Value(3)},
			wantAttrList: "Removing product_manager_id from the configuration",
			wantWireName: `"product_manager"`,
		},
		{
			name:  "a set attribute",
			plan:  &productResourceData{RegulationIds: types.SetNull(types.Int64Type)},
			state: &productResourceData{RegulationIds: types.SetValueMust(types.Int64Type, []attr.Value{types.Int64Value(1)})},

			wantAttrList: "Removing regulation_ids from the configuration",
			wantWireName: `"regulations"`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(`{"detail":"this field may not be null."}`))
			}))
			defer srv.Close()

			client, err := dd.NewClientWithResponses(srv.URL)
			if err != nil {
				t.Fatalf("NewClientWithResponses: %v", err)
			}

			ddResource := tc.plan.defectdojoResource()
			targets := clearedDdFields(tc.plan, tc.state, ddResource)
			if len(targets) != 1 {
				t.Fatalf("clearedDdFields = %v, want exactly one target", targets)
			}

			var diags diag.Diagnostics
			applyClearedFields(context.Background(), &diags, client, "defectdojo_product", 7, ddResource, targets)

			if !diags.HasError() {
				t.Fatal("applyClearedFields reported no error although the API answered 400")
			}
			detail := diags.Errors()[0].Detail()
			if !strings.Contains(detail, tc.wantAttrList) {
				t.Errorf("the diagnostic must name the argument the practitioner wrote.\n"+
					"want substring: %q\ngot:\n%s", tc.wantAttrList, detail)
			}
			if !strings.Contains(detail, tc.wantWireName) {
				t.Errorf("the echoed request body must stay as the server saw it.\n"+
					"want substring: %q\ngot:\n%s", tc.wantWireName, detail)
			}
		})
	}
}

func TestClearPatchBody(t *testing.T) {
	t.Parallel()

	raw, err := clearPatchBody([]clearTarget{
		{jsonName: "branch_tag"},
		{jsonName: "percent_complete"},
		{jsonName: "tags", isSlice: true},
	})
	if err != nil {
		t.Fatalf("clearPatchBody: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("clearPatchBody produced invalid JSON %q: %v", raw, err)
	}

	if v, ok := got["branch_tag"]; !ok || v != nil {
		t.Errorf("branch_tag = %v (present=%v), want explicit null", v, ok)
	}
	if v, ok := got["percent_complete"]; !ok || v != nil {
		t.Errorf("percent_complete = %v (present=%v), want explicit null", v, ok)
	}
	if v, ok := got["tags"]; !ok {
		t.Errorf("tags missing, want an empty array")
	} else if arr, isArr := v.([]any); !isArr || len(arr) != 0 {
		t.Errorf("tags = %v, want an empty array (the API types tags as a non-nullable array)", v)
	}
}

// TestStillSetAfterUpdate covers the guard that stops the provider sending a
// redundant clear PATCH.
//
// Some DefectDojo serializers declare a field default, so the full update
// already cleared the value; re-sending an explicit null is wasted at best, and
// defectdojo_metadata rejects a lone {"product": null} with "Metadata entries
// need either a product, endpoint, location or a finding" (verified on
// 3.1.101) even though the row is already correct.
func TestStillSetAfterUpdate(t *testing.T) {
	t.Parallel()

	branchTag := "main"
	tags := []string{"alpha"}
	ddResource := &ddTestDefectdojoResource{
		TestCreate: dd.TestCreate{
			BranchTag: &branchTag, // server kept it: still needs clearing
			BuildId:   nil,        // server already cleared it
			Tags:      &tags,      // still populated
		},
	}

	targets := []clearTarget{
		{jsonName: "branch_tag", ddFieldName: "BranchTag"},
		{jsonName: "build_id", ddFieldName: "BuildId"},
		{jsonName: "tags", ddFieldName: "Tags", isSlice: true},
		{jsonName: "made_up", ddFieldName: "NoSuchField"},
	}

	var got []string
	for _, target := range stillSetAfterUpdate(ddResource, targets) {
		got = append(got, target.jsonName)
	}
	sort.Strings(got)

	// build_id is dropped; made_up is kept because we cannot tell, so the PATCH
	// gets to decide rather than the provider silently skipping it.
	want := []string{"branch_tag", "made_up", "tags"}
	if len(got) != len(want) {
		t.Fatalf("stillSetAfterUpdate = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("stillSetAfterUpdate = %v, want %v", got, want)
		}
	}
}

// TestStillSetAfterUpdateDropsEmptyCollection guards the slice branch: an empty
// slice means the update already emptied it.
func TestStillSetAfterUpdateDropsEmptyCollection(t *testing.T) {
	t.Parallel()

	empty := []string{}
	ddResource := &ddTestDefectdojoResource{TestCreate: dd.TestCreate{Tags: &empty}}
	targets := []clearTarget{{jsonName: "tags", ddFieldName: "Tags", isSlice: true}}

	if got := stillSetAfterUpdate(ddResource, targets); len(got) != 0 {
		t.Errorf("stillSetAfterUpdate = %v, want empty: the update already emptied tags", got)
	}
}

// TestStillSetAfterUpdateKeepsZeroValuedScalars pins the difference between "the
// server dropped the value" and "the server stored a zero".
//
// DefectDojo keeps 0, false and "" verbatim; it does not coerce them to null.
// Verified on 3.1.101: PATCH {"percent_complete": 0, "branch_tag": ""} reads
// back as 0 and "". So a non-nil pointer to a zero value is a value the server
// chose to send, not evidence that the update already cleared the column. Only
// a nil pointer is that evidence.
//
// Judging emptiness by the pointee made this guard swallow the clear target,
// no explicit-null PATCH went out, the read path wrote the stored zero back
// over the planned null, and the apply died with the very inconsistency the
// clearing mechanism exists to prevent:
//
//	.percent_complete: was null, but now cty.NumberIntVal(0)
//	.branch_tag:       was null, but now cty.StringVal("")
//
// The Slice branch is a different convention and stays as it is: the API types
// collections as non-nullable arrays, so [] genuinely means "no elements".
func TestStillSetAfterUpdateKeepsZeroValuedScalars(t *testing.T) {
	t.Parallel()

	zeroInt := 0
	emptyString := ""
	falseBool := false

	for _, tc := range []struct {
		name       string
		ddResource defectdojoResource
		targets    []clearTarget
		want       []string
	}{
		{
			name: `*int holding 0 and *string holding ""`,
			ddResource: &ddTestDefectdojoResource{TestCreate: dd.TestCreate{
				PercentComplete: &zeroInt,
				BranchTag:       &emptyString,
				BuildId:         &emptyString,
			}},
			targets: []clearTarget{
				{jsonName: "percent_complete", ddFieldName: "PercentComplete"},
				{jsonName: "branch_tag", ddFieldName: "BranchTag"},
				{jsonName: "build_id", ddFieldName: "BuildId"},
			},
			want: []string{"branch_tag", "build_id", "percent_complete"},
		},
		{
			name: "*bool holding false",
			ddResource: &findingTemplateDefectdojoResource{FindingTemplate: dd.FindingTemplate{
				FixAvailable: &falseBool,
			}},
			targets: []clearTarget{{jsonName: "fix_available", ddFieldName: "FixAvailable"}},
			want:    []string{"fix_available"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var got []string
			for _, target := range stillSetAfterUpdate(tc.ddResource, tc.targets) {
				got = append(got, target.jsonName)
			}
			sort.Strings(got)

			if !slices.Equal(got, tc.want) {
				t.Errorf("stillSetAfterUpdate = %v, want %v: DefectDojo stores 0, false and \"\" "+
					"verbatim, so a non-nil pointer to one still needs an explicit-null PATCH",
					got, tc.want)
			}
		})
	}
}
