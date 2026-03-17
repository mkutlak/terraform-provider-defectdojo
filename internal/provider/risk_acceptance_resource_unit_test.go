package provider

import (
	"context"
	"testing"

	dd "github.com/doximity/terraform-provider-defectdojo/internal/ddclient"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gotest.tools/assert"
)

func TestRiskAcceptanceResourcePopulate(t *testing.T) {
	expectedId := 20
	expectedName := "Accept CVE-2023-1234"
	expectedOwner := 5
	expectedFindings := []int{100, 101}
	expectedDecision := dd.RiskAcceptanceDecision("A")

	ddResource := riskAcceptanceDefectdojoResource{
		RiskAcceptance: dd.RiskAcceptance{
			Id:               &expectedId,
			Name:             expectedName,
			Owner:            expectedOwner,
			AcceptedFindings: expectedFindings,
			Decision:         &expectedDecision,
		},
	}

	resourceData := riskAcceptanceResourceData{}
	var trd terraformResourceData = &resourceData

	populateResourceData(context.Background(), &diag.Diagnostics{}, &trd, &ddResource)

	assert.Equal(t, resourceData.Id.ValueString(), "20")
	assert.Equal(t, resourceData.Name.ValueString(), expectedName)
	assert.Equal(t, resourceData.Owner.ValueInt64(), int64(expectedOwner))
	assert.Equal(t, resourceData.Decision.ValueString(), "A")

	expectedFindingsSet := types.SetValueMust(
		types.Int64Type,
		[]attr.Value{
			types.Int64Value(100),
			types.Int64Value(101),
		},
	)
	assert.DeepEqual(t, resourceData.AcceptedFindings, expectedFindingsSet)
}

func TestRiskAcceptanceResourcePopulateDefectdojo(t *testing.T) {
	resourceData := riskAcceptanceResourceData{
		Name:  types.StringValue("Accept CVE-2023-1234"),
		Owner: types.Int64Value(5),
		AcceptedFindings: types.SetValueMust(
			types.Int64Type,
			[]attr.Value{
				types.Int64Value(100),
				types.Int64Value(101),
			},
		),
		Decision: types.StringValue("A"),
	}

	ddRes := resourceData.defectdojoResource()
	ddRisk := ddRes.(*riskAcceptanceDefectdojoResource)
	var trd terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, trd, &ddRes)

	assert.Equal(t, ddRisk.Name, "Accept CVE-2023-1234")
	assert.Equal(t, ddRisk.Owner, 5)
	assert.DeepEqual(t, ddRisk.AcceptedFindings, []int{100, 101})
}
