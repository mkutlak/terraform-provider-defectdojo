package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEngagementIdDataSource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("dox-test-eng-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEngagementDataSourceIdConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.defectdojo_engagement.test", "name", name),
				),
			},
		},
	})
}

func TestAccEngagementNameDataSource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("dox-test-eng-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEngagementDataSourceNameConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.defectdojo_engagement.test", "name", name),
					resource.TestCheckResourceAttrSet("data.defectdojo_engagement.test", "id"),
				),
			},
		},
	})
}

func testAccEngagementDataSourceIdConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test_product" {
  name            = %[1]q
  description     = "test product for engagement data source"
  product_type_id = 1
}
resource "defectdojo_engagement" "test" {
  product      = defectdojo_product.test_product.id
  target_start = "2025-01-01"
  target_end   = "2025-12-31"
  name         = %[1]q
}
data "defectdojo_engagement" "test" {
  id         = defectdojo_engagement.test.id
  depends_on = [defectdojo_engagement.test]
}
`, name)
}

func testAccEngagementDataSourceNameConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test_product" {
  name            = %[1]q
  description     = "test product for engagement data source"
  product_type_id = 1
}
resource "defectdojo_engagement" "test" {
  product      = defectdojo_product.test_product.id
  target_start = "2025-01-01"
  target_end   = "2025-12-31"
  name         = %[1]q
}
data "defectdojo_engagement" "test" {
  name       = %[1]q
  depends_on = [defectdojo_engagement.test]
}
`, name)
}
