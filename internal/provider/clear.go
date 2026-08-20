package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
)

// Every generated request struct field is a pointer with `omitempty`, so a nil
// pointer serialises to nothing and DRF leaves the column untouched. Omission
// therefore never clears a value.
//
// The only way to clear a column through this API is a PATCH carrying an
// explicit JSON null. Verified against DefectDojo 3.1.101: PUT with the field
// omitted leaves it unchanged, PATCH {"branch_tag": null} clears it. So Update
// diffs the plan against the prior state, and any attribute that held a value
// and is now null gets an explicit null in a follow-up PATCH (issue #30).
//
// This only concerns Optional-without-Computed attributes. For Optional+Computed
// ones, Terraform core fills a null config with the prior state value before the
// provider ever sees the plan, so "removed from config" and "unchanged" are
// indistinguishable by construction.

// partialUpdateFunc is the shape every generated `<Tag>PartialUpdateWithBody`
// method shares on dd.ClientInterface, which dd.ClientWithResponses embeds.
type partialUpdateFunc func(ctx context.Context, id int, contentType string, body io.Reader, reqEditors ...dd.RequestEditorFn) (*http.Response, error)

// nullClearer is an optional interface for defectdojoResource implementations
// that support clearing attributes via a PATCH with an explicit null body.
//
// Implementations name the raw client method, not the ...WithResponse wrapper.
// The wrapper parses the PATCH response into the resource's own model, which
// applyClearedFields then discards in favour of a readApiCall - and for
// defectdojo_notifications the parse fails outright, because
// ParseNotificationsPartialUpdateResponse unmarshals into dd.Notifications,
// whose ScanAddedEmpty is a scalar enum while the live API answers with an
// array. Naming the raw method keeps all fourteen resources uniform.
//
// TestResourcesWithClearableAttributesImplementNullClearer fails the build if a
// resource has an attribute that can be cleared but no implementation here.
type nullClearer interface {
	partialUpdateWithBody(client *dd.ClientWithResponses) partialUpdateFunc
}

// sendClearPatch issues the PATCH and returns the raw status and body.
func sendClearPatch(ctx context.Context, patch partialUpdateFunc, typeName string, idNumber int, body []byte) (int, []byte, error) {
	tflog.Info(ctx, typeName+" clearFieldsApiCall")

	httpResp, err := patch(ctx, idNumber, "application/json", bytes.NewReader(body))
	if err != nil {
		return 0, nil, err
	}
	defer func() { _ = httpResp.Body.Close() }()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return 0, nil, err
	}
	tflog.Info(ctx, fmt.Sprintf("response %s: %s", httpResp.Status, respBody))
	return httpResp.StatusCode, respBody, nil
}

// clearTarget is one attribute the practitioner removed from configuration,
// resolved down to the wire field that has to be cleared.
//
// tfsdkName is carried alongside jsonName because the two frequently disagree -
// `product_manager_id` is `product_manager` on the wire, `regulation_ids` is
// `regulations` - and a diagnostic naming the wire field sends the reader
// looking for an argument this provider does not have. Anything addressed to
// the practitioner uses tfsdkName; the PATCH body stays as the server saw it.
type clearTarget struct {
	tfsdkName   string // the Terraform attribute name, as written in configuration
	jsonName    string // the DefectDojo JSON field name
	ddFieldName string // the ddclient Go field name, used to re-check after the update
	isSlice     bool   // collections clear to [], scalars to null
}

// stillSetAfterUpdate narrows the clear list to fields the update response shows
// are actually still populated.
//
// Where a field is declared with an explicit DRF default, the serializer
// substitutes that default for an omitted field and the PUT has already cleared
// it. A redundant PATCH is then at best a wasted request and at worst a failure:
// defectdojo_metadata validates across fields, so a lone {"product": null} is
// rejected with "Metadata entries need either a product, endpoint, location or
// a finding" even though the row is already correct.
//
// The update response is authoritative here: updateApiCall assigns it back onto
// the wrapper, so a nil (or empty) field means the server already dropped it.
func stillSetAfterUpdate(ddResource defectdojoResource, targets []clearTarget) []clearTarget {
	ddVal := reflect.ValueOf(ddResource).Elem()

	var out []clearTarget
	for _, target := range targets {
		field := ddVal.FieldByName(target.ddFieldName)
		if !field.IsValid() {
			out = append(out, target) // cannot tell; let the PATCH decide
			continue
		}
		if ddFieldIsEmpty(field) {
			continue // the update already cleared it
		}
		out = append(out, target)
	}
	return out
}

