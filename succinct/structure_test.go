package succinct

import (
	"testing"
)

func loadSuccinctJSON(t *testing.T, path string) *DocSet {
	ds, err := DocSetFromJSON(path)
	if err != nil {
		t.Errorf("error getting DocSet from JSON: %s\n", err)
	}
	if ds == nil {
		t.Errorf("Returned docSet is nil")
	}
	return ds
}

func TestLoadSuccinctJSON(t *testing.T) {
	loadSuccinctJSON(t, "../test_data/serialize_example.json")
}

func TestHeaderBytesFromJSON(t *testing.T) {
	ds := loadSuccinctJSON(t, "../test_data/serialize_example.json")
	for docId := range ds.Docs {
		seq := ds.Docs[docId].Sequences[ds.Docs[docId].MainId]
		for _, block := range seq.Blocks {
			checkHeaderBytes(t, &block.BlockItems)
			checkHeaderBytes(t, &block.BlockGrafts)
			checkHeaderBytes(t, &block.BlockScope)
			checkHeaderBytes(t, &block.IncludedScopes)
			checkHeaderBytes(t, &block.OpenScopes)
		}
	}
}

func TestTokenChars(t *testing.T) {
	ds := loadSuccinctJSON(t, "../test_data/serialize_example.json")
	for docId := range ds.Docs {
		seq := ds.Docs[docId].Sequences[ds.Docs[docId].MainId]
		for _, block := range seq.Blocks {
			pos := 0
			ba := block.BlockItems
			for pos < len(ba.bytes) {
				itemLength, itemType, itemSubtype, err := ba.headerBytes(pos)
				if err != nil {
					t.Errorf("headerBytes threw error: %s", err)
				}
				if itemType > len(itemStrings) {
					t.Errorf("Unexpected itemType %d", itemType)
				}
				if itemType == 0 {
					tokenString, err := tokenChars(&ds.Enums, &ba, itemSubtype, pos)
					if err != nil {
						t.Errorf("Error from tokenChars: %s", err)
					}
					if len(tokenString) == 0 {
						t.Errorf("Empty string from tokenChars")
					}
				}
				pos += itemLength
			}
		}
	}
}