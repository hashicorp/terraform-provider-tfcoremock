provider "tfcoremock" {
  fail_on_create = ["iden"]
}

resource "tfcoremock_simple_resource" "resource" {
  id = "iden"
}
