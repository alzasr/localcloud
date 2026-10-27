package config_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/alzasr/localcloud/internal/domain/config"
)

// validSettings returns settings that pass validation. Every call returns
// fresh slices, so tests may modify the result.
func validSettings() config.Config {
	return config.Config{
		Infra: config.Infra{
			Network: config.Network{Name: "localcloud", Subnet: "172.28.0.0/16"},
			Services: []config.InfraService{
				{Name: "postgres", Kind: config.KindDatabase, Version: "16", Port: 5432},
				{Name: "rabbitmq", Kind: config.KindMessageBus, Port: 5672},
			},
		},
		Services: []config.Service{
			{Name: "api", Port: 8080, DependsOn: []string{"postgres", "rabbitmq"}},
		},
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		modify  func(*config.Config)
		wantErr error
	}{
		{
			name:   "valid settings",
			modify: func(*config.Config) {},
		},
		{
			name:    "network without name",
			modify:  func(c *config.Config) { c.Infra.Network.Name = "" },
			wantErr: config.ErrNetworkNameRequired,
		},
		{
			name:    "network without subnet",
			modify:  func(c *config.Config) { c.Infra.Network.Subnet = "" },
			wantErr: config.ErrInvalidSubnet,
		},
		{
			name:    "subnet is not a CIDR",
			modify:  func(c *config.Config) { c.Infra.Network.Subnet = "172.28.0.1" },
			wantErr: config.ErrInvalidSubnet,
		},
		{
			name:    "infra service without name",
			modify:  func(c *config.Config) { c.Infra.Services[0].Name = "" },
			wantErr: config.ErrServiceNameRequired,
		},
		{
			name:    "unknown infra service kind",
			modify:  func(c *config.Config) { c.Infra.Services[0].Kind = "queue" },
			wantErr: config.ErrUnknownKind,
		},
		{
			name:    "infra service without port",
			modify:  func(c *config.Config) { c.Infra.Services[0].Port = 0 },
			wantErr: config.ErrInvalidPort,
		},
		{
			name:    "duplicate infra service",
			modify:  func(c *config.Config) { c.Infra.Services[1].Name = "postgres" },
			wantErr: config.ErrDuplicateService,
		},
		{
			name:    "service without name",
			modify:  func(c *config.Config) { c.Services[0].Name = "" },
			wantErr: config.ErrServiceNameRequired,
		},
		{
			name:    "service port out of range",
			modify:  func(c *config.Config) { c.Services[0].Port = 70000 },
			wantErr: config.ErrInvalidPort,
		},
		{
			name:    "duplicate service",
			modify:  func(c *config.Config) { c.Services = append(c.Services, c.Services[0]) },
			wantErr: config.ErrDuplicateService,
		},
		{
			name:    "service depends on unknown infra service",
			modify:  func(c *config.Config) { c.Services[0].DependsOn = []string{"redis"} },
			wantErr: config.ErrUnknownDependency,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			settings := validSettings()
			tt.modify(&settings)

			if err := settings.Validate(); !errors.Is(err, tt.wantErr) {
				t.Errorf("Validate() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateNamesTheOffendingService verifies that the error tells which
// services are involved, so the problem can be fixed without reading the code.
func TestValidateNamesTheOffendingService(t *testing.T) {
	settings := validSettings()
	settings.Services[0].DependsOn = []string{"redis"}

	err := settings.Validate()
	if err == nil {
		t.Fatal("Validate() error = nil, want an error")
	}

	for _, want := range []string{"api", "redis"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("Validate() error = %q, want it to mention %q", err, want)
		}
	}
}
