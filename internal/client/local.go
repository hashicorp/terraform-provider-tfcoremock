package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-provider-mock/internal/data"
	"os"
	"path"
)

var _ Client = Local{}

type Local struct {
	ResourceDirectory string
	DataDirectory     string
}

func (local Local) ReadResource(ctx context.Context, id string) (*data.Resource, error) {
	tflog.Trace(ctx, "Local.ReadResource")

	jsonPath := path.Join(local.ResourceDirectory, fmt.Sprintf("%s.json", id))

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

func (local Local) WriteResource(ctx context.Context, value *data.Resource) error {
	tflog.Trace(ctx, "Local.WriteResource")

	jsonData, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(local.ResourceDirectory, 0700); err != nil {
		return err
	}

	jsonPath := path.Join(local.ResourceDirectory, fmt.Sprintf("%s.json", value.GetId()))
	if _, err := os.Stat(jsonPath); err == nil {
		return errors.New("resource with the specified id likely already exists")
	}

	if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
		return err
	}

	return nil
}

func (local Local) UpdateResource(ctx context.Context, value *data.Resource) error {
	jsonData, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}

	jsonPath := path.Join(local.ResourceDirectory, fmt.Sprintf("%s.json", value.GetId()))
	if _, err := os.Stat(jsonPath); err != nil {
		return err
	}

	if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
		return err
	}

	return nil
}

func (local Local) DeleteResource(ctx context.Context, id string) error {
	jsonPath := path.Join(local.ResourceDirectory, fmt.Sprintf("%s.json", id))
	if err := os.Remove(jsonPath); err != nil {
		return err
	}

	return nil
}

func (local Local) ReadDataSource(ctx context.Context, id string) (*data.Resource, error) {
	jsonPath := path.Join(local.DataDirectory, fmt.Sprintf("%s.json", id))

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
