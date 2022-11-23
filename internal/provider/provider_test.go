package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
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
		files, err := os.ReadDir("terraform.resource")
		if err != nil {
			if os.IsNotExist(err) {
				return // Then it's fine.
			}

			t.Fatalf("could not read the resource directory for cleanup: %v", err)
		}
		defer os.Remove("terraform.resource")

		if len(files) != 0 {
			t.Fatalf("failed to tidy up after test")
		}
	}
}
