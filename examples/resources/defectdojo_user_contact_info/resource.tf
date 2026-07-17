resource "defectdojo_user_contact_info" "example" {
  user                         = 1
  phone_number                 = "+15551234567"
  deduplication_execution_mode = "async"
}
