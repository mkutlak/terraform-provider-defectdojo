package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// testAccJiraInstanceDataSourceConfig seeds a jira instance and looks it up
// through the data source, keyed by whichever attribute the caller names.
//
// The lookup deliberately goes through the created resource's own attributes
// rather than a hardcoded id or url: these tests previously assumed an instance
// was already present on the server, so they could only ever have passed
// against a pre-seeded database.
func testAccJiraInstanceDataSourceConfig(name string, url string, lookup string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_jira_instance" "seed" {
  url                       = %[2]q
  username                  = "testuser"
  password                  = "testpassword"
  configuration_name        = %[1]q
  epic_name_id              = 10001
  open_status_key           = 11
  close_status_key          = 21
  info_mapping_severity     = "Trivial"
  low_mapping_severity      = "Minor"
  medium_mapping_severity   = "Major"
  high_mapping_severity     = "Critical"
  critical_mapping_severity = "Blocker"
}
data "defectdojo_jira_instance" "test" {
  %[3]s
}
`, name, url, lookup)
}

func TestAccJiraInstanceIdDataSource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test-ds-id-%s", uniqueId())
	url := fmt.Sprintf("https://jira-%s.example.com", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccJiraInstanceDataSourceConfig(name, url, "id = defectdojo_jira_instance.seed.id"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.defectdojo_jira_instance.test", tfjsonpath.New("url"), knownvalue.StringExact(url)),
					statecheck.ExpectKnownValue("data.defectdojo_jira_instance.test", tfjsonpath.New("configuration_name"), knownvalue.StringExact(name)),
				},
			},
		},
	})
}

func TestAccJiraInstanceUrlDataSource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test-ds-url-%s", uniqueId())
	url := fmt.Sprintf("https://jira-%s.example.com", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccJiraInstanceDataSourceConfig(name, url, "url = defectdojo_jira_instance.seed.url"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.defectdojo_jira_instance.test", tfjsonpath.New("url"), knownvalue.StringExact(url)),
					statecheck.ExpectKnownValue("data.defectdojo_jira_instance.test", tfjsonpath.New("id"), knownvalue.NotNull()),
				},
			},
		},
	})
}
