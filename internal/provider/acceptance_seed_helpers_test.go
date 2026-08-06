package provider

import (
	"context"
	"fmt"
	"testing"

	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
)

// seededLocation captures the ids and attributes of a URL location seeded via
// seedUrlLocationWithProduct, along with the product and location-product
// link it was seeded under.
type seededLocation struct {
	ProductId  int
	LocationId int
	LinkId     int
	Host       string
	Protocol   string
	Path       string
	Port       int
}

// seedUrlLocationWithProduct seeds a product, a URL location, and a link
// between them through the live DefectDojo API:
//
//  1. POST /api/v2/products/          - a product to attach the location to
//  2. POST /api/v2/url/               - a URL location (subtype of Location)
//  3. POST /api/v2/location/products/ - link the location to the product
//
// prefix is combined with uniqueId() to derive unique product name and host
// values. Cleanup of all three seeded objects is registered via t.Cleanup, in
// reverse creation order (link, then location, then product).
func seedUrlLocationWithProduct(t *testing.T, ctx context.Context, client *dd.ClientWithResponses, prefix string) seededLocation {
	t.Helper()

	// 1. Seed a product (product type 1 ships with the test instance and is
	// relied on by the other acceptance tests too).
	prodResp, err := client.ProductsCreateWithResponse(ctx, dd.ProductsCreateJSONRequestBody{
		Name:        fmt.Sprintf("%s-%s", prefix, uniqueId()),
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
	t.Cleanup(func() {
		resp, err := client.ProductsDestroy(ctx, productId)
		if err != nil {
			t.Logf("cleanup: error deleting seed product %d: %v", productId, err)
			return
		}
		_ = resp.Body.Close()
	})

	// 2. Seed a URL location.
	host := fmt.Sprintf("%s-%s.example.com", prefix, uniqueId())
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
	if urlResp.JSON201 == nil || urlResp.JSON201.Id == nil {
		t.Fatalf("unexpected response creating seed url location: %d\n%s", urlResp.StatusCode(), urlResp.Body)
	}
	locationId := *urlResp.JSON201.Id
	t.Cleanup(func() {
		resp, err := client.UrlDestroy(ctx, locationId)
		if err != nil {
			t.Logf("cleanup: error deleting seed url location %d: %v", locationId, err)
			return
		}
		_ = resp.Body.Close()
	})

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
	t.Cleanup(func() {
		resp, err := client.LocationProductsDestroy(ctx, linkId)
		if err != nil {
			t.Logf("cleanup: error deleting seed location-product link %d: %v", linkId, err)
			return
		}
		_ = resp.Body.Close()
	})

	return seededLocation{
		ProductId:  productId,
		LocationId: locationId,
		LinkId:     linkId,
		Host:       host,
		Protocol:   protocol,
		Path:       path,
		Port:       port,
	}
}
