package remote

import (
	"encoding/json"
	"os"
	"path"
)

const MetadataFile string = ".metadata"

type Metadata struct {
	Providers map[string]MetadataEntry `json:"providers"`
}

type MetadataEntry struct {
	Remote string `json:"remote"`
	Local  string `json:"local"`
}

func LoadMetadata(directory string) (Metadata, error) {
	metadata := Metadata{}

	data, err := os.ReadFile(path.Join(directory, MetadataFile))
	if err != nil {
		if os.IsNotExist(err) {
			metadata.Providers = make(map[string]MetadataEntry)
			return metadata, nil
		}
		return metadata, err
	}

	if err := json.Unmarshal(data, &metadata.Providers); err != nil {
		return metadata, err
	}
	return metadata, nil
}

func (metadata *Metadata) Save(directory string) error {
	data, err := json.Marshal(metadata.Providers)
	if err != nil {
		return err
	}

	file, err := os.Create(path.Join(directory, MetadataFile))
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		return err
	}
	return nil
}

func (metadata *Metadata) Set(key, remote, local string) {
	metadata.Providers[key] = MetadataEntry{remote, local}
}

func (metadata *Metadata) GetRemote(key string) string {
	return metadata.Providers[key].Remote
}

func (metadata *Metadata) GetLocal(key string) string {
	return metadata.Providers[key].Local
}
