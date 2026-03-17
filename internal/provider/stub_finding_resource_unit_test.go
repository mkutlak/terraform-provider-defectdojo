package provider

import (
	"context"
	"fmt"
	"testing"

	dd "github.com/doximity/terraform-provider-defectdojo/internal/ddclient"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gotest.tools/assert"
)

func TestStubFindingResourcePopulate(t *testing.T) {
	expectedId := 30
	expectedTitle := "Test Stub Finding"
	expectedTest := 7
	expectedSeverity := "Medium"
	expectedDescription := "A stub finding description"
	expectedReporter := 1

	ddObj := stubFindingDefectdojoResource{
		StubFindingCreate: dd.StubFindingCreate{
			Id:          &expectedId,
			Title:       expectedTitle,
			Test:        expectedTest,
			Severity:    &expectedSeverity,
			Description: &expectedDescription,
			Reporter:    &expectedReporter,
		},
	}

	resourceData := stubFindingResourceData{}
	var tfResource terraformResourceData = &resourceData
	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddObj)

	assert.Equal(t, resourceData.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, resourceData.Title.ValueString(), expectedTitle)
	assert.Equal(t, resourceData.Test.ValueInt64(), int64(expectedTest))
	assert.Equal(t, resourceData.Severity.ValueString(), expectedSeverity)
	assert.Equal(t, resourceData.Description.ValueString(), expectedDescription)
	assert.Equal(t, resourceData.Reporter.ValueInt64(), int64(expectedReporter))
}

func TestStubFindingResource_defectdojoResource(t *testing.T) {
	expectedTitle := "Test Stub Finding"
	expectedTest := int64(7)
	expectedSeverity := "Low"

	resourceData := stubFindingResourceData{
		Title:    types.StringValue(expectedTitle),
		Test:     types.Int64Value(expectedTest),
		Severity: types.StringValue(expectedSeverity),
	}

	ddResource := resourceData.defectdojoResource()
	var tfResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, tfResource, &ddResource)

	ddObj := ddResource.(*stubFindingDefectdojoResource)
	assert.Equal(t, ddObj.Title, expectedTitle)
	assert.Equal(t, int64(ddObj.Test), expectedTest)
	assert.Equal(t, *ddObj.Severity, expectedSeverity)
}
