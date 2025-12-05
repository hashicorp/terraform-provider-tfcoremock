// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-provider-tfcoremock/internal/data"
)

var _ Client = Local{}

type Local struct {
	ResourceDirectory string
	DataDirectory     string
}

func (local Local) ReadResource(ctx context.Context, id string) (*data.Resource, error) {
	tflog.Trace(ctx, "Local.ReadResource")

	jsonPath := filepath.Join(local.ResourceDirectory, fmt.Sprintf("%s.json", id))

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

	jsonPath := filepath.Join(local.ResourceDirectory, fmt.Sprintf("%s.json", value.GetId()))

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

	jsonPath := filepath.Join(local.ResourceDirectory, fmt.Sprintf("%s.json", value.GetId()))
	if _, err := os.Stat(jsonPath); err != nil {
		return err
	}

	if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
		return err
	}

	return nil
}

func (local Local) DeleteResource(ctx context.Context, id string) error {
	jsonPath := filepath.Join(local.ResourceDirectory, fmt.Sprintf("%s.json", id))
	if err := os.Remove(jsonPath); err != nil {
		return err
	}

	// If the directory is empty after we've deleted this resource, let's tidy
	// up and delete the directory as well.
	resources, err := os.Open(local.ResourceDirectory)
	if err != nil {
		// Something weird has happened, but we're not going to fail the whole
		// delete operation just cos we couldn't clean up the directory.
		tflog.Info(ctx, fmt.Sprintf("couldn't open resource directory at (%s) to tidy up: %v", local.ResourceDirectory, err))
		return nil
	}

	files, err := resources.Readdirnames(1)
	if len(files) > 0 {
		// Then we're not going to do anything, there are still other files or
		// resources within this directory.
		return nil
	}

	if err == io.EOF {
		// Then we returned an empty slice of files because the directory is
		// empty - let's delete the directory then. This is an acceptable
		// outcome, so we're not going to log anything.
		_ = os.Remove(local.ResourceDirectory)
		return nil
	}

	// Then something else caused us to return an empty slice. We'll be cautious
	// and log the error but not delete the directory.
	tflog.Info(ctx, fmt.Sprintf("failed to query if the resource directory at (%s) was empty: %v", local.ResourceDirectory, err))
	return nil
}

func (local Local) ReadDataSource(ctx context.Context, id string) (*data.Resource, error) {
	jsonPath := filepath.Join(local.DataDirectory, fmt.Sprintf("%s.json", id))

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

func (local Local) ListResources(ctx context.Context, typeName *string, id *string, yield func(resource *data.Resource, err error), limit int64) error {
	if id != nil {
		yield(local.ReadResource(ctx, *id))
		return nil
	}

	entries, err := os.ReadDir(local.ResourceDirectory)
	if err != nil {
		return err
	}

	var count int64
	for _, entry := range entries {
		if count == limit {
			break // only yield the exact number of responses
		}

		if entry.IsDir() {
			continue // no nested directories
		}

		ext := filepath.Ext(entry.Name())
		if ext != ".json" {
			continue // only read the json files
		}

		jsonData, err := os.ReadFile(filepath.Join(local.ResourceDirectory, entry.Name()))
		if err != nil {
			count++
			yield(nil, fmt.Errorf("failed to read %s: %w", entry.Name(), err))
			continue
		}

		var value data.Resource
		if err := json.Unmarshal(jsonData, &value); err != nil {
			count++
			yield(nil, fmt.Errorf("failed to unmarshal %s: %w", entry.Name(), err))
			continue
		}

		if typeName != nil && value.ResourceType != *typeName {
			continue // wrong type
		}

		count++
		yield(&value, nil)
	}

	return nil
}
