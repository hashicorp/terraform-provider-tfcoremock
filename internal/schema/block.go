// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"fmt"

	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/pkg/errors"
)

const (
	NestingModeList = "list"
	NestingModeSet  = "set"
)

// Block defines an internal representation of a Terraform block in a schema.
//
// It is designed to be read dynamically from a JSON object, allowing schemas,
// blocks and attributes to be defined dynamically by the user of the provider.
type Block struct {
	Description         string `json:"-"` // Dynamic resources don't need descriptions so hide them from the exposed JSON schema.
	MarkdownDescription string `json:"-"` // Dynamic resources don't need descriptions so hide them from the exposed JSON schema.

	Attributes map[string]Attribute `json:"attributes"`
	Blocks     map[string]Block     `json:"blocks"`
	Mode       string               `json:"mode"`
}

type ToListBlock[B any, A any] func(block Block, blocks map[string]B, attributes map[string]A) *B
type ToSetBlock[B any, A any] func(block Block, blocks map[string]B, attributes map[string]A) *B

// ToTerraformBlock converts our representation of a Block into a Terraform SDK
// block so it can be passed back to Terraform Core in a resource or data source
// schema.
func ToTerraformBlock[B, A any](b Block, toListBlock ToListBlock[B, A], toSetBlock ToSetBlock[B, A], attributeTypes *AttributeTypes[A]) (*B, error) {
	tfAttributes := make(map[string]A)
	tfBlocks := make(map[string]B)

	for name, attribute := range b.Attributes {
		attribute, err := ToTerraformAttribute(attribute, attributeTypes)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create attribute '%s'", name)
		}
		tfAttributes[name] = *attribute
	}

	for name, block := range b.Blocks {
		block, err := ToTerraformBlock(block, toListBlock, toSetBlock, attributeTypes)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create block '%s'", name)
		}
		tfBlocks[name] = *block
	}

	switch b.Mode {
	case "", NestingModeList:
		return toListBlock(b, tfBlocks, tfAttributes), nil
	case NestingModeSet:
		return toSetBlock(b, tfBlocks, tfAttributes), nil
	default:
		return nil, fmt.Errorf("invalid nesting mode '%s'", b.Mode)
	}
}

func blocksToTerraformResourceBlocks(blocks map[string]Block) (map[string]resource_schema.Block, error) {
	toListBlock := func(block Block, blocks map[string]resource_schema.Block, attributes map[string]resource_schema.Attribute) *resource_schema.Block {
		var tfBlock resource_schema.Block
		tfBlock = resource_schema.ListNestedBlock{
			Description:         block.Description,
			MarkdownDescription: block.MarkdownDescription,
			NestedObject: resource_schema.NestedBlockObject{
				Attributes: attributes,
				Blocks:     blocks,
			},
		}
		return &tfBlock
	}

	toSetBlock := func(block Block, blocks map[string]resource_schema.Block, attributes map[string]resource_schema.Attribute) *resource_schema.Block {
		var tfBlock resource_schema.Block
		tfBlock = resource_schema.SetNestedBlock{
			Description:         block.Description,
			MarkdownDescription: block.MarkdownDescription,
			NestedObject: resource_schema.NestedBlockObject{
				Attributes: attributes,
				Blocks:     blocks,
			},
		}
		return &tfBlock
	}

	tfBlocks := make(map[string]resource_schema.Block)
	for name, block := range blocks {
		block, err := ToTerraformBlock(block, toListBlock, toSetBlock, resources)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create block '%s'", name)
		}
		tfBlocks[name] = *block
	}
	return tfBlocks, nil
}

func blocksToTerraformDataSourceBlocks(blocks map[string]Block) (map[string]datasource_schema.Block, error) {
	toListBlock := func(block Block, blocks map[string]datasource_schema.Block, attributes map[string]datasource_schema.Attribute) *datasource_schema.Block {
		var tfBlock datasource_schema.Block
		tfBlock = datasource_schema.ListNestedBlock{
			Description:         block.Description,
			MarkdownDescription: block.MarkdownDescription,
			NestedObject: datasource_schema.NestedBlockObject{
				Attributes: attributes,
				Blocks:     blocks,
			},
		}
		return &tfBlock
	}

	toSetBlock := func(block Block, blocks map[string]datasource_schema.Block, attributes map[string]datasource_schema.Attribute) *datasource_schema.Block {
		var tfBlock datasource_schema.Block
		tfBlock = datasource_schema.SetNestedBlock{
			Description:         block.Description,
			MarkdownDescription: block.MarkdownDescription,
			NestedObject: datasource_schema.NestedBlockObject{
				Attributes: attributes,
				Blocks:     blocks,
			},
		}
		return &tfBlock
	}

	tfBlocks := make(map[string]datasource_schema.Block)
	for name, block := range blocks {
		block, err := ToTerraformBlock(block, toListBlock, toSetBlock, datasources)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create block '%s'", name)
		}
		tfBlocks[name] = *block
	}
	return tfBlocks, nil
}
