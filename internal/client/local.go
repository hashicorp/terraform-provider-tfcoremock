package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-provider-mock/internal/data"
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

	// Let's just do a quick sanity check here. We are expecting the stat to
	// return an os.IsNotExist error, we want to make sure we aren't trying to
	// create a resource that already exists. If we don't get an error then that
	// means we are trying to overwrite a resource when we shouldn't, and if we
	// get anything other than an os.IsNotExist error then something even
	// weirder is happening.
	if _, err := os.Stat(jsonPath); err == nil {
		return errors.New("resource with the specified id likely already exists")
	} else if err != nil && !os.IsNotExist(err) {
		return err
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
