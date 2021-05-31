package versification

import (
	"regexp"
	"strconv"
	"strings"
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
