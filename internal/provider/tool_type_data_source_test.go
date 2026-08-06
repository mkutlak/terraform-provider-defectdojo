package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccToolTypeIdDataSource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("dox-test-tt-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccToolTypeDataSourceIdConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.defectdojo_tool_type.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
				},
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
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccToolTypeDataSourceNameConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.defectdojo_tool_type.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.defectdojo_tool_type.test", tfjsonpath.New("id"), knownvalue.NotNull()),
				},
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
