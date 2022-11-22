provider "tfcoremock" {}

resource "tfcoremock_dynamic_resource" "test" {

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
