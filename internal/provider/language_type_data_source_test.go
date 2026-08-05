package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccLanguageTypeDataSource creates the resource, then reads it back through the
// data source and asserts every shared attribute matches. Both sides are
// populated by the same reflection engine, so any divergence is a mapping bug.
func TestAccLanguageTypeDataSource(t *testing.T) {
	t.Parallel()
	language := fmt.Sprintf("TestLang-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccLanguageTypeResourceConfig(language) + `
data "defectdojo_language_type" "test" {
  id = defectdojo_language_type.test.id
}
`,
				Check: testAccCheckDataSourceMatchesResource(
					"data.defectdojo_language_type.test", "defectdojo_language_type.test"),
			},
		},
	})
}
