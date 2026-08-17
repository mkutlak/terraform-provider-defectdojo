package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccEngagementResource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test-%s", uniqueId())
	updatedName := fmt.Sprintf("updated-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccEngagementResourceConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_engagement.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("defectdojo_engagement.test", tfjsonpath.New("target_start"), knownvalue.StringExact("2025-01-01")),
					statecheck.ExpectKnownValue("defectdojo_engagement.test", tfjsonpath.New("target_end"), knownvalue.StringExact("2025-12-31")),
					statecheck.ExpectKnownValue("defectdojo_engagement.test", tfjsonpath.New("product"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("defectdojo_engagement.test", tfjsonpath.New("description"), knownvalue.StringExact("")),
				},
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_engagement.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccEngagementResourceUpdatedConfig(name, updatedName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_engagement.test", tfjsonpath.New("name"), knownvalue.StringExact(updatedName)),
					statecheck.ExpectKnownValue("defectdojo_engagement.test", tfjsonpath.New("target_start"), knownvalue.StringExact("2025-06-01")),
					statecheck.ExpectKnownValue("defectdojo_engagement.test", tfjsonpath.New("target_end"), knownvalue.StringExact("2025-12-31")),
				},
			},
		},
	})
}

func testAccEngagementResourceConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test_product" {
  name            = %[1]q
  description     = "test product for engagement"
  product_type_id = 1
}
resource "defectdojo_engagement" "test" {
  product      = defectdojo_product.test_product.id
  target_start = "2025-01-01"
  target_end   = "2025-12-31"
  name         = %[1]q
}
`, name)
}

func testAccEngagementResourceUpdatedConfig(productName, engagementName string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test_product" {
  name            = %[1]q
  description     = "test product for engagement"
  product_type_id = 1
}
resource "defectdojo_engagement" "test" {
  product      = defectdojo_product.test_product.id
  target_start = "2025-06-01"
  target_end   = "2025-12-31"
  name         = %[2]q
}
`, productName, engagementName)
}

// TestAccEngagementResourceClearOptionalAttributes covers the clearing path on
// the resource with the second-largest Optional-only surface (issue #30).
func TestAccEngagementResourceClearOptionalAttributes(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test-clear-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccEngagementResourceClearConfig(name, `
  branch_tag    = "main"
  build_id      = "build-1"
  commit_hash   = "abc123"
  version       = "v1"
  reason        = "scheduled"
  tracker       = "https://tracker.example.com/1"
  test_strategy = "https://strategy.example.com/plan"
  tags          = ["alpha"]`),
			},
			{
				Config: testAccEngagementResourceClearConfig(name, ""),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_engagement.test", tfjsonpath.New("branch_tag"), knownvalue.Null()),
					statecheck.ExpectKnownValue("defectdojo_engagement.test", tfjsonpath.New("build_id"), knownvalue.Null()),
					statecheck.ExpectKnownValue("defectdojo_engagement.test", tfjsonpath.New("commit_hash"), knownvalue.Null()),
					statecheck.ExpectKnownValue("defectdojo_engagement.test", tfjsonpath.New("version"), knownvalue.Null()),
					statecheck.ExpectKnownValue("defectdojo_engagement.test", tfjsonpath.New("reason"), knownvalue.Null()),
					statecheck.ExpectKnownValue("defectdojo_engagement.test", tfjsonpath.New("tags"), knownvalue.Null()),
				},
			},
			{
				Config:   testAccEngagementResourceClearConfig(name, ""),
				PlanOnly: true,
			},
		},
	})
}

func testAccEngagementResourceClearConfig(name string, extra string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "clear_p" {
  name            = %[1]q
  description     = "clear-attributes regression"
  product_type_id = 1
}
resource "defectdojo_engagement" "test" {
  name         = %[1]q
  product      = defectdojo_product.clear_p.id
  target_start = "2026-01-01"
  target_end   = "2026-12-31"
%[2]s
}
`, name, extra)
}
