package dynamic

import "github.com/hashicorp/terraform-provider-tfcoremock/internal/schema"

type Resources struct {
	RemoteProviders  map[string]string        `json:"remote_providers"`
	DynamicResources map[string]schema.Schema `json:"dynamic_resources"`
}
