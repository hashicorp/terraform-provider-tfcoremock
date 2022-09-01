provider "mock" {}

resource "mock_dynamic_resource" "test" {
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
