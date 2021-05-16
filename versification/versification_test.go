package versification

import (
	"io/ioutil"
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

	m := VrsToForwardMappings(s)

	if len(m.MappedVerses) == 0 {
		t.Errorf("No vrs mappings were returned")
	}

	if v, present := m.MappedVerses["PSA 51:0"]; present {
		if len(v) != 2 {
			t.Errorf("Expected PSA 51:0 to have 2 mapped verses, but found %d", len(v))
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

	m := VrsToForwardMappings(s)
	r := ReverseVersification(m)

	if len(r.MappedVerses) == 0 {
		t.Errorf("No reverse mappings were returned")
	}

	for _, mv := range m.MappedVerses {
		if _, present := r.MappedVerses[mv[0]]; !present {
			t.Errorf("Expected mapped verse %s to be a key in reverse mappings, but not found.", mv[0])
		}
	}
}
