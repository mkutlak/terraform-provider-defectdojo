package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCredentialResource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test-cred-%s", uniqueId())
	updatedName := fmt.Sprintf("updated-cred-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCredentialResourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_credential.test", "name", name),
					resource.TestCheckResourceAttr("defectdojo_credential.test", "username", "testuser"),
					resource.TestCheckResourceAttr("defectdojo_credential.test", "url", "https://example.com"),
					resource.TestCheckResourceAttr("defectdojo_credential.test", "role", "viewer"),
					resource.TestCheckResourceAttr("defectdojo_credential.test", "environment", "1"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_credential.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccCredentialResourceUpdatedConfig(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_credential.test", "name", updatedName),
					resource.TestCheckResourceAttr("defectdojo_credential.test", "username", "updateduser"),
				),
			},
		},
	})
}

func testAccCredentialResourceConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_credential" "test" {
  name        = %[1]q
  environment = 1
  username    = "testuser"
  role        = "viewer"
  url         = "https://example.com"
  description = "A test credential"
}
`, name)
}

func testAccCredentialResourceUpdatedConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_credential" "test" {
  name        = %[1]q
  environment = 1
  username    = "updateduser"
  role        = "viewer"
  url         = "https://updated.example.com"
  description = "Updated credential"
}
`, name)
}
