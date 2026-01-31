package packager

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func validateOptionalJSONFile(dir string, fileName string) error {
	p := filepath.Join(dir, fileName)
	b, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read %s: %w", fileName, err)
	}
	var anyJSON any
	if err := json.Unmarshal(b, &anyJSON); err != nil {
		return fmt.Errorf("invalid %s: %w", fileName, err)
	}
	return nil
}

func Validate(dir string, m *Manifest) error {
	if m.EffectiveID() == "" || m.Version == "" || m.EffectiveEntrypointPath() == "" {
		return fmt.Errorf("manifest missing required fields")
	}
	ep := filepath.Join(dir, m.EffectiveEntrypointPath())
	if _, err := os.Stat(ep); err != nil {
		return fmt.Errorf("entrypoint not found: %s", m.EffectiveEntrypointPath())
	}

	// Optional (backward compatible) driver metadata specs.
	// New templates include these; older drivers may not.
	if err := validateOptionalJSONFile(dir, "capabilities.json"); err != nil {
		return err
	}
	if err := validateOptionalJSONFile(dir, "events.json"); err != nil {
		return err
	}
	return nil
}
