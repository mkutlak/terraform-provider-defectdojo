package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRegulationDataSource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccRegulationResourceConfig(name) + `
data "defectdojo_regulation" "test" {
  id = defectdojo_regulation.test.id
}
`,
				Check: testAccCheckDataSourceMatchesResource(
					"data.defectdojo_regulation.test", "defectdojo_regulation.test"),
			},
		},
	})
}
