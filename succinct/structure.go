package succinct

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Block map[string]string

type Sequence struct {
	Type   string
	Blocks []Block
}

type Doc struct {
	Headers   map[string]string
	MainId    string
	Sequences map[string]Sequence
}

type DocSet struct {
	Id    string
	Enums map[string]string
	Docs  map[string]Doc
}

func DocSetFromJSON(pathString string) (DocSet, error) {
	jsonFile, err := os.Open(pathString)
	if err != nil {
		return DocSet{}, err
	}
	defer jsonFile.Close()
	bytes, _ := ioutil.ReadAll(jsonFile)
	var suc DocSet
	err = json.Unmarshal(bytes, &suc)
	if err != nil {
		return DocSet{}, err
	}
	return suc, nil
}
