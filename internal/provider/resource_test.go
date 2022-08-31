package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccComplexResource(t *testing.T) {
	defer CleanupTestingDirectories(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(""),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/complex/create/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mock_complex_resource.test", "string", "hello"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "list.#", "1"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "list.0.string", "one"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "object.bool", "true"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "set.#", "2"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "set.0.string", "one"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "set.1.string", "zero")),
			},
			{
				Config: LoadFile(t, "testdata/complex/update/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mock_complex_resource.test", "string", "hello"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "list.#", "2"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "list.0.string", "zero"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "list.1.string", "one"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "object.string", "world"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "set.#", "1"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "set.0.string", "zero")),
			},
			{
				Config: LoadFile(t, "testdata/complex/delete/main.tf"),
			},
		},
	})
}

func TestAccComplexResourceWithBlocks(t *testing.T) {
	defer CleanupTestingDirectories(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(""),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/complex_block/create/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mock_complex_resource.test", "list_block.#", "2"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "list_block.0.integer", "0"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "list_block.1.integer", "1"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "set_block.#", "1"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "set_block.0.integer", "0")),
			},
			{
				Config: LoadFile(t, "testdata/complex_block/update/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mock_complex_resource.test", "list_block.#", "1"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "list_block.0.integer", "0"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "set_block.#", "2"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "set_block.0.integer", "0"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "set_block.1.integer", "1")),
			},
			{
				Config: LoadFile(t, "testdata/complex/delete/main.tf"),
			},
		},
	})
}

func TestAccDynamicResource(t *testing.T) {
	defer CleanupTestingDirectories(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(LoadFile(t, "testdata/dynamic/dynamic_resources.json")),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/dynamic/create/main.tf"),
				Check:  resource.TestCheckResourceAttr("mock_dynamic_resource.test", "integer", "0"),
			},
			{
				Config: LoadFile(t, "testdata/dynamic/update/main.tf"),
				Check:  resource.TestCheckResourceAttr("mock_dynamic_resource.test", "integer", "1"),
			},
			{
				Config: LoadFile(t, "testdata/dynamic/delete/main.tf"),
			},
		},
	})
}

func TestAccDynamicResourceNested(t *testing.T) {
	defer CleanupTestingDirectories(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(LoadFile(t, "testdata/dynamic_nested/dynamic_resources.json")),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/dynamic_nested/create/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mock_dynamic_resource.test", "string", "hello"),
					resource.TestCheckResourceAttr("mock_dynamic_resource.test", "list.#", "1"),
					resource.TestCheckResourceAttr("mock_dynamic_resource.test", "list.0.string", "one"),
					resource.TestCheckResourceAttr("mock_dynamic_resource.test", "object.bool", "true")),
			},
			{
				Config: LoadFile(t, "testdata/dynamic/delete/main.tf"),
			},
		},
	})
}

func TestAccDynamicResourceWithBlocks(t *testing.T) {
	defer CleanupTestingDirectories(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(LoadFile(t, "testdata/dynamic_block/dynamic_resources.json")),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/dynamic_block/create/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mock_dynamic_resource.test", "integer", "0"),
					resource.TestCheckResourceAttr("mock_dynamic_resource.test", "nested_list.#", "1"),
					resource.TestCheckResourceAttr("mock_dynamic_resource.test", "nested_list.0.integer", "0")),
			},
			{
				Config: LoadFile(t, "testdata/dynamic_block/update/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mock_dynamic_resource.test", "integer", "0"),
					resource.TestCheckResourceAttr("mock_dynamic_resource.test", "nested_list.#", "2"),
					resource.TestCheckResourceAttr("mock_dynamic_resource.test", "nested_list.0.integer", "0"),
					resource.TestCheckResourceAttr("mock_dynamic_resource.test", "nested_list.1.integer", "1")),
			},
			{
				Config: LoadFile(t, "testdata/dynamic/delete/main.tf"),
			},
		},
	})
}

func TestAccDynamicResourceWithId(t *testing.T) {
	defer CleanupTestingDirectories(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(LoadFile(t, "testdata/dynamic/dynamic_resources.json")),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/dynamic_with_id/create/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mock_dynamic_resource.test", "integer", "0"),
					resource.TestCheckResourceAttr("mock_dynamic_resource.test", "id", "my_id")),
			},
			{
				Config: LoadFile(t, "testdata/dynamic_with_id/update/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mock_dynamic_resource.test", "integer", "1"),
					resource.TestCheckResourceAttr("mock_dynamic_resource.test", "id", "my_id")),
			},
			{
				Config: LoadFile(t, "testdata/dynamic/delete/main.tf"),
			},
		},
	})
}

func TestAccSimpleDataSource(t *testing.T) {
	defer CleanupTestingDirectories(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(""),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/simple_datasource/get/main.tf"),
				Check:  resource.TestCheckResourceAttr("data.mock_simple_resource.test", "integer", "0"),
			},
		},
	})
}

func TestAccSimpleResource(t *testing.T) {
	defer CleanupTestingDirectories(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(""),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/simple/create/main.tf"),
				Check:  resource.TestCheckResourceAttr("mock_simple_resource.test", "integer", "0"),
			},
			{
				Config: LoadFile(t, "testdata/simple/update/main.tf"),
				Check:  resource.TestCheckResourceAttr("mock_simple_resource.test", "integer", "1"),
			},
			{
				Config: LoadFile(t, "testdata/simple/delete/main.tf"),
			},
		},
	})
}

func TestAccSimpleResourceWithDrift(t *testing.T) {
	defer CleanupTestingDirectories(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(""),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/simple/create/main.tf"),
				Check: func(state *terraform.State) error {
					id := state.RootModule().Resources["mock_simple_resource.test"].Primary.Attributes["id"]
					return os.Remove(fmt.Sprintf("terraform.resource/%s.json", id))
				},
				ExpectNonEmptyPlan: true,
			},
			{
				Config: LoadFile(t, "testdata/simple/update/main.tf"),
				Check:  resource.TestCheckResourceAttr("mock_simple_resource.test", "integer", "1"),
			},
			{
				Config: LoadFile(t, "testdata/simple/delete/main.tf"),
			},
		},
	})
}

func TestAccSimpleResourceWithId(t *testing.T) {
	defer CleanupTestingDirectories(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(""),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/simple_with_id/create/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mock_simple_resource.test", "string", "hello"),
					resource.TestCheckResourceAttr("mock_simple_resource.test", "id", "my_id")),
			},
			{
				Config: LoadFile(t, "testdata/simple_with_id/update/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mock_simple_resource.test", "string", "world"),
					resource.TestCheckResourceAttr("mock_simple_resource.test", "id", "my_id")),
			},
			{
				Config: LoadFile(t, "testdata/simple/delete/main.tf"),
			},
		},
	})
}
