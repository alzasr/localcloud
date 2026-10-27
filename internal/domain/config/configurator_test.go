package config_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/alzasr/localcloud/internal/domain/config"
)

// fakeRepository is an in-memory implementation of config.Repository.
type fakeRepository struct {
	settings config.Config
	loadErr  error
	saveErr  error

	loads int
	saved *config.Config
}

func (r *fakeRepository) Load() (config.Config, error) {
	r.loads++

	if r.loadErr != nil {
		return config.Config{}, r.loadErr
	}

	return r.settings, nil
}

func (r *fakeRepository) Save(settings config.Config) error {
	if r.saveErr != nil {
		return r.saveErr
	}

	r.saved = &settings

	return nil
}

func TestConfiguratorSettingsBeforeLoad(t *testing.T) {
	configurator := config.NewConfigurator(&fakeRepository{})

	if _, err := configurator.Settings(); !errors.Is(err, config.ErrNotConfigured) {
		t.Errorf("Settings() error = %v, want %v", err, config.ErrNotConfigured)
	}
}

func TestConfiguratorLoad(t *testing.T) {
	repository := &fakeRepository{settings: validSettings()}
	configurator := config.NewConfigurator(repository)

	if err := configurator.Load(); err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}

	if repository.loads != 1 {
		t.Errorf("repository loads = %d, want 1", repository.loads)
	}
	assertSettings(t, configurator, validSettings())
}

func TestConfiguratorLoadRejectsInvalidSettings(t *testing.T) {
	settings := validSettings()
	settings.Infra.Network.Name = ""

	configurator := config.NewConfigurator(&fakeRepository{settings: settings})
	if err := configurator.Load(); !errors.Is(err, config.ErrNetworkNameRequired) {
		t.Errorf("Load() error = %v, want %v", err, config.ErrNetworkNameRequired)
	}

	if _, err := configurator.Settings(); !errors.Is(err, config.ErrNotConfigured) {
		t.Errorf("Settings() error = %v, want %v", err, config.ErrNotConfigured)
	}
}

func TestConfiguratorLoadDropsPreviousSettingsOnError(t *testing.T) {
	repository := &fakeRepository{settings: validSettings()}
	configurator := config.NewConfigurator(repository)
	if err := configurator.Load(); err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}

	repository.loadErr = config.ErrMalformed
	if err := configurator.Load(); !errors.Is(err, config.ErrMalformed) {
		t.Errorf("Load() error = %v, want %v", err, config.ErrMalformed)
	}

	if _, err := configurator.Settings(); !errors.Is(err, config.ErrNotConfigured) {
		t.Errorf("Settings() error = %v, want %v", err, config.ErrNotConfigured)
	}
}

func TestConfiguratorSave(t *testing.T) {
	repository := &fakeRepository{}
	configurator := config.NewConfigurator(repository)

	settings := validSettings()
	if err := configurator.Save(settings); err != nil {
		t.Fatalf("Save() error = %v, want nil", err)
	}

	if repository.saved == nil {
		t.Fatal("Save() did not store the settings in the repository")
	}
	if !reflect.DeepEqual(*repository.saved, settings) {
		t.Errorf("stored settings = %+v, want %+v", *repository.saved, settings)
	}
	assertSettings(t, configurator, settings)
}

func TestConfiguratorSaveRejectsInvalidSettings(t *testing.T) {
	repository := &fakeRepository{}
	configurator := config.NewConfigurator(repository)

	settings := validSettings()
	settings.Services[0].Port = 0

	if err := configurator.Save(settings); !errors.Is(err, config.ErrInvalidPort) {
		t.Errorf("Save() error = %v, want %v", err, config.ErrInvalidPort)
	}
	if repository.saved != nil {
		t.Errorf("Save() stored invalid settings: %+v", *repository.saved)
	}
	if _, err := configurator.Settings(); !errors.Is(err, config.ErrNotConfigured) {
		t.Errorf("Settings() error = %v, want %v", err, config.ErrNotConfigured)
	}
}

func TestConfiguratorSavePropagatesRepositoryError(t *testing.T) {
	repository := &fakeRepository{saveErr: errors.New("disk is full")}
	configurator := config.NewConfigurator(repository)

	if err := configurator.Save(validSettings()); err == nil {
		t.Error("Save() error = nil, want an error")
	}
	if _, err := configurator.Settings(); !errors.Is(err, config.ErrNotConfigured) {
		t.Errorf("Settings() error = %v, want %v", err, config.ErrNotConfigured)
	}
}

// assertSettings verifies the settings kept by the configurator.
func assertSettings(t *testing.T, configurator *config.Configurator, want config.Config) {
	t.Helper()

	got, err := configurator.Settings()
	if err != nil {
		t.Fatalf("Settings() error = %v, want nil", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Settings() = %+v, want %+v", got, want)
	}
}
