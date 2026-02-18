package danish

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
)

const LANGUAGE string = "da"

type Danish struct {
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

func (l *Danish) Id() string                      { return l.code }
func (l *Danish) Name() string                    { return l.name }
func (l *Danish) Description() string             { return l.description }
func (l *Danish) Numerals() map[string][]string   { return l.numerals }
func (l *Danish) Cardinal() map[string]string     { return languages.NumeralMap(l.numerals, 0) }
func (l *Danish) Ordinal() map[string]string      { return languages.NumeralMap(l.numerals, 1) }
func (l *Danish) Graphemes() []string             { return l.graphemes }
func (l *Danish) Vowels() []string                { return l.vowels }
func (l *Danish) Misspellings() [][]string        { return l.misspellings }
func (l *Danish) Homophones() [][]string          { return l.homophones }
func (l *Danish) Antonyms() map[string][]string   { return l.antonyms }
func (l *Danish) Homoglyphs() map[string][]string { return l.homoglyphs }
func (l *Danish) SimilarChars(char string) []string {
	return languages.SimilarChars(l.homoglyphs, char)
}
func (l *Danish) SimilarSpellings(word string) []string {
	return languages.SimilarSpellings(l.misspellings, word)
}
func (l *Danish) SimilarSounds(word string) []string {
	return languages.SimilarSounds(l.homophones, word)
}
func (l *Danish) Keyboards() (boards []internal.Keyboard) {
	for _, b := range languages.Keyboards() {
		if b.Language() == l.code {
			boards = append(boards, b)
		}
	}
	return
}

var (
	daMisspellings = [][]string{}
	daHomophones   = [][]string{
		{"punktum", "."},
		{"snabela", "@"},
		{"bindestreg", "-"},
		{"skraastreg", "/", "skråstreg"},
	}
	daAntonyms = map[string][]string{
		"god":       {"dårlig"},
		"dårlig":    {"god"},
		"stor":      {"lille"},
		"lille":     {"stor"},
		"ja":        {"nej"},
		"nej":       {"ja"},
		"høj":       {"lav"},
		"lav":       {"høj"},
		"varm":      {"kold"},
		"kold":      {"varm"},
		"ny":        {"gammel"},
		"gammel":    {"ny"},
		"hurtig":    {"langsom"},
		"langsom":   {"hurtig"},
		"let":       {"svær"},
		"svær":      {"let"},
		"lys":       {"mørk"},
		"mørk":      {"lys"},
		"åben":      {"lukket"},
		"lukket":    {"åben"},
		"inde":      {"ude"},
		"ude":       {"inde"},
		"op":        {"ned"},
		"ned":       {"op"},
		"før":       {"efter"},
		"efter":     {"før"},
		"tidlig":    {"sent"},
		"sent":      {"tidlig"},
		"start":     {"slut"},
		"slut":      {"start"},
		"sand":      {"falsk"},
		"falsk":     {"sand"},
		"rig":       {"fattig"},
		"fattig":    {"rig"},
		"stærk":     {"svag"},
		"svag":      {"stærk"},
		"fuld":      {"tom"},
		"tom":       {"fuld"},
		"nær":       {"fjern"},
		"fjern":     {"nær"},
		"købe":      {"sælge"},
		"sælge":     {"købe"},
		"give":      {"tage"},
		"tage":      {"give"},
		"elske":     {"hade"},
		"hade":      {"elske"},
		"vinde":     {"tabe"},
		"tabe":      {"vinde"},
		"dag":       {"nat"},
		"nat":       {"dag"},
	}
	Language = Danish{
		code:        LANGUAGE,
		name:        "Danish",
		description: "Danish is a North Germanic language spoken primarily in Denmark",
		numerals: map[string][]string{
			"0":          {"nul"},
			"1":          {"en", "første"},
			"2":          {"to", "anden"},
			"3":          {"tre", "tredje"},
			"4":          {"fire", "fjerde"},
			"5":          {"fem", "femte"},
			"6":          {"seks", "sjette"},
			"7":          {"syv", "syvende"},
			"8":          {"otte", "ottende"},
			"9":          {"ni", "niende"},
			"10":         {"ti", "tiende"},
			"11":         {"elleve", "elfte"},
			"12":         {"tolv", "tolvte"},
			"13":         {"tretten", "trettende"},
			"14":         {"fjorten", "fjortende"},
			"15":         {"femten", "femtende"},
			"16":         {"seksten", "sekstende"},
			"17":         {"sytten", "syttende"},
			"18":         {"atten", "attende"},
			"19":         {"nitten", "nittende"},
			"20":         {"tyve", "tyvende"},
			"30":         {"tredive"},
			"40":         {"fyrre"},
			"50":         {"halvtreds"},
			"60":         {"tres"},
			"70":         {"halvfjerds"},
			"80":         {"firs"},
			"90":         {"halvfems"},
			"100":        {"hundrede"},
			"1000":       {"tusind"},
			"1000000":    {"million"},
			"1000000000": {"milliard"},
		},
		graphemes: []string{
			"a", "b", "c", "d", "e", "f", "g",
			"h", "i", "j", "k", "l", "m", "n",
			"o", "p", "q", "r", "s", "t", "u",
			"v", "w", "x", "y", "z", "æ", "ø", "å",
		},
		vowels:       []string{"a", "e", "i", "o", "u", "y", "æ", "ø", "å"},
		misspellings: daMisspellings,
		homophones:   daHomophones,
		antonyms:     daAntonyms,
		homoglyphs:   languages.DefaultLatinHomoglyphs(),
	}
)

func init() { languages.AddLanguage(LANGUAGE, func() internal.Language { return &Language }) }
