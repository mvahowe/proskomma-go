package versification

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestVrsToForwardMappings(t *testing.T) {
	jsonFile, err := os.Open("../test_data/truncated_versification.vrs")
	if err != nil {
		t.Error("Unable to open json test data file")
	}
	defer jsonFile.Close()
	bytes, _ := ioutil.ReadAll(jsonFile)
	s := string(bytes)

	m, err := VrsToForwardMappings(s)

	if err != nil {
		t.Errorf("Error running VrsToForwardMappings: %s", err)
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

func TestReverseVersification(t *testing.T) {
	jsonFile, err := os.Open("../test_data/truncated_versification.vrs")
	if err != nil {
		t.Error("Unable to open json test data file")
	}
	defer jsonFile.Close()
	bytes, _ := ioutil.ReadAll(jsonFile)
	s := string(bytes)

	m, err := VrsToForwardMappings(s)

	bb, err := json.Marshal(m)
	log.Printf("%s", bb)

	r, err := ReverseVersification(m)

	if err != nil {
		t.Errorf("Error running VrsToForwardMappings: %s", err)
	}

	if len(r.MappedVerses) == 0 {
		t.Errorf("No reverse mappings were returned")
	}

	b, err := json.Marshal(r)
	log.Printf("%s", b)

	//TODO remove
	if len(r.MappedVerses) != 0 {
		t.Errorf("No reverse mappings were returned")
	}
}
