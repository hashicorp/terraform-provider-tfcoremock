provider "tfcoremock" {}

resource "tfcoremock_dynamic_resource" "test" {
  integer = 1
}
