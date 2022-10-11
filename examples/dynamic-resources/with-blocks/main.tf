resource "mock_dynamic_resource" "example" {
  my_values {
    my_value = "Hello, "
  }

  my_values {
    my_value = "world!"
  }
}
