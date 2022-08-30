package data

import (
	"errors"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var _ tftypes.ValueConverter = &Resource{}
var _ tftypes.ValueCreator = &Resource{}

// Resource is the data structure that is actually written into our data stores.
//
// It currently only publicly contains the Values mapping of attribute names to
// actual values. It is designed as a bridge between the Terraform SDK
// representation of a value and a generic JSON representation that can be
// read/written externally. In theory, any terraform object can be represented
// as a Resource. In practice, there will probably be edge cases and types that
// have been missed.
//
// If we could write tftypes.Value into a human friendly format, and read back
// any changes from that then we wouldn't need this bridge. But, we can't do
// that using the current SDK so we handle it ourselves here.
//
// You must call the WithType function manually to attach the object type before
// attempting to convert a Resource into a Terraform SDK value.
//
// The types are attached automatically when converting from a Terraform SDK
// object.
type Resource struct {
	Values map[string]Value `json:"values"`

	objectType tftypes.Object
}

// GetId returns the ID of the resource.
//
// It assumes the ID value exists and is a string type.
func (r Resource) GetId() string {
	return *r.Values["id"].String
}

// WithType adds type information into a Resource as this is not stored as part
// of our external API.
//
// You must call this function to set the type information before using
// ToTerraform5Value(). The type information can usually be retrieved from the
// Terraform SDK, so this information should be readily available it just needs
// to be added after the Resource has been created.
func (r *Resource) WithType(objectType tftypes.Object) *Resource {
	r.objectType = objectType
	return r
}

// ToTerraform5Value ensures that Resource implements the tftypes.ValueCreator
// interface, and so can be converted into Terraform types easily.
func (r Resource) ToTerraform5Value() (interface{}, error) {
	return objectToTerraform5Value(r.Values, r.objectType)
}

// FromTerraform5Value ensures that Resource implements the
// tftypes.ValueConverter interface, and so can be converted from Terraform
// types easily.
func (r *Resource) FromTerraform5Value(value tftypes.Value) error {
	// It has to be an object we are converting from.
	if !value.Type().Is(tftypes.Object{}) {
		return errors.New("can only convert between object types")
	}

	values, err := FromTerraform5Value(value)
	if err != nil {
		return err
	}

	// We know these kinds of conversions are safe now, as we checked the type
	// at the beginning.
	r.Values = values.Object
	r.objectType = value.Type().(tftypes.Object)
	return nil
}
