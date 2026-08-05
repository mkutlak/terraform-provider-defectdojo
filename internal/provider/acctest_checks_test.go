package provider

import (
	"context"
	"fmt"
	"sort"
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

		// Data sources appear in the same map as managed resources, keyed with a
		// "data." prefix - the same discriminator Terraform itself uses
		// (terraform/resource_address.go:78-80). They are reads, not ownership:
		// nothing was created, so nothing should be expected to disappear. Two
		// distinct failures come from not skipping them: data-source-only types
		// (defectdojo_endpoint, _location, _user_profile,
		// _configuration_permission) have no resource model to derive, and types
		// shared with a resource (defectdojo_user) would assert that a
		// pre-existing object the test merely looked up had been deleted.
		if strings.HasPrefix(name, "data.") {
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

// testAccCheckDataSourceMatchesResource asserts that every attribute the data
// source exposes has the same value as the corresponding managed resource.
//
// Both sides are populated by the same reflection engine from the same model,
// so any divergence is a mapping bug. Comparing the whole overlap rather than a
// hand-picked attribute or two means a data source cannot quietly drop a field
// as the model grows - which is exactly how these tests rot.
//
// Attributes the data source does not expose at all are skipped: a data source
// schema is legitimately allowed to be a subset. Pass names in ignore for
// attributes that are expected to differ (e.g. write-only credentials that the
// API never echoes back).
func testAccCheckDataSourceMatchesResource(dataSourceName, resourceName string, ignore ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("data source not found in state: %s", dataSourceName)
		}
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", resourceName)
		}

		skip := make(map[string]bool, len(ignore))
		for _, k := range ignore {
			skip[k] = true
		}

		var mismatches []string
		compared := 0

		for key, resourceValue := range rs.Primary.Attributes {
			if skip[key] {
				continue
			}
			dataValue, present := ds.Primary.Attributes[key]
			if !present {
				continue
			}
			compared++
			if dataValue != resourceValue {
				mismatches = append(mismatches,
					fmt.Sprintf("%s: resource=%q data source=%q", key, resourceValue, dataValue))
			}
		}

		// A schema mismatch or a renamed attribute would otherwise make this
		// check pass by comparing nothing at all.
		if compared == 0 {
			return fmt.Errorf("%s vs %s: no overlapping attributes were compared; "+
				"this check would pass vacuously", dataSourceName, resourceName)
		}
		if len(mismatches) > 0 {
			sort.Strings(mismatches)
			return fmt.Errorf("%s does not match %s in %d of %d compared attributes:\n  %s",
				dataSourceName, resourceName, len(mismatches), compared, strings.Join(mismatches, "\n  "))
		}
		return nil
	}
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

		if strings.HasPrefix(resourceName, "data.") {
			return fmt.Errorf("%s: is a data source; the disappears pattern applies to managed resources only", resourceName)
		}

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
