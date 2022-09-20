provider "mock" {}

resource "mock_simple_resource" "test" {
  id = "my_id"
  string = "hello"
}
