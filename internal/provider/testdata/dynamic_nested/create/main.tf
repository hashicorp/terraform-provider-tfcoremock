provider "mock" {}

resource "mock_dynamic_resource" "test" {
  string = "hello"

  list = [
    {
      string = "one"
    }
  ]

  object = {
    bool = true
  }
}
