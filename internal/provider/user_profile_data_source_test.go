package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccUserProfileDataSource exercises the read-only defectdojo_user_profile
// data source. GET /api/v2/user_profile/ always describes the authenticated
// user - the caller identified by the provider's credentials - so no seeding
// is required; acceptance tests authenticate as the "admin" user.
func TestAccUserProfileDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "defectdojo" {}
data "defectdojo_user_profile" "me" {}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.defectdojo_user_profile.me", "id"),
					resource.TestCheckResourceAttr("data.defectdojo_user_profile.me", "username", "admin"),
					resource.TestCheckResourceAttr("data.defectdojo_user_profile.me", "is_superuser", "true"),
				),
			},
		},
	})
}
