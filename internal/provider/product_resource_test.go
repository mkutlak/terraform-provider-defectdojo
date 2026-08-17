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

// TestAccProductResourceRevenueLiteral applies a revenue that is valid per the
// schema regex but not in DefectDojo's canonical two-decimal-place form.
//
// DefectDojo stores revenue as a Django DecimalField(decimal_places=2), so it
// echoes "100" back as "100.00". Before the ddFormat:"decimal" tag, that
// rewrote state out from under the config and every apply failed with
// "Provider produced inconsistent result after apply". Every other product test
// hardcodes the already-canonical "100.00", which is why this went unnoticed.
func TestAccProductResourceRevenueLiteral(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test-revenue-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccProductResourceRevenueConfig(name, "100"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("revenue"), knownvalue.StringExact("100")),
				},
			},
			{
				Config:   testAccProductResourceRevenueConfig(name, "100"),
				PlanOnly: true,
			},
			// A single decimal place is canonicalised to two server-side too.
			{
				Config: testAccProductResourceRevenueConfig(name, "250.5"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("revenue"), knownvalue.StringExact("250.5")),
				},
			},
			// The canonical spelling must still round-trip untouched.
			{
				Config: testAccProductResourceRevenueConfig(name, "300.00"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("revenue"), knownvalue.StringExact("300.00")),
				},
			},
		},
	})
}

func testAccProductResourceRevenueConfig(name string, revenue string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test" {
  name            = %[1]q
  description     = "revenue round-trip test"
  product_type_id = 1
  revenue         = %[2]q
}
`, name, revenue)
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
  tags = ["foo", "needs review"]

  business_criticality = "something else"
  revenue = "a lot of money"
}
`, name)
}

// TestAccProductResourceClearOptionalAttributes covers the clearing path on a
// resource with both scalar and collection attributes (issue #30).
//
// Collections clear to an empty array rather than an explicit null, because
// the API types tags and regulations as non-nullable arrays.
func TestAccProductResourceClearOptionalAttributes(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test-clear-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccProductResourceClearConfig(name, `
  revenue              = "100.00"
  business_criticality = "high"
  user_records         = 1000
  prod_numeric_grade   = 90
  tags                 = ["alpha", "beta"]`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("revenue"), knownvalue.StringExact("100.00")),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("user_records"), knownvalue.Int64Exact(1000)),
				},
			},
			{
				Config: testAccProductResourceClearConfig(name, ""),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("revenue"), knownvalue.Null()),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("business_criticality"), knownvalue.Null()),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("user_records"), knownvalue.Null()),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("prod_numeric_grade"), knownvalue.Null()),
					statecheck.ExpectKnownValue("defectdojo_product.test", tfjsonpath.New("tags"), knownvalue.Null()),
				},
			},
			{
				Config:   testAccProductResourceClearConfig(name, ""),
				PlanOnly: true,
			},
		},
	})
}

func testAccProductResourceClearConfig(name string, extra string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test" {
  name            = %[1]q
  description     = "clear-attributes regression"
  product_type_id = 1
%[2]s
}
`, name, extra)
}

// TestAccProductResourceTagCaseCollision is the regression test for
// DefectDojo's instance-wide, case-insensitive tag table.
//
// The server answers with whichever spelling of a tag was registered FIRST
// anywhere on the instance, so a perfectly lower-case configuration breaks as
// soon as the capitalised spelling exists - from the UI, another tool, or
// another Terraform resource. Before preserveTagCase, the second apply below
// failed with:
//
//	.tags: planned set element cty.StringVal("zz-...") does not correlate with
//	any element in actual
//
// The tag name is unique per run so this cannot disturb any other test.
func TestAccProductResourceTagCaseCollision(t *testing.T) {
	t.Parallel()
	suffix := uniqueId()
	upper := fmt.Sprintf("Zz-%s", suffix)
	lower := fmt.Sprintf("zz-%s", suffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			// Register the capitalised spelling in the instance-wide tag table.
			{
				Config: testAccProductResourceTagCaseConfig(suffix, upper, upper),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_product.seed", tfjsonpath.New("tags"),
						knownvalue.SetExact([]knownvalue.Check{knownvalue.StringExact(upper)})),
				},
			},
			// A second product asks for the lower-case spelling. DefectDojo
			// stores the same tag and answers with the capitalised form; state
			// must keep what the practitioner wrote.
			{
				Config: testAccProductResourceTagCaseConfig(suffix, upper, lower),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_product.other", tfjsonpath.New("tags"),
						knownvalue.SetExact([]knownvalue.Check{knownvalue.StringExact(lower)})),
				},
			},
			// And the result has to be stable, not re-diff forever.
			{
				Config:   testAccProductResourceTagCaseConfig(suffix, upper, lower),
				PlanOnly: true,
			},
		},
	})
}

func testAccProductResourceTagCaseConfig(suffix, seedTag, otherTag string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "seed" {
  name            = "tagcase-seed-%[1]s"
  description     = "registers the capitalised tag spelling"
  product_type_id = 1
  tags            = [%[2]q]
}
resource "defectdojo_product" "other" {
  name            = "tagcase-other-%[1]s"
  description     = "asks for a different case of the same tag"
  product_type_id = 1
  tags            = [%[3]q]
  depends_on      = [defectdojo_product.seed]
}
`, suffix, seedTag, otherTag)
}
