package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccProductTypeIdDataSource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("dox-test-repo-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccProductTypeDataSourceIdConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.defectdojo_product_type.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.defectdojo_product_type.test", tfjsonpath.New("description"), knownvalue.StringExact("test")),
				},
			},
		},
	})
}

func TestAccProductTypeNameDataSource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("dox-test-repo-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccProductTypeDataSourceNameConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.defectdojo_product_type.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.defectdojo_product_type.test", tfjsonpath.New("description"), knownvalue.StringExact("test")),
				},
			},
		},
	})
}

func TestAccProductTypeBooleansDataSource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("dox-test-repo-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			// Test default values of our booleans
			{
				Config: testAccProductTypeBooleanChecksDefaultConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.defectdojo_product_type.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.defectdojo_product_type.test", tfjsonpath.New("description"), knownvalue.StringExact("test")),
					statecheck.ExpectKnownValue("data.defectdojo_product_type.test", tfjsonpath.New("critical_product"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("data.defectdojo_product_type.test", tfjsonpath.New("key_product"), knownvalue.Bool(false)),
				},
			},
			// Test our booleans when defined as true
			{
				Config: testAccProductTypeBooleanChecksConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.defectdojo_product_type.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.defectdojo_product_type.test", tfjsonpath.New("description"), knownvalue.StringExact("test")),
					statecheck.ExpectKnownValue("data.defectdojo_product_type.test", tfjsonpath.New("critical_product"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("data.defectdojo_product_type.test", tfjsonpath.New("key_product"), knownvalue.Bool(true)),
				},
			},
		},
	})
}

func testAccProductTypeDataSourceIdConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product_type" "test" {
	name = %[1]q
	description = "test"
}
data "defectdojo_product_type" "test" {
	id = defectdojo_product_type.test.id
	depends_on = [defectdojo_product_type.test]
}
`, name)
}

func testAccProductTypeDataSourceNameConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product_type" "test" {
	name = %[1]q
	description = "test"
}
data "defectdojo_product_type" "test" {
	name = %[2]q
	depends_on = [defectdojo_product_type.test]
}
`, name, name)
}

func testAccProductTypeBooleanChecksDefaultConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product_type" "test" {
	name = %[1]q
	description = "test"
}
data "defectdojo_product_type" "test" {
	name = %[1]q
	depends_on = [defectdojo_product_type.test]
}
`, name)
}

func testAccProductTypeBooleanChecksConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product_type" "test" {
	name = %[1]q
	description = "test"
	critical_product = true
	key_product = true
}
data "defectdojo_product_type" "test" {
	name = %[1]q
	depends_on = [defectdojo_product_type.test]
}
`, name)
}
