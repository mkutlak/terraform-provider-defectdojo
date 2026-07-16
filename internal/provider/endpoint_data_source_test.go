package provider

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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

	// 1. Seed a product (product type 1 ships with the test instance and is
	// relied on by the other acceptance tests too).
	prodResp, err := client.ProductsCreateWithResponse(ctx, dd.ProductsCreateJSONRequestBody{
		Name:        fmt.Sprintf("tf-acc-endpoint-ds-%s", uniqueId()),
		Description: "endpoint data source acceptance test",
		ProdType:    1,
	})
	if err != nil {
		t.Fatalf("error creating seed product: %v", err)
	}
	if prodResp.JSON201 == nil || prodResp.JSON201.Id == nil {
		t.Fatalf("unexpected response creating seed product: %d\n%s", prodResp.StatusCode(), prodResp.Body)
	}
	productId := *prodResp.JSON201.Id
	defer func() {
		if _, err := client.ProductsDestroy(ctx, productId); err != nil {
			t.Logf("cleanup: error deleting seed product %d: %v", productId, err)
		}
	}()

	// 2. Seed a URL location.
	host := fmt.Sprintf("tf-acc-%s.example.com", uniqueId())
	protocol := "https"
	port := 8443
	path := "tf-acc"
	urlResp, err := client.UrlCreateWithResponse(ctx, dd.UrlCreateJSONRequestBody{
		Host:     host,
		Protocol: &protocol,
		Port:     &port,
		Path:     &path,
	})
	if err != nil {
		t.Fatalf("error creating seed url location: %v", err)
	}
	if urlResp.StatusCode() == 404 {
		// The /api/v2/url/ endpoint only exists in DefectDojo 3.x; on 2.x the
		// endpoint projection under test does not exist either.
		t.Skip("Skipped: this DefectDojo instance does not expose the 3.x location/url APIs")
	}
	if urlResp.JSON201 == nil || urlResp.JSON201.Id == nil {
		t.Fatalf("unexpected response creating seed url location: %d\n%s", urlResp.StatusCode(), urlResp.Body)
	}
	locationId := *urlResp.JSON201.Id
	defer func() {
		if _, err := client.UrlDestroy(ctx, locationId); err != nil {
			t.Logf("cleanup: error deleting seed url location %d: %v", locationId, err)
		}
	}()

	// 3. Link the location to the product; the pair is what the legacy
	// endpoints API projects.
	linkResp, err := client.LocationProductsCreateWithResponse(ctx, dd.LocationProductsCreateJSONRequestBody{
		Location: locationId,
		Product:  productId,
	})
	if err != nil {
		t.Fatalf("error linking seed location to product: %v", err)
	}
	if linkResp.JSON201 == nil || linkResp.JSON201.Id == nil {
		t.Fatalf("unexpected response linking seed location to product: %d\n%s", linkResp.StatusCode(), linkResp.Body)
	}
	linkId := *linkResp.JSON201.Id
	defer func() {
		if _, err := client.LocationProductsDestroy(ctx, linkId); err != nil {
			t.Logf("cleanup: error deleting seed location-product link %d: %v", linkId, err)
		}
	}()

	// 4. Resolve the projected endpoint id via the legacy list API.
	epResp, err := client.EndpointsListWithResponse(ctx, &dd.EndpointsListParams{
		LocationId: &locationId,
		Product:    &productId,
	})
	if err != nil {
		t.Fatalf("error listing endpoints for seeded location: %v", err)
	}
	if epResp.JSON200 == nil || len(epResp.JSON200.Results) != 1 {
		t.Fatalf("expected exactly one projected endpoint for location %d and product %d, got response: %d\n%s",
			locationId, productId, epResp.StatusCode(), epResp.Body)
	}
	endpoint := epResp.JSON200.Results[0]
	if endpoint.Id == nil {
		t.Fatalf("projected endpoint has no id: %+v", endpoint)
	}
	endpointId := *endpoint.Id

	// 5. Read it back through the data source.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointDataSourceConfig(endpointId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.defectdojo_endpoint.test", "id", strconv.Itoa(endpointId)),
					resource.TestCheckResourceAttr("data.defectdojo_endpoint.test", "host", host),
					resource.TestCheckResourceAttr("data.defectdojo_endpoint.test", "protocol", protocol),
					resource.TestCheckResourceAttr("data.defectdojo_endpoint.test", "port", strconv.Itoa(port)),
					resource.TestCheckResourceAttr("data.defectdojo_endpoint.test", "path", path),
					resource.TestCheckResourceAttr("data.defectdojo_endpoint.test", "product", strconv.Itoa(productId)),
				),
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
