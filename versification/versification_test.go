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
			t.Errorf("Expected PSA 51:0 to have 2 mapped verses, but only found %d", len(v.Verses))
		}
	} else {
		t.Errorf("PSA 51:0 mapping not found")
	}
}

/*
{"mappedVerses":{"GEN 31:55":["GEN 32:1"],"GEN 32:1-32":["GEN 32:2-33"],"LEV 6:1-7":["LEV 5:20-26"],"LEV 6:8-30":["LEV 6:1-23"],"PSA 51:0":["PSA 51:1","PSA 51:2"],"PSA 51:1-19":["PSA 51:3-21"],"ACT 19:40":["ACT 19:40"],"ACT 19:41":["ACT 19:40"],"S3Y 1:1-29":["DAG 3:24-52"],"S3Y 1:30-31":["DAG 3:52-53"],"S3Y 1:32":["DAG 3:55"],"S3Y 1:33":["DAG 3:54"]}}
*/
