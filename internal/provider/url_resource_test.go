package provider

import (
	"fmt"
	"regexp"
	"strings"
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

// TestAccUrlResourceInvalid asserts the plan-time validators. Both the tag
// grammar and the protocol set are things DefectDojo answers with a 400, so
// catching them at plan time turns an apply-time failure into a config error.
//
// "HTTPS" is the interesting protocol spelling: it is the upper-case form of an
// accepted value, which makes it the natural mistake to pair with an upper-case
// host.
//
// These steps never reach the server, so they share one test case rather than
// paying for a lifecycle each.
func TestAccUrlResourceInvalid(t *testing.T) {
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
			{
				ExpectError: regexp.MustCompile(`.*Invalid\s+Attribute\s+Value\s+Match.*`),
				Config:      testAccUrlResourceProtocolConfig(host, "HTTPS"),
			},
			{
				ExpectError: regexp.MustCompile(`.*Invalid\s+Attribute\s+Value\s+Match.*`),
				Config:      testAccUrlResourceProtocolConfig(host, "wss"),
			},
		},
	})
}

// TestAccUrlResourceHostCase applies a host whose configured spelling is not
// the one DefectDojo stores.
//
// DefectDojo case-folds every host (dojo/url/models.py, URL.clean_host), so the
// read path used to write "api.example.com" over a configured "API.Example.COM"
// and every apply failed at create with "Provider produced inconsistent result
// after apply", after which the URL sat in state tainted and the next apply
// destroyed it, created it again and failed identically. Every other url test
// hardcodes an already-lower-case host, which is why this went unnoticed.
//
// There is deliberately no ImportState step: import has no configured value to
// preserve, so it stores the server's lower-case spelling, exactly like the
// decimal and datetime helpers. TestAccUrlResource still covers import with a
// lower-case host.
func TestAccUrlResourceHostCase(t *testing.T) {
	t.Parallel()
	host := fmt.Sprintf("Host-%s.Example.COM", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccUrlResourceHostOnlyConfig(host),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_url.test", tfjsonpath.New("host"), knownvalue.StringExact(host)),
				},
			},
			// The refresh must not rewrite state either, so the next plan is empty.
			{
				Config:   testAccUrlResourceHostOnlyConfig(host),
				PlanOnly: true,
			},
			// Rewriting the configured spelling is still a change Terraform
			// applies and converges on, rather than a permanent diff.
			{
				Config: testAccUrlResourceHostOnlyConfig(strings.ToLower(host)),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_url.test", tfjsonpath.New("host"), knownvalue.StringExact(strings.ToLower(host))),
				},
			},
		},
	})
}

func testAccUrlResourceProtocolConfig(host, protocol string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_url" "test" {
  host     = %[1]q
  protocol = %[2]q
}
`, host, protocol)
}

func testAccUrlResourceHostOnlyConfig(host string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_url" "test" {
  host = %[1]q
}
`, host)
}

func testAccUrlResourceInvalidTagsConfig(host string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_url" "test" {
  host = %[1]q
  tags = ["ok", "needs review"]
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
