package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUrlResource(t *testing.T) {
	t.Parallel()
	host := fmt.Sprintf("host-%s.example.com", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccUrlResourceConfig(host),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_url.test", "host", host),
					resource.TestCheckResourceAttr("defectdojo_url.test", "protocol", "https"),
					resource.TestCheckResourceAttr("defectdojo_url.test", "port", "8443"),
					resource.TestCheckResourceAttr("defectdojo_url.test", "path", "api/v1"),
					resource.TestCheckResourceAttr("defectdojo_url.test", "query", "foo=bar"),
					resource.TestCheckResourceAttr("defectdojo_url.test", "fragment", "section-1"),
					resource.TestCheckResourceAttr("defectdojo_url.test", "user_info", "user:pass"),
					resource.TestCheckResourceAttr("defectdojo_url.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("defectdojo_url.test", "tags.0", "bar"),
					resource.TestCheckResourceAttr("defectdojo_url.test", "tags.1", "foo"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_url.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccUrlResourceUpdatedConfig(host),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_url.test", "host", host),
					resource.TestCheckResourceAttr("defectdojo_url.test", "path", "api/v2"),
					resource.TestCheckResourceAttr("defectdojo_url.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("defectdojo_url.test", "tags.0", "updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccUrlDataSource(t *testing.T) {
	t.Parallel()
	host := fmt.Sprintf("host-%s.example.com", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			// Read by id
			{
				Config: testAccUrlDataSourceConfig(host),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.defectdojo_url.test", "host", host),
					resource.TestCheckResourceAttrPair("data.defectdojo_url.test", "id", "defectdojo_url.test", "id"),
				),
			},
			// Read by host
			{
				Config: testAccUrlDataSourceHostConfig(host),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.defectdojo_url.test", "host", host),
					resource.TestCheckResourceAttrPair("data.defectdojo_url.test", "id", "defectdojo_url.test", "id"),
				),
			},
		},
	})
}

func testAccUrlResourceConfig(host string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_url" "test" {
  host      = %[1]q
  protocol  = "https"
  port      = 8443
  path      = "api/v1"
  query     = "foo=bar"
  fragment  = "section-1"
  user_info = "user:pass"
  tags      = ["foo", "bar"]
}
`, host)
}

func testAccUrlResourceUpdatedConfig(host string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_url" "test" {
  host      = %[1]q
  protocol  = "https"
  port      = 8443
  path      = "api/v2"
  query     = "foo=bar"
  fragment  = "section-1"
  user_info = "user:pass"
  tags      = ["updated"]
}
`, host)
}

func testAccUrlDataSourceConfig(host string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_url" "test" {
  host = %[1]q
}
data "defectdojo_url" "test" {
  id         = defectdojo_url.test.id
  depends_on = [defectdojo_url.test]
}
`, host)
}

func testAccUrlDataSourceHostConfig(host string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_url" "test" {
  host = %[1]q
}
data "defectdojo_url" "test" {
  host       = %[1]q
  depends_on = [defectdojo_url.test]
}
`, host)
}
