// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dynamic

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/xeipuuv/gojsonschema"

	"github.com/hashicorp/terraform-provider-tfcoremock/internal/schema"
)

const (
	dynamicResourcesSchemaEnvVarName = "TFCOREMOCK_DYNAMIC_RESOURCES_SCHEMA"
)

type Reader interface {
	Read() (map[string]schema.Schema, error)
}

type FileReader struct {
	File string
}

type StringReader struct {
	Data string
}

func (r FileReader) Read() (map[string]schema.Schema, error) {
	dynamicResourcesSchema := "https://raw.githubusercontent.com/hashicorp/terraform-provider-tfcoremock/main/schema/dynamic_resources.json"
	if dynamicResourcesSchemaEnvVar := os.Getenv(dynamicResourcesSchemaEnvVarName); dynamicResourcesSchemaEnvVar != "" {
		dynamicResourcesSchema = dynamicResourcesSchemaEnvVar
	}

	schemaLoader := gojsonschema.NewReferenceLoader(dynamicResourcesSchema)

	data, err := os.ReadFile(r.File)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to read dynamic resources file")
	}

	documentLoader := gojsonschema.NewStringLoader(string(data))
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return nil, err
	}

	if result.Valid() {
		var dynamicResources map[string]schema.Schema
		if len(data) > 0 {
			if err := json.Unmarshal(data, &dynamicResources); err != nil {
				return nil, errors.Wrap(err, "failed to unmarshal dynamic resources json")
			}
		}

		return dynamicResources, nil
	}

	var errs []string
	for _, err := range result.Errors() {
		errs = append(errs, err.String())
	}

	return nil, fmt.Errorf("failed json schema check: %s", strings.Join(errs, ", "))
}

func (r StringReader) Read() (map[string]schema.Schema, error) {
	var dynamicResources map[string]schema.Schema
	if len(r.Data) > 0 {
		if err := json.Unmarshal([]byte(r.Data), &dynamicResources); err != nil {
			return nil, err
		}
	}
	return dynamicResources, nil
}
