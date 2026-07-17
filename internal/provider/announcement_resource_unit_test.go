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

func TestAnnouncementResourcePopulate(t *testing.T) {
	expectedId := 1
	expectedMessage := "This is a test announcement"
	expectedStyle := "info"
	expectedDismissable := true

	ddObj := announcementDefectdojoResource{
		Announcement: dd.Announcement{
			Id:          &expectedId,
			Message:     &expectedMessage,
			Style:       (*dd.AnnouncementStyle)(&expectedStyle),
			Dismissable: &expectedDismissable,
		},
	}

	resourceData := announcementResourceData{}
	var tfResource terraformResourceData = &resourceData
	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddObj)

	assert.Equal(t, resourceData.Id.ValueString(), fmt.Sprint(expectedId))
	assert.Equal(t, resourceData.Message.ValueString(), expectedMessage)
	assert.Equal(t, resourceData.Style.ValueString(), expectedStyle)
	assert.Equal(t, resourceData.Dismissable.ValueBool(), expectedDismissable)
}

func TestAnnouncementResourcePopulate_NilStyle(t *testing.T) {
	expectedId := 1
	expectedMessage := "This is a test announcement"

	ddObj := announcementDefectdojoResource{
		Announcement: dd.Announcement{
			Id:      &expectedId,
			Message: &expectedMessage,
		},
	}

	resourceData := announcementResourceData{}
	var tfResource terraformResourceData = &resourceData
	populateResourceData(context.Background(), &diag.Diagnostics{}, &tfResource, &ddObj)

	assert.Equal(t, resourceData.Style.IsNull(), true)
	assert.Equal(t, resourceData.Dismissable.IsNull(), true)
}

func TestAnnouncementResource_defectdojoResource(t *testing.T) {
	expectedMessage := "This is a test announcement"
	expectedStyle := "warning"
	expectedDismissable := true

	resourceData := announcementResourceData{
		Message:     types.StringValue(expectedMessage),
		Style:       types.StringValue(expectedStyle),
		Dismissable: types.BoolValue(expectedDismissable),
	}

	ddResource := resourceData.defectdojoResource()
	var tfResource terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, tfResource, &ddResource)

	ddObj := ddResource.(*announcementDefectdojoResource)
	assert.Equal(t, *ddObj.Message, expectedMessage)
	assert.Equal(t, string(*ddObj.Style), expectedStyle)
	assert.Equal(t, *ddObj.Dismissable, expectedDismissable)
}

func TestAnnouncementToRequest(t *testing.T) {
	expectedMessage := "This is a test announcement"
	expectedStyle := "danger"
	expectedDismissable := false

	obj := dd.Announcement{
		Message:     &expectedMessage,
		Style:       (*dd.AnnouncementStyle)(&expectedStyle),
		Dismissable: &expectedDismissable,
	}

	req := announcementToRequest(obj)

	assert.Equal(t, *req.Message, expectedMessage)
	assert.Equal(t, *req.Dismissable, expectedDismissable)
	if req.Style == nil {
		t.Fatal("expected req.Style to be non-nil")
	}
	assert.Equal(t, string(*req.Style), expectedStyle)
}

func TestAnnouncementToRequest_NilStyle(t *testing.T) {
	expectedMessage := "This is a test announcement"

	obj := dd.Announcement{
		Message: &expectedMessage,
	}

	req := announcementToRequest(obj)

	assert.Equal(t, *req.Message, expectedMessage)
	if req.Style != nil {
		t.Fatal("expected req.Style to be nil")
	}
}
