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

## [0.4.1](https://github.com/mkutlak/terraform-provider-defectdojo/compare/v0.4.0...v0.4.1) (2026-04-01)


### Miscellaneous

* **deps:** bump github.com/oapi-codegen/runtime from 1.3.0 to 1.3.1 ([01c1823](https://github.com/mkutlak/terraform-provider-defectdojo/commit/01c18232553c9d6144eab295c9bf3b348771681a))

## [0.4.0](https://github.com/mkutlak/terraform-provider-defectdojo/compare/v0.3.1...v0.4.0) (2026-03-24)


### Features

* support name-based lookup for multiple data sources ([8952ac9](https://github.com/mkutlak/terraform-provider-defectdojo/commit/8952ac96ec01513c91070c467385871be2d420f2))

## [0.3.1](https://github.com/mkutlak/terraform-provider-defectdojo/compare/v0.3.0...v0.3.1) (2026-03-23)


### Bug Fixes

* set default empty string for engagement description ([db72468](https://github.com/mkutlak/terraform-provider-defectdojo/commit/db72468bbd884d3d65fe5e2cca68ab9e9acdf945))


### Tests

* verify empty description in engagement acceptance test ([99a9c52](https://github.com/mkutlak/terraform-provider-defectdojo/commit/99a9c527ba3e7c0ead7c887e255cb6ce48e6e9f8))

## [0.3.0](https://github.com/mkutlak/terraform-provider-defectdojo/compare/v0.2.0...v0.3.0) (2026-03-23)


### Features

* improve provider error reporting and schema definitions ([04ff8c4](https://github.com/mkutlak/terraform-provider-defectdojo/commit/04ff8c4459a0190cd9473a13bcd7895ae8b1a8b7))


### Miscellaneous

* add multi-version DefectDojo compatibility testing and automation ([0c2eabb](https://github.com/mkutlak/terraform-provider-defectdojo/commit/0c2eabb8a5a5ae9613cf3d5f3beb5a6f706f5d90))
* update dependencies and add generate-docs target ([97744b2](https://github.com/mkutlak/terraform-provider-defectdojo/commit/97744b2f0e02eb69b3802790062ab3b3fd6c4eba))

## [0.2.0](https://github.com/mkutlak/terraform-provider-defectdojo/compare/v0.1.2...v0.2.0) (2026-03-18)


### Features

* add comprehensive set of resources and data sources ([20bd4d2](https://github.com/mkutlak/terraform-provider-defectdojo/commit/20bd4d28b0401055f5a2c5ff9bbb30e66b0cc6cf))
* implement 0.1.0-rc1 changes and refactor to a new namespace ([3da79fe](https://github.com/mkutlak/terraform-provider-defectdojo/commit/3da79fe3172989ccad8a31f6b87d5ba3dfb3e7df))


### Bug Fixes

* enable parallel execution for acceptance tests and update manifest to v0.1.1 ([191079f](https://github.com/mkutlak/terraform-provider-defectdojo/commit/191079f7b3c750bdde746e5eb96058e39b3fa7c3))
* improve api error handling and resource data population logic ([ffe884f](https://github.com/mkutlak/terraform-provider-defectdojo/commit/ffe884f0d4b67d2d0d4e25da5020815f2a5bc20c))


### Miscellaneous

* **master:** release 0.1.0 ([302f2cd](https://github.com/mkutlak/terraform-provider-defectdojo/commit/302f2cd095053b33ab1f50b7de85fcd47680e38a))
* **master:** release 0.1.0 ([f4b69ea](https://github.com/mkutlak/terraform-provider-defectdojo/commit/f4b69eab5f12b7c57f6b010aefba75c2d9f6336c))
* **master:** release 0.1.2 ([05e51fe](https://github.com/mkutlak/terraform-provider-defectdojo/commit/05e51fe7f38ef8be3b132db3cda8db9626952b79))
* **master:** release 0.1.2 ([503c3d8](https://github.com/mkutlak/terraform-provider-defectdojo/commit/503c3d8e337aeb423352f33e8335a8517b72ad38))
* update CI workflows, add CLAUDE.md, and update .gitignore ([8df7195](https://github.com/mkutlak/terraform-provider-defectdojo/commit/8df719594049f854041d06143c4d3c33004180a0))
* update dependencies in go.mod and go.sum ([6df6327](https://github.com/mkutlak/terraform-provider-defectdojo/commit/6df632715d65bc39d3d483b9e52e365fcabdef67))


### CI/CD

* add setup-terraform to test workflow ([ab474c7](https://github.com/mkutlak/terraform-provider-defectdojo/commit/ab474c7c979f9af0e36e228a0a108d701372734b))
* Add token to release-please workflow ([0656107](https://github.com/mkutlak/terraform-provider-defectdojo/commit/0656107a7e081eb4670bebb63e1888b609dde84c))
* Fix version from 0.1.0-rc1 to 0.0.0 ([4f3b2dc](https://github.com/mkutlak/terraform-provider-defectdojo/commit/4f3b2dc20a5c327b04712c0cf18c90d2aa94c4af))
* setup release-please and improve local testing environment ([416f653](https://github.com/mkutlak/terraform-provider-defectdojo/commit/416f653b22a763742565be279cbc284f0a2dde96))

## [0.1.2](https://github.com/mkutlak/terraform-provider-defectdojo/compare/v0.1.1...v0.1.2) (2026-03-18)


### Bug Fixes

* enable parallel execution for acceptance tests and update manifest to v0.1.1 ([191079f](https://github.com/mkutlak/terraform-provider-defectdojo/commit/191079f7b3c750bdde746e5eb96058e39b3fa7c3))
* improve api error handling and resource data population logic ([ffe884f](https://github.com/mkutlak/terraform-provider-defectdojo/commit/ffe884f0d4b67d2d0d4e25da5020815f2a5bc20c))


### Miscellaneous

* update dependencies in go.mod and go.sum ([6df6327](https://github.com/mkutlak/terraform-provider-defectdojo/commit/6df632715d65bc39d3d483b9e52e365fcabdef67))

## [0.1.0](https://github.com/mkutlak/terraform-provider-defectdojo/compare/v0.0.0...v0.1.0) (2026-03-18)


### Features

* implement 0.1.0-rc1 changes and refactor to a new namespace ([3da79fe](https://github.com/mkutlak/terraform-provider-defectdojo/commit/3da79fe3172989ccad8a31f6b87d5ba3dfb3e7df))


### Miscellaneous

* update CI workflows, add CLAUDE.md, and update .gitignore ([8df7195](https://github.com/mkutlak/terraform-provider-defectdojo/commit/8df719594049f854041d06143c4d3c33004180a0))


### CI/CD

* add setup-terraform to test workflow ([ab474c7](https://github.com/mkutlak/terraform-provider-defectdojo/commit/ab474c7c979f9af0e36e228a0a108d701372734b))
* Add token to release-please workflow ([0656107](https://github.com/mkutlak/terraform-provider-defectdojo/commit/0656107a7e081eb4670bebb63e1888b609dde84c))
* Fix version from 0.1.0-rc1 to 0.0.0 ([4f3b2dc](https://github.com/mkutlak/terraform-provider-defectdojo/commit/4f3b2dc20a5c327b04712c0cf18c90d2aa94c4af))
* setup release-please and improve local testing environment ([416f653](https://github.com/mkutlak/terraform-provider-defectdojo/commit/416f653b22a763742565be279cbc284f0a2dde96))

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
