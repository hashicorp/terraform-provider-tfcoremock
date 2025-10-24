package provider

import (
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccSimpleAction(t *testing.T) {
	t.Cleanup(CleanupTestingDirectories(t))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(""),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.14.0-beta1"))),
		},
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/actions/simple.tf"),
			},
		},
	})
}

func TestAccDynamicAction(t *testing.T) {
	t.Cleanup(CleanupTestingDirectories(t))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(LoadFile(t, "testdata/actions/dynamic_resources.json")),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.14.0-beta1"))),
		},
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/actions/dynamic.tf"),
			},
		},
	})
}
