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

//TODO Bcv probably isn't a good name...because this structure is used for cv/bcv...only bcv has the Book populated
//js code refers in vars to this as mappings, but we use mappings all over the place
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
	bookCodes := []string{
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
	//const ret = new ByteArray(64);

	for _, vm := range v {
		recordTypeStr := vm.MappingType
		fromVerseStart := vm.Verses[0]
		fromVerseEnd := vm.Verses[1]
		//verses := vm.Verses
		//for (const [recordTypeStr, [fromVerseStart, fromVerseEnd], mappings] of preSuccinctBC) {
		//   const pos = ret.length;
		pos := ba.Length()

		//   const recordType = recordTypeStr === 'bcv' ? bcvMappingType : cvMappingType;
		recordType := bcvMappingType
		if recordTypeStr == "cv" {
			recordType = cvMappingType
		}

		//   ret.pushNBytes([0, fromVerseStart, fromVerseEnd]);
		ba.PushNBytes([]uint32{0, uint32(fromVerseStart), uint32(fromVerseEnd)})

		//   if (recordType === bcvMappingType) {
		//       const bookIndex = bci[mappings[0][3]];
		//       ret.pushNByte(bookIndex);
		//   }
		if recordType == bcvMappingType {
			bookIndex := bci[vm.Bcv.Book]
			ba.PushNByte(uint32(bookIndex))
		}

		//   ret.pushNByte(mappings.length);
		if recordTypeStr == "cv" {
			ba.PushNByte(3)
		} else {
			ba.PushNByte(4)

		}

		//   for (const [ch, fromV] of mappings) {
		//       ret.pushNBytes([ch, fromV]);
		//   }
		ba.PushNBytes([]uint32{uint32(vm.Bcv.Chapter), uint32(vm.Bcv.FromVerse)})

		//   const recordLength = ret.length - pos;
		recordLength := ba.Length() - pos

		//   if (recordLength > 63) {
		//       throw new Error(`Mapping in succinctifyVerseMapping ${JSON.stringify(mappings)} is too long (${recordLength} bytes)`);
		//   }
		if recordLength > 63 {
			jsonMappings, _ := json.Marshal(vm.Bcv)

			err := fmt.Errorf("Mapping in succinctifyVerseMapping %s is too long (%d bytes)", jsonMappings, recordLength)
			return ba, err
		}
		//   ret.setByte(pos, makeMappingLengthByte(recordType, recordLength));
		err := ba.SetByte(pos, uint8(makeMappingLengthByte(recordType, recordLength)))
		if err != nil {
			return ba, err
		}
	}

	//ret.trim();
	err := ba.Trim()
	if err != nil {
		return ba, err
	}
	//return ret;

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
		//ret[book] = {};
		s.Mappings[book] = make(map[string]succinct.ByteArray)
		for chapter, verseMappings := range chapterMap {
			//ret[book][chapter] = succinctifyVerseMapping(mappings, bci);
			s.Mappings[book][chapter], err = succinctifyVerseMapping(verseMappings, bookCodeToIndex, indexToBookCode)
			if err != nil {
				return s, err
			}
		}
	}

	//p.BookMappings
	/*
	   const bci = bookCodeIndex();
	   for (const [book, chapters] of Object.entries(preSuccinctVerseMapping(preSuccinct))) {
	       ret[book] = {};
	       for (const [chapter, mappings] of Object.entries(chapters)) {
	           ret[book][chapter] = succinctifyVerseMapping(mappings, bci);
	       }
	   }
	*/
	return s, nil
}

/*
    "S3Y 1:1-29": [
      "DAG 3:24-52"
    ],
-----------------------
   toSpecs DAG 3:24-52
   fromSpec S3Y 1:1-29
   fromBook S3Y
   fromCVV 1:1-29
   toBook DAG
   record bcv
   fromCh 1
   fromV 1-29
   fromV includes a dash
   fromV 1
   toV 29
   record bcv,1,29,
   ----
      toCh 3
      fromV 24-52
      fromV includes a dash
      fromV 24
      toV 52
      not cv
      record bcv,1,29,3,24,52,DAG
   ----
   setting ret[fromBook] to empty object
   setting ret[fromBook][fromCh] to empty object
-----------------------
*/

