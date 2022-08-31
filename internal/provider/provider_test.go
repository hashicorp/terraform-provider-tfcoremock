package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/stretchr/testify/require"
)

func ProviderFactories(resources string) map[string]func() (tfprotov6.ProviderServer, error) {
	provider := NewForTesting("test", resources)()
	return map[string]func() (tfprotov6.ProviderServer, error){
		"mock": providerserver.NewProtocol6WithError(provider),
	}
}

func LoadFile(t *testing.T, file string) string {
	data, err := os.ReadFile(file)
	require.NoError(t, err)

	return string(data)
}

func CleanupTestingDirectories(t *testing.T) {
	files, err := os.ReadDir("terraform.resource")
	if err != nil {
		if os.IsNotExist(err) {
			return // Then it's fine.
		}

		require.NoError(t, err)
	}
	defer os.Remove("terraform.resource")

	if len(files) != 0 {
		require.Fail(t, "failed to tidy up after test")
	}
}
