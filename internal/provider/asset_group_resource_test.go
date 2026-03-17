package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAssetGroupResource(t *testing.T) {
	name := fmt.Sprintf("asset-group-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccAssetGroupResourceConfig(name, 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("defectdojo_asset_group.test", "asset"),
					resource.TestCheckResourceAttrSet("defectdojo_asset_group.test", "group"),
					resource.TestCheckResourceAttr("defectdojo_asset_group.test", "role", "1"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_asset_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccAssetGroupResourceConfig(name, 2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_asset_group.test", "role", "2"),
				),
			},
		},
	})
}

func testAccAssetGroupResourceConfig(name string, role int) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "ag_product" {
  name            = %[1]q
  description     = "test product for asset group"
  product_type_id = 1
}
resource "defectdojo_dojo_group" "ag_group" {
  name = %[1]q
}
resource "defectdojo_asset_group" "test" {
  asset = defectdojo_product.ag_product.id
  group = defectdojo_dojo_group.ag_group.id
  role  = %[2]d
}
`, name, role)
}
