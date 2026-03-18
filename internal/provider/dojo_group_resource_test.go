package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDojoGroupResource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test-group-%s", uniqueId())
	updatedName := fmt.Sprintf("updated-group-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDojoGroupResourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_dojo_group.test", "name", name),
					resource.TestCheckResourceAttr("defectdojo_dojo_group.test", "description", "A test group"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_dojo_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccDojoGroupResourceUpdatedConfig(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_dojo_group.test", "name", updatedName),
					resource.TestCheckResourceAttr("defectdojo_dojo_group.test", "description", "Updated description"),
				),
			},
		},
	})
}

func testAccDojoGroupResourceConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_dojo_group" "test" {
  name        = %[1]q
  description = "A test group"
}
`, name)
}

func testAccDojoGroupResourceUpdatedConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_dojo_group" "test" {
  name        = %[1]q
  description = "Updated description"
}
`, name)
}
