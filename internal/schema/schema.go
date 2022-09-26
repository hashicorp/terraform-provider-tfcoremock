package schema

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Schema defines an internal representation of a Terraform schema.
//
// It is designed to be read dynamically from a JSON object, allowing schemas,
// blocks and attributes to be defined dynamically by the user of the provider.
type Schema struct {
	Attributes map[string]Attribute `json:"attributes"`
	Blocks     map[string]Block     `json:"blocks"`
}

// AllAttributes returns the attributes for the dynamic schema, plus the
// required ID attribute that is attached to tfsdk.Schema objects automatically.
func (schema Schema) AllAttributes() map[string]Attribute {
	attributes := make(map[string]Attribute, 0)
	for key, attribute := range schema.Attributes {
		attributes[key] = attribute
	}
	attributes["id"] = Attribute{
		Type:     String,
		Optional: false,
		Required: false,
		Computed: true,
	}
	return attributes
}

// ToTerraformResourceSchema converts out representation of a Schema into a
// Terraform SDK tfsdk.Schema. It automatically creates and attaches a computed
// type called `id` that is required by every resource and data source in this
// provider.
func (schema Schema) ToTerraformResourceSchema() (tfsdk.Schema, error) {
	out := tfsdk.Schema{}

	var err error
	if out.Attributes, err = schema.getTerraformAttributes(); err != nil {
		return out, err
	}

	out.Attributes["id"] = tfsdk.Attribute{
		Required: false,
		Optional: true,
		Computed: true,
		PlanModifiers: tfsdk.AttributePlanModifiers{
			resource.UseStateForUnknown(),
			resource.RequiresReplace(),
		},
		Type: types.StringType,
	}

	if out.Blocks, err = schema.getTerraformBlocks(); err != nil {
		return out, err
	}

	return out, nil
}

// ToTerraformDataSourceSchema converts our representation of a Schema into a
// Terraform SDK tfsdk.Schema. It automatically creates and attaches a required
// attribute called `id` that is required by every resource and data source in
// this provider.
func (schema Schema) ToTerraformDataSourceSchema() (tfsdk.Schema, error) {
	out := tfsdk.Schema{}

	var err error
	if out.Attributes, err = schema.getTerraformAttributes(); err != nil {
		return out, err
	}

	out.Attributes["id"] = tfsdk.Attribute{
		Required: true,
		Optional: false,
		Computed: false,
		PlanModifiers: tfsdk.AttributePlanModifiers{
			resource.UseStateForUnknown(),
			resource.RequiresReplace(),
		},
		Type: types.StringType,
	}

	if out.Blocks, err = schema.getTerraformBlocks(); err != nil {
		return out, err
	}

	return out, nil
}

func (schema Schema) getTerraformAttributes() (map[string]tfsdk.Attribute, error) {
	if _, ok := schema.Attributes["id"]; ok {
		return nil, errors.New("top level dynamic objects cannot define a value called `id` as the provider will generate an identifier for them")
	}
	return attributesToTerraformAttributes(schema.Attributes)
}

func (schema Schema) getTerraformBlocks() (map[string]tfsdk.Block, error) {
	return blocksToTerraformBlocks(schema.Blocks)
}
