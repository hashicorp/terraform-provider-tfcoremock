package resource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	list_schema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-provider-tfcoremock/internal/client"
	"github.com/hashicorp/terraform-provider-tfcoremock/internal/data"
	"github.com/hashicorp/terraform-provider-tfcoremock/internal/schema"
)

var _ list.ListResource = ListResource{}

type ListResource struct {
	Name           string
	InternalSchema schema.Schema
	Client         client.Client
}

func (l ListResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = l.Name
}

func (l ListResource) ListResourceConfigSchema(ctx context.Context, request list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema.Description = l.InternalSchema.Description
	response.Schema.MarkdownDescription = l.InternalSchema.MarkdownDescription
	response.Schema.Attributes = map[string]list_schema.Attribute{
		"id": list_schema.StringAttribute{
			Optional: true,
		},
	}
}

func (l ListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream) {
	resource := &data.Resource{
		ResourceType: l.Name,
	}

	diags := request.Config.Get(ctx, &resource)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	stream.Results = func(yield func(list.ListResult) bool) {
		var id *string
		if value, ok := resource.Values["id"]; ok {
			id = value.String
		}

		err := l.Client.ListResources(ctx, client.Filter(l.Name), id, func(resource *data.Resource, err error) {
			result := request.NewListResult(ctx)
			if err != nil {
				result.Diagnostics.Append(diag.NewErrorDiagnostic("failed to query resource", err.Error()))
				return
			} else {
				result.DisplayName = resource.GetId()
				result.Diagnostics.Append(result.Identity.Set(ctx, resource.Identity())...)

				if request.IncludeResource {
					typ := request.ResourceSchema.Type().TerraformType(ctx)
					result.Diagnostics.Append(result.Resource.Set(ctx, resource.WithType(typ.(tftypes.Object)))...)
				}
			}
			yield(result)
		}, request.Limit)
		if err != nil {
			yield(list.ListResult{
				Diagnostics: diag.Diagnostics{
					diag.NewErrorDiagnostic("failed to query resources", err.Error()),
				},
			})
		}
	}
}