func preSuccinctVerseMapping(m map[string][]string) (PreSuccinctMappings, error) {
	p := NewPreSuccinctMappings()
	//for (let [fromSpec, toSpecs] of Object.entries(mappingJson)) {
	for k, mv := range m {
		//k is the fromSpec like GEN 31:55
		//mv is basically the toSpecs but is always an array

		//const [fromBook, fromCVV] = fromSpec.split(' ');
		s := strings.Split(k, " ")
		fromBook := s[0]
		fromCvv := s[1]

		//const toBook = toSpecs[0].split(' ')[0];
		toBook := strings.Split(mv[0], " ")[0]

		//const record = toBook === fromBook ? ["cv"] : ["bcv"];
		record := NewVerseMappings()
		record.MappingType = "bcv"
		if toBook == fromBook {
			record.MappingType = "cv"
		}

		//let [fromCh, fromV] = fromCVV.split(':');
		s = strings.Split(fromCvv, ":")
		fromCh := s[0]
		fromV := s[1]

		//let toV = fromV;
		toV := fromV

		//if (fromV.includes('-')) {
		//    const vBits = fromV.split('-');
		//    fromV = vBits[0];
		//    toV = vBits[1];
		//}
		if strings.Contains(fromV, "-") {
			s = strings.Split(fromV, "-")
			fromV = s[0]
			toV = s[1]
		}

		//record.push([parseInt(fromV), parseInt(toV)]);
		fromVInt, err := strconv.Atoi(fromV)
		if err != nil {
			return p, err
		}
		toVInt, err := strconv.Atoi(toV)
		if err != nil {
			return p, err
		}
		record.Verses = append(record.Verses, fromVInt, toVInt)
		//record.push([]);   //??? not sure if I need to do anything here at first
		//need to add the fromV and toV as ints to the record.Verses

		//for (const toCVV of toSpecs.map(ts => ts.split(' ')[1])) {
		for i := range mv {
			s = strings.Split(mv[i], " ")
			toCvv := s[1]
			//let [toCh, fromV] = toCVV.split(':');
			s = strings.Split(toCvv, ":")
			toCh := s[0]
			fromV := s[1]
			//let toV = fromV;
			toV = fromV
			//if (fromV.includes('-')) {
			//    const vBits = fromV.split('-');
			//    fromV = vBits[0];
			//    toV = vBits[1];
			//}

			if strings.Contains(fromV, "-") {
				s = strings.Split(fromV, "-")
				fromV = s[0]
				toV = s[1]
			}
			//if (record[0] === 'cv') {
			//    record[2].push([parseInt(toCh), parseInt(fromV), parseInt(toV)]);
			//} else {
			//    record[2].push([parseInt(toCh), parseInt(fromV), parseInt(toV), toBook]);
			//}
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
		//if (!(fromBook in ret)) {
		//    console.log("   setting ret[fromBook] to empty object")
		//    ret[fromBook] = {};
		//}
		if _, present := p.BookMappings[fromBook]; !present {
			//p.BookMappings[fromBook] = make([]ChapterMappings, 0)
			p.BookMappings[fromBook] = make(map[string][]VerseMappings)
		}

		//if (!(fromCh in ret[fromBook])) {
		//    console.log("   setting ret[fromBook][fromCh] to empty object")
		//    ret[fromBook][fromCh] = [];
		//}
		if _, present := p.BookMappings[fromBook][fromCh]; !present {
			p.BookMappings[fromBook][fromCh] = make([]VerseMappings, 0)
		}
		//ret[fromBook][fromCh].push(record);
		p.BookMappings[fromBook][fromCh] = append(p.BookMappings[fromBook][fromCh], record)

	}

	//b, err := json.Marshal(user)
	//if err != nil {
	//    fmt.Printf("Error: %s", err)
	//    return;
	// }
	//fmt.Println(string(b))

	return p, nil
}

//TODO unsuccinctifyVerseMapping
//TODO mapVerse
