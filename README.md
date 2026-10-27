# localcloud

`localcloud` is a command line tool for running a local cloud environment for
development. This repository holds the initial Go project: a small CLI that can
be installed with `go install`; the only third-party dependency is a YAML
parser used for the workspace settings.

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

## Configuration

The settings of a workspace are stored in `localcloud.yaml` in the workspace
directory:

```yaml
infra:
  network:
    name: localcloud
    subnet: 172.28.0.0/16
  services:
    - name: postgres
      kind: database
      version: "16"
      port: 5432
    - name: rabbitmq
      kind: message_bus
      port: 5672

services:
  - name: api
    port: 8080
    depends_on:
      - postgres
```

The settings are grouped into two blocks:

- `infra` — the stateful services, such as databases and message buses, and the
  network they share (`name` and `subnet` in CIDR notation);
- `services` — the stateless services running on top of the infrastructure,
  where `depends_on` lists the infra services a service needs.

The `kind` of an infra service is `database` or `message_bus`; `version` is
optional and falls back to the default version of that service.

The settings are handled by the configurator — the business logic module other
modules get their settings from. It validates them instead of running with
inconsistent input: empty or malformed network settings, repeated service
names, unknown service kinds, ports outside 1–65535 and dependencies on
services that are not described are reported as errors. Unknown keys are
rejected too, so a typo does not silently disappear. A missing, empty or
unreadable file is reported separately from invalid settings.

The configurator does not know about the file: it works with a repository port,
and reading and writing `localcloud.yaml` is the job of the
`internal/adapters/configfile` adapter.

## Architecture

The project follows clean architecture: layers are separated and dependencies
point inwards only — adapters know about use cases, use cases know about the
business logic, never the other way round.

| Layer | Package | Responsibility |
|-------|---------|----------------|
| Adapters | `internal/adapters/cli`, `internal/adapters/configfile` | Terminal client (argument parsing, help, output formatting, exit codes) and the YAML storage of the workspace settings. |
| Use cases | `internal/usecases` | User scenarios behind commands, for now `init`: they orchestrate the actions of the business logic. |
| Business logic | `internal/domain/workspace`, `internal/domain/config` | Rules of a local cloud workspace, the actions performed on it and the configurator that keeps the settings for the other modules. |

Rules:

- adapters contain no business rules and never decide what a valid workspace or
  a valid setting is;
- use cases do not parse arguments, do not print to the terminal and do not
  return exit codes;
- business logic prints nothing, knows neither adapters nor use cases, and
  declares the ports (such as the settings repository) that adapters implement;
- a new command becomes a use case that the terminal adapter calls; a new
  action belongs to the business logic and is called by a use case; a new way
  to store data becomes an adapter implementing a port of the business logic;
- adapters turn errors of the business logic into user messages and exit codes,
  so they import the business logic packages to tell the errors apart;
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
    │   ├── cli
    │   │   ├── cli.go               # terminal adapter: parsing, help, output, exit codes
    │   │   └── usage.go             # help texts of the commands
    │   └── configfile
    │       └── configfile.go        # localcloud.yaml storage of the workspace settings
    ├── buildinfo
    │   └── buildinfo.go             # version information of the binary
    ├── domain
    │   ├── config
    │   │   ├── config.go            # settings model: infra and services blocks
    │   │   ├── configurator.go      # module that keeps the settings, repository port
    │   │   ├── errors.go            # errors describing invalid settings
    │   │   └── validate.go          # rules the settings must follow
    │   └── workspace
    │       └── workspace.go         # business logic of a local cloud workspace
    └── usecases
        └── init.go                  # "init" user scenario
```
