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

// ManifestFileNames are the declarations a package may carry, in the order they
// win.
//
// driver.json is the version 2 format: one file holding the roles, the
// capability blocks and the endpoints, replacing the four that could previously
// disagree with each other. manifest.json remains for version 1 packages, and
// for a version 2 package that also wants to install on a controller predating
// driver.json.
var ManifestFileNames = []string{"driver.json", "manifest.json"}

func LoadManifest(dir string) (*Manifest, error) {
	var lastErr error
	for _, name := range ManifestFileNames {
		data, err := os.ReadFile(dir + "/" + name)
		if err != nil {
			lastErr = err
			continue
		}
		return ParseManifest(data)
	}
	return nil, lastErr
}

func ParseManifest(data []byte) (*Manifest, error) {
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}
