package provider

// This file contains a permanent COVERAGE audit that complements
// ddfield_audit_unit_test.go. Where that file checks "every ddField tag
// resolves and type-pairs correctly", this file checks the opposite
// direction: "every exported field on the generated ddclient struct is
// either ddField-mapped or has a documented, honest reason for being
// skipped".
//
// Why this exists: it is easy for a generated-client field to go completely
// unnoticed (never mapped, never excused) when a resource is first written,
// and easy for a NEW field to appear silently after `make regen-client`
// bumps to a newer DefectDojo spec. Without this test, such gaps are only
// found by manually diffing structs. This test makes them a `go test`
// failure: any exported field reachable by the reflection engine that is
// neither tagged via `ddField` nor listed in ddFieldKnownUnmapped fails the
// build.
//
// HOW TO FIX A FAILURE: figure out why the new/existing field isn't mapped.
// If it should be exposed as a Terraform attribute, map it (add the ddField
// tag, schema attribute, and wire the 4 unit test funcs + acceptance test).
// If it genuinely should not be exposed, add an entry to
// ddFieldKnownUnmapped with an honest, specific reason - verify against the
// generated struct and its corresponding *Request struct in
// internal/ddclient/client.gen.go before excusing anything as "read-only".

import (
	"reflect"
	"sort"
	"testing"
)

// ddFieldKnownUnmapped documents generated-client struct fields that are
// deliberately NOT mapped to Terraform attributes, with the reason.
var ddFieldKnownUnmapped = map[string]map[string]string{ // model -> DD field -> reason
	"ddTestResourceData": {
		"ApiScanConfiguration": "writable but not yet exposed — candidate for follow-up (see Phase 2 report)",
		"Created":              "server-managed timestamp",
		"Files":                "read-only server field, not in TestCreateRequest (files are attached via a separate upload endpoint)",
		"Notes":                "*[]*int — element type unsupported",
		"Updated":              "server-managed timestamp",
	},
	"endpointResourceData": {
		"ActiveFindingCount": "read-only server field, not in a request struct (endpoint is a read-only projection since DefectDojo 3.0; see CLAUDE.md exclusion table)",
		"Created":            "server-managed timestamp",
		"LocationId":         "read-only server field, not in a request struct (endpoint is a read-only projection since DefectDojo 3.0)",
		"Relationship":       "read-only server field, not in a request struct (endpoint is a read-only projection since DefectDojo 3.0)",
		"RelationshipData":   "interface{} (free-form JSON) — engine cannot map",
		"Status":             "read-only server field, not in a request struct (endpoint is a read-only projection since DefectDojo 3.0)",
		"Tags":               "read-only server field, not in a request struct (endpoint is a read-only projection since DefectDojo 3.0)",
		"Updated":            "server-managed timestamp",
	},
	"engagementPresetResourceData": {
		"Created":  "server-managed timestamp",
		"Prefetch": "prefetch embellishment, not a model field",
	},
	"engagementResourceData": {
		"Active":                     "read-only server field, not in EngagementRequest (derived from status)",
		"BuildServer":                "writable but not yet exposed — candidate for follow-up (see Phase 2 report)",
		"Created":                    "server-managed timestamp",
		"DoneTesting":                "read-only server field, not in EngagementRequest (derived from status)",
		"Files":                      "nested struct — engine cannot map",
		"Notes":                      "nested struct — engine cannot map",
		"OrchestrationEngine":        "writable but not yet exposed — candidate for follow-up (see Phase 2 report)",
		"Progress":                   "read-only server field, not in EngagementRequest",
		"RiskAcceptance":             "read-only server field, not in EngagementRequest (risk acceptances link to an engagement, not the reverse)",
		"SourceCodeManagementServer": "writable but not yet exposed — candidate for follow-up (see Phase 2 report)",
		"TmodelPath":                 "read-only server field, not in EngagementRequest",
		"Updated":                    "server-managed timestamp",
	},
	"findingTemplateResourceData": {
		"Endpoints":         "read-only server field, not in FindingTemplateRequest (endpoints_text is the writable form)",
		"LastUsed":          "server-managed timestamp",
		"NumericalSeverity": "read-only server field, not in FindingTemplateRequest (server-computed from severity)",
		"Tags":              "writable but not yet exposed — candidate for follow-up (see Phase 2 report)",
		"VulnerabilityIds":  "read-only server field, not in FindingTemplateRequest",
	},
	"jiraProductConfigurationResourceData": {
		"CloseTransitionFields":  "interface{} (free-form JSON) — engine cannot map",
		"CustomFields":           "interface{} (free-form JSON) — engine cannot map",
		"Prefetch":               "prefetch embellishment, not a model field",
		"ReopenTransitionFields": "interface{} (free-form JSON) — engine cannot map",
	},
	"locationProductResourceData": {
		"Created":          "server-managed timestamp",
		"RelationshipData": "interface{} (free-form JSON) — engine cannot map",
		"Updated":          "server-managed timestamp",
	},
	"metadataResourceData": {
		"Endpoint": "writable in MetaRequest, but the DefectDojo 3.1.101 API rejects this parent when set; intentionally not exposed (see metadata_resource.go schema description)",
		"Location": "writable in MetaRequest, but the DefectDojo 3.1.101 API silently ignores this parent when set; intentionally not exposed (see metadata_resource.go schema description)",
	},
	"notificationWebhookResourceData": {
		"FirstError": "server-managed timestamp",
		"LastError":  "server-managed timestamp",
		"Note":       "read-only server field, not in NotificationWebhooksRequest (server-populated error description)",
	},
	"notificationsResourceData": {
		"Prefetch": "prefetch embellishment, not a model field",
	},
	"productAPIScanConfigurationResourceData": {
		"Prefetch": "prefetch embellishment, not a model field",
	},
	"productResourceData": {
		"Created":       "server-managed timestamp",
		"FindingsCount": "read-only server field, not in ProductRequest",
		"FindingsList":  "read-only server field, not in ProductRequest",
		"Prefetch":      "prefetch embellishment, not a model field",
		"ProductMeta":   "nested struct — engine cannot map",
	},
	"productTypeResourceData": {
		"Created":  "server-managed timestamp",
		"Prefetch": "prefetch embellishment, not a model field",
		"Updated":  "server-managed timestamp",
	},
	"riskAcceptanceResourceData": {
		"Created":               "server-managed timestamp",
		"ExpirationDateHandled": "read-only server field per API docs (marked readonly in the OpenAPI schema; appears on RiskAcceptanceRequest only because oapi-codegen does not strip DRF read_only fields from generated request structs)",
		"ExpirationDateWarned":  "read-only server field per API docs (marked readonly in the OpenAPI schema; appears on RiskAcceptanceRequest only because oapi-codegen does not strip DRF read_only fields from generated request structs)",
		"Notes":                 "read-only server field, not in RiskAcceptanceRequest",
		"Path":                  "read-only server field, not in RiskAcceptanceRequest (server-generated path to the uploaded proof file)",
		"Updated":               "server-managed timestamp",
	},
	"toolConfigurationResourceData": {
		"Prefetch": "prefetch embellishment, not a model field",
	},
	"toolProductSettingsResourceData": {
		"Notes":    "read-only server field, not in ToolProductSettingsRequest",
		"Prefetch": "prefetch embellishment, not a model field",
	},
	"urlResourceData": {
		"String": "redundant read-only string form (URL.String; the resource already exposes the parsed host/port/path/protocol/query/fragment fields)",
	},
	"userContactInfoResourceData": {
		"PasswordLastReset": "server-managed timestamp",
		"Prefetch":          "prefetch embellishment, not a model field",
		"TokenLastReset":    "server-managed timestamp",
		"UiUseTailwind":     "UI-only preference",
		"UserProfile":       "nested struct — engine cannot map",
	},
	"userProfileResourceData": {
		"ConfigurationPermissions": "*[]*int — element type unsupported",
		"PasswordLastReset":        "server-managed timestamp",
		"TokenLastReset":           "server-managed timestamp",
	},
	"userResourceData": {
		"ConfigurationPermissions": "*[]*int — element type unsupported",
		"DateJoined":               "server-managed timestamp",
		"LastLogin":                "server-managed timestamp",
		"PasswordLastReset":        "server-managed timestamp",
		"TokenLastReset":           "server-managed timestamp",
	},
}

