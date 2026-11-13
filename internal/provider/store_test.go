// Copyright IBM Corp. 2022, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

/*
	TODO: `StateStore: true` doesn't test lock contention. How do we enable testing that?
*/

func TestAccStateStore_Fs_init(t *testing.T) {
	t.Skip("skipping this test as the tfcoremock repo doesn't enable running tests that include experimental features")
	
	// Setting this environment variable ensures TF core uses pluggable state storage during init.
	// This is only temporary while PSS is experimental.
	t.Setenv("TF_ENABLE_PLUGGABLE_STATE_STORAGE", "1")

	td := t.TempDir()

	t.Cleanup(CleanupTestingDirectories(t))

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(""),
		Steps: []resource.TestStep{
			{
				StateStore: true, // Activates a bunch of pre-packaged testing steps...
				Config:     testConfigWithFsStateStore(td),
			},
		},
	})
}

func testConfigWithFsStateStore(td string) string {
	return fmt.Sprintf(`
terraform {
  required_providers {
    tfcoremock = {
      source = "hashicorp/tfcoremock"
    }
  }

  state_store "tfcoremock_fs" {
    provider "tfcoremock" {
      # No config
    }

    workspace_dir = "%s"
  }
}
`, td)
}
