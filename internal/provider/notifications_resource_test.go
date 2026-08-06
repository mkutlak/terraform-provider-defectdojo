package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
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
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccNotificationsResourceConfig(productName, []string{"alert", "mail"}, []string{"alert"}, "alert"),
				// product is Int64 while id is String, so this pair is compared as
				// flatmap strings; statecheck.CompareValuePairs is type-strict and
				// reports "271 != 271" for values that are equal but differently typed.
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("defectdojo_notifications.test", "product", "defectdojo_product.test", "id"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_notifications.test", tfjsonpath.New("scan_added_empty"), knownvalue.StringExact("alert")),
					statecheck.ExpectKnownValue("defectdojo_notifications.test", tfjsonpath.New("scan_added"), knownvalue.SetSizeExact(2)),
					statecheck.ExpectKnownValue("defectdojo_notifications.test", tfjsonpath.New("scan_added"), knownvalue.SetPartial([]knownvalue.Check{knownvalue.StringExact("alert"), knownvalue.StringExact("mail")})),
					statecheck.ExpectKnownValue("defectdojo_notifications.test", tfjsonpath.New("sla_breach"), knownvalue.SetSizeExact(1)),
					statecheck.ExpectKnownValue("defectdojo_notifications.test", tfjsonpath.New("sla_breach"), knownvalue.SetPartial([]knownvalue.Check{knownvalue.StringExact("alert")})),
				},
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_notifications.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccNotificationsResourceConfig(productName, []string{"mail"}, []string{"alert", "webhooks"}, "webhooks"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_notifications.test", tfjsonpath.New("scan_added_empty"), knownvalue.StringExact("webhooks")),
					statecheck.ExpectKnownValue("defectdojo_notifications.test", tfjsonpath.New("scan_added"), knownvalue.SetSizeExact(1)),
					statecheck.ExpectKnownValue("defectdojo_notifications.test", tfjsonpath.New("scan_added"), knownvalue.SetPartial([]knownvalue.Check{knownvalue.StringExact("mail")})),
					statecheck.ExpectKnownValue("defectdojo_notifications.test", tfjsonpath.New("sla_breach"), knownvalue.SetSizeExact(2)),
					statecheck.ExpectKnownValue("defectdojo_notifications.test", tfjsonpath.New("sla_breach"), knownvalue.SetPartial([]knownvalue.Check{knownvalue.StringExact("alert"), knownvalue.StringExact("webhooks")})),
				},
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
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationsDataSourceConfig(productName),
				// product is Int64 while id is String, so this pair is compared as
				// flatmap strings; statecheck.CompareValuePairs is type-strict and
				// reports "271 != 271" for values that are equal but differently typed.
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.defectdojo_notifications.test", "product", "defectdojo_product.test", "id"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs("data.defectdojo_notifications.test", tfjsonpath.New("id"), "defectdojo_notifications.test", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("data.defectdojo_notifications.test", tfjsonpath.New("scan_added_empty"), knownvalue.StringExact("mail")),
					statecheck.ExpectKnownValue("data.defectdojo_notifications.test", tfjsonpath.New("scan_added"), knownvalue.SetSizeExact(1)),
					statecheck.ExpectKnownValue("data.defectdojo_notifications.test", tfjsonpath.New("scan_added"), knownvalue.SetPartial([]knownvalue.Check{knownvalue.StringExact("alert")})),
				},
			},
		},
	})
}

func testAccNotificationsResourceConfig(productName string, scanAdded []string, slaBreach []string, scanAddedEmpty string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test" {
  name            = %[1]q
  description     = "test"
  product_type_id = 1
}
resource "defectdojo_notifications" "test" {
  product          = defectdojo_product.test.id
  scan_added       = %[2]s
  sla_breach       = %[3]s
  scan_added_empty = %[4]q
}
`, productName, quotedList(scanAdded), quotedList(slaBreach), scanAddedEmpty)
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
  product          = defectdojo_product.test.id
  scan_added       = ["alert"]
  scan_added_empty = "mail"
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
	quoted := make([]string, len(values))
	for i, v := range values {
		quoted[i] = fmt.Sprintf("%q", v)
	}
	return "[" + strings.Join(quoted, ", ") + "]"
}
