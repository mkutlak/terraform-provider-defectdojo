package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRiskAcceptanceResource(t *testing.T) {
	t.Skip("Skipped: requires accepted_findings which needs Finding resources created outside Terraform")
	name := fmt.Sprintf("test-%s", uniqueId())
	updatedName := fmt.Sprintf("updated-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccRiskAcceptanceResourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_risk_acceptance.test", "name", name),
					resource.TestCheckResourceAttr("defectdojo_risk_acceptance.test", "owner", "1"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "defectdojo_risk_acceptance.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"accepted_findings"},
			},
			// Update and Read testing
			{
				Config: testAccRiskAcceptanceResourceUpdatedConfig(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_risk_acceptance.test", "name", updatedName),
				),
			},
		},
	})
}

func testAccRiskAcceptanceResourceConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_risk_acceptance" "test" {
  name             = %[1]q
  owner            = 1
  accepted_findings = []
  accepted_by      = "security-team"
  decision         = "A"
}
`, name)
}

func testAccRiskAcceptanceResourceUpdatedConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_risk_acceptance" "test" {
  name             = %[1]q
  owner            = 1
  accepted_findings = []
  accepted_by      = "security-team"
  decision         = "A"
}
`, name)
}
