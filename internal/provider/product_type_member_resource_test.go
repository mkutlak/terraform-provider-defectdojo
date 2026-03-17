package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProductTypeMemberResource(t *testing.T) {
	username := fmt.Sprintf("ptmember-%s", resource.UniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProductTypeMemberResourceConfig(username),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("defectdojo_product_type_member.test", "id"),
					resource.TestCheckResourceAttr("defectdojo_product_type_member.test", "product_type", "1"),
					resource.TestCheckResourceAttrSet("defectdojo_product_type_member.test", "user"),
					resource.TestCheckResourceAttr("defectdojo_product_type_member.test", "role", "3"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_product_type_member.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccProductTypeMemberResourceConfig(username string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_user" "ptm_user" {
  username = %[1]q
  email    = "ptmember@example.com"
  password = "TestPassword123!"
}
resource "defectdojo_product_type_member" "test" {
  product_type = 1
  user         = defectdojo_user.ptm_user.id
  role         = 3
}
`, username)
}
