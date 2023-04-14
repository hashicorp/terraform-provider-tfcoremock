// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	provider_schema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-provider-tfcoremock/internal/client"
	"github.com/hashicorp/terraform-provider-tfcoremock/internal/resource"
	"github.com/hashicorp/terraform-provider-tfcoremock/internal/schema/complex"
	"github.com/hashicorp/terraform-provider-tfcoremock/internal/schema/dynamic"
	"github.com/hashicorp/terraform-provider-tfcoremock/internal/schema/simple"
)

var _ provider.Provider = &tfcoremockProvider{}

const (
	description = `The 'tfcoremock' provider is intended to aid with testing the Terraform core libraries and the Terraform CLI. This provider should allow users to define all possible Terraform configurations and run them through the Terraform core platform.

The provider supplies two static resources:

- 'tfcoremock_simple_resource'
- 'tfcoremock_complex_resource'
 
Users can then define additional dynamic resources by supplying a 'dynamic_resources.json' file alongside their root Terraform configuration. These dynamic resources can be used to model any Terraform configuration not covered by the provided static resources.

By default, all resources created by the provider are then converted into a human-readable JSON format and written out to the resource directory. This behaviour can be disabled by turning on the 'use_only_state' flag in the provider schema (this is useful when running the provider in a Terraform Cloud environment). The resource directory defaults to 'terraform.resource'.

All resources supplied by the provider (including the simple and complex resource as well as any dynamic resources) are duplicated into data sources. The data sources should be supplied in the JSON format that resources are written into. The provider looks into the data directory, which defaults to 'terraform.data'.

Finally, all resources (and data sources) supplied by the provider have an 'id' attribute that is generated if not set by the configuration. Dynamic resources cannot define an 'id' attribute as the provider will create one for them. The 'id' attribute is used as name of the human-readable JSON files held in the resource and data directories.`

	markdownDescription = `The ''tfcoremock'' provider is intended to aid with testing the Terraform core libraries and the Terraform CLI. This provider should allow users to define all possible Terraform configurations and run them through the Terraform core platform.

The provider supplies two static resources:

- ''tfcoremock_simple_resource''
- ''tfcoremock_complex_resource''
 
Users can then define additional dynamic resources by supplying a ''dynamic_resources.json'' file alongside their root Terraform configuration. These dynamic resources can be used to model any Terraform configuration not covered by the provided static resources.

By default, all resources created by the provider are then converted into a human-readable JSON format and written out to the resource directory. This behaviour can be disabled by turning on the ''use_only_state'' flag in the provider schema (this is useful when running the provider in a Terraform Cloud environment). The resource directory defaults to ''terraform.resource''.

All resources supplied by the provider (including the simple and complex resource as well as any dynamic resources) are duplicated into data sources. The data sources should be supplied in the JSON format that resources are written into. The provider looks into the data directory, which defaults to ''terraform.data''.

Finally, all resources (and data sources) supplied by the provider have an ''id'' attribute that is generated if not set by the configuration. Dynamic resources cannot define an ''id'' attribute as the provider will create one for them. The ''id'' attribute is used as name of the human-readable JSON files held in the resource and data directories.`

	dynamicResourcesPathEnvVarName = "TFCOREMOCK_DYNAMIC_RESOURCES_FILE"
)

type tfcoremockProvider struct {
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

func (m *tfcoremockProvider) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
	var data providerData
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	if data.UseOnlyState.ValueBool() {
		directory := "terraform.data"
		if !data.DataDirectory.IsNull() {
			directory = data.DataDirectory.String()
		}

		m.client = client.State{
			DataDirectory: directory,
		}
	} else {
		dataDirectory := "terraform.data"
		resourceDirectory := "terraform.resource"

		if !data.DataDirectory.IsNull() {
			dataDirectory = data.DataDirectory.ValueString()
		}

		if !data.ResourceDirectory.IsNull() {
			resourceDirectory = data.ResourceDirectory.ValueString()
		}

		m.client = client.Local{
			ResourceDirectory: resourceDirectory,
			DataDirectory:     dataDirectory,
		}
	}
}

func (m *tfcoremockProvider) Metadata(ctx context.Context, request provider.MetadataRequest, response *provider.MetadataResponse) {
	response.Version = m.version
	response.TypeName = "tfcoremock"
}

