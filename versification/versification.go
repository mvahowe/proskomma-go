package versification

import (
	"regexp"
	"strings"
)

type ForwardMappings struct {
	MappedVerses map[string][]string `json:"mappedVerses"`
}

type ReverseMappings struct {
	MappedVerses map[string][]string `json:"reverseMappedVerses"`
}

func NewForwardMappings() ForwardMappings {
	var m ForwardMappings
	m.MappedVerses = make(map[string][]string)
	return m
}

func NewReverseMappings() ReverseMappings {
	var m ReverseMappings
	m.MappedVerses = make(map[string][]string)
	return m
}

func VrsToForwardMappings(s string) ForwardMappings {
	m := NewForwardMappings()
	lines := strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
	r, _ := regexp.Compile("^([A-Z1-6]{3} [0-9]+:[0-9]+(-[0-9]+)?) = ([A-Z1-6]{3} [0-9]+:[0-9]+[a-z]?(-[0-9]+)?)$")
	for i := range lines {
		lineBits := r.FindStringSubmatch(lines[i])
		if lineBits == nil {
			continue
		}

		verses := make([]string, 0, len(lineBits)-2)
		verses = append(verses, lineBits[3])
		if v, present := m.MappedVerses[lineBits[1]]; present {
			verses = append(v, verses...)
		}
		m.MappedVerses[lineBits[1]] = verses
	}

	return m
}

func ReverseVersification(m ForwardMappings) ReverseMappings {
	r := NewReverseMappings()
	for k, mv := range m.MappedVerses {
		for i := range mv {
			verses := make([]string, 0, 1)
			if v, present := r.MappedVerses[mv[i]]; present {
				verses = append(v, k)
			} else {
				verses = append(verses, k)
			}
			r.MappedVerses[mv[i]] = verses
		}
	}
	return r
}
