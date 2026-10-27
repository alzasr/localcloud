package config

import (
	"fmt"
	"net"
)

// Bounds of a valid TCP port.
const (
	minPort = 1
	maxPort = 65535
)

// Validate checks the settings against the rules of the domain: the
// infrastructure network must be described, every service must have a name, a
// known kind where it applies and a valid port, names must be unique, and a
// stateless service may only depend on described infrastructure services.
func (c Config) Validate() error {
	if c.Infra.Network.Name == "" {
		return ErrNetworkNameRequired
	}
	if _, _, err := net.ParseCIDR(c.Infra.Network.Subnet); err != nil {
		return fmt.Errorf("%w: %q", ErrInvalidSubnet, c.Infra.Network.Subnet)
	}

	infraNames := make(map[string]struct{}, len(c.Infra.Services))
	for _, service := range c.Infra.Services {
		if err := service.validate(); err != nil {
			return fmt.Errorf("infra service %q: %w", service.Name, err)
		}
		if _, ok := infraNames[service.Name]; ok {
			return fmt.Errorf("infra service %q: %w", service.Name, ErrDuplicateService)
		}
		infraNames[service.Name] = struct{}{}
	}

	names := make(map[string]struct{}, len(c.Services))
	for _, service := range c.Services {
		if err := service.validate(); err != nil {
			return fmt.Errorf("service %q: %w", service.Name, err)
		}
		if _, ok := names[service.Name]; ok {
			return fmt.Errorf("service %q: %w", service.Name, ErrDuplicateService)
		}
		names[service.Name] = struct{}{}

		for _, dependency := range service.DependsOn {
			if _, ok := infraNames[dependency]; !ok {
				return fmt.Errorf(
					"service %q depends on %q: %w",
					service.Name,
					dependency,
					ErrUnknownDependency,
				)
			}
		}
	}

	return nil
}

// validate checks a stateful infrastructure service.
func (s InfraService) validate() error {
	if s.Name == "" {
		return ErrServiceNameRequired
	}
	if s.Kind != KindDatabase && s.Kind != KindMessageBus {
		return fmt.Errorf("%w: %q", ErrUnknownKind, s.Kind)
	}
	if !validPort(s.Port) {
		return fmt.Errorf("%w: %d", ErrInvalidPort, s.Port)
	}

	return nil
}

// validate checks a stateless service.
func (s Service) validate() error {
	if s.Name == "" {
		return ErrServiceNameRequired
	}
	if !validPort(s.Port) {
		return fmt.Errorf("%w: %d", ErrInvalidPort, s.Port)
	}

	return nil
}

// validPort reports whether port is in the allowed range.
func validPort(port int) bool {
	return port >= minPort && port <= maxPort
}
