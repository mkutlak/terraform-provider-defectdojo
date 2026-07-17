# Always describes the authenticated user (the caller identified by the
# provider's credentials). Takes no lookup parameters.
data "defectdojo_user_profile" "me" {}

# Example: authorize the calling user on a product. The id attribute is a
# string, so it must be converted to a number for the authorized_users set.
# resource "defectdojo_product" "example" {
#   name             = "Example Product"
#   description      = "An example product"
#   prod_type        = 1
#   authorized_users = [tonumber(data.defectdojo_user_profile.me.id)]
# }
