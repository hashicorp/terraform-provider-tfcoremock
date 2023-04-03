resource "tfcoremock_simple_resource" "example_one" {
  string = "resource_module1"
}

data "tfcoremock_simple_resource" "example_data" {
  id = "simple_resource"
  depends_on = [
    tfcoremock_simple_resource.example_one
  ]
}

resource "tfcoremock_simple_resource" "example_two" {
  string = data.tfcoremock_simple_resource.example_data.string
}
