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

func TestUrlResourcePopulate(t *testing.T) {
	expectedId := 99
	expectedHost := "example.com"
	expectedProtocol := "https"
	expectedPort := 8443
	expectedPath := "/api/v1"
	expectedQuery := "foo=bar"
	expectedFragment := "section-1"
	expectedUserInfo := "user:pass"
	expectedHostValidationFailure := true
	expectedType := "domain"

	expectedTags := []string{"foo", "bar", "baz"}
	expectedTagsSet := types.SetValueMust(
		types.StringType,
		[]attr.Value{
			types.StringValue("foo"),
			types.StringValue("bar"),
			types.StringValue("baz"),
		},
	)

	ddUrl := urlDefectdojoResource{
		URL: dd.URL{
			Id:                    &expectedId,
			Host:                  expectedHost,
			Protocol:              &expectedProtocol,
			Port:                  &expectedPort,
			Path:                  &expectedPath,
			Query:                 &expectedQuery,
			Fragment:              &expectedFragment,
			UserInfo:              &expectedUserInfo,
			Tags:                  &expectedTags,
			HostValidationFailure: &expectedHostValidationFailure,
			Type:                  &expectedType,
		},
	}

	urlResource := urlResourceData{}
	var terraformResource terraformResourceData = &urlResource

	populateResourceData(context.Background(), &diag.Diagnostics{}, &terraformResource, &ddUrl)
	assert.Equal(t, urlResource.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, urlResource.Host.ValueString(), expectedHost)
	assert.Equal(t, urlResource.Protocol.ValueString(), expectedProtocol)
	assert.Equal(t, urlResource.Port.ValueInt64(), (int64)(expectedPort))
	assert.Equal(t, urlResource.Path.ValueString(), expectedPath)
	assert.Equal(t, urlResource.Query.ValueString(), expectedQuery)
	assert.Equal(t, urlResource.Fragment.ValueString(), expectedFragment)
	assert.Equal(t, urlResource.UserInfo.ValueString(), expectedUserInfo)
	assert.Equal(t, urlResource.HostValidationFailure.ValueBool(), expectedHostValidationFailure)
	assert.Equal(t, urlResource.Type.ValueString(), expectedType)
	assert.DeepEqual(t, urlResource.Tags, expectedTagsSet)

	ddUrl = urlDefectdojoResource{
		URL: dd.URL{},
	}
	populateResourceData(context.Background(), &diag.Diagnostics{}, &terraformResource, &ddUrl)

	nilStringSet := types.SetNull(types.StringType)

	assert.Equal(t, urlResource.Host.ValueString(), "")
	assert.Equal(t, urlResource.Protocol.IsNull(), true)
	assert.Equal(t, urlResource.Port.IsNull(), true)
	assert.Equal(t, urlResource.Path.IsNull(), true)
	assert.Equal(t, urlResource.Query.IsNull(), true)
	assert.Equal(t, urlResource.Fragment.IsNull(), true)
	assert.Equal(t, urlResource.UserInfo.IsNull(), true)
	assert.Equal(t, urlResource.HostValidationFailure.IsNull(), true)
	assert.Equal(t, urlResource.Type.IsNull(), true)
	assert.DeepEqual(t, urlResource.Tags, nilStringSet)
}

func TestUrlResourcePopulateNils(t *testing.T) {
	nilStringSet := types.SetNull(types.StringType)

	urlResource := urlResourceData{}
	var terraformResource terraformResourceData = &urlResource

	assert.Equal(t, urlResource.Host.ValueString(), "")
	assert.Equal(t, urlResource.Protocol.ValueString(), "")
	assert.Equal(t, urlResource.Port.ValueInt64(), (int64)(0))
	assert.Equal(t, urlResource.Path.ValueString(), "")
	assert.Equal(t, urlResource.Query.ValueString(), "")
	assert.Equal(t, urlResource.Fragment.ValueString(), "")
	assert.Equal(t, urlResource.UserInfo.ValueString(), "")
	assert.Equal(t, urlResource.HostValidationFailure.ValueBool(), false)
	assert.Equal(t, urlResource.Type.ValueString(), "")

	assert.DeepEqual(t, urlResource.Tags.Elements(), []attr.Value{})

	ddUrl := urlDefectdojoResource{
		URL: dd.URL{},
	}
	populateResourceData(context.Background(), &diag.Diagnostics{}, &terraformResource, &ddUrl)

	// still all empty/null values after running populate
	assert.Equal(t, urlResource.Host.ValueString(), "")
	assert.Equal(t, urlResource.Protocol.ValueString(), "")
	assert.Equal(t, urlResource.Port.ValueInt64(), (int64)(0))
	assert.Equal(t, urlResource.Path.ValueString(), "")
	assert.Equal(t, urlResource.Query.ValueString(), "")
	assert.Equal(t, urlResource.Fragment.ValueString(), "")
	assert.Equal(t, urlResource.UserInfo.ValueString(), "")
	assert.Equal(t, urlResource.HostValidationFailure.ValueBool(), false)
	assert.Equal(t, urlResource.Type.ValueString(), "")
	assert.DeepEqual(t, urlResource.Tags, nilStringSet)
}

