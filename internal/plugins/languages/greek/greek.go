package greek

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
)

const LANGUAGE string = "el"

type Greek struct {
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

func (l *Greek) Id() string                        { return l.code }
func (l *Greek) Name() string                      { return l.name }
func (l *Greek) Description() string               { return l.description }
func (l *Greek) Numerals() map[string][]string     { return l.numerals }
func (l *Greek) Cardinal() map[string]string       { return languages.NumeralMap(l.numerals, 0) }
func (l *Greek) Ordinal() map[string]string        { return languages.NumeralMap(l.numerals, 1) }
func (l *Greek) Graphemes() []string               { return l.graphemes }
func (l *Greek) Vowels() []string                  { return l.vowels }
func (l *Greek) Misspellings() [][]string          { return l.misspellings }
func (l *Greek) Homophones() [][]string            { return l.homophones }
func (l *Greek) Antonyms() map[string][]string     { return l.antonyms }
func (l *Greek) Homoglyphs() map[string][]string   { return l.homoglyphs }
func (l *Greek) SimilarChars(char string) []string { return languages.SimilarChars(l.homoglyphs, char) }
func (l *Greek) SimilarSpellings(word string) []string {
	return languages.SimilarSpellings(l.misspellings, word)
}
func (l *Greek) SimilarSounds(word string) []string {
	return languages.SimilarSounds(l.homophones, word)
}
func (l *Greek) Keyboards() (boards []internal.Keyboard) {
	for _, b := range languages.Keyboards() {
		if b.Language() == l.code {
			boards = append(boards, b)
		}
	}
	return
}

var (
	elHomoglyphs = func() map[string][]string {
		m := languages.DefaultLatinHomoglyphs()
		// Greek ↔ Latin confusables (key IDN phishing surface)
		m["Α"] = []string{"A"}
		m["Β"] = []string{"B"}
		m["Ε"] = []string{"E"}
		m["Ζ"] = []string{"Z"}
		m["Η"] = []string{"H"}
		m["Ι"] = []string{"I", "1", "l"}
		m["Κ"] = []string{"K"}
		m["Μ"] = []string{"M"}
		m["Ν"] = []string{"N"}
		m["Ο"] = []string{"O", "0"}
		m["Ρ"] = []string{"P", "ρ"}
		m["Τ"] = []string{"T"}
		m["Υ"] = []string{"Y"}
		m["Χ"] = []string{"X"}
		m["α"] = []string{"a"}
		m["β"] = []string{"b"}
		m["ε"] = []string{"e"}
		m["ι"] = []string{"i", "1", "l"}
		m["κ"] = []string{"k"}
		m["ο"] = []string{"o", "0"}
		m["ρ"] = []string{"p", "ρ"}
		m["υ"] = []string{"u", "v", "y"}
		m["χ"] = []string{"x"}
		return m
	}()

	elMisspellings = [][]string{
		{"ellada", "ελλαδα", "ελλάδα"},
	}
	elHomophones = [][]string{
		{"τελεια", ".", "τελεία"},
	}
	elAntonyms = map[string][]string{
		"καλό":      {"κακό"},
		"κακό":      {"καλό"},
		"μεγάλο":    {"μικρό"},
		"μικρό":     {"μεγάλο"},
		"ναι":       {"όχι", "οχι"},
		"όχι":       {"ναι"},
		"οχι":       {"ναι"},
		"ψηλό":      {"χαμηλό"},
		"χαμηλό":    {"ψηλό"},
		"ζεστό":     {"κρύο"},
		"κρύο":      {"ζεστό"},
		"νέο":       {"παλιό"},
		"παλιό":     {"νέο"},
		"γρήγορο":   {"αργό"},
		"αργό":      {"γρήγορο"},
		"εύκολο":    {"δύσκολο"},
		"δύσκολο":   {"εύκολο"},
		"φωτεινό":   {"σκοτεινό"},
		"σκοτεινό":  {"φωτεινό"},
		"ανοιχτό":   {"κλειστό"},
		"κλειστό":   {"ανοιχτό"},
		"μέσα":      {"έξω"},
		"έξω":       {"μέσα"},
		"πάνω":      {"κάτω"},
		"κάτω":      {"πάνω"},
		"πριν":      {"μετά"},
		"μετά":      {"πριν"},
		"νωρίς":     {"αργά"},
		"αργά":      {"νωρίς"},
		"αρχή":      {"τέλος"},
		"τέλος":     {"αρχή"},
		"αλήθεια":   {"ψέμα"},
		"ψέμα":      {"αλήθεια"},
		"πλούσιος":  {"φτωχός"},
		"φτωχός":    {"πλούσιος"},
		"δυνατός":   {"αδύναμος"},
		"αδύναμος":  {"δυνατός"},
		"γεμάτο":    {"άδειο"},
		"άδειο":     {"γεμάτο"},
		"κοντά":     {"μακριά"},
		"μακριά":    {"κοντά"},
		"αγοράζω":   {"πουλάω"},
		"πουλάω":    {"αγοράζω"},
		"δίνω":      {"παίρνω"},
		"παίρνω":    {"δίνω"},
		"αγαπώ":     {"μισώ"},
		"μισώ":      {"αγαπώ"},
		"κερδίζω":   {"χάνω"},
		"χάνω":      {"κερδίζω"},
		"μέρα":      {"νύχτα"},
		"νύχτα":     {"μέρα"},
		"ζωή":       {"θάνατος"},
		"θάνατος":   {"ζωή"},
	}

	Language = Greek{
		code:        LANGUAGE,
		name:        "Greek",
		description: "Greek is a Hellenic language and the official language of Greece and Cyprus",
		numerals: map[string][]string{
			"0":          {"μηδέν"},
			"1":          {"ένα", "πρώτος"},
			"2":          {"δύο", "δεύτερος"},
			"3":          {"τρία", "τρίτος"},
			"4":          {"τέσσερα", "τέταρτος"},
			"5":          {"πέντε", "πέμπτος"},
			"6":          {"έξι", "έκτος"},
			"7":          {"επτά", "έβδομος"},
			"8":          {"οκτώ", "όγδοος"},
			"9":          {"εννέα", "ένατος"},
			"10":         {"δέκα", "δέκατος"},
			"11":         {"έντεκα"},
			"12":         {"δώδεκα"},
			"13":         {"δεκατρία"},
			"14":         {"δεκατέσσερα"},
			"15":         {"δεκαπέντε"},
			"16":         {"δεκαέξι"},
			"17":         {"δεκαεπτά"},
			"18":         {"δεκαοκτώ"},
			"19":         {"δεκαεννέα"},
			"20":         {"είκοσι"},
			"30":         {"τριάντα"},
			"40":         {"σαράντα"},
			"50":         {"πενήντα"},
			"60":         {"εξήντα"},
			"70":         {"εβδομήντα"},
			"80":         {"ογδόντα"},
			"90":         {"ενενήντα"},
			"100":        {"εκατό"},
			"1000":       {"χίλια"},
			"1000000":    {"εκατομμύριο"},
			"1000000000": {"δισεκατομμύριο"},
		},
		graphemes: []string{
			"α", "β", "γ", "δ", "ε", "ζ", "η", "θ", "ι", "κ", "λ", "μ", "ν", "ξ", "ο", "π", "ρ", "σ", "τ", "υ", "φ", "χ", "ψ", "ω",
		},
		vowels:       []string{"α", "ε", "η", "ι", "ο", "υ", "ω"},
		misspellings: elMisspellings,
		homophones:   elHomophones,
		antonyms:     elAntonyms,
		homoglyphs:   elHomoglyphs,
	}
)

func init() { languages.AddLanguage(LANGUAGE, func() internal.Language { return &Language }) }
