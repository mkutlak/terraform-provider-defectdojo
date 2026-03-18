package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEndpointResource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("endpoint-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointResourceConfig(name, "example.com", "https"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_endpoint.test", "host", "example.com"),
					resource.TestCheckResourceAttr("defectdojo_endpoint.test", "protocol", "https"),
				),
			},
			{
				ResourceName:      "defectdojo_endpoint.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEndpointResourceConfig(name, "updated.example.com", "https"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_endpoint.test", "host", "updated.example.com"),
				),
			},
		},
	})
}

func testAccEndpointResourceConfig(name string, host string, protocol string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "ep_product" {
  name            = %[1]q
  description     = "test product for endpoint"
  product_type_id = 1
}
resource "defectdojo_endpoint" "test" {
  host     = %[2]q
  protocol = %[3]q
  product  = defectdojo_product.ep_product.id
}
`, name, host, protocol)
}
