package complex

import (
	"github.com/hashicorp/terraform-provider-mock/internal/schema"
)

func Schema(maxDepth int) schema.Schema {
	return schema.Schema{
		Attributes: attributes(0, maxDepth),
		Blocks:     blocks(0, maxDepth),
	}
}

func blocks(depth, maxDepth int) map[string]schema.Block {
	if depth == maxDepth {
		return nil
	}

	blks := make(map[string]schema.Block)

	blks["list_block"] = schema.Block{
		Attributes: attributes(depth+1, maxDepth),
		Blocks:     blocks(depth+1, maxDepth),
		Mode:       schema.NestingModeList,
	}

	blks["set_block"] = schema.Block{
		Attributes: attributes(depth+1, maxDepth),
		Blocks:     blocks(depth+1, maxDepth),
		Mode:       schema.NestingModeSet,
	}

	return blks
}

func attributes(depth, maxDepth int) map[string]schema.Attribute {
	attrs := map[string]schema.Attribute{
		"bool": {
			Optional: true,
			Type:     schema.Boolean,
		},
		"number": {
			Optional: true,
			Type:     schema.Number,
		},
		"string": {
			Optional: true,
			Type:     schema.String,
		},
		"float": {
			Optional: true,
			Type:     schema.Float,
		},
		"integer": {
			Optional: true,
			Type:     schema.Integer,
		},
	}

	if depth < maxDepth {
		attrs["list"] = schema.Attribute{
			Optional: true,
			Type:     schema.List,
			List: &schema.Attribute{
				Type:   schema.Object,
				Object: attributes(depth+1, maxDepth),
			},
		}
		attrs["map"] = schema.Attribute{
			Optional: true,
			Type:     schema.Map,
			Map: &schema.Attribute{
				Type:   schema.Object,
				Object: attributes(depth+1, maxDepth),
			},
		}
		attrs["object"] = schema.Attribute{
			Optional: true,
			Type:     schema.Object,
			Object:   attributes(depth+1, maxDepth),
		}
		attrs["set"] = schema.Attribute{
			Optional: true,
			Type:     schema.Set,
			Set: &schema.Attribute{
				Type:   schema.Object,
				Object: attributes(depth+1, maxDepth),
			},
		}
	}

	return attrs
}
