package ddclient

// To regenerate the client from an OpenAPI spec, run:
//
//	make regen-client DD_VERSION=<version>    # e.g. DD_VERSION=3.1.101
//
// Prerequisite: openapi-specs/<version>/defect_dojo.json must exist locally.
// Collected specs are LOCAL artifacts, not tracked in git - collect one first
// with `make dd-up && make dd-spec` (optionally DD_VERSION=<version>).
//
// The target encodes the full validated procedure:
//
//  1. Copy openapi-specs/<version>/defect_dojo.json to the repo root as
//     defect_dojo.json (the root copy is gitignored; do not commit it).
//  2. Run the oapi-codegen version pinned in go.mod (never a global binary):
//     go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@<pinned> \
//     -generate types,client,skip-fmt -package ddclient -o client.gen.go defect_dojo.json
//     skip-fmt is required: the raw output contains invalid Go that gofmt
//     rejects, which the post-processing below cleans up. (For the same
//     reason the yaml config in this directory cannot generate as-is.)
//  3. Delete invalid bare-nil enum constants: sed -i '/ = <nil>$/d' client.gen.go
//     String enum members with the quoted value "<nil>" are valid Go and are kept.
//  4. Remove the "Defines values for X" const blocks and Valid() methods for
//     enum types aliased to time.Time or openapi_types.Date - Go cannot have
//     constants of those types. This also removes the switch cases that would
//     otherwise be orphaned by step 3 (every bare-nil constant belongs to one
//     of these blocks).
//  5. Run: goimports -w client.gen.go
//
// All post-processing steps are idempotent; see the regen-client target in
// GNUmakefile for the exact sed/awk implementation.
