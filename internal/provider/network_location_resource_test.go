package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNetworkLocationResource(t *testing.T) {
	t.Parallel()
	location := fmt.Sprintf("VPN-%s", uniqueId())
	updatedLocation := fmt.Sprintf("Internal-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkLocationResourceConfig(location),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_network_location.test", "location", location),
				),
			},
			{
				ResourceName:      "defectdojo_network_location.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNetworkLocationResourceConfig(updatedLocation),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_network_location.test", "location", updatedLocation),
				),
			},
		},
	})
}

func testAccNetworkLocationResourceConfig(location string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_network_location" "test" {
  location = %[1]q
}
`, location)
}
