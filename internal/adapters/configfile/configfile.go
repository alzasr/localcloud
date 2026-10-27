// Package configfile stores the settings of a workspace in a YAML file placed
// in the workspace directory. It implements the repository port of the
// configurator, so the business logic stays independent of the file format.
package configfile

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/alzasr/localcloud/internal/domain/config"
)

// FileName is the name of the settings file inside a workspace directory.
const FileName = "localcloud.yaml"

// fileMode holds the permissions the settings file is created with.
const fileMode fs.FileMode = 0o644

// fileIndent is the number of spaces the settings file is indented with.
const fileIndent = 2

// Repository stores the settings of a workspace in <directory>/localcloud.yaml.
type Repository struct {
	path string
}

// New returns the repository of the settings file of the given workspace
// directory.
func New(directory string) Repository {
	return Repository{path: filepath.Join(directory, FileName)}
}

// Path returns the location of the settings file.
func (r Repository) Path() string {
	return r.path
}

// Load reads the settings from the file. A missing or empty file is reported
// as config.ErrNotConfigured, an unreadable one as config.ErrMalformed.
func (r Repository) Load() (config.Config, error) {
	data, err := os.ReadFile(r.path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return config.Config{}, fmt.Errorf("%w: no settings file %s", config.ErrNotConfigured, r.path)
		}

		return config.Config{}, fmt.Errorf("read settings file %s: %w", r.path, err)
	}

	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)

	var document documentYAML
	if err := decoder.Decode(&document); err != nil {
		if errors.Is(err, io.EOF) {
			return config.Config{}, fmt.Errorf("%w: settings file %s is empty", config.ErrNotConfigured, r.path)
		}

		return config.Config{}, fmt.Errorf("%w: settings file %s: %v", config.ErrMalformed, r.path, err)
	}

	return document.toConfig(), nil
}

// Save writes the settings to the file.
func (r Repository) Save(settings config.Config) error {
	var buffer bytes.Buffer

	encoder := yaml.NewEncoder(&buffer)
	encoder.SetIndent(fileIndent)

	if err := encoder.Encode(toDocument(settings)); err != nil {
		return fmt.Errorf("encode settings: %w", err)
	}
	if err := encoder.Close(); err != nil {
		return fmt.Errorf("encode settings: %w", err)
	}

	if err := os.WriteFile(r.path, buffer.Bytes(), fileMode); err != nil {
		return fmt.Errorf("write settings file %s: %w", r.path, err)
	}

	return nil
}

// documentYAML is the on-disk representation of the settings.
type documentYAML struct {
	Infra    infraYAML     `yaml:"infra"`
	Services []serviceYAML `yaml:"services,omitempty"`
}

// infraYAML is the on-disk representation of the infrastructure block.
type infraYAML struct {
	Network  networkYAML        `yaml:"network"`
	Services []infraServiceYAML `yaml:"services,omitempty"`
}

// networkYAML is the on-disk representation of the shared network.
type networkYAML struct {
	Name   string `yaml:"name,omitempty"`
	Subnet string `yaml:"subnet,omitempty"`
}

// infraServiceYAML is the on-disk representation of a stateful service.
type infraServiceYAML struct {
	Name    string `yaml:"name"`
	Kind    string `yaml:"kind"`
	Version string `yaml:"version,omitempty"`
	Port    int    `yaml:"port"`
}

// serviceYAML is the on-disk representation of a stateless service.
type serviceYAML struct {
	Name      string   `yaml:"name"`
	Port      int      `yaml:"port"`
	DependsOn []string `yaml:"depends_on,omitempty"`
}

// toConfig converts the parsed file into settings.
func (d documentYAML) toConfig() config.Config {
	settings := config.Config{
		Infra: config.Infra{
			Network: config.Network{
				Name:   d.Infra.Network.Name,
				Subnet: d.Infra.Network.Subnet,
			},
		},
	}

	for _, service := range d.Infra.Services {
		settings.Infra.Services = append(settings.Infra.Services, config.InfraService{
			Name:    service.Name,
			Kind:    config.Kind(service.Kind),
			Version: service.Version,
			Port:    service.Port,
		})
	}

	for _, service := range d.Services {
		settings.Services = append(settings.Services, config.Service{
			Name:      service.Name,
			Port:      service.Port,
			DependsOn: service.DependsOn,
		})
	}

	return settings
}

// toDocument converts settings into their on-disk representation.
func toDocument(settings config.Config) documentYAML {
	document := documentYAML{
		Infra: infraYAML{
			Network: networkYAML{
				Name:   settings.Infra.Network.Name,
				Subnet: settings.Infra.Network.Subnet,
			},
		},
	}

	for _, service := range settings.Infra.Services {
		document.Infra.Services = append(document.Infra.Services, infraServiceYAML{
			Name:    service.Name,
			Kind:    string(service.Kind),
			Version: service.Version,
			Port:    service.Port,
		})
	}

	for _, service := range settings.Services {
		document.Services = append(document.Services, serviceYAML{
			Name:      service.Name,
			Port:      service.Port,
			DependsOn: service.DependsOn,
		})
	}

	return document
}