// TestUrlResourcePopulateHost drives the ddFormat:"host" mechanism through the
// reflection engine itself, not just the helper underneath it.
//
// DefectDojo lower-cases every host it stores, so a configured "API.Example.COM"
// is echoed back as "api.example.com". Without the tag the read path wrote the
// server's spelling over the practitioner's, state disagreed with config, and
// Terraform failed the apply with "Provider produced inconsistent result after
// apply". The drift row is the other half of that contract: preservation must
// not hide an out-of-band change.
//
// clean_host() also punycodes an IDNA name and compresses an IPv6 address.
// Those are not case-only differences, so they reach state as the server
// spelled them.
func TestUrlResourcePopulateHost(t *testing.T) {
	for _, tc := range []struct {
		name       string
		configured string
		server     string
		want       string
	}{
		{"case-only difference keeps the configured spelling", "API.Example.COM", "api.example.com", "API.Example.COM"},
		{"a different host is reported as drift", "API.Example.COM", "other.example.com", "other.example.com"},
		{"IDNA punycoding reaches state verbatim", "Bücher.example", "xn--bcher-kva.example", "xn--bcher-kva.example"},
		{"IPv6 compression reaches state verbatim", "2001:db8:0:0:0:0:0:1", "2001:db8::1", "2001:db8::1"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ddUrl := urlDefectdojoResource{URL: dd.URL{Host: tc.server}}
			urlResource := urlResourceData{Host: types.StringValue(tc.configured)}
			var terraformResource terraformResourceData = &urlResource

			diags := diag.Diagnostics{}
			populateResourceData(context.Background(), &diags, &terraformResource, &ddUrl)

			assert.Equal(t, diags.HasError(), false)
			assert.Equal(t, urlResource.Host.ValueString(), tc.want)
		})
	}
}

func TestUrlResource__defectdojoResource(t *testing.T) {
	expectedHost := "example.com"
	expectedProtocol := "https"
	expectedPort := 8443
	expectedPath := "/api/v1"
	expectedQuery := "foo=bar"
	expectedFragment := "section-1"
	expectedUserInfo := "user:pass"

	expectedTags := []string{"foo", "bar", "baz"}
	expectedTagsSet := types.SetValueMust(
		types.StringType,
		[]attr.Value{
			types.StringValue("foo"),
			types.StringValue("bar"),
			types.StringValue("baz"),
		},
	)

	urlResource := urlResourceData{
		Host:     types.StringValue(expectedHost),
		Protocol: types.StringValue(expectedProtocol),
		Port:     types.Int64Value(int64(expectedPort)),
		Path:     types.StringValue(expectedPath),
		Query:    types.StringValue(expectedQuery),
		Fragment: types.StringValue(expectedFragment),
		UserInfo: types.StringValue(expectedUserInfo),

		Tags: expectedTagsSet,
	}

	ddResource := urlResource.defectdojoResource()
	ddUrl := ddResource.(*urlDefectdojoResource)
	var terraformResource terraformResourceData = &urlResource
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, terraformResource, &ddResource)

	assert.Equal(t, ddUrl.Host, expectedHost)
	assert.Equal(t, *ddUrl.Protocol, expectedProtocol)
	assert.Equal(t, *ddUrl.Port, expectedPort)
	assert.Equal(t, *ddUrl.Path, expectedPath)
	assert.Equal(t, *ddUrl.Query, expectedQuery)
	assert.Equal(t, *ddUrl.Fragment, expectedFragment)
	assert.Equal(t, *ddUrl.UserInfo, expectedUserInfo)

	assert.DeepEqual(t, *ddUrl.Tags, expectedTags)

	req := urlToRequest(ddUrl.URL)
	assert.Equal(t, req.Host, expectedHost)
	assert.Equal(t, *req.Protocol, expectedProtocol)
	assert.Equal(t, *req.Port, expectedPort)
	assert.Equal(t, *req.Path, expectedPath)
	assert.Equal(t, *req.Query, expectedQuery)
	assert.Equal(t, *req.Fragment, expectedFragment)
	assert.Equal(t, *req.UserInfo, expectedUserInfo)
	assert.DeepEqual(t, *req.Tags, expectedTags)
}

func TestUrlResource__defectdojoResource_Nulls(t *testing.T) {
	var nilString *string
	var nilInt *int
	var nilBool *bool
	var nilStringSlice *[]string

	urlResource := urlResourceData{
		Id:                    types.StringNull(),
		Host:                  types.StringNull(),
		Protocol:              types.StringNull(),
		Port:                  types.Int64Null(),
		Path:                  types.StringNull(),
		Query:                 types.StringNull(),
		Fragment:              types.StringNull(),
		UserInfo:              types.StringNull(),
		HostValidationFailure: types.BoolNull(),
		Type:                  types.StringNull(),

		Tags: types.SetNull(types.StringType),
	}

	ddResource := urlResource.defectdojoResource()
	ddUrl := ddResource.(*urlDefectdojoResource)
	var terraformResource terraformResourceData = &urlResource
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, terraformResource, &ddResource)

	assert.Equal(t, ddUrl.Id, nilInt)
	assert.Equal(t, ddUrl.Host, "")
	assert.Equal(t, ddUrl.Protocol, nilString)
	assert.Equal(t, ddUrl.Port, nilInt)
	assert.Equal(t, ddUrl.Path, nilString)
	assert.Equal(t, ddUrl.Query, nilString)
	assert.Equal(t, ddUrl.Fragment, nilString)
	assert.Equal(t, ddUrl.UserInfo, nilString)
	assert.Equal(t, ddUrl.HostValidationFailure, nilBool)
	assert.Equal(t, ddUrl.Type, nilString)

	// Null TF values are skipped, so pointer fields remain nil
	assert.Equal(t, ddUrl.Tags, nilStringSlice)
}
