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
	Description         string `json:"-"` // Dynamic resources don't need descriptions so hide them from the exposed JSON schema.
	MarkdownDescription string `json:"-"` // Dynamic resources don't need descriptions so hide them from the exposed JSON schema.

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
		return a.toSimpleTerraformAttribute(types.BoolType), nil
	case Float:
		return a.toSimpleTerraformAttribute(types.Float64Type), nil
	case Integer:
		return a.toSimpleTerraformAttribute(types.Int64Type), nil
	case Number:
		return a.toSimpleTerraformAttribute(types.NumberType), nil
	case String:
		return a.toSimpleTerraformAttribute(types.StringType), nil
	case List:
		if a.List.Type == Object {
			attributes, err := attributesToTerraformAttributes(a.List.Object)
			if err != nil {
				return tfsdk.Attribute{}, nil
			}

			return tfsdk.Attribute{
				Description:         a.Description,
				MarkdownDescription: a.MarkdownDescription,
				Optional:            a.Optional,
				Required:            a.Required,
				Computed:            a.Computed,
				Attributes:          tfsdk.ListNestedAttributes(attributes),
			}, nil
		}
		attribute, err := a.List.ToTerraformAttribute()
		if err != nil {
			return tfsdk.Attribute{}, nil
		}
		return a.toSimpleTerraformAttribute(types.ListType{ElemType: attribute.Type}), nil
	case Map:
		if a.Map.Type == Object {
			attributes, err := attributesToTerraformAttributes(a.Map.Object)
			if err != nil {
				return tfsdk.Attribute{}, nil
			}

			return tfsdk.Attribute{
				Description:         a.Description,
				MarkdownDescription: a.MarkdownDescription,
				Optional:            a.Optional,
				Required:            a.Required,
				Computed:            a.Computed,
				Attributes:          tfsdk.MapNestedAttributes(attributes),
			}, nil
		}
		attribute, err := a.Map.ToTerraformAttribute()
		if err != nil {
			return tfsdk.Attribute{}, nil
		}
		return a.toSimpleTerraformAttribute(types.MapType{ElemType: attribute.Type}), nil
	case Set:
		if a.Set.Type == Object {
			attributes, err := attributesToTerraformAttributes(a.Set.Object)
			if err != nil {
				return tfsdk.Attribute{}, nil
			}

			return tfsdk.Attribute{
				Description:         a.Description,
				MarkdownDescription: a.MarkdownDescription,
				Optional:            a.Optional,
				Required:            a.Required,
				Computed:            a.Computed,
				Attributes:          tfsdk.SetNestedAttributes(attributes),
			}, nil
		}
		attribute, err := a.Set.ToTerraformAttribute()
		if err != nil {
			return tfsdk.Attribute{}, nil
		}
		return a.toSimpleTerraformAttribute(types.SetType{ElemType: attribute.Type}), nil
	case Object:
		attributes, err := attributesToTerraformAttributes(a.Object)
		if err != nil {
			return tfsdk.Attribute{}, err
		}

		return tfsdk.Attribute{
			Description:         a.Description,
			MarkdownDescription: a.MarkdownDescription,
			Optional:            a.Optional,
			Required:            a.Required,
			Computed:            a.Computed,
			Attributes:          tfsdk.SingleNestedAttributes(attributes),
		}, nil
	default:
		return tfsdk.Attribute{}, errors.New("unrecognized attribute type: " + string(a.Type))
	}
}

func (a Attribute) toSimpleTerraformAttribute(t attr.Type) tfsdk.Attribute {
	return tfsdk.Attribute{
		Description:         a.Description,
		MarkdownDescription: a.MarkdownDescription,
		Optional:            a.Optional,
		Required:            a.Required,
		Computed:            a.Computed,
		Type:                t,
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
