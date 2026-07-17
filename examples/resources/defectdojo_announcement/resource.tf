# DefectDojo allows a single global announcement; if one already exists,
# import it instead of creating a new one (see import.sh in this directory).
resource "defectdojo_announcement" "example" {
  message     = "Scheduled maintenance this weekend. Expect brief downtime."
  style       = "warning"
  dismissable = true
}
