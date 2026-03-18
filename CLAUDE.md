# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Terraform provider for [DefectDojo](https://www.defectdojo.org/) - a security vulnerability management tool. The provider uses the Terraform Plugin Framework (not the older SDKv2) and communicates with DefectDojo via an internally generated Go client at `internal/ddclient/` (generated with `oapi-codegen` from DefectDojo's OpenAPI spec).

## Build & Test Commands

```shell
go install                              # Build the provider
go generate ./...                       # Regenerate docs (requires terraform CLI)
make testacc                            # Run all acceptance tests (needs live DefectDojo instance)
TESTARGS="-run TestFunctionName" make testacc  # Run a single acceptance test
TF_LOG="DEBUG" make testacc             # Acceptance tests with debug output
go test ./internal/provider/ -run TestProductResource  # Run unit tests (no TF_ACC needed)
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

The `ddField` tag value must match the exact Go field name in the corresponding `ddclient` struct (e.g., `ddclient.Product`). The reflection engine handles type conversions between Terraform types (`types.String`, `types.Bool`, `types.Int64`, `types.Set`) and Go types (string, *string, bool, *bool, int, *int, *[]int, *[]string).

### Resources & Data Sources

The provider implements 33 resources and 33 data sources. See `provider.go` `Resources()` and `DataSources()` for the full list.

### Notable Files

- `plan_modifiers.go` - Custom plan modifiers for default values (`stringDefault`, `boolDefault`)
- `internal/ref/main.go` - Generic `Of[E]()` helper to create pointers (used in data sources)

### Provider Authentication

Supports two auth modes (resolved in `provider.go`):
1. API key (`DEFECTDOJO_APIKEY` or `api_key` config)
2. Username/password (`DEFECTDOJO_USERNAME`/`DEFECTDOJO_PASSWORD`) - fetches a token via `ApiTokenAuthCreate`

## Intentionally Excluded Resources

The following DefectDojo API objects are **not** implemented as Terraform resources because they are not a good fit for infrastructure-as-code management. Do not re-add them:

| API Object | Reason |
|------------|--------|
| Finding | Created by scan tool imports, not manual IaC. Extremely complex (50+ fields). |
| StubFinding | Scan artifact, same as Finding. |
| EndpointStatus | Join table between endpoints and findings, managed by the system. |
| Technology (AppAnalysis) | Auto-detected from scan results, not manually managed. |
| Language | Auto-detected from scan results, not manually managed. |
| Announcement | DefectDojo only allows ONE global announcement. Cannot create/delete freely. |

## Release Process

1. Update `CHANGELOG.md` in a PR, merge to `master`
2. Tag `master` with `vX.Y.Z` and push tags
3. GoReleaser (via GitHub Actions) builds and publishes to Terraform Registry
