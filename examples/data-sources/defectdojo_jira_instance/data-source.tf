# Lookup by URL
data "defectdojo_jira_instance" "example" {
  url = "https://jira.example.com"
}

# Or lookup by ID
# data "defectdojo_jira_instance" "by_id" {
#   id = "1"
# }
