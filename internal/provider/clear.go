package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
)

// Clearing an attribute on Update.
//
// Every generated request struct field is a pointer with `omitempty`, so a nil
// pointer serialises to nothing at all and DRF leaves the corresponding column
// untouched. Omission therefore never clears a value: an attribute that was
// once set and is then deleted from configuration keeps its old server-side
// value, the read path writes that value back into state, and Terraform fails
// the apply with
//
//	Provider produced inconsistent result after apply:
//	.branch_tag: was null, but now cty.StringVal("main")
//
// The error is retry-proof and fires after the server has already been
// mutated (GitHub issue #30).
//
// The only way to clear a column through this API is a PATCH carrying an
// explicit JSON null - verified against DefectDojo 3.1.101, where PUT with the
// field omitted leaves it unchanged and PATCH {"branch_tag": null} clears it.
// So Update diffs the plan against the prior state, and any attribute that held
// a value and is now null gets an explicit null in a follow-up PATCH.
//
// This only concerns Optional-without-Computed attributes. For Optional+Computed
// ones, Terraform core fills a null config with the prior state value before the
// provider ever sees the plan, so "removed from config" and "unchanged" are
// indistinguishable by construction - which is what Optional+Computed means.

// nullClearer is an optional interface for defectdojoResource implementations
// that support clearing attributes via a PATCH with an explicit null body.
//
// Implementations only send the request; they deliberately do not parse the
// response into the wrapper. The engine follows a successful clear with
// readApiCall, which every resource already implements correctly, so this stays
// a uniform seven lines per resource even where the PATCH response model
// differs from the one the wrapper embeds.
//
// TestResourcesWithClearableAttributesImplementNullClearer fails the build if a
// resource has an attribute that can be cleared but no implementation here.
type nullClearer interface {
	clearFieldsApiCall(ctx context.Context, client *dd.ClientWithResponses, idNumber int, body []byte) (int, []byte, error)
}

// clearedDdFields returns the DefectDojo JSON field names for attributes that
// hold a value in the prior state but are null in the plan - that is, the ones
// the practitioner removed from configuration.
//
// Set-typed attributes are reported separately because an empty collection, not
// null, is how the API expresses "no elements": the spec types tags and the
// various *_ids as non-nullable arrays.
func clearedDdFields(plan, state terraformResourceData, ddResource defectdojoResource) (nulls []string, empties []string) {
	planVal := reflect.ValueOf(plan).Elem()
	stateVal := reflect.ValueOf(state).Elem()
	if planVal.Type() != stateVal.Type() {
		return nil, nil
	}
	ddType := reflect.ValueOf(ddResource).Elem().Type()

	for i := range planVal.NumField() {
		fieldDescriptor := planVal.Type().Field(i)
		ddFieldName := fieldDescriptor.Tag.Get("ddField")
		if ddFieldName == "" {
			continue
		}

		planAttr, ok := planVal.Field(i).Interface().(attr.Value)
		if !ok {
			continue
		}
		stateAttr, ok := stateVal.Field(i).Interface().(attr.Value)
		if !ok {
			continue
		}

		// Only a plan that is definitively null means "cleared". An unknown
		// plan value is one the server will fill in, and a state value that is
		// already null has nothing to clear.
		if planAttr.IsUnknown() || !planAttr.IsNull() {
			continue
		}
		if stateAttr.IsNull() || stateAttr.IsUnknown() {
			continue
		}

		ddField, ok := ddType.FieldByName(ddFieldName)
		if !ok {
			continue // reported by populateDefectdojoResource
		}
		jsonName, _, _ := strings.Cut(ddField.Tag.Get("json"), ",")
		if jsonName == "" || jsonName == "-" {
			continue
		}

		if isSliceDdField(ddField.Type) {
			empties = append(empties, jsonName)
		} else {
			nulls = append(nulls, jsonName)
		}
	}

	return nulls, empties
}

// isSliceDdField reports whether a ddclient field is a collection (or a pointer
// to one).
func isSliceDdField(t reflect.Type) bool {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Kind() == reflect.Slice
}

// clearPatchBody builds the PATCH payload that clears the given fields.
func clearPatchBody(nulls []string, empties []string) ([]byte, error) {
	body := map[string]any{}
	for _, name := range nulls {
		body[name] = nil
	}
	for _, name := range empties {
		body[name] = []any{}
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(body); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// applyClearedFields issues the follow-up PATCH that clears attributes removed
// from configuration, then refreshes the wrapper from the server so state
// reflects what was actually stored.
func applyClearedFields(
	ctx context.Context,
	diags *diag.Diagnostics,
	client *dd.ClientWithResponses,
	typeName string,
	idNumber int,
	ddResource defectdojoResource,
	nulls []string,
	empties []string,
) {
	all := append(append([]string{}, nulls...), empties...)

	clearer, ok := ddResource.(nullClearer)
	if !ok {
		diags.AddError(
			"Cannot Clear Attributes on "+typeName,
			fmt.Sprintf("The attributes %s were removed from the configuration, but %s does not "+
				"implement clearFieldsApiCall, so the provider cannot ask DefectDojo to clear them. "+
				"Omitting a field from an update request leaves it unchanged. This is a provider "+
				"bug; please report it.", strings.Join(all, ", "), typeName))
		return
	}

	body, err := clearPatchBody(nulls, empties)
	if err != nil {
		diags.AddError("Error Clearing Attributes on "+typeName, err.Error())
		return
	}

	statusCode, respBody, err := clearer.clearFieldsApiCall(ctx, client, idNumber, body)
	if err != nil {
		diags.AddError("Error Clearing Attributes on "+typeName, err.Error())
		return
	}
	if statusCode != 200 {
		diags.AddError(
			"API Error Clearing Attributes on "+typeName,
			fmt.Sprintf("Removing %s from the configuration requires DefectDojo to accept an "+
				"explicit null for those fields, but it answered %d.\n\nrequest:\n\n%s\n\nbody:\n\n%s",
				strings.Join(all, ", "), statusCode, string(body), string(respBody)))
		return
	}

	// The PATCH response model does not always match the struct the wrapper
	// embeds, so re-read rather than parsing it per resource.
	statusCode, respBody, err = ddResource.readApiCall(ctx, client, idNumber)
	if err != nil {
		diags.AddError("Error Re-reading "+typeName+" After Clearing Attributes", err.Error())
		return
	}
	if statusCode != 200 {
		diags.AddError(
			"API Error Re-reading "+typeName+" After Clearing Attributes",
			fmt.Sprintf("Unexpected response code from API: %d\n\nbody:\n\n%s", statusCode, string(respBody)))
	}
}