// ddFieldIsEmpty reports whether a ddclient field currently holds no value.
//
// "No value" means nil, never zero. DefectDojo stores 0, false and "" verbatim
// rather than coercing them to null - verified on 3.1.101, where
// PATCH {"percent_complete": 0, "branch_tag": ""} reads back as 0 and "" - so a
// non-nil pointer to a zero value is a value the server chose to send, and the
// clear target has to survive.
func ddFieldIsEmpty(field reflect.Value) bool {
	switch field.Kind() {
	case reflect.Ptr:
		// Nil is the whole test. The pointee still gets looked at because the
		// collections arrive as *[]T, and there an empty slice does mean empty.
		if field.IsNil() {
			return true
		}
		return ddFieldIsEmpty(field.Elem())
	case reflect.Slice, reflect.Map:
		// The spec types tags and the various id sets as non-nullable arrays,
		// so an empty array is how the server says "no elements".
		return field.IsNil() || field.Len() == 0
	default:
		// A dereferenced scalar: the pointer was not nil, so this is a value.
		return false
	}
}

// clearedDdFields returns the attributes that hold a value in the prior state
// but are null in the plan - that is, the ones the practitioner removed from
// configuration.
//
// Set-typed attributes are flagged separately because an empty collection, not
// null, is how the API expresses "no elements".
func clearedDdFields(plan, state terraformResourceData, ddResource defectdojoResource) []clearTarget {
	planVal := reflect.ValueOf(plan).Elem()
	stateVal := reflect.ValueOf(state).Elem()
	ddType := reflect.ValueOf(ddResource).Elem().Type()

	var targets []clearTarget

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

		tfsdkName := fieldDescriptor.Tag.Get("tfsdk")
		if tfsdkName == "" {
			tfsdkName = jsonName // nothing better to show; every field should have one
		}

		ddFieldType := ddField.Type
		if ddFieldType.Kind() == reflect.Ptr {
			ddFieldType = ddFieldType.Elem()
		}

		targets = append(targets, clearTarget{
			tfsdkName:   tfsdkName,
			jsonName:    jsonName,
			ddFieldName: ddFieldName,
			isSlice:     ddFieldType.Kind() == reflect.Slice,
		})
	}

	return targets
}

// clearPatchBody builds the PATCH payload that clears the given fields.
func clearPatchBody(targets []clearTarget) ([]byte, error) {
	body := map[string]any{}
	for _, target := range targets {
		if target.isSlice {
			body[target.jsonName] = []any{}
		} else {
			body[target.jsonName] = nil
		}
	}

	return json.Marshal(body)
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
	targets []clearTarget,
) {
	// Spelled the way the configuration spells them: the wire names below are
	// for the request dump, not for the practitioner.
	names := make([]string, 0, len(targets))
	for _, t := range targets {
		names = append(names, t.tfsdkName)
	}
	all := strings.Join(names, ", ")

	clearer, ok := ddResource.(nullClearer)
	if !ok {
		diags.AddError(
			"Cannot Clear Attributes on "+typeName,
			fmt.Sprintf("The attributes %s were removed from the configuration, but %s does not "+
				"implement nullClearer, so the provider cannot ask DefectDojo to clear them. "+
				"Omitting a field from an update request leaves it unchanged. This is a provider "+
				"bug; please report it.", all, typeName))
		return
	}

	body, err := clearPatchBody(targets)
	if err != nil {
		diags.AddError("Error Clearing Attributes on "+typeName, err.Error())
		return
	}

	statusCode, respBody, err := sendClearPatch(ctx, clearer.partialUpdateWithBody(client), typeName, idNumber, body)
	if err != nil {
		diags.AddError("Error Clearing Attributes on "+typeName, err.Error())
		return
	}
	if statusCode != 200 {
		diags.AddError(
			"API Error Clearing Attributes on "+typeName,
			fmt.Sprintf("Removing %s from the configuration requires DefectDojo to accept an "+
				"explicit null for those fields, but it answered %d. Put them back in the "+
				"configuration, or destroy and recreate the resource without them."+
				"\n\nrequest:\n\n%s\n\nbody:\n\n%s",
				all, statusCode, string(body), string(respBody)))
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
