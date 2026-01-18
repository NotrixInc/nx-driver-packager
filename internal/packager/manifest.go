package packager

import (
	"encoding/json"
	"fmt"
	"os"
)

type Entrypoint struct {
	Runtime string `json:"runtime"`
	Path    string `json:"path"`
}

func (e *Entrypoint) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	// Legacy format: "entrypoint": "bin/driver"
	if data[0] == '"' {
		var path string
		if err := json.Unmarshal(data, &path); err != nil {
			return err
		}
		e.Path = path
		return nil
	}
	// New format: "entrypoint": { "runtime": "go", "path": "bin/driver" }
	var obj struct {
		Runtime string `json:"runtime"`
		Path    string `json:"path"`
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	if obj.Path == "" {
		return fmt.Errorf("entrypoint.path is required")
	}
	*e = Entrypoint(obj)
	return nil
}

type Manifest struct {
	// Support both "id" (new) and "driver_id" (legacy)
	ID       string `json:"id"`
	DriverID string `json:"driver_id"`

	Name    string `json:"name"`
	Version string `json:"version"`

	Entrypoint Entrypoint `json:"entrypoint"`

	// Legacy top-level runtime (new manifests put runtime under entrypoint)
	Runtime string `json:"runtime"`
}

func (m *Manifest) EffectiveID() string {
	if m.ID != "" {
		return m.ID
	}
	return m.DriverID
}

func (m *Manifest) EffectiveRuntime() string {
	if m.Entrypoint.Runtime != "" {
		return m.Entrypoint.Runtime
	}
	return m.Runtime
}

func (m *Manifest) EffectiveEntrypointPath() string {
	return m.Entrypoint.Path
}

func LoadManifest(dir string) (*Manifest, error) {
	data, err := os.ReadFile(dir + "/manifest.json")
	if err != nil {
		return nil, err
	}
	return ParseManifest(data)
}

func ParseManifest(data []byte) (*Manifest, error) {
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}
