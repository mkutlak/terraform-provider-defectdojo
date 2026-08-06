package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccProductTypeResource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("dox-test-pt-%s", uniqueId())
	desc := fmt.Sprintf("dox test pt description %s", uniqueId())
	criticalProduct := "true"
	keyProduct := "true"
	updatedName := fmt.Sprintf("dox-new-pt-name-%s", uniqueId())
	updatedDesc := fmt.Sprintf("updated description %s", uniqueId())
	updatedCriticalProduct := "false"
	updatedKeyProduct := "false"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProductTypeResourceConfig(name, desc, criticalProduct, keyProduct),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_product_type.test", "critical_product", criticalProduct),
					resource.TestCheckResourceAttr("defectdojo_product_type.test", "key_product", keyProduct),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_product_type.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("defectdojo_product_type.test", tfjsonpath.New("description"), knownvalue.StringExact(desc)),
				},
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_product_type.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccProductTypeResourceConfig(updatedName, updatedDesc, updatedCriticalProduct, updatedKeyProduct),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_product_type.test", "critical_product", updatedCriticalProduct),
					resource.TestCheckResourceAttr("defectdojo_product_type.test", "key_product", updatedKeyProduct),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_product_type.test", tfjsonpath.New("name"), knownvalue.StringExact(updatedName)),
					statecheck.ExpectKnownValue("defectdojo_product_type.test", tfjsonpath.New("description"), knownvalue.StringExact(updatedDesc)),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccProductTypeResourceConfig(name string, desc string, criticalProduct string, keyProduct string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product_type" "test" {
  name = %[1]q
  description = %[2]q
  critical_product = %[3]q
  key_product = %[4]q
}
`, name, desc, criticalProduct, keyProduct)
}
