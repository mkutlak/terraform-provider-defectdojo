resource "defectdojo_notification_webhook" "example" {
  name = "Slack Security Alerts"
  url  = "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"
}
