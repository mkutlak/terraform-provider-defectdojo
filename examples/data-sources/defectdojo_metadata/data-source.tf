# Lookup by id
data "defectdojo_metadata" "example" {
  id = "1"
}

# Or lookup by name (names are only unique per parent object; if more than
# one parent object has metadata with the same name, look up by id instead)
# data "defectdojo_metadata" "by_name" {
#   name = "environment"
# }
