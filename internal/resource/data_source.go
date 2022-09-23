package resource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-provider-mock/internal/client"
	"github.com/hashicorp/terraform-provider-mock/internal/data"
	"github.com/hashicorp/terraform-provider-mock/internal/schema"
)

var _ datasource.DataSource = DataSource{}

type DataSource struct {
	Name   string
	Schema schema.Schema
	Client client.Client
}

func (d DataSource) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = d.Name
}

func (d DataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	var diags diag.Diagnostics

	schema, err := d.Schema.ToTerraformDataSourceSchema()
	if err != nil {
		diags.Append(diag.NewErrorDiagnostic("failed to build data source schea", err.Error()))
	}

	return schema, diags
}

func (d DataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	resource := &data.Resource{}

	response.Diagnostics.Append(request.Config.Get(ctx, &resource)...)
	if response.Diagnostics.HasError() {
		return
	}

	data, err := d.Client.ReadDataSource(ctx, resource.GetId())
	if err != nil {
		response.Diagnostics.AddError("failed to read data source", err.Error())
		return
	}

	if data == nil {
		response.Diagnostics.AddError(
			"target data source does not exist",
			fmt.Sprintf("data source at %s could not be found in data directory", resource.GetId()))
	}

	typ := request.Config.Schema.Type().TerraformType(ctx)
	response.Diagnostics.Append(response.State.Set(ctx, data.WithType(typ.(tftypes.Object)))...)
}
