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

func TestDojoGroupResourcePopulate(t *testing.T) {
	expectedId := 42
	expectedName := "Test Group"
	expectedDescription := "A test group"

	ddResource := dojoGroupDefectdojoResource{
		DojoGroup: dd.DojoGroup{
			Id:          &expectedId,
			Name:        expectedName,
			Description: &expectedDescription,
		},
	}

	resourceData := dojoGroupResourceData{}
	var terraformResource terraformResourceData = &resourceData

	populateResourceData(context.Background(), &diag.Diagnostics{}, &terraformResource, &ddResource)
	assert.Equal(t, resourceData.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, resourceData.Name.ValueString(), expectedName)
	assert.Equal(t, resourceData.Description.ValueString(), expectedDescription)
}

func TestDojoGroupResourcePopulateNils(t *testing.T) {
	ddResource := dojoGroupDefectdojoResource{
		DojoGroup: dd.DojoGroup{},
	}

	resourceData := dojoGroupResourceData{}
	var terraformResource terraformResourceData = &resourceData

	populateResourceData(context.Background(), &diag.Diagnostics{}, &terraformResource, &ddResource)
	assert.Equal(t, resourceData.Name.ValueString(), "")
	assert.Equal(t, resourceData.Description.IsNull(), true)
}

func TestDojoGroupResource__defectdojoResource(t *testing.T) {
	expectedName := "Test Group"
	expectedDescription := "A test group"

	resourceData := dojoGroupResourceData{
		Name:        types.StringValue(expectedName),
		Description: types.StringValue(expectedDescription),
	}

	ddRes := resourceData.defectdojoResource()
	ddGroup := ddRes.(*dojoGroupDefectdojoResource)
	var terraformResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, terraformResource, &ddRes)

	assert.Equal(t, ddGroup.Name, expectedName)
	assert.Equal(t, *ddGroup.Description, expectedDescription)
}
