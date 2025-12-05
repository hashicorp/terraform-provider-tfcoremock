// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package complex

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-tfcoremock/internal/schema"
)

var (
	description = `A complex resource that contains five basic attributes, four complex attributes, and two nested blocks.

The five basic attributes are boolean, number, string, float, and integer (as with the tfcoremock_simple_resource).

The complex attributes are a map, a list, a set, and an object. The object type contains the same set of attributes as the schema itself, making a recursive structure. The list, set and map all contain objects which are also recursive. Blocks cannot go into attributes, so the complex attributes do not recurse on the block types.

The blocks are a nested list and a nested set. The blocks contain the same set of attributes and blocks as the schema itself, also making a recursive structure. Note, blocks contain both attributes and more blocks so the block types are fully recursive.

The complex and block types are nested %d times, at the leaf level of recursion the complex attributes and blocks only contain the simple (ie. non-recursive) attributes. This prevents a potentially infinite level of recursion.`
	markdownDescription = `A complex resource that contains five basic attributes, four complex attributes, and two nested blocks.

The five basic attributes are ''boolean'', ''number'', ''string'', ''float'', and ''integer'' (as with the ''tfcoremock_simple_resource'').

The complex attributes are a ''map'', a ''list'', a ''set'', and an ''object''. The ''object'' type contains the same set of attributes as the schema itself, making a recursive structure. The ''list'', ''set'' and ''map'' all contain objects which are also recursive. Blocks cannot go into attributes, so the complex attributes do not recurse on the block types.

The blocks are a nested ''list_block'' and a nested ''set_block''. The blocks contain the same set of attributes and blocks as the schema itself, also making a recursive structure. Note, blocks contain both attributes and more blocks so the block types are fully recursive.

The complex and block types are nested %d times, at the leaf level of recursion the complex attributes and blocks only contain the simple (ie. non-recursive) attributes. This prevents a potentially infinite level of recursion.`
)

func Schema(maxDepth int) schema.Schema {
	return schema.Schema{
		Description:         fmt.Sprintf(description, maxDepth),
		MarkdownDescription: strings.ReplaceAll(fmt.Sprintf(markdownDescription, maxDepth), "''", "`"),
		Attributes:          attributes(0, maxDepth),
		Blocks:              blocks(0, maxDepth),
	}
}

func blocks(depth, maxDepth int) map[string]schema.Block {
	if depth == maxDepth {
		return nil
	}

	blks := make(map[string]schema.Block)

	blks["list_block"] = schema.Block{
		Description:         "A list block that contains the same attributes and blocks as the root schema, allowing nested blocks and objects to be modelled.",
		MarkdownDescription: "A list block that contains the same attributes and blocks as the root schema, allowing nested blocks and objects to be modelled.",
		Attributes:          attributes(depth+1, maxDepth),
		Blocks:              blocks(depth+1, maxDepth),
		Mode:                schema.NestingModeList,
	}

	blks["set_block"] = schema.Block{
		Description:         "A set block that contains the same attributes and blocks as the root schema, allowing nested blocks and objects to be modelled.",
		MarkdownDescription: "A set block that contains the same attributes and blocks as the root schema, allowing nested blocks and objects to be modelled.",
		Attributes:          attributes(depth+1, maxDepth),
		Blocks:              blocks(depth+1, maxDepth),
		Mode:                schema.NestingModeSet,
	}

	return blks
}

func attributes(depth, maxDepth int) map[string]schema.Attribute {
	attrs := map[string]schema.Attribute{
		"bool": {
			Description:         "An optional boolean attribute, can be true or false.",
			MarkdownDescription: "An optional boolean attribute, can be true or false.",
			Optional:            true,
			Type:                schema.Boolean,
		},
		"number": {
			Description:         "An optional number attribute, can be an integer or a float.",
			MarkdownDescription: "An optional number attribute, can be an integer or a float.",
			Optional:            true,
			Type:                schema.Number,
		},
		"string": {
			Description:         "An optional string attribute.",
			MarkdownDescription: "An optional string attribute.",
			Optional:            true,
			Type:                schema.String,
		},
		"float": {
			Description:         "An optional float attribute.",
			MarkdownDescription: "An optional float attribute.",
			Optional:            true,
			Type:                schema.Float,
		},
		"integer": {
			Description:         "An optional integer attribute.",
			MarkdownDescription: "An optional integer attribute.",
			Optional:            true,
			Type:                schema.Integer,
		},
	}

	if depth < maxDepth {
		attrs["list"] = schema.Attribute{
			Description:         "A list attribute that contains objects that match the root schema, allowing for nested collections and objects to be modelled.",
			MarkdownDescription: "A list attribute that contains objects that match the root schema, allowing for nested collections and objects to be modelled.",
			Optional:            true,
			Type:                schema.List,
			List: &schema.Attribute{
				Type:   schema.Object,
				Object: attributes(depth+1, maxDepth),
			},
		}
		attrs["map"] = schema.Attribute{
			Description:         "A map attribute that contains objects that match the root schema, allowing for nested collections and objects to be modelled.",
			MarkdownDescription: "A map attribute that contains objects that match the root schema, allowing for nested collections and objects to be modelled.",
			Optional:            true,
			Type:                schema.Map,
			Map: &schema.Attribute{
				Type:   schema.Object,
				Object: attributes(depth+1, maxDepth),
			},
		}
		attrs["object"] = schema.Attribute{
			Description:         "An object attribute that matches the root schema, allowing for nested collections and objects to be modelled.",
			MarkdownDescription: "An object attribute that matches the root schema, allowing for nested collections and objects to be modelled.",
			Optional:            true,
			Type:                schema.Object,
			Object:              attributes(depth+1, maxDepth),
		}
		attrs["set"] = schema.Attribute{
			Description:         "A set attribute that contains objects that match the root schema, allowing for nested collections and objects to be modelled.",
			MarkdownDescription: "A set attribute that contains objects that match the root schema, allowing for nested collections and objects to be modelled.",
			Optional:            true,
			Type:                schema.Set,
			Set: &schema.Attribute{
				Type:   schema.Object,
				Object: attributes(depth+1, maxDepth),
			},
		}
	}

	return attrs
}
