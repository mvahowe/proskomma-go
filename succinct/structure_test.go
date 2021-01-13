package succinct

import (
	"fmt"
	"testing"
)

func TestLoadSuccinctJSON(t *testing.T) {
	ds, err := DocSetFromJSON(
		"../test_data/serialize_example.json",
	)
	if err != nil {
		fmt.Printf("error getting DocSet from JSON: %s\n", err)
	}
	if ds == nil {
		fmt.Print("Returned docSet is nil")
	}
}

