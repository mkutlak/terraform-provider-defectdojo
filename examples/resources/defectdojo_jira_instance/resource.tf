resource "defectdojo_jira_instance" "example" {
  url                       = "https://jira.example.com"
  username                  = "jira-user@example.com"
  epic_name_id              = 10014
  open_status_key           = 11
  close_status_key          = 31
  info_mapping_severity     = "Lowest"
  low_mapping_severity      = "Low"
  medium_mapping_severity   = "Medium"
  high_mapping_severity     = "High"
  critical_mapping_severity = "Highest"
}
