package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSlaConfigurationResource(t *testing.T) {
	name := fmt.Sprintf("test-%s", uniqueId())
	updatedName := fmt.Sprintf("updated-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSlaConfigurationResourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_sla_configuration.test", "name", name),
				),
			},
			{
				ResourceName:      "defectdojo_sla_configuration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSlaConfigurationResourceConfig(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_sla_configuration.test", "name", updatedName),
				),
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
