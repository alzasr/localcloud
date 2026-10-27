package cli

const usage = `localcloud - a local cloud environment for development.

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
`

const versionUsage = `localcloud version - print version information about the binary.

Usage:
  localcloud version
`

const helpUsage = `localcloud help - print help for localcloud or one of its commands.

Usage:
  localcloud help [command]
`

const initUsage = `localcloud init - initialize a local cloud workspace.

Usage:
  localcloud init [directory]

Arguments:
  directory   Directory to initialize (default: current directory).

Flags:
  -h, --help   Print this help

The command is a stub: it validates its arguments and reports that workspace
initialization is not implemented yet.
`
