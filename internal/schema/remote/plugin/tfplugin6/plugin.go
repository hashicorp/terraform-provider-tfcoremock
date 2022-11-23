package tfplugin6

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-plugin"
	proto "github.com/hashicorp/terraform-provider-tfcoremock/internal/proto/tfplugin6"
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
	resp, err := p.client.GetProviderSchema(ctx, &proto.GetProviderSchema_Request{})
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
	if from.NestedType != nil {
		return createNestedAttribute(from)
	}

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
			last = createSimpleAttribute(concrete[ix].(string), &last)
		}
		return populateCommon(from, createSimpleAttribute(concrete[0].(string), &last))
	default:
		panic(fmt.Sprintf("unrecognized type: %T", fromType))
	}
}

func createNestedAttribute(from *proto.Schema_Attribute) schema.Attribute {
	attributes := make(map[string]schema.Attribute)
	for _, attr := range from.NestedType.Attributes {
		attributes[attr.Name] = fromAttribute(attr)
	}

	object := schema.Attribute{
		Type:   schema.Object,
		Object: attributes,
	}

	switch from.NestedType.Nesting {
	case proto.Schema_Object_LIST:
		return populateCommon(from, schema.Attribute{
			Type: schema.List,
			List: &object,
		})
	case proto.Schema_Object_MAP:
		return populateCommon(from, schema.Attribute{
			Type: schema.Map,
			Map:  &object,
		})
	case proto.Schema_Object_SINGLE:
		return populateCommon(from, object)
	case proto.Schema_Object_SET:
		return populateCommon(from, schema.Attribute{
			Type: schema.Set,
			Set:  &object,
		})
	default:
		panic("unrecognized nesting type: " + from.NestedType.Nesting.String())
	}
}

func populateCommon(from *proto.Schema_Attribute, to schema.Attribute) schema.Attribute {
	to.Optional = from.Optional
	to.Required = from.Required
	to.Computed = from.Computed
	return to
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
