// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package data

import (
	"errors"
	"math/big"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Value is the mock provider's representation of any generic Terraform Value.
//
// It can be converted from/to a tftypes.Value using the functions in this
// package, and it can be marshalled to/from JSON using the Golang JSON package.
//
// Only a single field in the struct will be set at a given time.
//
// We use pointers where appropriate to make sure the omitempty metadata works
// to keep the produced structs as small and relevant as possible as they are
// intended to be consumed by humans.
//
// We introduce pointers to the complex objects because there is a difference
// between unset (or nil) and an empty list and we want to record that
// difference.
type Value struct {
	Boolean *bool      `json:"boolean,omitempty"`
	Number  *big.Float `json:"number,omitempty"`
	String  *string    `json:"string,omitempty"`

	List   *[]Value          `json:"list,omitempty"`
	Map    *map[string]Value `json:"map,omitempty"`
	Object *map[string]Value `json:"object,omitempty"`
	Set    *[]Value          `json:"set,omitempty"`
}

// ToTerraform5Value accepts our representation of a Value alongside the
// tftypes.Type description, and returns a tftypes.Value object that can be
// passed into the Terraform SDK.
func ToTerraform5Value(v Value, t tftypes.Type) (tftypes.Value, error) {
	switch {
	case t.Is(tftypes.Bool):
		return tftypes.NewValue(tftypes.Bool, v.Boolean), nil
	case t.Is(tftypes.String):
		return tftypes.NewValue(tftypes.String, v.String), nil
	case t.Is(tftypes.Number):
		return tftypes.NewValue(tftypes.Number, v.Number), nil
	case t.Is(tftypes.List{}):
		return listToTerraform5Value(v.List, t.(tftypes.List))
	case t.Is(tftypes.Map{}):
		return mapToTerraform5Value(v.Map, t.(tftypes.Map))
	case t.Is(tftypes.Object{}):
		object, err := objectToTerraform5Value(v.Object, t.(tftypes.Object))
		if err != nil {
			return tftypes.Value{}, err
		}
		return tftypes.NewValue(t, object), nil
	case t.Is(tftypes.Set{}):
		return setToTerraform5Value(v.Set, t.(tftypes.Set))
	default:
		return tftypes.Value{}, errors.New("Unrecognized type: " + t.String())
	}
}

// FromTerraform5Value accepts a tftypes.Value and returns our representation
// of a Value.
//
// Note, that unlike the reverse ToTerraform5Value function we do not need to
// include the type information as this is not embedded in our representation of
// the type (the expectation is that the type information will always be
// provided by the SDK regardless of which direction we need to go).
func FromTerraform5Value(v tftypes.Value) (Value, error) {
	t := v.Type()
	switch {
	case t.Is(tftypes.Bool):
		ret := Value{}
		err := v.As(&ret.Boolean)
		return ret, err
	case t.Is(tftypes.String):
		ret := Value{}
		err := v.As(&ret.String)
		return ret, err
	case t.Is(tftypes.Number):
		ret := Value{}
		err := v.As(&ret.Number)
		return ret, err
	case t.Is(tftypes.List{}):
		return listFromTerraform5Value(v)
	case t.Is(tftypes.Map{}):
		return mapFromTerraform5Value(v)
	case t.Is(tftypes.Object{}):
		return objectFromTerraform5Value(v)
	case t.Is(tftypes.Set{}):
		return setFromTerraform5Value(v)
	default:
		return Value{}, errors.New("Unrecognized type: " + t.String())
	}
}

func listToTerraform5Value(values *[]Value, listType tftypes.List) (tftypes.Value, error) {
	if values == nil {
		return tftypes.NewValue(listType, nil), nil
	}

	children := make([]tftypes.Value, 0)
	for _, value := range *values {
		child, err := ToTerraform5Value(value, listType.ElementType)
		if err != nil {
			return tftypes.Value{}, err
		}
		children = append(children, child)
	}
	return tftypes.NewValue(listType, children), nil
}

func mapToTerraform5Value(values *map[string]Value, mapType tftypes.Map) (tftypes.Value, error) {
	if values == nil {
		return tftypes.NewValue(mapType, nil), nil
	}

	children := make(map[string]tftypes.Value)
	for name, value := range *values {
		child, err := ToTerraform5Value(value, mapType.ElementType)
		if err != nil {
			return tftypes.Value{}, err
		}
		children[name] = child
	}
	return tftypes.NewValue(mapType, children), nil
}

// objectToTerraform5Value is a bit of a special case as we return the inner
// value of the value instead of the value directly (as with the other
// functions). This is because we use this function as part of our
// implementation of the ValueCreator and ValueConverter of the Resource type
// which expects the underlying structure instead of being already converted.
func objectToTerraform5Value(values *map[string]Value, objectType tftypes.Object) (interface{}, error) {
	if values == nil {
		return nil, nil
	}

	children := make(map[string]tftypes.Value)
	for name, childType := range objectType.AttributeTypes {

		// It is possible that this child type exists in the type representation
		// but not in the actual value (this is because we can have optional
		// attributes in objects). So we try and retrieve the child from the
		// values but if it is not there we don't fail, instead we just set an
		// empty value in its place.

		var err error
		if value, ok := (*values)[name]; ok {
			if children[name], err = ToTerraform5Value(value, childType); err != nil {
				return nil, err
			}
			continue
		}

		// Otherwise we just set a nil value.
		children[name] = tftypes.NewValue(childType, nil)
	}
	return children, nil
}

func setToTerraform5Value(values *[]Value, setType tftypes.Set) (tftypes.Value, error) {
	if values == nil {
		return tftypes.NewValue(setType, nil), nil
	}

	children := make([]tftypes.Value, 0)
	for _, value := range *values {
		child, err := ToTerraform5Value(value, setType.ElementType)
		if err != nil {
			return tftypes.Value{}, err
		}
		children = append(children, child)
	}
	return tftypes.NewValue(setType, children), nil
}

func listFromTerraform5Value(v tftypes.Value) (Value, error) {
	var children []tftypes.Value
	if err := v.As(&children); err != nil {
		return Value{}, err
	}

	// There is a difference between a list being empty and being null from
	// Terraform's perspective. So we want to create a list of length 0 rather
	// than leaving it as null.
	list := make([]Value, 0)
	for _, child := range children {
		parsed, err := FromTerraform5Value(child)
		if err != nil {
			return Value{}, err
		}
		list = append(list, parsed)
	}

	return Value{
		List: &list,
	}, nil
}

func mapFromTerraform5Value(v tftypes.Value) (Value, error) {
	var children map[string]tftypes.Value
	if err := v.As(&children); err != nil {
		return Value{}, err
	}

	values := make(map[string]Value)
	for name, child := range children {
		parsed, err := FromTerraform5Value(child)
		if err != nil {
			return Value{}, err
		}
		values[name] = parsed
	}

	return Value{
		Map: &values,
	}, nil
}

func objectFromTerraform5Value(v tftypes.Value) (Value, error) {
	var children map[string]tftypes.Value
	if err := v.As(&children); err != nil {
		return Value{}, err
	}

	values := make(map[string]Value)
	for name, child := range children {
		if child.IsNull() || !child.IsKnown() {
			// Terraform handles unset objects differently to us. We just don't
			// add unset attributes to our objects while terraform adds them
			// but sets them to null. If this child value is null in the
			// Terraform representation we just skip it.
			//
			// Note, the reverse implementation in objectToTerrafrom5Value. We
			// check the type information and set any missing attributes as null
			// when converting into the terraform representation.
			//
			// For now, we also treat unknown values the same as null by just
			// skipping them. Any computed values will be filled in later.
			continue
		}

		parsed, err := FromTerraform5Value(child)
		if err != nil {
			return Value{}, err
		}
		values[name] = parsed

	}

	return Value{
		Object: &values,
	}, nil
}

func setFromTerraform5Value(v tftypes.Value) (Value, error) {
	var children []tftypes.Value
	if err := v.As(&children); err != nil {
		return Value{}, err
	}

	set := make([]Value, 0)
	for _, child := range children {
		parsed, err := FromTerraform5Value(child)
		if err != nil {
			return Value{}, err
		}
		set = append(set, parsed)
	}

	return Value{
		Set: &set,
	}, nil
}
