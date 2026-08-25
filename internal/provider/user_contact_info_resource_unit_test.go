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

func TestUserContactInfoResourcePopulate(t *testing.T) {
	expectedId := 10
	expectedUser := 99
	expectedTitle := "Dr."
	expectedPhoneNumber := "+1234567890"
	expectedCellNumber := "+0987654321"
	expectedTwitterUsername := "twitteruser"
	expectedGithubUsername := "githubuser"
	expectedSlackUsername := "slackuser@example.com"
	expectedSlackUserId := "U12345"
	expectedBlockExecution := true
	expectedForcePasswordReset := false
	expectedDeduplicationExecutionMode := dd.UserContactInfoDeduplicationExecutionMode("async_wait")

	// PhoneNumber and CellNumber are oapi-codegen oneOf-string wrappers as of
	// DefectDojo 3.2 (see isOapiUnionStringType in resource.go); they cannot
	// be built with a plain &string literal.
	phoneNumber := &dd.UserContactInfo_PhoneNumber{}
	assert.NilError(t, phoneNumber.FromUserContactInfoPhoneNumber0(expectedPhoneNumber))
	cellNumber := &dd.UserContactInfo_CellNumber{}
	assert.NilError(t, cellNumber.FromUserContactInfoCellNumber0(expectedCellNumber))

	ddResource := userContactInfoDefectdojoResource{
		UserContactInfo: dd.UserContactInfo{
			Id:                         &expectedId,
			User:                       expectedUser,
			Title:                      &expectedTitle,
			PhoneNumber:                phoneNumber,
			CellNumber:                 cellNumber,
			TwitterUsername:            &expectedTwitterUsername,
			GithubUsername:             &expectedGithubUsername,
			SlackUsername:              &expectedSlackUsername,
			SlackUserId:                &expectedSlackUserId,
			BlockExecution:             &expectedBlockExecution,
			ForcePasswordReset:         &expectedForcePasswordReset,
			DeduplicationExecutionMode: &expectedDeduplicationExecutionMode,
		},
	}

	resourceData := userContactInfoResourceData{}
	var terraformResource terraformResourceData = &resourceData

	populateResourceData(context.Background(), &diag.Diagnostics{}, &terraformResource, &ddResource)
	assert.Equal(t, resourceData.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, resourceData.User.ValueInt64(), int64(expectedUser))
	assert.Equal(t, resourceData.Title.ValueString(), expectedTitle)
	assert.Equal(t, resourceData.PhoneNumber.ValueString(), expectedPhoneNumber)
	assert.Equal(t, resourceData.CellNumber.ValueString(), expectedCellNumber)
	assert.Equal(t, resourceData.TwitterUsername.ValueString(), expectedTwitterUsername)
	assert.Equal(t, resourceData.GithubUsername.ValueString(), expectedGithubUsername)
	assert.Equal(t, resourceData.SlackUsername.ValueString(), expectedSlackUsername)
	assert.Equal(t, resourceData.SlackUserId.ValueString(), expectedSlackUserId)
	assert.Equal(t, resourceData.BlockExecution.ValueBool(), expectedBlockExecution)
	assert.Equal(t, resourceData.ForcePasswordReset.ValueBool(), expectedForcePasswordReset)
	assert.Equal(t, resourceData.DeduplicationExecutionMode.ValueString(), string(expectedDeduplicationExecutionMode))
}

// TestUserContactInfoResourcePopulateNils covers the read path when
// PhoneNumber and CellNumber are absent server-side, so their wrapper
// pointers are nil. Both are oapi-codegen oneOf-string wrappers (see
// isOapiUnionStringType in resource.go).
func TestUserContactInfoResourcePopulateNils(t *testing.T) {
	ddResource := userContactInfoDefectdojoResource{
		UserContactInfo: dd.UserContactInfo{User: 99},
	}

	resourceData := userContactInfoResourceData{}
	var terraformResource terraformResourceData = &resourceData
	populateResourceData(context.Background(), &diag.Diagnostics{}, &terraformResource, &ddResource)

	assert.Equal(t, resourceData.PhoneNumber.IsNull(), true)
	assert.Equal(t, resourceData.CellNumber.IsNull(), true)
}

