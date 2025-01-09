provider "tfcoremock" {
  fail_on_read = ["iden"]
}

resource "tfcoremock_simple_resource" "resource" {
  id = "iden"
}
