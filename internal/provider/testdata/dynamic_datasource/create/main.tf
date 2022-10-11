provider "mock" {}

data "mock_simple_resource" "test" {
  id = "simple_resource"
}

resource "mock_dynamic_resource" "test" {
  id = "my_dynamic_resource"
  my_value = data.mock_simple_resource.test.integer
}
