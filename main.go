package main

import (
	"fmt"

	"github.com/mvahowe/proskomma-go/succinct"
)

func main() {
	ds, err := succinct.DocSetFromJSON(
		"./test_data/serialize_example.json",
	)
	if err != nil {
		fmt.Printf("error getting DocSet from JSON: %s\n", err)
	}

	fmt.Printf("Enums: %+v\n", ds.Enums)

	for docId := range ds.Docs {
		seq := ds.Docs[docId].Sequences[ds.Docs[docId].MainId]
		for _, block := range seq.Blocks {
			fmt.Printf("%+v\n", block)
		}
	}
}
