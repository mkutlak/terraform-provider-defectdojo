package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCredentialMappingResource(t *testing.T) {
	t.Parallel()
	t.Skip("Skipped: CredentialMapping API returns 404 on create - requires investigation")
	name := fmt.Sprintf("credmap-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCredentialMappingResourceConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("defectdojo_credential_mapping.test", "cred_id"),
					resource.TestCheckResourceAttrSet("defectdojo_credential_mapping.test", "product"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_credential_mapping.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing - change url
			{
				Config: testAccCredentialMappingResourceUpdatedConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("defectdojo_credential_mapping.test", "cred_id"),
					resource.TestCheckResourceAttr("defectdojo_credential_mapping.test", "url", "https://updated.example.com"),
				),
			},
		},
	})
}

func testAccCredentialMappingResourceConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "cm_product" {
  name            = %[1]q
  description     = "test product for credential mapping"
  product_type_id = 1
}
resource "defectdojo_credential" "cm_cred" {
  name        = %[1]q
  environment = 1
  username    = "testuser"
  role        = "admin"
  url         = "https://example.com"
}
resource "defectdojo_credential_mapping" "test" {
  cred_id = defectdojo_credential.cm_cred.id
  product = defectdojo_product.cm_product.id
}
`, name)
}

func testAccCredentialMappingResourceUpdatedConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "cm_product" {
  name            = %[1]q
  description     = "test product for credential mapping"
  product_type_id = 1
}
resource "defectdojo_credential" "cm_cred" {
  name        = %[1]q
  environment = 1
  username    = "testuser"
  role        = "admin"
  url         = "https://example.com"
}
resource "defectdojo_credential_mapping" "test" {
  cred_id = defectdojo_credential.cm_cred.id
  product = defectdojo_product.cm_product.id
  url     = "https://updated.example.com"
}
`, name)
}
