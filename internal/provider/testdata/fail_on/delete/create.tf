provider "tfcoremock" {
  fail_on_delete = ["iden"]
}

resource "tfcoremock_simple_resource" "resource" {
  id = "iden"
}
