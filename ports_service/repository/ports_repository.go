package repository

import (
	"portsservice/ports"
)

// StorageI is an interface one should implement
// to provide ports storage
type StorageI interface {
	Upsert(port *ports.Port) error
	AllPorts() ([]*ports.Port, error)
	FindByID(string) (*ports.Port, error)
	Init() error
	Close() error
}
