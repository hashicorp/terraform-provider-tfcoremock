data "tfcoremock_simple_resource" "example" {
  id = "data_source"
}

resource "tfcoremock_dynamic_resource" "example" {
  my_value = data.tfcoremock_simple_resource.example.integer
}
