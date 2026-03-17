package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccStubFindingResource(t *testing.T) {
	name := fmt.Sprintf("test-%s", resource.UniqueId())
	updatedName := fmt.Sprintf("updated-%s", resource.UniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStubFindingResourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_stub_finding.test", "title", name),
					resource.TestCheckResourceAttr("defectdojo_stub_finding.test", "severity", "Medium"),
				),
			},
			{
				ResourceName:      "defectdojo_stub_finding.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStubFindingResourceConfig(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_stub_finding.test", "title", updatedName),
				),
			},
		},
	})
}

func testAccStubFindingResourceConfig(title string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}

resource "defectdojo_product_type" "test" {
  name = "stub-finding-test-product-type"
}

resource "defectdojo_product" "test" {
  name        = "stub-finding-test-product"
  description = "Test product for stub finding"
  prod_type   = defectdojo_product_type.test.id
}

resource "defectdojo_engagement" "test" {
  name         = "stub-finding-test-engagement"
  product      = defectdojo_product.test.id
  target_start = "2023-01-01"
  target_end   = "2023-12-31"
}

resource "defectdojo_dd_test" "test" {
  engagement   = defectdojo_engagement.test.id
  test_type    = 1
  target_start = "2023-01-01T00:00:00Z"
  target_end   = "2023-12-31T23:59:59Z"
}

resource "defectdojo_stub_finding" "test" {
  title    = %[1]q
  test     = defectdojo_dd_test.test.id
  severity = "Medium"
}
`, title)
}
