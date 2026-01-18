package packager

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type packageIndex struct {
	files map[string]tar.Header
}

func indexPackage(pkgFile string) (*packageIndex, []byte, error) {
	f, err := os.Open(pkgFile)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, nil, err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	idx := &packageIndex{files: make(map[string]tar.Header)}

	var manifest []byte

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}

		name := strings.TrimPrefix(hdr.Name, "./")
		name = filepath.ToSlash(name)
		if name == "" {
			continue
		}

		switch hdr.Typeflag {
		case tar.TypeReg, tar.TypeRegA:
			idx.files[name] = *hdr
			if name == "manifest.json" {
				b, err := io.ReadAll(tr)
				if err != nil {
					return nil, nil, err
				}
				manifest = b
			}
		default:
			// ignore dirs and other types
		}
	}

	return idx, manifest, nil
}

func listInputFiles(inputDir string) (map[string]struct{}, error) {
	files := make(map[string]struct{})
	root, err := filepath.Abs(inputDir)
	if err != nil {
		return nil, err
	}

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == root {
			return nil
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)

		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("symlinks not allowed: %s", rel)
		}

		if info.IsDir() {
			return nil
		}

		files[rel] = struct{}{}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func VerifyPackage(pkgFile string, inputDir string) error {
	idx, manifestBytes, err := indexPackage(pkgFile)
	if err != nil {
		return err
	}
	if manifestBytes == nil {
		return fmt.Errorf("manifest.json not found in package")
	}

	m, err := ParseManifest(manifestBytes)
	if err != nil {
		return fmt.Errorf("invalid manifest.json: %w", err)
	}

	if m.EffectiveID() == "" || m.Version == "" || m.EffectiveEntrypointPath() == "" {
		return fmt.Errorf("manifest missing required fields")
	}

	ep := filepath.ToSlash(strings.TrimPrefix(m.EffectiveEntrypointPath(), "./"))
	if _, ok := idx.files[ep]; !ok {
		return fmt.Errorf("entrypoint not found in package: %s", m.EffectiveEntrypointPath())
	}

	if inputDir != "" {
		inFiles, err := listInputFiles(inputDir)
		if err != nil {
			return err
		}

		// Compare file paths only (regular files). Package index already tracks regular files.
		var missing []string
		var extra []string

		for f := range inFiles {
			if _, ok := idx.files[f]; !ok {
				missing = append(missing, f)
			}
		}
		for f := range idx.files {
			if _, ok := inFiles[f]; !ok {
				extra = append(extra, f)
			}
		}

		if len(missing) > 0 || len(extra) > 0 {
			// Keep output short but actionable.
			errMsg := "package contents mismatch"
			if len(missing) > 0 {
				errMsg += fmt.Sprintf("; missing %d files", len(missing))
			}
			if len(extra) > 0 {
				errMsg += fmt.Sprintf("; extra %d files", len(extra))
			}
			return fmt.Errorf(errMsg)
		}
	}

	return nil
}
