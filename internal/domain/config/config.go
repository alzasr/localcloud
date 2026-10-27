// Package config is the configurator of a local cloud workspace: it
// deserializes the settings of the workspace, keeps them for the other
// business logic modules and checks that they follow the rules of the domain.
//
// The settings are grouped into blocks: the infrastructure with the stateful
// services and the network they share, and the stateless services that run on
// top of it.
//
// The package knows nothing about the format the settings are stored in:
// reading and writing them is the job of a Repository implementation from the
// adapters layer. The package never prints anything.
package config

// Kind is the kind of a stateful infrastructure service.
type Kind string

// Kinds of the stateful infrastructure services known at the moment.
const (
	// KindDatabase is a database server, such as PostgreSQL.
	KindDatabase Kind = "database"
	// KindMessageBus is a message broker, such as RabbitMQ.
	KindMessageBus Kind = "message_bus"
)

// Config holds the settings of a local cloud workspace.
type Config struct {
	// Infra is the block with the stateful services and their network.
	Infra Infra
	// Services is the block with the stateless services running on top of
	// the infrastructure.
	Services []Service
}

// Infra is the block with the stateful services, such as databases and
// message buses, and the network they share.
type Infra struct {
	Network  Network
	Services []InfraService
}

// Network describes the network shared by the infrastructure services.
type Network struct {
	// Name is the name of the network.
	Name string
	// Subnet is the address range of the network in CIDR notation.
	Subnet string
}

// InfraService is a stateful service: it keeps its data between restarts.
type InfraService struct {
	Name string
	Kind Kind
	// Version is the version of the service to run; empty means the default
	// one.
	Version string
	Port    int
}

// Service is a stateless service that runs on top of the infrastructure.
type Service struct {
	Name string
	Port int
	// DependsOn lists the names of the infra services the service needs.
	DependsOn []string
}