// collectDdFields recursively collects the set of exported field names
// reachable on ddStruct by the reflection engine in resource.go: fields
// declared directly on the struct, plus - for anonymous (embedded) struct
// fields (or pointers to structs) - their exported fields too. This mirrors
// what reflect.Value.FieldByName's promotion sees.
func collectDdFields(t reflect.Type, seen map[string]bool, out *[]string) {
	for f := range t.Fields() {
		if f.Anonymous {
			et := f.Type
			if et.Kind() == reflect.Ptr {
				et = et.Elem()
			}
			if et.Kind() == reflect.Struct {
				collectDdFields(et, seen, out)
				continue
			}
		}
		if f.IsExported() && !seen[f.Name] {
			seen[f.Name] = true
			*out = append(*out, f.Name)
		}
	}
}

// TestDdFieldCoverage asserts that every exported field on every
// defectdojoResource() wrapper struct in ddFieldAuditTable is either mapped
// to a Terraform attribute via a `ddField` tag, or has a documented reason in
// ddFieldKnownUnmapped for why it is intentionally not mapped. This is a
// permanent regression gate: after `make regen-client` picks up new or
// renamed fields from a newer DefectDojo spec, any newly-appeared field must
// be explicitly triaged (mapped or excused) before this test passes again.
func TestDdFieldCoverage(t *testing.T) {
	totalAudited, totalMapped, totalExcused := 0, 0, 0

	for name, data := range ddFieldAuditTable {
		t.Run(name, func(t *testing.T) {
			tfType := reflect.TypeOf(data)
			if tfType.Kind() != reflect.Ptr || tfType.Elem().Kind() != reflect.Struct {
				t.Fatalf("%s: table entry must be a pointer to a struct, got %s", name, tfType)
			}
			tfStruct := tfType.Elem()

			ddStruct := ddStructType(t, name, data.defectdojoResource())

			mappedTags := map[string]bool{}
			for tfField := range tfStruct.Fields() {
				if tag := tfField.Tag.Get("ddField"); tag != "" {
					mappedTags[tag] = true
				}
			}

			var ddFields []string
			collectDdFields(ddStruct, map[string]bool{}, &ddFields)
			sort.Strings(ddFields)

			for _, field := range ddFields {
				totalAudited++
				if mappedTags[field] {
					totalMapped++
					continue
				}
				if reason, ok := ddFieldKnownUnmapped[name][field]; ok && reason != "" {
					totalExcused++
					continue
				}
				t.Errorf("%s: dd field %s is neither ddField-mapped nor excused in ddFieldKnownUnmapped", name, field)
			}
		})
	}

	t.Logf("audited %d dd fields across %d models (%d mapped, %d excused)", totalAudited, len(ddFieldAuditTable), totalMapped, totalExcused)
}
