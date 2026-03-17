package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLanguageTypeResource(t *testing.T) {
	language := fmt.Sprintf("TestLang-%s", resource.UniqueId())
	updatedLanguage := fmt.Sprintf("UpdatedLang-%s", resource.UniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccLanguageTypeResourceConfig(language),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_language_type.test", "language", language),
				),
			},
			{
				ResourceName:      "defectdojo_language_type.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLanguageTypeResourceConfig(updatedLanguage),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_language_type.test", "language", updatedLanguage),
				),
			},
		},
	})
}

func testAccLanguageTypeResourceConfig(language string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_language_type" "test" {
  language = %[1]q
}
`, language)
}
