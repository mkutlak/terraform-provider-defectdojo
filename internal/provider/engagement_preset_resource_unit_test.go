package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
	"gotest.tools/assert"
)

func TestEngagementPresetResourcePopulate(t *testing.T) {
	expectedId := 11
	expectedTitle := "Security Preset"
	expectedProduct := 3
	expectedNotes := "Test these specific areas"
	expectedScope := "192.168.1.0/24, https://app.example.com"

	ddResource := engagementPresetDefectdojoResource{
		EngagementPresets: dd.EngagementPresets{
			Id:      &expectedId,
			Title:   &expectedTitle,
			Product: expectedProduct,
			Notes:   &expectedNotes,
			Scope:   &expectedScope,
		},
	}

	resourceData := engagementPresetResourceData{}
	var tfResource terraformResourceData = &resourceData

	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddResource)

	assert.Equal(t, resourceData.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, resourceData.Title.ValueString(), expectedTitle)
	assert.Equal(t, resourceData.Product.ValueInt64(), int64(expectedProduct))
	assert.Equal(t, resourceData.Notes.ValueString(), expectedNotes)
	assert.Equal(t, resourceData.Scope.ValueString(), expectedScope)
}

func TestEngagementPresetResourcePopulateNils(t *testing.T) {
	ddResource := engagementPresetDefectdojoResource{
		EngagementPresets: dd.EngagementPresets{
			// Product is a non-pointer required int
			Product: 0,
		},
	}

	resourceData := engagementPresetResourceData{}
	var tfResource terraformResourceData = &resourceData

	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddResource)

	assert.Equal(t, resourceData.Id.IsNull(), true)
	assert.Equal(t, resourceData.Title.IsNull(), true)
	assert.Equal(t, resourceData.Product.ValueInt64(), int64(0))
	assert.Equal(t, resourceData.Notes.IsNull(), true)
	assert.Equal(t, resourceData.Scope.IsNull(), true)
}

func TestEngagementPresetResource__defectdojoResource(t *testing.T) {
	expectedTitle := "Security Preset"
	expectedProduct := 3
	expectedNotes := "Test these specific areas"
	expectedScope := "192.168.1.0/24"

	resourceData := engagementPresetResourceData{
		Title:   types.StringValue(expectedTitle),
		Product: types.Int64Value(int64(expectedProduct)),
		Notes:   types.StringValue(expectedNotes),
		Scope:   types.StringValue(expectedScope),
	}

	ddRes := resourceData.defectdojoResource()
	ddPreset := ddRes.(*engagementPresetDefectdojoResource)
	var tfResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, tfResource, &ddRes)

	assert.Equal(t, *ddPreset.Title, expectedTitle)
	assert.Equal(t, ddPreset.Product, expectedProduct)
	assert.Equal(t, *ddPreset.Notes, expectedNotes)
	assert.Equal(t, *ddPreset.Scope, expectedScope)
}
