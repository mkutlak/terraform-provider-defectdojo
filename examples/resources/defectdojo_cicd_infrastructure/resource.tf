resource "defectdojo_cicd_infrastructure" "example" {
  name                = "Jenkins"
  infrastructure_type = "build_server"
  description         = "Primary build server"
  url                 = "https://jenkins.example.com"
}
