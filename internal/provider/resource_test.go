// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccComplexResource(t *testing.T) {
	t.Cleanup(CleanupTestingDirectories(t))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(""),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/complex/create/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tfcoremock_complex_resource.test", "string", "hello"),
					resource.TestCheckResourceAttr("tfcoremock_complex_resource.test", "list.#", "1"),
					resource.TestCheckResourceAttr("tfcoremock_complex_resource.test", "list.0.string", "one"),
					resource.TestCheckResourceAttr("tfcoremock_complex_resource.test", "object.bool", "true"),
					resource.TestCheckResourceAttr("tfcoremock_complex_resource.test", "set.#", "2"),
					resource.TestCheckResourceAttr("tfcoremock_complex_resource.test", "set.0.string", "one"),
					resource.TestCheckResourceAttr("tfcoremock_complex_resource.test", "set.1.string", "zero")),
			},
			{
				Config: LoadFile(t, "testdata/complex/update/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tfcoremock_complex_resource.test", "string", "hello"),
					resource.TestCheckResourceAttr("tfcoremock_complex_resource.test", "list.#", "2"),
					resource.TestCheckResourceAttr("tfcoremock_complex_resource.test", "list.0.string", "zero"),
					resource.TestCheckResourceAttr("tfcoremock_complex_resource.test", "list.1.string", "one"),
					resource.TestCheckResourceAttr("tfcoremock_complex_resource.test", "object.string", "world"),
					resource.TestCheckResourceAttr("tfcoremock_complex_resource.test", "set.#", "1"),
					resource.TestCheckResourceAttr("tfcoremock_complex_resource.test", "set.0.string", "zero")),
			},
			{
				Config: LoadFile(t, "testdata/complex/delete/main.tf"),
			},
		},
	})
}

func TestAccComplexResourceWithBlocks(t *testing.T) {
	t.Cleanup(CleanupTestingDirectories(t))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(""),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/complex_block/create/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tfcoremock_complex_resource.test", "list_block.#", "2"),
					resource.TestCheckResourceAttr("tfcoremock_complex_resource.test", "list_block.0.integer", "0"),
					resource.TestCheckResourceAttr("tfcoremock_complex_resource.test", "list_block.1.integer", "1"),
					resource.TestCheckResourceAttr("tfcoremock_complex_resource.test", "set_block.#", "1"),
					resource.TestCheckResourceAttr("tfcoremock_complex_resource.test", "set_block.0.integer", "0")),
			},
			{
				Config: LoadFile(t, "testdata/complex_block/update/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tfcoremock_complex_resource.test", "list_block.#", "1"),
					resource.TestCheckResourceAttr("tfcoremock_complex_resource.test", "list_block.0.integer", "0"),
					resource.TestCheckResourceAttr("tfcoremock_complex_resource.test", "set_block.#", "2"),
					resource.TestCheckResourceAttr("tfcoremock_complex_resource.test", "set_block.0.integer", "0"),
					resource.TestCheckResourceAttr("tfcoremock_complex_resource.test", "set_block.1.integer", "1")),
			},
			{
				Config: LoadFile(t, "testdata/complex/delete/main.tf"),
			},
		},
	})
}

func TestAccDynamicResource(t *testing.T) {
	t.Cleanup(CleanupTestingDirectories(t))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(LoadFile(t, "testdata/dynamic/dynamic_resources.json")),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/dynamic/create/main.tf"),
				Check:  resource.TestCheckResourceAttr("tfcoremock_dynamic_resource.test", "integer", "0"),
			},
			{
				Config: LoadFile(t, "testdata/dynamic/update/main.tf"),
				Check:  resource.TestCheckResourceAttr("tfcoremock_dynamic_resource.test", "integer", "1"),
			},
			{
				Config: LoadFile(t, "testdata/dynamic/delete/main.tf"),
			},
		},
	})
}

