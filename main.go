package main

import (
	"./succinct"
	"fmt"
)

func main() {
	ds, _ := succinct.DocSetFromJSON(
		"./test_data/serialize_example.json",
	)
	for docId := range ds.Docs {
		seq := ds.Docs[docId].Sequences[ds.Docs[docId].MainId]
		for blockMap := range seq.BlockArrayMaps {
			fmt.Println("here", blockMap)
		}
	}
}
