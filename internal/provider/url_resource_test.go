package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_url.test", tfjsonpath.New("host"), knownvalue.StringExact(host)),
					statecheck.ExpectKnownValue("defectdojo_url.test", tfjsonpath.New("protocol"), knownvalue.StringExact("https")),
					statecheck.ExpectKnownValue("defectdojo_url.test", tfjsonpath.New("port"), knownvalue.Int64Exact(8443)),
					statecheck.ExpectKnownValue("defectdojo_url.test", tfjsonpath.New("path"), knownvalue.StringExact("api/v1")),
					statecheck.ExpectKnownValue("defectdojo_url.test", tfjsonpath.New("query"), knownvalue.StringExact("foo=bar")),
					statecheck.ExpectKnownValue("defectdojo_url.test", tfjsonpath.New("fragment"), knownvalue.StringExact("section-1")),
					statecheck.ExpectKnownValue("defectdojo_url.test", tfjsonpath.New("user_info"), knownvalue.StringExact("user:pass")),
					statecheck.ExpectKnownValue("defectdojo_url.test", tfjsonpath.New("tags"), knownvalue.SetExact([]knownvalue.Check{knownvalue.StringExact("bar"), knownvalue.StringExact("foo")})),
				},
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_url.test", tfjsonpath.New("host"), knownvalue.StringExact(host)),
					statecheck.ExpectKnownValue("defectdojo_url.test", tfjsonpath.New("path"), knownvalue.StringExact("api/v2")),
					statecheck.ExpectKnownValue("defectdojo_url.test", tfjsonpath.New("tags"), knownvalue.SetExact([]knownvalue.Check{knownvalue.StringExact("updated")})),
				},
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.defectdojo_url.test", tfjsonpath.New("host"), knownvalue.StringExact(host)),
					statecheck.CompareValuePairs("data.defectdojo_url.test", tfjsonpath.New("id"), "defectdojo_url.test", tfjsonpath.New("id"), compare.ValuesSame()),
				},
			},
			// Read by host
			{
				Config: testAccUrlDataSourceHostConfig(host),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.defectdojo_url.test", tfjsonpath.New("host"), knownvalue.StringExact(host)),
					statecheck.CompareValuePairs("data.defectdojo_url.test", tfjsonpath.New("id"), "defectdojo_url.test", tfjsonpath.New("id"), compare.ValuesSame()),
				},
			},
		},
	})
}

// TestAccUrlResourceInvalidTags asserts the tag grammar is enforced during
// plan. Before defectdojo_url shared product's validator, an uppercase tag was
// accepted here and only failed at apply - after the URL had been created -
// because DefectDojo answers with the lower-cased form.
//
// This mirrors TestAccProductResourceInvalid, which has covered the same
// grammar on defectdojo_product since commit a54fb89.
func TestAccUrlResourceInvalidTags(t *testing.T) {
	t.Parallel()
	host := fmt.Sprintf("host-%s.example.com", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				ExpectError: regexp.MustCompile(`.*Invalid\s+Attribute.*`),
				Config:      testAccUrlResourceInvalidTagsConfig(host),
			},
		},
	})
}

func testAccUrlResourceInvalidTagsConfig(host string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_url" "test" {
  host = %[1]q
  tags = ["ok", "BAR"]
}
`, host)
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
