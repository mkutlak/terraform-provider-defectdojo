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

func TestMetadataResourcePopulate(t *testing.T) {
	expectedId := 99
	expectedName := "A Name"
	expectedValue := "A Value"
	expectedProduct := 42
	expectedFinding := 45

	ddMetadata := metadataDefectdojoResource{
		Meta: dd.Meta{
			Id:      &expectedId,
			Name:    expectedName,
			Value:   expectedValue,
			Product: &expectedProduct,
			Finding: &expectedFinding,
		},
	}

	metadataResource := metadataResourceData{}
	var terraformResource terraformResourceData = &metadataResource

	populateResourceData(context.Background(), &diag.Diagnostics{}, &terraformResource, &ddMetadata)
	assert.Equal(t, metadataResource.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, metadataResource.Name.ValueString(), expectedName)
	assert.Equal(t, metadataResource.Value.ValueString(), expectedValue)
	assert.Equal(t, metadataResource.Product.ValueInt64(), (int64)(expectedProduct))
	assert.Equal(t, metadataResource.Finding.ValueInt64(), (int64)(expectedFinding))

	ddMetadata = metadataDefectdojoResource{
		Meta: dd.Meta{},
	}
	populateResourceData(context.Background(), &diag.Diagnostics{}, &terraformResource, &ddMetadata)

	// Name/Value are bare (non-pointer) strings on dd.Meta, so their zero
	// value is an empty string, not null.
	assert.Equal(t, metadataResource.Name.ValueString(), "")
	assert.Equal(t, metadataResource.Value.ValueString(), "")
	assert.Equal(t, metadataResource.Id.IsNull(), true)
	assert.Equal(t, metadataResource.Product.IsNull(), true)
	assert.Equal(t, metadataResource.Finding.IsNull(), true)
}

func TestMetadataResourcePopulateNils(t *testing.T) {
	metadataResource := metadataResourceData{}
	var terraformResource terraformResourceData = &metadataResource

	assert.Equal(t, metadataResource.Name.ValueString(), "")
	assert.Equal(t, metadataResource.Value.ValueString(), "")
	assert.Equal(t, metadataResource.Product.ValueInt64(), (int64)(0))
	assert.Equal(t, metadataResource.Finding.ValueInt64(), (int64)(0))

	ddMetadata := metadataDefectdojoResource{
		Meta: dd.Meta{},
	}
	populateResourceData(context.Background(), &diag.Diagnostics{}, &terraformResource, &ddMetadata)

	// still all empty/null values after running populate
	assert.Equal(t, metadataResource.Name.ValueString(), "")
	assert.Equal(t, metadataResource.Value.ValueString(), "")
	assert.Equal(t, metadataResource.Id.IsNull(), true)
	assert.Equal(t, metadataResource.Product.IsNull(), true)
	assert.Equal(t, metadataResource.Finding.IsNull(), true)
}

func TestMetadataResource__defectdojoResource(t *testing.T) {
	expectedName := "A Name"
	expectedValue := "A Value"
	expectedProduct := 42
	expectedFinding := 45

	metadataResource := metadataResourceData{
		Name:    types.StringValue(expectedName),
		Value:   types.StringValue(expectedValue),
		Product: types.Int64Value(int64(expectedProduct)),
		Finding: types.Int64Value(int64(expectedFinding)),
	}

	ddResource := metadataResource.defectdojoResource()
	ddMetadata := ddResource.(*metadataDefectdojoResource)
	var terraformResource terraformResourceData = &metadataResource
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, terraformResource, &ddResource)

	assert.Equal(t, ddMetadata.Name, expectedName)
	assert.Equal(t, ddMetadata.Value, expectedValue)
	assert.Equal(t, *ddMetadata.Product, expectedProduct)
	assert.Equal(t, *ddMetadata.Finding, expectedFinding)

	req := metadataToRequest(ddMetadata.Meta)
	assert.Equal(t, req.Name, expectedName)
	assert.Equal(t, req.Value, expectedValue)
	assert.Equal(t, *req.Product, expectedProduct)
	assert.Equal(t, *req.Finding, expectedFinding)
}

func TestMetadataResource__defectdojoResource_Nulls(t *testing.T) {
	var nilInt *int

	metadataResource := metadataResourceData{
		Id:      types.StringNull(),
		Name:    types.StringNull(),
		Value:   types.StringNull(),
		Product: types.Int64Null(),
		Finding: types.Int64Null(),
	}

	ddResource := metadataResource.defectdojoResource()
	ddMetadata := ddResource.(*metadataDefectdojoResource)
	var terraformResource terraformResourceData = &metadataResource
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, terraformResource, &ddResource)

	// Null TF values are skipped by populateDefectdojoResource, so the
	// underlying dd.Meta fields remain at their zero values.
	assert.Equal(t, ddMetadata.Id, nilInt)
	assert.Equal(t, ddMetadata.Name, "")
	assert.Equal(t, ddMetadata.Value, "")
	assert.Equal(t, ddMetadata.Product, nilInt)
	assert.Equal(t, ddMetadata.Finding, nilInt)
}
