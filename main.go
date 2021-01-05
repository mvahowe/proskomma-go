package main

import (
	"fmt"

	"github.com/mvahowe/proskomma-go/succinct"
)

func main() {
	ds, _ := succinct.DocSetFromJSON(
		"./test_data/serialize_example.json",
	)

	for docId := range ds.Docs {
		seq := ds.Docs[docId].Sequences[ds.Docs[docId].MainId]
		for _, block := range seq.Blocks {
			fmt.Printf("%+v\n", block)
		}
	}
}
