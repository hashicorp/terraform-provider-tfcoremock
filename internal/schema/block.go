package schema

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
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
	Attributes map[string]Attribute `json:"attributes"`
	Blocks     map[string]Block     `json:"blocks"`
	Mode       string               `json:"mode"`
}

// ToTerraformBlock converts our representation of a Block into a Terraform SDK
// block so it can be passed back to Terraform Core in a resource or data source
// schema.
func (b Block) ToTerraformBlock() (tfsdk.Block, error) {
	tfAttributes := make(map[string]tfsdk.Attribute)
	tfBlocks := make(map[string]tfsdk.Block)

	for name, attribute := range b.Attributes {
		tfAttribute, err := attribute.ToTerraformAttribute()
		if err != nil {
			return tfsdk.Block{}, err
		}
		tfAttributes[name] = tfAttribute
	}

	for name, block := range b.Blocks {
		tfBlock, err := block.ToTerraformBlock()
		if err != nil {
			return tfsdk.Block{}, err
		}
		tfBlocks[name] = tfBlock
	}

	switch b.Mode {
	case "", NestingModeList:
		return tfsdk.Block{
			Attributes:  tfAttributes,
			Blocks:      tfBlocks,
			NestingMode: tfsdk.BlockNestingModeList,
		}, nil
	case NestingModeSet:
		return tfsdk.Block{
			Attributes:  tfAttributes,
			Blocks:      tfBlocks,
			NestingMode: tfsdk.BlockNestingModeSet,
		}, nil
	default:
		return tfsdk.Block{}, errors.New("invalid nesting mode: " + b.Mode)
	}
}

func blocksToTerraformBlocks(blocks map[string]Block) (map[string]tfsdk.Block, error) {
	tfBlocks := make(map[string]tfsdk.Block)
	for name, block := range blocks {
		tfBlock, err := block.ToTerraformBlock()
		if err != nil {
			return nil, err
		}
		tfBlocks[name] = tfBlock
	}
	return tfBlocks, nil
}
