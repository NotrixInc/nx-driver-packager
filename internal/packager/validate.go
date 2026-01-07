
package packager

import (
	"fmt"
	"os"
	"path/filepath"
)

func Validate(dir string, m *Manifest) error {
	if m.DriverID == "" || m.Version == "" || m.Entrypoint == "" {
		return fmt.Errorf("manifest missing required fields")
	}
	ep := filepath.Join(dir, m.Entrypoint)
	if _, err := os.Stat(ep); err != nil {
		return fmt.Errorf("entrypoint not found: %s", m.Entrypoint)
	}
	return nil
}
