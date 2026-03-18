package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFindingTemplateResource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test-%s", uniqueId())
	updatedName := fmt.Sprintf("updated-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFindingTemplateResourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_finding_template.test", "title", name),
					resource.TestCheckResourceAttr("defectdojo_finding_template.test", "severity", "High"),
				),
			},
			{
				ResourceName:      "defectdojo_finding_template.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFindingTemplateResourceConfig(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_finding_template.test", "title", updatedName),
				),
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
