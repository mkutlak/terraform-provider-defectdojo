# Lookup by username
data "defectdojo_user" "example" {
  username = "admin"
}

# Or lookup by ID
# data "defectdojo_user" "by_id" {
#   id = "1"
# }
