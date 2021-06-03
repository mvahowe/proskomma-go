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

func TestPreSuccinctVerseMapping(t *testing.T) {
	jsonFile, err := os.Open("../test_data/truncated_versification.vrs")
	if err != nil {
		t.Error("Unable to open json test data file")
	}
	defer jsonFile.Close()
	bytes, _ := ioutil.ReadAll(jsonFile)
	s := string(bytes)

	m := VrsToForwardMappings(s)

	p, err := preSuccinctVerseMapping(m.MappedVerses)
	if err != nil {
		t.Errorf("preSuccinctVerseMapping failed %s", err)
	}

	preSuccinctBooks := []string{"GEN", "LEV", "PSA", "ACT", "S3Y"}
	if len(preSuccinctBooks) != len(p.BookMappings) {
		t.Errorf("Expected preSuccinct mappings to have %d books, but found %d", len(preSuccinctBooks), len(p.BookMappings))
	}

	for _, b := range preSuccinctBooks {
		if _, present := p.BookMappings[b]; !present {
			t.Errorf("Expected book %s to be a key in preSuccinct mappings but not found.", b)
		}
	}

	if _, present := p.BookMappings["GEN"]["31"]; !present {
		t.Error("Expected book/chapter mapping GEN 31 to be present, but it was not.")
	}
	if _, present := p.BookMappings["GEN"]["32"]; !present {
		t.Error("Expected book/chapter mapping GEN 32 to be present, but it was not.")
	}

	if vm, present := p.BookMappings["S3Y"]["1"]; present {
		if vm[0].Bcv.Book != "DAG" {
			t.Errorf("Expected to find mapping to book DAG, but found %s", vm[0].Bcv.Book)
		}
	} else {
		t.Error("Expected book/chapter mapping S3Y 1 to be present, but it was not.")
	}

	r := ReverseVersification(m)

	if len(r.MappedVerses) == 0 {
		t.Errorf("No reverse mappings were returned")
	}
	pr, err := preSuccinctVerseMapping(r.MappedVerses)
	if err != nil {
		t.Errorf("preSuccinctVerseMapping failed on reverse mappings %s", err)
	}

	preSuccinctBooks = []string{"GEN", "LEV", "PSA", "ACT", "DAG"}
	if len(preSuccinctBooks) != len(pr.BookMappings) {
		t.Errorf("Expected preSuccinct reverse mappings to have %d books, but found %d", len(preSuccinctBooks), len(pr.BookMappings))
	}

	for _, b := range preSuccinctBooks {
		if _, present := pr.BookMappings[b]; !present {
			t.Errorf("Expected book %s to be a key in preSuccinct reverse mappings but not found.", b)
		}
	}

	if _, present := pr.BookMappings["LEV"]["5"]; !present {
		t.Error("Expected book/chapter mapping LEV 5 to be present, but it was not.")
	}
	if _, present := pr.BookMappings["LEV"]["6"]; !present {
		t.Error("Expected book/chapter mapping LEV 6 to be present, but it was not.")
	}

	if vm, present := pr.BookMappings["DAG"]["3"]; present {
		if vm[0].Bcv.Book != "S3Y" {
			t.Errorf("Expected to find mapping to book S3Y, but found %s", vm[0].Bcv.Book)
		}
	} else {
		t.Error("Expected book/chapter mapping DAG 3 to be present, but it was not.")
	}

	//j, err := json.Marshal(pr)
	//t.Errorf("%s", string(j))

}
