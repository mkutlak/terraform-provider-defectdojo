package provider

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
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
	populateDefectdojoResource(ctx, &diags, data, &ddResource)

	statusCode, body, err := ddResource.createApiCall(ctx, r.client)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating "+r.typeName,
			err.Error())
		return
	}

	if statusCode == 201 {
		populateResourceData(ctx, &diags, &data, ddResource)
	} else {
		resp.Diagnostics.AddError(
			"API Error Creating "+r.typeName,
			fmt.Sprintf("Unexpected response code from API: %d", statusCode)+
				fmt.Sprintf("\n\nbody:\n\n%s", string(body)),
		)
		return
	}

	tflog.Trace(ctx, "resource created")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
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
	populateDefectdojoResource(ctx, &diags, data, &ddResource)

	statusCode, body, err := ddResource.readApiCall(ctx, r.client, idNumber)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving "+r.typeName,
			err.Error())
		return
	}

	switch statusCode {
	case 200:
		populateResourceData(ctx, &diags, &data, ddResource)
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

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
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
	populateDefectdojoResource(ctx, &diags, data, &ddResource)

	statusCode, body, err := ddResource.updateApiCall(ctx, r.client, idNumber)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating "+r.typeName,
			err.Error())
		return
	}

	if statusCode == 200 {
		populateResourceData(ctx, &diags, &data, ddResource)
	} else {
		resp.Diagnostics.AddError(
			"API Error Updating "+r.typeName,
			fmt.Sprintf("Unexpected response code from API: %d", statusCode)+
				fmt.Sprintf("\n\nbody:\n\n%+v", string(body)),
		)
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
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
					t, err := time.Parse(time.RFC3339, str)
					if err != nil {
						diags.AddError("Error converting value", fmt.Sprintf("Could not parse datetime value %s: %s", str, err))
						continue
					}
					ddFieldValue.Set(reflect.ValueOf(t))
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem() == reflect.TypeFor[time.Time]() {
					str := fieldValue.MethodByName("ValueString").Call(nil)[0].String()
					t, err := time.Parse(time.RFC3339, str)
					if err != nil {
						diags.AddError("Error converting value", fmt.Sprintf("Could not parse datetime value %s: %s", str, err))
						continue
					}
					ddFieldValue.Set(reflect.ValueOf(&t))
				} else if ddFieldDescriptor.Type == reflect.TypeFor[openapi_types.Date]() {
					str := fieldValue.MethodByName("ValueString").Call(nil)[0].String()
					t, err := time.Parse("2006-01-02", str)
					if err != nil {
						diags.AddError("Error converting value", fmt.Sprintf("Could not parse date value %s: %s", str, err))
						continue
					}
					ddFieldValue.Set(reflect.ValueOf(openapi_types.Date{Time: t}))
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem() == reflect.TypeFor[openapi_types.Date]() {
					str := fieldValue.MethodByName("ValueString").Call(nil)[0].String()
					t, err := time.Parse("2006-01-02", str)
					if err != nil {
						diags.AddError("Error converting value", fmt.Sprintf("Could not parse date value %s: %s", str, err))
						continue
					}
					d := openapi_types.Date{Time: t}
					ddFieldValue.Set(reflect.ValueOf(&d))
				} else {
					tflog.Warn(ctx, fmt.Sprintf("WARN [populateDefectdojoResource]: Don't know how to assign type %s to type %s\n", fieldDescriptor.Type, ddFieldDescriptor.Type))
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
					tflog.Warn(ctx, fmt.Sprintf("WARN [populateDefectdojoResource]: Don't know how to assign type %s to type %s\n", fieldDescriptor.Type, ddFieldDescriptor.Type))
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
					tflog.Warn(ctx, fmt.Sprintf("WARN [populateDefectdojoResource]: Don't know how to assign type %s to type %s\n", fieldDescriptor.Type, ddFieldDescriptor.Type))
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
					tflog.Warn(ctx, fmt.Sprintf("WARN [populateDefectdojoResource]: Don't know how to assign type %s to type %s\n", fieldDescriptor.Type, ddFieldDescriptor.Type))
				}

			case typeOfTypesSet:
				if ddFieldDescriptor.Type.Kind() == reflect.Slice {
					// the destination field is a direct slice (e.g. []int, []string)
					if ddFieldDescriptor.Type.Elem().Kind() == reflect.Int {
						int64s := []int64{}
						diags_ := fieldValue.Interface().(types.Set).ElementsAs(context.Background(), &int64s, false)
						if len(diags_) > 0 {
							diags.Append(diags_...)
							continue
						}
						ints := make([]int, 0, len(int64s))
						for _, val := range int64s {
							ints = append(ints, (int)(val))
						}
						ddFieldValue.Set(reflect.ValueOf(ints))
					} else if ddFieldDescriptor.Type.Elem().Kind() == reflect.String {
						strs := []string{}
						diags_ := fieldValue.Interface().(types.Set).ElementsAs(context.Background(), &strs, false)
						if len(diags_) > 0 {
							diags.Append(diags_...)
							continue
						}
						ddFieldValue.Set(reflect.ValueOf(strs))
					}
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem().Kind() == reflect.Slice {
					// the destination field is a pointer to a slice
					if ddFieldDescriptor.Type.Elem().Elem().Kind() == reflect.Int {
						// it's a slice of int
						int64s := []int64{}
						diags_ := fieldValue.Interface().(types.Set).ElementsAs(context.Background(), &int64s, false)
						if len(diags_) > 0 {
							diags.Append(diags_...)
							continue
						}
						ints := make([]int, 0, len(int64s))
						for _, val := range int64s {
							ints = append(ints, (int)(val))
						}
						destVal := reflect.New(ddFieldDescriptor.Type.Elem())
						destVal.Elem().Set(reflect.ValueOf(ints))
						ddFieldValue.Set(destVal)
					} else if ddFieldDescriptor.Type.Elem().Elem().Kind() == reflect.String {
						// it's a slice of string
						strs := []string{}
						diags_ := fieldValue.Interface().(types.Set).ElementsAs(context.Background(), &strs, false)
						if len(diags_) > 0 {
							diags.Append(diags_...)
							continue
						}
						destVal := reflect.New(ddFieldDescriptor.Type.Elem())
						destVal.Elem().Set(reflect.ValueOf(strs))
						ddFieldValue.Set(destVal)
					}
				} else {
					tflog.Warn(ctx, fmt.Sprintf("WARN [populateDefectdojoResource]: Don't know how to assign type %s to type %s\n", fieldDescriptor.Type, ddFieldDescriptor.Type))
				}

			default:
				tflog.Warn(ctx, fmt.Sprintf("WARN [populateDefectdojoResource]: Don't know how to assign anything (type was %s) to type %s\n", fieldDescriptor.Type, ddFieldDescriptor.Type))
			}
		}
	}
}

