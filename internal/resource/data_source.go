// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-provider-tfcoremock/internal/client"
	"github.com/hashicorp/terraform-provider-tfcoremock/internal/data"
	"github.com/hashicorp/terraform-provider-tfcoremock/internal/schema"
)

var _ datasource.DataSource = DataSource{}

type DataSource struct {
	Name           string
	InternalSchema schema.Schema
	Client         client.Client
}

func (d DataSource) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = d.Name
}

func (d DataSource) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	var err error
	if response.Schema, err = d.InternalSchema.ToTerraformDataSourceSchema(); err != nil {
		response.Diagnostics.Append(diag.NewErrorDiagnostic(fmt.Sprintf("failed to build data source schema for '%s'", d.Name), err.Error()))
	}
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
