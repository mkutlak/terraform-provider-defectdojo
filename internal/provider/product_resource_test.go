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

func TestAccProductBaseResource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("dox-test-repo-%s", uniqueId())
	updatedName := fmt.Sprintf("dox-new-name-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProductResourceConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("description"), knownvalue.StringExact("test")),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("product_type_id"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("business_criticality"), knownvalue.StringExact("high")),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("enable_full_risk_acceptance"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("enable_skip_risk_acceptance"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("external_audience"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("internet_accessible"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("life_cycle"), knownvalue.StringExact("production")),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("origin"), knownvalue.StringExact("internal")),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("platform"), knownvalue.StringExact("web")),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("prod_numeric_grade"), knownvalue.Int64Exact(100)),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("revenue"), knownvalue.StringExact("100.00")),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("user_records"), knownvalue.Int64Exact(1000000)),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("tags"), knownvalue.SetExact([]knownvalue.Check{knownvalue.StringExact("bar"), knownvalue.StringExact("foo")})),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("regulation_ids"), knownvalue.SetSizeExact(0)),
				},
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_product.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccProductResourceUpdatedConfig(updatedName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("name"), knownvalue.StringExact(updatedName)),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("description"), knownvalue.StringExact("updated")),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("product_type_id"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("business_criticality"), knownvalue.StringExact("medium")),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("enable_full_risk_acceptance"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("enable_skip_risk_acceptance"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("external_audience"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("internet_accessible"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("life_cycle"), knownvalue.StringExact("retirement")),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("origin"), knownvalue.StringExact("third party library")),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("platform"), knownvalue.StringExact("desktop")),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("prod_numeric_grade"), knownvalue.Int64Exact(50)),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("revenue"), knownvalue.StringExact("200.00")),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("user_records"), knownvalue.Int64Exact(500000)),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("tags"), knownvalue.SetExact([]knownvalue.Check{knownvalue.StringExact("updated")})),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
func TestAccProductResourceNoTags(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("dox-test-repo-%s", uniqueId())
	updatedName := fmt.Sprintf("dox-new-name-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProductResourceNoTagsConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("description"), knownvalue.StringExact("test")),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("product_type_id"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("tags"), knownvalue.Null()),
				},
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_product.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccProductResourceNoTagsConfig(updatedName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("name"), knownvalue.StringExact(updatedName)),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccProductResourceEmptyTags(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("dox-test-repo-%s", uniqueId())
	updatedName := fmt.Sprintf("dox-new-name-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProductResourceEmptyTagsConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("description"), knownvalue.StringExact("test")),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("product_type_id"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("tags"), knownvalue.SetSizeExact(0)),
				},
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_product.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccProductResourceEmptyTagsConfig(updatedName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("name"), knownvalue.StringExact(updatedName)),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccProductResourceDeleteDrift(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("dox-delete-%s", uniqueId())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProductResourceConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("description"), knownvalue.StringExact("test")),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("product_type_id"), knownvalue.Int64Exact(1)),
				},
			},
			// Delete the underlying resource and see that it detects it has been deleted
			{
				ExpectNonEmptyPlan: true,
				Config:             testAccProductResourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDisappears("defectdojo_product.test"),
				),
			},
			{
				Config: testAccProductResourceConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("description"), knownvalue.StringExact("test")),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("product_type_id"), knownvalue.Int64Exact(1)),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccProductResourceInvalid(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("dox-invalid-%s", uniqueId())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				ExpectError: regexp.MustCompile(`.*Invalid\s+Attribute.*`),
				Config:      testAccProductResourceInvalidConfig(name),
			},
		},
	})
}

func testAccProductResourceNoTagsConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test" {
  name = %[1]q
  description = trimspace(<<-DOC
	  test
  DOC
	)
  product_type_id = 1

  business_criticality = "high"
  enable_full_risk_acceptance = false
  enable_skip_risk_acceptance = true
  external_audience = true
  internet_accessible = true
  life_cycle = "production"
  origin = "internal"
  platform = "web"
  prod_numeric_grade = 100
  regulation_ids = []
  revenue = "100.00"
  user_records = 1000000
}
`, name)
}

func testAccProductResourceEmptyTagsConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test" {
  name = %[1]q
  description = trimspace(<<-DOC
	  test
  DOC
	)
  product_type_id = 1

  business_criticality = "high"
  enable_full_risk_acceptance = false
  enable_skip_risk_acceptance = true
  external_audience = true
  internet_accessible = true
  life_cycle = "production"
  origin = "internal"
  platform = "web"
  prod_numeric_grade = 100
  regulation_ids = []
  revenue = "100.00"
  user_records = 1000000
	tags = []
}
`, name)
}

func testAccProductResourceConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test" {
  name = %[1]q
  description = trimspace(<<-DOC
	  test
  DOC
	)
  product_type_id = 1
  tags = ["foo", "bar"]

  business_criticality = "high"
  enable_full_risk_acceptance = true
  enable_skip_risk_acceptance = true
  external_audience = true
  internet_accessible = true
  life_cycle = "production"
  origin = "internal"
  platform = "web"
  prod_numeric_grade = 100
  regulation_ids = []
  revenue = "100.00"
  user_records = 1000000
}
`, name)
}
func testAccProductResourceUpdatedConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test" {
  name = %[1]q
  description = "updated"
  product_type_id = 1
  tags = ["updated"]

  business_criticality = "medium"
  enable_full_risk_acceptance = false
  enable_skip_risk_acceptance = false
  external_audience = false
  internet_accessible = false
  life_cycle = "retirement"
  origin = "third party library"
  platform = "desktop"
  prod_numeric_grade = 50
  regulation_ids = []
  revenue = "200.00"
  user_records = 500000
}
`, name)
}

func testAccProductResourceInvalidConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test" {
  name = %[1]q
  description = "test"
  product_type_id = 1
  tags = ["foo", "BAR"]

  business_criticality = "something else"
  revenue = "a lot of money"
}
`, name)
}
