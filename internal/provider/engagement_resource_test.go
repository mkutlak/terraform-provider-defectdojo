package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccEngagementResource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test-%s", uniqueId())
	updatedName := fmt.Sprintf("updated-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccEngagementResourceConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_engagement.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("defectdojo_engagement.test", tfjsonpath.New("target_start"), knownvalue.StringExact("2025-01-01")),
					statecheck.ExpectKnownValue("defectdojo_engagement.test", tfjsonpath.New("target_end"), knownvalue.StringExact("2025-12-31")),
					statecheck.ExpectKnownValue("defectdojo_engagement.test", tfjsonpath.New("product"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("defectdojo_engagement.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
				},
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_engagement.test", tfjsonpath.New("name"), knownvalue.StringExact(updatedName)),
					statecheck.ExpectKnownValue("defectdojo_engagement.test", tfjsonpath.New("target_start"), knownvalue.StringExact("2025-06-01")),
					statecheck.ExpectKnownValue("defectdojo_engagement.test", tfjsonpath.New("target_end"), knownvalue.StringExact("2025-12-31")),
				},
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
