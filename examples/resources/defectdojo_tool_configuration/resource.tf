resource "defectdojo_tool_configuration" "example" {
  name                = "Snyk Configuration"
  tool_type           = 1
  authentication_type = "API"
}
