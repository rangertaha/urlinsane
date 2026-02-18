package italian

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
)

const LANGUAGE string = "it"

type Italian struct {
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

func (l *Italian) Id() string                      { return l.code }
func (l *Italian) Name() string                    { return l.name }
func (l *Italian) Description() string             { return l.description }
func (l *Italian) Numerals() map[string][]string   { return l.numerals }
func (l *Italian) Cardinal() map[string]string     { return languages.NumeralMap(l.numerals, 0) }
func (l *Italian) Ordinal() map[string]string      { return languages.NumeralMap(l.numerals, 1) }
func (l *Italian) Graphemes() []string             { return l.graphemes }
func (l *Italian) Vowels() []string                { return l.vowels }
func (l *Italian) Misspellings() [][]string        { return l.misspellings }
func (l *Italian) Homophones() [][]string          { return l.homophones }
func (l *Italian) Antonyms() map[string][]string   { return l.antonyms }
func (l *Italian) Homoglyphs() map[string][]string { return l.homoglyphs }
func (l *Italian) SimilarChars(char string) []string {
	return languages.SimilarChars(l.homoglyphs, char)
}
func (l *Italian) SimilarSpellings(word string) []string {
	return languages.SimilarSpellings(l.misspellings, word)
}
func (l *Italian) SimilarSounds(word string) []string {
	return languages.SimilarSounds(l.homophones, word)
}
func (l *Italian) Keyboards() (boards []internal.Keyboard) {
	for _, b := range languages.Keyboards() {
		if b.Language() == l.code {
			boards = append(boards, b)
		}
	}
	return
}

var (
	itMisspellings = [][]string{
		{"perche", "perché"},
		{"cosi", "così"},
		{"piu", "più"},
	}
	itHomophones = [][]string{
		{"punto", "."},
		{"virgola", ","},
		{"chiocciola", "@"},
		{"trattino", "-"},
		{"slash", "/"},
		{"anno", "hanno"},
	}
	itAntonyms = map[string][]string{
		"buono":     {"cattivo"},
		"cattivo":   {"buono"},
		"grande":    {"piccolo"},
		"piccolo":   {"grande"},
		"alto":      {"basso"},
		"basso":     {"alto"},
		"si":        {"no"},
		"no":        {"si"},
		"caldo":     {"freddo"},
		"freddo":    {"caldo"},
		"nuovo":     {"vecchio"},
		"vecchio":   {"nuovo"},
		"veloce":    {"lento"},
		"lento":     {"veloce"},
		"facile":    {"difficile"},
		"difficile": {"facile"},
		"chiaro":    {"scuro"},
		"scuro":     {"chiaro"},
		"aperto":    {"chiuso"},
		"chiuso":    {"aperto"},
		"dentro":    {"fuori"},
		"fuori":     {"dentro"},
		"sopra":     {"sotto"},
		"sotto":     {"sopra"},
		"prima":     {"dopo"},
		"dopo":      {"prima"},
		"presto":    {"tardi"},
		"tardi":     {"presto"},
		"iniziare":  {"finire"},
		"finire":    {"iniziare"},
		"vero":      {"falso"},
		"falso":     {"vero"},
		"ricco":     {"povero"},
		"povero":    {"ricco"},
		"forte":     {"debole"},
		"debole":    {"forte"},
		"pieno":     {"vuoto"},
		"vuoto":     {"pieno"},
		"vicino":    {"lontano"},
		"lontano":   {"vicino"},
		"comprare":  {"vendere"},
		"vendere":   {"comprare"},
		"dare":      {"prendere"},
		"prendere":  {"dare"},
		"amare":     {"odiare"},
		"odiare":    {"amare"},
		"vincere":   {"perdere"},
		"perdere":   {"vincere"},
		"giorno":    {"notte"},
		"notte":     {"giorno"},
	}
	Language = Italian{
		code:        LANGUAGE,
		name:        "Italian",
		description: "Italian is a Romance language spoken primarily in Italy and parts of Switzerland",
		numerals: map[string][]string{
			"0":          {"zero"},
			"1":          {"uno", "primo"},
			"2":          {"due", "secondo"},
			"3":          {"tre", "terzo"},
			"4":          {"quattro", "quarto"},
			"5":          {"cinque", "quinto"},
			"6":          {"sei", "sesto"},
			"7":          {"sette", "settimo"},
			"8":          {"otto", "ottavo"},
			"9":          {"nove", "nono"},
			"10":         {"dieci", "decimo"},
			"11":         {"undici", "undicesimo"},
			"12":         {"dodici", "dodicesimo"},
			"13":         {"tredici", "tredicesimo"},
			"14":         {"quattordici", "quattordicesimo"},
			"15":         {"quindici", "quindicesimo"},
			"16":         {"sedici", "sedicesimo"},
			"17":         {"diciassette", "diciassettesimo"},
			"18":         {"diciotto", "diciottesimo"},
			"19":         {"diciannove", "diciannovesimo"},
			"20":         {"venti", "ventesimo"},
			"30":         {"trenta", "trentesimo"},
			"40":         {"quaranta", "quarantesimo"},
			"50":         {"cinquanta", "cinquantesimo"},
			"60":         {"sessanta", "sessantesimo"},
			"70":         {"settanta", "settantesimo"},
			"80":         {"ottanta", "ottantesimo"},
			"90":         {"novanta", "novantesimo"},
			"100":        {"cento", "centesimo"},
			"1000":       {"mille", "millesimo"},
			"1000000":    {"milione", "milionesimo"},
			"1000000000": {"miliardo", "miliardesimo"},
		},
		graphemes: []string{
			"a", "b", "c", "d", "e", "f", "g",
			"h", "i", "j", "k", "l", "m", "n",
			"o", "p", "q", "r", "s", "t", "u",
			"v", "w", "x", "y", "z", "à", "è", "é", "ì", "ò", "ù",
		},
		vowels:       []string{"a", "e", "i", "o", "u"},
		misspellings: itMisspellings,
		homophones:   itHomophones,
		antonyms:     itAntonyms,
		homoglyphs:   languages.DefaultLatinHomoglyphs(),
	}
)

func init() { languages.AddLanguage(LANGUAGE, func() internal.Language { return &Language }) }
