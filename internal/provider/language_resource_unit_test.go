package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gotest.tools/assert"
)

func TestLanguageResourcePopulate(t *testing.T) {
	expectedId := 60
	expectedLanguage := 2
	expectedProduct := 3
	expectedUser := 1
	expectedFiles := 42
	expectedCode := 1000
	expectedBlank := 100
	expectedComment := 50

	ddObj := languageDefectdojoResource{
		ddLanguageWrapper: ddLanguageWrapper{
			Id:             &expectedId,
			LanguageTypeId: expectedLanguage,
			Product:        expectedProduct,
			User:           &expectedUser,
			Files:          &expectedFiles,
			Code:           &expectedCode,
			Blank:          &expectedBlank,
			Comment:        &expectedComment,
		},
	}

	resourceData := languageResourceData{}
	var tfResource terraformResourceData = &resourceData
	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddObj)

	assert.Equal(t, resourceData.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, resourceData.Language.ValueInt64(), int64(expectedLanguage)) // Language field maps to language_type tfsdk attribute
	assert.Equal(t, resourceData.Product.ValueInt64(), int64(expectedProduct))
	assert.Equal(t, resourceData.User.ValueInt64(), int64(expectedUser))
	assert.Equal(t, resourceData.Files.ValueInt64(), int64(expectedFiles))
	assert.Equal(t, resourceData.Code.ValueInt64(), int64(expectedCode))
	assert.Equal(t, resourceData.Blank.ValueInt64(), int64(expectedBlank))
	assert.Equal(t, resourceData.Comment.ValueInt64(), int64(expectedComment))
}

func TestLanguageResource_defectdojoResource(t *testing.T) {
	expectedLanguage := int64(2)
	expectedProduct := int64(3)
	expectedCode := int64(500)

	resourceData := languageResourceData{
		Language: types.Int64Value(expectedLanguage),
		Product:  types.Int64Value(expectedProduct),
		Code:     types.Int64Value(expectedCode),
	}

	ddResource := resourceData.defectdojoResource()
	var tfResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, tfResource, &ddResource)

	ddObj := ddResource.(*languageDefectdojoResource)
	assert.Equal(t, int64(ddObj.LanguageTypeId), expectedLanguage)
	assert.Equal(t, int64(ddObj.Product), expectedProduct)
	assert.Equal(t, int64(*ddObj.Code), expectedCode)
}
