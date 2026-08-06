package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccDevelopmentEnvironmentResource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test-%s", uniqueId())
	updatedName := fmt.Sprintf("updated-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccDevelopmentEnvironmentResourceConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_development_environment.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
				},
			},
			{
				ResourceName:      "defectdojo_development_environment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDevelopmentEnvironmentResourceConfig(updatedName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_development_environment.test", tfjsonpath.New("name"), knownvalue.StringExact(updatedName)),
				},
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
