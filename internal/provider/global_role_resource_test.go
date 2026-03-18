package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGlobalRoleResource(t *testing.T) {
	t.Parallel()
	username := fmt.Sprintf("globalrole-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccGlobalRoleResourceConfig(username),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("defectdojo_global_role.test", "id"),
					resource.TestCheckResourceAttrSet("defectdojo_global_role.test", "user"),
					resource.TestCheckResourceAttr("defectdojo_global_role.test", "role", "4"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_global_role.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGlobalRoleResourceConfig(username string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_user" "gr_user" {
  username = %[1]q
  email    = "globalrole@example.com"
  password = "TestPassword123!"
}
resource "defectdojo_global_role" "test" {
  user = defectdojo_user.gr_user.id
  role = 4
}
`, username)
}
