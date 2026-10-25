// Package cli implements the command line interface of the localcloud binary.
package cli

import (
	"fmt"
	"io"

	"github.com/alzasr/localcloud/internal/buildinfo"
)

// Name is the name of the command as it is invoked from the shell.
const Name = "localcloud"

// Exit codes returned by Run.
const (
	exitOK    = 0
	exitUsage = 2
)

const usage = `localcloud - a local cloud environment for development.

Usage:
  localcloud <command>

Commands:
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

// Run executes the localcloud CLI with the given arguments, without the
// program name, writing regular output to stdout and errors to stderr. It
// returns the process exit code.
func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stdout, usage)
		return exitOK
	}

	command, rest := args[0], args[1:]
	switch command {
	case "help", "-h", "--help":
		return runHelp(rest, stdout, stderr)
	case "version", "-v", "--version":
		return runVersion(rest, stdout, stderr)
	default:
		fmt.Fprintf(stderr, "%s: unknown command %q\n\nRun '%s help' for usage.\n", Name, command, Name)
		return exitUsage
	}
}

// runVersion implements the "version" command and its flags.
func runVersion(args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 {
		fmt.Fprintf(stderr, "%s version: unexpected argument %q\n", Name, args[0])
		return exitUsage
	}

	fmt.Fprintf(stdout, "%s %s\n", Name, buildinfo.Summary())

	return exitOK
}

// runHelp implements the "help" command and its flags.
func runHelp(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stdout, usage)
		return exitOK
	}

	if len(args) > 1 {
		fmt.Fprintf(stderr, "%s help: unexpected argument %q\n", Name, args[1])
		return exitUsage
	}

	var text string
	switch args[0] {
	case "help":
		text = helpUsage
	case "version":
		text = versionUsage
	default:
		fmt.Fprintf(stderr, "%s help: unknown command %q\n\nRun '%s help' for usage.\n", Name, args[0], Name)
		return exitUsage
	}

	fmt.Fprint(stdout, text)

	return exitOK
}
