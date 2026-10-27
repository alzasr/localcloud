// Package cli is the terminal adapter of localcloud: it parses command line
// arguments, calls the user scenarios of the usecases package and turns their
// outcome into output and exit codes. It holds no business rules.
package cli

import (
	"errors"
	"fmt"
	"io"

	"github.com/alzasr/localcloud/internal/buildinfo"
	"github.com/alzasr/localcloud/internal/domain/workspace"
	"github.com/alzasr/localcloud/internal/usecases"
)

// Name is the name of the command as it is invoked from the shell.
const Name = "localcloud"

// Exit codes returned by Run.
const (
	exitOK      = 0
	exitFailure = 1
	exitUsage   = 2
)

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
	case "init":
		return runInit(rest, stdout, stderr)
	case "version", "-v", "--version":
		return runVersion(rest, stdout, stderr)
	default:
		fmt.Fprintf(stderr, "%s: unknown command %q\n\nRun '%s help' for usage.\n", Name, command, Name)
		return exitUsage
	}
}

// runInit implements the "init" command: it parses the command line arguments,
// hands the directory to the Init use case and reports the outcome.
func runInit(args []string, stdout, stderr io.Writer) int {
	if hasHelpFlag(args) {
		fmt.Fprint(stdout, initUsage)
		return exitOK
	}

	if len(args) > 1 {
		fmt.Fprintf(stderr, "%s init: unexpected argument %q\n", Name, args[1])
		return exitUsage
	}

	directory := "."
	if len(args) == 1 {
		directory = args[0]
	}

	result, err := usecases.Init{}.Execute(directory)
	if err != nil {
		return reportInitError(err, directory, stderr)
	}

	fmt.Fprintf(stdout, "%s init: initialized a local cloud workspace in %q\n", Name, result.Directory)

	return exitOK
}

// reportInitError writes the error of the init use case to stderr and returns
// the exit code: invalid input is a usage error, everything else is a failure.
func reportInitError(err error, directory string, stderr io.Writer) int {
	switch {
	case errors.Is(err, workspace.ErrEmptyDirectory):
		fmt.Fprintf(stderr, "%s init: %v\n", Name, err)
		return exitUsage
	case errors.Is(err, workspace.ErrNotImplemented):
		fmt.Fprintf(stderr, "%s init: creating a local cloud workspace in %q is not implemented yet\n", Name, directory)
		return exitFailure
	default:
		fmt.Fprintf(stderr, "%s init: %v\n", Name, err)
		return exitFailure
	}
}

// hasHelpFlag reports whether the arguments ask for the command help.
func hasHelpFlag(args []string) bool {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			return true
		}
	}

	return false
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
	case "init":
		text = initUsage
	case "version":
		text = versionUsage
	default:
		fmt.Fprintf(stderr, "%s help: unknown command %q\n\nRun '%s help' for usage.\n", Name, args[0], Name)
		return exitUsage
	}

	fmt.Fprint(stdout, text)

	return exitOK
}
