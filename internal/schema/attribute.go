package schema

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

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
}

// ToTerraformAttribute converts our representation of an Attribute into a
// Terraform SDK attribute so it can be passed back to Terraform Core in a
// resource or data source schema.
func (a Attribute) ToTerraformAttribute() (tfsdk.Attribute, error) {
	switch a.Type {
	case Boolean:
		return withType(a.getTerraformAttribute(), types.BoolType), nil
	case Float:
		return withType(a.getTerraformAttribute(), types.Float64Type), nil
	case Integer:
		return withType(a.getTerraformAttribute(), types.Int64Type), nil
	case Number:
		return withType(a.getTerraformAttribute(), types.NumberType), nil
	case String:
		return withType(a.getTerraformAttribute(), types.StringType), nil
	case List:
		if a.List.Type == Object {
			attributes, err := attributesToTerraformAttributes(a.List.Object)
			if err != nil {
				return tfsdk.Attribute{}, nil
			}
			return asList(a.getTerraformAttribute(), attributes), nil
		}
		attribute, err := a.List.ToTerraformAttribute()
		if err != nil {
			return tfsdk.Attribute{}, nil
		}
		return withType(a.getTerraformAttribute(), types.ListType{ElemType: attribute.Type}), nil
	case Map:
		if a.Map.Type == Object {
			attributes, err := attributesToTerraformAttributes(a.Map.Object)
			if err != nil {
				return tfsdk.Attribute{}, nil
			}
			return asMap(a.getTerraformAttribute(), attributes), nil
		}
		attribute, err := a.Map.ToTerraformAttribute()
		if err != nil {
			return tfsdk.Attribute{}, nil
		}
		return withType(a.getTerraformAttribute(), types.MapType{ElemType: attribute.Type}), nil
	case Set:
		if a.Set.Type == Object {
			attributes, err := attributesToTerraformAttributes(a.Set.Object)
			if err != nil {
				return tfsdk.Attribute{}, nil
			}
			return asSet(a.getTerraformAttribute(), attributes), nil
		}
		attribute, err := a.Set.ToTerraformAttribute()
		if err != nil {
			return tfsdk.Attribute{}, nil
		}
		return withType(a.getTerraformAttribute(), types.SetType{ElemType: attribute.Type}), nil
	case Object:
		attributes, err := attributesToTerraformAttributes(a.Object)
		if err != nil {
			return tfsdk.Attribute{}, err
		}
		return asObject(a.getTerraformAttribute(), attributes), nil
	default:
		return tfsdk.Attribute{}, errors.New("unrecognized attribute type: " + string(a.Type))
	}
}

func (a Attribute) getTerraformAttribute() tfsdk.Attribute {
	attribute := tfsdk.Attribute{
		Description:         a.Description,
		MarkdownDescription: a.MarkdownDescription,
		Optional:            a.Optional,
		Required:            a.Required,
		Computed:            a.Computed,
		Sensitive:           a.Sensitive,
	}

	if a.Computed {
		attribute.PlanModifiers = append(attribute.PlanModifiers, resource.UseStateForUnknown())
	}

	return attribute
}

func withType(attribute tfsdk.Attribute, t attr.Type) tfsdk.Attribute {
	attribute.Type = t
	return attribute
}

func asObject(attribute tfsdk.Attribute, attributes map[string]tfsdk.Attribute) tfsdk.Attribute {
	attribute.Attributes = tfsdk.SingleNestedAttributes(attributes)
	return attribute
}

func asList(attribute tfsdk.Attribute, attributes map[string]tfsdk.Attribute) tfsdk.Attribute {
	attribute.Attributes = tfsdk.ListNestedAttributes(attributes)
	return attribute
}

func asSet(attribute tfsdk.Attribute, attributes map[string]tfsdk.Attribute) tfsdk.Attribute {
	attribute.Attributes = tfsdk.SetNestedAttributes(attributes)
	return attribute
}

func asMap(attribute tfsdk.Attribute, attributes map[string]tfsdk.Attribute) tfsdk.Attribute {
	attribute.Attributes = tfsdk.MapNestedAttributes(attributes)
	return attribute
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
