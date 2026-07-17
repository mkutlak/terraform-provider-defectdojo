package provider

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
)

// TestAccSystemSettingsResource exercises the defectdojo_system_settings
// singleton resource (adopted via singletonAdopter; see resource.go and
// system_settings_resource.go).
//
// Unlike every other acceptance test in this provider, this test
// deliberately does NOT call t.Parallel(): /api/v2/system_settings/ is a
// single row shared by the entire DefectDojo instance, so running this test
// concurrently with anything else that also mutates system settings would
// race and produce flaky failures.
//
// Because system_settings already exists before Terraform ever touches it
// (there is no create/destroy, only adopt/update), the test:
//  1. Reads the current values of the two fields it manages (team_name,
//     disclaimer_notes) directly via the raw API client.
//  2. Registers a t.Cleanup that restores those exact prior values via
//     PATCH, regardless of test outcome - "terraform destroy" for this
//     resource only forgets it in Terraform state and never touches the
//     server, so nothing else will put the values back.
//  3. Runs the usual create/update/import resource.Test steps against
//     unique-per-run values so repeated runs don't collide.
func TestAccSystemSettingsResource(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance test skipped unless TF_ACC is set (test mutates shared, instance-wide system settings)")
	}
	testAccPreCheck(t)

	ctx := context.Background()
	client, err := newClient(ctx,
		os.Getenv("DEFECTDOJO_BASEURL"),
		os.Getenv("DEFECTDOJO_APIKEY"),
		os.Getenv("DEFECTDOJO_USERNAME"),
		os.Getenv("DEFECTDOJO_PASSWORD"))
	if err != nil {
		t.Fatalf("could not build API client for test seeding: %v", err)
	}

	listResp, err := client.SystemSettingsListWithResponse(ctx, &dd.SystemSettingsListParams{})
	if err != nil {
		t.Fatalf("error listing system settings: %v", err)
	}
	if listResp.JSON200 == nil || len(listResp.JSON200.Results) == 0 || listResp.JSON200.Results[0].Id == nil {
		t.Fatalf("expected exactly one system settings row, got response: %d\n%s", listResp.StatusCode(), listResp.Body)
	}
	current := listResp.JSON200.Results[0]
	settingsId := *current.Id

	// Snapshot the exact prior values of the fields this test manages so
	// they can be restored after the test, whatever its outcome.
	var priorTeamName *string
	if current.TeamName != nil {
		v := *current.TeamName
		priorTeamName = &v
	}
	var priorDisclaimerNotes *string
	if current.DisclaimerNotes != nil {
		v := *current.DisclaimerNotes
		priorDisclaimerNotes = &v
	}

	t.Cleanup(func() {
		restoreReq := dd.PatchedSystemSettingsRequest{
			TeamName:        priorTeamName,
			DisclaimerNotes: priorDisclaimerNotes,
		}
		resp, err := client.SystemSettingsPartialUpdateWithResponse(context.Background(), settingsId, restoreReq)
		if err != nil {
			t.Errorf("error restoring prior system settings values: %v", err)
			return
		}
		if resp.StatusCode() != 200 {
			t.Errorf("unexpected status restoring prior system settings values: %d\n%s", resp.StatusCode(), resp.Body)
		}
	})

	teamName := fmt.Sprintf("tf-acc-team-%s", uniqueId())
	disclaimerNotes := fmt.Sprintf("tf-acc-notes-%s", uniqueId())
	updatedTeamName := fmt.Sprintf("tf-acc-team-updated-%s", uniqueId())
	updatedDisclaimerNotes := fmt.Sprintf("tf-acc-notes-updated-%s", uniqueId())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Adopt (Create) and Read testing
			{
				Config: testAccSystemSettingsResourceConfig(teamName, disclaimerNotes),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_system_settings.test", "team_name", teamName),
					resource.TestCheckResourceAttr("defectdojo_system_settings.test", "disclaimer_notes", disclaimerNotes),
					resource.TestCheckResourceAttr("defectdojo_system_settings.test", "id", fmt.Sprintf("%d", settingsId)),
				),
			},
			// Update and Read testing
			{
				Config: testAccSystemSettingsResourceConfig(updatedTeamName, updatedDisclaimerNotes),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("defectdojo_system_settings.test", "team_name", updatedTeamName),
					resource.TestCheckResourceAttr("defectdojo_system_settings.test", "disclaimer_notes", updatedDisclaimerNotes),
				),
			},
			// ImportState testing. ImportStateVerify is intentionally left
			// false (rather than using ImportStateVerifyIgnore): import
			// reads back all ~70 system_settings attributes from the live
			// API, while this test's config only manages two of them
			// (team_name, disclaimer_notes). Every other Optional+Computed
			// attribute would legitimately differ between "value absent
			// from config, populated by a prior Read" and "value freshly
			// populated by Import" in ways ImportStateVerify can't
			// distinguish from a real bug, so a full state comparison isn't
			// meaningful for this singleton resource.
			{
				ResourceName:      "defectdojo_system_settings.test",
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Delete testing automatically occurs in TestCase; for this
			// singleton resource it only removes defectdojo_system_settings
			// from Terraform state (see the t.Cleanup above, which restores
			// the real server-side values).
		},
	})
}

func testAccSystemSettingsResourceConfig(teamName, disclaimerNotes string) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
resource "defectdojo_system_settings" "test" {
  team_name        = %q
  disclaimer_notes = %q
}
`, teamName, disclaimerNotes)
}
