provider "mock" {}

resource "mock_complex_resource" "test" {
  string = "hello"

  list_block {
    integer = 0
  }

  set_block {
    integer = 0
  }

  set_block {
    integer = 1
  }
}
