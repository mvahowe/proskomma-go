package succinct

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type BlockStringMap map[string]string

type BlockArrayMap map[string]*ByteArray

type Sequence struct {
	Type            string
	BlockStringMaps []BlockStringMap `json:"blocks"`
	BlockArrayMaps  []BlockArrayMap
}

type Doc struct {
	Headers   map[string]string
	MainId    string
	Sequences map[string]Sequence
}

type EnumArrayMap map[string]*ByteArray

type DocSet struct {
	Id           string
	EnumStrings  map[string]string `json:"enums"`
	EnumArrayMap EnumArrayMap
	Docs         map[string]Doc
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
	if err != nil {
		return nil, err
	}
	suc.EnumArrayMap = make(EnumArrayMap)
	for enumLabel, enumB64 := range suc.EnumStrings {
		ba := NewByteArray(256)
		ba.fromBase64(enumB64)
		_ = ba.Trim()
		suc.EnumArrayMap[enumLabel] = &ba
	}
	for docId := range suc.Docs {
		doc := suc.Docs[docId]
		for seqId := range doc.Sequences {
			seq := doc.Sequences[seqId]
			for blockCount := range seq.BlockStringMaps {
				fmt.Println("BLOCK")
				blockArrayMap := BlockArrayMap{}
				seq.BlockArrayMaps = append(seq.BlockArrayMaps, blockArrayMap)
				blockStringMap := seq.BlockStringMaps[blockCount]
				for blockFieldKey, blockFieldValue := range blockStringMap {
					ba := NewByteArray(256)
					ba.fromBase64(blockFieldValue)
					_ = ba.Trim()
					blockArrayMap[blockFieldKey] = &ba
					fmt.Println(blockFieldKey, blockArrayMap[blockFieldKey], len(blockArrayMap), len(seq.BlockArrayMaps))
				}
			}
		}
	}
	return &suc, nil
}
