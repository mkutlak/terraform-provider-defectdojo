package provider

import (
	"context"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type terraformResourceData interface {
	id() types.String
	setId(types.String)
	defectdojoResource() defectdojoResource
}

type defectdojoResource interface {
	createApiCall(context.Context, *dd.ClientWithResponses) (int, []byte, error)
	readApiCall(context.Context, *dd.ClientWithResponses, int) (int, []byte, error)
	updateApiCall(context.Context, *dd.ClientWithResponses, int) (int, []byte, error)
	deleteApiCall(context.Context, *dd.ClientWithResponses, int) (int, []byte, error)
}

// singletonAdopter is an optional interface for defectdojoResource
// implementations backed by an API object that always exists server-side and
// can be neither created nor deleted (e.g. system_settings). On Create the
// engine "adopts" the existing object: it resolves the id via adoptApiCall and
// then applies the plan with updateApiCall instead of calling createApiCall.
// On Delete the engine only removes the resource from Terraform state; the
// server-side object keeps its last-applied values.
type singletonAdopter interface {
	// adoptApiCall resolves the singleton's id (e.g. via the list endpoint)
	// without modifying the object. It must return 200 on success.
	adoptApiCall(context.Context, *dd.ClientWithResponses) (id int, statusCode int, body []byte, err error)
}
type dataProvider interface {
	getData(context.Context, dataGetter) (terraformResourceData, diag.Diagnostics)
}

type terraformResource struct {
	client   *dd.ClientWithResponses
	typeName string
	dataProvider
}

type dataGetter interface {
	Get(context.Context, any) diag.Diagnostics
}

var typeOfTypesString = reflect.TypeFor[types.String]()
var typeOfTypesBool = reflect.TypeFor[types.Bool]()
var typeOfTypesInt64 = reflect.TypeFor[types.Int64]()
var typeOfTypesFloat64 = reflect.TypeFor[types.Float64]()
var typeOfTypesSet = reflect.TypeFor[types.Set]()

// ddFormatDecimal marks a types.String attribute whose DefectDojo column is a
// Django DecimalField. The server answers in one canonical form, so the read
// path keeps the configured literal when it denotes the same amount rather
// than rewriting state and tripping "Provider produced inconsistent result
// after apply". See decimal.go.
const ddFormatDecimal = "decimal"

// knownDdFormats is the set of `ddFormat` struct tag values populateResourceData
// understands. TestDdFormatTagsAreKnown turns a typo into a `go test` failure
// instead of an apply-time diagnostic nobody reads.
var knownDdFormats = map[string]bool{ddFormatDecimal: true}

// addUnsupportedMappingError reports a (Terraform type, ddclient type) pairing
// the reflection engine cannot convert.
//
// Every such pairing is a provider bug - almost always generated-client drift
// after `make regen-client` - and quietly logging it drops the practitioner's
// value into a file nobody reads. That silent-failure half of the issue #23
// root cause is why these are diagnostics rather than tflog.Warn calls.
//
// Making them loud is safe precisely because TestDdFieldAuditTagsResolve proves
// no currently-mapped field can reach them; that audit is load-bearing, not
// advisory.
func addUnsupportedMappingError(diags *diag.Diagnostics, fn, tfsdkName string, tfType, ddType reflect.Type) {
	diags.AddError("Unsupported ddField Mapping", fmt.Sprintf(
		"%s cannot convert between the Terraform type %s and the DefectDojo client type %s "+
			"for attribute %q. This is a provider bug, usually generated-client drift after "+
			"`make regen-client`; please report it, including the DefectDojo version in use.",
		fn, tfType, ddType, tfsdkName))
}

// addUnhandledTerraformTypeError reports a Terraform attribute type the engine
// does not handle at all, as opposed to a type it handles but cannot pair with
// the given ddclient field.
func addUnhandledTerraformTypeError(diags *diag.Diagnostics, fn, tfsdkName string, tfType, ddType reflect.Type) {
	diags.AddError("Unsupported Terraform Attribute Type", fmt.Sprintf(
		"%s does not handle the Terraform type %s at all (attribute %q, DefectDojo client type "+
			"%s). Only types.String, types.Bool, types.Int64, types.Float64 and types.Set are "+
			"supported. This is a provider bug; please report it.",
		fn, tfType, tfsdkName, ddType))
}

// renderStringValue converts a server-provided string into the value to store
// in state, honouring the optional `ddFormat` struct tag.
//
// An unrecognised ddFormat is an error rather than a silent pass-through: a tag
// that does not do what it says is exactly the kind of quiet degradation that
// let issue #23 survive.
func renderStringValue(diags *diag.Diagnostics, tag reflect.StructTag, current types.String, server string) types.String {
	switch format := tag.Get("ddFormat"); format {
	case "":
		return types.StringValue(server)
	case ddFormatDecimal:
		return preserveDecimalLiteral(current, server)
	default:
		diags.AddError("Unknown ddFormat Tag", fmt.Sprintf(
			"Attribute %q carries ddFormat:%q, which populateResourceData does not understand "+
				"(valid values: %s). This is a provider bug; please report it.",
			tag.Get("tfsdk"), format, strings.Join(slices.Sorted(maps.Keys(knownDdFormats)), ", ")))
		return types.StringValue(server)
	}
}

func (r *terraformResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*dd.ClientWithResponses)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected dd.ClientWithResponses, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r terraformResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	data, diags := r.getData(ctx, req.Plan)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured HTTP Client",
			"Expected configured HTTP client. Please report this issue to the provider developers.",
		)

		return
	}

	ddResource := data.defectdojoResource()
	populateDefectdojoResource(ctx, &resp.Diagnostics, data, &ddResource)

	// A conversion failure (e.g. an unparseable datetime) must abort before the
	// API call: the ddclient struct would otherwise carry a zero value that
	// DefectDojo accepts, producing a state/config mismatch.
	if resp.Diagnostics.HasError() {
		return
	}

	if sa, ok := ddResource.(singletonAdopter); ok {
		id, statusCode, body, err := sa.adoptApiCall(ctx, r.client)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Adopting "+r.typeName,
				err.Error())
			return
		}
		if statusCode != 200 {
			resp.Diagnostics.AddError(
				"API Error Adopting "+r.typeName,
				fmt.Sprintf("Unexpected response code from API: %d", statusCode)+
					fmt.Sprintf("\n\nbody:\n\n%s", string(body)),
			)
			return
		}

		statusCode, body, err = ddResource.updateApiCall(ctx, r.client, id)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Updating Adopted "+r.typeName,
				err.Error())
			return
		}
		if statusCode != 200 {
			resp.Diagnostics.AddError(
				"API Error Updating Adopted "+r.typeName,
				fmt.Sprintf("Unexpected response code from API: %d", statusCode)+
					fmt.Sprintf("\n\nbody:\n\n%s", string(body)),
			)
			return
		}

		populateResourceData(ctx, &resp.Diagnostics, &data, ddResource)

		tflog.Trace(ctx, "singleton resource adopted")

		// The singleton has already been updated server-side, so state must be
		// written even if populateResourceData reported problems - bailing out
		// here would lose the adoption.
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	statusCode, body, err := ddResource.createApiCall(ctx, r.client)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating "+r.typeName,
			err.Error())
		return
	}

	if statusCode == 201 {
		populateResourceData(ctx, &resp.Diagnostics, &data, ddResource)
	} else {
		resp.Diagnostics.AddError(
			"API Error Creating "+r.typeName,
			fmt.Sprintf("Unexpected response code from API: %d", statusCode)+
				fmt.Sprintf("\n\nbody:\n\n%s", string(body)),
		)
		return
	}

	tflog.Trace(ctx, "resource created")

	// The resource exists server-side from here on; write state even if
	// populateResourceData reported problems, so it can be managed/destroyed.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *terraformResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data, diags := r.getData(ctx, req.State)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.id().IsNull() {
		resp.Diagnostics.AddError(
			"Could not Retrieve "+r.typeName,
			"The Id field was null but it is required to retrieve the resource")
		return
	}

	idNumber, err := strconv.Atoi(data.id().ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not Retrieve "+r.typeName,
			"Error while parsing the resource ID from state: "+err.Error())
		return
	}

	ddResource := data.defectdojoResource()
	populateDefectdojoResource(ctx, &resp.Diagnostics, data, &ddResource)

	if resp.Diagnostics.HasError() {
		return
	}

	statusCode, body, err := ddResource.readApiCall(ctx, r.client, idNumber)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving "+r.typeName,
			err.Error())
		return
	}

	switch statusCode {
	case 200:
		populateResourceData(ctx, &resp.Diagnostics, &data, ddResource)
	case 404:
		resp.State.RemoveResource(ctx)
		return
	default:
		resp.Diagnostics.AddError(
			"API Error Retrieving "+r.typeName,
			fmt.Sprintf("Unexpected response code from API: %d", statusCode)+
				fmt.Sprintf("\n\nbody:\n\n%+v", string(body)),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r terraformResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data, diags := r.getData(ctx, req.Plan)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured HTTP Client",
			"Expected configured HTTP client. Please report this issue to the provider developers.",
		)

		return
	}

	if data.id().IsNull() {
		resp.Diagnostics.AddError(
			"Could not Update "+r.typeName,
			"The Id field was null but it is required to update the resource")
		return
	}

	idNumber, err := strconv.Atoi(data.id().ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not Update "+r.typeName,
			"Error while parsing the resource ID from state: "+err.Error())
		return
	}

	ddResource := data.defectdojoResource()
	populateDefectdojoResource(ctx, &resp.Diagnostics, data, &ddResource)

	if resp.Diagnostics.HasError() {
		return
	}

	// Attributes the practitioner removed from configuration need an explicit
	// null; omitting them from the update request would leave them unchanged.
	// See clear.go and GitHub issue #30.
	priorState, stateDiags := r.getData(ctx, req.State)
	resp.Diagnostics.Append(stateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	clearNulls, clearEmpties := clearedDdFields(data, priorState, ddResource)

	statusCode, body, err := ddResource.updateApiCall(ctx, r.client, idNumber)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating "+r.typeName,
			err.Error())
		return
	}

	if statusCode == 200 {
		if len(clearNulls) > 0 || len(clearEmpties) > 0 {
			applyClearedFields(ctx, &resp.Diagnostics, r.client, r.typeName, idNumber, ddResource, clearNulls, clearEmpties)
			if resp.Diagnostics.HasError() {
				return
			}
		}
		populateResourceData(ctx, &resp.Diagnostics, &data, ddResource)
	} else {
		resp.Diagnostics.AddError(
			"API Error Updating "+r.typeName,
			fmt.Sprintf("Unexpected response code from API: %d", statusCode)+
				fmt.Sprintf("\n\nbody:\n\n%+v", string(body)),
		)
		return
	}

	// The resource exists server-side; write state even if populateResourceData
	// reported problems, so it can be managed/destroyed.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r terraformResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data, diags := r.getData(ctx, req.State)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured HTTP Client",
			"Expected configured HTTP client. Please report this issue to the provider developers.",
		)

		return
	}

	if data.id().IsNull() {
		resp.Diagnostics.AddError(
			"Could not Delete "+r.typeName,
			"The Id field was null but it is required to delete the resource")
		return
	}

	idNumber, err := strconv.Atoi(data.id().ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not Delete "+r.typeName,
			"Error while parsing the resource ID from state: "+err.Error())
		return
	}

	ddResource := data.defectdojoResource()

	if _, ok := ddResource.(singletonAdopter); ok {
		// Singletons cannot be deleted server-side; only forget them from
		// Terraform state. The object keeps its last-applied values.
		resp.State.RemoveResource(ctx)
		return
	}

	statusCode, body, err := ddResource.deleteApiCall(ctx, r.client, idNumber)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting "+r.typeName,
			err.Error())
		return
	}

	if statusCode != 204 {
		resp.Diagnostics.AddError(
			"API Error Deleting "+r.typeName,
			fmt.Sprintf("Unexpected response code from API: %d", statusCode)+
				fmt.Sprintf("\n\nbody:\n\n%+v", string(body)),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r terraformResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func populateDefectdojoResource(ctx context.Context, diags *diag.Diagnostics, resourceData terraformResourceData, ddResource *defectdojoResource) {
	resourceVal := reflect.ValueOf(resourceData).Elem()
	resourceType := resourceVal.Type()
	ddVal := reflect.ValueOf(*ddResource).Elem()

	for i := 0; i < resourceVal.NumField(); i++ {
		fieldDescriptor := resourceType.Field(i)
		tag := fieldDescriptor.Tag
		ddFieldName := tag.Get("ddField")
		if ddFieldName != "" {
			fieldValue := resourceVal.Field(i)

			// Skip fields that are null or unknown - they should not overwrite
			// existing values (e.g., values read from the API before an update).
			isNull := fieldValue.MethodByName("IsNull").Call(nil)[0].Bool()
			isUnknown := fieldValue.MethodByName("IsUnknown").Call(nil)[0].Bool()
			if isNull || isUnknown {
				continue
			}

			ddFieldDescriptor, ok := ddVal.Type().FieldByName(ddFieldName)
			if !ok {
				diags.AddError("Error: No such field", fmt.Sprintf("A field named %s was specified to look sync data from the defectdojo client type, but no such field was found.", ddFieldName))
				continue
			}
			ddFieldValue := ddVal.FieldByName(ddFieldName)

			switch fieldDescriptor.Type {

			case typeOfTypesString:
				if ddFieldDescriptor.Type.Kind() == reflect.String {
					// if the destination field is a string (or named string type), we can grab the `Value` field and assign it directly
					ddFieldValue.Set(fieldValue.MethodByName("ValueString").Call(nil)[0].Convert(ddFieldDescriptor.Type))
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem().Kind() == reflect.String {
					// the destination field is a *string (or compatible/alias) so we have to set it to a pointer
					destType := ddFieldDescriptor.Type.Elem()
					destVal := reflect.New(destType)
					destVal.Elem().Set(fieldValue.MethodByName("ValueString").Call(nil)[0].Convert(destType))
					ddFieldValue.Set(destVal)
				} else if ddFieldDescriptor.Type.Kind() == reflect.Int {
					// the destination field is an int
					srcVal := fieldValue.MethodByName("ValueString").Call(nil)[0]
					strVal := srcVal.Interface().(string)
					intVal, err := strconv.Atoi(strVal)

					if err != nil {
						diags.AddError("Error converting value", fmt.Sprintf("Could not convert string value %s to *int: %e", strVal, err))
						continue
					}
					ddFieldValue.Set(reflect.ValueOf(intVal))
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem().Kind() == reflect.Int {
					// the destination field is a *int
					destType := ddFieldDescriptor.Type.Elem()
					destVal := reflect.New(destType)
					str := fieldValue.MethodByName("ValueString").Call(nil)[0].String()
					num, err := strconv.Atoi(str)
					if err != nil {
						diags.AddError("Error converting value", fmt.Sprintf("Could not convert string value %s to *int: %e", str, err))
						continue
					}
					destVal.Elem().Set(reflect.ValueOf(num))
					ddFieldValue.Set(destVal)
				} else if ddFieldDescriptor.Type == reflect.TypeFor[time.Time]() {
					str := fieldValue.MethodByName("ValueString").Call(nil)[0].String()
					t, err := parseDateTime(str)
					if err != nil {
						diags.AddError("Error converting value", fmt.Sprintf("Could not parse %s: %s", tag.Get("tfsdk"), err))
						continue
					}
					ddFieldValue.Set(reflect.ValueOf(t))
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem() == reflect.TypeFor[time.Time]() {
					str := fieldValue.MethodByName("ValueString").Call(nil)[0].String()
					t, err := parseDateTime(str)
					if err != nil {
						diags.AddError("Error converting value", fmt.Sprintf("Could not parse %s: %s", tag.Get("tfsdk"), err))
						continue
					}
					ddFieldValue.Set(reflect.ValueOf(&t))
				} else if ddFieldDescriptor.Type == reflect.TypeFor[openapi_types.Date]() {
					str := fieldValue.MethodByName("ValueString").Call(nil)[0].String()
					t, err := parseDate(str)
					if err != nil {
						diags.AddError("Error converting value", fmt.Sprintf("Could not parse %s: %s", tag.Get("tfsdk"), err))
						continue
					}
					ddFieldValue.Set(reflect.ValueOf(openapi_types.Date{Time: t}))
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem() == reflect.TypeFor[openapi_types.Date]() {
					str := fieldValue.MethodByName("ValueString").Call(nil)[0].String()
					t, err := parseDate(str)
					if err != nil {
						diags.AddError("Error converting value", fmt.Sprintf("Could not parse %s: %s", tag.Get("tfsdk"), err))
						continue
					}
					d := openapi_types.Date{Time: t}
					ddFieldValue.Set(reflect.ValueOf(&d))
				} else {
					addUnsupportedMappingError(diags, "populateDefectdojoResource", tag.Get("tfsdk"), fieldDescriptor.Type, ddFieldDescriptor.Type)
				}

			case typeOfTypesBool:
				if ddFieldDescriptor.Type.Kind() == reflect.Bool {
					ddFieldValue.Set(fieldValue.MethodByName("ValueBool").Call(nil)[0])
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem().Kind() == reflect.Bool {
					destType := ddFieldDescriptor.Type.Elem()
					destVal := reflect.New(destType)
					destVal.Elem().Set(fieldValue.MethodByName("ValueBool").Call(nil)[0].Convert(destType))
					ddFieldValue.Set(destVal)
				} else {
					addUnsupportedMappingError(diags, "populateDefectdojoResource", tag.Get("tfsdk"), fieldDescriptor.Type, ddFieldDescriptor.Type)
				}

			case typeOfTypesInt64:
				if ddFieldDescriptor.Type.Kind() == reflect.Int {
					destVal := reflect.New(ddFieldDescriptor.Type)
					destVal.Elem().Set(fieldValue.MethodByName("ValueInt64").Call(nil)[0].Convert(ddFieldDescriptor.Type))
					ddFieldValue.Set(destVal.Elem())
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem().Kind() == reflect.Int {
					destType := ddFieldDescriptor.Type.Elem()
					destVal := reflect.New(destType)
					destVal.Elem().Set(fieldValue.MethodByName("ValueInt64").Call(nil)[0].Convert(destType))
					ddFieldValue.Set(destVal)
				} else if ddFieldDescriptor.Type.Kind() == reflect.Int32 {
					v := int32(fieldValue.MethodByName("ValueInt64").Call(nil)[0].Int())
					ddFieldValue.Set(reflect.ValueOf(v))
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem().Kind() == reflect.Int32 {
					v := int32(fieldValue.MethodByName("ValueInt64").Call(nil)[0].Int())
					ddFieldValue.Set(reflect.ValueOf(&v))
				} else {
					addUnsupportedMappingError(diags, "populateDefectdojoResource", tag.Get("tfsdk"), fieldDescriptor.Type, ddFieldDescriptor.Type)
				}

			case typeOfTypesFloat64:
				if ddFieldDescriptor.Type.Kind() == reflect.Float64 {
					ddFieldValue.Set(fieldValue.MethodByName("ValueFloat64").Call(nil)[0])
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem().Kind() == reflect.Float64 {
					v := fieldValue.MethodByName("ValueFloat64").Call(nil)[0].Float()
					ddFieldValue.Set(reflect.ValueOf(&v))
				} else if ddFieldDescriptor.Type.Kind() == reflect.Float32 {
					v := float32(fieldValue.MethodByName("ValueFloat64").Call(nil)[0].Float())
					ddFieldValue.Set(reflect.ValueOf(v))
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem().Kind() == reflect.Float32 {
					v := float32(fieldValue.MethodByName("ValueFloat64").Call(nil)[0].Float())
					ddFieldValue.Set(reflect.ValueOf(&v))
				} else {
					addUnsupportedMappingError(diags, "populateDefectdojoResource", tag.Get("tfsdk"), fieldDescriptor.Type, ddFieldDescriptor.Type)
				}

			case typeOfTypesSet:
				if ddFieldDescriptor.Type.Kind() == reflect.Slice {
					// the destination field is a direct slice (e.g. []int, []string,
					// or a slice of a defined int/string type); elements are
					// converted one by one so defined element types work too
					if ddFieldDescriptor.Type.Elem().Kind() == reflect.Int {
						int64s := []int64{}
						diags_ := fieldValue.Interface().(types.Set).ElementsAs(ctx, &int64s, false)
						if len(diags_) > 0 {
							diags.Append(diags_...)
							continue
						}
						out := reflect.MakeSlice(ddFieldDescriptor.Type, len(int64s), len(int64s))
						for i, val := range int64s {
							out.Index(i).Set(reflect.ValueOf((int)(val)).Convert(ddFieldDescriptor.Type.Elem()))
						}
						ddFieldValue.Set(out)
					} else if ddFieldDescriptor.Type.Elem().Kind() == reflect.String {
						strs := []string{}
						diags_ := fieldValue.Interface().(types.Set).ElementsAs(ctx, &strs, false)
						if len(diags_) > 0 {
							diags.Append(diags_...)
							continue
						}
						out := reflect.MakeSlice(ddFieldDescriptor.Type, len(strs), len(strs))
						for i, val := range strs {
							out.Index(i).Set(reflect.ValueOf(val).Convert(ddFieldDescriptor.Type.Elem()))
						}
						ddFieldValue.Set(out)
					} else {
						// A slice whose element kind is neither int nor string is not
						// convertible by the loops above. Without this branch the value
						// was dropped with no diagnostic and no log line at all.
						addUnsupportedMappingError(diags, "populateDefectdojoResource", tag.Get("tfsdk"), fieldDescriptor.Type, ddFieldDescriptor.Type)
					}
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem().Kind() == reflect.Slice {
					// the destination field is a pointer to a slice; elements are
					// converted one by one so defined element types (e.g. enum
					// strings like dd.NotificationsRequestScanAdded) work too
					if ddFieldDescriptor.Type.Elem().Elem().Kind() == reflect.Int {
						// it's a slice of int
						int64s := []int64{}
						diags_ := fieldValue.Interface().(types.Set).ElementsAs(ctx, &int64s, false)
						if len(diags_) > 0 {
							diags.Append(diags_...)
							continue
						}
						sliceType := ddFieldDescriptor.Type.Elem()
						out := reflect.MakeSlice(sliceType, len(int64s), len(int64s))
						for i, val := range int64s {
							out.Index(i).Set(reflect.ValueOf((int)(val)).Convert(sliceType.Elem()))
						}
						destVal := reflect.New(sliceType)
						destVal.Elem().Set(out)
						ddFieldValue.Set(destVal)
					} else if ddFieldDescriptor.Type.Elem().Elem().Kind() == reflect.String {
						// it's a slice of string
						strs := []string{}
						diags_ := fieldValue.Interface().(types.Set).ElementsAs(ctx, &strs, false)
						if len(diags_) > 0 {
							diags.Append(diags_...)
							continue
						}
						sliceType := ddFieldDescriptor.Type.Elem()
						out := reflect.MakeSlice(sliceType, len(strs), len(strs))
						for i, val := range strs {
							out.Index(i).Set(reflect.ValueOf(val).Convert(sliceType.Elem()))
						}
						destVal := reflect.New(sliceType)
						destVal.Elem().Set(out)
						ddFieldValue.Set(destVal)
					} else {
						// A slice whose element kind is neither int nor string is not
						// convertible by the loops above. Without this branch the value
						// was dropped with no diagnostic and no log line at all.
						addUnsupportedMappingError(diags, "populateDefectdojoResource", tag.Get("tfsdk"), fieldDescriptor.Type, ddFieldDescriptor.Type)
					}
				} else {
					addUnsupportedMappingError(diags, "populateDefectdojoResource", tag.Get("tfsdk"), fieldDescriptor.Type, ddFieldDescriptor.Type)
				}

			default:
				addUnhandledTerraformTypeError(diags, "populateDefectdojoResource", tag.Get("tfsdk"), fieldDescriptor.Type, ddFieldDescriptor.Type)
			}
		}
	}
}

func populateResourceData(ctx context.Context, diags *diag.Diagnostics, d *terraformResourceData, ddResource defectdojoResource) {
	tflog.Info(ctx, "populateResourceData")

	resourceVal := reflect.ValueOf(*d).Elem()
	resourceType := resourceVal.Type()
	ddVal := reflect.ValueOf(ddResource).Elem()

	for i := 0; i < resourceVal.NumField(); i++ {
		fieldDescriptor := resourceType.Field(i)
		tag := fieldDescriptor.Tag
		ddFieldName := tag.Get("ddField")
		if ddFieldName != "" {
			fieldValue := resourceVal.Field(i)

			ddFieldDescriptor, ok := ddVal.Type().FieldByName(ddFieldName)
			if !ok {
				diags.AddError("Error: No such field", fmt.Sprintf("A field named %s was specified to look sync data from the defectdojo client type, but no such field was found.", ddFieldName))
				continue
			}
			ddFieldValue := ddVal.FieldByName(ddFieldName)

			switch fieldDescriptor.Type {

			case typeOfTypesString:
				if ddFieldDescriptor.Type.Kind() == reflect.String {
					// if the source field is a string, we can use it directly,
					// subject to any `ddFormat` normalisation the attribute asks for
					current := fieldValue.Interface().(types.String)
					fieldValue.Set(reflect.ValueOf(renderStringValue(diags, tag, current, ddFieldValue.String())))
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem().Kind() == reflect.String {
					// if the source field is a pointer, make sure it's a pointer to a string, and then we can grab the pointed-to value,
					// but only if the pointer is not nil
					if !ddFieldValue.IsNil() {
						current := fieldValue.Interface().(types.String)
						fieldValue.Set(reflect.ValueOf(renderStringValue(diags, tag, current, ddFieldValue.Elem().String())))
					} else {
						fieldValue.Set(reflect.ValueOf(types.StringNull()))
					}
				} else if ddFieldDescriptor.Type.Kind() == reflect.Int {
					fieldValue.Set(reflect.ValueOf(types.StringValue(fmt.Sprint(ddFieldValue.Int()))))
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem().Kind() == reflect.Int {
					if !ddFieldValue.IsNil() {
						fieldValue.Set(reflect.ValueOf(types.StringValue(fmt.Sprint(ddFieldValue.Elem().Int()))))
					} else {
						fieldValue.Set(reflect.ValueOf(types.StringNull()))
					}
				} else if ddFieldDescriptor.Type == reflect.TypeFor[time.Time]() {
					t := ddFieldValue.Interface().(time.Time)
					if !t.IsZero() {
						// Keep the configured/prior literal when it denotes the
						// same instant, so state matches config (see #23).
						current := fieldValue.Interface().(types.String)
						fieldValue.Set(reflect.ValueOf(preserveDateTimeLiteral(current, t)))
					} else {
						fieldValue.Set(reflect.ValueOf(types.StringNull()))
					}
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem() == reflect.TypeFor[time.Time]() {
					if !ddFieldValue.IsNil() {
						t := ddFieldValue.Elem().Interface().(time.Time)
						current := fieldValue.Interface().(types.String)
						fieldValue.Set(reflect.ValueOf(preserveDateTimeLiteral(current, t)))
					} else {
						fieldValue.Set(reflect.ValueOf(types.StringNull()))
					}
				} else if ddFieldDescriptor.Type == reflect.TypeFor[openapi_types.Date]() {
					d := ddFieldValue.Interface().(openapi_types.Date)
					if !d.IsZero() {
						fieldValue.Set(reflect.ValueOf(types.StringValue(d.Format("2006-01-02"))))
					} else {
						fieldValue.Set(reflect.ValueOf(types.StringNull()))
					}
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem() == reflect.TypeFor[openapi_types.Date]() {
					if !ddFieldValue.IsNil() {
						d := ddFieldValue.Elem().Interface().(openapi_types.Date)
						fieldValue.Set(reflect.ValueOf(types.StringValue(d.Format("2006-01-02"))))
					} else {
						fieldValue.Set(reflect.ValueOf(types.StringNull()))
					}
				} else {
					addUnsupportedMappingError(diags, "populateResourceData", tag.Get("tfsdk"), fieldDescriptor.Type, ddFieldDescriptor.Type)
				}

			case typeOfTypesBool:
				if ddFieldDescriptor.Type.Kind() == reflect.Bool {
					// if the source field is a bool, we can use it directly
					fieldValue.Set(reflect.ValueOf(types.BoolValue(ddFieldValue.Bool())))
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem().Kind() == reflect.Bool {
					// if the source field is a pointer, make sure it's a pointer to a bool, and then we can grab the pointed-to value,
					// but only if the pointer is not nil
					if !ddFieldValue.IsNil() {
						fieldValue.Set(reflect.ValueOf(types.BoolValue(ddFieldValue.Elem().Bool())))
					} else {
						fieldValue.Set(reflect.ValueOf(types.BoolNull()))
					}
				} else {
					addUnsupportedMappingError(diags, "populateResourceData", tag.Get("tfsdk"), fieldDescriptor.Type, ddFieldDescriptor.Type)
				}

			case typeOfTypesInt64:
				// int64 / *int64 are accepted here but NOT by
				// populateDefectdojoResource, and ddFieldPairingSupported
				// (ddfield_audit_unit_test.go) rejects the pairing for exactly that
				// reason: the ddField contract is bidirectional.
				//
				// Do not "fix" the audit to match this branch. An int64 field would
				// then round-trip cleanly on Read and be silently dropped on
				// Create/Update - the issue #23 shape precisely. If the generated
				// client ever grows an int64 field worth mapping, add the write-path
				// branch first, then relax the audit.
				if ddFieldDescriptor.Type.Kind() == reflect.Int64 || ddFieldDescriptor.Type.Kind() == reflect.Int {
					// if the source field is an int or int64, we can cast and use it directly
					fieldValue.Set(reflect.ValueOf(types.Int64Value((int64)(ddFieldValue.Int()))))
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && (ddFieldDescriptor.Type.Elem().Kind() == reflect.Int64 || ddFieldDescriptor.Type.Elem().Kind() == reflect.Int) {
					// if the source field is a pointer, make sure it's a pointer to an int64, and then we can grab the pointed-to value,
					// but only if the pointer is not nil
					if !ddFieldValue.IsNil() {
						fieldValue.Set(reflect.ValueOf(types.Int64Value((int64)(ddFieldValue.Elem().Int()))))
					} else {
						fieldValue.Set(reflect.ValueOf(types.Int64Null()))
					}
				} else if ddFieldDescriptor.Type.Kind() == reflect.Int32 {
					fieldValue.Set(reflect.ValueOf(types.Int64Value(int64(ddFieldValue.Int()))))
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem().Kind() == reflect.Int32 {
					if !ddFieldValue.IsNil() {
						fieldValue.Set(reflect.ValueOf(types.Int64Value(int64(ddFieldValue.Elem().Int()))))
					} else {
						fieldValue.Set(reflect.ValueOf(types.Int64Null()))
					}
				} else {
					addUnsupportedMappingError(diags, "populateResourceData", tag.Get("tfsdk"), fieldDescriptor.Type, ddFieldDescriptor.Type)
				}

			case typeOfTypesFloat64:
				if ddFieldDescriptor.Type.Kind() == reflect.Float64 {
					fieldValue.Set(reflect.ValueOf(types.Float64Value(ddFieldValue.Float())))
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem().Kind() == reflect.Float64 {
					if !ddFieldValue.IsNil() {
						fieldValue.Set(reflect.ValueOf(types.Float64Value(ddFieldValue.Elem().Float())))
					} else {
						fieldValue.Set(reflect.ValueOf(types.Float64Null()))
					}
				} else if ddFieldDescriptor.Type.Kind() == reflect.Float32 {
					fieldValue.Set(reflect.ValueOf(types.Float64Value(float64(ddFieldValue.Interface().(float32)))))
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem().Kind() == reflect.Float32 {
					if !ddFieldValue.IsNil() {
						fieldValue.Set(reflect.ValueOf(types.Float64Value(float64(ddFieldValue.Elem().Interface().(float32)))))
					} else {
						fieldValue.Set(reflect.ValueOf(types.Float64Null()))
					}
				} else {
					addUnsupportedMappingError(diags, "populateResourceData", tag.Get("tfsdk"), fieldDescriptor.Type, ddFieldDescriptor.Type)
				}

			case typeOfTypesSet:
				if ddFieldDescriptor.Type.Kind() == reflect.Slice {
					// the source field is a direct slice (e.g. []int, []string)
					if ddFieldDescriptor.Type.Elem().Kind() == reflect.Int {
						if ddFieldValue.Len() > 0 || !fieldValue.MethodByName("IsNull").Call(nil)[0].Bool() {
							elems := []attr.Value{}
							for i := 0; i < ddFieldValue.Len(); i++ {
								elems = append(elems, types.Int64Value(ddFieldValue.Index(i).Int()))
							}
							destVal, dgs := types.SetValue(types.Int64Type, elems)
							diags.Append(dgs...)
							fieldValue.Set(reflect.ValueOf(destVal))
						} else {
							fieldValue.Set(reflect.ValueOf(types.SetNull(types.Int64Type)))
						}
					} else if ddFieldDescriptor.Type.Elem().Kind() == reflect.String {
						if ddFieldValue.Len() > 0 || !fieldValue.MethodByName("IsNull").Call(nil)[0].Bool() {
							elems := []attr.Value{}
							for i := 0; i < ddFieldValue.Len(); i++ {
								elems = append(elems, types.StringValue(ddFieldValue.Index(i).String()))
							}
							destVal, dgs := types.SetValue(types.StringType, elems)
							diags.Append(dgs...)
							fieldValue.Set(reflect.ValueOf(destVal))
						} else {
							fieldValue.Set(reflect.ValueOf(types.SetNull(types.StringType)))
						}
					} else {
						// A slice whose element kind is neither int nor string is not
						// convertible by the loops above. Without this branch the value
						// was dropped with no diagnostic and no log line at all.
						addUnsupportedMappingError(diags, "populateResourceData", tag.Get("tfsdk"), fieldDescriptor.Type, ddFieldDescriptor.Type)
					}
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem().Kind() == reflect.Slice {
					// the source field is a pointer to a slice
					if ddFieldDescriptor.Type.Elem().Elem().Kind() == reflect.Int {
						// it's a slice of int

						if !ddFieldValue.IsZero() && (ddFieldValue.Elem().Len() > 0 || !fieldValue.MethodByName("IsNull").Call(nil)[0].Bool()) {
							elems := []attr.Value{}
							for i := 0; i < ddFieldValue.Elem().Len(); i++ {
								elems = append(elems, types.Int64Value(ddFieldValue.Elem().Index(i).Int()))
							}
							destVal, dgs := types.SetValue(types.Int64Type, elems)
							diags.Append(dgs...)
							fieldValue.Set(reflect.ValueOf(destVal))
						} else {
							destVal := types.SetNull(types.Int64Type)
							fieldValue.Set(reflect.ValueOf(destVal))
						}
					} else if ddFieldDescriptor.Type.Elem().Elem().Kind() == reflect.String {
						// it's a slice of string

						if !ddFieldValue.IsZero() && (ddFieldValue.Elem().Len() > 0 || !fieldValue.MethodByName("IsNull").Call(nil)[0].Bool()) {
							elems := []attr.Value{}
							for i := 0; i < ddFieldValue.Elem().Len(); i++ {
								elems = append(elems, types.StringValue(ddFieldValue.Elem().Index(i).String()))
							}
							destVal, dgs := types.SetValue(types.StringType, elems)
							diags.Append(dgs...)
							fieldValue.Set(reflect.ValueOf(destVal))
						} else {
							destVal := types.SetNull(types.StringType)
							fieldValue.Set(reflect.ValueOf(destVal))
						}
					} else {
						// A slice whose element kind is neither int nor string is not
						// convertible by the loops above. Without this branch the value
						// was dropped with no diagnostic and no log line at all.
						addUnsupportedMappingError(diags, "populateResourceData", tag.Get("tfsdk"), fieldDescriptor.Type, ddFieldDescriptor.Type)
					}
				} else {
					addUnsupportedMappingError(diags, "populateResourceData", tag.Get("tfsdk"), fieldDescriptor.Type, ddFieldDescriptor.Type)
				}
			default:
				addUnhandledTerraformTypeError(diags, "populateResourceData", tag.Get("tfsdk"), fieldDescriptor.Type, ddFieldDescriptor.Type)
			}
		}
	}
}
