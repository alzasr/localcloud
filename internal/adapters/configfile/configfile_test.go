package configfile_test

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/alzasr/localcloud/internal/adapters/configfile"
	"github.com/alzasr/localcloud/internal/domain/config"
)

// settings returns the settings used by the tests. Every call returns fresh
// slices.
func settings() config.Config {
	return config.Config{
		Infra: config.Infra{
			Network: config.Network{Name: "localcloud", Subnet: "172.28.0.0/16"},
			Services: []config.InfraService{
				{Name: "postgres", Kind: config.KindDatabase, Version: "16", Port: 5432},
				{Name: "rabbitmq", Kind: config.KindMessageBus, Port: 5672},
			},
		},
		Services: []config.Service{
			{Name: "api", Port: 8080, DependsOn: []string{"postgres"}},
		},
	}
}

// example is the settings file as it is described in the documentation.
const example = `infra:
  network:
    name: localcloud
    subnet: 172.28.0.0/16
  services:
    - name: postgres
      kind: database
      version: "16"
      port: 5432
    - name: rabbitmq
      kind: message_bus
      port: 5672
services:
  - name: api
    port: 8080
    depends_on:
      - postgres
`

func TestRepositorySaveAndLoad(t *testing.T) {
	directory := t.TempDir()
	repository := configfile.New(directory)

	wantPath := filepath.Join(directory, configfile.FileName)
	if got := repository.Path(); got != wantPath {
		t.Errorf("Path() = %q, want %q", got, wantPath)
	}

	if err := repository.Save(settings()); err != nil {
		t.Fatalf("Save() error = %v, want nil", err)
	}

	loaded, err := repository.Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}
	if !reflect.DeepEqual(loaded, settings()) {
		t.Errorf("Load() = %+v, want %+v", loaded, settings())
	}
}

// TestRepositoryLoadExample reads the settings file in the documented format.
func TestRepositoryLoadExample(t *testing.T) {
	repository := configfile.New(t.TempDir())
	writeFile(t, repository.Path(), example)

	loaded, err := repository.Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}
	if !reflect.DeepEqual(loaded, settings()) {
		t.Errorf("Load() = %+v, want %+v", loaded, settings())
	}
}

// TestRepositorySaveWritesDocumentedFormat verifies that Save writes the file
// exactly in the format shown in the documentation.
func TestRepositorySaveWritesDocumentedFormat(t *testing.T) {
	repository := configfile.New(t.TempDir())
	if err := repository.Save(settings()); err != nil {
		t.Fatalf("Save() error = %v, want nil", err)
	}

	data, err := os.ReadFile(repository.Path())
	if err != nil {
		t.Fatalf("read settings file: %v", err)
	}
	if string(data) != example {
		t.Errorf("settings file = %q, want %q", data, example)
	}
}

func TestRepositoryLoadMissingFile(t *testing.T) {
	repository := configfile.New(t.TempDir())

	_, err := repository.Load()
	if !errors.Is(err, config.ErrNotConfigured) {
		t.Fatalf("Load() error = %v, want %v", err, config.ErrNotConfigured)
	}
	if !strings.Contains(err.Error(), configfile.FileName) {
		t.Errorf("Load() error = %q, want it to mention %q", err, configfile.FileName)
	}
}

func TestRepositoryLoadEmptyFile(t *testing.T) {
	repository := configfile.New(t.TempDir())
	writeFile(t, repository.Path(), "")

	if _, err := repository.Load(); !errors.Is(err, config.ErrNotConfigured) {
		t.Errorf("Load() error = %v, want %v", err, config.ErrNotConfigured)
	}
}

func TestRepositoryLoadMalformedFile(t *testing.T) {
	repository := configfile.New(t.TempDir())
	writeFile(t, repository.Path(), "infra: [\n")

	if _, err := repository.Load(); !errors.Is(err, config.ErrMalformed) {
		t.Errorf("Load() error = %v, want %v", err, config.ErrMalformed)
	}
}

// TestRepositoryLoadUnknownField verifies that unknown keys are not silently
// ignored: a misspelled setting would otherwise disappear without a trace.
func TestRepositoryLoadUnknownField(t *testing.T) {
	repository := configfile.New(t.TempDir())
	writeFile(t, repository.Path(), `infra:
  network:
    name: localcloud
    subnet: 172.28.0.0/16
  unknown: true
`)

	if _, err := repository.Load(); !errors.Is(err, config.ErrMalformed) {
		t.Errorf("Load() error = %v, want %v", err, config.ErrMalformed)
	}
}

// TestConfiguratorUsesFileRepository wires the business logic and the adapter
// together: the configurator stores settings with the file repository, reads
// them back and validates them.
func TestConfiguratorUsesFileRepository(t *testing.T) {
	repository := configfile.New(t.TempDir())
	configurator := config.NewConfigurator(repository)

	if err := configurator.Load(); !errors.Is(err, config.ErrNotConfigured) {
		t.Errorf("Load() of an unconfigured workspace error = %v, want %v", err, config.ErrNotConfigured)
	}

	if err := configurator.Save(settings()); err != nil {
		t.Fatalf("Save() error = %v, want nil", err)
	}

	reloaded := config.NewConfigurator(repository)
	if err := reloaded.Load(); err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}

	got, err := reloaded.Settings()
	if err != nil {
		t.Fatalf("Settings() error = %v, want nil", err)
	}
	if !reflect.DeepEqual(got, settings()) {
		t.Errorf("Settings() = %+v, want %+v", got, settings())
	}
}

// writeFile writes content to path, failing the test on error.
func writeFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
