package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccJiraInstanceIdDataSource(t *testing.T) {
	t.Parallel()
	t.Skip("Skipped: JiraInstance requires write-only password handling (not yet implemented)")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: `
provider "defectdojo" {}
data "defectdojo_jira_instance" "test" {
  id = "1"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.defectdojo_jira_instance.test", tfjsonpath.New("url"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccJiraInstanceUrlDataSource(t *testing.T) {
	t.Parallel()
	t.Skip("Skipped: JiraInstance requires write-only password handling (not yet implemented)")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: `
provider "defectdojo" {}
data "defectdojo_jira_instance" "test" {
  url = "https://jira.example.com"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.defectdojo_jira_instance.test", tfjsonpath.New("url"), knownvalue.StringExact("https://jira.example.com")),
					statecheck.ExpectKnownValue("data.defectdojo_jira_instance.test", tfjsonpath.New("id"), knownvalue.NotNull()),
				},
			},
		},
	})
}
