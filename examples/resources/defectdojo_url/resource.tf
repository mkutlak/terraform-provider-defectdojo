resource "defectdojo_url" "example" {
  host      = "example.com"
  protocol  = "https"
  port      = 443
  path      = "/api/v1"
  query     = "foo=bar"
  fragment  = "section-1"
  user_info = "user:pass"
  tags      = ["prod"]
}
