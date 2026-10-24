// Package buildinfo exposes information about how the localcloud binary was
// built.
//
// Values stamped by the Go toolchain (module version, VCS revision and build
// time) are used when available, so binaries installed with
// "go install github.com/alzasr/localcloud@latest" report a meaningful
// version. Local builds can override the values with -ldflags.
package buildinfo

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

// The variables below are set at link time, for example:
//
//	go build -ldflags "\
//		-X github.com/alzasr/localcloud/internal/buildinfo.version=v0.1.0 \
//		-X github.com/alzasr/localcloud/internal/buildinfo.commit=1a2b3c4 \
//		-X github.com/alzasr/localcloud/internal/buildinfo.date=2025-01-02T15:04:05Z"
var (
	version = ""
	commit  = ""
	date    = ""
)

const (
	// devVersion is reported for local builds that carry no version stamp.
	devVersion = "dev"
	// unknownValue is reported when a build property is not available.
	unknownValue = "unknown"
	// shortCommitLen is the length of the abbreviated VCS revision.
	shortCommitLen = 12
)

// Version returns the module version of the binary, or "dev" for local builds
// that are not stamped with a version.
func Version() string {
	if version != "" {
		return version
	}

	info, ok := debug.ReadBuildInfo()
	if !ok || info.Main.Version == "" || info.Main.Version == "(devel)" {
		return devVersion
	}

	return info.Main.Version
}

// Commit returns the abbreviated VCS revision the binary was built from.
func Commit() string {
	if commit != "" {
		return commit
	}

	revision, ok := setting("vcs.revision")
	if !ok {
		return unknownValue
	}
	if len(revision) > shortCommitLen {
		return revision[:shortCommitLen]
	}

	return revision
}

// BuildDate returns the time the binary was built in RFC 3339 format.
func BuildDate() string {
	if date != "" {
		return date
	}

	buildTime, ok := setting("vcs.time")
	if !ok {
		return unknownValue
	}

	return buildTime
}

// GoVersion returns the version of the Go toolchain used to build the binary.
func GoVersion() string {
	return runtime.Version()
}

// Platform returns the target platform in "os/arch" form.
func Platform() string {
	return runtime.GOOS + "/" + runtime.GOARCH
}

// Summary returns a single line description of the build.
func Summary() string {
	return fmt.Sprintf(
		"%s (commit %s, built %s, %s, %s)",
		Version(),
		Commit(),
		BuildDate(),
		GoVersion(),
		Platform(),
	)
}

// setting looks up a build setting, such as "vcs.revision", in the build
// information embedded by the Go toolchain.
func setting(key string) (string, bool) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "", false
	}

	for _, s := range info.Settings {
		if s.Key == key && s.Value != "" {
			return s.Value, true
		}
	}

	return "", false
}
