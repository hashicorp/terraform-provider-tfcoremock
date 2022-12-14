provider "tfcoremock" {}

resource "tfcoremock_complex_resource" "test" {
  string = "hello"

  list = [
    {
      string = "zero"
    },
    {
      string = "one"
    }
  ]

  object = {
    string = "world"
  }

  set = [
    {
      string = "zero"
    },
  ]
}
