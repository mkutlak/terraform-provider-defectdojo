# Singleton resource: DefectDojo always has exactly one system_settings row.
# `terraform apply` adopts that existing row (there is no create/destroy on
# the API - only List and Update/PartialUpdate) and updates it in place.
# `terraform destroy` only forgets it in Terraform state; the settings keep
# their last-applied values on the server.
#
# Only attributes you set below are managed by Terraform - any attribute you
# omit keeps whatever value is already configured in DefectDojo.
resource "defectdojo_system_settings" "example" {
  team_name        = "Example Security Team"
  disclaimer_notes = "Internal use only."

  enable_deduplication = true
  delete_duplicates    = true
  max_dupes            = 5

  enable_finding_sla         = true
  engagement_auto_close      = true
  engagement_auto_close_days = 3

  # jira_minimum_severity = "High"
}
