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

data "mock_simple_resource" "example" {
  id = "data_source"
}

resource "mock_dynamic_resource" "example" {
  my_value = data.mock_simple_resource.example.integer
}
