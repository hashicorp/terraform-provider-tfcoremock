package schema

import (
	"errors"

	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	resource_schema_planmodifier "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	resource_schema_stringplanmodifier "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// Schema defines an internal representation of a Terraform schema.
//
// It is designed to be read dynamically from a JSON object, allowing schemas,
// blocks and attributes to be defined dynamically by the user of the provider.
type Schema struct {
	Description         string               `json:"-"` // Dynamic resources don't need descriptions so hide them from the exposed JSON schema.
	MarkdownDescription string               `json:"-"` // Dynamic resources don't need descriptions so hide them from the exposed JSON schema.
	Attributes          map[string]Attribute `json:"attributes"`
	Blocks              map[string]Block     `json:"blocks"`
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
func (schema Schema) ToTerraformResourceSchema() (resource_schema.Schema, error) {
	out := resource_schema.Schema{
		Description:         schema.Description,
		MarkdownDescription: schema.MarkdownDescription,
	}

	var err error
	if err = schema.validateAttributes(); err != nil {
		return out, err
	}

	if out.Attributes, err = attributesToTerraformResourceAttributes(schema.Attributes); err != nil {
		return out, err
	}
	out.Attributes["id"] = resource_schema.StringAttribute{
		Required: false,
		Optional: true,
		Computed: true,
		PlanModifiers: []resource_schema_planmodifier.String{
			resource_schema_stringplanmodifier.UseStateForUnknown(),
			resource_schema_stringplanmodifier.RequiresReplace(),
		},
	}

	if out.Blocks, err = blocksToTerraformResourceBlocks(schema.Blocks); err != nil {
		return out, err
	}

	return out, nil
}

// ToTerraformDataSourceSchema converts our representation of a Schema into a
// Terraform SDK tfsdk.Schema. It automatically creates and attaches a required
// attribute called `id` that is required by every resource and data source in
// this provider.
func (schema Schema) ToTerraformDataSourceSchema() (datasource_schema.Schema, error) {
	out := datasource_schema.Schema{
		Description:         schema.Description,
		MarkdownDescription: schema.MarkdownDescription,
	}

	var err error
	if err = schema.validateAttributes(); err != nil {
		return out, err
	}

	if out.Attributes, err = attributesToTerraformDataSourceAttributes(schema.Attributes); err != nil {
		return out, err
	}

	out.Attributes["id"] = datasource_schema.StringAttribute{
		Required: true,
		Optional: false,
		Computed: false,
	}

	if out.Blocks, err = blocksToTerraformDataSourceBlocks(schema.Blocks); err != nil {
		return out, err
	}

	return out, nil
}

func (schema Schema) validateAttributes() error {
	if _, ok := schema.Attributes["id"]; ok {
		return errors.New("top level dynamic objects cannot define a value called `id` as the provider will generate an identifier for them")
	}
	return nil
}
