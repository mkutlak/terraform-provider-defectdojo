package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// TestAccTestTypeDataSource exercises the read-only defectdojo_test_type data
// source. Test Types cannot be created through this provider (the API has no
// DELETE and its update body cannot change the name), so this test reads a
// built-in Test Type that DefectDojo seeds on startup rather than creating
// one. "Burp Scan" is assumed to be a stable, always-present built-in Test
// Type across supported DefectDojo versions.
func TestAccTestTypeDataSource(t *testing.T) {
	t.Parallel()
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance test skipped unless TF_ACC is set (reads a built-in Test Type from a live instance)")
	}
	testAccPreCheck(t)

	const testTypeName = "Burp Scan"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccTestTypeDataSourceConfig(testTypeName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.defectdojo_test_type.test", tfjsonpath.New("name"), knownvalue.StringExact(testTypeName)),
					statecheck.ExpectKnownValue("data.defectdojo_test_type.test", tfjsonpath.New("id"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func testAccTestTypeDataSourceConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
data "defectdojo_test_type" "test" {
  name = %[1]q
}
`, name)
}
