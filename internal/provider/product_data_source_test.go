package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProductDataSource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("dox-test-repo-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccProductDataSourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.defectdojo_product.test", "name", name),
					resource.TestCheckResourceAttr("data.defectdojo_product.test", "description", "test"),
				),
			},
		},
	})
}

func TestAccProductNameDataSource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("dox-test-repo-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProductDataSourceNameConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.defectdojo_product.test", "name", name),
					resource.TestCheckResourceAttr("data.defectdojo_product.test", "description", "test"),
				),
			},
		},
	})
}

func TestAccProductDataSourceConflict(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("dox-test-repo-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccProductDataSourceConflictConfig(name),
				ExpectError: regexp.MustCompile("Conflicting Lookup Parameters"),
			},
		},
	})
}

func testAccProductDataSourceConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test" {
  name = %[1]q
	description = "test"
  product_type_id = 1
}
data "defectdojo_product" "test" {
  id = defectdojo_product.test.id
  depends_on = [defectdojo_product.test]
}
`, name)
}

func testAccProductDataSourceNameConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test" {
  name            = %[1]q
  description     = "test"
  product_type_id = 1
}
data "defectdojo_product" "test" {
  name       = %[1]q
  depends_on = [defectdojo_product.test]
}
`, name)
}

func testAccProductDataSourceConflictConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test" {
  name            = %[1]q
  description     = "test"
  product_type_id = 1
}
data "defectdojo_product" "test" {
  id         = defectdojo_product.test.id
  name       = %[1]q
  depends_on = [defectdojo_product.test]
}
`, name)
}
