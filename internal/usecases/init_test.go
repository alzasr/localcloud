package usecases_test

import (
	"errors"
	"testing"

	"github.com/alzasr/localcloud/internal/domain/workspace"
	"github.com/alzasr/localcloud/internal/usecases"
)

func TestInitExecute(t *testing.T) {
	tests := []struct {
		name          string
		directory     string
		wantDirectory string
		wantErr       error
	}{
		{
			name:          "default directory",
			directory:     ".",
			wantDirectory: ".",
			wantErr:       workspace.ErrNotImplemented,
		},
		{
			name:          "explicit directory",
			directory:     "env",
			wantDirectory: "env",
			wantErr:       workspace.ErrNotImplemented,
		},
		{
			name:      "empty directory",
			directory: "",
			wantErr:   workspace.ErrEmptyDirectory,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := usecases.Init{}.Execute(tt.directory)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Execute(%q) error = %v, want %v", tt.directory, err, tt.wantErr)
			}
			if got.Directory != tt.wantDirectory {
				t.Errorf("Execute(%q) directory = %q, want %q", tt.directory, got.Directory, tt.wantDirectory)
			}
		})
	}
}
