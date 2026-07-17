package provider

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
	"gotest.tools/assert"
)

func TestLocationDataSourcePopulate(t *testing.T) {
	expectedId := 10
	expectedLocationType := "URL"
	expectedLocationValue := "https://example.com"
	expectedTags := []string{"foo", "bar"}
	expectedInheritedTags := []string{"baz"}
	expectedCreated := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	expectedUpdated := time.Date(2024, 6, 7, 8, 9, 10, 0, time.UTC)

	expectedTagsSet := types.SetValueMust(
		types.StringType,
		[]attr.Value{
			types.StringValue("foo"),
			types.StringValue("bar"),
		},
	)
	expectedInheritedTagsSet := types.SetValueMust(
		types.StringType,
		[]attr.Value{
			types.StringValue("baz"),
		},
	)

	ddLocation := locationDefectdojoResource{
		Location: dd.Location{
			Id:            &expectedId,
			LocationType:  &expectedLocationType,
			LocationValue: &expectedLocationValue,
			Tags:          &expectedTags,
			InheritedTags: &expectedInheritedTags,
			Created:       &expectedCreated,
			Updated:       &expectedUpdated,
		},
	}

	locationResource := locationResourceData{}
	var tfResource terraformResourceData = &locationResource
	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddLocation)

	assert.Equal(t, locationResource.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, locationResource.LocationType.ValueString(), expectedLocationType)
	assert.Equal(t, locationResource.LocationValue.ValueString(), expectedLocationValue)
	assert.DeepEqual(t, locationResource.Tags, expectedTagsSet)
	assert.DeepEqual(t, locationResource.InheritedTags, expectedInheritedTagsSet)
	assert.Equal(t, locationResource.Created.ValueString(), expectedCreated.Format(time.RFC3339))
	assert.Equal(t, locationResource.Updated.ValueString(), expectedUpdated.Format(time.RFC3339))
}

func TestLocationDataSourcePopulateNils(t *testing.T) {
	nilStringSet := types.SetNull(types.StringType)

	locationResource := locationResourceData{}
	var tfResource terraformResourceData = &locationResource

	ddLocation := locationDefectdojoResource{
		Location: dd.Location{},
	}
	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddLocation)

	assert.Equal(t, locationResource.Id.IsNull(), true)
	assert.Equal(t, locationResource.LocationType.IsNull(), true)
	assert.Equal(t, locationResource.LocationValue.IsNull(), true)
	assert.Equal(t, locationResource.Created.IsNull(), true)
	assert.Equal(t, locationResource.Updated.IsNull(), true)
	assert.DeepEqual(t, locationResource.Tags, nilStringSet)
	assert.DeepEqual(t, locationResource.InheritedTags, nilStringSet)
}

func TestLocationDataSource__defectdojoResource(t *testing.T) {
	expectedLocationType := "URL"
	expectedLocationValue := "https://example.com"
	expectedTags := []string{"foo", "bar"}
	expectedInheritedTags := []string{"baz"}
	expectedCreated := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	expectedUpdated := time.Date(2024, 6, 7, 8, 9, 10, 0, time.UTC)

	expectedTagsSet := types.SetValueMust(
		types.StringType,
		[]attr.Value{
			types.StringValue("foo"),
			types.StringValue("bar"),
		},
	)
	expectedInheritedTagsSet := types.SetValueMust(
		types.StringType,
		[]attr.Value{
			types.StringValue("baz"),
		},
	)

	locationResource := locationResourceData{
		LocationType:  types.StringValue(expectedLocationType),
		LocationValue: types.StringValue(expectedLocationValue),
		Tags:          expectedTagsSet,
		InheritedTags: expectedInheritedTagsSet,
		Created:       types.StringValue(expectedCreated.Format(time.RFC3339)),
		Updated:       types.StringValue(expectedUpdated.Format(time.RFC3339)),
	}

	ddResource := locationResource.defectdojoResource()
	ddLocation := ddResource.(*locationDefectdojoResource)
	var tfResource terraformResourceData = &locationResource
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, tfResource, &ddResource)

	assert.Equal(t, *ddLocation.LocationType, expectedLocationType)
	assert.Equal(t, *ddLocation.LocationValue, expectedLocationValue)
	assert.DeepEqual(t, *ddLocation.Tags, expectedTags)
	assert.DeepEqual(t, *ddLocation.InheritedTags, expectedInheritedTags)
	assert.Assert(t, ddLocation.Created.Equal(expectedCreated))
	assert.Assert(t, ddLocation.Updated.Equal(expectedUpdated))
}

func TestLocationDataSource__defectdojoResource_Nulls(t *testing.T) {
	var nilString *string
	var nilTime *time.Time
	var nilStringSlice *[]string

	locationResource := locationResourceData{
		Id:            types.StringNull(),
		LocationType:  types.StringNull(),
		LocationValue: types.StringNull(),
		Tags:          types.SetNull(types.StringType),
		InheritedTags: types.SetNull(types.StringType),
		Created:       types.StringNull(),
		Updated:       types.StringNull(),
	}

	ddResource := locationResource.defectdojoResource()
	ddLocation := ddResource.(*locationDefectdojoResource)
	var tfResource terraformResourceData = &locationResource
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, tfResource, &ddResource)

	assert.Equal(t, ddLocation.Id, (*int)(nil))
	assert.Equal(t, ddLocation.LocationType, nilString)
	assert.Equal(t, ddLocation.LocationValue, nilString)
	assert.Equal(t, ddLocation.Created, nilTime)
	assert.Equal(t, ddLocation.Updated, nilTime)
	assert.Equal(t, ddLocation.Tags, nilStringSlice)
	assert.Equal(t, ddLocation.InheritedTags, nilStringSlice)
}

func TestLocationDataSource_writeStubsAreReadOnly(t *testing.T) {
	ddObj := &locationDefectdojoResource{}

	status, body, err := ddObj.createApiCall(context.Background(), nil)
	assert.Equal(t, status, 0)
	assert.Assert(t, body == nil)
	assert.ErrorContains(t, err, "read-only in DefectDojo 3.x")

	status, body, err = ddObj.updateApiCall(context.Background(), nil, 1)
	assert.Equal(t, status, 0)
	assert.Assert(t, body == nil)
	assert.ErrorContains(t, err, "read-only in DefectDojo 3.x")

	status, body, err = ddObj.deleteApiCall(context.Background(), nil, 1)
	assert.Equal(t, status, 0)
	assert.Assert(t, body == nil)
	assert.ErrorContains(t, err, "read-only in DefectDojo 3.x")
}
