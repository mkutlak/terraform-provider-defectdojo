package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDojoGroupIdDataSource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("dox-test-group-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDojoGroupDataSourceIdConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.defectdojo_dojo_group.test", "name", name),
				),
			},
		},
	})
}

func TestAccDojoGroupNameDataSource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("dox-test-group-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDojoGroupDataSourceNameConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.defectdojo_dojo_group.test", "name", name),
					resource.TestCheckResourceAttrSet("data.defectdojo_dojo_group.test", "id"),
				),
			},
		},
	})
}

func testAccDojoGroupDataSourceIdConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_dojo_group" "test" {
  name = %[1]q
}
data "defectdojo_dojo_group" "test" {
  id         = defectdojo_dojo_group.test.id
  depends_on = [defectdojo_dojo_group.test]
}
`, name)
}

func testAccDojoGroupDataSourceNameConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_dojo_group" "test" {
  name = %[1]q
}
data "defectdojo_dojo_group" "test" {
  name       = %[1]q
  depends_on = [defectdojo_dojo_group.test]
}
`, name)
}
