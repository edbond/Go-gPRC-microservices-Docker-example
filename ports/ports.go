package ports

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// Port is a domain model of port
type Port struct {
	Key         string    `json:"key,omitempty"`
	Name        string    `json:"name,omitempty"`
	City        string    `json:"city,omitempty"`
	Country     string    `json:"country,omitempty"`
	Alias       []string  `json:"alias,omitempty"`
	Regions     []string  `json:"regions,omitempty"`
	Coordinates []float64 `json:"coordinates,omitempty"`
	Province    string    `json:"province,omitempty"`
	Timezone    string    `json:"timezone,omitempty"`
	Unlocks     []string  `json:"unlocks,omitempty"`
	Code        string    `json:"code,omitempty"`
}

var (
	// ErrPortParse error means that a port record in json file is invalid
	ErrPortParse = errors.New("error parsing port json")
)

// PortCallback is a function that will be called when port parsed
// from JSON
type PortCallback func(*Port)

// LoadFromJSON reads JSON file from a reader and parses ports
// Calls onPort callback function for each port parsed
func LoadFromJSON(log *logrus.Entry, filename string, onPort PortCallback) error {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening ports json file: %w", err)
	}

	// Close file when function exits
	defer jsonFile.Close()

	reader := bufio.NewReader(jsonFile)

	dec := json.NewDecoder(reader)

	// read open bracket
	_, err = dec.Token()
	if err != nil {
		return fmt.Errorf("error reading open bracket in json: %w", err)
	}

	// while the object contains values
	for dec.More() {
		var port Port

		// decode json object key, this will be PortCode
		portCode, err := dec.Token()
		if err != nil {
			return fmt.Errorf("error reading port object key: %w", err)
		}

		// decode an object value (Port)
		err = dec.Decode(&port)
		if err != nil {
			return fmt.Errorf("%s: %w", ErrPortParse, err)
		}

		port.Key = portCode.(string)
		onPort(&port)
	}

	// read closing bracket
	_, err = dec.Token()
	if err != nil {
		log.Errorf("no closing bracker in json: %v", err)
	}

	err = jsonFile.Close()
	if err != nil {
		log.Warnf("error closing json file: %v", err)
	}

	return nil
}

// PortToProto converts Port struct to PortProto
func PortToProto(port *Port) *PortProto {
	return &PortProto{
		Key:         port.Key,
		Name:        port.Name,
		City:        port.City,
		Country:     port.Country,
		Alias:       port.Alias,
		Regions:     port.Regions,
		Coordinates: port.Coordinates,
		Province:    port.Province,
		Timezone:    port.Timezone,
		Unlocks:     port.Unlocks,
		Code:        port.Code,
	}
}

// ProtoToPort converts Port from protobuf to struct
func ProtoToPort(port *PortProto) *Port {
	return &Port{
		Key:         port.Key,
		Name:        port.Name,
		City:        port.City,
		Country:     port.Country,
		Alias:       port.Alias,
		Regions:     port.Regions,
		Coordinates: port.Coordinates,
		Province:    port.Province,
		Timezone:    port.Timezone,
		Unlocks:     port.Unlocks,
		Code:        port.Code,
	}
}
