package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-provider-mock/internal/client"
	"github.com/hashicorp/terraform-provider-mock/internal/resource"
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

func (m *mockProvider) Resources(ctx context.Context) []func() tfresource.Resource {
	schemas, err := m.reader.Read()
	if err != nil {
		tflog.Error(ctx, err.Error())
		return nil
	}

	resources := []func() tfresource.Resource{
		func() tfresource.Resource {
			return resource.Resource{
				Name:   "mock_complex_resource",
				Schema: complex.Schema(3),
				Client: m.client,
			}
		},
		func() tfresource.Resource {
			return resource.Resource{
				Name:   "mock_simple_resource",
				Schema: simple.Schema,
				Client: m.client,
			}
		},
	}

	for name, schema := range schemas {
		resources = append(resources, func() tfresource.Resource {
			return resource.Resource{
				Name:   name,
				Schema: schema,
				Client: m.client,
			}
		})
	}

	return resources
}

func (m *mockProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	schemas, err := m.reader.Read()
	if err != nil {
		tflog.Error(ctx, err.Error())
		return nil
	}

	datasources := []func() datasource.DataSource{
		func() datasource.DataSource {
			return resource.DataSource{
				Name:   "mock_complex_resource",
				Schema: complex.Schema(3),
				Client: m.client,
			}
		},
		func() datasource.DataSource {
			return resource.DataSource{
				Name:   "mock_simple_resource",
				Schema: simple.Schema,
				Client: m.client,
			}
		},
	}

	for name, schema := range schemas {
		datasources = append(datasources, func() datasource.DataSource {
			return resource.DataSource{
				Name:   name,
				Schema: schema,
				Client: m.client,
			}
		})
	}

	return datasources
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
