package versification

import (
	"log"
	"regexp"
	"strings"
)

type ForwardMappings struct {
	MappedVerses map[string]MappedVerse `json:"mappedVerses"`
}

type ReverseMappings struct {
	MappedVerses map[string]MappedVerse `json:"reverseMappedVerses"`
}

type MappedVerse struct {
	Verses []string
}

func NewForwardMappings() ForwardMappings {
	var m ForwardMappings
	m.MappedVerses = make(map[string]MappedVerse)
	return m
}

func NewReverseMappings() ReverseMappings {
	var m ReverseMappings
	m.MappedVerses = make(map[string]MappedVerse)
	return m
}

func VrsToForwardMappings(s string) (ForwardMappings, error) {
	mappings := NewForwardMappings()
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

	return mappings, nil
}

func ReverseVersification(m ForwardMappings) (ReverseMappings, error) {
	mappings := NewReverseMappings()
	/*
	   // Assumes each verse is only mapped from once
	   const ret = {};
	   for (const [fromSpec, toSpecs] of Object.entries(vrsJson.mappedVerses)) {
	       for (const toSpec of toSpecs) {
	           toSpec in ret ? ret[toSpec].push(fromSpec) : ret[toSpec] = [fromSpec];
	       }
	   }
	   return {reverseMappedVerses: ret};
	*/
	for k, mv := range m.MappedVerses {
		log.Printf("Key: %s", k)
		for i := range mv.Verses {
			log.Printf("   %s", mv.Verses[i])
			verses := make([]string, 0, 1)
			if v, present := mappings.MappedVerses[mv.Verses[i]]; present {
				verses = append(v.Verses, k)
			} else {
				verses = append(verses, k)
			}
			mappedVerse := MappedVerse{Verses: verses}
			mappings.MappedVerses[mv.Verses[i]] = mappedVerse
		}
	}

	return mappings, nil
}
