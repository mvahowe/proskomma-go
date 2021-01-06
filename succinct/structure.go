package succinct

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Enums struct {
	IDs         ByteArray `json:"ids"`
	WordLike    ByteArray `json:"wordLike"`
	NotWordLike ByteArray `json:"notWordLike"`
	ScopeBits   ByteArray `json:"scopeBits"`
	GraftTypes  ByteArray `json:"graftTypes"`
}

type Block struct {
	BlockScope     ByteArray `json:"bs"`
	BlockGrafts    ByteArray `json:"bg"`
	BlockItems     ByteArray `json:"c"`
	OpenScopes     ByteArray `json:"os"`
	IncludedScopes ByteArray `json:"is"`
}

type Sequence struct {
	Type   string  `json:"type"`
	Blocks []Block `json:"blocks"`
}

type Doc struct {
	Headers   map[string]string   `json:"headers"`
	MainId    string              `json:"mainId"`
	Sequences map[string]Sequence `json:"sequences"`
}

type DocSet struct {
	Id    string         `json:"id"`
	Enums Enums          `json:"enums"`
	Docs  map[string]Doc `json:"docs"`
}

func DocSetFromJSON(pathString string) (*DocSet, error) {
	jsonFile, err := os.Open(pathString)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	bytes, _ := ioutil.ReadAll(jsonFile)
	var suc DocSet
	err = json.Unmarshal(bytes, &suc)
	return &suc, err
}
