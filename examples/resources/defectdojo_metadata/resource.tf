resource "defectdojo_metadata" "example" {
  name  = "environment"
  value = "production"

  # Exactly one of product or finding must be set.
  # Attaching metadata to a product is the recommended approach.
  product = data.defectdojo_product.example.id



  # Discouraged: findings are import-managed; attaching metadata couples state to scan artifacts.
  # finding = 1
}
