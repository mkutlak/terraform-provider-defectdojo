package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProductMemberResource(t *testing.T) {
	t.Parallel()
	username := fmt.Sprintf("prodmember-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProductMemberResourceConfig(username),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("defectdojo_product_member.test", "id"),
					resource.TestCheckResourceAttrSet("defectdojo_product_member.test", "product"),
					resource.TestCheckResourceAttrSet("defectdojo_product_member.test", "user"),
					resource.TestCheckResourceAttr("defectdojo_product_member.test", "role", "3"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_product_member.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccProductMemberResourceConfig(username string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "pm_product" {
  name            = %[1]q
  description     = "test product for product member"
  product_type_id = 1
}
resource "defectdojo_user" "pm_user" {
  username = %[1]q
  email    = "prodmember@example.com"
  password = "TestPassword123!"
}
resource "defectdojo_product_member" "test" {
  product = defectdojo_product.pm_product.id
  user    = defectdojo_user.pm_user.id
  role    = 3
}
`, username)
}
