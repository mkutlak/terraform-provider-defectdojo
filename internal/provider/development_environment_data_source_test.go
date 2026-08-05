package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccDevelopmentEnvironmentDataSource creates the resource, then reads it back through the
// data source and asserts every shared attribute matches. Both sides are
// populated by the same reflection engine, so any divergence is a mapping bug.
func TestAccDevelopmentEnvironmentDataSource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccDevelopmentEnvironmentResourceConfig(name) + `
data "defectdojo_development_environment" "test" {
  id = defectdojo_development_environment.test.id
}
`,
				Check: testAccCheckDataSourceMatchesResource(
					"data.defectdojo_development_environment.test", "defectdojo_development_environment.test"),
			},
		},
	})
}
