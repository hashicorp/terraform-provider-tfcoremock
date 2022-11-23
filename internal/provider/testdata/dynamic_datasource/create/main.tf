provider "tfcoremock" {}

data "tfcoremock_simple_resource" "test" {
  id = "simple_resource"
}

resource "tfcoremock_dynamic_resource" "test" {
  id = "my_dynamic_resource"
  my_value = data.tfcoremock_simple_resource.test.integer
}
