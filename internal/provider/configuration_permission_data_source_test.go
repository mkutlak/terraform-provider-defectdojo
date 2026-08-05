package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccConfigurationPermissionDataSource exercises the read-only
// defectdojo_configuration_permission data source. Configuration permissions
// are DefectDojo's own built-in permission registry (not something the
// provider can create), so this test looks up a permission that is known to
// exist in every DefectDojo instance: "add_user", a stable Django
// permission codename.
func TestAccConfigurationPermissionDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigurationPermissionDataSourceCodenameConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.defectdojo_configuration_permission.by_codename", "codename", "add_user"),
					resource.TestCheckResourceAttrSet("data.defectdojo_configuration_permission.by_codename", "id"),
					resource.TestCheckResourceAttrSet("data.defectdojo_configuration_permission.by_codename", "name"),
				),
			},
			{
				Config: testAccConfigurationPermissionDataSourceIdConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.defectdojo_configuration_permission.by_id", "id"),
					resource.TestCheckResourceAttrSet("data.defectdojo_configuration_permission.by_id", "name"),
					resource.TestCheckResourceAttrPair(
						"data.defectdojo_configuration_permission.by_id", "id",
						"data.defectdojo_configuration_permission.by_codename", "id"),
				),
			},
		},
	})
}

func testAccConfigurationPermissionDataSourceCodenameConfig() string {
	return `
provider "defectdojo" {}
data "defectdojo_configuration_permission" "by_codename" {
  codename = "add_user"
}
`
}

func testAccConfigurationPermissionDataSourceIdConfig() string {
	return `
provider "defectdojo" {}
data "defectdojo_configuration_permission" "by_codename" {
  codename = "add_user"
}
data "defectdojo_configuration_permission" "by_id" {
  id = data.defectdojo_configuration_permission.by_codename.id
}
`
}
