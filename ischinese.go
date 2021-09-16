package ischinese

import (
	"bufio"
	"embed"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"log"
	"strings"
)

var debugFlag = false

func debug(v ...interface{}) {
	if debugFlag {
		log.Println(v...)
	}
}

var commonRange = [][]rune{
	// https://en.wikipedia.org/wiki/CJK_Unified_Ideographs
	{
		'\u4E00', '\u9FFC', // CJK Unified Ideographs
	},
	{
		'\u3400', '\u4DBF', // CJK Unified Ideographs Extension A
	},
	{
		'\U00020000', '\U0002A6DD', // CJK Unified Ideographs Extension B
	},
	{
		'\U0002A700', '\U0002B734', // CJK Unified Ideographs Extension C
	},
	{
		'\U0002B740', '\U0002B81D', // CJK Unified Ideographs Extension D
	},
	{
		'\U0002B820', '\U0002CEA1', // CJK Unified Ideographs Extension E
	},
	{
		'\U0002CEB0', '\U0002EBE0', // CJK Unified Ideographs Extension F
	},
	{
		'\U00030000', '\U0003134F', // CJK Unified Ideographs Extension G
	},
	{
		'\uFA0E', '\uFA0F', // CJK Compatibility Ideographs
	},
	{
		'\uFA11', '\uFA11', // CJK Compatibility Ideographs
	},
	{
		'\uFA13', '\uFA14', // CJK Compatibility Ideographs
	},
	{
		'\uFA1F', '\uFA1F', // CJK Compatibility Ideographs
	},
	{
		'\uFA21', '\uFA21', // CJK Compatibility Ideographs
	},
	{
		'\uFA23', '\uFA24', // CJK Compatibility Ideographs
	},
	{
		'\uFA27', '\uFA29', // CJK Compatibility Ideographs
	},
	{
		'\u3300', '\u33FF', // Other CJK ideographs in Unicode, not Unified
	},
	{
		'\uFE30', '\uFE4F', // Other CJK ideographs in Unicode, not Unified
	},
	{
		'\uF900', '\uFAFF', // Other CJK ideographs in Unicode, not Unified
	},
	{
		'\U0002F800', '\U0002FA1F', // Other CJK ideographs in Unicode, not Unified
	},
	// https://en.wikipedia.org/wiki/CJK_Symbols_and_Punctuation
	{
		'\u3000', '\u303F',
	},
	// https://en.wikipedia.org/wiki/Chinese_punctuation
	{
		'\uFF0C',
		'\uFF0C',
	},
	{
		'\uFF01',
		'\uFF01',
	},
	{
		'\uFF1F',
		'\uFF1F',
	},
	{
		'\uFF1A',
		'\uFF1B',
	},
	{
		'\uFF08',
		'\uFF09',
	},
	{
		'\uFF3B',
		'\uFF3B',
	},
	{
		'\uFF3D',
		'\uFF3D',
	},
	{
		'\u3010',
		'\u3011',
	},
}

var simplifiedDict map[rune]struct{}
var traditionalDict map[rune]struct{}

func init() {
	simplifiedDict = make(map[rune]struct{})
	traditionalDict = make(map[rune]struct{})
	err := buildDictionary(simplifiedDict, traditionalDict)
	if err != nil {
		panic(err)
	}
}

//go:embed Unihan_Variants.txt
var fs embed.FS

func buildDictionary(simplifiedDict, traditionalDict map[rune]struct{}) error {
	file, err := fs.Open("Unihan_Variants.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	parseUnicode := func(s string, dict map[rune]struct{}) {
		r, err := parseUnicodeString(s)
		if err != nil {
			// eat err
			return
		}
		dict[r] = struct{}{}
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// skip comments
		if strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		switch fields[1] {
		case "kSimplifiedVariant":
			parseUnicode(fields[0], traditionalDict)
			for _, field := range fields[2:] {
				parseUnicode(field, simplifiedDict)
			}
		case "kTraditionalVariant":
			parseUnicode(fields[0], simplifiedDict)
			for _, field := range fields[2:] {
				parseUnicode(field, traditionalDict)
			}
		default:
			continue
		}
	}
	return scanner.Err()
}

const unicodeStringPrefix = "U+"

func parseUnicodeString(s string) (rune, error) {
	s = strings.TrimPrefix(s, unicodeStringPrefix)
	// unicode use 4 bytes
	times := 8 - len(s)
	if times < 0 {
		return 0, errors.New("invalid unicode")
	}
	// prefix padding
	for i := 0; i < times; i++ {
		s = "0" + s
	}
	bs, err := hex.DecodeString(s)
	if err != nil {
		return 0, err
	}
	return rune(binary.BigEndian.Uint32(bs)), nil
}

func isChineseChar(r rune) bool {
	for _, runes := range commonRange {
		if runes[0] <= r && r <= runes[1] {
			return true
		}
	}
	return false
}

func isSimplifiedChineseChar(r rune) bool {
	if !isChineseChar(r) {
		return false
	}
	if _, ok := simplifiedDict[r]; ok {
		return true
	}
	if _, ok := traditionalDict[r]; ok {
		return false
	}
	return true
}

func isTraditionalChineseChar(r rune) bool {
	if !isChineseChar(r) {
		return false
	}
	if _, ok := traditionalDict[r]; ok {
		return true
	}
	if _, ok := simplifiedDict[r]; ok {
		return false
	}
	return true
}

// IsChinese true if more than 50% of unicode code points are Chinese unicode
func IsChinese(s string) bool {
	return nonPureFuncHelper(s, isChineseChar)
}

// IsSimplifiedChinese true if more than 50% of unicode code points are simplified Chinese unicode
func IsSimplifiedChinese(s string) bool {
	return nonPureFuncHelper(s, isSimplifiedChineseChar)
}

// IsTraditionalChinese true if more than 50% of unicode code points are traditional Chinese unicode
func IsTraditionalChinese(s string) bool {
	return nonPureFuncHelper(s, isTraditionalChineseChar)
}
func nonPureFuncHelper(s string, f func(rune) bool) bool {
	if len(s) == 0 {
		return true
	}
	var counter float32
	var total float32
	for _, r := range s {
		total++
		if !f(r) {
			debug(string([]rune{r}))
		} else {
			counter++
		}
	}
	return counter/total > 0.5
}

// IsPureChinese true if 100% of unicode code points are Chinese unicode
func IsPureChinese(s string) bool {
	return pureFuncHelper(s, isChineseChar)
}

// IsPureSimplifiedChinese true if 100% of unicode code points are simplified Chinese unicode
func IsPureSimplifiedChinese(s string) bool {
	return pureFuncHelper(s, isSimplifiedChineseChar)
}

// IsPureTraditionalChinese true if 100% of unicode code points are traditional Chinese unicode
func IsPureTraditionalChinese(s string) bool {
	return pureFuncHelper(s, isTraditionalChineseChar)
}

func pureFuncHelper(s string, f func(rune) bool) bool {
	for _, r := range s {
		if !f(r) {
			debug(string([]rune{r}))
			return false
		}
	}
	return true
}
