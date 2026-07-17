resource "defectdojo_location_product" "example" {
  location     = defectdojo_url.example.id
  product      = defectdojo_product.example.id
  relationship = "owned_by"
  status       = "Active"
}
