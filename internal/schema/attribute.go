package schema

import (
	"errors"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Attribute defines an internal representation of a Terraform attribute in a
// schema.
//
// It is designed to be read dynamically from a JSON object, allowing schemas,
// blocks and attributes to be defined dynamically by the user of the provider.
type Attribute struct {
	Type     Type `json:"type"`
	Optional bool `json:"optional"`
	Required bool `json:"required"`
	Computed bool `json:"computed"`

	List   *Attribute           `json:"list,omitempty"`
	Map    *Attribute           `json:"map,omitempty"`
	Object map[string]Attribute `json:"object,omitempty"`
	Set    *Attribute           `json:"set,omitempty"`
}

// ToTerraformAttribute converts our representation of an Attribute into a
// Terraform SDK attribute so it can be passed back to Terraform Core in a
// resource or data source schema.
func (a Attribute) ToTerraformAttribute() (tfsdk.Attribute, error) {
	switch a.Type {
	case Boolean:
		return a.toTerraformAttribute(types.BoolType), nil
	case Float:
		return a.toTerraformAttribute(types.Float64Type), nil
	case Integer:
		return a.toTerraformAttribute(types.Int64Type), nil
	case Number:
		return a.toTerraformAttribute(types.NumberType), nil
	case String:
		return a.toTerraformAttribute(types.StringType), nil
	case List:
		attribute, err := a.List.ToTerraformAttribute()
		if err != nil {
			return tfsdk.Attribute{}, nil
		}
		return a.toTerraformAttribute(types.ListType{ElemType: attribute.Type}), nil
	case Map:
		attribute, err := a.Map.ToTerraformAttribute()
		if err != nil {
			return tfsdk.Attribute{}, nil
		}
		return a.toTerraformAttribute(types.MapType{ElemType: attribute.Type}), nil
	case Set:
		attribute, err := a.Set.ToTerraformAttribute()
		if err != nil {
			return tfsdk.Attribute{}, nil
		}
		return a.toTerraformAttribute(types.SetType{ElemType: attribute.Type}), nil
	case Object:
		attributes, err := attributesToTerraformAttributes(a.Object)
		if err != nil {
			return tfsdk.Attribute{}, err
		}
		return tfsdk.Attribute{
			Optional:   a.Optional,
			Required:   a.Required,
			Computed:   a.Computed,
			Attributes: tfsdk.SingleNestedAttributes(attributes),
		}, nil
	default:
		return tfsdk.Attribute{}, errors.New("unrecognized attribute type: " + string(a.Type))
	}
}

func (a Attribute) toTerraformAttribute(t attr.Type) tfsdk.Attribute {
	return tfsdk.Attribute{
		Optional: a.Optional,
		Required: a.Required,
		Computed: a.Computed,
		Type:     t,
	}
}

func attributesToTerraformAttributes(attributes map[string]Attribute) (map[string]tfsdk.Attribute, error) {
	tfAttributes := make(map[string]tfsdk.Attribute)
	for name, attribute := range attributes {
		tfAttribute, err := attribute.ToTerraformAttribute()
		if err != nil {
			return nil, err
		}
		tfAttributes[name] = tfAttribute
	}
	return tfAttributes, nil
}
