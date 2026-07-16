resource "defectdojo_product" "example" {
  name            = "An example name"
  description     = "An example description"
  product_type_id = data.defectdojo_product_type.example.id

  # IDs of users authorized on this product (DefectDojo 3.x)
  # authorized_users = [1, 2]
}