func (m *tfcoremockProvider) Resources(ctx context.Context) []func() tfresource.Resource {
	resources := []func() tfresource.Resource{
		func() tfresource.Resource {
			return resource.Resource{
				Name:           "tfcoremock_complex_resource",
				InternalSchema: complex.Schema(3),
				Client:         m.client,
			}
		},
		func() tfresource.Resource {
			return resource.Resource{
				Name:           "tfcoremock_simple_resource",
				InternalSchema: simple.Schema,
				Client:         m.client,
			}
		},
	}

	schemas, err := m.reader.Read()
	if err != nil {
		// This isn't ideal, as the plugin will tell the user this is a problem
		// with the provider. It's not though, this means the provided dynamic
		// resources file either wasn't valid JSON or didn't match our schema.
		//
		// We don't have a way to raise an error through the plugin at this
		// point in time though, so the only thing we can really do is panic.
		//
		// We add a lot of context to this panic to try and make the caller
		// realise exactly what the problem is.
		panic(fmt.Sprintf("The tfcoremock provider either failed to parse or failed to validate your dynamic resources file. "+
			"Terraform will say this is a problem in the provider, but in this case it is a problem in your dynamic resources file. "+
			"We have the following error from the parser, hopefully it provides additional context about the problem but these errors are not always helpful."+
			"\n\n%s\n", err.Error()))
	}

	for name, schema := range schemas {
		resourceName := name
		resourceSchema := schema
		resources = append(resources, func() tfresource.Resource {
			return resource.Resource{
				Name:           resourceName,
				InternalSchema: resourceSchema,
				Client:         m.client,
			}
		})
	}

	return resources
}

func (m *tfcoremockProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	datasources := []func() datasource.DataSource{
		func() datasource.DataSource {
			return resource.DataSource{
				Name:           "tfcoremock_complex_resource",
				InternalSchema: complex.Schema(3),
				Client:         m.client,
			}
		},
		func() datasource.DataSource {
			return resource.DataSource{
				Name:           "tfcoremock_simple_resource",
				InternalSchema: simple.Schema,
				Client:         m.client,
			}
		},
	}

	schemas, err := m.reader.Read()
	if err != nil {
		// This isn't ideal, as the plugin will tell the user this is a problem
		// with the provider. It's not though, this means the provided dynamic
		// resources file either wasn't valid JSON or didn't match our schema.
		//
		// We don't have a way to raise an error through the plugin at this
		// point in time though, so the only thing we can really do is panic.
		//
		// We add a lot of context to this panic to try and make the caller
		// realise exactly what the problem is.
		panic(fmt.Sprintf("The tfcoremock provider either failed to parse or failed to validate your dynamic resources file. "+
			"Terraform will say this is a problem in the provider, but in this case it is a problem in your dynamic resources file. "+
			"We have the following error from the parser, hopefully it provides additional context about the problem but these errors are not always helpful."+
			"\n\n%s\n", err.Error()))
	}

	for name, schema := range schemas {
		datasourceName := name
		datasourceSchema := schema
		datasources = append(datasources, func() datasource.DataSource {
			return resource.DataSource{
				Name:           datasourceName,
				InternalSchema: datasourceSchema,
				Client:         m.client,
			}
		})
	}

	return datasources
}

func (m *tfcoremockProvider) Schema(ctx context.Context, request provider.SchemaRequest, response *provider.SchemaResponse) {
	response.Schema = provider_schema.Schema{
		Description:         description,
		MarkdownDescription: strings.ReplaceAll(markdownDescription, "''", "`"),
		Attributes: map[string]provider_schema.Attribute{
			"resource_directory": provider_schema.StringAttribute{
				Description:         "The directory that the provider should use to write the human-readable JSON files for each managed resource. If `use_only_state` is set to `true` then this value does not matter. Defaults to `terraform.resource`.",
				MarkdownDescription: "The directory that the provider should use to write the human-readable JSON files for each managed resource. If `use_only_state` is set to `true` then this value does not matter. Defaults to `terraform.resource`.",
				Optional:            true,
			},
			"data_directory": provider_schema.StringAttribute{
				Description:         "The directory that the provider should use to read the human-readable JSON files for each requested data source. Defaults to `data.resource`.",
				MarkdownDescription: "The directory that the provider should use to read the human-readable JSON files for each requested data source. Defaults to `data.resource`.",
				Optional:            true,
			},
			"use_only_state": provider_schema.BoolAttribute{
				Description:         "If set to true the provider will rely only on the Terraform state file to load managed resources and will not write anything to disk. Defaults to `false`.",
				MarkdownDescription: "If set to true the provider will rely only on the Terraform state file to load managed resources and will not write anything to disk. Defaults to `false`.",
				Optional:            true,
			},
		},
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		dynamicResourcesPath := "dynamic_resources.json"
		if dynamicResourcesPathEnvVar := os.Getenv(dynamicResourcesPathEnvVarName); len(dynamicResourcesPathEnvVar) > 0 {
			dynamicResourcesPath = dynamicResourcesPathEnvVar
		}

		return &tfcoremockProvider{
			version: version,
			reader:  dynamic.FileReader{File: dynamicResourcesPath},
		}
	}
}

func NewForTesting(version string, resources string) func() provider.Provider {
	return func() provider.Provider {
		return &tfcoremockProvider{
			version: version,
			reader:  dynamic.StringReader{Data: resources},
		}
	}
}
