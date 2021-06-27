package versification

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/mvahowe/proskomma-go/succinct"
)

type ForwardMappings struct {
	MappedVerses map[string][]string `json:"mappedVerses"`
}

type ReverseMappings struct {
	MappedVerses map[string][]string `json:"reverseMappedVerses"`
}

type PreSuccinctMappings struct {
	BookMappings map[string]map[string][]VerseMappings
}

type VerseMappings struct {
	MappingType string
	Verses      []int
	Bcv         Bcv
}

type Bcv struct {
	Chapter   int
	FromVerse int
	ToVerse   int
	Book      string
}

type SuccinctMappings struct {
	Mappings map[string]map[string]succinct.ByteArray
}

const cvMappingType = 2
const bcvMappingType = 3

func bookCodeIndex() (map[string]int, map[int]string) {
	// From Paratext via Scripture Burrito
	bookCodes := [...]string{
		"GEN",
		"EXO",
		"LEV",
		"NUM",
		"DEU",
		"JOS",
		"JDG",
		"RUT",
		"1SA",
		"2SA",
		"1KI",
		"2KI",
		"1CH",
		"2CH",
		"EZR",
		"NEH",
		"EST",
		"JOB",
		"PSA",
		"PRO",
		"ECC",
		"SNG",
		"ISA",
		"JER",
		"LAM",
		"EZK",
		"DAN",
		"HOS",
		"JOL",
		"AMO",
		"OBA",
		"JON",
		"MIC",
		"NAM",
		"HAB",
		"ZEP",
		"HAG",
		"ZEC",
		"MAL",
		"MAT",
		"MRK",
		"LUK",
		"JHN",
		"ACT",
		"ROM",
		"1CO",
		"2CO",
		"GAL",
		"EPH",
		"PHP",
		"COL",
		"1TH",
		"2TH",
		"1TI",
		"2TI",
		"TIT",
		"PHM",
		"HEB",
		"JAS",
		"1PE",
		"2PE",
		"1JN",
		"2JN",
		"3JN",
		"JUD",
		"REV",
		"TOB",
		"JDT",
		"ESG",
		"WIS",
		"SIR",
		"BAR",
		"LJE",
		"S3Y",
		"SUS",
		"BEL",
		"1MA",
		"2MA",
		"3MA",
		"4MA",
		"1ES",
		"2ES",
		"MAN",
		"PS2",
		"ODA",
		"PSS",
		"JSA",
		"JDB",
		"TBS",
		"SST",
		"DNT",
		"BLT",
		"EZA",
		"5EZ",
		"6EZ",
		"DAG",
		"PS3",
		"2BA",
		"LBA",
		"JUB",
		"ENO",
		"1MQ",
		"2MQ",
		"3MQ",
		"REP",
		"4BA",
		"LAO"}

	bookCodeToIndex := make(map[string]int)
	indexToBookCode := make(map[int]string)
	for i := range bookCodes {
		bookCodeToIndex[bookCodes[i]] = i
		indexToBookCode[i] = bookCodes[i]
	}
	return bookCodeToIndex, indexToBookCode
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

func NewPreSuccinctMappings() PreSuccinctMappings {
	var p PreSuccinctMappings
	p.BookMappings = make(map[string]map[string][]VerseMappings)
	return p
}

func NewVerseMappings() VerseMappings {
	var v VerseMappings
	v.Verses = make([]int, 0)
	return v
}

func NewSuccinctMappings() SuccinctMappings {
	var s SuccinctMappings
	s.Mappings = make(map[string]map[string]succinct.ByteArray)
	return s
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

func makeMappingLengthByte(r int, l int) int {
	return l + (r * 64)
}

func succinctifyVerseMapping(v []VerseMappings, bci map[string]int, ibc map[int]string) (succinct.ByteArray, error) {
	ba := succinct.NewByteArray(64)

	for _, vm := range v {
		recordTypeStr := vm.MappingType
		fromVerseStart := vm.Verses[0]
		fromVerseEnd := vm.Verses[1]
		pos := ba.Length()

		recordType := bcvMappingType
		if recordTypeStr == "cv" {
			recordType = cvMappingType
		}

		ba.PushNBytes([]uint32{0, uint32(fromVerseStart), uint32(fromVerseEnd)})

		if recordType == bcvMappingType {
			bookIndex := bci[vm.Bcv.Book]
			ba.PushNByte(uint32(bookIndex))
		}

		if recordTypeStr == "cv" {
			ba.PushNByte(3)
		} else {
			ba.PushNByte(4)

		}

		ba.PushNBytes([]uint32{uint32(vm.Bcv.Chapter), uint32(vm.Bcv.FromVerse)})

		recordLength := ba.Length() - pos
		if recordLength > 63 {
			jsonMappings, _ := json.Marshal(vm.Bcv)

			err := fmt.Errorf("Mapping in succinctifyVerseMapping %s is too long (%d bytes)", jsonMappings, recordLength)
			return ba, err
		}
		err := ba.SetByte(pos, uint8(makeMappingLengthByte(recordType, recordLength)))
		if err != nil {
			return ba, err
		}
	}

	err := ba.Trim()
	if err != nil {
		return ba, err
	}

	return ba, nil
}

func SuccinctifyVerseMappings(m map[string][]string) (SuccinctMappings, error) {
	s := NewSuccinctMappings()
	bookCodeToIndex, indexToBookCode := bookCodeIndex()
	p, err := preSuccinctVerseMapping(m)
	if err != nil {
		return s, err
	}
	for book, chapterMap := range p.BookMappings {
		s.Mappings[book] = make(map[string]succinct.ByteArray)
		for chapter, verseMappings := range chapterMap {
			s.Mappings[book][chapter], err = succinctifyVerseMapping(verseMappings, bookCodeToIndex, indexToBookCode)
			if err != nil {
				return s, err
			}
		}
	}

	return s, nil
}

func preSuccinctVerseMapping(m map[string][]string) (PreSuccinctMappings, error) {
	p := NewPreSuccinctMappings()
	for k, mv := range m {
		s := strings.Split(k, " ")
		fromBook := s[0]
		fromCvv := s[1]
		toBook := strings.Split(mv[0], " ")[0]

		record := NewVerseMappings()
		record.MappingType = "bcv"
		if toBook == fromBook {
			record.MappingType = "cv"
		}

		s = strings.Split(fromCvv, ":")
		fromCh := s[0]
		fromV := s[1]
		toV := fromV

		if strings.Contains(fromV, "-") {
			s = strings.Split(fromV, "-")
			fromV = s[0]
			toV = s[1]
		}

		fromVInt, err := strconv.Atoi(fromV)
		if err != nil {
			return p, err
		}
		toVInt, err := strconv.Atoi(toV)
		if err != nil {
			return p, err
		}
		record.Verses = append(record.Verses, fromVInt, toVInt)

		for i := range mv {
			s = strings.Split(mv[i], " ")
			toCvv := s[1]
			s = strings.Split(toCvv, ":")
			toCh := s[0]
			fromV := s[1]
			toV = fromV

			if strings.Contains(fromV, "-") {
				s = strings.Split(fromV, "-")
				fromV = s[0]
				toV = s[1]
			}
			toChInt, err := strconv.Atoi(toCh)
			if err != nil {
				return p, err
			}
			fromVInt, err := strconv.Atoi(fromV)
			if err != nil {
				return p, err
			}
			toVInt, err := strconv.Atoi(toV)
			if err != nil {
				return p, err
			}
			if record.MappingType == "cv" {
				record.Bcv = Bcv{
					Chapter:   toChInt,
					FromVerse: fromVInt,
					ToVerse:   toVInt,
				}
			} else {
				record.Bcv = Bcv{
					Chapter:   toChInt,
					FromVerse: fromVInt,
					ToVerse:   toVInt,
					Book:      toBook,
				}
			}
		}
		if _, present := p.BookMappings[fromBook]; !present {
			p.BookMappings[fromBook] = make(map[string][]VerseMappings)
		}

		if _, present := p.BookMappings[fromBook][fromCh]; !present {
			p.BookMappings[fromBook][fromCh] = make([]VerseMappings, 0)
		}
		p.BookMappings[fromBook][fromCh] = append(p.BookMappings[fromBook][fromCh], record)
	}

	return p, nil
}
