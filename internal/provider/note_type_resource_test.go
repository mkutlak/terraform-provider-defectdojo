package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNoteTypeResource(t *testing.T) {
	name := fmt.Sprintf("test-%s", uniqueId())
	updatedName := fmt.Sprintf("updated-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNoteTypeResourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_note_type.test", "name", name),
					resource.TestCheckResourceAttr("defectdojo_note_type.test", "description", "Test note type description"),
				),
			},
			{
				ResourceName:      "defectdojo_note_type.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNoteTypeResourceConfig(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_note_type.test", "name", updatedName),
				),
			},
		},
	})
}

func testAccNoteTypeResourceConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_note_type" "test" {
  name        = %[1]q
  description = "Test note type description"
}
`, name)
}
