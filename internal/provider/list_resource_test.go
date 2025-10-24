package provider

import (
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccSimpleResourceList(t *testing.T) {
	t.Cleanup(CleanupTestingDirectories(t))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(""),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.14.0"))),
		},
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/list/simple/main.tf"),
			},
			{
				Query:  true,
				Config: LoadFile(t, "testdata/list/simple/main.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectIdentity("tfcoremock_simple_resource.resource", map[string]knownvalue.Check{
						"id": knownvalue.StringExact("one"),
					}),
					querycheck.ExpectIdentity("tfcoremock_simple_resource.resource", map[string]knownvalue.Check{
						"id": knownvalue.StringExact("two"),
					}),
				},
			},
		},
	})
}

func TestAccDynamicResourceList(t *testing.T) {
	t.Cleanup(CleanupTestingDirectories(t))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(LoadFile(t, "testdata/list/dynamic/dynamic_resources.json")),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.14.0"))),
		},
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/list/dynamic/main.tf"),
			},
			{
				Query:  true,
				Config: LoadFile(t, "testdata/list/dynamic/main.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectIdentity("tfcoremock_dynamic_resource.resource", map[string]knownvalue.Check{
						"id": knownvalue.StringExact("one"),
					}),
					querycheck.ExpectKnownValue("tfcoremock_dynamic_resource.resource", "one", tfjsonpath.New("value"), knownvalue.StringExact("hello, world")),
				},
			},
		},
	})
}
