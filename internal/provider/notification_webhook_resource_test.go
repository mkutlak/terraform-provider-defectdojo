package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationWebhookResource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test-%s", uniqueId())
	updatedName := fmt.Sprintf("updated-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccNotificationWebhookResourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_notification_webhook.test", "name", name),
					resource.TestCheckResourceAttr("defectdojo_notification_webhook.test", "url", "https://hooks.example.com/webhook"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_notification_webhook.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccNotificationWebhookResourceUpdatedConfig(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_notification_webhook.test", "name", updatedName),
				),
			},
		},
	})
}

func testAccNotificationWebhookResourceConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_notification_webhook" "test" {
  name = %[1]q
  url  = "https://hooks.example.com/webhook"
}
`, name)
}

func testAccNotificationWebhookResourceUpdatedConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_notification_webhook" "test" {
  name = %[1]q
  url  = "https://hooks.example.com/webhook"
}
`, name)
}
