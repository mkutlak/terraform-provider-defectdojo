package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccCicdInfrastructureIdDataSource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("dox-test-ci-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCicdInfrastructureDataSourceIdConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.defectdojo_cicd_infrastructure.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.defectdojo_cicd_infrastructure.test", tfjsonpath.New("infrastructure_type"), knownvalue.StringExact("scm_server")),
				},
				Check: testAccCheckDataSourceMatchesResource("data.defectdojo_cicd_infrastructure.test", "defectdojo_cicd_infrastructure.test"),
			},
		},
	})
}

func TestAccCicdInfrastructureNameDataSource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("dox-test-ci-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCicdInfrastructureDataSourceNameConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.defectdojo_cicd_infrastructure.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.defectdojo_cicd_infrastructure.test", tfjsonpath.New("id"), knownvalue.NotNull()),
				},
				Check: testAccCheckDataSourceMatchesResource("data.defectdojo_cicd_infrastructure.test", "defectdojo_cicd_infrastructure.test"),
			},
		},
	})
}

func testAccCicdInfrastructureDataSourceIdConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_cicd_infrastructure" "test" {
  name                = %[1]q
  infrastructure_type = "scm_server"
}
data "defectdojo_cicd_infrastructure" "test" {
  id         = defectdojo_cicd_infrastructure.test.id
  depends_on = [defectdojo_cicd_infrastructure.test]
}
`, name)
}

func testAccCicdInfrastructureDataSourceNameConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_cicd_infrastructure" "test" {
  name                = %[1]q
  infrastructure_type = "scm_server"
}
data "defectdojo_cicd_infrastructure" "test" {
  name       = %[1]q
  depends_on = [defectdojo_cicd_infrastructure.test]
}
`, name)
}
