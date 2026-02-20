// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package computed

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-provider-tfcoremock/internal/data"
	"github.com/hashicorp/terraform-provider-tfcoremock/internal/schema"
)

// GenerateComputedValues steps through the resource and uses the schema to
// populate any computed values.
//
// Computed values have a sensible default for all primitive types, and can be
// specified using a data.Value object as part of the dynamic schema.
//
// Objects are complicated as you can have nested objects with required values
// so the default value for a computed object is to generate an object with all
// the required and computed values populated using a default.
func GenerateComputedValues(resource *data.Resource, schema schema.Schema) error {
	if err := generateComputedValuesForObject(&resource.Values, schema.AllAttributes()); err != nil {
		return err
	}

	if err := generateComputedValuesForBlocks(&resource.Values, schema.Blocks); err != nil {
		return err
	}

	return nil
}

func generateComputedValuesForBlocks(values *map[string]data.Value, blocks map[string]schema.Block) error {
	for key, block := range blocks {
		var err error
		switch block.Mode {
		case schema.NestingModeSet:
			err = generateComputedValuesForBlock((*values)[key].Set, block)
		case "", schema.NestingModeList:
			err = generateComputedValuesForBlock((*values)[key].List, block)
		default:
			return errors.New("unrecognized block type: " + block.Mode)
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func generateComputedValuesForBlock(values *[]data.Value, block schema.Block) error {
	if values == nil {
		return nil
	}

	for ix, value := range *values {
		if err := generateComputedValuesForObject(value.Object, block.Attributes); err != nil {
			return err
		}

		if err := generateComputedValuesForBlocks(value.Object, block.Blocks); err != nil {
			return err
		}

		(*values)[ix] = value
	}

	return nil
}

func generateComputedValue(value data.Value, attribute *schema.Attribute) (data.Value, error) {
	var err error
	switch attribute.Type {
	case schema.Boolean, schema.Float, schema.Integer, schema.Number, schema.String:
		// For these types we don't need to do anything, they have a value
		// set and we're all good to leave them as is.
	case schema.List:
		err = generateComputedValuesForList(value.List, attribute.List)
	case schema.Set:
		err = generateComputedValuesForSet(value.Set, attribute.Set)
	case schema.Map:
		err = generateComputedValuesForMap(value.Map, attribute.Map)
	case schema.Object:
		err = generateComputedValuesForObject(value.Object, attribute.Object)
	default:
		return value, errors.New("unrecognized attribute type: " + string(attribute.Type))
	}

	return value, err
}

func generateComputedValuesForList(values *[]data.Value, attribute *schema.Attribute) error {
	for ix, value := range *values {
		// Then we're going to go through each value and check if it has any
		// attributes that need to be computed.
		newValue, err := generateComputedValue(value, attribute)
		if err != nil {
			return err
		}
		(*values)[ix] = newValue
	}
	return nil
}

func generateComputedValuesForSet(values *[]data.Value, attribute *schema.Attribute) error {
	for ix, value := range *values {
		// Then we're going to go through each value and check if it has any
		// attributes that need to be computed.
		newValue, err := generateComputedValue(value, attribute)
		if err != nil {
			return err
		}
		(*values)[ix] = newValue
	}
	return nil
}

func generateComputedValuesForMap(values *map[string]data.Value, attribute *schema.Attribute) error {
	for key, value := range *values {
		// Then we're going to go through each value and check if it has any
		// attributes that need to be computed.
		newValue, err := generateComputedValue(value, attribute)
		if err != nil {
			return err
		}
		(*values)[key] = newValue
	}
	return nil
}

func generateComputedValuesForObject(values *map[string]data.Value, attributes map[string]schema.Attribute) error {
	for key, attribute := range attributes {
		if value, ok := (*values)[key]; ok {
			// This means we already have a value for this attribute, so we're
			// not going to generate a new one completely. But we do need to
			// recurse down into any objects as they maybe have generated
			// attributes.
			var err error
			if (*values)[key], err = generateComputedValue(value, &attribute); err != nil {
				return err
			}
			continue
		}

		if attribute.Value != nil {
			if !attribute.Computed {
				// If we didn't check this, it would just cause another error
				// later but at least here we can return a nice error message.
				return fmt.Errorf("attribute %s has specified a value in the json schema without being marked as computed", key)
			}
			var err error
			if (*values)[key], err = generateComputedValue(*attribute.Value, &attribute); err != nil {
				return err
			}
		}
	}

	return nil
}
