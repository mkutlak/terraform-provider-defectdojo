package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccCicdInfrastructureResource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test-%s", uniqueId())
	updatedName := fmt.Sprintf("updated-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCicdInfrastructureResourceConfig(name, "scm_server", "initial description"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_cicd_infrastructure.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("defectdojo_cicd_infrastructure.test", tfjsonpath.New("infrastructure_type"), knownvalue.StringExact("scm_server")),
					statecheck.ExpectKnownValue("defectdojo_cicd_infrastructure.test", tfjsonpath.New("description"), knownvalue.StringExact("initial description")),
				},
			},
			{
				ResourceName:      "defectdojo_cicd_infrastructure.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// In-place update: name and description change, infrastructure_type
				// stays put.
				Config: testAccCicdInfrastructureResourceConfig(updatedName, "scm_server", "updated description"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_cicd_infrastructure.test", tfjsonpath.New("name"), knownvalue.StringExact(updatedName)),
					statecheck.ExpectKnownValue("defectdojo_cicd_infrastructure.test", tfjsonpath.New("description"), knownvalue.StringExact("updated description")),
				},
			},
			{
				// infrastructure_type carries RequiresReplace() (see
				// cicd_infrastructure_resource.go): DefectDojo silently keeps the
				// original value on an in-place PUT, so the provider must destroy
				// and recreate instead of updating.
				Config: testAccCicdInfrastructureResourceConfig(updatedName, "build_server", "updated description"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("defectdojo_cicd_infrastructure.test", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_cicd_infrastructure.test", tfjsonpath.New("infrastructure_type"), knownvalue.StringExact("build_server")),
				},
			},
		},
	})
}

func testAccCicdInfrastructureResourceConfig(name, infrastructureType, description string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_cicd_infrastructure" "test" {
  name                 = %[1]q
  infrastructure_type  = %[2]q
  description          = %[3]q
}
`, name, infrastructureType, description)
}
