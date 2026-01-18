# nx-driver-packager

Minimal packaging + verification CLI for Notrix driver bundles.

This tool produces a `.nxpkg` file, which is a `tar.gz` archive of a driver directory.
It also supports verifying a package (and optionally comparing it to a source directory).

## Repo Layout

- `cmd/nx-driver-packager/`: CLI entrypoint (`pack`, `verify`)
- `internal/packager/`: manifest parsing/validation, pack + verify logic
- `internal/archive/targz/`: tar writer used for packaging
- `driver-notrix-vort-dcdimmer-go/`: example driver directory

## Build

From repo root:

```powershell
go build -o .\nx-driver-packager.exe .\cmd\nx-driver-packager
```

## Usage

### Pack

Creates a `.nxpkg` (a gzipped tar of the input directory).

```powershell
.\nx-driver-packager.exe pack --input .\driver-notrix-vort-dcdimmer-go --out com.notrix.vort.dcdimmer-1.0.0.nxpkg
```

### Verify

Validates that:

- `manifest.json` exists in the package and parses
- required manifest fields exist
- the manifest entrypoint path exists in the archive

```powershell
.\nx-driver-packager.exe verify --pkg .\com.notrix.vort.dcdimmer-1.0.0.nxpkg
```

Optionally compare the package file list to a source directory (detect missing/extra files):

```powershell
.\nx-driver-packager.exe verify --pkg .\com.notrix.vort.dcdimmer-1.0.0.nxpkg --input .\driver-notrix-vort-dcdimmer-go
```

## Driver Directory Requirements

The `--input` directory is archived as-is (excluding symlinks; they are rejected).
At minimum it must contain:

- `manifest.json`
- the entrypoint file referenced by the manifest (e.g. `bin/driver`)

## Manifest Format

The packager supports two manifest styles.

### New format (recommended)

```json
{
	"id": "com.example.driver",
	"name": "Example Driver",
	"version": "1.0.0",
	"entrypoint": {
		"runtime": "go",
		"path": "bin/driver"
	}
}
```

### Legacy format (supported)

```json
{
	"driver_id": "com.example.driver",
	"name": "Example Driver",
	"version": "1.0.0",
	"runtime": "go",
	"entrypoint": "bin/driver"
}
```

### Fields

- `id` / `driver_id`: required (either one)
- `version`: required
- `entrypoint.path` (new) or `entrypoint` (legacy): required
- `entrypoint.runtime` (new) or `runtime` (legacy): optional (currently informational)

## Windows / PowerShell Notes

- To run a local executable from the current directory in PowerShell, prefix with `./` or `.\`:
	- `.\nx-driver-packager.exe ...`
- PowerShell does not use `\` for line continuation (use backtick `` ` `` or write the command on one line).

## Troubleshooting

- `error: manifest missing required fields`
	- Ensure `manifest.json` contains `id`/`driver_id`, `version`, and an entrypoint path.
- `error: entrypoint not found: ...`
	- Ensure the referenced entrypoint file exists under the input directory.
- `symlinks not allowed: ...`
	- Replace symlinks with real files before packaging.
