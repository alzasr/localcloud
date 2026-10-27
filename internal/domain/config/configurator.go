package config

// Repository reads and writes the settings of a workspace. Where and how the
// settings are stored is an implementation detail of the adapters layer.
type Repository interface {
	// Load returns the stored settings. It returns ErrNotConfigured when the
	// workspace has no settings yet.
	Load() (Config, error)
	// Save stores the settings.
	Save(Config) error
}

// Configurator is the module the other business logic modules get the settings
// of a workspace from. It deserializes the settings, keeps them and checks
// them against the rules of the domain.
type Configurator struct {
	repository Repository

	settings Config
	loaded   bool
}

// NewConfigurator returns a configurator that keeps the settings in the given
// repository.
func NewConfigurator(repository Repository) *Configurator {
	return &Configurator{repository: repository}
}

// Load reads the settings from the repository, validates them and keeps them
// for the other modules. Settings kept by a previous Load or Save are dropped
// when loading fails.
func (c *Configurator) Load() error {
	c.settings, c.loaded = Config{}, false

	settings, err := c.repository.Load()
	if err != nil {
		return err
	}
	if err := settings.Validate(); err != nil {
		return err
	}

	c.settings, c.loaded = settings, true

	return nil
}

// Settings returns the settings kept by the configurator. It returns
// ErrNotConfigured when no settings were loaded or saved yet. Callers must
// treat the returned value as read-only.
func (c *Configurator) Settings() (Config, error) {
	if !c.loaded {
		return Config{}, ErrNotConfigured
	}

	return c.settings, nil
}

// Save validates the settings, stores them in the repository and keeps them
// for the other modules.
func (c *Configurator) Save(settings Config) error {
	if err := settings.Validate(); err != nil {
		return err
	}
	if err := c.repository.Save(settings); err != nil {
		return err
	}

	c.settings, c.loaded = settings, true

	return nil
}
