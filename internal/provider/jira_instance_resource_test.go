package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccJiraInstanceResource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test-%s", uniqueId())
	updatedName := fmt.Sprintf("updated-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccJiraInstanceResourceConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_jira_instance.test", tfjsonpath.New("configuration_name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("defectdojo_jira_instance.test", tfjsonpath.New("url"), knownvalue.StringExact("https://jira.example.com")),
					statecheck.ExpectKnownValue("defectdojo_jira_instance.test", tfjsonpath.New("username"), knownvalue.StringExact("testuser")),
				},
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_jira_instance.test", tfjsonpath.New("configuration_name"), knownvalue.StringExact(updatedName)),
				},
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
