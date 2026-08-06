package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccNetworkLocationResource(t *testing.T) {
	t.Parallel()
	location := fmt.Sprintf("VPN-%s", uniqueId())
	updatedLocation := fmt.Sprintf("Internal-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkLocationResourceConfig(location),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_network_location.test", tfjsonpath.New("location"), knownvalue.StringExact(location)),
				},
			},
			{
				ResourceName:      "defectdojo_network_location.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNetworkLocationResourceConfig(updatedLocation),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_network_location.test", tfjsonpath.New("location"), knownvalue.StringExact(updatedLocation)),
				},
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
