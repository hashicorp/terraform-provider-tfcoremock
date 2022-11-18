## 0.1.0

First release of the Mock terraform provider.

FEATURES:

* `mock_simple_resource`: Resource and data source for a simple resource that can model numbers, strings, and booleans.
* `mock_complex_resource`: Resource and data source for a complex resource that can model nested blocks, lists, sets, maps and objects.
* Reads a `dynamic_resources.json` file to allow the user to specify additional resources and data sources dynamically.
  * Add support for computed attributes within dynamic resources. ([#5](https://github.com/hashicorp/terraform-provider-mock/pull/5))
