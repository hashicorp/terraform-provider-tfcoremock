provider "tfcoremock" {}

resource "tfcoremock_dynamic_resource" "test" {
  set = [
    {
      "custom_value": "zero",
    },
    {
      "custom_value": "one",
    },
    {
      "custom_value": "two",
    }
  ]
}
