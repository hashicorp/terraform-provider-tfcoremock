// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"

	"github.com/hashicorp/terraform-provider-tfcoremock/internal/data"
)

type Client interface {
	ReadResource(ctx context.Context, id string) (*data.Resource, error)
	WriteResource(ctx context.Context, value *data.Resource) error
	UpdateResource(ctx context.Context, value *data.Resource) error
	DeleteResource(ctx context.Context, id string) error
	ReadDataSource(ctx context.Context, id string) (*data.Resource, error)
}
