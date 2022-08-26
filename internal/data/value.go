package data

import (
	"errors"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"math/big"
)

// Value is the mock provider's representation of any generic Terraform value.
//
// It can be converted from/to a tftypes.Value using the functions in this
// package, and it can be marshalled to/from JSON using the Golang JSON package.
//
// Only a single field in the struct will be set at a given time.
//
// We use pointers where appropriate to make sure the omitempty metadata works
// to keep the produced structs as small and relevant as possible as they are
// intended to be consumed by humans.
type Value struct {
	Boolean *bool      `json:"boolean,omitempty"`
	Number  *big.Float `json:"number,omitempty"`
	String  *string    `json:"string,omitempty"`

	List   []Value          `json:"list,omitempty"`
	Map    map[string]Value `json:"map,omitempty"`
	Object map[string]Value `json:"object,omitempty"`
	Set    []Value          `json:"set,omitempty"`
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
func FromTerraform5Value(value tftypes.Value) (Value, error) {
	t := value.Type()
	switch {
	case t.Is(tftypes.Bool):
		v := Value{}
		err := value.As(&v.Boolean)
		return v, err
	case t.Is(tftypes.String):
		v := Value{}
		err := value.As(&v.String)
		return v, err
	case t.Is(tftypes.Number):
		v := Value{}
		err := value.As(&v.Number)
		return v, err
	case t.Is(tftypes.List{}):
		return listFromTerraform5Value(value)
	case t.Is(tftypes.Map{}):
		return mapFromTerraform5Value(value)
	case t.Is(tftypes.Object{}):
		return objectFromTerraform5Value(value)
	case t.Is(tftypes.Set{}):
		return setFromTerraform5Value(value)
	default:
		return Value{}, errors.New("Unrecognized type: " + t.String())
	}
}

func listToTerraform5Value(values []Value, listType tftypes.List) (tftypes.Value, error) {
	var children []tftypes.Value
	for _, value := range values {
		child, err := ToTerraform5Value(value, listType.ElementType)
		if err != nil {
			return tftypes.Value{}, err
		}
		children = append(children, child)
	}
	return tftypes.NewValue(listType, children), nil
}

func mapToTerraform5Value(values map[string]Value, mapType tftypes.Map) (tftypes.Value, error) {
	children := make(map[string]tftypes.Value)
	for name, value := range values {
		child, err := ToTerraform5Value(value, mapType.ElementType)
		if err != nil {
			return tftypes.Value{}, err
		}
		children[name] = child
	}
	return tftypes.NewValue(mapType, children), nil
}

// objectToTerraform5Value is a bit of a special case as we return the inner
// value of the Value instead of the Value directly (as with the other
// functions). This is because we use this function as part of our
// implementation of the ValueCreator and ValueConverter of the Resource type
// which expects the underlying structure instead of being already converted.
func objectToTerraform5Value(values map[string]Value, objectType tftypes.Object) (map[string]tftypes.Value, error) {
	children := make(map[string]tftypes.Value)
	for name, childType := range objectType.AttributeTypes {

		// It is possible that this child type exists in the type representation
		// but not in the actual (this is because we can have optional
		// attributes in objects). So we try and retrieve the child from the
		// values but if it is not there we don't fail, instead we just set an
		// empty value in its place.

		var err error
		if value, ok := values[name]; ok {
			if children[name], err = ToTerraform5Value(value, childType); err != nil {
				return nil, err
			}
			continue
		}

		if children[name], err = ToTerraform5Value(Value{}, childType); err != nil {
			return nil, err
		}
	}
	return children, nil
}

func setToTerraform5Value(values []Value, setType tftypes.Set) (tftypes.Value, error) {
	var children []tftypes.Value
	for _, value := range values {
		child, err := ToTerraform5Value(value, setType.ElementType)
		if err != nil {
			return tftypes.Value{}, err
		}
		children = append(children, child)
	}
	return tftypes.NewValue(setType, children), nil
}

func listFromTerraform5Value(value tftypes.Value) (Value, error) {
	var children []tftypes.Value
	if err := value.As(&children); err != nil {
		return Value{}, err
	}

	var list []Value
	for _, child := range children {
		parsed, err := FromTerraform5Value(child)
		if err != nil {
			return Value{}, err
		}
		list = append(list, parsed)
	}

	return Value{
		List: list,
	}, nil
}

func mapFromTerraform5Value(value tftypes.Value) (Value, error) {
	var children map[string]tftypes.Value
	if err := value.As(&children); err != nil {
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
		Map: values,
	}, nil
}

func objectFromTerraform5Value(value tftypes.Value) (Value, error) {
	var children map[string]tftypes.Value
	if err := value.As(&children); err != nil {
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
		Object: values,
	}, nil
}

func setFromTerraform5Value(value tftypes.Value) (Value, error) {
	var children []tftypes.Value
	if err := value.As(&children); err != nil {
		return Value{}, err
	}

	var set []Value
	for _, child := range children {
		parsed, err := FromTerraform5Value(child)
		if err != nil {
			return Value{}, err
		}
		set = append(set, parsed)
	}

	return Value{
		Set: set,
	}, nil
}
