// Package usecases holds the user scenarios of the localcloud interface. A
// scenario describes what a command does by orchestrating the actions of the
// business logic; it neither parses arguments nor writes to the terminal.
package usecases

import "github.com/alzasr/localcloud/internal/domain/workspace"

// Init is the use case behind the "init" command: it prepares a local cloud
// workspace in the requested directory.
type Init struct{}

// InitResult describes the outcome of the Init use case.
type InitResult struct {
	// Directory is the workspace directory the scenario worked with.
	Directory string
}

// Execute validates the directory and initializes the workspace there.
//
// The returned error is the error of the business logic, so callers can tell
// invalid input from a failed action with errors.Is.
func (Init) Execute(directory string) (InitResult, error) {
	ws, err := workspace.New(directory)
	if err != nil {
		return InitResult{}, err
	}

	result := InitResult{Directory: ws.Directory()}
	if err := ws.Initialize(); err != nil {
		return result, err
	}

	return result, nil
}
