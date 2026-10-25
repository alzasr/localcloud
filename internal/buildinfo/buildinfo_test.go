package buildinfo

import (
	"runtime"
	"strings"
	"testing"
)

func TestVersionIsSet(t *testing.T) {
	if got := Version(); got == "" {
		t.Fatal("Version() is empty")
	}
}

func TestCommitIsSet(t *testing.T) {
	if got := Commit(); got == "" {
		t.Fatal("Commit() is empty")
	}
}

func TestBuildDateIsSet(t *testing.T) {
	if got := BuildDate(); got == "" {
		t.Fatal("BuildDate() is empty")
	}
}

func TestGoVersion(t *testing.T) {
	if got, want := GoVersion(), runtime.Version(); got != want {
		t.Errorf("GoVersion() = %q, want %q", got, want)
	}
}

func TestPlatform(t *testing.T) {
	want := runtime.GOOS + "/" + runtime.GOARCH
	if got := Platform(); got != want {
		t.Errorf("Platform() = %q, want %q", got, want)
	}
}

func TestSummary(t *testing.T) {
	got := Summary()

	for _, want := range []string{Version(), Commit(), BuildDate(), GoVersion(), Platform()} {
		if !strings.Contains(got, want) {
			t.Errorf("Summary() = %q, want it to contain %q", got, want)
		}
	}
}
