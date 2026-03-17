package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDevelopmentEnvironmentResource(t *testing.T) {
	name := fmt.Sprintf("test-%s", uniqueId())
	updatedName := fmt.Sprintf("updated-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDevelopmentEnvironmentResourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_development_environment.test", "name", name),
				),
			},
			{
				ResourceName:      "defectdojo_development_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDevelopmentEnvironmentResourceConfig(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_development_environment.test", "name", updatedName),
				),
			},
		},
	})
}

func testAccDevelopmentEnvironmentResourceConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_development_environment" "test" {
  name = %[1]q
}
`, name)
}
