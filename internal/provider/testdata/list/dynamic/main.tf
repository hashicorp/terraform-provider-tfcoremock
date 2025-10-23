
provider "tfcoremock" {}

resource "tfcoremock_dynamic_resource" "one" {
  id = "one"
  value = "hello, world"
}

resource "tfcoremock_dynamic_resource" "two" {
  id = "two"
  value = "goodbye, world"
}
