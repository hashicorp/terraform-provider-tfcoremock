package computed

import (
	"errors"
	"math/big"

	"github.com/hashicorp/go-uuid"

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
		case "", schema.NestingModeSingle:
			next := []data.Value{(*values)[key]}
			err = generateComputedValuesForBlock(&next, block)
		case schema.NestingModeSet:
			err = generateComputedValuesForBlock((*values)[key].Set, block)
		case schema.NestingModeList:
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

		// If we get here, then it means we don't have a value for the target
		// attribute so we're going to have to get one from somewhere.

		if attribute.Value != nil {
			// The user has provided us with a value, so this is very easy.
			(*values)[key] = *attribute.Value
			continue
		}

		if !attribute.Required && !attribute.Computed {
			// We don't actually need to provide anything, so let's just skip
			// over this.
			//
			// Note, that we are actually creating values for required
			// attributes within computed objects at this point. We are relying
			// on Terraform Core and the SDK to properly validate the config, so
			// we do trust that any required but non-computed attributes have
			// been properly set in the configuration.
			continue
		}

		// Finally, we are in the unhappy position of having to generate a
		// value.

		switch attribute.Type {
		case schema.Boolean:
			value := false
			(*values)[key] = data.Value{
				Boolean: &value,
			}
		case schema.Float, schema.Integer, schema.Number:
			(*values)[key] = data.Value{
				Number: big.NewFloat(0),
			}
		case schema.String:
			value, err := uuid.GenerateUUID()
			if err != nil {
				return err
			}
			(*values)[key] = data.Value{
				String: &value,
			}
		case schema.List:
			(*values)[key] = data.Value{
				List: &[]data.Value{},
			}
		case schema.Map:
			(*values)[key] = data.Value{
				Map: &map[string]data.Value{},
			}
		case schema.Set:
			(*values)[key] = data.Value{
				Set: &[]data.Value{},
			}
		case schema.Object:
			// The object is the only tricky one, as we can't just set it as
			// empty.
			childValues := make(map[string]data.Value)
			if err := generateComputedValuesForObject(&childValues, attribute.Object); err != nil {
				return err
			}
			(*values)[key] = data.Value{
				Object: &childValues,
			}
		default:
			return errors.New("unrecognized attribute type: " + string(attribute.Type))
		}
	}

	return nil
}
