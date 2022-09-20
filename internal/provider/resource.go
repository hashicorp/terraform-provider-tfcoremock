package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	"github.com/hashicorp/terraform-provider-mock/internal/resource"
	"github.com/hashicorp/terraform-provider-mock/internal/schema"
)

var _ provider.ResourceType = resourceType{}

type resourceType struct {
	Schema schema.Schema
}

func (c resourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	schema, err := c.Schema.ToTerraformResourceSchema()
	if err != nil {
		diags := diag.Diagnostics{}
		diags.Append(diag.NewErrorDiagnostic("could not generate resource schema", err.Error()))
		return tfsdk.Schema{}, diags
	}
	return schema, nil
}

func (c resourceType) NewResource(ctx context.Context, p provider.Provider) (tfresource.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(p)
	if diags.HasError() {
		return nil, diags
	}

	return resource.Resource{
		Schema: c.Schema,
		Client: provider.client,
	}, nil
}
