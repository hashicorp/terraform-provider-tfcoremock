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

data "mock_string_data_source" "example" {
  id = "data_source"
}

resource "mock_string_resource" "example" {
  my_value = data.mock_string_data_source.example.my_value
}
