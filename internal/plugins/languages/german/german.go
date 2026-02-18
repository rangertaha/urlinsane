package german

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
)

const LANGUAGE string = "de"

type German struct {
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

func (l *German) Id() string                      { return l.code }
func (l *German) Name() string                    { return l.name }
func (l *German) Description() string             { return l.description }
func (l *German) Numerals() map[string][]string   { return l.numerals }
func (l *German) Cardinal() map[string]string     { return languages.NumeralMap(l.numerals, 0) }
func (l *German) Ordinal() map[string]string      { return languages.NumeralMap(l.numerals, 1) }
func (l *German) Graphemes() []string             { return l.graphemes }
func (l *German) Vowels() []string                { return l.vowels }
func (l *German) Misspellings() [][]string        { return l.misspellings }
func (l *German) Homophones() [][]string          { return l.homophones }
func (l *German) Antonyms() map[string][]string   { return l.antonyms }
func (l *German) Homoglyphs() map[string][]string { return l.homoglyphs }
func (l *German) SimilarChars(char string) []string {
	return languages.SimilarChars(l.homoglyphs, char)
}
func (l *German) SimilarSpellings(word string) []string {
	return languages.SimilarSpellings(l.misspellings, word)
}
func (l *German) SimilarSounds(word string) []string {
	return languages.SimilarSounds(l.homophones, word)
}
func (l *German) Keyboards() (boards []internal.Keyboard) {
	for _, b := range languages.Keyboards() {
		if b.Language() == l.code {
			boards = append(boards, b)
		}
	}
	return
}

var (
	deMisspellings = [][]string{
		{"uber", "über"},
		{"strasse", "straße"},
		{"fuenf", "fünf"},
		{"dass", "das"},
	}
	deHomophones = [][]string{
		{"punkt", "."},
		{"klammeraffe", "@"},
		{"bindestrich", "-"},
		{"slash", "/"},
		{"seit", "seid"},
	}
	deAntonyms = map[string][]string{
		"gut":        {"schlecht"},
		"schlecht":   {"gut"},
		"groß":       {"klein"},
		"gross":      {"klein"},
		"klein":      {"groß", "gross"},
		"hoch":       {"niedrig"},
		"niedrig":    {"hoch"},
		"ja":         {"nein"},
		"nein":       {"ja"},
		"heiß":       {"kalt"},
		"kalt":       {"heiß"},
		"neu":        {"alt"},
		"alt":        {"neu"},
		"schnell":    {"langsam"},
		"langsam":    {"schnell"},
		"leicht":     {"schwer"},
		"schwer":     {"leicht"},
		"hell":       {"dunkel"},
		"dunkel":     {"hell"},
		"offen":      {"geschlossen"},
		"geschlossen": {"offen"},
		"innen":      {"außen"},
		"außen":      {"innen"},
		"oben":       {"unten"},
		"unten":      {"oben"},
		"vor":        {"nach"},
		"nach":       {"vor"},
		"früh":       {"spät"},
		"spät":       {"früh"},
		"anfangen":   {"aufhören"},
		"aufhören":   {"anfangen"},
		"wahr":       {"falsch"},
		"falsch":     {"wahr"},
		"reich":      {"arm"},
		"arm":        {"reich"},
		"stark":      {"schwach"},
		"schwach":    {"stark"},
		"voll":       {"leer"},
		"leer":       {"voll"},
		"nah":        {"fern"},
		"fern":       {"nah"},
		"kaufen":     {"verkaufen"},
		"verkaufen":  {"kaufen"},
		"geben":      {"nehmen"},
		"nehmen":     {"geben"},
		"lieben":     {"hassen"},
		"hassen":     {"lieben"},
		"gewinnen":   {"verlieren"},
		"verlieren":  {"gewinnen"},
		"tag":        {"nacht"},
		"nacht":      {"tag"},
	}

	Language = German{
		code:        LANGUAGE,
		name:        "German",
		description: "German is a West Germanic language spoken primarily in Central Europe",
		numerals: map[string][]string{
			"0":          {"null"},
			"1":          {"eins", "erste"},
			"2":          {"zwei", "zweite"},
			"3":          {"drei", "dritte"},
			"4":          {"vier", "vierte"},
			"5":          {"fünf", "fünfte"},
			"6":          {"sechs", "sechste"},
			"7":          {"sieben", "siebte"},
			"8":          {"acht", "achte"},
			"9":          {"neun", "neunte"},
			"10":         {"zehn", "zehnte"},
			"11":         {"elf", "elfte"},
			"12":         {"zwölf", "zwölfte"},
			"13":         {"dreizehn", "dreizehnte"},
			"14":         {"vierzehn", "vierzehnte"},
			"15":         {"fünfzehn", "fünfzehnte"},
			"16":         {"sechzehn", "sechzehnte"},
			"17":         {"siebzehn", "siebzehnte"},
			"18":         {"achtzehn", "achtzehnte"},
			"19":         {"neunzehn", "neunzehnte"},
			"20":         {"zwanzig", "zwanzigste"},
			"30":         {"dreißig", "dreißigste"},
			"40":         {"vierzig", "vierzigste"},
			"50":         {"fünfzig", "fünfzigste"},
			"60":         {"sechzig", "sechzigste"},
			"70":         {"siebzig", "siebzigste"},
			"80":         {"achtzig", "achtzigste"},
			"90":         {"neunzig", "neunzigste"},
			"100":        {"hundert", "hundertste"},
			"1000":       {"tausend", "tausendste"},
			"1000000":    {"million", "millionste"},
			"1000000000": {"milliarde", "milliardste"},
		},
		graphemes: []string{
			"a", "b", "c", "d", "e", "f", "g",
			"h", "i", "j", "k", "l", "m", "n",
			"o", "p", "q", "r", "s", "t", "u",
			"v", "w", "x", "y", "z", "ä", "ö", "ü", "ß",
		},
		vowels:       []string{"a", "e", "i", "o", "u", "ä", "ö", "ü"},
		misspellings: deMisspellings,
		homophones:   deHomophones,
		antonyms:     deAntonyms,
		homoglyphs:   languages.DefaultLatinHomoglyphs(),
	}
)

func init() {
	languages.AddLanguage(LANGUAGE, func() internal.Language { return &Language })
}
