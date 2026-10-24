package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/alzasr/localcloud/internal/buildinfo"
	"github.com/alzasr/localcloud/internal/cli"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantCode   int
		wantStdout string
		wantStderr string
	}{
		{
			name:       "no arguments prints usage",
			args:       nil,
			wantCode:   0,
			wantStdout: "Usage:",
		},
		{
			name:       "help prints usage",
			args:       []string{"help"},
			wantCode:   0,
			wantStdout: "localcloud <command>",
		},
		{
			name:       "short help flag prints usage",
			args:       []string{"-h"},
			wantCode:   0,
			wantStdout: "Commands:",
		},
		{
			name:       "long help flag prints usage",
			args:       []string{"--help"},
			wantCode:   0,
			wantStdout: "Commands:",
		},
		{
			name:       "help for version",
			args:       []string{"help", "version"},
			wantCode:   0,
			wantStdout: "localcloud version",
		},
		{
			name:       "help for help",
			args:       []string{"help", "help"},
			wantCode:   0,
			wantStdout: "localcloud help",
		},
		{
			name:       "version command",
			args:       []string{"version"},
			wantCode:   0,
			wantStdout: cli.Name + " ",
		},
		{
			name:       "short version flag",
			args:       []string{"-v"},
			wantCode:   0,
			wantStdout: cli.Name + " ",
		},
		{
			name:       "long version flag",
			args:       []string{"--version"},
			wantCode:   0,
			wantStdout: cli.Name + " ",
		},
		{
			name:       "unknown command",
			args:       []string{"deploy"},
			wantCode:   2,
			wantStderr: `unknown command "deploy"`,
		},
		{
			name:       "empty command",
			args:       []string{""},
			wantCode:   2,
			wantStderr: `unknown command ""`,
		},
		{
			name:       "version rejects extra arguments",
			args:       []string{"version", "extra"},
			wantCode:   2,
			wantStderr: `unexpected argument "extra"`,
		},
		{
			name:       "help rejects unknown command",
			args:       []string{"help", "deploy"},
			wantCode:   2,
			wantStderr: `unknown command "deploy"`,
		},
		{
			name:       "help rejects extra arguments",
			args:       []string{"help", "version", "extra"},
			wantCode:   2,
			wantStderr: `unexpected argument "extra"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer

			gotCode := cli.Run(tt.args, &stdout, &stderr)
			if gotCode != tt.wantCode {
				t.Errorf("Run(%q) exit code = %d, want %d", tt.args, gotCode, tt.wantCode)
			}

			assertOutput(t, "stdout", stdout.String(), tt.wantStdout)
			assertOutput(t, "stderr", stderr.String(), tt.wantStderr)
		})
	}
}

func TestRunVersionPrintsBuildInfo(t *testing.T) {
	var stdout, stderr bytes.Buffer

	if code := cli.Run([]string{"version"}, &stdout, &stderr); code != 0 {
		t.Fatalf("Run(version) exit code = %d, want 0 (stderr: %q)", code, stderr.String())
	}

	want := cli.Name + " " + buildinfo.Summary()
	if got := strings.TrimSpace(stdout.String()); got != want {
		t.Errorf("version output = %q, want %q", got, want)
	}
	if stderr.Len() != 0 {
		t.Errorf("unexpected stderr: %q", stderr.String())
	}
}

// assertOutput verifies that got contains want, or is empty when want is empty.
func assertOutput(t *testing.T, stream, got, want string) {
	t.Helper()

	if want == "" {
		if got != "" {
			t.Errorf("unexpected %s: %q", stream, got)
		}
		return
	}

	if !strings.Contains(got, want) {
		t.Errorf("%s = %q, want it to contain %q", stream, got, want)
	}
}
