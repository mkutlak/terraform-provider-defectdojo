package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProductAPIScanConfigurationResource(t *testing.T) {
	name := fmt.Sprintf("test-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProductAPIScanConfigurationResourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("defectdojo_product_api_scan_configuration.test", "tool_configuration"),
					resource.TestCheckResourceAttrSet("defectdojo_product_api_scan_configuration.test", "product"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_product_api_scan_configuration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccProductAPIScanConfigurationResourceUpdatedConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_product_api_scan_configuration.test", "service_key_1", "updated-key"),
				),
			},
		},
	})
}

func testAccProductAPIScanConfigurationResourceConfig(name string) string {
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
  description     = "test product for api scan configuration"
  product_type_id = 1
}
resource "defectdojo_product_api_scan_configuration" "test" {
  product            = defectdojo_product.test_product.id
  tool_configuration = defectdojo_tool_configuration.test_tool_config.id
  service_key_1      = "test-service-key"
}
`, name)
}

func testAccProductAPIScanConfigurationResourceUpdatedConfig(name string) string {
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
  description     = "test product for api scan configuration"
  product_type_id = 1
}
resource "defectdojo_product_api_scan_configuration" "test" {
  product            = defectdojo_product.test_product.id
  tool_configuration = defectdojo_tool_configuration.test_tool_config.id
  service_key_1      = "updated-key"
}
`, name)
}
