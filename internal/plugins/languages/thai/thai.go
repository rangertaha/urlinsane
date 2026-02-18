package thai

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
)

const LANGUAGE string = "th"

type Thai struct {
	code         string
	name         string
	description  string
	numerals     map[string][]string
	graphemes    []string
	vowels       []string
	misspellings [][]string
	homophones   [][]string
	antonyms     map[string][]string
	homoglyphs   map[string][]string
}

func (l *Thai) Id() string                        { return l.code }
func (l *Thai) Name() string                      { return l.name }
func (l *Thai) Description() string               { return l.description }
func (l *Thai) Numerals() map[string][]string     { return l.numerals }
func (l *Thai) Cardinal() map[string]string       { return languages.NumeralMap(l.numerals, 0) }
func (l *Thai) Ordinal() map[string]string        { return languages.NumeralMap(l.numerals, 1) }
func (l *Thai) Graphemes() []string               { return l.graphemes }
func (l *Thai) Vowels() []string                  { return l.vowels }
func (l *Thai) Misspellings() [][]string          { return l.misspellings }
func (l *Thai) Homophones() [][]string            { return l.homophones }
func (l *Thai) Antonyms() map[string][]string     { return l.antonyms }
func (l *Thai) Homoglyphs() map[string][]string   { return l.homoglyphs }
func (l *Thai) SimilarChars(char string) []string { return languages.SimilarChars(l.homoglyphs, char) }
func (l *Thai) SimilarSpellings(word string) []string {
	return languages.SimilarSpellings(l.misspellings, word)
}
func (l *Thai) SimilarSounds(word string) []string {
	return languages.SimilarSounds(l.homophones, word)
}
func (l *Thai) Keyboards() (boards []internal.Keyboard) {
	for _, b := range languages.Keyboards() {
		if b.Language() == l.code {
			boards = append(boards, b)
		}
	}
	return
}

var (
	thHomophones = [][]string{
		{"จุด", "."},
	}
	thAntonyms = map[string][]string{
		"ดี":        {"เลว"},
		"เลว":       {"ดี"},
		"ใหญ่":      {"เล็ก"},
		"เล็ก":      {"ใหญ่"},
		"สูง":       {"ต่ำ"},
		"ต่ำ":       {"สูง"},
		"ยาว":       {"สั้น"},
		"สั้น":      {"ยาว"},
		"มาก":       {"น้อย"},
		"น้อย":      {"มาก"},
		"ใหม่":      {"เก่า"},
		"เก่า":      {"ใหม่"},
		"เร็ว":      {"ช้า"},
		"ช้า":       {"เร็ว"},
		"ง่าย":      {"ยาก"},
		"ยาก":       {"ง่าย"},
		"ร้อน":      {"เย็น"},
		"เย็น":      {"ร้อน"},
		"สว่าง":     {"มืด"},
		"มืด":       {"สว่าง"},
		"แข็งแรง":   {"อ่อนแอ"},
		"อ่อนแอ":    {"แข็งแรง"},
		"หนัก":      {"เบา"},
		"เบา":       {"หนัก"},
		"หนา":       {"บาง"},
		"บาง":       {"หนา"},
		"เปิด":      {"ปิด"},
		"ปิด":       {"เปิด"},
		"เข้า":      {"ออก"},
		"ออก":       {"เข้า"},
		"ขึ้น":      {"ลง"},
		"ลง":        {"ขึ้น"},
		"หน้า":      {"หลัง"},
		"หลัง":      {"หน้า"},
		"ซ้าย":      {"ขวา"},
		"ขวา":       {"ซ้าย"},
		"ใกล้":      {"ไกล"},
		"ไกล":       {"ใกล้"},
		"ถูก":       {"ผิด"},
		"ผิด":       {"ถูก"},
		"จริง":      {"เท็จ"},
		"เท็จ":      {"จริง"},
		"ชนะ":       {"แพ้"},
		"แพ้":       {"ชนะ"},
		"ซื้อ":      {"ขาย"},
		"ขาย":       {"ซื้อ"},
		"มา":        {"ไป"},
		"ไป":        {"มา"},
		"มี":        {"ไม่มี"},
		"ไม่มี":     {"มี"},
		"ใช่":       {"ไม่"},
		"ไม่":       {"ใช่"},
	}

	Language = Thai{
		code:        LANGUAGE,
		name:        "Thai",
		description: "Thai is a Kra–Dai language spoken primarily in Thailand",
		numerals: map[string][]string{
			"0":          {"ศูนย์"},
			"1":          {"หนึ่ง", "ที่หนึ่ง"},
			"2":          {"สอง", "ที่สอง"},
			"3":          {"สาม", "ที่สาม"},
			"4":          {"สี่", "ที่สี่"},
			"5":          {"ห้า", "ที่ห้า"},
			"6":          {"หก", "ที่หก"},
			"7":          {"เจ็ด", "ที่เจ็ด"},
			"8":          {"แปด", "ที่แปด"},
			"9":          {"เก้า", "ที่เก้า"},
			"10":         {"สิบ", "ที่สิบ"},
			"11":         {"สิบเอ็ด", "ที่สิบเอ็ด"},
			"12":         {"สิบสอง", "ที่สิบสอง"},
			"13":         {"สิบสาม", "ที่สิบสาม"},
			"14":         {"สิบสี่", "ที่สิบสี่"},
			"15":         {"สิบห้า", "ที่สิบห้า"},
			"16":         {"สิบหก", "ที่สิบหก"},
			"17":         {"สิบเจ็ด", "ที่สิบเจ็ด"},
			"18":         {"สิบแปด", "ที่สิบแปด"},
			"19":         {"สิบเก้า", "ที่สิบเก้า"},
			"20":         {"ยี่สิบ", "ที่ยี่สิบ"},
			"30":         {"สามสิบ"},
			"40":         {"สี่สิบ"},
			"50":         {"ห้าสิบ"},
			"60":         {"หกสิบ"},
			"70":         {"เจ็ดสิบ"},
			"80":         {"แปดสิบ"},
			"90":         {"เก้าสิบ"},
			"100":        {"หนึ่งร้อย"},
			"1000":       {"หนึ่งพัน"},
			"10000":      {"หนึ่งหมื่น"},
			"1000000":    {"หนึ่งล้าน"},
			"1000000000": {"หนึ่งพันล้าน"},
		},
		graphemes:    []string{},
		vowels:       []string{},
		misspellings: [][]string{},
		homophones:   thHomophones,
		antonyms:     thAntonyms,
		homoglyphs:   languages.DefaultLatinHomoglyphs(),
	}
)

func init() { languages.AddLanguage(LANGUAGE, func() internal.Language { return &Language }) }
