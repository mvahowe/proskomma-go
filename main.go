package main

import (
	"fmt"

	"./succinct"
)

func main() {
	ds, _ := succinct.DocSetFromJSON(
		"./test_data/serialize_example.json",
	)

	for docId := range ds.Docs {
		seq := ds.Docs[docId].Sequences[ds.Docs[docId].MainId]
		for i, blockMap := range seq.BlockArrayMaps {
			fmt.Println("here", i)
			for k, v := range blockMap {
				fmt.Println(k, v)
			}
		}
	}
}
