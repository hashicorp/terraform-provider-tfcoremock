## v0.6.0 (Unreleased)

ENHANCEMENTS:

* Add support for mocking actions. ([#191](https://github.com/hashicorp/terraform-provider-tfcoremock/pull/191))
* Introduce `defer_changes` attributes to the provider configuration. This allows controlling if resources should defer there changes during the current operation. ([#190](https://github.com/hashicorp/terraform-provider-tfcoremock/pull/190))

## v0.5.0 (15 Apr 2025)

NOTES:

* Update dependencies.

## v0.4.0 (9 Jan 2025)

ENHANCEMENTS:

* Introduce `fail_on_create`, `fail_on_delete`, `fail_on_read`, `fail_on_update` attributes to the provider configuration. This allows controlling if resources should fail during certain operations. ([#154](https://github.com/hashicorp/terraform-provider-tfcoremock/pull/154))

## v0.3.0 (26 Aug 2024)

ENHANCEMENTS:

* Destroying the last managed resource in a workspace will now cause the provider to also tidy up and remove the resource directory itself. ([#54](https://github.com/hashicorp/terraform-provider-tfcoremock/issues/54))

## v0.2.0 (14 Apr 2023)

ENHANCEMENTS:

* Computed attributes in dynamic resources will no longer create default values, but will return null values by default. Users can still specify concrete values for computed attributes to return. ([#51](https://github.com/hashicorp/terraform-provider-tfcoremock/issues/51))

BUG FIXES:

* Fix bug in which custom values for the resource and data directories were being interpreted incorrectly, meaning custom resource and data directories were unusable. ([#52](https://github.com/hashicorp/terraform-provider-tfcoremock/issues/52))

## v0.1.3 (03 Apr 2023)

BUG FIXES:

* Fix bug in which data sources that were not setting the computed status for their attributes. ([#48](https://github.com/hashicorp/terraform-provider-tfcoremock/issues/48))

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
