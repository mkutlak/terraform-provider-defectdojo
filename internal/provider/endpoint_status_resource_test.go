package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEndpointStatusResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointStatusResourceConfig(false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("defectdojo_endpoint_status.test", "endpoint"),
					resource.TestCheckResourceAttrSet("defectdojo_endpoint_status.test", "finding"),
					resource.TestCheckResourceAttr("defectdojo_endpoint_status.test", "false_positive", "false"),
					resource.TestCheckResourceAttr("defectdojo_endpoint_status.test", "mitigated", "false"),
				),
			},
			{
				ResourceName:      "defectdojo_endpoint_status.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEndpointStatusResourceConfig(true, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_endpoint_status.test", "mitigated", "true"),
				),
			},
		},
	})
}

func testAccEndpointStatusResourceConfig(mitigated bool, falsePositive bool) string {
	return fmt.Sprintf(`
provider "defectdojo" {}

resource "defectdojo_product_type" "test" {
  name = "endpoint-status-test-product-type"
}

resource "defectdojo_product" "test" {
  name         = "endpoint-status-test-product"
  description  = "Test product for endpoint status"
  prod_type    = defectdojo_product_type.test.id
}

resource "defectdojo_endpoint" "test" {
  host    = "endpoint-status-test.example.com"
  product = defectdojo_product.test.id
}

resource "defectdojo_engagement" "test" {
  name        = "endpoint-status-test-engagement"
  product     = defectdojo_product.test.id
  target_start = "2023-01-01"
  target_end   = "2023-12-31"
}

resource "defectdojo_dd_test" "test" {
  engagement  = defectdojo_engagement.test.id
  test_type   = 1
  target_start = "2023-01-01T00:00:00Z"
  target_end   = "2023-12-31T23:59:59Z"
}

resource "defectdojo_stub_finding" "test" {
  title = "endpoint-status-test-finding"
  test  = defectdojo_dd_test.test.id
}

resource "defectdojo_endpoint_status" "test" {
  endpoint       = defectdojo_endpoint.test.id
  finding        = defectdojo_stub_finding.test.id
  mitigated      = %[1]t
  false_positive = %[2]t
}
`, mitigated, falsePositive)
}
