provider "mock" {}

resource "mock_dynamic_resource" "test" {
  id = "my_id"
  integer = 0
}
