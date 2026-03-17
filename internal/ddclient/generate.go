package ddclient

// To regenerate the client from the OpenAPI spec:
//   1. Run: oapi-codegen -generate types,client,skip-fmt -package ddclient -o client.gen.go ../../defect_dojo.json
//   2. Fix invalid bare <nil> constants: sed -i '/ = <nil>$/d' client.gen.go
//   3. Fix orphaned return statements in switch blocks after <nil> removal
//   4. Remove time.Time-based const blocks and their Valid() methods (time.Time cannot be a Go constant)
//   5. Run: goimports -w client.gen.go
//
// See the Makefile or CLAUDE.md for the full regeneration procedure.
