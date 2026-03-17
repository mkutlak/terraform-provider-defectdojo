package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProductGroupResource(t *testing.T) {
	groupName := fmt.Sprintf("prodgroup-%s", resource.UniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProductGroupResourceConfig(groupName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("defectdojo_product_group.test", "id"),
					resource.TestCheckResourceAttr("defectdojo_product_group.test", "product", "1"),
					resource.TestCheckResourceAttrSet("defectdojo_product_group.test", "group"),
					resource.TestCheckResourceAttr("defectdojo_product_group.test", "role", "3"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_product_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccProductGroupResourceConfig(groupName string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_dojo_group" "pg_group" {
  name = %[1]q
}
resource "defectdojo_product_group" "test" {
  product = 1
  group   = defectdojo_dojo_group.pg_group.id
  role    = 3
}
`, groupName)
}
