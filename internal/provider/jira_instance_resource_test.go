package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccJiraInstanceResource(t *testing.T) {
	t.Skip("Skipped: JiraInstance requires write-only password handling (not yet implemented)")
	name := fmt.Sprintf("test-%s", uniqueId())
	updatedName := fmt.Sprintf("updated-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccJiraInstanceResourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_jira_instance.test", "configuration_name", name),
					resource.TestCheckResourceAttr("defectdojo_jira_instance.test", "url", "https://jira.example.com"),
					resource.TestCheckResourceAttr("defectdojo_jira_instance.test", "username", "testuser"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "defectdojo_jira_instance.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			// Update and Read testing
			{
				Config: testAccJiraInstanceResourceUpdatedConfig(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_jira_instance.test", "configuration_name", updatedName),
				),
			},
		},
	})
}

func testAccJiraInstanceResourceConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_jira_instance" "test" {
  url                      = "https://jira.example.com"
  username                 = "testuser"
  password                 = "testpassword"
  configuration_name       = %[1]q
  epic_name_id             = 10001
  open_status_key          = 11
  close_status_key         = 21
  info_mapping_severity    = "Trivial"
  low_mapping_severity     = "Minor"
  medium_mapping_severity  = "Major"
  high_mapping_severity    = "Critical"
  critical_mapping_severity = "Blocker"
}
`, name)
}

func testAccJiraInstanceResourceUpdatedConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_jira_instance" "test" {
  url                      = "https://jira.example.com"
  username                 = "testuser"
  password                 = "testpassword"
  configuration_name       = %[1]q
  epic_name_id             = 10001
  open_status_key          = 11
  close_status_key         = 21
  info_mapping_severity    = "Trivial"
  low_mapping_severity     = "Minor"
  medium_mapping_severity  = "Major"
  high_mapping_severity    = "Critical"
  critical_mapping_severity = "Blocker"
}
`, name)
}
