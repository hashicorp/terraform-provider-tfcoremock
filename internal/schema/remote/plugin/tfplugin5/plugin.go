package tfplugin5

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-plugin"
	proto "github.com/hashicorp/terraform-provider-tfcoremock/internal/proto/tfplugin5"
	"github.com/hashicorp/terraform-provider-tfcoremock/internal/schema"
	"google.golang.org/grpc"
)

type ProviderPlugin struct {
	plugin.Plugin
}

func (p *ProviderPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &Provider{
		client: proto.NewProviderClient(c),
	}, nil
}

func (p *ProviderPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	panic("not implemented")
}

type Provider struct {
	client proto.ProviderClient
}

func (p *Provider) GetSchema(ctx context.Context) (map[string]schema.Schema, map[string]schema.Schema, error) {
	resp, err := p.client.GetSchema(ctx, &proto.GetProviderSchema_Request{})
	if err != nil {
		return nil, nil, err
	}

	resources := make(map[string]schema.Schema)
	dataSources := make(map[string]schema.Schema)
	for key, schema := range resp.ResourceSchemas {
		resources[key] = fromSchema(schema)
	}

	for key, schema := range resp.DataSourceSchemas {
		dataSources[key] = fromSchema(schema)
	}
	return resources, dataSources, nil
}

func fromSchema(from *proto.Schema) schema.Schema {
	block := fromBlock(from.Block)
	return schema.Schema{
		Attributes: block.Attributes,
		Blocks:     block.Blocks,
	}
}

func fromBlock(from *proto.Schema_Block) schema.Block {
	attributes := make(map[string]schema.Attribute)
	for _, attribute := range from.Attributes {
		attributes[attribute.Name] = fromAttribute(attribute)
	}

	blocks := make(map[string]schema.Block)
	for _, block := range from.BlockTypes {
		newBlock := fromBlock(block.Block)
		switch block.Nesting {
		case proto.Schema_NestedBlock_SET:
			newBlock.Mode = schema.NestingModeSet
		case proto.Schema_NestedBlock_LIST:
			newBlock.Mode = schema.NestingModeList
		case proto.Schema_NestedBlock_GROUP, proto.Schema_NestedBlock_SINGLE:
			newBlock.Mode = schema.NestingModeSingle
		case proto.Schema_NestedBlock_MAP:
			panic("the developers of this plugin are very interested to find a provider that has used this type of nesting, please report this crash as a bug in our repositor")
		default:
			panic("unrecognized nested block: " + block.Nesting.String())
		}
		blocks[block.TypeName] = newBlock
	}

	return schema.Block{
		Attributes: attributes,
		Blocks:     blocks,
	}
}

func fromAttribute(from *proto.Schema_Attribute) schema.Attribute {
	var fromType interface{}
	if err := json.Unmarshal(from.Type, &fromType); err != nil {
		panic("incompatible type: " + string(from.Type))
	}

	switch concrete := fromType.(type) {
	case string:
		return populateCommon(from, createSimpleAttribute(concrete, nil))
	case []interface{}:
		var last schema.Attribute
		for ix := len(concrete) - 1; ix > 0; ix-- {
			last = createComplexAttribute(concrete[ix], &last)
		}
		return populateCommon(from, createComplexAttribute(concrete[0], &last))
	default:
		panic(fmt.Sprintf("unrecognized type: %T", fromType))
	}
}

func populateCommon(from *proto.Schema_Attribute, to schema.Attribute) schema.Attribute {
	to.Optional = from.Optional
	to.Required = from.Required
	to.Computed = from.Computed
	return to
}

func createComplexAttribute(t interface{}, last *schema.Attribute) schema.Attribute {
	switch concrete := t.(type) {
	case string:
		return createSimpleAttribute(concrete, last)
	case []interface{}:
		switch concrete[0].(string) {
		case "object":
			concreteMap, ok := concrete[1].(map[string]interface{})
			if !ok {
				panic(fmt.Sprintf("unrecognized object type: %v", concrete[1]))
			}

			attributes := make(map[string]schema.Attribute)
			for key, value := range concreteMap {
				attribute := createComplexAttribute(value, nil)
				attribute.Optional = true
				attributes[key] = attribute
			}
			return schema.Attribute{
				Type:   schema.Object,
				Object: attributes,
			}
		default:
			var last schema.Attribute
			for ix := len(concrete) - 1; ix > 0; ix-- {
				last = createComplexAttribute(concrete[ix], &last)
			}
			return createComplexAttribute(concrete[0], &last)
		}
	default:
		panic(fmt.Sprintf("unrecognized type: %T", t))
	}
}

func createSimpleAttribute(t string, last *schema.Attribute) schema.Attribute {
	switch t {
	case "bool":
		return schema.Attribute{
			Type: schema.Boolean,
		}
	case "string":
		return schema.Attribute{
			Type: schema.String,
		}
	case "number":
		return schema.Attribute{
			Type: schema.Number,
		}
	case "list":
		return schema.Attribute{
			Type: schema.List,
			List: last,
		}
	case "set":
		return schema.Attribute{
			Type: schema.Set,
			Set:  last,
		}
	case "map":
		return schema.Attribute{
			Type: schema.Map,
			Map:  last,
		}
	default:
		panic("unrecognized type: " + t)
	}
}
