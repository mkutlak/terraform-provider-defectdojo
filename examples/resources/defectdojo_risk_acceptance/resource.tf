resource "defectdojo_risk_acceptance" "example" {
  name              = "Accepted Low-Risk Findings Q1"
  owner             = 1
  accepted_findings = [1, 2]
}
