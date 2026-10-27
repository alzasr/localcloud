package workspace_test

import (
	"errors"
	"testing"

	"github.com/alzasr/localcloud/internal/domain/workspace"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		directory string
		wantDir   string
		wantErr   error
	}{
		{
			name:      "empty directory",
			directory: "",
			wantErr:   workspace.ErrEmptyDirectory,
		},
		{
			name:      "current directory",
			directory: ".",
			wantDir:   ".",
		},
		{
			name:      "relative directory",
			directory: "env",
			wantDir:   "env",
		},
		{
			name:      "absolute directory",
			directory: "/tmp/env",
			wantDir:   "/tmp/env",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ws, err := workspace.New(tt.directory)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("New(%q) error = %v, want %v", tt.directory, err, tt.wantErr)
			}
			if got := ws.Directory(); got != tt.wantDir {
				t.Errorf("Directory() = %q, want %q", got, tt.wantDir)
			}
		})
	}
}

func TestWorkspaceInitializeIsNotImplemented(t *testing.T) {
	ws, err := workspace.New("env")
	if err != nil {
		t.Fatalf("New(env) error = %v, want nil", err)
	}

	if err := ws.Initialize(); !errors.Is(err, workspace.ErrNotImplemented) {
		t.Errorf("Initialize() error = %v, want %v", err, workspace.ErrNotImplemented)
	}
}
