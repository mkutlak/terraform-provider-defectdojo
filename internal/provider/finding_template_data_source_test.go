package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccFindingTemplateDataSource creates the resource, then reads it back through the
// data source and asserts every shared attribute matches. Both sides are
// populated by the same reflection engine, so any divergence is a mapping bug.
func TestAccFindingTemplateDataSource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccFindingTemplateResourceConfig(name) + `
data "defectdojo_finding_template" "test" {
  id = defectdojo_finding_template.test.id
}
`,
				Check: testAccCheckDataSourceMatchesResource(
					"data.defectdojo_finding_template.test", "defectdojo_finding_template.test"),
			},
		},
	})
}
