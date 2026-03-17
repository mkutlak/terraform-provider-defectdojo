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

func TestLanguageTypeResourcePopulate(t *testing.T) {
	expectedId := 99
	expectedLanguage := "Go"
	expectedColor := "#00ADD8"

	ddObj := languageTypeDefectdojoResource{
		LanguageType: dd.LanguageType{
			Id:       &expectedId,
			Language: expectedLanguage,
			Color:    &expectedColor,
		},
	}

	resourceData := languageTypeResourceData{}
	var tfResource terraformResourceData = &resourceData
	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddObj)

	assert.Equal(t, resourceData.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, resourceData.Language.ValueString(), expectedLanguage)
	assert.Equal(t, resourceData.Color.ValueString(), expectedColor)
}

func TestLanguageTypeResource_defectdojoResource(t *testing.T) {
	expectedLanguage := "Go"
	expectedColor := "#00ADD8"

	resourceData := languageTypeResourceData{
		Language: types.StringValue(expectedLanguage),
		Color:    types.StringValue(expectedColor),
	}

	ddResource := resourceData.defectdojoResource()
	var tfResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, tfResource, &ddResource)

	ddObj := ddResource.(*languageTypeDefectdojoResource)
	assert.Equal(t, ddObj.Language, expectedLanguage)
	assert.Equal(t, *ddObj.Color, expectedColor)
}
