package dynamic

import (
	"encoding/json"
	"os"

	"github.com/hashicorp/terraform-provider-tfcoremock/internal/schema"
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
	data, err := os.ReadFile(r.File)
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
