package provider

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
)

// TestAccEndpointDataSource exercises the read-only defectdojo_endpoint data
// source. Since DefectDojo 3.x endpoints can no longer be created through the
// endpoints API (they are a projection of Locations), the test seeds one
// through the new APIs directly:
//
//  1. POST /api/v2/products/          - a product to attach the location to
//  2. POST /api/v2/url/               - a URL location (subtype of Location)
//  3. POST /api/v2/location/products/ - link the location to the product
//
// GET /api/v2/endpoints/ then projects the (location, product) pair as a
// legacy endpoint, which the data source reads by id.
func TestAccEndpointDataSource(t *testing.T) {
	t.Parallel()
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance test skipped unless TF_ACC is set (test seeds data via the live API)")
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

	seeded := seedUrlLocationWithProduct(t, ctx, client, "tf-acc-endpoint-ds")

	// Resolve the projected endpoint id via the legacy list API.
	epResp, err := client.EndpointsListWithResponse(ctx, &dd.EndpointsListParams{
		LocationId: &seeded.LocationId,
		Product:    &seeded.ProductId,
	})
	if err != nil {
		t.Fatalf("error listing endpoints for seeded location: %v", err)
	}
	if epResp.JSON200 == nil || len(epResp.JSON200.Results) != 1 {
		t.Fatalf("expected exactly one projected endpoint for location %d and product %d, got response: %d\n%s",
			seeded.LocationId, seeded.ProductId, epResp.StatusCode(), epResp.Body)
	}
	endpoint := epResp.JSON200.Results[0]
	if endpoint.Id == nil {
		t.Fatalf("projected endpoint has no id: %+v", endpoint)
	}
	endpointId := *endpoint.Id

	// Read it back through the data source.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointDataSourceConfig(endpointId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.defectdojo_endpoint.test", "port", strconv.Itoa(seeded.Port)),
					resource.TestCheckResourceAttr("data.defectdojo_endpoint.test", "product", strconv.Itoa(seeded.ProductId)),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.defectdojo_endpoint.test", tfjsonpath.New("id"), knownvalue.StringExact(strconv.Itoa(endpointId))),
					statecheck.ExpectKnownValue("data.defectdojo_endpoint.test", tfjsonpath.New("host"), knownvalue.StringExact(seeded.Host)),
					statecheck.ExpectKnownValue("data.defectdojo_endpoint.test", tfjsonpath.New("protocol"), knownvalue.StringExact(seeded.Protocol)),
					statecheck.ExpectKnownValue("data.defectdojo_endpoint.test", tfjsonpath.New("path"), knownvalue.StringExact(seeded.Path)),
				},
			},
		},
	})
}

func testAccEndpointDataSourceConfig(endpointId int) string {
	return fmt.Sprintf(`
provider "defectdojo" {}
data "defectdojo_endpoint" "test" {
  id = "%d"
}
`, endpointId)
}
