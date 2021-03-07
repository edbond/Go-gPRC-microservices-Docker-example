package storage

// Implementation of Storage interface by in-memory map

import (
	"errors"
	"portsservice/ports"
)

var (
	// ErrNotFound returned when Port not found in repository
	ErrNotFound = errors.New("port not found in repository")
)

// MemoryRepository stores all ports in memory map
type MemoryRepository struct {
	ports map[string]ports.Port
}

// Init initialize map in memory
func (repository *MemoryRepository) Init() error {
	repository.ports = make(map[string]ports.Port)
	return nil
}

// Close clears map
func (repository *MemoryRepository) Close() error {
	repository.ports = make(map[string]ports.Port)
	return nil
}

// Upsert adds or replace Port in repository by Port Key
func (repository *MemoryRepository) Upsert(port *ports.Port) error {
	repository.ports[port.Key] = *port
	return nil
}

// AllPorts returns all Ports in repository
func (repository *MemoryRepository) AllPorts() ([]ports.Port, error) {
	allPorts := make([]ports.Port, 0, len(repository.ports))
	for k := range repository.ports {
		p := repository.ports[k]
		allPorts = append(allPorts, p)
	}
	return allPorts, nil
}

// FindByID lookups Port by Key in repository
func (repository *MemoryRepository) FindByID(key string) (*ports.Port, error) {
	port, found := repository.ports[key]
	if found {
		return &port, nil
	}

	return nil, ErrNotFound
}
