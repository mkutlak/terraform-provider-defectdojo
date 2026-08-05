package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccLocationDataSource exercises the read-only defectdojo_location data
// source. Locations have no dedicated create API in DefectDojo 3.x: every
// concrete location (URL, network location, ...) shares the same underlying
// id, so this test seeds a defectdojo_url resource and reads it back through
// defectdojo_location using that same id.
func TestAccLocationDataSource(t *testing.T) {
	t.Parallel()
	host := fmt.Sprintf("tf-acc-location-ds-%s.example.com", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccLocationDataSourceConfig(host),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.defectdojo_location.test", "id", "defectdojo_url.test", "id"),
					resource.TestCheckResourceAttrSet("data.defectdojo_location.test", "location_type"),
					resource.TestCheckResourceAttrSet("data.defectdojo_location.test", "location_value"),
				),
			},
		},
	})
}

func testAccLocationDataSourceConfig(host string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_url" "test" {
  host = %[1]q
}
data "defectdojo_location" "test" {
  id         = defectdojo_url.test.id
  depends_on = [defectdojo_url.test]
}
`, host)
}
