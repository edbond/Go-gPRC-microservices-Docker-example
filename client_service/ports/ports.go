package ports

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

var (
	// ErrPortParse error means that a port record in json file is invalid
	ErrPortParse = errors.New("error parsing port json")
)

// PortCallback is a function that will be called when port parsed
// from JSON
type PortCallback func(port *Port)

// LoadFromJSON reads JSON file from a reader and parses ports
// Calls onPort callback function for each port parsed
func LoadFromJSON(logger *zerolog.Logger, filename string, onPort PortCallback) error {
	var err error

	jsonFile, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening ports json file: %w", err)
	}

	// Close file when function exits
	defer func() {
		err = jsonFile.Close()
		if err != nil {
			logger.Err(err).Msg("error closing json file")
		}
	}()

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
		logger.Err(err).Msg("no closing bracket in json")
	}

	return nil
}
