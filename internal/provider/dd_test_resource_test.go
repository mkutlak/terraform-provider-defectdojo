package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccDDTestResource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDDTestResourceConfig(name, "Test Title"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_test.test", tfjsonpath.New("title"), knownvalue.StringExact("Test Title")),
					statecheck.ExpectKnownValue("defectdojo_test.test", tfjsonpath.New("test_type"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("defectdojo_test.test", tfjsonpath.New("target_start"), knownvalue.StringExact("2025-01-01T10:00:00Z")),
					statecheck.ExpectKnownValue("defectdojo_test.test", tfjsonpath.New("target_end"), knownvalue.StringExact("2025-01-01T18:00:00Z")),
					statecheck.ExpectKnownValue("defectdojo_test.test", tfjsonpath.New("engagement"), knownvalue.NotNull()),
				},
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_test.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccDDTestResourceConfig(name, "Updated Test Title"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_test.test", tfjsonpath.New("title"), knownvalue.StringExact("Updated Test Title")),
				},
			},
		},
	})
}

func testAccDDTestResourceConfig(name, title string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test_product" {
  name            = %[1]q
  description     = "test product for test resource"
  product_type_id = 1
}
resource "defectdojo_engagement" "test_engagement" {
  product      = defectdojo_product.test_product.id
  target_start = "2025-01-01"
  target_end   = "2025-12-31"
  name         = %[1]q
}
resource "defectdojo_test" "test" {
  test_type    = 1
  engagement   = defectdojo_engagement.test_engagement.id
  target_start = "2025-01-01T10:00:00Z"
  target_end   = "2025-01-01T18:00:00Z"
  title        = %[2]q
}
`, name, title)
}

// TestAccDDTestResourceDateOnlyTargets covers GitHub issue #23:
// `defectdojo_test.target_start`/`target_end` are backed by time.Time and used
// to be parsed as RFC3339 only, so a date-only literal such as one produced by
// formatdate("YYYY-MM-DD", ...) (which works fine on defectdojo_engagement)
// caused "Provider produced inconsistent result after apply" because the
// value round-tripped through the DD API and came back null. The provider now
// accepts date-only input on datetime-backed attributes and preserves the
// practitioner's literal in state (see preserveDateTimeLiteral in datetime.go).
//
// This test deliberately has NO ImportState step, unlike TestAccDDTestResource.
// Import cannot know the practitioner's literal, so it records the canonical
// RFC3339 form returned by the server; running ImportStateVerify against a
// date-only config would therefore report a false-positive mismatch. That is
// expected, documented behaviour, not a bug.
func TestAccDDTestResourceDateOnlyTargets(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test-dateonly-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			// Create with date-only target_start/target_end. Before the fix,
			// this step failed with "inconsistent result after apply" raised
			// by Terraform core itself, so this step alone reproduces the issue.
			{
				Config: testAccDDTestResourceDateOnlyConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_test.dateonly", "target_start", "2026-07-28"),
					resource.TestCheckResourceAttr("defectdojo_test.dateonly", "target_end", "2031-07-28"),
					resource.TestCheckResourceAttrSet("defectdojo_test.dateonly", "id"),
				),
			},
			// Re-apply the same config as plan-only: proves preserveDateTimeLiteral
			// keeps the literal on Read, so the refresh path does not drift into a
			// perpetual diff.
			{
				Config:   testAccDDTestResourceDateOnlyConfig(name),
				PlanOnly: true,
			},
		},
	})
}

func testAccDDTestResourceDateOnlyConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "dateonly_product" {
  name            = %[1]q
  description     = "test product for date-only test resource"
  product_type_id = 1
}
resource "defectdojo_engagement" "dateonly_engagement" {
  product      = defectdojo_product.dateonly_product.id
  target_start = "2025-01-01"
  target_end   = "2025-12-31"
  name         = %[1]q
}
resource "defectdojo_test" "dateonly" {
  test_type    = 1
  engagement   = defectdojo_engagement.dateonly_engagement.id
  target_start = "2026-07-28"
  target_end   = "2031-07-28"
  title        = %[1]q
}
`, name)
}

// TestAccDdTestResourceClearOptionalAttributes is the regression test for
// GitHub issue #30.
//
// Removing a previously-set Optional attribute from configuration used to fail
// with "Provider produced inconsistent result after apply: .branch_tag: was
// null, but now cty.StringVal(...)". Every generated request struct field is a
// pointer with `omitempty`, so the update request simply omitted the field and
// DefectDojo left the column untouched; the read path then wrote the stale
// value back over the planned null.
//
// The error was retry-proof and fired after the server had already been
// mutated. Clearing now goes through an explicit-null PATCH (see clear.go).
func TestAccDdTestResourceClearOptionalAttributes(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test-clear-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccDdTestResourceClearConfig(name, `
  branch_tag       = "main"
  build_id         = "build-1"
  commit_hash      = "abc123"
  percent_complete = 50
  version          = "v1"
  tags             = ["alpha"]`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_test.t", tfjsonpath.New("branch_tag"), knownvalue.StringExact("main")),
					statecheck.ExpectKnownValue("defectdojo_test.t", tfjsonpath.New("percent_complete"), knownvalue.Int64Exact(50)),
				},
			},
			// Every optional attribute removed from config: they must actually
			// be cleared server-side, not silently retained.
			{
				Config: testAccDdTestResourceClearConfig(name, ""),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_test.t", tfjsonpath.New("branch_tag"), knownvalue.Null()),
					statecheck.ExpectKnownValue("defectdojo_test.t", tfjsonpath.New("build_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("defectdojo_test.t", tfjsonpath.New("commit_hash"), knownvalue.Null()),
					statecheck.ExpectKnownValue("defectdojo_test.t", tfjsonpath.New("percent_complete"), knownvalue.Null()),
					statecheck.ExpectKnownValue("defectdojo_test.t", tfjsonpath.New("version"), knownvalue.Null()),
					statecheck.ExpectKnownValue("defectdojo_test.t", tfjsonpath.New("tags"), knownvalue.Null()),
				},
			},
			// The cleared state must be stable, not re-diff on every plan.
			{
				Config:   testAccDdTestResourceClearConfig(name, ""),
				PlanOnly: true,
			},
		},
	})
}

func testAccDdTestResourceClearConfig(name string, extra string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "p" {
  name            = %[1]q
  description     = "issue 30 regression"
  product_type_id = 1
}
resource "defectdojo_engagement" "e" {
  name         = %[1]q
  product      = defectdojo_product.p.id
  target_start = "2026-01-01"
  target_end   = "2026-12-31"
}
resource "defectdojo_test" "t" {
  engagement   = defectdojo_engagement.e.id
  test_type    = 1
  target_start = "2026-01-01T00:00:00Z"
  target_end   = "2026-01-02T00:00:00Z"
%[2]s
}
`, name, extra)
}
