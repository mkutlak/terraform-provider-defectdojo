resource "defectdojo_notifications" "example" {
  # Scope this row to a single product. `product` and `user` are mutually
  # exclusive; leave both unset only when importing the pre-existing global
  # default row (see import.sh).
  product = data.defectdojo_product.example.id

  scan_added                 = ["alert", "mail"]
  sla_breach                 = ["alert"]
  sla_breach_combined        = ["alert"]
  risk_acceptance_expiration = ["mail"]
  engagement_added           = ["alert"]
  close_engagement           = ["alert"]
  auto_close_engagement      = ["alert"]
  stale_engagement           = ["mail"]
  upcoming_engagement        = ["mail"]
  test_added                 = ["alert"]
  review_requested           = ["alert", "mail"]
  code_review                = ["alert"]
  user_mentioned             = ["alert"]
  product_added              = ["mail"]
  product_type_added         = ["mail"]
  jira_update                = ["alert"]
  other                      = []

  # Triggered whenever an (re-)import has been done, even if it created,
  # updated, or closed no findings. Single-value, not a set.
  scan_added_empty = "alert"
}

# NOTE: a fresh DefectDojo instance pre-creates a global default notification
# row (product and user both unset) and a per-user default row for every
# existing user (e.g. admin). Attempting to create a new row for a scope
# (product/user combination) that already exists fails with:
#   "Notification for user and product already exists"
# Manage those pre-existing rows by importing them instead - see import.sh.
