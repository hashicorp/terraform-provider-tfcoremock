
provider "tfcoremock" {}

resource "tfcoremock_simple_resource" "one" {
  id = "one"
}

resource "tfcoremock_simple_resource" "two" {
  id = "two"
}
