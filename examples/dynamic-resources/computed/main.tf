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

resource "mock_dynamic_resource" "example" {}
