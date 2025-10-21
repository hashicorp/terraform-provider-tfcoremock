
resource "tfcoremock_simple_resource" "main" {}

resource "tfcoremock_simple_resource" "other" {
  id = tfcoremock_simple_resource.main.id
}
