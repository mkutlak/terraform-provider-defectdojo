package provider

import (
	"context"
	"testing"

	dd "github.com/doximity/terraform-provider-defectdojo/internal/ddclient"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gotest.tools/assert"
)

func TestAnnouncementResourcePopulate(t *testing.T) {
	expectedId := 1
	expectedMessage := "System maintenance at midnight"
	expectedStyle := dd.AnnouncementStyle("warning")
	expectedDismissable := true

	ddResource := announcementDefectdojoResource{
		Announcement: dd.Announcement{
			Id:          &expectedId,
			Message:     &expectedMessage,
			Style:       &expectedStyle,
			Dismissable: &expectedDismissable,
		},
	}

	resourceData := announcementResourceData{}
	var trd terraformResourceData = &resourceData

	populateResourceData(context.Background(), &diag.Diagnostics{}, &trd, &ddResource)

	assert.Equal(t, resourceData.Id.ValueString(), "1")
	assert.Equal(t, resourceData.Message.ValueString(), expectedMessage)
	assert.Equal(t, resourceData.Style.ValueString(), "warning")
	assert.Equal(t, resourceData.Dismissable.ValueBool(), expectedDismissable)
}

func TestAnnouncementResourcePopulateDefectdojo(t *testing.T) {
	resourceData := announcementResourceData{
		Message:     types.StringValue("System maintenance at midnight"),
		Style:       types.StringValue("warning"),
		Dismissable: types.BoolValue(true),
	}

	ddRes := resourceData.defectdojoResource()
	ddAnnouncement := ddRes.(*announcementDefectdojoResource)
	var trd terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, trd, &ddRes)

	assert.Equal(t, *ddAnnouncement.Message, "System maintenance at midnight")
	assert.Equal(t, *ddAnnouncement.Dismissable, true)
}
