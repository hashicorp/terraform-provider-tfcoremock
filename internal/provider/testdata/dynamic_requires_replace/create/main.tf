provider "mock" {}

resource "mock_list_of_objects" "test" {
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
