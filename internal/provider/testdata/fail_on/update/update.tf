provider "tfcoremock" {
  fail_on_update = ["iden"]
}

resource "tfcoremock_simple_resource" "resource" {
  id = "iden"
  number = 0
}
