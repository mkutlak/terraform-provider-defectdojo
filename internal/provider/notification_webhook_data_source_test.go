package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccNotificationWebhookDataSource creates the resource, then reads it back through the
// data source and asserts every shared attribute matches. Both sides are
// populated by the same reflection engine, so any divergence is a mapping bug.
func TestAccNotificationWebhookDataSource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationWebhookResourceConfig(name) + `
data "defectdojo_notification_webhook" "test" {
  id = defectdojo_notification_webhook.test.id
}
`,
				Check: testAccCheckDataSourceMatchesResource(
					"data.defectdojo_notification_webhook.test", "defectdojo_notification_webhook.test"),
			},
		},
	})
}
