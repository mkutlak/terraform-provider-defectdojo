package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserResource(t *testing.T) {
	username := fmt.Sprintf("testuser-%s", resource.UniqueId())
	updatedUsername := fmt.Sprintf("updated-%s", resource.UniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccUserResourceConfig(username),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_user.test", "username", username),
					resource.TestCheckResourceAttr("defectdojo_user.test", "email", "test@example.com"),
					resource.TestCheckResourceAttr("defectdojo_user.test", "first_name", "Test"),
					resource.TestCheckResourceAttr("defectdojo_user.test", "last_name", "User"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "defectdojo_user.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			// Update and Read testing
			{
				Config: testAccUserResourceUpdatedConfig(updatedUsername),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_user.test", "username", updatedUsername),
					resource.TestCheckResourceAttr("defectdojo_user.test", "email", "updated@example.com"),
				),
			},
		},
	})
}

func testAccUserResourceConfig(username string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_user" "test" {
  username   = %[1]q
  email      = "test@example.com"
  first_name = "Test"
  last_name  = "User"
  password   = "TestPassword123!"
}
`, username)
}

func testAccUserResourceUpdatedConfig(username string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_user" "test" {
  username   = %[1]q
  email      = "updated@example.com"
  first_name = "Updated"
  last_name  = "Name"
  password   = "UpdatedPassword123!"
}
`, username)
}
