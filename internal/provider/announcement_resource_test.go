package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// TestAccAnnouncementResource is intentionally NOT run in parallel: DefectDojo
// allows at most one Announcement instance-wide, and running this test
// concurrently with other announcement tests (or another test run against the
// same instance) would race on that single global slot.
func TestAccAnnouncementResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccAnnouncementResourceConfig("This is a test announcement", "info", true),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_announcement.test", tfjsonpath.New("message"), knownvalue.StringExact("This is a test announcement")),
					statecheck.ExpectKnownValue("defectdojo_announcement.test", tfjsonpath.New("style"), knownvalue.StringExact("info")),
					statecheck.ExpectKnownValue("defectdojo_announcement.test", tfjsonpath.New("dismissable"), knownvalue.Bool(true)),
				},
			},
			{
				ResourceName:      "defectdojo_announcement.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAnnouncementResourceConfig("This is an updated test announcement", "warning", false),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_announcement.test", tfjsonpath.New("message"), knownvalue.StringExact("This is an updated test announcement")),
					statecheck.ExpectKnownValue("defectdojo_announcement.test", tfjsonpath.New("style"), knownvalue.StringExact("warning")),
					statecheck.ExpectKnownValue("defectdojo_announcement.test", tfjsonpath.New("dismissable"), knownvalue.Bool(false)),
				},
			},
		},
	})
}

func testAccAnnouncementResourceConfig(message string, style string, dismissable bool) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_announcement" "test" {
  message     = %[1]q
  style       = %[2]q
  dismissable = %[3]t
}
`, message, style, dismissable)
}

// TestAccAnnouncementDataSource is intentionally NOT run in parallel for the
// same reason as TestAccAnnouncementResource: only one Announcement may
// exist instance-wide.
func TestAccAnnouncementDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccAnnouncementDataSourceConfig("This is a test announcement", "info", true),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.defectdojo_announcement.test", tfjsonpath.New("message"), knownvalue.StringExact("This is a test announcement")),
					statecheck.ExpectKnownValue("data.defectdojo_announcement.test", tfjsonpath.New("style"), knownvalue.StringExact("info")),
					statecheck.ExpectKnownValue("data.defectdojo_announcement.test", tfjsonpath.New("dismissable"), knownvalue.Bool(true)),
				},
			},
		},
	})
}

func testAccAnnouncementDataSourceConfig(message string, style string, dismissable bool) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_announcement" "test" {
  message     = %[1]q
  style       = %[2]q
  dismissable = %[3]t
}
data "defectdojo_announcement" "test" {
  id         = defectdojo_announcement.test.id
  depends_on = [defectdojo_announcement.test]
}
`, message, style, dismissable)
}
