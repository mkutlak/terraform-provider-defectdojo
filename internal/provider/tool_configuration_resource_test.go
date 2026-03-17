package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccToolConfigurationResource(t *testing.T) {
	name := fmt.Sprintf("test-%s", uniqueId())
	updatedName := fmt.Sprintf("updated-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccToolConfigurationResourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_tool_configuration.test", "name", name),
				),
			},
			{
				ResourceName:      "defectdojo_tool_configuration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccToolConfigurationResourceConfig(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_tool_configuration.test", "name", updatedName),
				),
			},
		},
	})
}

func testAccToolConfigurationResourceConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_tool_type" "test_type" {
  name = %[1]q
}
resource "defectdojo_tool_configuration" "test" {
  name      = %[1]q
  tool_type = defectdojo_tool_type.test_type.id
}
`, name)
}
