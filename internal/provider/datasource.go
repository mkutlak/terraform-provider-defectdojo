package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dd "github.com/mkutlak/terraform-provider-defectdojo/internal/ddclient"
)

// nameFilterable is an optional interface that dataProviders can implement
// to enable lookup by name (or another unique field) instead of numeric ID.
type nameFilterable interface {
	nameFromData(data terraformResourceData) (string, bool)
	listByName(ctx context.Context, client *dd.ClientWithResponses, name string, data terraformResourceData) error
}

type terraformDatasource struct {
	client *dd.ClientWithResponses
	dataProvider
}

func (r *terraformDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r terraformDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data, diags := r.getData(ctx, req.Config)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Check for conflicting lookup parameters (both id and name specified)
	if !data.id().IsNull() {
		if nf, ok := r.dataProvider.(nameFilterable); ok {
			if _, hasName := nf.nameFromData(data); hasName {
				resp.Diagnostics.AddError(
					"Conflicting Lookup Parameters",
					"Specify either id or the lookup field (e.g. name, username, url), not both")
				return
			}
		}
	}

	// If id is null, try name-based resolution
	if data.id().IsNull() {
		nf, ok := r.dataProvider.(nameFilterable)
		if !ok {
			resp.Diagnostics.AddError(
				"Could not Retrieve Resource",
				"The Id field was null but it is required to retrieve the resource")
			return
		}

		name, hasName := nf.nameFromData(data)
		if !hasName {
			resp.Diagnostics.AddError(
				"Could not Retrieve Resource",
				"Either id or the lookup field (e.g. name, username, url) must be specified")
			return
		}

		if err := nf.listByName(ctx, r.client, name, data); err != nil {
			resp.Diagnostics.AddError(
				"Error Looking Up Resource By Name",
				err.Error())
			return
		}
	}

	idNumber, err := strconv.Atoi(data.id().ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not Retrieve Resource",
			"Error while parsing the resource ID from state: "+err.Error())
		return
	}

	ddResource := data.defectdojoResource()
	populateDefectdojoResource(ctx, &diags, data, &ddResource)

	statusCode, body, err := ddResource.readApiCall(ctx, r.client, idNumber)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Resource",
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
			"API Error Retrieving Resource",
			fmt.Sprintf("Unexpected response code from API: %d", statusCode)+
				fmt.Sprintf("\n\nbody:\n\n%+v", string(body)),
		)
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
