resource "defectdojo_endpoint" "example" {
  protocol = "https"
  host     = "app.example.com"
  path     = "api/v1"
  product  = 1
}
