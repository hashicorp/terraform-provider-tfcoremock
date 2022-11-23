provider "tfcoremock" {}

resource "tfcoremock_dynamic_resource" "test" {
  string = "hello"

  list = [
    {
      string = "one"
    }
  ]

  metadata_free_list = [
    {
      string = "other"
    }
  ]

  object = {
    bool = true
  }
}
