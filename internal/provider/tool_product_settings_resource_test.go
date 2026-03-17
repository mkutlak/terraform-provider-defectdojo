package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccToolProductSettingsResource(t *testing.T) {
	name := fmt.Sprintf("test-%s", uniqueId())
	updatedName := fmt.Sprintf("updated-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccToolProductSettingsResourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_tool_product_settings.test", "name", name),
					resource.TestCheckResourceAttr("defectdojo_tool_product_settings.test", "setting_url", "https://tool.example.com/project/test"),
					resource.TestCheckResourceAttrSet("defectdojo_tool_product_settings.test", "tool_configuration"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_tool_product_settings.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccToolProductSettingsResourceUpdatedConfig(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_tool_product_settings.test", "name", updatedName),
				),
			},
		},
	})
}

func testAccToolProductSettingsResourceConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_tool_type" "test_tool_type" {
  name = %[1]q
}
resource "defectdojo_tool_configuration" "test_tool_config" {
  name      = %[1]q
  tool_type = defectdojo_tool_type.test_tool_type.id
}
resource "defectdojo_product" "test_product" {
  name            = %[1]q
  description     = "test product for tool product settings"
  product_type_id = 1
}
resource "defectdojo_tool_product_settings" "test" {
  name               = %[1]q
  setting_url        = "https://tool.example.com/project/test"
  product            = defectdojo_product.test_product.id
  tool_configuration = defectdojo_tool_configuration.test_tool_config.id
  tool_project_id    = "test-project-id"
}
`, name)
}

func testAccToolProductSettingsResourceUpdatedConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_tool_type" "test_tool_type" {
  name = %[1]q
}
resource "defectdojo_tool_configuration" "test_tool_config" {
  name      = %[1]q
  tool_type = defectdojo_tool_type.test_tool_type.id
}
resource "defectdojo_product" "test_product" {
  name            = %[1]q
  description     = "test product for tool product settings"
  product_type_id = 1
}
resource "defectdojo_tool_product_settings" "test" {
  name               = %[1]q
  setting_url        = "https://tool.example.com/project/test"
  product            = defectdojo_product.test_product.id
  tool_configuration = defectdojo_tool_configuration.test_tool_config.id
  tool_project_id    = "updated-project-id"
}
`, name)
}
