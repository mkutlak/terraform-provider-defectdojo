package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccNetworkLocationDataSource creates the resource, then reads it back through the
// data source and asserts every shared attribute matches. Both sides are
// populated by the same reflection engine, so any divergence is a mapping bug.
func TestAccNetworkLocationDataSource(t *testing.T) {
	t.Parallel()
	location := fmt.Sprintf("VPN-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkLocationResourceConfig(location) + `
data "defectdojo_network_location" "test" {
  id = defectdojo_network_location.test.id
}
`,
				Check: testAccCheckDataSourceMatchesResource(
					"data.defectdojo_network_location.test", "defectdojo_network_location.test"),
			},
		},
	})
}
