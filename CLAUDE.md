# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Terraform provider for [DefectDojo](https://www.defectdojo.org/) - a security vulnerability management tool. The provider uses the Terraform Plugin Framework (not the older SDKv2) and communicates with DefectDojo via an internally generated Go client at `internal/ddclient/` (generated with `oapi-codegen` from DefectDojo's OpenAPI spec).

All prose - schema descriptions, diagnostics, comments, docs - follows [ASD-STE100](https://www.asd-ste100.org/): one idea per sentence, 25 words at most, active voice, one term per concept.

## Build & Test Commands

```shell
go install                              # Build the provider
go generate ./...                       # Regenerate docs (requires terraform CLI)
make testacc                            # Run all acceptance tests (needs live DefectDojo instance)
TESTARGS="-run TestFunctionName" make testacc  # Run a single acceptance test
TF_LOG="DEBUG" make testacc             # Acceptance tests with debug output
go test ./internal/provider/ -run TestProductResource  # Run unit tests (no TF_ACC needed)
DD_VERSION=3.1.101 make dd-up           # Start a specific DefectDojo version (3.x only; 2.x is unsupported)
make dd-spec                            # Fetch OpenAPI spec from running instance
make regen-client                       # Regenerate internal/ddclient from openapi-specs/$DD_VERSION spec
                                        # (spec is a local artifact, not tracked in git; collect it first via `make dd-up && make dd-spec`)
make dd-compat                          # Run multi-version compat checks (spec collection)
make dd-compat-test                     # Run compat checks + acceptance tests
```

### Test Types

- **Unit tests** (`*_unit_test.go`): Test reflection-based data mapping without a DefectDojo instance. Run with plain `go test`.
- **Acceptance tests** (`*_test.go`, excluding unit): Require a live DefectDojo instance. Set `TF_ACC=1` and configure `DEFECTDOJO_BASEURL` + either `DEFECTDOJO_APIKEY` or `DEFECTDOJO_USERNAME`/`DEFECTDOJO_PASSWORD`.

## Architecture

All provider logic lives in `internal/provider/`. There is no separate client wrapper - the provider uses the generated `ddclient` package directly.

### Reflection-Based CRUD Engine

The core abstraction is a **reflection-driven generic CRUD system** that avoids per-resource boilerplate:

- **`resource.go`** defines `terraformResource` which implements Create/Read/Update/Delete/ImportState generically using two interfaces:
  - `terraformResourceData` - wraps a resource's Terraform state struct, provides `id()` and `defectdojoResource()`
  - `defectdojoResource` - wraps a DD API struct, provides `createApiCall`/`readApiCall`/`updateApiCall`/`deleteApiCall`
- **`populateDefectdojoResource()`** uses reflection to copy Terraform state -> DD API struct (guided by `ddField` struct tags)
- **`populateResourceData()`** uses reflection to copy DD API response -> Terraform state (guided by `ddField` struct tags)
- **`datasource.go`** defines `terraformDatasource` with a similar generic Read

### Adding a New Resource

Each resource follows this pattern (see `product_resource.go` as the most complete example):

1. Define a `*ResourceData` struct with `tfsdk` and `ddField` struct tags mapping to the DD client struct fields
2. Define a `*DefectdojoResource` struct embedding the DD client type, implementing the 4 API call methods
3. Define a `*Resource` struct embedding `terraformResource`, implement `Schema()` and `Metadata()`
4. Create a `*DataProvider` implementing `getData()`, and constructor `New*Resource()` wiring the data provider
5. Register in `provider.go` `Resources()` or `DataSources()`

### Key Struct Tag Convention

```go
Name types.String `tfsdk:"name" ddField:"Name"`
```

The `ddField` tag value must match the exact Go field name in the corresponding `ddclient` struct (e.g., `ddclient.Product`). The reflection engine handles type conversions between Terraform types (`types.String`, `types.Bool`, `types.Int64`, `types.Set`) and Go types (string, *string, bool, *bool, int, _int, _[]int, \*[]string).

### Resources & Data Sources

The provider implements 28 resources and 32 data sources. See `provider.go` `Resources()` and `DataSources()` for the full list. Notable special cases: `endpoint` is data-source-only and deprecated (use `url`/`location`); `location`, `user_profile`, `test_type`, and `configuration_permission` are data-source-only; `system_settings` is a singleton resource (adopt-on-create via the engine's `singletonAdopter` interface, destroy = state-remove-only); `announcement` is limited to one instance server-side.

### Notable Files

- Default attribute values use the framework-native `stringdefault`/`booldefault` packages (there is no custom plan-modifier file)
- `scripts/dd-version-compat.sh` - Multi-version compatibility test automation
- `openapi-specs/<version>/defect_dojo.json` - Collected OpenAPI specs per DD version (collected locally via `make dd-up && make dd-spec`, not tracked in git)

### Provider Authentication

Supports two auth modes (resolved in `provider.go`):

1. API key (`DEFECTDOJO_APIKEY` or `api_key` config)
2. Username/password (`DEFECTDOJO_USERNAME`/`DEFECTDOJO_PASSWORD`) - fetches a token via `ApiTokenAuthCreate`

## Intentionally Excluded Resources

The following DefectDojo API objects are **not** implemented as Terraform resources because they are not a good fit for infrastructure-as-code management. Do not re-add them:

| API Object               | Reason                                                                        |
| ------------------------ | ----------------------------------------------------------------------------- |
| Finding                  | Created by scan tool imports, not manual IaC. Extremely complex (50+ fields). |
| StubFinding              | Scan artifact, same as Finding.                                               |
| EndpointStatus           | Join table between endpoints and findings, managed by the system.             |
| Technology (AppAnalysis) | Auto-detected from scan results, not manually managed.                        |
| Language                 | Auto-detected from scan results, not manually managed.                        |
| Credential               | API removed in DefectDojo 3.0.                                                |
| CredentialMapping        | API removed in DefectDojo 3.0.                                                |
| DojoGroup                | API removed in DefectDojo 3.0 RBAC overhaul — replaced by authorized_users field. |
| DojoGroupMember          | API removed in DefectDojo 3.0 RBAC overhaul — replaced by authorized_users field. |
| GlobalRole               | API removed in DefectDojo 3.0 RBAC overhaul — replaced by authorized_users field. |
| ProductMember            | API removed in DefectDojo 3.0 RBAC overhaul — replaced by authorized_users field. |
| ProductGroup             | API removed in DefectDojo 3.0 RBAC overhaul — replaced by authorized_users field. |
| ProductTypeMember        | API removed in DefectDojo 3.0 RBAC overhaul — replaced by authorized_users field. |
| ProductTypeGroup         | API removed in DefectDojo 3.0 RBAC overhaul — replaced by authorized_users field. |
| AssetGroup               | API removed in DefectDojo 3.0 RBAC overhaul — replaced by authorized_users field. |
| Endpoint (as resource)   | Read-only projection since DefectDojo 3.0; data source retained (deprecated — use url/location). |
| LocationFinding          | Join between locations and import-managed findings; system-managed.           |
| SonarQubeIssue / SonarQubeTransition | Scanner-managed artifacts of the SonarQube integration.           |
| TestImport               | Historical records of import runs — artifacts, not desired state.             |
| Notes (incl. object-scoped notes) | Conversational/audit artifacts, not desired state.                   |
| Asset / Organization routes | v3 route aliases of products / product_types — would double-manage state.  |
| jira_configurations / jira_projects routes | Legacy aliases of jira_instances / jira_product_configurations. |
| import-scan / reimport-scan / endpoint_meta_import and other RPC-style endpoints | Actions and artifacts, not resources. |
| TestType (as resource)   | API has no DELETE and update cannot rename; destroy would leave permanent server-side leftovers. Data source available. |
| Metadata location/endpoint parents | Broken upstream in 3.1.101: location parent silently ignored, endpoint parent rejected. Only product/finding parents exposed. |

## Release Process

This project uses [release-please](https://github.com/googleapis/release-please) for automated releases driven by [Conventional Commits](https://www.conventionalcommits.org/).

### Commit Message Format

```
type(scope): description

[optional body]

[optional BREAKING CHANGE: description]
```

**Types:** `feat`, `fix`, `chore`, `docs`, `ci`, `refactor`, `perf`, `test`
**Scopes (optional):** resource name, e.g. `product`, `engagement`, `endpoint`

### Semver Rules

- `feat:` → minor bump (new resource, data source, attribute)
- `fix:` / `chore:` / `docs:` → patch bump
- `feat!:` or `BREAKING CHANGE:` footer → major bump
- Resource/attribute removal or rename = breaking change

### How It Works

1. Merge PRs to `master` using **rebase-merge**
2. release-please automatically opens/updates a Release PR with changelog + version bump
3. Review and merge the Release PR to trigger a release
4. GoReleaser (via GitHub Actions) builds, signs, and publishes to the Terraform Registry

Rebase-merge replays **every** commit onto `master`, so release-please parses all of them
and the PR title does not affect versioning (it is only validated by
`.github/workflows/pr-title-check.yml`). So: every commit must be conventional and
correctly typed, the highest bump in the branch wins, a single `!` or `BREAKING CHANGE:`
footer anywhere cuts a major, and a commit that a later one in the same branch corrects
still reaches the changelog and `git bisect` — squash such pairs before opening the PR.

Commits are GPG-signed and carry `Signed-off-by`, so a rewrite must preserve both and
needs `git push --force-with-lease`. Audit the whole branch, not just the title:

```shell
git log --format="%h %G? %s" master..HEAD      # every subject, and signing after a rewrite
git log master..HEAD | grep "^BREAKING CHANGE" # every body
```

### Configuration

- `release-please-config.json` — release-please settings
- `.release-please-manifest.json` — current version tracker
