package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEngagementResource(t *testing.T) {
	name := fmt.Sprintf("test-%s", uniqueId())
	updatedName := fmt.Sprintf("updated-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccEngagementResourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_engagement.test", "name", name),
					resource.TestCheckResourceAttr("defectdojo_engagement.test", "target_start", "2025-01-01"),
					resource.TestCheckResourceAttr("defectdojo_engagement.test", "target_end", "2025-12-31"),
					resource.TestCheckResourceAttrSet("defectdojo_engagement.test", "product"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_engagement.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccEngagementResourceUpdatedConfig(name, updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_engagement.test", "name", updatedName),
					resource.TestCheckResourceAttr("defectdojo_engagement.test", "target_start", "2025-06-01"),
					resource.TestCheckResourceAttr("defectdojo_engagement.test", "target_end", "2025-12-31"),
				),
			},
		},
	})
}

func testAccEngagementResourceConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test_product" {
  name            = %[1]q
  description     = "test product for engagement"
  product_type_id = 1
}
resource "defectdojo_engagement" "test" {
  product      = defectdojo_product.test_product.id
  target_start = "2025-01-01"
  target_end   = "2025-12-31"
  name         = %[1]q
}
`, name)
}

func testAccEngagementResourceUpdatedConfig(productName, engagementName string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test_product" {
  name            = %[1]q
  description     = "test product for engagement"
  product_type_id = 1
}
resource "defectdojo_engagement" "test" {
  product      = defectdojo_product.test_product.id
  target_start = "2025-06-01"
  target_end   = "2025-12-31"
  name         = %[2]q
}
`, productName, engagementName)
}
