package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMetadataResource(t *testing.T) {
	t.Parallel()
	productName := fmt.Sprintf("dox-test-repo-%s", uniqueId())
	metadataName := fmt.Sprintf("dox-meta-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccMetadataResourceConfig(productName, metadataName, "bar"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_metadata.test", "name", metadataName),
					resource.TestCheckResourceAttr("defectdojo_metadata.test", "value", "bar"),
					resource.TestCheckResourceAttrPair("defectdojo_metadata.test", "product", "defectdojo_product.test", "id"),
				),
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
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_metadata.test", "name", metadataName),
					resource.TestCheckResourceAttr("defectdojo_metadata.test", "value", "baz"),
				),
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
		Steps: []resource.TestStep{
			// Read by id
			{
				Config: testAccMetadataDataSourceConfig(productName, metadataName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.defectdojo_metadata.by_id", "name", metadataName),
					resource.TestCheckResourceAttr("data.defectdojo_metadata.by_id", "value", "bar"),
					resource.TestCheckResourceAttrPair("data.defectdojo_metadata.by_id", "product", "defectdojo_product.test", "id"),
				),
			},
			// Read by name
			{
				Config: testAccMetadataDataSourceNameConfig(productName, metadataName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.defectdojo_metadata.by_name", "name", metadataName),
					resource.TestCheckResourceAttr("data.defectdojo_metadata.by_name", "value", "bar"),
					resource.TestCheckResourceAttrPair("data.defectdojo_metadata.by_name", "product", "defectdojo_product.test", "id"),
				),
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
