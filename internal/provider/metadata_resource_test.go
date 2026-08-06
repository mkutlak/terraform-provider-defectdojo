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

func TestAccMetadataResource(t *testing.T) {
	t.Parallel()
	productName := fmt.Sprintf("dox-test-repo-%s", uniqueId())
	metadataName := fmt.Sprintf("dox-meta-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccMetadataResourceConfig(productName, metadataName, "bar"),
				// product is Int64 while id is String, so this pair is compared as
				// flatmap strings; statecheck.CompareValuePairs is type-strict and
				// reports "271 != 271" for values that are equal but differently typed.
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("defectdojo_metadata.test", "product", "defectdojo_product.test", "id"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_metadata.test", tfjsonpath.New("name"), knownvalue.StringExact(metadataName)),
					statecheck.ExpectKnownValue("defectdojo_metadata.test", tfjsonpath.New("value"), knownvalue.StringExact("bar")),
				},
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_metadata.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccMetadataResourceConfig(productName, metadataName, "baz"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_metadata.test", tfjsonpath.New("name"), knownvalue.StringExact(metadataName)),
					statecheck.ExpectKnownValue("defectdojo_metadata.test", tfjsonpath.New("value"), knownvalue.StringExact("baz")),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccMetadataResourceValidatorConflict(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("dox-meta-conflict-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config:      testAccMetadataResourceConflictConfig(name),
				ExpectError: regexp.MustCompile(`(?s)Invalid Attribute Combination`),
			},
		},
	})
}

func TestAccMetadataDataSource(t *testing.T) {
	t.Parallel()
	productName := fmt.Sprintf("dox-test-repo-%s", uniqueId())
	metadataName := fmt.Sprintf("dox-meta-ds-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			// Read by id
			{
				Config: testAccMetadataDataSourceConfig(productName, metadataName),
				// product is Int64 while id is String, so this pair is compared as
				// flatmap strings; statecheck.CompareValuePairs is type-strict and
				// reports "271 != 271" for values that are equal but differently typed.
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.defectdojo_metadata.by_id", "product", "defectdojo_product.test", "id"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.defectdojo_metadata.by_id", tfjsonpath.New("name"), knownvalue.StringExact(metadataName)),
					statecheck.ExpectKnownValue("data.defectdojo_metadata.by_id", tfjsonpath.New("value"), knownvalue.StringExact("bar")),
				},
			},
			// Read by name
			{
				Config: testAccMetadataDataSourceNameConfig(productName, metadataName),
				// product is Int64 while id is String, so this pair is compared as
				// flatmap strings; statecheck.CompareValuePairs is type-strict and
				// reports "271 != 271" for values that are equal but differently typed.
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.defectdojo_metadata.by_name", "product", "defectdojo_product.test", "id"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.defectdojo_metadata.by_name", tfjsonpath.New("name"), knownvalue.StringExact(metadataName)),
					statecheck.ExpectKnownValue("data.defectdojo_metadata.by_name", tfjsonpath.New("value"), knownvalue.StringExact("bar")),
				},
			},
		},
	})
}

func testAccMetadataResourceConfig(productName, metadataName, value string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test" {
  name            = %[1]q
  description     = "test"
  product_type_id = 1
}
resource "defectdojo_metadata" "test" {
  name    = %[2]q
  value   = %[3]q
  product = defectdojo_product.test.id
}
`, productName, metadataName, value)
}

func testAccMetadataResourceConflictConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test" {
  name            = "%[1]s-product"
  description     = "test"
  product_type_id = 1
}
resource "defectdojo_metadata" "test" {
  name     = %[1]q
  value    = "bar"
  product  = defectdojo_product.test.id
  finding  = 1
}
`, name)
}

func testAccMetadataDataSourceConfig(productName, metadataName string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test" {
  name            = %[1]q
  description     = "test"
  product_type_id = 1
}
resource "defectdojo_metadata" "test" {
  name    = %[2]q
  value   = "bar"
  product = defectdojo_product.test.id
}
data "defectdojo_metadata" "by_id" {
  id         = defectdojo_metadata.test.id
  depends_on = [defectdojo_metadata.test]
}
`, productName, metadataName)
}

func testAccMetadataDataSourceNameConfig(productName, metadataName string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test" {
  name            = %[1]q
  description     = "test"
  product_type_id = 1
}
resource "defectdojo_metadata" "test" {
  name    = %[2]q
  value   = "bar"
  product = defectdojo_product.test.id
}
data "defectdojo_metadata" "by_name" {
  name       = %[2]q
  depends_on = [defectdojo_metadata.test]
}
`, productName, metadataName)
}