func TestAccDynamicResourceNested(t *testing.T) {
	t.Cleanup(CleanupTestingDirectories(t))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(LoadFile(t, "testdata/dynamic_nested/dynamic_resources.json")),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/dynamic_nested/create/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tfcoremock_dynamic_resource.test", "string", "hello"),
					resource.TestCheckResourceAttr("tfcoremock_dynamic_resource.test", "list.#", "1"),
					resource.TestCheckResourceAttr("tfcoremock_dynamic_resource.test", "list.0.string", "one"),
					resource.TestCheckResourceAttr("tfcoremock_dynamic_resource.test", "metadata_free_list.#", "1"),
					resource.TestCheckResourceAttr("tfcoremock_dynamic_resource.test", "metadata_free_list.0.string", "other"),
					resource.TestCheckResourceAttr("tfcoremock_dynamic_resource.test", "object.bool", "true")),
			},
			{
				Config: LoadFile(t, "testdata/dynamic/delete/main.tf"),
			},
		},
	})
}

func TestAccDynamicResourceWithBlocks(t *testing.T) {
	t.Cleanup(CleanupTestingDirectories(t))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(LoadFile(t, "testdata/dynamic_block/dynamic_resources.json")),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/dynamic_block/create/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tfcoremock_dynamic_resource.test", "integer", "0"),
					resource.TestCheckResourceAttr("tfcoremock_dynamic_resource.test", "nested_list.#", "1"),
					resource.TestCheckResourceAttr("tfcoremock_dynamic_resource.test", "nested_list.0.integer", "0")),
			},
			{
				Config: LoadFile(t, "testdata/dynamic_block/update/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tfcoremock_dynamic_resource.test", "integer", "0"),
					resource.TestCheckResourceAttr("tfcoremock_dynamic_resource.test", "nested_list.#", "2"),
					resource.TestCheckResourceAttr("tfcoremock_dynamic_resource.test", "nested_list.0.integer", "0"),
					resource.TestCheckResourceAttr("tfcoremock_dynamic_resource.test", "nested_list.1.integer", "1")),
			},
			{
				Config: LoadFile(t, "testdata/dynamic/delete/main.tf"),
			},
		},
	})
}

func TestAccDynamicResourceWithComputed(t *testing.T) {
	t.Cleanup(CleanupTestingDirectories(t))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(LoadFile(t, "testdata/dynamic_computed/dynamic_resources.json")),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/dynamic_computed/create/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tfcoremock_dynamic_resource.test", "object_with_value.boolean", "true"),
					resource.TestCheckResourceAttr("tfcoremock_dynamic_resource.test", "object_with_value.string", "hello"),
					resource.TestCheckResourceAttr("tfcoremock_dynamic_resource.test", "computed_list.#", "0"),
					resource.TestCheckResourceAttr("tfcoremock_dynamic_resource.test", "set.#", "3")),
			},
			{
				Config: LoadFile(t, "testdata/dynamic/delete/main.tf"),
			},
		},
	})
}

func TestAccDynamicResourceWithComputedBlocks(t *testing.T) {
	t.Cleanup(CleanupTestingDirectories(t))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(LoadFile(t, "testdata/dynamic_computed_block/dynamic_resources.json")),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/dynamic_computed_block/create/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tfcoremock_dynamic_resource.test", "other.#", "2"),
					resource.TestCheckResourceAttr("tfcoremock_dynamic_resource.test", "other.0.id", "my-id"),
					resource.TestCheckResourceAttr("tfcoremock_dynamic_resource.test", "other.0.nested.#", "2"),
					resource.TestCheckResourceAttr("tfcoremock_dynamic_resource.test", "object.#", "0")),
			},
			{
				Config: LoadFile(t, "testdata/dynamic/delete/main.tf"),
			},
		},
	})
}

func TestAccDynamicResourceWithComputedSetBlocks(t *testing.T) {
	t.Cleanup(CleanupTestingDirectories(t))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(LoadFile(t, "testdata/dynamic_computed_block_set/dynamic_resources.json")),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/dynamic_computed_block_set/create/main.tf"),
			},
			{
				Config: LoadFile(t, "testdata/dynamic/delete/main.tf"),
			},
		},
	})
}

