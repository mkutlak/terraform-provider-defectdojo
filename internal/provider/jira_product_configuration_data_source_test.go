package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccJiraProductConfigurationDataSource creates the resource, then reads it back through the
// data source and asserts every shared attribute matches. Both sides are
// populated by the same reflection engine, so any divergence is a mapping bug.
func TestAccJiraProductConfigurationDataSource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("dox-test-repo-%s", uniqueId())
	jirakey := fmt.Sprintf("APPSEC%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccJiraProductConfigurationResourceConfig(name, jirakey) + `
data "defectdojo_jira_product_configuration" "test" {
  id = defectdojo_jira_product_configuration.test.id
}
`,
				Check: testAccCheckDataSourceMatchesResource(
					"data.defectdojo_jira_product_configuration.test", "defectdojo_jira_product_configuration.test"),
			},
		},
	})
}
