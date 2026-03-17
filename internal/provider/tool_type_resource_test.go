package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccToolTypeResource(t *testing.T) {
	name := fmt.Sprintf("test-%s", resource.UniqueId())
	updatedName := fmt.Sprintf("updated-%s", resource.UniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccToolTypeResourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_tool_type.test", "name", name),
				),
			},
			{
				ResourceName:      "defectdojo_tool_type.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccToolTypeResourceConfig(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_tool_type.test", "name", updatedName),
				),
			},
		},
	})
}

func testAccToolTypeResourceConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_tool_type" "test" {
  name = %[1]q
}
`, name)
}
