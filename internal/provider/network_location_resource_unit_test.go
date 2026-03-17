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

func TestNetworkLocationResourcePopulate(t *testing.T) {
	expectedId := 99
	expectedLocation := "VPN"

	ddObj := networkLocationDefectdojoResource{
		NetworkLocations: dd.NetworkLocations{
			Id:       &expectedId,
			Location: expectedLocation,
		},
	}

	resourceData := networkLocationResourceData{}
	var tfResource terraformResourceData = &resourceData
	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddObj)

	assert.Equal(t, resourceData.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, resourceData.Location.ValueString(), expectedLocation)
}

func TestNetworkLocationResource_defectdojoResource(t *testing.T) {
	expectedLocation := "VPN"

	resourceData := networkLocationResourceData{
		Location: types.StringValue(expectedLocation),
	}

	ddResource := resourceData.defectdojoResource()
	var tfResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, tfResource, &ddResource)

	ddObj := ddResource.(*networkLocationDefectdojoResource)
	assert.Equal(t, ddObj.Location, expectedLocation)
}
