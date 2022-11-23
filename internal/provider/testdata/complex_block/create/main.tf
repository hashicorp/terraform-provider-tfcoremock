provider "tfcoremock" {}

resource "tfcoremock_complex_resource" "test" {
  string = "hello"

  list_block {
    integer = 0
  }

  list_block {
    integer = 1
  }

  set_block {
    integer = 0
  }
}
