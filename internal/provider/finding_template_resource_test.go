package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccFindingTemplateResource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test-%s", uniqueId())
	updatedName := fmt.Sprintf("updated-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccFindingTemplateResourceConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_finding_template.test", tfjsonpath.New("title"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("defectdojo_finding_template.test", tfjsonpath.New("severity"), knownvalue.StringExact("High")),
				},
			},
			{
				ResourceName:      "defectdojo_finding_template.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFindingTemplateResourceConfig(updatedName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_finding_template.test", tfjsonpath.New("title"), knownvalue.StringExact(updatedName)),
				},
			},
		},
	})
}

func testAccFindingTemplateResourceConfig(title string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_finding_template" "test" {
  title       = %[1]q
  severity    = "High"
  description = "A test finding template description"
}
`, title)
}
