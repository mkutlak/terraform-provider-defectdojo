package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// testAccCheckDestroyed verifies that every defectdojo resource left in state
// is really gone server-side after `terraform destroy`.
//
// It deliberately goes through the SAME code path the resource's Read uses -
// defectdojoResource().readApiCall - so the destroy check can never drift from
// Read. This mirrors the discipline the AWS provider gets by having CheckDestroy
// call the resource's own finder, except that the reflection engine lets one
// function cover every resource instead of one per resource.
//
// Resources implementing singletonAdopter are skipped: the engine turns their
// Delete into a state-removal (resource.go:339-344), so the API object is
// expected to survive.
func testAccCheckDestroyed(s *terraform.State) error {
	ctx := context.Background()

	client, err := testAccClient(ctx)
	if err != nil {
		return fmt.Errorf("building client for destroy check: %w", err)
	}

	var survivors []string

	for name, rs := range s.RootModule().Resources {
		if !strings.HasPrefix(rs.Type, "defectdojo_") {
			continue
		}

		data, ok := testAccResourceModel(rs.Type)
		if !ok {
			// Loud rather than silent: an unknown type means the derivation in
			// acctest_registry_test.go missed a resource, which would otherwise
			// quietly disable destroy verification for it.
			return fmt.Errorf("%s (%s): no model could be derived; "+
				"destroy cannot be verified (see TestAccResourceRegistryIsDerivable)", name, rs.Type)
		}

		if testAccIsUndestroyable(data) {
			continue
		}

		if rs.Primary == nil || rs.Primary.ID == "" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("%s: parsing id %q: %w", name, rs.Primary.ID, err)
		}

		statusCode, body, err := data.defectdojoResource().readApiCall(ctx, client, id)
		if err != nil {
			return fmt.Errorf("%s: reading %s/%d during destroy check: %w", name, rs.Type, id, err)
		}

		switch statusCode {
		case 404:
			// Correctly destroyed.
		case 200:
			survivors = append(survivors, fmt.Sprintf("%s (%s/%d)", name, rs.Type, id))
		default:
			return fmt.Errorf("%s: unexpected status %d reading %s/%d during destroy check\n\nbody:\n\n%s",
				name, statusCode, rs.Type, id, string(body))
		}
	}

	if len(survivors) > 0 {
		return fmt.Errorf("resources still exist after destroy: %s", strings.Join(survivors, ", "))
	}
	return nil
}

// testAccCheckDisappears deletes a resource out-of-band, so a follow-up step can
// assert the provider notices the drift and plans a recreate. Pair it with
// ExpectNonEmptyPlan on the same step.
//
// This replaces the previous hardcoded two-type regexp switch: it is driven by
// the derived registry, so it works for every resource rather than just
// defectdojo_product and defectdojo_jira_product_configuration.
func testAccCheckDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := context.Background()

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found in state: %s", resourceName)
		}
		if rs.Primary == nil || rs.Primary.ID == "" {
			return fmt.Errorf("%s: id is not set", resourceName)
		}

		data, ok := testAccResourceModel(rs.Type)
		if !ok {
			return fmt.Errorf("%s (%s): no model could be derived", resourceName, rs.Type)
		}
		if testAccIsUndestroyable(data) {
			return fmt.Errorf("%s (%s): implements singletonAdopter and cannot be deleted server-side; "+
				"the disappears pattern does not apply", resourceName, rs.Type)
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("%s: parsing id %q: %w", resourceName, rs.Primary.ID, err)
		}

		client, err := testAccClient(ctx)
		if err != nil {
			return fmt.Errorf("building client for disappears check: %w", err)
		}

		statusCode, body, err := data.defectdojoResource().deleteApiCall(ctx, client, id)
		if err != nil {
			return fmt.Errorf("%s: deleting %s/%d out of band: %w", resourceName, rs.Type, id, err)
		}
		if statusCode != 204 {
			return fmt.Errorf("%s: unexpected status %d deleting %s/%d out of band\n\nbody:\n\n%s",
				resourceName, statusCode, rs.Type, id, string(body))
		}
		return nil
	}
}
