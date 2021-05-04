package versification

import (
	"encoding/json"
	"log"
	"regexp"
	"strings"
)

type VrsMappings struct {
	MappedVerses map[string]MappedVerse `json:"mappedVerses"`
}

type MappedVerse struct {
	Verses []string
}

func NewVrsMappings() VrsMappings {
	var m VrsMappings
	m.MappedVerses = make(map[string]MappedVerse)
	return m
}

func VrsToMappings(s string) (VrsMappings, error) {
	mappings := NewVrsMappings()
	lines := strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
	r, _ := regexp.Compile("^([A-Z1-6]{3} [0-9]+:[0-9]+(-[0-9]+)?) = ([A-Z1-6]{3} [0-9]+:[0-9]+[a-z]?(-[0-9]+)?)$")
	for i := range lines {
		lineBits := r.FindStringSubmatch(lines[i])
		if lineBits == nil {
			continue
		}

		verses := make([]string, 0, len(lineBits)-2)
		verses = append(verses, lineBits[3])
		if v, present := mappings.MappedVerses[lineBits[1]]; present {
			verses = append(v.Verses, verses...)
		}
		mappedVerse := MappedVerse{Verses: verses}
		mappings.MappedVerses[lineBits[1]] = mappedVerse
	}

	b, _ := json.Marshal(mappings)
	log.Printf("%s", string(b))

	return mappings, nil
}
