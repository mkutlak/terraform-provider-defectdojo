package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccLocationProductResource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test-%s", uniqueId())
	host := fmt.Sprintf("%s.example.com", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccLocationProductResourceConfig(name, host),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_location_product.test", tfjsonpath.New("location"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("defectdojo_location_product.test", tfjsonpath.New("product"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("defectdojo_location_product.test", tfjsonpath.New("relationship"), knownvalue.StringExact("owned_by")),
				},
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_location_product.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccLocationProductResourceUpdatedConfig(name, host),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_location_product.test", tfjsonpath.New("relationship"), knownvalue.StringExact("used_by")),
				},
			},
		},
	})
}

func testAccLocationProductResourceConfig(name string, host string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test" {
  name            = %[1]q
  description     = "test product for location product"
  product_type_id = 1
}
resource "defectdojo_url" "test" {
  host = %[2]q
}
resource "defectdojo_location_product" "test" {
  location     = defectdojo_url.test.id
  product      = defectdojo_product.test.id
  relationship = "owned_by"
}
`, name, host)
}

func testAccLocationProductResourceUpdatedConfig(name string, host string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test" {
  name            = %[1]q
  description     = "test product for location product"
  product_type_id = 1
}
resource "defectdojo_url" "test" {
  host = %[2]q
}
resource "defectdojo_location_product" "test" {
  location     = defectdojo_url.test.id
  product      = defectdojo_product.test.id
  relationship = "used_by"
}
`, name, host)
}

func TestAccLocationProductDataSource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test-%s", uniqueId())
	host := fmt.Sprintf("%s.example.com", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccLocationProductDataSourceConfig(name, host),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.defectdojo_location_product.test", tfjsonpath.New("location"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.defectdojo_location_product.test", tfjsonpath.New("product"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.defectdojo_location_product.test", tfjsonpath.New("relationship"), knownvalue.StringExact("owned_by")),
				},
			},
		},
	})
}

func testAccLocationProductDataSourceConfig(name string, host string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test" {
  name            = %[1]q
  description     = "test product for location product data source"
  product_type_id = 1
}
resource "defectdojo_url" "test" {
  host = %[2]q
}
resource "defectdojo_location_product" "test" {
  location     = defectdojo_url.test.id
  product      = defectdojo_product.test.id
  relationship = "owned_by"
}
data "defectdojo_location_product" "test" {
  id = defectdojo_location_product.test.id
}
`, name, host)
}
