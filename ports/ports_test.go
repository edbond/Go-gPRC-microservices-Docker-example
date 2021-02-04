package ports

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sirupsen/logrus"
)

func TestLoadJSON(t *testing.T) {
	testCases := []struct {
		filename      string
		expected      []Port
		expectedError error
	}{
		{
			filename: "valid_ports.json",
			expected: []Port{
				{
					Key:         "AEAJM",
					Name:        "Ajman",
					City:        "Ajman",
					Country:     "United Arab Emirates",
					Regions:     []string{"Ajman region"},
					Alias:       []string{"Ajman alias"},
					Coordinates: []float64{55.5136433, 25.4052165},
					Province:    "Ajman",
					Timezone:    "Asia/Dubai",
					Code:        "52000",
				},
			},
			expectedError: nil,
		},
		{
			filename:      "invalid.json",
			expected:      []Port{},
			expectedError: ErrPortParse,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.filename, func(t *testing.T) {
			var err error

			log := logrus.New().WithFields(logrus.Fields{
				"Test": tC.filename,
			})

			ports := []Port{}
			callback := func(port *Port) {
				ports = append(ports, *port)
			}
			err = LoadFromJSON(log, tC.filename, callback)
			if err != nil {
				if tC.expectedError != nil && strings.Contains(err.Error(), tC.expectedError.Error()) {
					// OK, we expect this error
				} else {
					t.Fatal(err)
				}
			}

			if !cmp.Equal(ports, tC.expected) {
				t.Fatalf(`Ports parsed and expected are different:  
				%v`, cmp.Diff(ports, tC.expected))
			}
		})
	}
}
