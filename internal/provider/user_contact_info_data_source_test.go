package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccUserContactInfoDataSource creates the resource, then reads it back through the
// data source and asserts every shared attribute matches. Both sides are
// populated by the same reflection engine, so any divergence is a mapping bug.
func TestAccUserContactInfoDataSource(t *testing.T) {
	t.Parallel()
	username := fmt.Sprintf("contacttest-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccUserContactInfoResourceConfig(username) + `
data "defectdojo_user_contact_info" "test" {
  id = defectdojo_user_contact_info.test.id
}
`,
				Check: testAccCheckDataSourceMatchesResource(
					"data.defectdojo_user_contact_info.test", "defectdojo_user_contact_info.test"),
			},
		},
	})
}
