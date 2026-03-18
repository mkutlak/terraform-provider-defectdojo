package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserContactInfoResource(t *testing.T) {
	t.Parallel()
	username := fmt.Sprintf("contacttest-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccUserContactInfoResourceConfig(username),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("defectdojo_user_contact_info.test", "id"),
					resource.TestCheckResourceAttrSet("defectdojo_user_contact_info.test", "user"),
					resource.TestCheckResourceAttr("defectdojo_user_contact_info.test", "title", "Dr."),
					resource.TestCheckResourceAttr("defectdojo_user_contact_info.test", "phone_number", "+1234567890"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_user_contact_info.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccUserContactInfoResourceUpdatedConfig(username),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_user_contact_info.test", "title", "Prof."),
					resource.TestCheckResourceAttr("defectdojo_user_contact_info.test", "phone_number", "+0987654321"),
				),
			},
		},
	})
}

func testAccUserContactInfoResourceConfig(username string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_user" "contactinfo_user" {
  username = %[1]q
  email    = "contactinfo@example.com"
  password = "TestPassword123!"
}
resource "defectdojo_user_contact_info" "test" {
  user         = defectdojo_user.contactinfo_user.id
  title        = "Dr."
  phone_number = "+1234567890"
}
`, username)
}

func testAccUserContactInfoResourceUpdatedConfig(username string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_user" "contactinfo_user" {
  username = %[1]q
  email    = "contactinfo@example.com"
  password = "TestPassword123!"
}
resource "defectdojo_user_contact_info" "test" {
  user         = defectdojo_user.contactinfo_user.id
  title        = "Prof."
  phone_number = "+0987654321"
}
`, username)
}
