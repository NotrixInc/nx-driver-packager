
package packager

import (
	"compress/gzip"
	"os"

	"github.com/NotrixInc/nx-driver-packager/internal/archive/targz"
)

func Pack(inputDir, outFile string) error {
	m, err := LoadManifest(inputDir)
	if err != nil {
		return err
	}
	if err := Validate(inputDir, m); err != nil {
		return err
	}

	f, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer f.Close()

	gz := gzip.NewWriter(f)
	defer gz.Close()

	return targz.WriteDir(inputDir, gz)
}
