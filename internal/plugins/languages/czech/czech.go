package czech

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
)

const LANGUAGE string = "cs"

type Czech struct {
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

func (l *Czech) Id() string                        { return l.code }
func (l *Czech) Name() string                      { return l.name }
func (l *Czech) Description() string               { return l.description }
func (l *Czech) Numerals() map[string][]string     { return l.numerals }
func (l *Czech) Cardinal() map[string]string       { return languages.NumeralMap(l.numerals, 0) }
func (l *Czech) Ordinal() map[string]string        { return languages.NumeralMap(l.numerals, 1) }
func (l *Czech) Graphemes() []string               { return l.graphemes }
func (l *Czech) Vowels() []string                  { return l.vowels }
func (l *Czech) Misspellings() [][]string          { return l.misspellings }
func (l *Czech) Homophones() [][]string            { return l.homophones }
func (l *Czech) Antonyms() map[string][]string     { return l.antonyms }
func (l *Czech) Homoglyphs() map[string][]string   { return l.homoglyphs }
func (l *Czech) SimilarChars(char string) []string { return languages.SimilarChars(l.homoglyphs, char) }
func (l *Czech) SimilarSpellings(word string) []string {
	return languages.SimilarSpellings(l.misspellings, word)
}
func (l *Czech) SimilarSounds(word string) []string {
	return languages.SimilarSounds(l.homophones, word)
}
func (l *Czech) Keyboards() (boards []internal.Keyboard) {
	for _, b := range languages.Keyboards() {
		if b.Language() == l.code {
			boards = append(boards, b)
		}
	}
	return
}

var (
	csMisspellings = [][]string{
		{"praha", "práha"},
	}
	csHomophones = [][]string{
		{"tecka", "."},
		{"tečka", "."},
	}
	csAntonyms = map[string][]string{
		"dobrý":     {"špatný"},
		"špatný":    {"dobrý"},
		"velký":     {"malý"},
		"malý":      {"velký"},
		"ano":       {"ne"},
		"ne":        {"ano"},
		"vysoký":    {"nízký"},
		"nízký":     {"vysoký"},
		"horký":     {"studený"},
		"studený":   {"horký"},
		"nový":      {"starý"},
		"starý":     {"nový"},
		"rychlý":    {"pomalý"},
		"pomalý":    {"rychlý"},
		"snadný":    {"těžký"},
		"těžký":     {"snadný"},
		"světlý":    {"tmavý"},
		"tmavý":     {"světlý"},
		"otevřený":  {"zavřený"},
		"zavřený":   {"otevřený"},
		"uvnitř":    {"venku"},
		"venku":     {"uvnitř"},
		"nahoře":    {"dole"},
		"dole":      {"nahoře"},
		"před":      {"po"},
		"po":        {"před"},
		"brzy":      {"pozdě"},
		"pozdě":     {"brzy"},
		"začátek":   {"konec"},
		"konec":     {"začátek"},
		"pravda":    {"lež"},
		"lež":       {"pravda"},
		"bohatý":    {"chudý"},
		"chudý":     {"bohatý"},
		"silný":     {"slabý"},
		"slabý":     {"silný"},
		"plný":      {"prázdný"},
		"prázdný":   {"plný"},
		"blízko":    {"daleko"},
		"daleko":    {"blízko"},
		"koupit":    {"prodat"},
		"prodat":    {"koupit"},
		"dát":       {"vzít"},
		"vzít":      {"dát"},
		"milovat":   {"nenávidět"},
		"nenávidět": {"milovat"},
		"vyhrát":    {"prohrát"},
		"prohrát":   {"vyhrát"},
		"den":       {"noc"},
		"noc":       {"den"},
	}

	Language = Czech{
		code:        LANGUAGE,
		name:        "Czech",
		description: "Czech is a West Slavic language spoken primarily in the Czech Republic",
		numerals: map[string][]string{
			"0":          {"nula"},
			"1":          {"jeden", "první"},
			"2":          {"dva", "druhý"},
			"3":          {"tři", "třetí"},
			"4":          {"čtyři", "čtvrtý"},
			"5":          {"pět", "pátý"},
			"6":          {"šest", "šestý"},
			"7":          {"sedm", "sedmý"},
			"8":          {"osm", "osmý"},
			"9":          {"devět", "devátý"},
			"10":         {"deset", "desátý"},
			"11":         {"jedenáct", "jedenáctý"},
			"12":         {"dvanáct", "dvanáctý"},
			"13":         {"třináct", "třináctý"},
			"14":         {"čtrnáct", "čtrnáctý"},
			"15":         {"patnáct", "patnáctý"},
			"16":         {"šestnáct", "šestnáctý"},
			"17":         {"sedmnáct", "sedmnáctý"},
			"18":         {"osmnáct", "osmnáctý"},
			"19":         {"devatenáct", "devatenáctý"},
			"20":         {"dvacet", "dvacátý"},
			"30":         {"třicet", "třicátý"},
			"40":         {"čtyřicet", "čtyřicátý"},
			"50":         {"padesát", "padesátý"},
			"60":         {"šedesát", "šedesátý"},
			"70":         {"sedmdesát", "sedmdesátý"},
			"80":         {"osmdesát", "osmdesátý"},
			"90":         {"devadesát", "devadesátý"},
			"100":        {"sto", "stý"},
			"1000":       {"tisíc", "tisící"},
			"1000000":    {"milion", "miliontý"},
			"1000000000": {"miliarda", "miliardtý"},
		},
		graphemes: []string{
			"a", "á", "b", "c", "č", "d", "ď", "e", "é", "ě", "f", "g", "h", "ch", "i", "í",
			"j", "k", "l", "m", "n", "ň", "o", "ó", "p", "q", "r", "ř", "s", "š", "t", "ť",
			"u", "ú", "ů", "v", "w", "x", "y", "ý", "z", "ž",
		},
		vowels:       []string{"a", "e", "i", "o", "u", "y", "á", "é", "í", "ó", "ú", "ů", "ý", "ě"},
		misspellings: csMisspellings,
		homophones:   csHomophones,
		antonyms:     csAntonyms,
		homoglyphs:   languages.DefaultLatinHomoglyphs(),
	}
)

func init() { languages.AddLanguage(LANGUAGE, func() internal.Language { return &Language }) }
