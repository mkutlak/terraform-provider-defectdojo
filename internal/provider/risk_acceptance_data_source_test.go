package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccRiskAcceptanceDataSource creates the resource, then reads it back through the
// data source and asserts every shared attribute matches. Both sides are
// populated by the same reflection engine, so any divergence is a mapping bug.
func TestAccRiskAcceptanceDataSource(t *testing.T) {
	t.Parallel()
	t.Skip("Skipped: requires accepted_findings which needs Finding resources created outside Terraform")
	name := fmt.Sprintf("test-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccRiskAcceptanceResourceConfig(name) + `
data "defectdojo_risk_acceptance" "test" {
  id = defectdojo_risk_acceptance.test.id
}
`,
				Check: testAccCheckDataSourceMatchesResource(
					"data.defectdojo_risk_acceptance.test", "defectdojo_risk_acceptance.test"),
			},
		},
	})
}
