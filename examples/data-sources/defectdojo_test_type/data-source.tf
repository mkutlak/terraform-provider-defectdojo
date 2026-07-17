# Lookup by name
data "defectdojo_test_type" "example" {
  name = "Burp Scan"
}

# Or lookup by ID
# data "defectdojo_test_type" "by_id" {
#   id = "1"
# }
