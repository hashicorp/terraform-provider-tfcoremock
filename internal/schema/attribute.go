// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"fmt"

	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/pkg/errors"

	"github.com/hashicorp/terraform-provider-tfcoremock/internal/data"
)

// Attribute defines an internal representation of a Terraform attribute in a
// schema.
//
// It is designed to be read dynamically from a JSON object, allowing schemas,
// blocks and attributes to be defined dynamically by the user of the provider.
type Attribute struct {
	Description         string `json:"-"` // Dynamic resources don't need descriptions so hide them from the exposed JSON schema.
	MarkdownDescription string `json:"-"` // Dynamic resources don't need descriptions so hide them from the exposed JSON schema.

	Type     Type `json:"type"`
	Optional bool `json:"optional"`
	Required bool `json:"required"`
	Computed bool `json:"computed"`

	Value *data.Value `json:"value,omitempty"`

	List   *Attribute           `json:"list,omitempty"`
	Map    *Attribute           `json:"map,omitempty"`
	Object map[string]Attribute `json:"object,omitempty"`
	Set    *Attribute           `json:"set,omitempty"`

	Sensitive bool `json:"sensitive"` // True if values for this attribute should be hidden in the plan.
	Replace   bool `json:"replace"`   // True if the resource should be replaced when this attribute changes.

	// SkipNestedMetadata instructs the dynamic resource to not use the nested
	// attribute field when building element and attribute types of complex
	// attributes (list, map, object, and set).
	SkipNestedMetadata bool `json:"skip_nested_metadata"`
}

// AttributeTypes contains functions that map provider attributes into Terraform
// resource or datasource attributes.
type AttributeTypes[A any] struct {
	asBoolean func(attribute Attribute) (*A, error)
	asFloat   func(attribute Attribute) (*A, error)
	asInteger func(attribute Attribute) (*A, error)
	asNumber  func(attribute Attribute) (*A, error)
	asString  func(attribute Attribute) (*A, error)

	asList       func(attribute Attribute) (*A, error)
	asNestedList func(attribute Attribute) (*A, error)

	asMap       func(attribute Attribute) (*A, error)
	asNestedMap func(attribute Attribute) (*A, error)

	asSet       func(attribute Attribute) (*A, error)
	asNestedSet func(attribute Attribute) (*A, error)

	asObject       func(attribute Attribute) (*A, error)
	asNestedObject func(attribute Attribute) (*A, error)
}

// ToTerraformAttribute converts our representation of an Attribute into a
// Terraform SDK attribute so it can be passed back to Terraform Core in a
// resource or data source schema.
func ToTerraformAttribute[A any](a Attribute, types *AttributeTypes[A]) (*A, error) {
	switch a.Type {
	case Boolean:
		return types.asBoolean(a)
	case Float:
		return types.asFloat(a)
	case Integer:
		return types.asInteger(a)
	case Number:
		return types.asNumber(a)
	case String:
		return types.asString(a)
	case List:
		if !a.SkipNestedMetadata && a.List.Type == Object {
			return types.asNestedList(a)
		}
		return types.asList(a)
	case Map:
		if !a.SkipNestedMetadata && a.Map.Type == Object {
			return types.asNestedMap(a)
		}
		return types.asMap(a)
	case Set:
		if !a.SkipNestedMetadata && a.Set.Type == Object {
			return types.asNestedSet(a)
		}
		return types.asNestedSet(a)
	case Object:
		if a.SkipNestedMetadata {
			return types.asObject(a)
		}

		return types.asNestedObject(a)
	case "":
		return nil, fmt.Errorf("missing attribute type")
	default:
		return nil, fmt.Errorf("unrecognized attribute type '%s'", a.Type)
	}
}

func attributesToTerraformResourceAttributes(attributes map[string]Attribute) (map[string]resource_schema.Attribute, error) {
	tfAttributes := make(map[string]resource_schema.Attribute)
	for name, attribute := range attributes {
		attribute, err := ToTerraformAttribute(attribute, resources)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create attribute '%s'", name)
		}
		tfAttributes[name] = *attribute
	}
	return tfAttributes, nil
}

func attributesToTerraformDataSourceAttributes(attributes map[string]Attribute) (map[string]datasource_schema.Attribute, error) {
	tfAttributes := make(map[string]datasource_schema.Attribute)
	for name, attribute := range attributes {
		attribute, err := ToTerraformAttribute(attribute, resources)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create attribute '%s'", name)
		}
		tfAttributes[name] = *attribute
	}
	return tfAttributes, nil
}
