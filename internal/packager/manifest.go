
package packager

import (
	"encoding/json"
	"os"
)

type Manifest struct {
	DriverID   string `json:"driver_id"`
	Name       string `json:"name"`
	Version    string `json:"version"`
	Entrypoint string `json:"entrypoint"`
	Runtime    string `json:"runtime"`
}

func LoadManifest(dir string) (*Manifest, error) {
	data, err := os.ReadFile(dir + "/manifest.json")
	if err != nil {
		return nil, err
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}
