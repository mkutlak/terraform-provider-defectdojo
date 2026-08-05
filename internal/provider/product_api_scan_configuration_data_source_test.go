package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccProductAPIScanConfigurationDataSource creates the resource, then reads it back through the
// data source and asserts every shared attribute matches. Both sides are
// populated by the same reflection engine, so any divergence is a mapping bug.
func TestAccProductAPIScanConfigurationDataSource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccProductAPIScanConfigurationResourceConfig(name) + `
data "defectdojo_product_api_scan_configuration" "test" {
  id = defectdojo_product_api_scan_configuration.test.id
}
`,
				Check: testAccCheckDataSourceMatchesResource(
					"data.defectdojo_product_api_scan_configuration.test", "defectdojo_product_api_scan_configuration.test"),
			},
		},
	})
}
