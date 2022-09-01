package resource

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-provider-mock/internal/computed"
	"os"

	"github.com/hashicorp/terraform-provider-mock/internal/client"
	"github.com/hashicorp/terraform-provider-mock/internal/data"
	"github.com/hashicorp/terraform-provider-mock/internal/schema"
)

var _ resource.Resource = Resource{}
var _ resource.ResourceWithImportState = Resource{}

type Resource struct {
	Schema schema.Schema
	Client client.Client
}

func (r Resource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	resource := &data.Resource{}
	response.Diagnostics.Append(request.Config.Get(ctx, &resource)...)
	if response.Diagnostics.HasError() {
		return
	}

	if err := computed.GenerateComputedValues(resource, r.Schema); err != nil {
		response.Diagnostics.Append(diag.NewErrorDiagnostic("failed to generate computed values", err.Error()))
		return
	}

	if err := r.Client.WriteResource(ctx, resource); err != nil {
		response.Diagnostics.Append(diag.NewErrorDiagnostic("failed to write resource", err.Error()))
		return
	}

	diags := response.State.Set(ctx, resource)
	response.Diagnostics.Append(diags...)
}

func (r Resource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	resource := &data.Resource{}
	response.Diagnostics.Append(request.State.Get(ctx, &resource)...)
	if response.Diagnostics.HasError() {
		return
	}

	data, err := r.Client.ReadResource(ctx, resource.GetId())
	if err != nil {
		if os.IsNotExist(err) {
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.Append(diag.NewErrorDiagnostic("failed to read resource", err.Error()))
		return
	}

	if data == nil {
		// The client returned a nil object with no error. This means it is
		// telling us to just rely on the state.
		data = resource
	}

	typ := request.State.Schema.Type().TerraformType(ctx)
	response.Diagnostics.Append(response.State.Set(ctx, data.WithType(typ.(tftypes.Object)))...)
}

func (r Resource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	resource := &data.Resource{}
	response.Diagnostics.Append(request.Plan.Get(ctx, &resource)...)
	if response.Diagnostics.HasError() {
		return
	}

	if err := r.Client.UpdateResource(ctx, resource); err != nil {
		response.Diagnostics.Append(diag.NewErrorDiagnostic("failed to update resource", err.Error()))
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, resource)...)
}

func (r Resource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	resource := &data.Resource{}
	response.Diagnostics.Append(request.State.Get(ctx, &resource)...)
	if response.Diagnostics.HasError() {
		return
	}

	if err := r.Client.DeleteResource(ctx, resource.GetId()); err != nil {
		response.Diagnostics.Append(diag.NewErrorDiagnostic("failed to delete resource", err.Error()))
	}
}

func (r Resource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}
