package dynamic

import (
	"encoding/json"
	"os"
)

type Reader interface {
	Read() (*Resources, error)
}

type FileReader struct {
	File string
}

type StringReader struct {
	Data string
}

func (r FileReader) Read() (*Resources, error) {
	data, err := os.ReadFile(r.File)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var resources Resources
	if len(data) > 0 {
		if err := json.Unmarshal(data, &resources); err != nil {
			return nil, err
		}
	}

	return &resources, nil
}

func (r StringReader) Read() (*Resources, error) {
	var resources Resources
	if len(r.Data) > 0 {
		if err := json.Unmarshal([]byte(r.Data), &resources); err != nil {
			return nil, err
		}
	}
	return &resources, nil
}
