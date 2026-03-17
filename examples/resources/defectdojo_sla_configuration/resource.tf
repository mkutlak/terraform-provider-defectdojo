resource "defectdojo_sla_configuration" "example" {
  name     = "Default SLA"
  critical = 7
  high     = 30
  medium   = 90
  low      = 120
}
