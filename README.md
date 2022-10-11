# Terraform `mock` Provider

The `mock` provider is intended to aid with testing the Terraform core libraries
and the Terraform CLI. This provider should allow users to define all possible 
Terraform configurations and run them through the Terraform core platform.

The provider supplies two static resources:

- `mock_simple_resource`
- `mock_complex_resource`
 
Users can then define additional dynamic resources by supplying a 
`dynamic_resources.json` file alongside their root Terraform configuration. 
These dynamic resources can be used to model any Terraform configuration not
covered by the provided static resources.

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

Finally, all resources (and data sources) supplied by the provider have an `id` 
attribute that is generated if not set by the configuration. Dynamic resources 
cannot define an `id` attribute as the provider will create one for them. The 
`id` attribute is used as name of the human-readable JSON files held in the
resource and data directories.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.18

## Using the provider

We provide a simple example here. View the [examples](./examples) and 
[docs](./docs) subdirectories for more examples.

In this example, we have a `mock_simple_resource` defined as a data source with
an identifier of `my_simple_resource`. This means we create a file 
`terraform.data/my_simple_resource.json` which defines a simple resource with
a single integer set. We then define a dynamic resource called 
`mock_dynamic_resource`. The dynamic resource holds a single integer, and is 
defined in the `dynamic_resources.json` file. Note, that we do not define an 
`id` field for this resource when we provide the definition. Despite this, we
can still provide a value for the `id` in the configuration because the provider
ensures that all resources have this attribute. In this example, we do provide
a value for the `id` field. If we didn't the provider would generate one for us.

The following subsections show the Terraform configuration pre-apply and then
show the extra files created post-apply.

### Pre-apply

#### **./main.tf**
```hcl
terraform {
  required_providers {
    mock = {
      source  = "hashicorp/mock"
      version = "1.0.0"
    }
  }
}

provider "mock" {
  
}

data "mock_simple_resource" "my_simple_resource" {
  id = "my_simple_resource"
}

resource "mock_dynamic_resource" "my_dynamic_resource" {
  id = "my_dynamic_resource"
  my_value = data.mock_simple_resource.my_simple_resource.integer
}
```

#### **./terraform.data/my_data_source.json**
```json
{
  "values": {
    "integer": {
      "integer": 0
    },
    "id": {
      "string": "my_simple_resource"
    }
  }
}
```

#### **./dynamic_resources.json**
```json
{
  "mock_dynamic_resource": {
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

#### **./terraform.resource/my_dynamic_resource.json**
```json
{
  "values": {
    "id": {
      "string": "my_dynamic_resource"
    },
    "my_value": {
      "number": "0"
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
