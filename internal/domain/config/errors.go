package config

import "errors"

// Errors returned by the package. They describe what is wrong with the
// settings; user interface wording is chosen by the adapters.
var (
	// ErrNotConfigured reports that the workspace has no settings yet.
	ErrNotConfigured = errors.New("workspace is not configured")
	// ErrMalformed reports that the stored settings cannot be deserialized.
	ErrMalformed = errors.New("settings are malformed")
	// ErrNetworkNameRequired reports that the infrastructure network has no
	// name.
	ErrNetworkNameRequired = errors.New("infra network name must not be empty")
	// ErrInvalidSubnet reports that the infrastructure network subnet is not
	// an address range in CIDR notation.
	ErrInvalidSubnet = errors.New("infra network subnet must be a valid CIDR")
	// ErrServiceNameRequired reports that a service has no name.
	ErrServiceNameRequired = errors.New("service name must not be empty")
	// ErrDuplicateService reports that a service name is used more than once.
	ErrDuplicateService = errors.New("service name is used more than once")
	// ErrUnknownKind reports that an infrastructure service has a kind that is
	// not known.
	ErrUnknownKind = errors.New("unknown infra service kind")
	// ErrInvalidPort reports that a service port is out of the allowed range.
	ErrInvalidPort = errors.New("port must be between 1 and 65535")
	// ErrUnknownDependency reports that a service depends on an infrastructure
	// service that is not described in the settings.
	ErrUnknownDependency = errors.New("depends on an unknown infra service")
)
