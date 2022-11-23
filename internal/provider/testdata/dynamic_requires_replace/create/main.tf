provider "tfcoremock" {}

resource "tfcoremock_list_of_objects" "test" {
  list = [
    {
      key   = "one"
      value = "first value"
    },
    {
      key   = "two"
      value = "second value"
    },
  ]
}
