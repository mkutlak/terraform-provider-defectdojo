package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRegulationResource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test-%s", uniqueId())
	updatedName := fmt.Sprintf("updated-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccRegulationResourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_regulation.test", "name", name),
					resource.TestCheckResourceAttr("defectdojo_regulation.test", "acronym", testAccRegulationAcronym(name)),
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
  acronym      = %[2]q
  category     = "other"
  jurisdiction = "US"
}
`, name, testAccRegulationAcronym(name))
}

// testAccRegulationAcronym derives a unique acronym from an already-unique
// name. acronym is unique-constrained server-side and capped at 20 characters,
// so a hardcoded value makes the resource test and the data source test collide
// whenever they run concurrently ("acronym already exists").
func testAccRegulationAcronym(name string) string {
	acronym := strings.ToUpper(strings.ReplaceAll(name, "-", ""))
	if len(acronym) > 20 {
		acronym = acronym[len(acronym)-20:]
	}
	return acronym
}
