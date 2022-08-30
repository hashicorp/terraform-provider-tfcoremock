package client

import (
	"context"
	"github.com/hashicorp/terraform-provider-mock/internal/data"
)

type Client interface {
	ReadResource(ctx context.Context, id string) (*data.Resource, error)
	WriteResource(ctx context.Context, value *data.Resource) error
	UpdateResource(ctx context.Context, value *data.Resource) error
	DeleteResource(ctx context.Context, id string) error
	ReadDataSource(ctx context.Context, id string) (*data.Resource, error)
}
