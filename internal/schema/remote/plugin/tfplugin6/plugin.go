package tfplugin6

import (
	"context"
	"github.com/hashicorp/go-plugin"
	proto "github.com/hashicorp/terraform-provider-mock/internal/proto/tfplugin6"
	"github.com/hashicorp/terraform-provider-mock/internal/schema"
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
	for range from.Attributes {
		//attributes[attribute.Name] = fromAttribute(attribute)
	}

	blocks := make(map[string]schema.Block)
	for _, block := range from.BlockTypes {
		newBlock := fromBlock(block.Block)
		switch block.Nesting {
		case proto.Schema_NestedBlock_SET:
			newBlock.Mode = schema.NestingModeSet
		default:
			// This isn't quite right, but let's default to list for everything
			// else instead of just failing.
			newBlock.Mode = schema.NestingModeList
		}
		blocks[block.TypeName] = newBlock
	}

	return schema.Block{
		Attributes: attributes,
		Blocks:     blocks,
	}
}
