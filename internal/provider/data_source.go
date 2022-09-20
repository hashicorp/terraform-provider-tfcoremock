package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	"github.com/hashicorp/terraform-provider-mock/internal/resource"
	"github.com/hashicorp/terraform-provider-mock/internal/schema"
)

var _ provider.DataSourceType = dataSourceType{}

type dataSourceType struct {
	Schema schema.Schema
}

func (c dataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	schema, err := c.Schema.ToTerraformDataSourceSchema()
	if err != nil {
		var diags diag.Diagnostics
		diags.Append(diag.NewErrorDiagnostic("could not generate data source schema", err.Error()))
		return tfsdk.Schema{}, diags
	}
	return schema, nil
}

func (c dataSourceType) NewDataSource(ctx context.Context, p provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(p)
	if diags.HasError() {
		return nil, diags
	}

	return resource.DataSource{
		Schema: c.Schema,
		Client: provider.client,
	}, nil
}
