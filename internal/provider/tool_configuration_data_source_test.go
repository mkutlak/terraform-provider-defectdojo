package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccToolConfigurationDataSource creates the resource, then reads it back through the
// data source and asserts every shared attribute matches. Both sides are
// populated by the same reflection engine, so any divergence is a mapping bug.
func TestAccToolConfigurationDataSource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccToolConfigurationResourceConfig(name) + `
data "defectdojo_tool_configuration" "test" {
  id = defectdojo_tool_configuration.test.id
}
`,
				Check: testAccCheckDataSourceMatchesResource(
					"data.defectdojo_tool_configuration.test", "defectdojo_tool_configuration.test"),
			},
		},
	})
}
