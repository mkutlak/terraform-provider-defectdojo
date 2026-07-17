# Lookup by codename
data "defectdojo_configuration_permission" "example" {
  codename = "add_product"
}

# Or lookup by ID
# data "defectdojo_configuration_permission" "by_id" {
#   id = "1"
# }
