package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// TestAccEngagementPresetResource opens on a config that omits scope.
//
// scope backs a NOT NULL column with no database default, so a request that
// leaves it out returns 500 (IntegrityError) - the provider must put an empty
// string on the wire, which is what the schema Default does. Sending "" is
// accepted, unlike title, where the serializer also rejects a blank value.
// title is therefore set in every step: it is Required precisely because there
// is no value the provider could supply on the practitioner's behalf.
//
// The minimal config has to come first, or it becomes an update and never
// exercises the create path that returned the 500.
func TestAccEngagementPresetResource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test-%s", uniqueId())
	updatedTitle := fmt.Sprintf("updated-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			// Create with scope omitted
			{
				Config: testAccEngagementPresetResourceMinimalConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_engagement_preset.test", tfjsonpath.New("scope"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("defectdojo_engagement_preset.test", tfjsonpath.New("product"), knownvalue.NotNull()),
				},
			},
			{
				Config:   testAccEngagementPresetResourceMinimalConfig(name),
				PlanOnly: true,
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_engagement_preset.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccEngagementPresetResourceConfig(name, "Test Preset Title"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_engagement_preset.test", tfjsonpath.New("title"), knownvalue.StringExact("Test Preset Title")),
				},
			},
			{
				Config: testAccEngagementPresetResourceConfig(name, updatedTitle),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_engagement_preset.test", tfjsonpath.New("title"), knownvalue.StringExact(updatedTitle)),
				},
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
  title   = "minimal preset"
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
