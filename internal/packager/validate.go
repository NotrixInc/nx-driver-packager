package packager

import (
	"fmt"
	"os"
	"path/filepath"
)

func Validate(dir string, m *Manifest) error {
	if m.EffectiveID() == "" || m.Version == "" || m.EffectiveEntrypointPath() == "" {
		return fmt.Errorf("manifest missing required fields")
	}
	ep := filepath.Join(dir, m.EffectiveEntrypointPath())
	if _, err := os.Stat(ep); err != nil {
		return fmt.Errorf("entrypoint not found: %s", m.EffectiveEntrypointPath())
	}
	return nil
}
