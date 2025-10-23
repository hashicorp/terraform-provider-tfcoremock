// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-provider-tfcoremock/internal/data"
)

var _ Client = State{}

type State struct {
	DataDirectory string
}

func (state State) ReadResource(ctx context.Context, id string) (*data.Resource, error) {
	return nil, nil
}

func (state State) WriteResource(ctx context.Context, value *data.Resource) error {
	return nil
}

func (state State) UpdateResource(ctx context.Context, value *data.Resource) error {
	return nil
}

func (state State) DeleteResource(ctx context.Context, id string) error {
	return nil
}

func (state State) ReadDataSource(ctx context.Context, id string) (*data.Resource, error) {
	tflog.Trace(ctx, "Local.ReadDataSource")

	jsonPath := path.Join(state.DataDirectory, fmt.Sprintf("%s.json", id))

	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, err
	}

	var value data.Resource
	if err := json.Unmarshal(jsonData, &value); err != nil {
		return nil, err
	}

	return &value, nil
}

func (state State) ListResources(ctx context.Context, typeName *string, id *string, yield func(resource *data.Resource, err error), limit int64) error {
	return nil
}
