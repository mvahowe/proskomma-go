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
	Mappings VrsMappings `json:"mappedVerses"`
}

type ReverseMappings struct {
	Mappings VrsMappings `json:"reverseMappedVerses"`
}

type VrsMappings struct {
	MappedVerses map[string][]string
	Keys         []string
}

type PreSuccinctMappings struct {
	BookMappings map[string]map[string][]VerseMappings
}

type VerseMappings struct {
	MappingType string
	Verses      []int
	Mappings    []Bcv
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

type UnsuccinctRecord struct {
	FromVerseStart int
	FromVerseEnd   int
	BookCode       string
	Mappings       []ChapterVerseStart
}

type ChapterVerseStart struct {
	Ch         int
	VerseStart int
}

const cvMappingType = 2
const bcvMappingType = 3

func bookCodeIndex() (map[string]int, map[int]string) {
	// From Paratext via Scripture Burrito
	bookCodes := [...]string{
		"GEN", "EXO", "LEV", "NUM", "DEU", "JOS", "JDG", "RUT", "1SA", "2SA", "1KI", "2KI",
		"1CH", "2CH", "EZR", "NEH", "EST", "JOB", "PSA", "PRO", "ECC", "SNG", "ISA", "JER",
		"LAM", "EZK", "DAN", "HOS", "JOL", "AMO", "OBA", "JON", "MIC", "NAM", "HAB", "ZEP",
		"HAG", "ZEC", "MAL", "MAT", "MRK", "LUK", "JHN", "ACT", "ROM", "1CO", "2CO", "GAL",
		"EPH", "PHP", "COL", "1TH", "2TH", "1TI", "2TI", "TIT", "PHM", "HEB", "JAS", "1PE",
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
	m.Mappings.MappedVerses = make(map[string][]string)
	m.Mappings.Keys = make([]string, 0)
	return m
}

func NewReverseMappings() ReverseMappings {
	var m ReverseMappings
	m.Mappings.MappedVerses = make(map[string][]string)
	m.Mappings.Keys = make([]string, 0)
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
	v.Mappings = make([]Bcv, 0)
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
		if v, present := m.Mappings.MappedVerses[lineBits[1]]; present {
			verses = append(v, verses...)
		} else {
			m.Mappings.Keys = append(m.Mappings.Keys, lineBits[1])
		}
		m.Mappings.MappedVerses[lineBits[1]] = verses
	}

	return m
}

func ReverseVersification(m ForwardMappings) ReverseMappings {
	r := NewReverseMappings()

	for _, k := range m.Mappings.Keys {
		//for k, mv := range m.MappedVerses {
		mv := m.Mappings.MappedVerses[k]
		for i := range mv {
			verses := make([]string, 0, 1)
			if v, present := r.Mappings.MappedVerses[mv[i]]; present {
				verses = append(v, k)
			} else {
				verses = append(verses, k)
				r.Mappings.Keys = append(r.Mappings.Keys, mv[i])
			}
			r.Mappings.MappedVerses[mv[i]] = verses
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
			bookIndex := bci[vm.Mappings[0].Book]
			ba.PushNByte(uint32(bookIndex))
		}

		ba.PushNByte(uint32(len(vm.Mappings)))

		for i, _ := range vm.Mappings {
			ba.PushNBytes([]uint32{uint32(vm.Mappings[i].Chapter), uint32(vm.Mappings[i].FromVerse)})
		}

		recordLength := ba.Length() - pos
		if recordLength > 63 {
			jsonMappings, _ := json.Marshal(vm.Mappings)

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

func SuccinctifyVerseMappings(m VrsMappings) (SuccinctMappings, error) {
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

func preSuccinctVerseMapping(m VrsMappings) (PreSuccinctMappings, error) {
	p := NewPreSuccinctMappings()
	for _, k := range m.Keys {
		mv := m.MappedVerses[k]
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

		//for (const toCVV of toSpecs.map(ts => ts.split(' ')[1])) {
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
				bcv := Bcv{
					Chapter:   toChInt,
					FromVerse: fromVInt,
					ToVerse:   toVInt,
				}
				record.Mappings = append(record.Mappings, bcv)
			} else {
				cv := Bcv{
					Chapter:   toChInt,
					FromVerse: fromVInt,
					ToVerse:   toVInt,
					Book:      toBook,
				}
				record.Mappings = append(record.Mappings, cv)
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

func mappingLengthByte(s succinct.ByteArray, p int) (uint8, uint8, error) {
	sByte, err := s.Byte(p)
	if err != nil {
		return 0, 0, err
	}
	t := sByte >> 6
	l := sByte % 64
	return t, l, nil
}

func UnsuccinctifyVerseMapping(s succinct.ByteArray, c string) ([]UnsuccinctRecord, error) {
	records := make([]UnsuccinctRecord, 0)
	_, indexToBookCode := bookCodeIndex()
	pos := 0
	for pos < s.Length() {
		recordPos := pos
		u := UnsuccinctRecord{}
		recordType, recordLenth, err := mappingLengthByte(s, pos)
		if err != nil {
			return records, err
		}
		recordPos++
		fromVerseStart, err := s.NByte(recordPos)
		if err != nil {
			return records, err
		}
		u.FromVerseStart = int(fromVerseStart)
		recordPos += s.NByteLength(u.FromVerseStart)
		fromVerseEnd, err := s.NByte(recordPos)
		if err != nil {
			return records, err
		}
		u.FromVerseEnd = int(fromVerseEnd)
		recordPos += s.NByteLength(u.FromVerseEnd)
		u.BookCode = c
		if recordType == bcvMappingType {
			bookIndex, err := s.NByte(recordPos)
			if err != nil {
				return records, err
			}
			u.BookCode = indexToBookCode[int(bookIndex)]
			recordPos += s.NByteLength(int(bookIndex))
		}
		nMappings, err := s.NByte(recordPos)
		if err != nil {
			return records, err
		}
		recordPos += s.NByteLength(int(nMappings))
		mappings := make([]ChapterVerseStart, 0)
		for len(mappings) < int(nMappings) {
			m := ChapterVerseStart{}
			ch, err := s.NByte(recordPos)
			if err != nil {
				return records, err
			}
			m.Ch = int(ch)
			recordPos += s.NByteLength(m.Ch)
			verseStart, err := s.NByte(recordPos)
			if err != nil {
				return records, err
			}
			m.VerseStart = int(verseStart)
			recordPos += s.NByteLength(m.VerseStart)
			mappings = append(mappings, m)
		}
		u.Mappings = mappings
		records = append(records, u)
		pos += int(recordLenth)
	}
	return records, nil
}
