package polish

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
)

const LANGUAGE string = "pl"

type Polish struct {
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

func (l *Polish) Id() string                      { return l.code }
func (l *Polish) Name() string                    { return l.name }
func (l *Polish) Description() string             { return l.description }
func (l *Polish) Numerals() map[string][]string   { return l.numerals }
func (l *Polish) Cardinal() map[string]string     { return languages.NumeralMap(l.numerals, 0) }
func (l *Polish) Ordinal() map[string]string      { return languages.NumeralMap(l.numerals, 1) }
func (l *Polish) Graphemes() []string             { return l.graphemes }
func (l *Polish) Vowels() []string                { return l.vowels }
func (l *Polish) Misspellings() [][]string        { return l.misspellings }
func (l *Polish) Homophones() [][]string          { return l.homophones }
func (l *Polish) Antonyms() map[string][]string   { return l.antonyms }
func (l *Polish) Homoglyphs() map[string][]string { return l.homoglyphs }
func (l *Polish) SimilarChars(char string) []string {
	return languages.SimilarChars(l.homoglyphs, char)
}
func (l *Polish) SimilarSpellings(word string) []string {
	return languages.SimilarSpellings(l.misspellings, word)
}
func (l *Polish) SimilarSounds(word string) []string {
	return languages.SimilarSounds(l.homophones, word)
}
func (l *Polish) Keyboards() (boards []internal.Keyboard) {
	for _, b := range languages.Keyboards() {
		if b.Language() == l.code {
			boards = append(boards, b)
		}
	}
	return
}

var (
	plMisspellings = [][]string{
		{"lodz", "łódź"},
		{"wroclaw", "wrocław"},
	}
	plHomophones = [][]string{
		{"kropka", "."},
	}
	plAntonyms = map[string][]string{
		"dobry":    {"zły"},
		"zły":      {"dobry"},
		"duży":     {"mały"},
		"mały":     {"duży"},
		"tak":      {"nie"},
		"nie":      {"tak"},
		"wysoki":   {"niski"},
		"niski":    {"wysoki"},
		"gorący":   {"zimny"},
		"zimny":    {"gorący"},
		"nowy":     {"stary"},
		"stary":    {"nowy"},
		"szybki":   {"wolny"},
		"wolny":    {"szybki"},
		"łatwy":    {"trudny"},
		"trudny":   {"łatwy"},
		"jasny":    {"ciemny"},
		"ciemny":   {"jasny"},
		"otwarty":  {"zamknięty"},
		"zamknięty": {"otwarty"},
		"wewnątrz": {"zewnątrz"},
		"zewnątrz": {"wewnątrz"},
		"góra":     {"dół"},
		"dół":      {"góra"},
		"przed":    {"po"},
		"po":       {"przed"},
		"wcześnie": {"późno"},
		"późno":    {"wcześnie"},
		"początek": {"koniec"},
		"koniec":   {"początek"},
		"prawda":   {"fałsz"},
		"fałsz":    {"prawda"},
		"bogaty":   {"biedny"},
		"biedny":   {"bogaty"},
		"silny":    {"słaby"},
		"słaby":    {"silny"},
		"pełny":    {"pusty"},
		"pusty":    {"pełny"},
		"blisko":   {"daleko"},
		"daleko":   {"blisko"},
		"kupić":    {"sprzedać"},
		"sprzedać": {"kupić"},
		"dać":      {"wziąć"},
		"wziąć":    {"dać"},
		"kochać":   {"nienawidzić"},
		"nienawidzić": {"kochać"},
		"wygrać":   {"przegrać"},
		"przegrać": {"wygrać"},
		"dzień":    {"noc"},
		"noc":      {"dzień"},
	}
	Language = Polish{
		code:        LANGUAGE,
		name:        "Polish",
		description: "Polish is a West Slavic language spoken primarily in Poland",
		numerals: map[string][]string{
			"0":          {"zero"},
			"1":          {"jeden", "pierwszy"},
			"2":          {"dwa", "drugi"},
			"3":          {"trzy", "trzeci"},
			"4":          {"cztery", "czwarty"},
			"5":          {"pięć", "piąty"},
			"6":          {"sześć", "szósty"},
			"7":          {"siedem", "siódmy"},
			"8":          {"osiem", "ósmy"},
			"9":          {"dziewięć", "dziewiąty"},
			"10":         {"dziesięć", "dziesiąty"},
			"11":         {"jedenaście", "jedenasty"},
			"12":         {"dwanaście", "dwunasty"},
			"13":         {"trzynaście", "trzynasty"},
			"14":         {"czternaście", "czternasty"},
			"15":         {"piętnaście", "piętnasty"},
			"16":         {"szesnaście", "szesnasty"},
			"17":         {"siedemnaście", "siedemnasty"},
			"18":         {"osiemnaście", "osiemnasty"},
			"19":         {"dziewiętnaście", "dziewiętnasty"},
			"20":         {"dwadzieścia", "dwudziesty"},
			"30":         {"trzydzieści", "trzydziesty"},
			"40":         {"czterdzieści", "czterdziesty"},
			"50":         {"pięćdziesiąt", "pięćdziesiąty"},
			"60":         {"sześćdziesiąt", "sześćdziesiąty"},
			"70":         {"siedemdziesiąt", "siedemdziesiąty"},
			"80":         {"osiemdziesiąt", "osiemdziesiąty"},
			"90":         {"dziewięćdziesiąt", "dziewięćdziesiąty"},
			"100":        {"sto", "setny"},
			"1000":       {"tysiąc", "tysięczny"},
			"1000000":    {"milion", "milionowy"},
			"1000000000": {"miliard", "miliardowy"},
		},
		graphemes: []string{
			"a", "ą", "b", "c", "ć", "d", "e", "ę", "f", "g", "h", "i", "j", "k", "l", "ł", "m",
			"n", "ń", "o", "ó", "p", "r", "s", "ś", "t", "u", "w", "y", "z", "ź", "ż",
		},
		vowels:       []string{"a", "e", "i", "o", "u", "y", "ą", "ę", "ó"},
		misspellings: plMisspellings,
		homophones:   plHomophones,
		antonyms:     plAntonyms,
		homoglyphs:   languages.DefaultLatinHomoglyphs(),
	}
)

func init() { languages.AddLanguage(LANGUAGE, func() internal.Language { return &Language }) }
