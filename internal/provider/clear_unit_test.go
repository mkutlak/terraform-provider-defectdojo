package provider

import (
	"encoding/json"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

	nulls, empties := clearedDdFields(plan, state, plan.defectdojoResource())
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
			nulls, empties := clearedDdFields(plan, state, plan.defectdojoResource())
			if len(nulls) != 0 || len(empties) != 0 {
				t.Errorf("clearedDdFields = (%v, %v), want empty: nothing was removed from config", nulls, empties)
			}
		})
	}

	t.Run("state already null", func(t *testing.T) {
		plan := &ddTestResourceData{BranchTag: types.StringNull()}
		state := &ddTestResourceData{BranchTag: types.StringNull()}
		nulls, empties := clearedDdFields(plan, state, plan.defectdojoResource())
		if len(nulls) != 0 || len(empties) != 0 {
			t.Errorf("clearedDdFields = (%v, %v), want empty: there was nothing to clear", nulls, empties)
		}
	})
}

func TestClearPatchBody(t *testing.T) {
	t.Parallel()

	raw, err := clearPatchBody([]string{"branch_tag", "percent_complete"}, []string{"tags"})
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
