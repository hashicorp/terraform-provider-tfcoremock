# Terraform `tfcoremock` Provider

The `tfcoremock` provider is intended to aid with testing the Terraform core libraries
and the Terraform CLI. This provider should allow users to define all possible 
Terraform configurations and run them through the Terraform core platform.

The provider supplies two static resources:

- `tfcoremock_simple_resource`
- `tfcoremock_complex_resource`
 
Users can then define additional dynamic resources by supplying a 
`dynamic_resources.json` file alongside their root Terraform configuration. 
These dynamic resources can be used to model any Terraform configuration not
covered by the provided static resources.

Use the `TFCOREMOCK_DYNAMIC_RESOURCES_FILE` environment variable to customise 
the location of the `dynamic_resources.json` file. By default, the provider 
looks for the `dynamic_resources.json` file in the same directory as the 
Terraform config files, but using this environment variable allows the dynamic
resources to be defined by any file on the system. For example: 
`TFCOREMOCK_DYNAMIC_RESOURCES_FILE=/path/to/resources.json terraform plan`

By default, all resources created by the provider are then converted into a 
human-readable JSON format and written out to the resource directory. This 
behaviour can be disabled by turning on the `use_only_state` flag in the 
provider schema (this is useful when running the provider in a Terraform Cloud
environment). The resource directory defaults to `terraform.resource`.

All resources supplied by the provider (including the simple and 
complex resource as well as any dynamic resources) are duplicated into data 
sources. The data sources should be supplied in the JSON format that resources
are written into. The provider looks into the data directory, which defaults to
`terraform.data`.

All resources (and data sources) supplied by the provider have an `id` 
attribute that is generated if not set by the configuration. Dynamic resources 
cannot define an `id` attribute as the provider will create one for them. The 
`id` attribute is used as name of the human-readable JSON files held in the
resource and data directories.

Additionally, all resources are available to be queried via `list` blocks. For
now only the `id` attribute is supported as a field to retrieve a specific 
instance. It is optional, so all resources of the specified type will be 
returned if the field is left blank.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Using the provider

We provide a simple example here. View the [examples](./examples) and 
[docs](./docs) subdirectories for more examples.

In this example, we have a `tfcoremock_simple_resource` defined as a data source with
an `id` of `my-simple-resource`. This means we create a file 
`terraform.data/my-simple-resource.json` which defines a simple resource with
a single integer set. We then define a dynamic resource called 
`tfcoremock_dynamic_resource`. The dynamic resource holds a single integer, and is 
defined in the `dynamic_resources.json` file.

Note, that we do not define an  `id` field for this resource when we provide the
definition. Despite this, we can still provide a value for the `id` in the
configuration because the provider ensures that all resources have this attribute.
In this example, we do provide a value for the `id` field. If we didn't, the provider
would generate one for us.

The following subsections show the Terraform configuration pre-apply and then
show the extra files created post-apply.

### Pre-apply

#### **./main.tf**
```hcl
terraform {
  required_providers {
    tfcoremock = {
      source  = "hashicorp/tfcoremock"
      version = "0.1.2"
    }
  }
}

data "tfcoremock_simple_resource" "my_simple_resource" {
  id = "my-simple-resource"
}

resource "tfcoremock_dynamic_resource" "my_dynamic_resource" {
  id = "my-dynamic-resource"
  my_value = data.tfcoremock_simple_resource.my_simple_resource.integer + 1
}
```

#### **./terraform.data/my-simple-resource.json**
```json
{
  "values": {
    "integer": {
      "number": "0"
    },
    "id": {
      "string": "my-simple-resource"
    }
  }
}
```

#### **./dynamic_resources.json**
```json
{
  "tfcoremock_dynamic_resource": {
    "attributes": {
      "my_value": {
        "type": "integer",
        "required": true
      }
    }
  }
}
```

### Post apply

In addition to the normal Terraform state and lock files, you will see the new
resource we created has been written into the resource directory.

#### **./terraform.resource/my-dynamic-resource.json**
```json
{
  "values": {
    "id": {
      "string": "my-dynamic-resource"
    },
    "my_value": {
      "number": "1"
    }
  }
}
```


## Developing the Provider

If you wish to work on the provider, you'll first need 
[Go](http://www.golang.org) installed on your machine 
(see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put 
the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.
