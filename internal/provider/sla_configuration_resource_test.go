package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// TestAccSlaConfigurationResource opens on a config that omits every SLA day
// count. DefectDojo fills them from its own model defaults and answers with
// them, so before critical/high/medium/low were marked Computed this failed on
// CREATE with "Provider produced inconsistent result after apply: .critical: was
// null, but now cty.NumberIntVal(...)". Every other SLA config sets all four
// explicitly, which is exactly why the defect went unnoticed - so the minimal
// config has to come first, or it becomes an update and proves nothing.
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
				Config: testAccSlaConfigurationResourceMinimalConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					// NotNull rather than a literal: the model defaults are
					// version-dependent, which is why no Default is declared
					// on these attributes.
					statecheck.ExpectKnownValue("defectdojo_sla_configuration.test", tfjsonpath.New("critical"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("defectdojo_sla_configuration.test", tfjsonpath.New("high"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("defectdojo_sla_configuration.test", tfjsonpath.New("medium"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("defectdojo_sla_configuration.test", tfjsonpath.New("low"), knownvalue.NotNull()),
				},
			},
			// Re-planning the same minimal config must be a no-op. This is what
			// proves Computed-without-Default is stable, rather than merely
			// surviving the first apply.
			{
				Config:   testAccSlaConfigurationResourceMinimalConfig(name),
				PlanOnly: true,
			},
			{
				ResourceName:      "defectdojo_sla_configuration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSlaConfigurationResourceConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_sla_configuration.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("defectdojo_sla_configuration.test", tfjsonpath.New("critical"), knownvalue.Int64Exact(7)),
				},
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

func testAccSlaConfigurationResourceMinimalConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_sla_configuration" "test" {
  name = %[1]q
}
`, name)
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
