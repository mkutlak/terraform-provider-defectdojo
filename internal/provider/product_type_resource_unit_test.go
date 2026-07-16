package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
	"gotest.tools/assert"
)

func TestProductTypeResourcePopulate(t *testing.T) {
	expectedId := 7
	expectedName := "Test Product Type"
	expectedDescription := "A test product type description"
	expectedCriticalProduct := true
	expectedKeyProduct := false
	expectedAuthorizedUsers := []int{7, 8}
	expectedAuthorizedUsersSet := types.SetValueMust(
		types.Int64Type,
		[]attr.Value{
			types.Int64Value(7),
			types.Int64Value(8),
		},
	)

	ddResource := productTypeDefectdojoResource{
		ProductType: dd.ProductType{
			Id:              &expectedId,
			Name:            expectedName,
			Description:     &expectedDescription,
			CriticalProduct: &expectedCriticalProduct,
			KeyProduct:      &expectedKeyProduct,
			AuthorizedUsers: &expectedAuthorizedUsers,
		},
	}

	resourceData := productTypeResourceData{}
	var tfResource terraformResourceData = &resourceData

	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddResource)

	assert.Equal(t, resourceData.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, resourceData.Name.ValueString(), expectedName)
	assert.Equal(t, resourceData.Description.ValueString(), expectedDescription)
	assert.Equal(t, resourceData.CriticalProduct.ValueBool(), expectedCriticalProduct)
	assert.Equal(t, resourceData.KeyProduct.ValueBool(), expectedKeyProduct)
	assert.DeepEqual(t, resourceData.AuthorizedUsers, expectedAuthorizedUsersSet)
}

func TestProductTypeResourcePopulateNils(t *testing.T) {
	ddResource := productTypeDefectdojoResource{
		ProductType: dd.ProductType{},
	}

	resourceData := productTypeResourceData{}
	var tfResource terraformResourceData = &resourceData

	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddResource)

	assert.Equal(t, resourceData.Id.IsNull(), true)
	assert.Equal(t, resourceData.Name.ValueString(), "")
	assert.Equal(t, resourceData.Description.IsNull(), true)
	assert.Equal(t, resourceData.CriticalProduct.IsNull(), true)
	assert.Equal(t, resourceData.KeyProduct.IsNull(), true)
	assert.DeepEqual(t, resourceData.AuthorizedUsers, types.SetNull(types.Int64Type))
}

func TestProductTypeResource__defectdojoResource(t *testing.T) {
	expectedName := "Test Product Type"
	expectedDescription := "A test product type description"
	expectedCriticalProduct := true
	expectedKeyProduct := false
	expectedAuthorizedUsers := []int{7, 8}
	expectedAuthorizedUsersSet := types.SetValueMust(
		types.Int64Type,
		[]attr.Value{
			types.Int64Value(7),
			types.Int64Value(8),
		},
	)

	resourceData := productTypeResourceData{
		Name:            types.StringValue(expectedName),
		Description:     types.StringValue(expectedDescription),
		CriticalProduct: types.BoolValue(expectedCriticalProduct),
		KeyProduct:      types.BoolValue(expectedKeyProduct),
		AuthorizedUsers: expectedAuthorizedUsersSet,
	}

	ddRes := resourceData.defectdojoResource()
	ddPT := ddRes.(*productTypeDefectdojoResource)
	var tfResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, tfResource, &ddRes)

	assert.Equal(t, ddPT.Name, expectedName)
	assert.Equal(t, *ddPT.Description, expectedDescription)
	assert.Equal(t, *ddPT.CriticalProduct, expectedCriticalProduct)
	assert.Equal(t, *ddPT.KeyProduct, expectedKeyProduct)
	assert.DeepEqual(t, *ddPT.AuthorizedUsers, expectedAuthorizedUsers)

	req := productTypeToRequest(ddPT.ProductType)
	assert.DeepEqual(t, *req.AuthorizedUsers, expectedAuthorizedUsers)
}
