package korean

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
)

const LANGUAGE string = "ko"

type Korean struct {
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

func (l *Korean) Id() string                      { return l.code }
func (l *Korean) Name() string                    { return l.name }
func (l *Korean) Description() string             { return l.description }
func (l *Korean) Numerals() map[string][]string   { return l.numerals }
func (l *Korean) Cardinal() map[string]string     { return languages.NumeralMap(l.numerals, 0) }
func (l *Korean) Ordinal() map[string]string      { return languages.NumeralMap(l.numerals, 1) }
func (l *Korean) Graphemes() []string             { return l.graphemes }
func (l *Korean) Vowels() []string                { return l.vowels }
func (l *Korean) Misspellings() [][]string        { return l.misspellings }
func (l *Korean) Homophones() [][]string          { return l.homophones }
func (l *Korean) Antonyms() map[string][]string   { return l.antonyms }
func (l *Korean) Homoglyphs() map[string][]string { return l.homoglyphs }
func (l *Korean) SimilarChars(char string) []string {
	return languages.SimilarChars(l.homoglyphs, char)
}
func (l *Korean) SimilarSpellings(word string) []string {
	return languages.SimilarSpellings(l.misspellings, word)
}
func (l *Korean) SimilarSounds(word string) []string {
	return languages.SimilarSounds(l.homophones, word)
}
func (l *Korean) Keyboards() (boards []internal.Keyboard) {
	for _, b := range languages.Keyboards() {
		if b.Language() == l.code {
			boards = append(boards, b)
		}
	}
	return
}

var (
	koHomophones = [][]string{
		{"점", ".", "마침표"},
		{"골뱅이", "@"},
		{"대시", "-"},
		{"슬래시", "/"},
	}
	koAntonyms = map[string][]string{
		"좋다":   {"나쁘다"},
		"나쁘다":  {"좋다"},
		"크다":   {"작다"},
		"작다":   {"크다"},
		"높다":   {"낮다"},
		"낮다":   {"높다"},
		"길다":   {"짧다"},
		"짧다":   {"길다"},
		"많다":   {"적다"},
		"적다":   {"많다"},
		"새롭다":  {"낡다"},
		"낡다":   {"새롭다"},
		"빠르다":  {"느리다"},
		"느리다":  {"빠르다"},
		"쉽다":   {"어렵다"},
		"어렵다":  {"쉽다"},
		"뜨겁다":  {"차갑다"},
		"차갑다":  {"뜨겁다"},
		"밝다":   {"어둡다"},
		"어둡다":  {"밝다"},
		"강하다":  {"약하다"},
		"약하다":  {"강하다"},
		"무겁다":  {"가볍다"},
		"가볍다":  {"무겁다"},
		"두껍다":  {"얇다"},
		"얇다":   {"두껍다"},
		"굵다":   {"가늘다"},
		"가늘다":  {"굵다"},
		"열다":   {"닫다"},
		"닫다":   {"열다"},
		"들어가다": {"나가다"},
		"나가다":  {"들어가다"},
		"위":    {"아래"},
		"아래":   {"위"},
		"앞":    {"뒤"},
		"뒤":    {"앞"},
		"왼쪽":   {"오른쪽"},
		"오른쪽":  {"왼쪽"},
		"가깝다":  {"멀다"},
		"멀다":   {"가깝다"},
		"맞다":   {"틀리다"},
		"틀리다":  {"맞다"},
		"참":    {"거짓"},
		"거짓":   {"참"},
		"이기다":  {"지다"},
		"지다":   {"이기다"},
		"사다":   {"팔다"},
		"팔다":   {"사다"},
		"오다":   {"가다"},
		"가다":   {"오다"},
		"있다":   {"없다"},
		"없다":   {"있다"},
		"네":    {"아니요"},
		"아니요":  {"네"},
	}

	Language = Korean{
		code:        LANGUAGE,
		name:        "Korean",
		description: "Korean is a Koreanic language written mainly with Hangul",
		numerals: map[string][]string{
			"0":          {"영"},
			"1":          {"일", "첫째"},
			"2":          {"이", "둘째"},
			"3":          {"삼", "셋째"},
			"4":          {"사", "넷째"},
			"5":          {"오", "다섯째"},
			"6":          {"육", "여섯째"},
			"7":          {"칠", "일곱째"},
			"8":          {"팔", "여덟째"},
			"9":          {"구", "아홉째"},
			"10":         {"십", "열째"},
			"11":         {"십일"},
			"12":         {"십이"},
			"13":         {"십삼"},
			"14":         {"십사"},
			"15":         {"십오"},
			"16":         {"십육"},
			"17":         {"십칠"},
			"18":         {"십팔"},
			"19":         {"십구"},
			"20":         {"이십"},
			"30":         {"삼십"},
			"40":         {"사십"},
			"50":         {"오십"},
			"60":         {"육십"},
			"70":         {"칠십"},
			"80":         {"팔십"},
			"90":         {"구십"},
			"100":        {"백"},
			"1000":       {"천"},
			"10000":      {"만"},
			"1000000":    {"백만"},
			"1000000000": {"십억"},
		},
		graphemes:    []string{"ㄱ", "ㄴ", "ㄷ", "ㄹ", "ㅁ", "ㅂ", "ㅅ", "ㅇ", "ㅈ", "ㅊ", "ㅋ", "ㅌ", "ㅍ", "ㅎ"},
		vowels:       []string{"ㅏ", "ㅑ", "ㅓ", "ㅕ", "ㅗ", "ㅛ", "ㅜ", "ㅠ", "ㅡ", "ㅣ"},
		misspellings: [][]string{},
		homophones:   koHomophones,
		antonyms:     koAntonyms,
		homoglyphs:   languages.DefaultLatinHomoglyphs(),
	}
)

func init() { languages.AddLanguage(LANGUAGE, func() internal.Language { return &Language }) }
