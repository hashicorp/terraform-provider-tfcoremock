provider mock {}

data "mock_simple_resource" "data" {
  id = "simple_resource"
}

resource "mock_simple_resource" "test" {
  integer = data.mock_simple_resource.data.integer
}
