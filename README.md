# terraform-provider-defectdojo

[DefectDojo API Terraform Provider](https://registry.terraform.io/providers/mkutlak/defectdojo)

Terraform provider for managing [DefectDojo](https://www.defectdojo.org/) resources. See the [provider documentation](https://registry.terraform.io/providers/mkutlak/defectdojo/latest/docs) for supported DefectDojo versions.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.8
- [Go](https://golang.org/doc/install) >= 1.25 (for building from source)

## Using the provider

Configure the provider via environment variables:

```shell
export DEFECTDOJO_BASEURL="https://defectdojo.my-company.com/api/v2"
export DEFECTDOJO_APIKEY="my-api-key"
```

Or with username/password:

```shell
export DEFECTDOJO_BASEURL="https://defectdojo.my-company.com/api/v2"
export DEFECTDOJO_USERNAME="admin"
export DEFECTDOJO_PASSWORD="my-password"
```

Or in the Terraform configuration:

```hcl
provider "defectdojo" {
  base_url = "https://defectdojo.my-company.com/api/v2"
  api_key  = var.dd_api_key
}
```

### Example usage

```hcl
data "defectdojo_product_type" "this" {
  name = "My Product Type"
}

resource "defectdojo_product" "this" {
  name            = "My Application"
  description     = "Managed by Terraform"
  product_type_id = data.defectdojo_product_type.this.id
}

resource "defectdojo_engagement" "security_assessment" {
  product      = defectdojo_product.this.id
  name         = "Security Assessment"
  target_start = "2025-01-01"
  target_end   = "2025-12-31"
}
```

## Supported Resources & Data Sources

### Core

| Resource                  | Data Source               |
| ------------------------- | ------------------------- |
| `defectdojo_product`      | `defectdojo_product`      |
| `defectdojo_product_type` | `defectdojo_product_type` |

### Infrastructure

| Resource                             | Data Source                          |
| ------------------------------------ | ------------------------------------ |
| `defectdojo_development_environment` | `defectdojo_development_environment` |
| `defectdojo_regulation`              | `defectdojo_regulation`              |
| `defectdojo_tool_type`               | `defectdojo_tool_type`               |
| `defectdojo_tool_configuration`      | `defectdojo_tool_configuration`      |
| `defectdojo_sla_configuration`       | `defectdojo_sla_configuration`       |
| `defectdojo_note_type`               | `defectdojo_note_type`               |
| `defectdojo_network_location`        | `defectdojo_network_location`        |
| `defectdojo_language_type`           | `defectdojo_language_type`           |

### Security & Access Control

| Resource                         | Data Source                      |
| -------------------------------- | -------------------------------- |
| `defectdojo_user`                | `defectdojo_user`                |
| `defectdojo_user_contact_info`   | `defectdojo_user_contact_info`   |
| `defectdojo_dojo_group`          | `defectdojo_dojo_group`          |
| `defectdojo_dojo_group_member`   | `defectdojo_dojo_group_member`   |
| `defectdojo_global_role`         | `defectdojo_global_role`         |
| `defectdojo_product_member`      | `defectdojo_product_member`      |
| `defectdojo_product_group`       | `defectdojo_product_group`       |
| `defectdojo_product_type_member` | `defectdojo_product_type_member` |
| `defectdojo_product_type_group`  | `defectdojo_product_type_group`  |
| `defectdojo_credential`          | `defectdojo_credential`          |

### Vulnerability Management

| Resource                       | Data Source                    |
| ------------------------------ | ------------------------------ |
| `defectdojo_engagement`        | `defectdojo_engagement`        |
| `defectdojo_engagement_preset` | `defectdojo_engagement_preset` |
| `defectdojo_test`              | `defectdojo_test`              |
| `defectdojo_finding`           | `defectdojo_finding`           |
| `defectdojo_finding_template`  | `defectdojo_finding_template`  |
| `defectdojo_endpoint`          | `defectdojo_endpoint`          |
| `defectdojo_endpoint_status`   | `defectdojo_endpoint_status`   |
| `defectdojo_stub_finding`      | `defectdojo_stub_finding`      |
| `defectdojo_technology`        | `defectdojo_technology`        |
| `defectdojo_language`          | `defectdojo_language`          |

### Integrations

| Resource                                    | Data Source                                 |
| ------------------------------------------- | ------------------------------------------- |
| `defectdojo_jira_instance`                  | `defectdojo_jira_instance`                  |
| `defectdojo_jira_product_configuration`     | `defectdojo_jira_product_configuration`     |
| `defectdojo_tool_product_settings`          | `defectdojo_tool_product_settings`          |
| `defectdojo_product_api_scan_configuration` | `defectdojo_product_api_scan_configuration` |
| `defectdojo_credential_mapping`             | `defectdojo_credential_mapping`             |
| `defectdojo_risk_acceptance`                | `defectdojo_risk_acceptance`                |
| `defectdojo_notification_webhook`           | `defectdojo_notification_webhook`           |
| `defectdojo_announcement`                   | `defectdojo_announcement`                   |
| `defectdojo_asset_group`                    | `defectdojo_asset_group`                    |

## Developing the Provider

### Building

```shell
go install
```

### Running Tests

Unit tests (no DefectDojo instance required):

```shell
go test ./internal/provider/ -run Test -count=1
```

Acceptance tests (requires a live DefectDojo instance):

```shell
make testacc
```

Run a single test:

```shell
TESTARGS="-run TestFunctionName" make testacc
```

### Local DefectDojo for Testing

A Docker Compose setup is included for running a local DefectDojo instance:

```shell
make dd-up                  # Start DefectDojo (default v2.54.3)
DD_VERSION=2.42.0 make dd-up  # Start a specific version
make dd-spec                # Fetch OpenAPI spec from running instance
make testacc-local          # Run acceptance tests against local instance
make dd-down                # Stop and clean up
make dd-compat              # Run compat checks against supported versions
make dd-compat-test         # Compat checks + acceptance tests
```

Default credentials: `admin` / `testpassword`

### Generating Documentation

```shell
go generate ./...
```

## Releasing

1. Update `CHANGELOG.md` in a PR
2. Tag with `vX.Y.Z` and push tags
3. GoReleaser (via GitHub Actions) builds and publishes to Terraform Registry

Pre-release versions can be tagged from any commit with a suffix (e.g., `v0.1.0-rc1`).

## Acknowledgments

This project is a fork of the original [Doximity DefectDojo Terraform Provider](https://github.com/doximity/terraform-provider-defectdojo). We are grateful for their initial work and contributions to the community.

## License

Licensed under the Apache v2 license. See [LICENSE.md](./LICENSE.md).
