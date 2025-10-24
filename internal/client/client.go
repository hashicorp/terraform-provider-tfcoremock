// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"

	"github.com/hashicorp/terraform-provider-tfcoremock/internal/data"
)

func Filter(value string) *string {
	return &value
}

type Client interface {
	ReadResource(ctx context.Context, id string) (*data.Resource, error)
	WriteResource(ctx context.Context, value *data.Resource) error
	UpdateResource(ctx context.Context, value *data.Resource) error
	DeleteResource(ctx context.Context, id string) error
	ListResources(ctx context.Context, typeName *string, id *string, yield func(resource *data.Resource, err error), limit int64) error
	ReadDataSource(ctx context.Context, id string) (*data.Resource, error)
}
