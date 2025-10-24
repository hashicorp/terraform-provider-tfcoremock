
resource "tfcoremock_simple_resource" "resource" {
  lifecycle {
    action_trigger {
      events = [before_create, before_update]
      actions = [action.tfcoremock_dynamic_action.action]
    }
  }
}

action "tfcoremock_dynamic_action" "action" {
  config {
    integer = 0
  }
}
