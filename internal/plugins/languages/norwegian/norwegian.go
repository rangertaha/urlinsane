package norwegian

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
)

const LANGUAGE string = "no"

type Norwegian struct {
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

func (l *Norwegian) Id() string                      { return l.code }
func (l *Norwegian) Name() string                    { return l.name }
func (l *Norwegian) Description() string             { return l.description }
func (l *Norwegian) Numerals() map[string][]string   { return l.numerals }
func (l *Norwegian) Cardinal() map[string]string     { return languages.NumeralMap(l.numerals, 0) }
func (l *Norwegian) Ordinal() map[string]string      { return languages.NumeralMap(l.numerals, 1) }
func (l *Norwegian) Graphemes() []string             { return l.graphemes }
func (l *Norwegian) Vowels() []string                { return l.vowels }
func (l *Norwegian) Misspellings() [][]string        { return l.misspellings }
func (l *Norwegian) Homophones() [][]string          { return l.homophones }
func (l *Norwegian) Antonyms() map[string][]string   { return l.antonyms }
func (l *Norwegian) Homoglyphs() map[string][]string { return l.homoglyphs }
func (l *Norwegian) SimilarChars(char string) []string {
	return languages.SimilarChars(l.homoglyphs, char)
}
func (l *Norwegian) SimilarSpellings(word string) []string {
	return languages.SimilarSpellings(l.misspellings, word)
}
func (l *Norwegian) SimilarSounds(word string) []string {
	return languages.SimilarSounds(l.homophones, word)
}
func (l *Norwegian) Keyboards() (boards []internal.Keyboard) {
	for _, b := range languages.Keyboards() {
		if b.Language() == l.code {
			boards = append(boards, b)
		}
	}
	return
}

var (
	noMisspellings = [][]string{
		{"norge", "Norge"},
	}
	noHomophones = [][]string{
		{"punktum", "."},
		{"krøllalfa", "@"},
		{"bindestrek", "-"},
		{"skråstrek", "/"},
	}
	noAntonyms = map[string][]string{
		"bra":       {"dårlig"},
		"dårlig":    {"bra"},
		"stor":      {"liten"},
		"liten":     {"stor"},
		"ja":        {"nei"},
		"nei":       {"ja"},
		"høy":       {"lav"},
		"lav":       {"høy"},
		"varm":      {"kald"},
		"kald":      {"varm"},
		"ny":        {"gammel"},
		"gammel":    {"ny"},
		"rask":      {"langsom"},
		"langsom":   {"rask"},
		"lett":      {"vanskelig"},
		"vanskelig": {"lett"},
		"lys":       {"mørk"},
		"mørk":      {"lys"},
		"åpen":      {"lukket"},
		"lukket":    {"åpen"},
		"inne":      {"ute"},
		"ute":       {"inne"},
		"opp":       {"ned"},
		"ned":       {"opp"},
		"før":       {"etter"},
		"etter":     {"før"},
		"tidlig":    {"sent"},
		"sent":      {"tidlig"},
		"start":     {"slutt"},
		"slutt":     {"start"},
		"sann":      {"falsk"},
		"falsk":     {"sann"},
		"rik":       {"fattig"},
		"fattig":    {"rik"},
		"sterk":     {"svak"},
		"svak":      {"sterk"},
		"full":      {"tom"},
		"tom":       {"full"},
		"nær":       {"fjern"},
		"fjern":     {"nær"},
		"kjøpe":     {"selge"},
		"selge":     {"kjøpe"},
		"gi":        {"ta"},
		"ta":        {"gi"},
		"elske":     {"hate"},
		"hate":      {"elske"},
		"vinne":     {"tape"},
		"tape":      {"vinne"},
		"dag":       {"natt"},
		"natt":      {"dag"},
	}
	Language = Norwegian{
		code:        LANGUAGE,
		name:        "Norwegian",
		description: "Norwegian is a North Germanic language spoken primarily in Norway",
		numerals: map[string][]string{
			"0":          {"null"},
			"1":          {"en", "første"},
			"2":          {"to", "andre"},
			"3":          {"tre", "tredje"},
			"4":          {"fire", "fjerde"},
			"5":          {"fem", "femte"},
			"6":          {"seks", "sjette"},
			"7":          {"sju", "sjuende"},
			"8":          {"åtte", "åttende"},
			"9":          {"ni", "niende"},
			"10":         {"ti", "tiende"},
			"11":         {"elleve", "ellevte"},
			"12":         {"tolv", "tolvte"},
			"13":         {"tretten", "trettende"},
			"14":         {"fjorten", "fjortende"},
			"15":         {"femten", "femtende"},
			"16":         {"seksten", "sekstende"},
			"17":         {"sytten", "syttende"},
			"18":         {"atten", "attende"},
			"19":         {"nitten", "nittende"},
			"20":         {"tjue", "tjuende"},
			"30":         {"tretti", "trettiende"},
			"40":         {"førti", "førtiende"},
			"50":         {"femti", "femtiende"},
			"60":         {"seksti", "sekstiende"},
			"70":         {"sytti", "syttiende"},
			"80":         {"åtti", "åttiende"},
			"90":         {"nitti", "nittiende"},
			"100":        {"hundre", "hundrede"},
			"1000":       {"tusen", "tusende"},
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
		misspellings: noMisspellings,
		homophones:   noHomophones,
		antonyms:     noAntonyms,
		homoglyphs:   languages.DefaultLatinHomoglyphs(),
	}
)

func init() { languages.AddLanguage(LANGUAGE, func() internal.Language { return &Language }) }
