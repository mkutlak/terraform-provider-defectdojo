package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccUserContactInfoResource(t *testing.T) {
	t.Parallel()
	username := fmt.Sprintf("contacttest-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccUserContactInfoResourceConfig(username),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_user_contact_info.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("defectdojo_user_contact_info.test", tfjsonpath.New("user"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("defectdojo_user_contact_info.test", tfjsonpath.New("title"), knownvalue.StringExact("Dr.")),
					statecheck.ExpectKnownValue("defectdojo_user_contact_info.test", tfjsonpath.New("phone_number"), knownvalue.StringExact("+1234567890")),
					statecheck.ExpectKnownValue("defectdojo_user_contact_info.test", tfjsonpath.New("deduplication_execution_mode"), knownvalue.StringExact("async")),
				},
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_user_contact_info.test", tfjsonpath.New("title"), knownvalue.StringExact("Prof.")),
					statecheck.ExpectKnownValue("defectdojo_user_contact_info.test", tfjsonpath.New("phone_number"), knownvalue.StringExact("+0987654321")),
					statecheck.ExpectKnownValue("defectdojo_user_contact_info.test", tfjsonpath.New("deduplication_execution_mode"), knownvalue.StringExact("sync")),
				},
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
  user                         = defectdojo_user.contactinfo_user.id
  title                        = "Dr."
  phone_number                 = "+1234567890"
  deduplication_execution_mode = "async"
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
  user                         = defectdojo_user.contactinfo_user.id
  title                        = "Prof."
  phone_number                 = "+0987654321"
  deduplication_execution_mode = "sync"
}
`, username)
}
