package repository

import "ports.services.com/ports"

// PortsRepository is an interface one should implement
// to provide ports storage
type PortsRepository interface {
	Upsert(ports.Port) error
	AllPorts() ([]ports.Port, error)
	FindByID(string) (*ports.Port, error)
	Init() error
	Close() error
}
