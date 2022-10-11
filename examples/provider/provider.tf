provider "mock" {
  resource_directory = "terraform.resource"
  data_directory     = "terraform.data"
  use_only_state     = false
}
