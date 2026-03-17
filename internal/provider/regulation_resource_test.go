package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRegulationResource(t *testing.T) {
	name := fmt.Sprintf("test-%s", uniqueId())
	updatedName := fmt.Sprintf("updated-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRegulationResourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_regulation.test", "name", name),
					resource.TestCheckResourceAttr("defectdojo_regulation.test", "acronym", "TST"),
					resource.TestCheckResourceAttr("defectdojo_regulation.test", "category", "other"),
					resource.TestCheckResourceAttr("defectdojo_regulation.test", "jurisdiction", "US"),
				),
			},
			{
				ResourceName:      "defectdojo_regulation.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRegulationResourceConfig(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_regulation.test", "name", updatedName),
				),
			},
		},
	})
}

func testAccRegulationResourceConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_regulation" "test" {
  name         = %[1]q
  acronym      = "TST"
  category     = "other"
  jurisdiction = "US"
}
`, name)
}