func TestAccDynamicResourceWithId(t *testing.T) {
	t.Cleanup(CleanupTestingDirectories(t))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(LoadFile(t, "testdata/dynamic/dynamic_resources.json")),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/dynamic_with_id/create/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tfcoremock_dynamic_resource.test", "integer", "0"),
					resource.TestCheckResourceAttr("tfcoremock_dynamic_resource.test", "id", "my_id")),
			},
			{
				Config: LoadFile(t, "testdata/dynamic_with_id/update/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tfcoremock_dynamic_resource.test", "integer", "1"),
					resource.TestCheckResourceAttr("tfcoremock_dynamic_resource.test", "id", "my_id")),
			},
			{
				Config: LoadFile(t, "testdata/dynamic/delete/main.tf"),
			},
		},
	})
}

func TestAccDynamicResourceWithRequiresReplace(t *testing.T) {
	t.Cleanup(CleanupTestingDirectories(t))

	var originalId string

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(LoadFile(t, "testdata/dynamic_requires_replace/dynamic_resources.json")),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/dynamic_requires_replace/create/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tfcoremock_list_of_objects.test", "list.0.key", "one"),
					resource.TestCheckResourceAttr("tfcoremock_list_of_objects.test", "list.0.value", "first value"),
					resource.TestCheckResourceAttr("tfcoremock_list_of_objects.test", "list.1.key", "two"),
					resource.TestCheckResourceAttr("tfcoremock_list_of_objects.test", "list.1.value", "second value"),
					SaveResourceId("tfcoremock_list_of_objects.test", &originalId)),
			},
			{
				Config: LoadFile(t, "testdata/dynamic_requires_replace/update/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tfcoremock_list_of_objects.test", "list.0.key", "three"),
					resource.TestCheckResourceAttr("tfcoremock_list_of_objects.test", "list.0.value", "first value"),
					resource.TestCheckResourceAttr("tfcoremock_list_of_objects.test", "list.1.key", "two"),
					resource.TestCheckResourceAttr("tfcoremock_list_of_objects.test", "list.1.value", "second value"),
					CheckResourceIdChanged("tfcoremock_list_of_objects.test", &originalId)),
			},
			{
				Config: LoadFile(t, "testdata/dynamic/delete/main.tf"),
			},
		},
	})
}

func TestAccMultipleDynamicResources(t *testing.T) {
	t.Cleanup(CleanupTestingDirectories(t))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(LoadFile(t, "testdata/multiple_dynamic_resources/dynamic_resources.json")),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/multiple_dynamic_resources/create/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tfcoremock_integer.integer", "integer", "404"),
					resource.TestCheckResourceAttr("tfcoremock_string.string", "string", "Hello, world!")),
			},
			{
				Config: LoadFile(t, "testdata/multiple_dynamic_resources/delete/main.tf"),
			},
		},
	})
}

func TestAccDynamicResourceWithDataSource(t *testing.T) {
	t.Cleanup(CleanupTestingDirectories(t))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(LoadFile(t, "testdata/dynamic_datasource/dynamic_resources.json")),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/dynamic_datasource/create/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.tfcoremock_simple_resource.test", "integer", "0"),
					resource.TestCheckResourceAttr("tfcoremock_dynamic_resource.test", "id", "my_dynamic_resource"),
					resource.TestCheckResourceAttr("tfcoremock_dynamic_resource.test", "my_value", "0")),
			},
			{
				Config: LoadFile(t, "testdata/dynamic_datasource/delete/main.tf"),
			},
		},
	})
}

func TestAccSimpleDataSource(t *testing.T) {
	t.Cleanup(CleanupTestingDirectories(t))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(""),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/simple_datasource/get/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.tfcoremock_simple_resource.data", "integer", "0"),
					resource.TestCheckResourceAttr("tfcoremock_simple_resource.test", "integer", "0")),
			},
		},
	})
}

