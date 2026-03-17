## 0.1.0-rc1

BREAKING CHANGES:
 - Upgraded to Go 1.25 and updated all dependencies.
 - Replaced external `defect-dojo-client-go` with in-repo generated client (`internal/ddclient/`) via oapi-codegen v2, targeting DefectDojo API v2.54.3.
 - Upgraded Terraform Plugin Framework to latest version.

FEATURES:
 - **Infrastructure resources & data sources:**
   - `defectdojo_development_environment`
   - `defectdojo_regulation`
   - `defectdojo_tool_type`
   - `defectdojo_tool_configuration`
   - `defectdojo_sla_configuration`
   - `defectdojo_note_type`
   - `defectdojo_network_location`
   - `defectdojo_language_type`
 - **Security & Access Control resources & data sources:**
   - `defectdojo_user`
   - `defectdojo_user_contact_info`
   - `defectdojo_dojo_group`
   - `defectdojo_dojo_group_member`
   - `defectdojo_global_role`
   - `defectdojo_product_member`
   - `defectdojo_product_group`
   - `defectdojo_product_type_member`
   - `defectdojo_product_type_group`
   - `defectdojo_credential`
 - **Vulnerability Management resources & data sources:**
   - `defectdojo_engagement`
   - `defectdojo_engagement_preset`
   - `defectdojo_test`
   - `defectdojo_finding_template`
   - `defectdojo_endpoint`
 - **Integration resources & data sources:**
   - `defectdojo_jira_instance`
   - `defectdojo_tool_product_settings`
   - `defectdojo_product_api_scan_configuration`
   - `defectdojo_credential_mapping`
   - `defectdojo_risk_acceptance`
   - `defectdojo_notification_webhook`
   - `defectdojo_asset_group`
 - Add `defectdojo_jira_product_configuration` data source (previously resource-only).
 - Add new fields to `defectdojo_product`: `disable_sla_breach_notifications`, `enable_product_tag_inheritance`, `sla_configuration`.
 - Extended reflection engine to support `time.Time`, `Date`, `float64`, `int32`, non-pointer slices, and named string types.

IMPROVEMENTS:
 - Add Docker Compose setup for local DefectDojo acceptance testing.
 - Add `testAccPreCheck()` with environment variable validation.
 - Replace deprecated `resource.UniqueId()` in all test files.
 - Remove commented-out debug code from reflection engine.
 - Update CI/CD workflows: Go 1.25, actions v6, Terraform matrix 1.8/1.9/1.10.

## 0.0.13

FEATURES
  - Add two attributes to `defectdojo_product_type` resource and data source:
    - `critical_product`
    - `key_product`

## 0.0.12

BUGFIX:
 - A product with no tags specified would cause a provider error from terraform.

## 0.0.11

FEATURES:
 - Add the following fields to `defectdojo_product` resource:
   - `business_criticality`
   - `enable_full_risk_acceptance`
   - `enable_skip_risk_acceptance`
   - `external_audience`
   - `internet_accessible`
   - `lifecycle`
   - `origin`
   - `platform`
   - `prod_numeric_grade`
   - `regulation_ids`
   - `revenue`
   - `user_records`

## 0.0.10

FEATURES:
 - Add `jira_product_configuration` resource.

## 0.0.9

BUGFIX:
 - Fix delete-drift detection in `product` and `product_type` resources. If the object was deleted outside terraform we remove it from the state.

## 0.0.8

BUGFIX:
 - Don't continue processing after encountering an error, that cause a panic.

## 0.0.7

Initial public release

## 0.0.4

FEATURES:
 - Add basic support for Product Type resource and data source

## 0.0.3

FEATURES:
 - First working version.
 - Basic support for Product resource and data source.
