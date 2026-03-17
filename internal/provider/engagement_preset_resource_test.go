package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEngagementPresetResource(t *testing.T) {
	name := fmt.Sprintf("test-%s", uniqueId())
	updatedTitle := fmt.Sprintf("updated-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccEngagementPresetResourceConfig(name, "Test Preset Title"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_engagement_preset.test", "title", "Test Preset Title"),
					resource.TestCheckResourceAttrSet("defectdojo_engagement_preset.test", "product"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_engagement_preset.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccEngagementPresetResourceConfig(name, updatedTitle),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_engagement_preset.test", "title", updatedTitle),
				),
			},
		},
	})
}

func testAccEngagementPresetResourceConfig(productName, title string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test_product" {
  name            = %[1]q
  description     = "test product for engagement preset"
  product_type_id = 1
}
resource "defectdojo_engagement_preset" "test" {
  product = defectdojo_product.test_product.id
  title   = %[2]q
  scope   = "All endpoints"
}
`, productName, title)
}
