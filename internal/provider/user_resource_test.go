package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccUserResource(t *testing.T) {
	t.Parallel()
	username := fmt.Sprintf("testuser-%s", uniqueId())
	updatedUsername := fmt.Sprintf("updated-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccUserResourceConfig(username),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_user.test", tfjsonpath.New("username"), knownvalue.StringExact(username)),
					statecheck.ExpectKnownValue("defectdojo_user.test", tfjsonpath.New("email"), knownvalue.StringExact("test@example.com")),
					statecheck.ExpectKnownValue("defectdojo_user.test", tfjsonpath.New("first_name"), knownvalue.StringExact("Test")),
					statecheck.ExpectKnownValue("defectdojo_user.test", tfjsonpath.New("last_name"), knownvalue.StringExact("User")),
					statecheck.ExpectKnownValue("defectdojo_user.test", tfjsonpath.New("is_staff"), knownvalue.Bool(true)),
				},
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_user.test", tfjsonpath.New("username"), knownvalue.StringExact(updatedUsername)),
					statecheck.ExpectKnownValue("defectdojo_user.test", tfjsonpath.New("email"), knownvalue.StringExact("updated@example.com")),
					statecheck.ExpectKnownValue("defectdojo_user.test", tfjsonpath.New("is_staff"), knownvalue.Bool(false)),
				},
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
  is_staff   = true
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
  is_staff   = false
}
`, username)
}
