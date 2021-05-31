package versification

import (
	"encoding/json"
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
	//const vrsString = fse.readFileSync(path.resolve(__dirname, '../test_data/truncated_versification.vrs')).toString();
	//const vrsJson = vrs2json(vrsString);
	//console.log("brad here is where I want to look")
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
		t.Error("Unable to open json test data file")
	}
	j, err := json.Marshal(p)
	t.Errorf("%s", string(j))

	//const preSuccinct = preSuccinctVerseMapping(vrsJson.mappedVerses);
	//console.log("brad done here is where I want to look")
	//let preSuccinctBooks = ['GEN', 'LEV', 'PSA', 'ACT', 'S3Y'];
	//t.equal(Object.keys(preSuccinct).length, preSuccinctBooks.length);
	//f//or (const book of preSuccinctBooks) {
	//	t.ok(book in preSuccinct);
	//}
	//t.ok('31' in preSuccinct['GEN']);
	//t.ok('32' in preSuccinct['GEN']);
	//t.ok(preSuccinct['S3Y']['1'][0][2][0].includes('DAG'))

	//const reversed = reverseVersification(vrsJson);
	//const preSuccinctReversed = preSuccinctVerseMapping(reversed.reverseMappedVerses);
	//preSuccinctBooks = ['GEN', 'LEV', 'PSA', 'ACT', 'DAG'];
	//t.equal(Object.keys(preSuccinctReversed).length, preSuccinctBooks.length);
	//for (const book of preSuccinctBooks) {
	//    t.ok(book in preSuccinctReversed);
	// }
	// t.ok('5' in preSuccinctReversed['LEV']);
	// t.ok('6' in preSuccinctReversed['LEV']);
	// t.ok(preSuccinctReversed['DAG']['3'][0][2][0].includes('S3Y'))
	// console.log(JSON.stringify(preSuccinctReversed, null, 2));

}
