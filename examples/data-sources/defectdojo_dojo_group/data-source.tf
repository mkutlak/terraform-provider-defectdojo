# Lookup by name
data "defectdojo_dojo_group" "example" {
  name = "My Group"
}

# Or lookup by ID
# data "defectdojo_dojo_group" "by_id" {
#   id = "1"
# }
