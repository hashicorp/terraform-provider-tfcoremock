// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var (
	datasources = &AttributeTypes[schema.Attribute]{}
)

func init() {
	datasources.asBoolean = asDataSourceBool
	datasources.asFloat = asDataSourceFloat
	datasources.asInteger = asDataSourceInteger
	datasources.asNumber = asDataSourceNumber
	datasources.asString = asDataSourceString
	datasources.asList = asDataSourceList
	datasources.asNestedList = asDataSourceNestedList
	datasources.asMap = asDataSourceMap
	datasources.asNestedMap = asDataSourceNestedMap
	datasources.asSet = asDataSourceSet
	datasources.asNestedSet = asDataSourceNestedSet
	datasources.asObject = asDataSourceObject
	datasources.asNestedObject = asDataSourceNestedObject
}

func asDataSourceBool(attribute Attribute) (*schema.Attribute, error) {
	tfAttribute := schema.BoolAttribute{
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
		Optional:            attribute.Optional,
		Required:            attribute.Required,
		Computed:            attribute.Computed,
		Sensitive:           attribute.Sensitive,
	}

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asDataSourceFloat(attribute Attribute) (*schema.Attribute, error) {
	tfAttribute := schema.Float64Attribute{
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
		Optional:            attribute.Optional,
		Required:            attribute.Required,
		Computed:            attribute.Computed,
		Sensitive:           attribute.Sensitive,
	}

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asDataSourceInteger(attribute Attribute) (*schema.Attribute, error) {
	tfAttribute := schema.Int64Attribute{
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
		Optional:            attribute.Optional,
		Required:            attribute.Required,
		Computed:            attribute.Computed,
		Sensitive:           attribute.Sensitive,
	}

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asDataSourceNumber(attribute Attribute) (*schema.Attribute, error) {
	tfAttribute := schema.NumberAttribute{
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
		Optional:            attribute.Optional,
		Required:            attribute.Required,
		Computed:            attribute.Computed,
		Sensitive:           attribute.Sensitive,
	}

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asDataSourceString(attribute Attribute) (*schema.Attribute, error) {
	tfAttribute := schema.StringAttribute{
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
		Optional:            attribute.Optional,
		Required:            attribute.Required,
		Computed:            attribute.Computed,
		Sensitive:           attribute.Sensitive,
	}

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asDataSourceList(attribute Attribute) (*schema.Attribute, error) {
	tfAttribute := schema.ListAttribute{
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
		Optional:            attribute.Optional,
		Required:            attribute.Required,
		Computed:            attribute.Computed,
		Sensitive:           attribute.Sensitive,
	}

	elem, err := ToTerraformAttribute(*attribute.List, datasources)
	if err != nil {
		return nil, err
	}
	tfAttribute.ElementType = (*elem).GetType()

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asDataSourceNestedList(attribute Attribute) (*schema.Attribute, error) {
	tfAttribute := schema.ListNestedAttribute{
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
		Optional:            attribute.Optional,
		Required:            attribute.Required,
		Computed:            attribute.Computed,
		Sensitive:           attribute.Sensitive,
	}

	var err error
	if tfAttribute.NestedObject.Attributes, err = attributesToTerraformDataSourceAttributes(attribute.List.Object); err != nil {
		return nil, err
	}

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asDataSourceMap(attribute Attribute) (*schema.Attribute, error) {
	tfAttribute := schema.MapAttribute{
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
		Optional:            attribute.Optional,
		Required:            attribute.Required,
		Computed:            attribute.Computed,
		Sensitive:           attribute.Sensitive,
	}

	elem, err := ToTerraformAttribute(*attribute.Map, datasources)
	if err != nil {
		return nil, err
	}
	tfAttribute.ElementType = (*elem).GetType()

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asDataSourceNestedMap(attribute Attribute) (*schema.Attribute, error) {
	tfAttribute := schema.MapNestedAttribute{
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
		Optional:            attribute.Optional,
		Required:            attribute.Required,
		Computed:            attribute.Computed,
		Sensitive:           attribute.Sensitive,
	}

	var err error
	if tfAttribute.NestedObject.Attributes, err = attributesToTerraformDataSourceAttributes(attribute.Map.Object); err != nil {
		return nil, err
	}

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asDataSourceSet(attribute Attribute) (*schema.Attribute, error) {
	tfAttribute := schema.SetAttribute{
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
		Optional:            attribute.Optional,
		Required:            attribute.Required,
		Computed:            attribute.Computed,
		Sensitive:           attribute.Sensitive,
	}

	elem, err := ToTerraformAttribute(*attribute.Set, datasources)
	if err != nil {
		return nil, err
	}
	tfAttribute.ElementType = (*elem).GetType()

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asDataSourceNestedSet(attribute Attribute) (*schema.Attribute, error) {
	tfAttribute := schema.SetNestedAttribute{
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
		Optional:            attribute.Optional,
		Required:            attribute.Required,
		Computed:            attribute.Computed,
		Sensitive:           attribute.Sensitive,
	}

	var err error
	if tfAttribute.NestedObject.Attributes, err = attributesToTerraformDataSourceAttributes(attribute.Set.Object); err != nil {
		return nil, err
	}

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}

func asDataSourceObject(attribute Attribute) (*schema.Attribute, error) {
	tfAttribute := schema.ObjectAttribute{
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
		Optional:            attribute.Optional,
		Required:            attribute.Required,
		Computed:            attribute.Computed,
		Sensitive:           attribute.Sensitive,
	}

	types := make(map[string]attr.Type)
	attributes, err := attributesToTerraformDataSourceAttributes(attribute.Object)
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

func asDataSourceNestedObject(attribute Attribute) (*schema.Attribute, error) {
	tfAttribute := schema.SingleNestedAttribute{
		Description:         attribute.Description,
		MarkdownDescription: attribute.MarkdownDescription,
		Optional:            attribute.Optional,
		Required:            attribute.Required,
		Computed:            attribute.Computed,
		Sensitive:           attribute.Sensitive,
	}

	var err error
	if tfAttribute.Attributes, err = attributesToTerraformDataSourceAttributes(attribute.Object); err != nil {
		return nil, err
	}

	var out schema.Attribute
	out = tfAttribute
	return &out, nil
}
