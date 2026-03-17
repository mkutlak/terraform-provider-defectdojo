resource "defectdojo_tool_product_settings" "example" {
  name               = "Snyk for My App"
  setting_url        = "https://app.snyk.io/org/my-org/project/abc123"
  product            = 1
  tool_configuration = 1
}
