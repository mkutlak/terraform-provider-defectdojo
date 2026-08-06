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

func TestAccJiraProductConfigurationResource(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("dox-test-repo-%s", uniqueId())
	jirakey := fmt.Sprintf("APPSEC%s", uniqueId())
	newjirakey := fmt.Sprintf("APPSEC%s", uniqueId())
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccJiraProductConfigurationResourceConfig(name, jirakey),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_jira_product_configuration.test", tfjsonpath.New("project_key"), knownvalue.StringExact(jirakey)),
					statecheck.ExpectKnownValue("defectdojo_jira_product_configuration.test", tfjsonpath.New("issue_template_dir"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("defectdojo_jira_product_configuration.test", tfjsonpath.New("push_all_issues"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("defectdojo_jira_product_configuration.test", tfjsonpath.New("enable_engagement_epic_mapping"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("defectdojo_jira_product_configuration.test", tfjsonpath.New("push_notes"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("defectdojo_jira_product_configuration.test", tfjsonpath.New("product_jira_sla_notification"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("defectdojo_jira_product_configuration.test", tfjsonpath.New("risk_acceptance_expiration_notification"), knownvalue.Bool(false)),
				},
			},
			// ImportState testing
			{
				ResourceName:      "defectdojo_jira_product_configuration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccJiraProductConfigurationResourceUpdateConfig(name, newjirakey),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_jira_product_configuration.test", tfjsonpath.New("project_key"), knownvalue.StringExact(newjirakey)),
					statecheck.ExpectKnownValue("defectdojo_jira_product_configuration.test", tfjsonpath.New("issue_template_dir"), knownvalue.StringExact("some/dir")),
					statecheck.ExpectKnownValue("defectdojo_jira_product_configuration.test", tfjsonpath.New("push_all_issues"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("defectdojo_jira_product_configuration.test", tfjsonpath.New("enable_engagement_epic_mapping"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("defectdojo_jira_product_configuration.test", tfjsonpath.New("push_notes"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("defectdojo_jira_product_configuration.test", tfjsonpath.New("product_jira_sla_notification"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("defectdojo_jira_product_configuration.test", tfjsonpath.New("risk_acceptance_expiration_notification"), knownvalue.Bool(true)),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccJiraProductConfigurationResourceDeleteDrift(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("dox-delete-%s", uniqueId())
	jirakey := fmt.Sprintf("APPSEC%s", uniqueId())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccJiraProductConfigurationResourceConfig(name, jirakey),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_jira_product_configuration.test", tfjsonpath.New("project_key"), knownvalue.StringExact(jirakey)),
				},
			},
			// Delete the underlying resource and see that it detects it has been deleted
			{
				ExpectNonEmptyPlan: true,
				Config:             testAccJiraProductConfigurationResourceConfig(name, jirakey),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDisappears("defectdojo_jira_product_configuration.test"),
				),
			},
			{
				Config: testAccJiraProductConfigurationResourceConfig(name, jirakey),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("defectdojo_jira_product_configuration.test", tfjsonpath.New("project_key"), knownvalue.StringExact(jirakey)),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccJiraProductConfigurationResourceInvalid(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("dox-delete-%s", uniqueId())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				ExpectError: regexp.MustCompile(`Invalid\s+Resource`),
				Config:      testAccInvalidJiraProductConfigurationResourceConfig(name),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccInvalidJiraProductConfigurationResourceConfig(name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_jira_product_configuration" "test" {
  project_key = %[1]q
}
`, name)
}

func testAccJiraProductConfigurationResourceConfig(productname string, name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test" {
  name = %[1]q
	description = "test"
  product_type_id = 1
}
resource "defectdojo_jira_product_configuration" "test" {
  product_id = defectdojo_product.test.id
  project_key = %[2]q
}
`, productname, name)
}

func testAccJiraProductConfigurationResourceUpdateConfig(productname string, name string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_product" "test" {
  name = %[1]q
	description = "test"
  product_type_id = 1
}
resource "defectdojo_jira_product_configuration" "test" {
  product_id = defectdojo_product.test.id
  project_key = %[2]q
  issue_template_dir = "some/dir"
  push_all_issues = true
  enable_engagement_epic_mapping = true
  push_notes = true
  product_jira_sla_notification = true
  risk_acceptance_expiration_notification = true
}
`, productname, name)
}
