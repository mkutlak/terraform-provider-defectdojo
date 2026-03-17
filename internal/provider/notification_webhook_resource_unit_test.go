package provider

import (
	"context"
	"testing"

	dd "github.com/doximity/terraform-provider-defectdojo/internal/ddclient"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gotest.tools/assert"
)

func TestNotificationWebhookResourcePopulate(t *testing.T) {
	expectedId := 15
	expectedName := "My Webhook"
	expectedUrl := "https://hooks.example.com/notify"
	expectedHeaderName := "X-Webhook-Secret"
	expectedStatus := "Active"

	ddResource := notificationWebhookDefectdojoResource{
		NotificationWebhooks: dd.NotificationWebhooks{
			Id:         &expectedId,
			Name:       &expectedName,
			Url:        &expectedUrl,
			HeaderName: &expectedHeaderName,
			Status:     &expectedStatus,
		},
	}

	resourceData := notificationWebhookResourceData{}
	var trd terraformResourceData = &resourceData

	populateResourceData(context.Background(), &diag.Diagnostics{}, &trd, &ddResource)

	assert.Equal(t, resourceData.Id.ValueString(), "15")
	assert.Equal(t, resourceData.Name.ValueString(), expectedName)
	assert.Equal(t, resourceData.Url.ValueString(), expectedUrl)
	assert.Equal(t, resourceData.HeaderName.ValueString(), expectedHeaderName)
	assert.Equal(t, resourceData.Status.ValueString(), expectedStatus)
}

func TestNotificationWebhookResourcePopulateDefectdojo(t *testing.T) {
	resourceData := notificationWebhookResourceData{
		Name:       types.StringValue("My Webhook"),
		Url:        types.StringValue("https://hooks.example.com/notify"),
		HeaderName: types.StringValue("X-Webhook-Secret"),
	}

	ddRes := resourceData.defectdojoResource()
	ddWebhook := ddRes.(*notificationWebhookDefectdojoResource)
	var trd terraformResourceData = &resourceData
	populateDefectdojoResource(context.Background(), &diag.Diagnostics{}, trd, &ddRes)

	assert.Equal(t, *ddWebhook.Name, "My Webhook")
	assert.Equal(t, *ddWebhook.Url, "https://hooks.example.com/notify")
	assert.Equal(t, *ddWebhook.HeaderName, "X-Webhook-Secret")
}