// TestUserContactInfoResourcePopulate_NumbersEmpty exercises the read path
// for the oneOf's second arm (maxLength: 0), which DefectDojo sends when
// PhoneNumber or CellNumber is blank. Empty-string coverage was flagged as
// the top gap in the isOapiUnionStringType engine (resource.go).
func TestUserContactInfoResourcePopulate_NumbersEmpty(t *testing.T) {
	phoneNumber := &dd.UserContactInfo_PhoneNumber{}
	assert.NilError(t, phoneNumber.FromUserContactInfoPhoneNumber1(""))
	cellNumber := &dd.UserContactInfo_CellNumber{}
	assert.NilError(t, cellNumber.FromUserContactInfoCellNumber1(""))

	ddResource := userContactInfoDefectdojoResource{
		UserContactInfo: dd.UserContactInfo{
			User:        99,
			PhoneNumber: phoneNumber,
			CellNumber:  cellNumber,
		},
	}

	resourceData := userContactInfoResourceData{}
	var terraformResource terraformResourceData = &resourceData
	populateResourceData(context.Background(), &diag.Diagnostics{}, &terraformResource, &ddResource)

	assert.Equal(t, resourceData.PhoneNumber.IsNull(), false)
	assert.Equal(t, resourceData.PhoneNumber.ValueString(), "")
	assert.Equal(t, resourceData.CellNumber.IsNull(), false)
	assert.Equal(t, resourceData.CellNumber.ValueString(), "")
}

// TestUserContactInfoResource_defectdojoResource_NumbersEmpty exercises the
// write path when PhoneNumber and CellNumber are configured as an explicit
// empty string, the other arm of the union a null round-trip does not cover.
func TestUserContactInfoResource_defectdojoResource_NumbersEmpty(t *testing.T) {
	resourceData := userContactInfoResourceData{
		User:        types.Int64Value(99),
		PhoneNumber: types.StringValue(""),
		CellNumber:  types.StringValue(""),
	}

	ddRes := resourceData.defectdojoResource()
	var terraformResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, terraformResource, &ddRes)

	ddUserContactInfo := ddRes.(*userContactInfoDefectdojoResource)
	assert.Assert(t, ddUserContactInfo.PhoneNumber != nil)
	assert.Assert(t, ddUserContactInfo.CellNumber != nil)

	phoneRaw, err := ddUserContactInfo.PhoneNumber.MarshalJSON()
	assert.NilError(t, err)
	assert.Equal(t, string(phoneRaw), `""`)

	cellRaw, err := ddUserContactInfo.CellNumber.MarshalJSON()
	assert.NilError(t, err)
	assert.Equal(t, string(cellRaw), `""`)
}

// TestUserContactInfoResource_defectdojoResource_NumbersNull exercises the
// write path when PhoneNumber and CellNumber are left null in configuration:
// the wrapper pointers must stay nil, or the request would send a value the
// practitioner never configured.
func TestUserContactInfoResource_defectdojoResource_NumbersNull(t *testing.T) {
	resourceData := userContactInfoResourceData{
		User:        types.Int64Value(99),
		PhoneNumber: types.StringNull(),
		CellNumber:  types.StringNull(),
	}

	ddRes := resourceData.defectdojoResource()
	var terraformResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, terraformResource, &ddRes)

	ddUserContactInfo := ddRes.(*userContactInfoDefectdojoResource)
	assert.Assert(t, ddUserContactInfo.PhoneNumber == nil)
	assert.Assert(t, ddUserContactInfo.CellNumber == nil)
}

func TestUserContactInfoResource__defectdojoResource(t *testing.T) {
	expectedUser := 99
	expectedTitle := "Dr."
	expectedPhoneNumber := "+1234567890"
	expectedDeduplicationExecutionMode := "async_wait"

	resourceData := userContactInfoResourceData{
		User:                       types.Int64Value(int64(expectedUser)),
		Title:                      types.StringValue(expectedTitle),
		PhoneNumber:                types.StringValue(expectedPhoneNumber),
		DeduplicationExecutionMode: types.StringValue(expectedDeduplicationExecutionMode),
	}

	ddRes := resourceData.defectdojoResource()
	ddUserContactInfo := ddRes.(*userContactInfoDefectdojoResource)
	var terraformResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, terraformResource, &ddRes)

	assert.Equal(t, ddUserContactInfo.User, expectedUser)
	assert.Equal(t, *ddUserContactInfo.Title, expectedTitle)
	assert.Equal(t, string(*ddUserContactInfo.DeduplicationExecutionMode), expectedDeduplicationExecutionMode)

	// PhoneNumber is an oapi-codegen oneOf-string wrapper as of DefectDojo 3.2
	// (see isOapiUnionStringType in resource.go); confirm the write path sets
	// it correctly, then confirm userContactInfoToRequest can carry it across
	// to the distinct UserContactInfoRequest_PhoneNumber wrapper type.
	gotPhoneNumber, err := ddUserContactInfo.PhoneNumber.AsUserContactInfoPhoneNumber0()
	assert.NilError(t, err)
	assert.Equal(t, gotPhoneNumber, expectedPhoneNumber)

	req := userContactInfoToRequest(ddUserContactInfo.UserContactInfo)
	gotReqPhoneNumber, err := req.PhoneNumber.AsUserContactInfoRequestPhoneNumber0()
	assert.NilError(t, err)
	assert.Equal(t, gotReqPhoneNumber, expectedPhoneNumber)
}
