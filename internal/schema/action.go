// Copyright IBM Corp. 2022, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/attr"
)

var (
	actions = &AttributeTypes[schema.Attribute]{}
)

func init() {
	actions.asBoolean = asActionBool
	actions.asFloat = asActionFloat
	actions.asInteger = asActionInteger
	actions.asNumber = asActionNumber
	actions.asString = asActionString
	actions.asList = asActionList
	actions.asNestedList = asActionNestedList
	actions.asMap = asActionMap
	actions.asNestedMap = asActionNestedMap
	actions.asSet = asActionSet
	actions.asNestedSet = asActionNestedSet
	actions.asObject = asActionObject
	actions.asNestedObject = asActionNestedObject
}

func asActionBool(attribute Attribute) (*schema.Attribute, error) {
	// action schemas don't have computed, but we share this definition with
	// resources and data sources. therefore, we set optional to true if the
	// attribute is computed.

	tfAttribute := schema.BoolAttribute{
		Required:            attribute.Required,
		Optional:            attribute.Optional || attribute.Computed,
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
	}

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asActionFloat(attribute Attribute) (*schema.Attribute, error) {
	// action schemas don't have computed, but we share this definition with
	// resources and data sources. therefore, we set optional to true if the
	// attribute is computed.

	tfAttribute := schema.Float64Attribute{
		Required:            attribute.Required,
		Optional:            attribute.Optional || attribute.Computed,
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
	}

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asActionInteger(attribute Attribute) (*schema.Attribute, error) {
	// action schemas don't have computed, but we share this definition with
	// resources and data sources. therefore, we set optional to true if the
	// attribute is computed.

	tfAttribute := schema.Int64Attribute{
		Required:            attribute.Required,
		Optional:            attribute.Optional || attribute.Computed,
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
	}

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asActionNumber(attribute Attribute) (*schema.Attribute, error) {
	// action schemas don't have computed, but we share this definition with
	// resources and data sources. therefore, we set optional to true if the
	// attribute is computed.

	tfAttribute := schema.NumberAttribute{
		Required:            attribute.Required,
		Optional:            attribute.Optional || attribute.Computed,
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
	}

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asActionString(attribute Attribute) (*schema.Attribute, error) {
	// action schemas don't have computed, but we share this definition with
	// resources and data sources. therefore, we set optional to true if the
	// attribute is computed.

	tfAttribute := schema.StringAttribute{
		Required:            attribute.Required,
		Optional:            attribute.Optional || attribute.Computed,
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
	}

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asActionList(attribute Attribute) (*schema.Attribute, error) {
	// action schemas don't have computed, but we share this definition with
	// resources and data sources. therefore, we set optional to true if the
	// attribute is computed.

	tfAttribute := schema.ListAttribute{
		Required:            attribute.Required,
		Optional:            attribute.Optional || attribute.Computed,
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
	}

	elem, err := ToTerraformAttribute(*attribute.List, actions)
	if err != nil {
		return nil, err
	}
	tfAttribute.ElementType = (*elem).GetType()

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asActionNestedList(attribute Attribute) (*schema.Attribute, error) {
	// action schemas don't have computed, but we share this definition with
	// resources and data sources. therefore, we set optional to true if the
	// attribute is computed.

	tfAttribute := schema.ListNestedAttribute{
		Required:            attribute.Required,
		Optional:            attribute.Optional || attribute.Computed,
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
	}

	var err error
	if tfAttribute.NestedObject.Attributes, err = attributesToTerraformActionAttributes(attribute.List.Object); err != nil {
		return nil, err
	}

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asActionMap(attribute Attribute) (*schema.Attribute, error) {
	// action schemas don't have computed, but we share this definition with
	// resources and data sources. therefore, we set optional to true if the
	// attribute is computed.

	tfAttribute := schema.MapAttribute{
		Required:            attribute.Required,
		Optional:            attribute.Optional || attribute.Computed,
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
	}

	elem, err := ToTerraformAttribute(*attribute.Map, actions)
	if err != nil {
		return nil, err
	}
	tfAttribute.ElementType = (*elem).GetType()

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asActionNestedMap(attribute Attribute) (*schema.Attribute, error) {
	// action schemas don't have computed, but we share this definition with
	// resources and data sources. therefore, we set optional to true if the
	// attribute is computed.

	tfAttribute := schema.MapNestedAttribute{
		Required:            attribute.Required,
		Optional:            attribute.Optional || attribute.Computed,
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
	}

	var err error
	if tfAttribute.NestedObject.Attributes, err = attributesToTerraformActionAttributes(attribute.Map.Object); err != nil {
		return nil, err
	}

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asActionSet(attribute Attribute) (*schema.Attribute, error) {
	// action schemas don't have computed, but we share this definition with
	// resources and data sources. therefore, we set optional to true if the
	// attribute is computed.

	tfAttribute := schema.SetAttribute{
		Required:            attribute.Required,
		Optional:            attribute.Optional || attribute.Computed,
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
	}

	elem, err := ToTerraformAttribute(*attribute.Set, actions)
	if err != nil {
		return nil, err
	}
	tfAttribute.ElementType = (*elem).GetType()

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asActionNestedSet(attribute Attribute) (*schema.Attribute, error) {
	// action schemas don't have computed, but we share this definition with
	// resources and data sources. therefore, we set optional to true if the
	// attribute is computed.

	tfAttribute := schema.SetNestedAttribute{
		Required:            attribute.Required,
		Optional:            attribute.Optional || attribute.Computed,
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
	}

	var err error
	if tfAttribute.NestedObject.Attributes, err = attributesToTerraformActionAttributes(attribute.Set.Object); err != nil {
		return nil, err
	}

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asActionObject(attribute Attribute) (*schema.Attribute, error) {
	// action schemas don't have computed, but we share this definition with
	// resources and data sources. therefore, we set optional to true if the
	// attribute is computed.

	tfAttribute := schema.ObjectAttribute{
		Required:            attribute.Required,
		Optional:            attribute.Optional || attribute.Computed,
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
	}

	types := make(map[string]attr.Type)
	attributes, err := attributesToTerraformActionAttributes(attribute.Object)
	if err != nil {
		return nil, err
	}
	for key, value := range attributes {
		types[key] = value.GetType()
	}
	tfAttribute.AttributeTypes = types

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asActionNestedObject(attribute Attribute) (*schema.Attribute, error) {
	// action schemas don't have computed, but we share this definition with
	// resources and data sources. therefore, we set optional to true if the
	// attribute is computed.

	tfAttribute := schema.SingleNestedAttribute{
		Required:            attribute.Required,
		Optional:            attribute.Optional || attribute.Computed,
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
	}

	var err error
	if tfAttribute.Attributes, err = attributesToTerraformActionAttributes(attribute.Object); err != nil {
		return nil, err
	}

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}
