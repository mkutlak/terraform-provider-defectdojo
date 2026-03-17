resource "defectdojo_credential" "example" {
  name        = "Web App Admin"
  environment = 1
  username    = "admin"
  role        = "admin"
  url         = "https://app.example.com/login"
}
