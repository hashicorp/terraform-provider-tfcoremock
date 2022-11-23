provider "tfcoremock" {}

resource "tfcoremock_list_of_objects" "test" {
  list = [
    {
      key   = "three"
      value = "first value"
    },
    {
      key   = "two"
      value = "second value"
    },
  ]
}
