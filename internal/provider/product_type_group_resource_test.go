package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProductTypeGroupResource(t *testing.T) {
	t.Parallel()
	groupName := fmt.Sprintf("ptgroup-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProductTypeGroupResourceConfig(groupName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("defectdojo_product_type_group.test", "id"),
					resource.TestCheckResourceAttr("defectdojo_product_type_group.test", "product_type", "1"),
					resource.TestCheckResourceAttrSet("defectdojo_product_type_group.test", "group"),
					resource.TestCheckResourceAttr("defectdojo_product_type_group.test", "role", "3"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_product_type_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccProductTypeGroupResourceConfig(groupName string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_dojo_group" "ptg_group" {
  name = %[1]q
}
resource "defectdojo_product_type_group" "test" {
  product_type = 1
  group        = defectdojo_dojo_group.ptg_group.id
  role         = 3
}
`, groupName)
}
