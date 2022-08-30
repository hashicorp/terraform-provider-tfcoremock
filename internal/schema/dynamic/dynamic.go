package dynamic

import (
	"encoding/json"
	"github.com/hashicorp/terraform-provider-mock/internal/schema"
	"os"
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
	// TODO(liamcervante): Turn this into an environment variable?
	data, err := os.ReadFile("dynamic_resources.json")
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var dynamicResources map[string]schema.Schema
	if len(data) > 0 {
		if err := json.Unmarshal(data, &dynamicResources); err != nil {
			return nil, err
		}
	}

	return dynamicResources, nil
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
