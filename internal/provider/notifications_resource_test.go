package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccNotificationsResource exercises a product-scoped
// defectdojo_notifications row: create, import, and update the notification
// channels.
func TestAccNotificationsResource(t *testing.T) {
	t.Parallel()
	productName := fmt.Sprintf("dox-test-repo-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccNotificationsResourceConfig(productName, []string{"alert", "mail"}, []string{"alert"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("defectdojo_notifications.test", "product", "defectdojo_product.test", "id"),
					resource.TestCheckResourceAttr("defectdojo_notifications.test", "scan_added.#", "2"),
					resource.TestCheckTypeSetElemAttr("defectdojo_notifications.test", "scan_added.*", "alert"),
					resource.TestCheckTypeSetElemAttr("defectdojo_notifications.test", "scan_added.*", "mail"),
					resource.TestCheckResourceAttr("defectdojo_notifications.test", "sla_breach.#", "1"),
					resource.TestCheckTypeSetElemAttr("defectdojo_notifications.test", "sla_breach.*", "alert"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_notifications.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccNotificationsResourceConfig(productName, []string{"mail"}, []string{"alert", "webhooks"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_notifications.test", "scan_added.#", "1"),
					resource.TestCheckTypeSetElemAttr("defectdojo_notifications.test", "scan_added.*", "mail"),
					resource.TestCheckResourceAttr("defectdojo_notifications.test", "sla_breach.#", "2"),
					resource.TestCheckTypeSetElemAttr("defectdojo_notifications.test", "sla_breach.*", "alert"),
					resource.TestCheckTypeSetElemAttr("defectdojo_notifications.test", "sla_breach.*", "webhooks"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccNotificationsDataSource(t *testing.T) {
	t.Parallel()
	productName := fmt.Sprintf("dox-test-repo-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationsDataSourceConfig(productName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.defectdojo_notifications.test", "id", "defectdojo_notifications.test", "id"),
					resource.TestCheckResourceAttrPair("data.defectdojo_notifications.test", "product", "defectdojo_product.test", "id"),
					resource.TestCheckResourceAttr("data.defectdojo_notifications.test", "scan_added.#", "1"),
					resource.TestCheckTypeSetElemAttr("data.defectdojo_notifications.test", "scan_added.*", "alert"),
				),
			},
		},
	})
}

func testAccNotificationsResourceConfig(productName string, scanAdded []string, slaBreach []string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test" {
  name            = %[1]q
  description     = "test"
  product_type_id = 1
}
resource "defectdojo_notifications" "test" {
  product    = defectdojo_product.test.id
  scan_added = %[2]s
  sla_breach = %[3]s
}
`, productName, quotedList(scanAdded), quotedList(slaBreach))
}

func testAccNotificationsDataSourceConfig(productName string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test" {
  name            = %[1]q
  description     = "test"
  product_type_id = 1
}
resource "defectdojo_notifications" "test" {
  product    = defectdojo_product.test.id
  scan_added = ["alert"]
}
data "defectdojo_notifications" "test" {
  id         = defectdojo_notifications.test.id
  depends_on = [defectdojo_notifications.test]
}
`, productName)
}

// quotedList renders a Go string slice as an HCL list-of-strings literal,
// e.g. []string{"alert", "mail"} -> `["alert", "mail"]`.
func quotedList(values []string) string {
	out := "["
	for i, v := range values {
		if i > 0 {
			out += ", "
		}
		out += fmt.Sprintf("%q", v)
	}
	out += "]"
	return out
}
