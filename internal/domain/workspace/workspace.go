// Package workspace holds the business logic of a local cloud workspace: the
// rules that apply to a workspace directory and the actions that can be
// performed on it.
//
// The package knows nothing about the command line interface, about user
// scenarios that call it or about the terminal. It returns descriptive errors
// and never prints anything.
package workspace

import "errors"

// Errors returned by the package. They describe the violated rule; user
// interface wording is chosen by the adapters.
var (
	// ErrEmptyDirectory reports that a workspace cannot be created for an
	// empty directory.
	ErrEmptyDirectory = errors.New("directory must not be empty")
	// ErrNotImplemented reports that the requested action has no
	// implementation yet.
	ErrNotImplemented = errors.New("workspace initialization is not implemented")
)

// Workspace is a directory that holds a local cloud environment.
type Workspace struct {
	directory string
}

// New returns the workspace for directory after validating it.
func New(directory string) (Workspace, error) {
	if directory == "" {
		return Workspace{}, ErrEmptyDirectory
	}

	return Workspace{directory: directory}, nil
}

// Directory returns the directory of the workspace.
func (w Workspace) Directory() string {
	return w.directory
}

// Initialize creates the local cloud environment in the workspace directory.
func (w Workspace) Initialize() error {
	return ErrNotImplemented
}
