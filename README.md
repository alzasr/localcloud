# localcloud

`localcloud` is a command line tool for running a local cloud environment for
development. This repository holds the initial Go project: a small,
dependency-free CLI that can be installed with `go install`.

## Requirements

- Go 1.24 or newer

## Installation

Install the latest version straight from the repository:

```sh
go install github.com/alzasr/localcloud@latest
```

Go resolves `@latest` to the newest release tag, or to the latest commit on the
default branch (pseudo-version) when no tag exists yet. To install a concrete
release:

```sh
go install github.com/alzasr/localcloud@v0.1.0
```

The binary is placed in `$(go env GOPATH)/bin` (usually `~/go/bin`). Make sure
that directory is on your `PATH`:

```sh
export PATH="$(go env GOPATH)/bin:$PATH"
```

### Build from a checkout

```sh
git clone git@github.com:alzasr/localcloud.git
cd localcloud
go build -o bin/localcloud .
./bin/localcloud --help
```

## Usage

```text
$ localcloud help
localcloud - a local cloud environment for development.

Usage:
  localcloud <command>

Commands:
  version   Print version information about the binary
  help      Print this help

Flags:
  -h, --help      Print this help
  -v, --version   Print version information

Use "localcloud help <command>" for more information about a command.
```

```text
$ localcloud version
localcloud v0.1.0 (commit 1a2b3c4d5e6f, built 2025-01-02T15:04:05Z, go1.24.0, linux/amd64)
```

The `version` output is derived from the build information embedded by the Go
toolchain. A binary installed with `go install ...@vX.Y.Z` reports that module
version; a build inside a Git checkout reports a pseudo-version of the current
commit (with a `+dirty` suffix when the work tree has uncommitted changes), and
a build without any VCS information reports `dev`.

## Development

```sh
go test ./...
go vet ./...
```

Project tasks are tracked as Markdown files in [`issues`](issues/README.md);
the naming and body conventions are described in
[`issues/README.md`](issues/README.md).

### Version stamping

You can also set the build metadata explicitly with linker flags:

```sh
go build -ldflags "\
  -X github.com/alzasr/localcloud/internal/buildinfo.version=$(git describe --tags --always) \
  -X github.com/alzasr/localcloud/internal/buildinfo.commit=$(git rev-parse --short HEAD) \
  -X github.com/alzasr/localcloud/internal/buildinfo.date=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  -o bin/localcloud .
```

### Project layout

```text
.
├── main.go                      # entry point, wires os.Args to the CLI package
├── go.mod
└── internal
    ├── buildinfo
    │   └── buildinfo.go         # version information of the binary
    └── cli
        └── cli.go               # command parsing, help and version output
```
