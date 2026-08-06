package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccLanguageTypeResource(t *testing.T) {
	t.Parallel()
	language := fmt.Sprintf("TestLang-%s", uniqueId())
	updatedLanguage := fmt.Sprintf("UpdatedLang-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccLanguageTypeResourceConfig(language),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_language_type.test", tfjsonpath.New("language"), knownvalue.StringExact(language)),
				},
			},
			{
				ResourceName:      "defectdojo_language_type.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLanguageTypeResourceConfig(updatedLanguage),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_language_type.test", tfjsonpath.New("language"), knownvalue.StringExact(updatedLanguage)),
				},
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
