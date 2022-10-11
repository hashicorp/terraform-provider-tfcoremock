data "mock_simple_resource" "example" {
  id = "data_source"
}

resource "mock_dynamic_resource" "example" {
  my_value = data.mock_simple_resource.example.integer
}
