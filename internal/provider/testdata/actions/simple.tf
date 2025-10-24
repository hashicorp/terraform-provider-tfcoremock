
resource "tfcoremock_simple_resource" "resource" {
  lifecycle {
    action_trigger {
      events = [before_create, before_update]
      actions = [action.tfcoremock_simple_resource.action]
    }
  }
}

action "tfcoremock_simple_resource" "action" {}
