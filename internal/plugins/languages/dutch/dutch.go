package dutch

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
)

const LANGUAGE string = "nl"

type Dutch struct {
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

func (l *Dutch) Id() string                        { return l.code }
func (l *Dutch) Name() string                      { return l.name }
func (l *Dutch) Description() string               { return l.description }
func (l *Dutch) Numerals() map[string][]string     { return l.numerals }
func (l *Dutch) Cardinal() map[string]string       { return languages.NumeralMap(l.numerals, 0) }
func (l *Dutch) Ordinal() map[string]string        { return languages.NumeralMap(l.numerals, 1) }
func (l *Dutch) Graphemes() []string               { return l.graphemes }
func (l *Dutch) Vowels() []string                  { return l.vowels }
func (l *Dutch) Misspellings() [][]string          { return l.misspellings }
func (l *Dutch) Homophones() [][]string            { return l.homophones }
func (l *Dutch) Antonyms() map[string][]string     { return l.antonyms }
func (l *Dutch) Homoglyphs() map[string][]string   { return l.homoglyphs }
func (l *Dutch) SimilarChars(char string) []string { return languages.SimilarChars(l.homoglyphs, char) }
func (l *Dutch) SimilarSpellings(word string) []string {
	return languages.SimilarSpellings(l.misspellings, word)
}
func (l *Dutch) SimilarSounds(word string) []string {
	return languages.SimilarSounds(l.homophones, word)
}
func (l *Dutch) Keyboards() (boards []internal.Keyboard) {
	for _, b := range languages.Keyboards() {
		if b.Language() == l.code {
			boards = append(boards, b)
		}
	}
	return
}

var (
	nlMisspellings = [][]string{
		{"zo", "zó"},
	}
	nlHomophones = [][]string{
		{"punt", "."},
		{"apenstaartje", "@"},
		{"streepje", "-"},
		{"slash", "/"},
	}
	nlAntonyms = map[string][]string{
		"goed":      {"slecht"},
		"slecht":    {"goed"},
		"groot":     {"klein"},
		"klein":     {"groot"},
		"hoog":      {"laag"},
		"laag":      {"hoog"},
		"ja":        {"nee"},
		"nee":       {"ja"},
		"heet":      {"koud"},
		"koud":      {"heet"},
		"nieuw":     {"oud"},
		"oud":       {"nieuw"},
		"snel":      {"langzaam"},
		"langzaam":  {"snel"},
		"makkelijk": {"moeilijk"},
		"moeilijk":  {"makkelijk"},
		"licht":     {"donker"},
		"donker":    {"licht"},
		"open":      {"gesloten"},
		"gesloten":  {"open"},
		"binnen":    {"buiten"},
		"buiten":    {"binnen"},
		"boven":     {"onder"},
		"onder":     {"boven"},
		"voor":      {"na"},
		"na":        {"voor"},
		"vroeg":     {"laat"},
		"laat":      {"vroeg"},
		"beginnen":  {"eindigen"},
		"eindigen":  {"beginnen"},
		"waar":      {"onwaar"},
		"onwaar":    {"waar"},
		"rijk":      {"arm"},
		"arm":       {"rijk"},
		"sterk":     {"zwak"},
		"zwak":      {"sterk"},
		"vol":       {"leeg"},
		"leeg":      {"vol"},
		"dichtbij":  {"ver"},
		"ver":       {"dichtbij"},
		"kopen":     {"verkopen"},
		"verkopen":  {"kopen"},
		"geven":     {"nemen"},
		"nemen":     {"geven"},
		"liefhebben": {"haten"},
		"haten":      {"liefhebben"},
		"winnen":    {"verliezen"},
		"verliezen": {"winnen"},
		"dag":       {"nacht"},
		"nacht":     {"dag"},
	}
	Language = Dutch{
		code:        LANGUAGE,
		name:        "Dutch",
		description: "Dutch is a West Germanic language spoken in the Netherlands and Belgium",
		numerals: map[string][]string{
			"0":          {"nul"},
			"1":          {"een", "eerste"},
			"2":          {"twee", "tweede"},
			"3":          {"drie", "derde"},
			"4":          {"vier", "vierde"},
			"5":          {"vijf", "vijfde"},
			"6":          {"zes", "zesde"},
			"7":          {"zeven", "zevende"},
			"8":          {"acht", "achtste"},
			"9":          {"negen", "negende"},
			"10":         {"tien", "tiende"},
			"11":         {"elf", "elfde"},
			"12":         {"twaalf", "twaalfde"},
			"13":         {"dertien", "dertiende"},
			"14":         {"veertien", "veertiende"},
			"15":         {"vijftien", "vijftiende"},
			"16":         {"zestien", "zestiende"},
			"17":         {"zeventien", "zeventiende"},
			"18":         {"achttien", "achttiende"},
			"19":         {"negentien", "negentiende"},
			"20":         {"twintig", "twintigste"},
			"30":         {"dertig", "dertigste"},
			"40":         {"veertig", "veertigste"},
			"50":         {"vijftig", "vijftigste"},
			"60":         {"zestig", "zestigste"},
			"70":         {"zeventig", "zeventigste"},
			"80":         {"tachtig", "tachtigste"},
			"90":         {"negentig", "negentigste"},
			"100":        {"honderd", "honderdste"},
			"1000":       {"duizend", "duizendste"},
			"1000000":    {"miljoen", "miljoenste"},
			"1000000000": {"miljard", "miljardste"},
		},
		graphemes: []string{
			"a", "b", "c", "d", "e", "f", "g",
			"h", "i", "j", "k", "l", "m", "n",
			"o", "p", "q", "r", "s", "t", "u",
			"v", "w", "x", "y", "z", "ë", "ï", "á", "é", "í", "ó", "ú",
		},
		vowels:       []string{"a", "e", "i", "o", "u", "y"},
		misspellings: nlMisspellings,
		homophones:   nlHomophones,
		antonyms:     nlAntonyms,
		homoglyphs:   languages.DefaultLatinHomoglyphs(),
	}
)

func init() { languages.AddLanguage(LANGUAGE, func() internal.Language { return &Language }) }
