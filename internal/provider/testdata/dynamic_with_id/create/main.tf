provider "tfcoremock" {}

resource "tfcoremock_dynamic_resource" "test" {
  id = "my_id"
  integer = 0
}
