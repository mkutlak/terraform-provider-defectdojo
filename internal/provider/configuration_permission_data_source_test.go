package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.defectdojo_configuration_permission.by_codename", tfjsonpath.New("codename"), knownvalue.StringExact("add_user")),
					statecheck.ExpectKnownValue("data.defectdojo_configuration_permission.by_codename", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.defectdojo_configuration_permission.by_codename", tfjsonpath.New("name"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccConfigurationPermissionDataSourceIdConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.defectdojo_configuration_permission.by_id", "id",
						"data.defectdojo_configuration_permission.by_codename", "id"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.defectdojo_configuration_permission.by_id", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.defectdojo_configuration_permission.by_id", tfjsonpath.New("name"), knownvalue.NotNull()),
				},
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
