package ddclient

// To regenerate the client from an OpenAPI spec, run:
//
//	make regen-client DD_VERSION=<version>    # e.g. DD_VERSION=3.2.300
//
// Prerequisite: openapi-specs/<version>/defect_dojo.json must exist locally.
// Collected specs are LOCAL artifacts, not tracked in git - collect one first
// with `make dd-up && make dd-spec` (optionally DD_VERSION=<version>).
//
// The target encodes the full validated procedure:
//
//  1. Copy openapi-specs/<version>/defect_dojo.json to the repo root as
//     defect_dojo.json (the root copy is gitignored; do not commit it).
//  2. In that copy, remove the bare null from each enum that oapi-codegen
//     writes without quotes: non-string schemas and date/date-time strings.
//     oapi-codegen v2.8.1 and later rejects such a null. String enum members
//     keep it as the quoted value "<nil>", which is valid Go.
//  3. Run the oapi-codegen version pinned in go.mod (never a global binary):
//     go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@<pinned> \
//     -generate types,client,skip-fmt -package ddclient -o client.gen.go defect_dojo.json
//     The raw output does not compile until step 4 removes the constants
//     of time types. (For the same reason the yaml config in this directory
//     cannot generate as-is.)
//  4. Remove the "Defines values for X" const blocks and Valid() methods for
//     enum types aliased to time.Time or openapi_types.Date - Go cannot have
//     constants of those types.
//  5. Run: goimports -w client.gen.go
//
// All steps are idempotent; see the regen-client target in
// GNUmakefile for the exact python/awk implementation.
