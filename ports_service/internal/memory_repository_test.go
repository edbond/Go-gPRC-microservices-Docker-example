package internal

import (
	"testing"

	"ports.services.com/ports"
)

// Test upsert adds Port to repository
// updates data by key
func TestUpsert(t *testing.T) {
	ajman := ports.Port{
		Key:     "AEAJM",
		Name:    "Ajman",
		City:    "Ajman",
		Country: "United Arab Emirates",
		Alias:   []string{"Ajman alias"},
		Regions: []string{"Ajman region"},
		Coordinates: []float64{
			55.5136433,
			25.4052165,
		},
		Province: "Ajman",
		Timezone: "Asia/Dubai",
		Unlocks:  []string{"AEAJM"},
		Code:     "52000",
	}

	repo := MemoryRepository{}
	err := repo.Init()
	if err != nil {
		t.Fatalf("repository initialize failure: %s", err.Error())
	}

	err = repo.Upsert(ajman)
	if err != nil {
		t.Fatalf("upsert error: %s", err.Error())
	}

	// Verify we have 1 port in repository now
	allPorts, err := repo.AllPorts()
	if len(allPorts) != 1 {
		t.Fatalf("repository should contains 1 port, has %d instead", len(allPorts))
	}

	// Update city, upsert should update port data by the key
	ajman.City = "Kyiv"

	err = repo.Upsert(ajman)
	if err != nil {
		t.Fatalf("upsert error: %s", err.Error())
	}

	// Should have still 1 port in repository
	allPorts, err = repo.AllPorts()
	if len(allPorts) != 1 {
		t.Fatalf("repository should contains 1 port after upsert, has %d instead", len(allPorts))
	}

	if allPorts[0].City != "Kyiv" {
		t.Fatalf("after updating city name, port should have city name Kyiv, has %s instead", allPorts[0].City)
	}

	// After changing the key we should have 2 ports in repository
	ajman.Key = "ABCD"
	err = repo.Upsert(ajman)
	if err != nil {
		t.Fatalf("upsert error: %s", err.Error())
	}

	// Should have 2 port in repository now
	allPorts, err = repo.AllPorts()
	if len(allPorts) != 2 {
		t.Fatalf("repository should contains 2 port after upsert, has %d instead", len(allPorts))
	}
}
