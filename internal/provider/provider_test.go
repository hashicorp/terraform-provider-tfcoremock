// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"errors"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func ProviderFactories(resources string) map[string]func() (tfprotov6.ProviderServer, error) {
	provider := NewForTesting("test", resources)()
	return map[string]func() (tfprotov6.ProviderServer, error){
		"tfcoremock": providerserver.NewProtocol6WithError(provider),
	}
}

func LoadFile(t *testing.T, file string) string {
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("could not read file %s: %v", file, err.Error())
	}

	return string(data)
}

func CleanupTestingDirectories(t *testing.T) func() {
	return func() {
		_, err := os.ReadDir("terraform.resource")
		if err != nil {
			if os.IsNotExist(err) {
				return // Then it's fine.
			}

			t.Fatalf("could not read the resource directory for cleanup: %v", err)
		}
		t.Fatalf("test should have deleted the resource directory on completion")
	}
}

func SaveResourceId(name string, id *string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		module := state.RootModule()
		rs, ok := module.Resources[name]
		if !ok {
			return errors.New("missing resource " + name)
		}

		*id = rs.Primary.Attributes["id"]
		return nil
	}
}

func CheckResourceIdChanged(name string, id *string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		module := state.RootModule()
		rs, ok := module.Resources[name]
		if !ok {
			return errors.New("missing resource " + name)
		}

		if *id == rs.Primary.Attributes["id"] {
			return errors.New("id value for " + name + " has not changed")
		}
		return nil
	}
}
