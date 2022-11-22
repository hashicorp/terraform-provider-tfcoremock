provider "tfcoremock" {}

resource "tfcoremock_simple_resource" "test" {
  id = "my_id"
  string = "world"
}
