data "mock_string_data_source" "example" {
  id = "data_source"
}

resource "mock_string_resource" "example" {
  my_value = data.mock_string_data_source.example.my_value
}
