package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccJiraInstanceIdDataSource(t *testing.T) {
	t.Parallel()
	t.Skip("Skipped: JiraInstance requires write-only password handling (not yet implemented)")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "defectdojo" {}
data "defectdojo_jira_instance" "test" {
  id = "1"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.defectdojo_jira_instance.test", "url"),
				),
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
		Steps: []resource.TestStep{
			{
				Config: `
provider "defectdojo" {}
data "defectdojo_jira_instance" "test" {
  url = "https://jira.example.com"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.defectdojo_jira_instance.test", "url", "https://jira.example.com"),
					resource.TestCheckResourceAttrSet("data.defectdojo_jira_instance.test", "id"),
				),
			},
		},
	})
}
