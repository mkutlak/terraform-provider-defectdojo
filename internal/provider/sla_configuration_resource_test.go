package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccSlaConfigurationResource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test-%s", uniqueId())
	updatedName := fmt.Sprintf("updated-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccSlaConfigurationResourceConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_sla_configuration.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
				},
			},
			{
				ResourceName:      "defectdojo_sla_configuration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSlaConfigurationResourceConfig(updatedName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_sla_configuration.test", tfjsonpath.New("name"), knownvalue.StringExact(updatedName)),
				},
			},
		},
	})
}

func testAccSlaConfigurationResourceConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_sla_configuration" "test" {
  name     = %[1]q
  critical = 7
  high     = 30
  medium   = 90
  low      = 180
}
`, name)
}
