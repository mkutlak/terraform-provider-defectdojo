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

func TestCicdInfrastructureResourcePopulate(t *testing.T) {
	expectedId := 42
	expectedName := "Test SCM Server"
	expectedDescription := "A test cicd infrastructure"
	expectedType := dd.CICDInfrastructureInfrastructureTypeScmServer
	expectedUrl := "https://scm.example.com"

	var ddUrl dd.CICDInfrastructure_Url
	assert.NilError(t, ddUrl.FromCICDInfrastructureUrl0(expectedUrl))

	ddObj := cicdInfrastructureDefectdojoResource{
		CICDInfrastructure: dd.CICDInfrastructure{
			Id:                 &expectedId,
			Name:               expectedName,
			Description:        &expectedDescription,
			InfrastructureType: expectedType,
			Url:                &ddUrl,
		},
	}

	resourceData := cicdInfrastructureResourceData{}
	var tfResource terraformResourceData = &resourceData
	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddObj)

	assert.Equal(t, resourceData.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, resourceData.Name.ValueString(), expectedName)
	assert.Equal(t, resourceData.Description.ValueString(), expectedDescription)
	assert.Equal(t, resourceData.InfrastructureType.ValueString(), string(expectedType))
	assert.Equal(t, resourceData.Url.ValueString(), expectedUrl)
}

// TestCicdInfrastructureResourcePopulate_UrlNull exercises the read path when
// the server has no url set at all: the wrapper pointer itself is nil.
func TestCicdInfrastructureResourcePopulate_UrlNull(t *testing.T) {
	expectedId := 43
	expectedName := "Test Build Server"
	expectedType := dd.CICDInfrastructureInfrastructureTypeBuildServer

	ddObj := cicdInfrastructureDefectdojoResource{
		CICDInfrastructure: dd.CICDInfrastructure{
			Id:                 &expectedId,
			Name:               expectedName,
			InfrastructureType: expectedType,
			Url:                nil,
		},
	}

	resourceData := cicdInfrastructureResourceData{}
	var tfResource terraformResourceData = &resourceData
	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddObj)

	assert.Assert(t, resourceData.Url.IsNull())
}

// TestCicdInfrastructureResourcePopulate_UrlEmpty exercises the read path for
// the oneOf's second arm (maxLength: 0), which DefectDojo sends when url is
// blank. This is the newest consumer of isOapiUnionStringType (resource.go),
// and empty-string/null coverage was flagged as the top gap in that engine.
func TestCicdInfrastructureResourcePopulate_UrlEmpty(t *testing.T) {
	expectedId := 44
	expectedName := "Test Orchestration"
	expectedType := dd.CICDInfrastructureInfrastructureTypeOrchestration

	var ddUrl dd.CICDInfrastructure_Url
	assert.NilError(t, ddUrl.FromCICDInfrastructureUrl1(""))

	ddObj := cicdInfrastructureDefectdojoResource{
		CICDInfrastructure: dd.CICDInfrastructure{
			Id:                 &expectedId,
			Name:               expectedName,
			InfrastructureType: expectedType,
			Url:                &ddUrl,
		},
	}

	resourceData := cicdInfrastructureResourceData{}
	var tfResource terraformResourceData = &resourceData
	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddObj)

	assert.Equal(t, resourceData.Url.ValueString(), "")
}

func TestCicdInfrastructureResource_defectdojoResource(t *testing.T) {
	expectedName := "Test SCM Server"
	expectedDescription := "A test cicd infrastructure"
	expectedType := "scm_server"
	expectedUrl := "https://scm.example.com"

	resourceData := cicdInfrastructureResourceData{
		Name:               types.StringValue(expectedName),
		Description:        types.StringValue(expectedDescription),
		InfrastructureType: types.StringValue(expectedType),
		Url:                types.StringValue(expectedUrl),
	}

	ddResource := resourceData.defectdojoResource()
	var tfResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, tfResource, &ddResource)

	ddObj := ddResource.(*cicdInfrastructureDefectdojoResource)
	assert.Equal(t, ddObj.Name, expectedName)
	assert.Equal(t, *ddObj.Description, expectedDescription)
	assert.Equal(t, string(ddObj.InfrastructureType), expectedType)

	gotUrl, err := ddObj.Url.AsCICDInfrastructureUrl0()
	assert.NilError(t, err)
	assert.Equal(t, gotUrl, expectedUrl)
}

// TestCicdInfrastructureResource_defectdojoResource_UrlEmpty exercises the
// write path when url is configured as an explicit empty string, the other
// arm of the union that a null round-trip would not cover.
func TestCicdInfrastructureResource_defectdojoResource_UrlEmpty(t *testing.T) {
	expectedName := "Test SCM Server"
	expectedType := "scm_server"

	resourceData := cicdInfrastructureResourceData{
		Name:               types.StringValue(expectedName),
		Description:        types.StringValue(""),
		InfrastructureType: types.StringValue(expectedType),
		Url:                types.StringValue(""),
	}

	ddResource := resourceData.defectdojoResource()
	var tfResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, tfResource, &ddResource)

	ddObj := ddResource.(*cicdInfrastructureDefectdojoResource)
	assert.Assert(t, ddObj.Url != nil)

	raw, err := ddObj.Url.MarshalJSON()
	assert.NilError(t, err)
	assert.Equal(t, string(raw), `""`)
}

// TestCicdInfrastructureResource_defectdojoResource_UrlNull exercises the write
// path when url is left null in configuration: the wrapper pointer itself must
// stay nil, or the request would send an explicit value the practitioner never
// configured.
func TestCicdInfrastructureResource_defectdojoResource_UrlNull(t *testing.T) {
	expectedName := "Test SCM Server"
	expectedType := "scm_server"

	resourceData := cicdInfrastructureResourceData{
		Name:               types.StringValue(expectedName),
		Description:        types.StringValue(""),
		InfrastructureType: types.StringValue(expectedType),
		Url:                types.StringNull(),
	}

	ddResource := resourceData.defectdojoResource()
	var tfResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, tfResource, &ddResource)

	ddObj := ddResource.(*cicdInfrastructureDefectdojoResource)
	assert.Assert(t, ddObj.Url == nil)
}

func TestCicdInfrastructureToRequest(t *testing.T) {
	description := "A test cicd infrastructure"
	var ddUrl dd.CICDInfrastructure_Url
	assert.NilError(t, ddUrl.FromCICDInfrastructureUrl0("https://scm.example.com"))

	obj := dd.CICDInfrastructure{
		Name:               "Test SCM Server",
		Description:        &description,
		InfrastructureType: dd.CICDInfrastructureInfrastructureTypeScmServer,
		Url:                &ddUrl,
	}

	req := cicdInfrastructureToRequest(obj)

	assert.Equal(t, req.Name, obj.Name)
	assert.Equal(t, *req.Description, description)
	assert.Equal(t, string(req.InfrastructureType), string(obj.InfrastructureType))
	gotUrl, err := req.Url.AsCICDInfrastructureRequestUrl0()
	assert.NilError(t, err)
	assert.Equal(t, gotUrl, "https://scm.example.com")
}
