
list "tfcoremock_dynamic_resource" "resource" {
  provider = tfcoremock
  include_resource = true
  config {
    id = "one"
  }
}
