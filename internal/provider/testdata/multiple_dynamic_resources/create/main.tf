provider "mock" {}

resource "mock_integer" "integer" {
  integer = 404
}

resource "mock_string" "string" {
  string = "Hello, world!"
}