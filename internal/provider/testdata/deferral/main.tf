
provider "tfcoremock" {
  defer_changes = ["defer_me"]
}

resource "tfcoremock_simple_resource" "resource" {
  id = "defer_me"
}
