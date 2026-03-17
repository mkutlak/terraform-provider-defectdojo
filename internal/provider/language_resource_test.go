package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLanguageResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccLanguageResourceConfig(1000, 100, 50),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("defectdojo_language.test", "language_type"),
					resource.TestCheckResourceAttrSet("defectdojo_language.test", "product"),
					resource.TestCheckResourceAttr("defectdojo_language.test", "code", "1000"),
					resource.TestCheckResourceAttr("defectdojo_language.test", "blank", "100"),
					resource.TestCheckResourceAttr("defectdojo_language.test", "comment", "50"),
				),
			},
			{
				ResourceName:      "defectdojo_language.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLanguageResourceConfig(2000, 200, 100),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_language.test", "code", "2000"),
				),
			},
		},
	})
}

func testAccLanguageResourceConfig(code int, blank int, comment int) string {
	return fmt.Sprintf(`
provider "defectdojo" {}

resource "defectdojo_product_type" "test" {
  name = "language-test-product-type"
}

resource "defectdojo_product" "test" {
  name        = "language-test-product"
  description = "Test product for language"
  prod_type   = defectdojo_product_type.test.id
}

resource "defectdojo_language_type" "test" {
  language = "Go"
}

resource "defectdojo_language" "test" {
  language_type = defectdojo_language_type.test.id
  product  = defectdojo_product.test.id
  code     = %[1]d
  blank    = %[2]d
  comment  = %[3]d
}
`, code, blank, comment)
}
