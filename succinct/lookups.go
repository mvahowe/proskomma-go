package succinct

import "sort"

func reverseIntLookup(m map[string]int) []string {
	var reversed = make(map[int]string)
	for s, i := range m {
		reversed[i] = s
	}
	var keys []int
	for i := range reversed {
		keys = append(keys, i)
	}
	sort.Ints(keys)
	var ret []string
	for i := range keys {
		ret = append(ret, reversed[i])
	}
	return ret
}

var itemString2Int = map[string]int{
	"token":      0,
	"graft":      1,
	"startScope": 2,
	"endScope":   3,
}

var itemStrings = reverseIntLookup(itemString2Int)

var graftLocation = map[string]string{
	"heading":    "block",
	"title":      "block",
	"endTitle":   "block",
	"remark":     "block",
	"footnote":   "inline",
	"xref":       "inline",
	"noteCaller": "inline",
	"esbCat":     "inline",
}

var scopeString2Int = map[string]int{
	"blockTag":     0,
	"inline":       1,
	"chapter":      2,
	"pubChapter":   3,
	"altChapter":   4,
	"verses":       5,
	"verse":        6,
	"pubVerse":     7,
	"altVerse":     8,
	"esbCat":       9,
	"span":         10,
	"table":        11,
	"cell":         12,
	"milestone":    13,
	"spanWithAtts": 14,
	"attribute":    15,
	"hangingGraft": 16,
	"orphanTokens": 17,
}

var scopeStrings = reverseIntLookup(scopeString2Int)

var tokenString2Int = map[string]int{
	"wordLike":      0,
	"punctuation":   1,
	"lineSpace":     2,
	"eol":           3,
	"softLineBreak": 4,
	"noBreakSpace":  5,
	"bareSlash":     6,
	"unknown":       7,
}

var tokenStrings = reverseIntLookup(tokenString2Int)

var tokenCategory = []string{
	"wordLike",
	"notWordLike",
	"notWordLike",
	"notWordLike",
	"notWordLike",
	"notWordLike",
	"notWordLike",
	"notWordLike",
}
