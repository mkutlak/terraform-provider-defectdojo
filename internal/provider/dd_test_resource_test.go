package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDDTestResource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDDTestResourceConfig(name, "Test Title"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_test.test", "title", "Test Title"),
					resource.TestCheckResourceAttr("defectdojo_test.test", "test_type", "1"),
					resource.TestCheckResourceAttr("defectdojo_test.test", "target_start", "2025-01-01T10:00:00Z"),
					resource.TestCheckResourceAttr("defectdojo_test.test", "target_end", "2025-01-01T18:00:00Z"),
					resource.TestCheckResourceAttrSet("defectdojo_test.test", "engagement"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_test.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccDDTestResourceConfig(name, "Updated Test Title"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_test.test", "title", "Updated Test Title"),
				),
			},
		},
	})
}

func testAccDDTestResourceConfig(name, title string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test_product" {
  name            = %[1]q
  description     = "test product for test resource"
  product_type_id = 1
}
resource "defectdojo_engagement" "test_engagement" {
  product      = defectdojo_product.test_product.id
  target_start = "2025-01-01"
  target_end   = "2025-12-31"
  name         = %[1]q
}
resource "defectdojo_test" "test" {
  test_type    = 1
  engagement   = defectdojo_engagement.test_engagement.id
  target_start = "2025-01-01T10:00:00Z"
  target_end   = "2025-01-01T18:00:00Z"
  title        = %[2]q
}
`, name, title)
}
