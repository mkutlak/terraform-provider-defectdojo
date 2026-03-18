package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDojoGroupMemberResource(t *testing.T) {
	t.Parallel()
	groupName := fmt.Sprintf("member-group-%s", uniqueId())
	username := fmt.Sprintf("member-user-%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDojoGroupMemberResourceConfig(groupName, username),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("defectdojo_dojo_group_member.test", "id"),
					resource.TestCheckResourceAttrSet("defectdojo_dojo_group_member.test", "group"),
					resource.TestCheckResourceAttrSet("defectdojo_dojo_group_member.test", "user"),
					resource.TestCheckResourceAttr("defectdojo_dojo_group_member.test", "role", "3"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_dojo_group_member.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDojoGroupMemberResourceConfig(groupName, username string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_dojo_group" "member_group" {
  name = %[1]q
}
resource "defectdojo_user" "member_user" {
  username = %[2]q
  email    = "member@example.com"
  password = "TestPassword123!"
}
resource "defectdojo_dojo_group_member" "test" {
  group = defectdojo_dojo_group.member_group.id
  user  = defectdojo_user.member_user.id
  role  = 3
}
`, groupName, username)
}
