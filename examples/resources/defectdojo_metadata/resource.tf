resource "defectdojo_metadata" "example" {
  name  = "environment"
  value = "production"

  # Exactly one of product, location, endpoint, or finding must be set.
  # Attaching metadata to a product is the recommended approach.
  product = data.defectdojo_product.example.id

  # location = defectdojo_url.example.id

  # Discouraged: endpoints are scan-managed projections; prefer product or location.
  # endpoint = 1

  # Discouraged: findings are import-managed; attaching metadata couples state to scan artifacts.
  # finding = 1
}
