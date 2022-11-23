provider "tfcoremock" {}

resource "tfcoremock_dynamic_resource" "test" {
  set = [
    {
      "custom-value": "zero",
    },
    {
      "custom-value": "one",
    },
    {
      "custom-value": "two",
    }
  ]
}
