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
// attribute but no nullClearer implementation would fail at apply time with a
// provider bug, so the omission is a build failure instead.
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
				"defectdojoResource() does not implement nullClearer. Removing one of them from "+
				"configuration would fail the apply, because omitting a field from an update "+
				"request leaves it unchanged. Add partialUpdateWithBody - see clear.go for the "+
				"three-line shape.", tfTypeName, len(clearable), clearable)
		}
	}

	if needed == 0 {
		t.Fatal("no resource was found to have clearable attributes; the join is broken and this " +
			"test would pass vacuously")
	}
	t.Logf("%d resources have clearable attributes and implement nullClearer", needed)
}

// TestClearedDdFieldsDetectsRemovedAttributes drives the plan-vs-state diff over
// a real model. dd_test is the resource from issue #30.
//
// The zero-value rows are the state-side half of the distinction
// stillSetAfterUpdate has to make. This diff reads Terraform values, not Go
// ones, and there null is its own sentinel: types.Int64Value(0) and
// types.StringValue("") are non-null, so an attribute that held a zero and was
// then removed from configuration is a clear target like any other. It has to
// be, or removing branch_tag = "" would leave the empty string on the server.
func TestClearedDdFieldsDetectsRemovedAttributes(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name        string
		plan, state *ddTestResourceData
		wantNulls   []string
		wantEmpties []string
	}{
		{
			name: "attributes removed from configuration",
			plan: &ddTestResourceData{
				BranchTag:       types.StringNull(), // removed from config
				BuildId:         types.StringNull(), // removed from config
				CommitHash:      types.StringNull(), // removed from config
				PercentComplete: types.Int64Null(),  // removed from config
				Version:         types.StringValue("v2"),
				Title:           types.StringNull(), // never set: nothing to clear
				Tags:            types.SetNull(types.StringType),
			},
			state: &ddTestResourceData{
				BranchTag:       types.StringValue("main"),
				BuildId:         types.StringValue("b-1"),
				CommitHash:      types.StringValue("abc123"),
				PercentComplete: types.Int64Value(50),
				Version:         types.StringValue("v1"),
				Title:           types.StringNull(),
				Tags:            types.SetValueMust(types.StringType, []attr.Value{types.StringValue("a")}),
			},
			wantNulls: []string{"branch_tag", "build_id", "commit_hash", "percent_complete"},
			// A Set that held elements and is now null clears to [], not null:
			// the spec types tags as a non-nullable array.
			wantEmpties: []string{"tags"},
		},
		{
			name: "a removed zero is still a clear target",
			plan: &ddTestResourceData{
				BranchTag:       types.StringNull(),
				PercentComplete: types.Int64Null(),
				Tags:            types.SetNull(types.StringType),
			},
			state: &ddTestResourceData{
				BranchTag:       types.StringValue(""),
				PercentComplete: types.Int64Value(0),
				Tags:            types.SetValueMust(types.StringType, []attr.Value{}),
			},
			wantNulls:   []string{"branch_tag", "percent_complete"},
			wantEmpties: []string{"tags"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var nulls, empties []string
			for _, target := range clearedDdFields(tc.plan, tc.state, tc.plan.defectdojoResource()) {
				if target.isSlice {
					empties = append(empties, target.jsonName)
				} else {
					nulls = append(nulls, target.jsonName)
				}
			}
			sort.Strings(nulls)
			sort.Strings(empties)

			if !slices.Equal(nulls, tc.wantNulls) {
				t.Errorf("clearedDdFields nulls = %v, want %v", nulls, tc.wantNulls)
			}
			if !slices.Equal(empties, tc.wantEmpties) {
				t.Errorf("clearedDdFields empties = %v, want %v", empties, tc.wantEmpties)
			}
		})
	}
}

// TestClearedDdFieldsIgnoresUnchangedAndUnknown guards the two ways this diff
// could do damage: clearing something the practitioner never removed, and
// clearing something the server is about to fill in.
func TestClearedDdFieldsIgnoresUnchangedAndUnknown(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name        string
		plan, state types.String
	}{
		{"value unchanged", types.StringValue("main"), types.StringValue("main")},
		{"value changed", types.StringValue("release"), types.StringValue("main")},
		{"plan unknown (server will fill it)", types.StringUnknown(), types.StringValue("main")},
		{"state already null", types.StringNull(), types.StringNull()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			plan := &ddTestResourceData{BranchTag: tc.plan}
			state := &ddTestResourceData{BranchTag: tc.state}
			if got := clearedDdFields(plan, state, plan.defectdojoResource()); len(got) != 0 {
				t.Errorf("clearedDdFields = %v, want empty: nothing was removed from config", got)
			}
		})
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
			name:         "a set attribute",
			plan:         &productResourceData{RegulationIds: types.SetNull(types.Int64Type)},
			state:        &productResourceData{RegulationIds: types.SetValueMust(types.Int64Type, []attr.Value{types.Int64Value(1)})},
			wantAttrList: "Removing regulation_ids from the configuration",
			wantWireName: `"regulations"`,
		},
	} {
		// Not parallel: these subtests share the httptest server above, which
		// the deferred Close would shut down before a parallel subtest resumed.
		t.Run(tc.name, func(t *testing.T) {
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
// need either a product, endpoint, location or a finding" (verified on 3.1.101)
// even though the row is already correct.
//
// The zero-value rows pin the difference between "the server dropped the value"
// and "the server stored a zero". DefectDojo keeps 0, false and "" verbatim; it
// does not coerce them to null. Judging emptiness by the pointee instead of the
// pointer made this guard swallow the clear target, so no PATCH went out and the
// apply died with the very inconsistency the mechanism exists to prevent:
//
//	.percent_complete: was null, but now cty.NumberIntVal(0)
//	.branch_tag:       was null, but now cty.StringVal("")
//
// The Slice branch is a different convention and stays as it is: the API types
// collections as non-nullable arrays, so [] genuinely means "no elements".
func TestStillSetAfterUpdate(t *testing.T) {
	t.Parallel()

	branchTag := "main"
	tags := []string{"alpha"}
	empty := []string{}
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
			// build_id is dropped; made_up is kept because we cannot tell, so
			// the PATCH decides rather than the provider silently skipping it.
			name: "nil pointer means the update already cleared it",
			ddResource: &ddTestDefectdojoResource{TestCreate: dd.TestCreate{
				BranchTag: &branchTag, // server kept it: still needs clearing
				BuildId:   nil,        // server already cleared it
				Tags:      &tags,      // still populated
			}},
			targets: []clearTarget{
				{jsonName: "branch_tag", ddFieldName: "BranchTag"},
				{jsonName: "build_id", ddFieldName: "BuildId"},
				{jsonName: "tags", ddFieldName: "Tags", isSlice: true},
				{jsonName: "made_up", ddFieldName: "NoSuchField"},
			},
			want: []string{"branch_tag", "made_up", "tags"},
		},
		{
			name:       "an empty collection means the update already emptied it",
			ddResource: &ddTestDefectdojoResource{TestCreate: dd.TestCreate{Tags: &empty}},
			targets:    []clearTarget{{jsonName: "tags", ddFieldName: "Tags", isSlice: true}},
			want:       nil,
		},
		{
			name: `*int holding 0 and *string holding "" are still values`,
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
			name: "*bool holding false is still a value",
			ddResource: &findingTemplateDefectdojoResource{FindingTemplate: dd.FindingTemplate{
				FixAvailable: &falseBool,
			}},
			targets: []clearTarget{{jsonName: "fix_available", ddFieldName: "FixAvailable"}},
			want:    []string{"fix_available"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

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