func TestAccSimpleResource(t *testing.T) {
	t.Cleanup(CleanupTestingDirectories(t))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(""),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/simple/create/main.tf"),
				Check:  resource.TestCheckResourceAttr("tfcoremock_simple_resource.test", "integer", "0"),
			},
			{
				Config: LoadFile(t, "testdata/simple/update/main.tf"),
				Check:  resource.TestCheckResourceAttr("tfcoremock_simple_resource.test", "integer", "1"),
			},
			{
				Config: LoadFile(t, "testdata/simple/delete/main.tf"),
			},
		},
	})
}

func TestAccSimpleResourceWithDrift(t *testing.T) {
	t.Cleanup(CleanupTestingDirectories(t))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(""),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/simple/create/main.tf"),
				Check: func(state *terraform.State) error {
					id := state.RootModule().Resources["tfcoremock_simple_resource.test"].Primary.Attributes["id"]
					return os.Remove(fmt.Sprintf("terraform.resource/%s.json", id))
				},
				ExpectNonEmptyPlan: true,
			},
			{
				Config: LoadFile(t, "testdata/simple/update/main.tf"),
				Check:  resource.TestCheckResourceAttr("tfcoremock_simple_resource.test", "integer", "1"),
			},
			{
				Config: LoadFile(t, "testdata/simple/delete/main.tf"),
			},
		},
	})
}

func TestAccSimpleResourceWithId(t *testing.T) {
	t.Cleanup(CleanupTestingDirectories(t))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(""),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/simple_with_id/create/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tfcoremock_simple_resource.test", "string", "hello"),
					resource.TestCheckResourceAttr("tfcoremock_simple_resource.test", "id", "my_id")),
			},
			{
				Config: LoadFile(t, "testdata/simple_with_id/update/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tfcoremock_simple_resource.test", "string", "world"),
					resource.TestCheckResourceAttr("tfcoremock_simple_resource.test", "id", "my_id")),
			},
			{
				Config: LoadFile(t, "testdata/simple/delete/main.tf"),
			},
		},
	})
}

func TestAccSimpleResourceWithDependsOn(t *testing.T) {
	t.Cleanup(CleanupTestingDirectories(t))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(""),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/depends_on/create/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tfcoremock_simple_resource.example_one", "string", "resource_module1"),
					resource.TestCheckResourceAttr("tfcoremock_simple_resource.example_two", "string", "data")),
			},
		},
	})
}

func TestAccSimpleResourceFailsOnRead(t *testing.T) {
	t.Cleanup(CleanupTestingDirectories(t))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(""),
		Steps: []resource.TestStep{
			{
				Config:      LoadFile(t, "testdata/fail_on/read/main.tf"),
				ExpectError: regexp.MustCompile("forced failure"),
			},
		},
	})
}

func TestAccSimpleResourceFailsOnUpdate(t *testing.T) {
	t.Cleanup(CleanupTestingDirectories(t))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(""),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/fail_on/update/create.tf"),
			},
			{
				Config:      LoadFile(t, "testdata/fail_on/update/update.tf"),
				ExpectError: regexp.MustCompile("forced failure"),
			},
		},
	})
}

func TestAccSimpleResourceFailsOnCreate(t *testing.T) {
	t.Cleanup(CleanupTestingDirectories(t))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(""),
		Steps: []resource.TestStep{
			{
				Config:      LoadFile(t, "testdata/fail_on/create/main.tf"),
				ExpectError: regexp.MustCompile("forced failure"),
			},
		},
	})
}

func TestAccSimpleResourceFailsOnDelete(t *testing.T) {
	t.Cleanup(CleanupTestingDirectories(t))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(""),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/fail_on/delete/create.tf"),
			},
			{
				Config:      LoadFile(t, "testdata/fail_on/delete/delete.tf"),
				ExpectError: regexp.MustCompile("forced failure"),
			},
			{
				// We need to update the provider configuration to remove the
				// failing resource so that the test can complete.
				Config: LoadFile(t, "testdata/simple/delete/main.tf"),
			},
		},
	})
}
