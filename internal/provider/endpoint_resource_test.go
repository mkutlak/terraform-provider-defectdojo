package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEndpointResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointResourceConfig("example.com", "https"),
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
				Config: testAccEndpointResourceConfig("updated.example.com", "https"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_endpoint.test", "host", "updated.example.com"),
				),
			},
		},
	})
}

func testAccEndpointResourceConfig(host string, protocol string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_endpoint" "test" {
  host     = %[1]q
  protocol = %[2]q
}
`, host, protocol)
}
