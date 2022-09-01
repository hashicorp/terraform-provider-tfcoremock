provider "mock" {}

resource "mock_dynamic_resource" "test" {

  other {
    id = "my-id"

    nested {

    }

    nested {

    }
  }

  other {

  }
}
