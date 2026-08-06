package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccProductAPIScanConfigurationResource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProductAPIScanConfigurationResourceConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_product_api_scan_configuration.test", tfjsonpath.New("tool_configuration"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("defectdojo_product_api_scan_configuration.test", tfjsonpath.New("product"), knownvalue.NotNull()),
				},
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_product_api_scan_configuration.test", tfjsonpath.New("service_key_1"), knownvalue.StringExact("updated-key")),
				},
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
