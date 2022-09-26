terraform {
  required_providers {
    mock = {
      source  = "terraform.local/local/mock"
      version = "0.0.1"
    }
  }
}

provider "mock" {

}

resource "mock_string_resource" "example" {
  my_value = "Hello, world!"
}

resource "mock_integer_resource" "example" {
  my_value = 0
}
