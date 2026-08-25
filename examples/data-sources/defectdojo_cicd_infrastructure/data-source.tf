# Lookup by name
data "defectdojo_cicd_infrastructure" "example" {
  name = "Jenkins"
}

# Or lookup by ID
# data "defectdojo_cicd_infrastructure" "by_id" {
#   id = "1"
# }
