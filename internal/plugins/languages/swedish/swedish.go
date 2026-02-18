package swedish

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
)

const LANGUAGE string = "sv"

type Swedish struct {
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

func (l *Swedish) Id() string                      { return l.code }
func (l *Swedish) Name() string                    { return l.name }
func (l *Swedish) Description() string             { return l.description }
func (l *Swedish) Numerals() map[string][]string   { return l.numerals }
func (l *Swedish) Cardinal() map[string]string     { return languages.NumeralMap(l.numerals, 0) }
func (l *Swedish) Ordinal() map[string]string      { return languages.NumeralMap(l.numerals, 1) }
func (l *Swedish) Graphemes() []string             { return l.graphemes }
func (l *Swedish) Vowels() []string                { return l.vowels }
func (l *Swedish) Misspellings() [][]string        { return l.misspellings }
func (l *Swedish) Homophones() [][]string          { return l.homophones }
func (l *Swedish) Antonyms() map[string][]string   { return l.antonyms }
func (l *Swedish) Homoglyphs() map[string][]string { return l.homoglyphs }
func (l *Swedish) SimilarChars(char string) []string {
	return languages.SimilarChars(l.homoglyphs, char)
}
func (l *Swedish) SimilarSpellings(word string) []string {
	return languages.SimilarSpellings(l.misspellings, word)
}
func (l *Swedish) SimilarSounds(word string) []string {
	return languages.SimilarSounds(l.homophones, word)
}
func (l *Swedish) Keyboards() (boards []internal.Keyboard) {
	for _, b := range languages.Keyboards() {
		if b.Language() == l.code {
			boards = append(boards, b)
		}
	}
	return
}

var (
	svMisspellings = [][]string{}
	svHomophones   = [][]string{
		{"punkt", "."},
	}
	svAntonyms = map[string][]string{
		"bra":      {"dålig"},
		"dålig":    {"bra"},
		"stor":     {"liten"},
		"liten":    {"stor"},
		"ja":       {"nej"},
		"nej":      {"ja"},
		"hög":      {"låg"},
		"låg":      {"hög"},
		"varm":     {"kall"},
		"kall":     {"varm"},
		"ny":       {"gammal"},
		"gammal":   {"ny"},
		"snabb":    {"långsam"},
		"långsam":  {"snabb"},
		"lätt":     {"svår"},
		"svår":     {"lätt"},
		"ljus":     {"mörk"},
		"mörk":     {"ljus"},
		"öppen":    {"stängd"},
		"stängd":   {"öppen"},
		"inne":     {"ute"},
		"ute":      {"inne"},
		"upp":      {"ner"},
		"ner":      {"upp"},
		"före":     {"efter"},
		"efter":    {"före"},
		"tidig":    {"sen"},
		"sen":      {"tidig"},
		"start":    {"slut"},
		"slut":     {"start"},
		"sann":     {"falsk"},
		"falsk":    {"sann"},
		"rik":      {"fattig"},
		"fattig":   {"rik"},
		"stark":    {"svag"},
		"svag":     {"stark"},
		"full":     {"tom"},
		"tom":      {"full"},
		"nära":     {"fjärran"},
		"fjärran":  {"nära"},
		"köpa":     {"sälja"},
		"sälja":    {"köpa"},
		"ge":       {"ta"},
		"ta":       {"ge"},
		"älska":    {"hata"},
		"hata":     {"älska"},
		"vinna":    {"förlora"},
		"förlora":  {"vinna"},
		"dag":      {"natt"},
		"natt":     {"dag"},
	}
	Language = Swedish{
		code:        LANGUAGE,
		name:        "Swedish",
		description: "Swedish is a North Germanic language spoken primarily in Sweden",
		numerals: map[string][]string{
			"0":          {"noll"},
			"1":          {"ett", "första"},
			"2":          {"två", "andra"},
			"3":          {"tre", "tredje"},
			"4":          {"fyra", "fjärde"},
			"5":          {"fem", "femte"},
			"6":          {"sex", "sjätte"},
			"7":          {"sju", "sjunde"},
			"8":          {"åtta", "åttonde"},
			"9":          {"nio", "nionde"},
			"10":         {"tio", "tionde"},
			"11":         {"elva", "elfte"},
			"12":         {"tolv", "tolfte"},
			"13":         {"tretton", "trettonde"},
			"14":         {"fjorton", "fjortonde"},
			"15":         {"femton", "femtonde"},
			"16":         {"sexton", "sextonde"},
			"17":         {"sjutton", "sjuttonde"},
			"18":         {"arton", "artonde"},
			"19":         {"nitton", "nittonde"},
			"20":         {"tjugo", "tjugonde"},
			"30":         {"trettio", "trettionde"},
			"40":         {"fyrtio", "fyrtionde"},
			"50":         {"femtio", "femtionde"},
			"60":         {"sextio", "sextionde"},
			"70":         {"sjuttio", "sjuttionde"},
			"80":         {"åttio", "åttionde"},
			"90":         {"nittio", "nittionde"},
			"100":        {"hundra", "hundrade"},
			"1000":       {"tusen", "tusende"},
			"1000000":    {"miljon"},
			"1000000000": {"miljard"},
		},
		graphemes: []string{
			"a", "b", "c", "d", "e", "f", "g",
			"h", "i", "j", "k", "l", "m", "n",
			"o", "p", "q", "r", "s", "t", "u",
			"v", "w", "x", "y", "z", "å", "ä", "ö",
		},
		vowels:       []string{"a", "e", "i", "o", "u", "y", "å", "ä", "ö"},
		misspellings: svMisspellings,
		homophones:   svHomophones,
		antonyms:     svAntonyms,
		homoglyphs:   languages.DefaultLatinHomoglyphs(),
	}
)

func init() { languages.AddLanguage(LANGUAGE, func() internal.Language { return &Language }) }
