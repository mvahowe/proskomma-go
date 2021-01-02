package main

import (
	"./succinct"
	"fmt"
)

func main() {
	fmt.Println(
		succinct.DocSetFromJSON(
			"./test_data/serialize_example.json",
		),
	)
}
