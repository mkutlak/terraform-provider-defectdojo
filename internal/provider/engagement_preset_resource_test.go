package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccEngagementPresetResource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test-%s", uniqueId())
	updatedTitle := fmt.Sprintf("updated-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccEngagementPresetResourceConfig(name, "Test Preset Title"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_engagement_preset.test", tfjsonpath.New("title"), knownvalue.StringExact("Test Preset Title")),
					statecheck.ExpectKnownValue("defectdojo_engagement_preset.test", tfjsonpath.New("product"), knownvalue.NotNull()),
				},
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_engagement_preset.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccEngagementPresetResourceConfig(name, updatedTitle),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_engagement_preset.test", tfjsonpath.New("title"), knownvalue.StringExact(updatedTitle)),
				},
			},
		},
	})
}

// TestAccEngagementPresetResourceMinimal applies a config that omits both title
// and scope. DefectDojo stores an empty string for each and answers with it, so
// before they were marked Computed this failed with "Provider produced
// inconsistent result after apply: .title: was null, but now cty.StringVal("")".
func TestAccEngagementPresetResourceMinimal(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test-minimal-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccEngagementPresetResourceMinimalConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_engagement_preset.test", tfjsonpath.New("title"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("defectdojo_engagement_preset.test", tfjsonpath.New("scope"), knownvalue.StringExact("")),
				},
			},
			{
				Config:   testAccEngagementPresetResourceMinimalConfig(name),
				PlanOnly: true,
			},
		},
	})
}

func testAccEngagementPresetResourceMinimalConfig(productName string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test_product" {
  name            = %[1]q
  description     = "test product for minimal engagement preset"
  product_type_id = 1
}
resource "defectdojo_engagement_preset" "test" {
  product = defectdojo_product.test_product.id
}
`, productName)
}

func testAccEngagementPresetResourceConfig(productName, title string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test_product" {
  name            = %[1]q
  description     = "test product for engagement preset"
  product_type_id = 1
}
resource "defectdojo_engagement_preset" "test" {
  product = defectdojo_product.test_product.id
  title   = %[2]q
  scope   = "All endpoints"
}
`, productName, title)
}
