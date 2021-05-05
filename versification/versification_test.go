package versification

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestVrsToMappings(t *testing.T) {
	jsonFile, err := os.Open("../test_data/truncated_versification.vrs")
	if err != nil {
		t.Error("Unable to open json test data file")
	}
	defer jsonFile.Close()
	bytes, _ := ioutil.ReadAll(jsonFile)
	s := string(bytes)

	m, err := VrsToMappings(s)

	if err != nil {
		t.Errorf("Error running VrsToMappings: %s", err)
	}

	if len(m.MappedVerses) == 0 {
		t.Errorf("No vrs mappings were returned")
	}

	if v, present := m.MappedVerses["PSA 51:0"]; present {
		if len(v.Verses) != 2 {
			t.Errorf("Expected PSA 51:0 to have 2 mapped verses, but found %d", len(v.Verses))
		}
	} else {
		t.Errorf("PSA 51:0 mapping not found")
	}
}
