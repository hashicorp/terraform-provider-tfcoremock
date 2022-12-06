## v0.1.2 (06 Dec 2022)

FEATURES:

* Introduce the `TFCOREMOCK_DYNAMIC_RESOURCES_FILE` environment variable. The location of the `dynamic_resources.json` file is now customisable.

## v0.1.1 (24 Nov 2022)

FEATURES:

* `sensitive`: Resource and data source attributes can be marked as sensitive, meaning they will be elided in Terraform plans and logs.
* `replace`: Resource and data source attributes can be marked as forcing a replacement, meaning that when these attributes are modified the resource will be destroyed and recreated instead of just updated.
* `skip_nested_metadata`: Resource and data source complex attributes can be created without embedded metadata. This doesn't change anything when editing Terraform config, but it changes the underlying format of the attributes and removes optional, sensitive, replacement metadata from attributes nested within complex attributes marked with this field.

## v0.1.0 (22 Nov 2022)

First release of the Terraform Core Mock terraform provider.

FEATURES:

* `tfcoremock_simple_resource`: Resource and data source for a simple resource that can model numbers, strings, and booleans.
* `tfcoremock_complex_resource`: Resource and data source for a complex resource that can model nested blocks, lists, sets, maps and objects.
* Reads a `dynamic_resources.json` file to allow the user to specify additional resources and data sources dynamically.
  * Add support for computed attributes within dynamic resources. ([#5](https://github.com/hashicorp/terraform-provider-tfcoremock/pull/5))
