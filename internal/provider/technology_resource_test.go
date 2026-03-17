package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTechnologyResource(t *testing.T) {
	name := fmt.Sprintf("test-%s", resource.UniqueId())
	updatedName := fmt.Sprintf("updated-%s", resource.UniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTechnologyResourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_technology.test", "name", name),
					resource.TestCheckResourceAttrSet("defectdojo_technology.test", "product"),
					resource.TestCheckResourceAttrSet("defectdojo_technology.test", "user"),
				),
			},
			{
				ResourceName:      "defectdojo_technology.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTechnologyResourceConfig(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_technology.test", "name", updatedName),
				),
			},
		},
	})
}

func testAccTechnologyResourceConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}

resource "defectdojo_product_type" "test" {
  name = "technology-test-product-type"
}

resource "defectdojo_product" "test" {
  name        = "technology-test-product"
  description = "Test product for technology"
  prod_type   = defectdojo_product_type.test.id
}

data "defectdojo_user" "admin" {
  username = "admin"
}

resource "defectdojo_technology" "test" {
  name    = %[1]q
  product = defectdojo_product.test.id
  user    = data.defectdojo_user.admin.id
}
`, name)
}
