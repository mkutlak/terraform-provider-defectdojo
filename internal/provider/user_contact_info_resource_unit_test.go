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

	ddResource := userContactInfoDefectdojoResource{
		UserContactInfo: dd.UserContactInfo{
			Id:                 &expectedId,
			User:               expectedUser,
			Title:              &expectedTitle,
			PhoneNumber:        &expectedPhoneNumber,
			CellNumber:         &expectedCellNumber,
			TwitterUsername:    &expectedTwitterUsername,
			GithubUsername:     &expectedGithubUsername,
			SlackUsername:      &expectedSlackUsername,
			SlackUserId:        &expectedSlackUserId,
			BlockExecution:     &expectedBlockExecution,
			ForcePasswordReset: &expectedForcePasswordReset,
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
}

func TestUserContactInfoResource__defectdojoResource(t *testing.T) {
	expectedUser := 99
	expectedTitle := "Dr."
	expectedPhoneNumber := "+1234567890"

	resourceData := userContactInfoResourceData{
		User:        types.Int64Value(int64(expectedUser)),
		Title:       types.StringValue(expectedTitle),
		PhoneNumber: types.StringValue(expectedPhoneNumber),
	}

	ddRes := resourceData.defectdojoResource()
	ddUserContactInfo := ddRes.(*userContactInfoDefectdojoResource)
	var terraformResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, terraformResource, &ddRes)

	assert.Equal(t, ddUserContactInfo.User, expectedUser)
	assert.Equal(t, *ddUserContactInfo.Title, expectedTitle)
	assert.Equal(t, *ddUserContactInfo.PhoneNumber, expectedPhoneNumber)
}
