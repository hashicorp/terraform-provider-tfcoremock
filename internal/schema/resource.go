// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/numberplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var (
	resources = &AttributeTypes[schema.Attribute]{}
)

func init() {
	resources.asBoolean = asResourceBool
	resources.asFloat = asResourceFloat
	resources.asInteger = asResourceInteger
	resources.asNumber = asResourceNumber
	resources.asString = asResourceString
	resources.asList = asResourceList
	resources.asNestedList = asResourceNestedList
	resources.asMap = asResourceMap
	resources.asNestedMap = asResourceNestedMap
	resources.asSet = asResourceSet
	resources.asNestedSet = asResourceNestedSet
	resources.asObject = asResourceObject
	resources.asNestedObject = asResourceNestedObject
}

func asResourceBool(attribute Attribute) (*schema.Attribute, error) {
	tfAttribute := schema.BoolAttribute{
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
		Optional:            attribute.Optional,
		Required:            attribute.Required,
		Computed:            attribute.Computed,
		Sensitive:           attribute.Sensitive,
	}

	if attribute.Computed {
		tfAttribute.PlanModifiers = append(tfAttribute.PlanModifiers, boolplanmodifier.UseStateForUnknown())
	}

	if attribute.Replace {
		tfAttribute.PlanModifiers = append(tfAttribute.PlanModifiers, boolplanmodifier.RequiresReplace())
	}

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asResourceFloat(attribute Attribute) (*schema.Attribute, error) {
	tfAttribute := schema.Float64Attribute{
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
		Optional:            attribute.Optional,
		Required:            attribute.Required,
		Computed:            attribute.Computed,
		Sensitive:           attribute.Sensitive,
	}

	if attribute.Computed {
		tfAttribute.PlanModifiers = append(tfAttribute.PlanModifiers, float64planmodifier.UseStateForUnknown())
	}

	if attribute.Replace {
		tfAttribute.PlanModifiers = append(tfAttribute.PlanModifiers, float64planmodifier.RequiresReplace())
	}

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asResourceInteger(attribute Attribute) (*schema.Attribute, error) {
	tfAttribute := schema.Int64Attribute{
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
		Optional:            attribute.Optional,
		Required:            attribute.Required,
		Computed:            attribute.Computed,
		Sensitive:           attribute.Sensitive,
	}

	if attribute.Computed {
		tfAttribute.PlanModifiers = append(tfAttribute.PlanModifiers, int64planmodifier.UseStateForUnknown())
	}

	if attribute.Replace {
		tfAttribute.PlanModifiers = append(tfAttribute.PlanModifiers, int64planmodifier.RequiresReplace())
	}

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asResourceNumber(attribute Attribute) (*schema.Attribute, error) {
	tfAttribute := schema.NumberAttribute{
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
		Optional:            attribute.Optional,
		Required:            attribute.Required,
		Computed:            attribute.Computed,
		Sensitive:           attribute.Sensitive,
	}

	if attribute.Computed {
		tfAttribute.PlanModifiers = append(tfAttribute.PlanModifiers, numberplanmodifier.UseStateForUnknown())
	}

	if attribute.Replace {
		tfAttribute.PlanModifiers = append(tfAttribute.PlanModifiers, numberplanmodifier.RequiresReplace())
	}

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asResourceString(attribute Attribute) (*schema.Attribute, error) {
	tfAttribute := schema.StringAttribute{
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
		Optional:            attribute.Optional,
		Required:            attribute.Required,
		Computed:            attribute.Computed,
		Sensitive:           attribute.Sensitive,
	}

	if attribute.Computed {
		tfAttribute.PlanModifiers = append(tfAttribute.PlanModifiers, stringplanmodifier.UseStateForUnknown())
	}

	if attribute.Replace {
		tfAttribute.PlanModifiers = append(tfAttribute.PlanModifiers, stringplanmodifier.RequiresReplace())
	}

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asResourceList(attribute Attribute) (*schema.Attribute, error) {
	tfAttribute := schema.ListAttribute{
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
		Optional:            attribute.Optional,
		Required:            attribute.Required,
		Computed:            attribute.Computed,
		Sensitive:           attribute.Sensitive,
	}

	elem, err := ToTerraformAttribute(*attribute.List, resources)
	if err != nil {
		return nil, err
	}
	tfAttribute.ElementType = (*elem).GetType()

	if attribute.Computed {
		tfAttribute.PlanModifiers = append(tfAttribute.PlanModifiers, listplanmodifier.UseStateForUnknown())
	}

	if attribute.Replace {
		tfAttribute.PlanModifiers = append(tfAttribute.PlanModifiers, listplanmodifier.RequiresReplace())
	}

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asResourceNestedList(attribute Attribute) (*schema.Attribute, error) {
	tfAttribute := schema.ListNestedAttribute{
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
		Optional:            attribute.Optional,
		Required:            attribute.Required,
		Computed:            attribute.Computed,
		Sensitive:           attribute.Sensitive,
	}

	var err error
	if tfAttribute.NestedObject.Attributes, err = attributesToTerraformResourceAttributes(attribute.List.Object); err != nil {
		return nil, err
	}

	if attribute.Computed {
		tfAttribute.PlanModifiers = append(tfAttribute.PlanModifiers, listplanmodifier.UseStateForUnknown())
	}

	if attribute.Replace {
		tfAttribute.PlanModifiers = append(tfAttribute.PlanModifiers, listplanmodifier.RequiresReplace())
	}

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asResourceMap(attribute Attribute) (*schema.Attribute, error) {
	tfAttribute := schema.MapAttribute{
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
		Optional:            attribute.Optional,
		Required:            attribute.Required,
		Computed:            attribute.Computed,
		Sensitive:           attribute.Sensitive,
	}

	elem, err := ToTerraformAttribute(*attribute.Map, resources)
	if err != nil {
		return nil, err
	}
	tfAttribute.ElementType = (*elem).GetType()

	if attribute.Computed {
		tfAttribute.PlanModifiers = append(tfAttribute.PlanModifiers, mapplanmodifier.UseStateForUnknown())
	}

	if attribute.Replace {
		tfAttribute.PlanModifiers = append(tfAttribute.PlanModifiers, mapplanmodifier.RequiresReplace())
	}

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asResourceNestedMap(attribute Attribute) (*schema.Attribute, error) {
	tfAttribute := schema.MapNestedAttribute{
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
		Optional:            attribute.Optional,
		Required:            attribute.Required,
		Computed:            attribute.Computed,
		Sensitive:           attribute.Sensitive,
	}

	var err error
	if tfAttribute.NestedObject.Attributes, err = attributesToTerraformResourceAttributes(attribute.Map.Object); err != nil {
		return nil, err
	}

	if attribute.Computed {
		tfAttribute.PlanModifiers = append(tfAttribute.PlanModifiers, mapplanmodifier.UseStateForUnknown())
	}

	if attribute.Replace {
		tfAttribute.PlanModifiers = append(tfAttribute.PlanModifiers, mapplanmodifier.RequiresReplace())
	}

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asResourceSet(attribute Attribute) (*schema.Attribute, error) {
	tfAttribute := schema.SetAttribute{
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
		Optional:            attribute.Optional,
		Required:            attribute.Required,
		Computed:            attribute.Computed,
		Sensitive:           attribute.Sensitive,
	}

	elem, err := ToTerraformAttribute(*attribute.Set, resources)
	if err != nil {
		return nil, err
	}
	tfAttribute.ElementType = (*elem).GetType()

	if attribute.Computed {
		tfAttribute.PlanModifiers = append(tfAttribute.PlanModifiers, setplanmodifier.UseStateForUnknown())
	}

	if attribute.Replace {
		tfAttribute.PlanModifiers = append(tfAttribute.PlanModifiers, setplanmodifier.RequiresReplace())
	}

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asResourceNestedSet(attribute Attribute) (*schema.Attribute, error) {
	tfAttribute := schema.SetNestedAttribute{
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
		Optional:            attribute.Optional,
		Required:            attribute.Required,
		Computed:            attribute.Computed,
		Sensitive:           attribute.Sensitive,
	}

	var err error
	if tfAttribute.NestedObject.Attributes, err = attributesToTerraformResourceAttributes(attribute.Set.Object); err != nil {
		return nil, err
	}

	if attribute.Computed {
		tfAttribute.PlanModifiers = append(tfAttribute.PlanModifiers, setplanmodifier.UseStateForUnknown())
	}

	if attribute.Replace {
		tfAttribute.PlanModifiers = append(tfAttribute.PlanModifiers, setplanmodifier.RequiresReplace())
	}

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asResourceObject(attribute Attribute) (*schema.Attribute, error) {
	tfAttribute := schema.ObjectAttribute{
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
		Optional:            attribute.Optional,
		Required:            attribute.Required,
		Computed:            attribute.Computed,
		Sensitive:           attribute.Sensitive,
	}

	types := make(map[string]attr.Type)
	attributes, err := attributesToTerraformResourceAttributes(attribute.Object)
	if err != nil {
		return nil, err
	}
	for key, value := range attributes {
		types[key] = value.GetType()
	}
	tfAttribute.AttributeTypes = types

	if attribute.Computed {
		tfAttribute.PlanModifiers = append(tfAttribute.PlanModifiers, objectplanmodifier.UseStateForUnknown())
	}

	if attribute.Replace {
		tfAttribute.PlanModifiers = append(tfAttribute.PlanModifiers, objectplanmodifier.RequiresReplace())
	}

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asResourceNestedObject(attribute Attribute) (*schema.Attribute, error) {
	tfAttribute := schema.SingleNestedAttribute{
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
		Optional:            attribute.Optional,
		Required:            attribute.Required,
		Computed:            attribute.Computed,
		Sensitive:           attribute.Sensitive,
	}

	var err error
	if tfAttribute.Attributes, err = attributesToTerraformResourceAttributes(attribute.Object); err != nil {
		return nil, err
	}

	if attribute.Computed {
		tfAttribute.PlanModifiers = append(tfAttribute.PlanModifiers, objectplanmodifier.UseStateForUnknown())
	}

	if attribute.Replace {
		tfAttribute.PlanModifiers = append(tfAttribute.PlanModifiers, objectplanmodifier.RequiresReplace())
	}

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}
