package provider

import (
	"fmt"
	"regexp"
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

// TestAccEngagementResourceProductTagInheritance covers a child of a product
// with enable_product_tag_inheritance set.
//
// DefectDojo merges the product's tags into every Engagement, Test and Finding
// underneath it, and the create response then lists the child's own tag twice.
// A Terraform set cannot hold that, so the apply used to die on a framework
// diagnostic that named neither DefectDojo nor tag inheritance:
//
//	Error: Duplicate Set Element
//	  with defectdojo_engagement.child,
//	  This attribute contains duplicate values of: tftypes.String<"tfsprint">
//
// The server's list is deduplicated now, so what surfaces is the real problem
// underneath: the server holds a tag the configuration never asked for.
// Inherited tags are also sticky - PATCHing them away re-adds them - so the
// only configuration DefectDojo can satisfy is one that lists them, which is
// the second step here.
func TestAccEngagementResourceProductTagInheritance(t *testing.T) {
	t.Parallel()
	suffix := uniqueId()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			// Not listing the inherited tag cannot converge. This pins the
			// error the practitioner actually gets; it must never go back to
			// being "Duplicate Set Element".
			{
				Config:      testAccEngagementInheritanceConfig(suffix, fmt.Sprintf(`"sprint-%s"`, suffix)),
				ExpectError: regexp.MustCompile(`(?s)inconsistent result after apply`),
			},
			// Listing it does converge, and stays converged.
			{
				Config: testAccEngagementInheritanceConfig(suffix,
					fmt.Sprintf(`"sprint-%[1]s", "team-%[1]s"`, suffix)),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_engagement.child", tfjsonpath.New("tags"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact(fmt.Sprintf("sprint-%s", suffix)),
							knownvalue.StringExact(fmt.Sprintf("team-%s", suffix)),
						})),
				},
			},
			{
				Config: testAccEngagementInheritanceConfig(suffix,
					fmt.Sprintf(`"sprint-%[1]s", "team-%[1]s"`, suffix)),
				PlanOnly: true,
			},
		},
	})
}

func testAccEngagementInheritanceConfig(suffix, childTags string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "inherit_p" {
  name                           = "inherit-%[1]s"
  description                    = "propagates its tags to every child"
  product_type_id                = 1
  enable_product_tag_inheritance = true
  tags                           = ["team-%[1]s"]
}
resource "defectdojo_engagement" "child" {
  name         = "inherit-child-%[1]s"
  product      = defectdojo_product.inherit_p.id
  target_start = "2026-01-01"
  target_end   = "2026-12-31"
  tags         = [%[2]s]
}
`, suffix, childTags)
}
