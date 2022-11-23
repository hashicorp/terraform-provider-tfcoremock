provider "tfcoremock" {}

resource "tfcoremock_integer" "integer" {
  integer = 404
}

resource "tfcoremock_string" "string" {
  string = "Hello, world!"
}