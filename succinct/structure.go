package succinct

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type EnumArrayMap map[string]*ByteArray

func (e *EnumArrayMap) UnmarshalJSON(b []byte) error {
	*e = make(EnumArrayMap)
	var enumStrings map[string]string
	err := json.Unmarshal(b, &enumStrings)
	if err != nil {
		return err
	}
	for enumLabel, enumB64 := range enumStrings {
		ba := NewByteArray(256)
		ba.fromBase64(enumB64)
		_ = ba.Trim()
		(*e)[enumLabel] = &ba
	}
	return nil
}

type BlockArrayMap map[string]*ByteArray

func (bam *BlockArrayMap) UnmarshalJSON(b []byte) error {
	*bam = make(BlockArrayMap)
	var blocksStrings map[string]string
	err := json.Unmarshal(b, &blocksStrings)
	if err != nil {
		return err
	}
	for blockFieldKey, blockFieldValue := range blocksStrings {
		ba := NewByteArray(256)
		ba.fromBase64(blockFieldValue)
		_ = ba.Trim()
		(*bam)[blockFieldKey] = &ba
	}
	return nil
}

type Sequence struct {
	Type           string          `json:"type"`
	BlockArrayMaps []BlockArrayMap `json:"blocks"`
}

type Doc struct {
	Headers   map[string]string   `json:"headers"`
	MainId    string              `json:"mainId"`
	Sequences map[string]Sequence `json:"sequences"`
}

type DocSet struct {
	Id           string         `json:"id"`
	EnumArrayMap EnumArrayMap   `json:"enums"`
	Docs         map[string]Doc `json:"docs"`
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
