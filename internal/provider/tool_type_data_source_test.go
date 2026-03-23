package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccToolTypeIdDataSource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("dox-test-tt-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccToolTypeDataSourceIdConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.defectdojo_tool_type.test", "name", name),
				),
			},
		},
	})
}

func TestAccToolTypeNameDataSource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("dox-test-tt-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccToolTypeDataSourceNameConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.defectdojo_tool_type.test", "name", name),
					resource.TestCheckResourceAttrSet("data.defectdojo_tool_type.test", "id"),
				),
			},
		},
	})
}

func testAccToolTypeDataSourceIdConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_tool_type" "test" {
  name = %[1]q
}
data "defectdojo_tool_type" "test" {
  id         = defectdojo_tool_type.test.id
  depends_on = [defectdojo_tool_type.test]
}
`, name)
}

func testAccToolTypeDataSourceNameConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_tool_type" "test" {
  name = %[1]q
}
data "defectdojo_tool_type" "test" {
  name       = %[1]q
  depends_on = [defectdojo_tool_type.test]
}
`, name)
}
