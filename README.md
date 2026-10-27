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
  init      Initialize a local cloud workspace (not implemented yet)
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

```text
$ localcloud init ./workspace
localcloud init: creating a local cloud workspace in "./workspace" is not implemented yet
```

The `init` command is a stub: it validates its arguments, reports the directory
it would initialize and exits with code 1. Run `localcloud help init` (or
`localcloud init --help`) for its syntax. Workspace scaffolding itself will be
added in a follow-up task.

## Architecture

The project follows clean architecture: layers are separated and dependencies
point inwards only — adapters know about use cases, use cases know about the
business logic, never the other way round.

| Layer | Package | Responsibility |
|-------|---------|----------------|
| Adapters | `internal/adapters/cli` | Terminal client: argument parsing, help, output formatting and exit codes. |
| Use cases | `internal/usecases` | User scenarios behind commands, for now `init`: they orchestrate the actions of the business logic. |
| Business logic | `internal/domain/workspace` | Rules of a local cloud workspace and the concrete actions performed on it. |

Rules:

- adapters contain no business rules and never decide what a valid workspace is;
- use cases do not parse arguments, do not print to the terminal and do not
  return exit codes;
- business logic prints nothing and knows neither adapters nor use cases;
- a new command becomes a use case that the terminal adapter calls; a new
  action belongs to the business logic and is called by a use case;
- the terminal adapter turns errors of the business logic into user messages
  and exit codes, so it imports `internal/domain/workspace` to tell them apart;
- dependencies always point inwards: adapters → use cases → business logic.

`internal/buildinfo` is cross-cutting infrastructure that describes how the
binary was built; the terminal adapter uses it to print `localcloud version`.

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
├── main.go                          # entry point, wires os.Args to the CLI adapter
├── go.mod
└── internal
    ├── adapters
    │   └── cli
    │       ├── cli.go               # terminal adapter: parsing, help, output, exit codes
    │       └── usage.go             # help texts of the commands
    ├── buildinfo
    │   └── buildinfo.go             # version information of the binary
    ├── domain
    │   └── workspace
    │       └── workspace.go         # business logic of a local cloud workspace
    └── usecases
        └── init.go                  # "init" user scenario
```
