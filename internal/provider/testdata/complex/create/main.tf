provider "mock" {}

resource "mock_complex_resource" "test" {
  string = "hello"

  list = [
    {
      string = "one"
    }
  ]

  object = {
    bool = true
  }

  set = [
    {
      string = "zero"
    },
    {
      string = "one"
    }
  ]
}
