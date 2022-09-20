package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-provider-mock/internal/client"
	"github.com/hashicorp/terraform-provider-mock/internal/schema/complex"
	"github.com/hashicorp/terraform-provider-mock/internal/schema/dynamic"
	"github.com/hashicorp/terraform-provider-mock/internal/schema/simple"
)

var _ provider.Provider = &mockProvider{}

type mockProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string

	// reader will read the dynamic resource definitions in the GetResource and
	// GetDataSources functions.
	reader dynamic.Reader

	// client is provided to the actual resources so that their states can be
	// recorded and written to a backend other than the terraform state.
	client client.Client
}

type providerData struct {
	ResourceDirectory types.String `tfsdk:"resource_directory"`
	DataDirectory     types.String `tfsdk:"data_directory"`
	UseOnlyState      types.Bool   `tfsdk:"use_only_state"`
}

func (m *mockProvider) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
	var data providerData
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	if data.UseOnlyState.Value {
		directory := "terraform.data"
		if !data.DataDirectory.IsNull() {
			directory = data.DataDirectory.Value
		}

		m.client = client.State{
			DataDirectory: directory,
		}
	} else {
		dataDirectory := "terraform.data"
		resourceDirectory := "terraform.resource"

		if !data.DataDirectory.IsNull() {
			dataDirectory = data.DataDirectory.Value
		}

		if !data.ResourceDirectory.IsNull() {
			resourceDirectory = data.ResourceDirectory.Value
		}

		m.client = client.Local{
			ResourceDirectory: resourceDirectory,
			DataDirectory:     dataDirectory,
		}
	}
}

func (m *mockProvider) GetResources(ctx context.Context) (map[string]provider.ResourceType, diag.Diagnostics) {
	schemas, err := m.reader.Read()
	if err != nil {
		var diags diag.Diagnostics
		diags.Append(diag.NewErrorDiagnostic("could not read dynamic resources", err.Error()))
		return nil, diags
	}

	resources := make(map[string]provider.ResourceType)
	for name, schema := range schemas {
		resources[name] = resourceType{
			Schema: schema,
		}
	}

	resources["mock_complex_resource"] = resourceType{
		Schema: complex.Schema(3),
	}
	resources["mock_simple_resource"] = resourceType{
		Schema: simple.Schema,
	}

	return resources, nil
}

func (m *mockProvider) GetDataSources(ctx context.Context) (map[string]provider.DataSourceType, diag.Diagnostics) {
	schemas, err := m.reader.Read()
	if err != nil {
		var diags diag.Diagnostics
		diags.Append(diag.NewErrorDiagnostic("could not read dynamic resources", err.Error()))
		return nil, diags
	}

	sources := make(map[string]provider.DataSourceType)
	for name, schema := range schemas {
		sources[name] = dataSourceType{
			Schema: schema,
		}
	}

	sources["mock_complex_resource"] = dataSourceType{
		Schema: complex.Schema(3),
	}
	sources["mock_simple_resource"] = dataSourceType{
		Schema: simple.Schema,
	}

	return sources, nil
}

func (m *mockProvider) GetSchema(context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"resource_directory": {
				Optional: true,
				Type:     types.StringType,
			},
			"data_directory": {
				Optional: true,
				Type:     types.StringType,
			},
			"use_only_state": {
				Optional: true,
				Type:     types.BoolType,
			},
		},
	}, nil
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &mockProvider{
			version: version,
			// TODO(liamcervante): Turn this into an environment variable?
			reader: dynamic.FileReader{File: "dynamic_resources.json"},
		}
	}
}

func NewForTesting(version string, resources string) func() provider.Provider {
	return func() provider.Provider {
		return &mockProvider{
			version: version,
			reader:  dynamic.StringReader{Data: resources},
		}
	}
}

// convertProviderType is a helper function for NewResource and NewDataSource
// implementations to associate the concrete provider type. Alternatively,
// this helper can be skipped and the provider type can be directly type
// asserted (e.g. provider: in.(*scaffoldingProvider)), however using this can prevent
// potential panics.
func convertProviderType(in provider.Provider) (mockProvider, diag.Diagnostics) {
	var diags diag.Diagnostics

	p, ok := in.(*mockProvider)

	if !ok {
		diags.AddError(
			"Unexpected Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received. This is always a bug in the provider code and should be reported to the provider developers.", p),
		)
		return mockProvider{}, diags
	}

	if p == nil {
		diags.AddError(
			"Unexpected Provider Instance Type",
			"While creating the data source or resource, an unexpected empty provider instance was received. This is always a bug in the provider code and should be reported to the provider developers.",
		)
		return mockProvider{}, diags
	}

	return *p, diags
}
