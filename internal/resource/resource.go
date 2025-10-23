// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"fmt"
	"os"
	"slices"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-provider-tfcoremock/internal/computed"

	"github.com/hashicorp/terraform-provider-tfcoremock/internal/client"
	"github.com/hashicorp/terraform-provider-tfcoremock/internal/data"
	"github.com/hashicorp/terraform-provider-tfcoremock/internal/schema"
)

var _ resource.Resource = Resource{}
var _ resource.ResourceWithIdentity = Resource{}
var _ resource.ResourceWithImportState = Resource{}
var _ resource.ResourceWithModifyPlan = Resource{}

type Resource struct {
	Name           string
	InternalSchema schema.Schema
	Client         client.Client

	FailOnDelete []string
	FailOnCreate []string
	FailOnRead   []string
	FailOnUpdate []string
	DeferChanges []string
}

func (r Resource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = r.Name
}

func (r Resource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	var err error
	if response.Schema, err = r.InternalSchema.ToTerraformResourceSchema(); err != nil {
		response.Diagnostics.Append(diag.NewErrorDiagnostic(fmt.Sprintf("failed to build resource schema for '%s'", r.Name), err.Error()))
	}
}

func (r Resource) IdentitySchema(ctx context.Context, request resource.IdentitySchemaRequest, response *resource.IdentitySchemaResponse) {
	response.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
				Description:       "The ID of the resource.",
			},
		},
	}
}

func (r Resource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	resource := &data.Resource{}
	response.Diagnostics.Append(request.Plan.Get(ctx, &resource)...)
	if response.Diagnostics.HasError() {
		return
	}
	resource.ResourceType = r.Name

	// The root ID is a special computed value.
	if _, ok := resource.Values["id"]; !ok {
		id, err := uuid.GenerateUUID()
		if err != nil {
			response.Diagnostics.Append(diag.NewErrorDiagnostic("failed to generate id", err.Error()))
			return
		}
		resource.Values["id"] = data.Value{
			String: &id,
		}
	}

	// Now go and do the rest of the computed values.
	if err := computed.GenerateComputedValues(resource, r.InternalSchema); err != nil {
		response.Diagnostics.Append(diag.NewErrorDiagnostic("failed to generate computed values", err.Error()))
		return
	}

	if slices.Contains(r.FailOnCreate, resource.GetId()) {
		response.Diagnostics.Append(diag.NewErrorDiagnostic("failed to create resource", "forced failure"))
		return
	}

	if err := r.Client.WriteResource(ctx, resource); err != nil {
		response.Diagnostics.Append(diag.NewErrorDiagnostic("failed to write resource", err.Error()))
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, resource)...)
	response.Diagnostics.Append(response.Identity.Set(ctx, resource.Identity())...)
}

func (r Resource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	resource := &data.Resource{}
	response.Diagnostics.Append(request.State.Get(ctx, &resource)...)
	if response.Diagnostics.HasError() {
		return
	}

	if slices.Contains(r.FailOnRead, resource.GetId()) {
		response.Diagnostics.Append(diag.NewErrorDiagnostic("failed to read resource", "forced failure"))
		return
	}

	data, err := r.Client.ReadResource(ctx, resource.GetId())
	if err != nil {
		if os.IsNotExist(err) {
			// This is a bit of weird one as it means we tried to read a file
			// that doesn't exist but Terraform thinks it does. We treat this
			// as "drift" and let the Terraform framework handle it.
			response.State.RemoveResource(ctx)
			response.Diagnostics.Append(response.Identity.Set(ctx, resource.Identity())...)
			return
		}
		response.Diagnostics.AddError("failed to read resource", err.Error())
		return
	}

	if data == nil {
		// The client returned a nil object with no error. This means it is
		// telling us to just rely on the state.
		data = resource
	}

	typ := request.State.Schema.Type().TerraformType(ctx)
	response.Diagnostics.Append(response.State.Set(ctx, data.WithType(typ.(tftypes.Object)))...)
	response.Diagnostics.Append(response.Identity.Set(ctx, data.Identity())...)
}

func (r Resource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	resource := &data.Resource{}
	response.Diagnostics.Append(request.Plan.Get(ctx, &resource)...)
	if response.Diagnostics.HasError() {
		return
	}
	resource.ResourceType = r.Name

	if err := computed.GenerateComputedValues(resource, r.InternalSchema); err != nil {
		response.Diagnostics.Append(diag.NewErrorDiagnostic("failed to generate computed values", err.Error()))
		return
	}

	if slices.Contains(r.FailOnUpdate, resource.GetId()) {
		response.Diagnostics.Append(diag.NewErrorDiagnostic("failed to update resource", "forced failure"))
		return
	}

	if err := r.Client.UpdateResource(ctx, resource); err != nil {
		response.Diagnostics.AddError("failed to update resource", err.Error())
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, resource)...)
	response.Diagnostics.Append(response.Identity.Set(ctx, resource.Identity())...)
}

func (r Resource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	resource := &data.Resource{}
	response.Diagnostics.Append(request.State.Get(ctx, &resource)...)
	if response.Diagnostics.HasError() {
		return
	}

	if slices.Contains(r.FailOnDelete, resource.GetId()) {
		response.Diagnostics.Append(diag.NewErrorDiagnostic("failed to delete resource", "forced failure"))
		return
	}

	if err := r.Client.DeleteResource(ctx, resource.GetId()); err != nil {
		response.Diagnostics.AddError("failed to delete resource", err.Error())
		return
	}

	response.State.RemoveResource(ctx)
	response.Diagnostics.Append(response.Identity.Set(ctx, resource.Identity())...)
}

func (r Resource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

func (r Resource) ModifyPlan(ctx context.Context, request resource.ModifyPlanRequest, response *resource.ModifyPlanResponse) {
	res := &data.Resource{}
	response.Diagnostics.Append(request.Plan.Get(ctx, &res)...)
	if response.Diagnostics.HasError() {
		return
	}

	if _, ok := res.Values["id"]; !ok {
		// then resource is unknown or something, which means we can't check
		// it
		return
	}

	id := res.GetId()
	if slices.Contains(r.DeferChanges, id) {
		// Then we want to defer this change!

		if !request.ClientCapabilities.DeferralAllowed {
			response.Diagnostics.AddAttributeError(path.Root("id"), "Invalid resource deferral", "This `id` was marked as \"should be deferred\", but the current version of Terraform does not support deferrals.")
			return
		}

		response.Deferred = &resource.Deferred{
			// not technically true, but the closest we have
			Reason: resource.DeferredReasonResourceConfigUnknown,
		}
	}
}