func populateResourceData(ctx context.Context, diags *diag.Diagnostics, d *terraformResourceData, ddResource defectdojoResource) {
	tflog.Info(context.Background(), "populateResourceData")

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
					// if the source field is a string, we can use it directly
					fieldValue.Set(reflect.ValueOf(types.StringValue(ddFieldValue.String())))
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem().Kind() == reflect.String {
					// if the source field is a pointer, make sure it's a pointer to a string, and then we can grab the pointed-to value,
					// but only if the pointer is not nil
					if !ddFieldValue.IsNil() {
						fieldValue.Set(reflect.ValueOf(types.StringValue(ddFieldValue.Elem().String())))
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
						fieldValue.Set(reflect.ValueOf(types.StringValue(t.Format(time.RFC3339))))
					} else {
						fieldValue.Set(reflect.ValueOf(types.StringNull()))
					}
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem() == reflect.TypeFor[time.Time]() {
					if !ddFieldValue.IsNil() {
						t := ddFieldValue.Elem().Interface().(time.Time)
						fieldValue.Set(reflect.ValueOf(types.StringValue(t.Format(time.RFC3339))))
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
					tflog.Warn(ctx, fmt.Sprintf("WARN [populateResourceData]: Don't know how to assign type %s to type %s\n", ddFieldDescriptor.Type, fieldDescriptor.Type))
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
					tflog.Warn(ctx, fmt.Sprintf("WARN [populateResourceData]: Don't know how to assign type %s to type %s\n", ddFieldDescriptor.Type, fieldDescriptor.Type))
				}

			case typeOfTypesInt64:
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
					tflog.Warn(ctx, fmt.Sprintf("WARN [populateResourceData]: Don't know how to assign type %s to type %s\n", ddFieldDescriptor.Type, fieldDescriptor.Type))
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
					tflog.Warn(ctx, fmt.Sprintf("WARN [populateResourceData]: Don't know how to assign type %s to type %s\n", ddFieldDescriptor.Type, fieldDescriptor.Type))
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
							diags.Append(dgs.Errors()...)
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
							diags.Append(dgs.Errors()...)
							fieldValue.Set(reflect.ValueOf(destVal))
						} else {
							fieldValue.Set(reflect.ValueOf(types.SetNull(types.StringType)))
						}
					}
				} else if ddFieldDescriptor.Type.Kind() == reflect.Ptr && ddFieldDescriptor.Type.Elem().Kind() == reflect.Slice {
					// the source field is a pointer to a slice
					if ddFieldDescriptor.Type.Elem().Elem().Kind() == reflect.Int {
						// it's a slice of int

						if !ddFieldValue.IsZero() && (ddFieldValue.Elem().Len() > 0 || !fieldValue.MethodByName("IsNull").Call(nil)[0].Bool()) {
							elems := []attr.Value{}
							for _, val := range ddFieldValue.Elem().Interface().([]int) {
								elems = append(elems, types.Int64Value((int64)(val)))
							}
							destVal, dgs := types.SetValue(types.Int64Type, elems)
							diags.Append(dgs.Errors()...)
							fieldValue.Set(reflect.ValueOf(destVal))
						} else {
							destVal := types.SetNull(types.Int64Type)
							fieldValue.Set(reflect.ValueOf(destVal))
						}
					} else if ddFieldDescriptor.Type.Elem().Elem().Kind() == reflect.String {
						// it's a slice of string

						if !ddFieldValue.IsZero() && (ddFieldValue.Elem().Len() > 0 || !fieldValue.MethodByName("IsNull").Call(nil)[0].Bool()) {
							elems := []attr.Value{}
							for _, val := range ddFieldValue.Elem().Interface().([]string) {
								elems = append(elems, types.StringValue((string)(val)))
							}
							destVal, dgs := types.SetValue(types.StringType, elems)
							diags.Append(dgs.Errors()...)
							fieldValue.Set(reflect.ValueOf(destVal))
						} else {
							destVal := types.SetNull(types.StringType)
							fieldValue.Set(reflect.ValueOf(destVal))
						}
					}
				} else {
					tflog.Warn(ctx, fmt.Sprintf("WARN [populateResourceData]: Don't know how to assign type %s to type %s\n", ddFieldDescriptor.Type, fieldDescriptor.Type))
				}
			default:
				tflog.Warn(ctx, fmt.Sprintf("WARN [populateResourceData]: Don't know how to assign anything (type was %s) to type %s\n", ddFieldDescriptor.Type, fieldDescriptor.Type))
			}
		}
	}
}
