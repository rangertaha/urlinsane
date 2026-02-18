package japanese

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
)

const LANGUAGE string = "ja"

type Japanese struct {
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

func (l *Japanese) Id() string                      { return l.code }
func (l *Japanese) Name() string                    { return l.name }
func (l *Japanese) Description() string             { return l.description }
func (l *Japanese) Numerals() map[string][]string   { return l.numerals }
func (l *Japanese) Cardinal() map[string]string     { return languages.NumeralMap(l.numerals, 0) }
func (l *Japanese) Ordinal() map[string]string      { return languages.NumeralMap(l.numerals, 1) }
func (l *Japanese) Graphemes() []string             { return l.graphemes }
func (l *Japanese) Vowels() []string                { return l.vowels }
func (l *Japanese) Misspellings() [][]string        { return l.misspellings }
func (l *Japanese) Homophones() [][]string          { return l.homophones }
func (l *Japanese) Antonyms() map[string][]string   { return l.antonyms }
func (l *Japanese) Homoglyphs() map[string][]string { return l.homoglyphs }
func (l *Japanese) SimilarChars(char string) []string {
	return languages.SimilarChars(l.homoglyphs, char)
}
func (l *Japanese) SimilarSpellings(word string) []string {
	return languages.SimilarSpellings(l.misspellings, word)
}
func (l *Japanese) SimilarSounds(word string) []string {
	return languages.SimilarSounds(l.homophones, word)
}
func (l *Japanese) Keyboards() (boards []internal.Keyboard) {
	for _, b := range languages.Keyboards() {
		if b.Language() == l.code {
			boards = append(boards, b)
		}
	}
	return
}

var (
	jaHomophones = [][]string{
		{"てん", ".", "点"},
	}
	jaAntonyms = map[string][]string{
		"良い":  {"悪い"},
		"悪い":  {"良い"},
		"大きい": {"小さい"},
		"小さい": {"大きい"},
		"高い":  {"低い"},
		"低い":  {"高い"},
		"長い":  {"短い"},
		"短い":  {"長い"},
		"多い":  {"少ない"},
		"少ない": {"多い"},
		"新しい": {"古い"},
		"古い":  {"新しい"},
		"早い":  {"遅い"},
		"遅い":  {"早い"},
		"熱い":  {"冷たい"},
		"冷たい": {"熱い"},
		"明るい": {"暗い"},
		"暗い":  {"明るい"},
		"強い":  {"弱い"},
		"弱い":  {"強い"},
		"重い":  {"軽い"},
		"軽い":  {"重い"},
		"太い":  {"細い"},
		"細い":  {"太い"},
		"太る":  {"痩せる"},
		"痩せる": {"太る"},
		"開く":  {"閉じる"},
		"閉じる": {"開く"},
		"入る":  {"出る"},
		"出る":  {"入る"},
		"上":   {"下"},
		"下":   {"上"},
		"前":   {"後ろ"},
		"後ろ":  {"前"},
		"左":   {"右"},
		"右":   {"左"},
		"近い":  {"遠い"},
		"遠い":  {"近い"},
		"正しい": {"間違い"},
		"間違い": {"正しい"},
		"真":   {"偽"},
		"偽":   {"真"},
		"勝つ":  {"負ける"},
		"負ける": {"勝つ"},
		"買う":  {"売る"},
		"売る":  {"買う"},
		"来る":  {"行く"},
		"行く":  {"来る"},
		"ある":  {"ない"},
		"ない":  {"ある"},
	}

	Language = Japanese{
		code:        LANGUAGE,
		name:        "Japanese",
		description: "Japanese is a Japonic language written with kana and kanji",
		numerals: map[string][]string{
			"0":          {"零", "第零"},
			"1":          {"一", "第一"},
			"2":          {"二", "第二"},
			"3":          {"三", "第三"},
			"4":          {"四", "第四"},
			"5":          {"五", "第五"},
			"6":          {"六", "第六"},
			"7":          {"七", "第七"},
			"8":          {"八", "第八"},
			"9":          {"九", "第九"},
			"10":         {"十", "第十"},
			"11":         {"十一", "第十一"},
			"12":         {"十二", "第十二"},
			"13":         {"十三", "第十三"},
			"14":         {"十四", "第十四"},
			"15":         {"十五", "第十五"},
			"16":         {"十六", "第十六"},
			"17":         {"十七", "第十七"},
			"18":         {"十八", "第十八"},
			"19":         {"十九", "第十九"},
			"20":         {"二十", "第二十"},
			"30":         {"三十"},
			"40":         {"四十"},
			"50":         {"五十"},
			"60":         {"六十"},
			"70":         {"七十"},
			"80":         {"八十"},
			"90":         {"九十"},
			"100":        {"百", "第百"},
			"1000":       {"千", "第千"},
			"10000":      {"万", "第万"},
			"1000000":    {"百万", "第百万"},
			"1000000000": {"十億", "第十億"},
		},
		graphemes:    []string{"あ", "い", "う", "え", "お"},
		vowels:       []string{"あ", "い", "う", "え", "お"},
		misspellings: [][]string{},
		homophones:   jaHomophones,
		antonyms:     jaAntonyms,
		homoglyphs:   languages.DefaultLatinHomoglyphs(),
	}
)

func init() { languages.AddLanguage(LANGUAGE, func() internal.Language { return &Language }) }
