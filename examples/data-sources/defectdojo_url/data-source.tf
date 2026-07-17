# Lookup by host (only works when the host is unique across URL records;
# use id instead when multiple URL records share the same host)
data "defectdojo_url" "example" {
  host = "example.com"
}

# Or lookup by ID
# data "defectdojo_url" "by_id" {
#   id = "1"
# }
