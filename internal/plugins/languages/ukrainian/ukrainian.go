package ukrainian

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
)

const LANGUAGE string = "uk"

type Ukrainian struct {
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

func (l *Ukrainian) Id() string                      { return l.code }
func (l *Ukrainian) Name() string                    { return l.name }
func (l *Ukrainian) Description() string             { return l.description }
func (l *Ukrainian) Numerals() map[string][]string   { return l.numerals }
func (l *Ukrainian) Cardinal() map[string]string     { return languages.NumeralMap(l.numerals, 0) }
func (l *Ukrainian) Ordinal() map[string]string      { return languages.NumeralMap(l.numerals, 1) }
func (l *Ukrainian) Graphemes() []string             { return l.graphemes }
func (l *Ukrainian) Vowels() []string                { return l.vowels }
func (l *Ukrainian) Misspellings() [][]string        { return l.misspellings }
func (l *Ukrainian) Homophones() [][]string          { return l.homophones }
func (l *Ukrainian) Antonyms() map[string][]string   { return l.antonyms }
func (l *Ukrainian) Homoglyphs() map[string][]string { return l.homoglyphs }
func (l *Ukrainian) SimilarChars(char string) []string {
	return languages.SimilarChars(l.homoglyphs, char)
}
func (l *Ukrainian) SimilarSpellings(word string) []string {
	return languages.SimilarSpellings(l.misspellings, word)
}
func (l *Ukrainian) SimilarSounds(word string) []string {
	return languages.SimilarSounds(l.homophones, word)
}
func (l *Ukrainian) Keyboards() (boards []internal.Keyboard) {
	for _, b := range languages.Keyboards() {
		if b.Language() == l.code {
			boards = append(boards, b)
		}
	}
	return
}

var (
	ukHomoglyphs = func() map[string][]string {
		m := languages.DefaultLatinHomoglyphs()
		// Cyrillic ↔ Latin confusables
		m["а"] = []string{"a", "à", "á", "â", "ã", "ä", "å", "ɑ"}
		m["е"] = []string{"e", "é", "è", "ê", "ë"}
		m["о"] = []string{"o", "0", "Ο", "О"}
		m["р"] = []string{"p", "ρ"}
		m["с"] = []string{"c", "ϲ"}
		m["х"] = []string{"x"}
		m["у"] = []string{"y"}
		m["і"] = []string{"i", "1", "l"}
		m["ї"] = []string{"ï", "i"}
		m["є"] = []string{"e", "є"}
		return m
	}()

	ukMisspellings = [][]string{
		{"kyiv", "kiev", "київ"},
	}
	ukHomophones = [][]string{
		{"крапка", "."},
	}
	ukAntonyms = map[string][]string{
		"добре":    {"погано"},
		"погано":   {"добре"},
		"великий":  {"малий"},
		"малий":    {"великий"},
		"так":      {"ні"},
		"ні":       {"так"},
		"високий":  {"низький"},
		"низький":  {"високий"},
		"гарячий":  {"холодний"},
		"холодний": {"гарячий"},
		"новий":    {"старий"},
		"старий":   {"новий"},
		"швидкий":  {"повільний"},
		"повільний": {"швидкий"},
		"легкий":   {"важкий"},
		"важкий":   {"легкий"},
		"світлий":  {"темний"},
		"темний":   {"світлий"},
		"відкритий": {"закритий"},
		"закритий":  {"відкритий"},
		"всередині": {"зовні"},
		"зовні":     {"всередині"},
		"вгорі":     {"внизу"},
		"внизу":     {"вгорі"},
		"до":        {"після"},
		"після":     {"до"},
		"рано":      {"пізно"},
		"пізно":     {"рано"},
		"початок":   {"кінець"},
		"кінець":    {"початок"},
		"правда":    {"брехня"},
		"брехня":    {"правда"},
		"багатий":   {"бідний"},
		"бідний":    {"багатий"},
		"сильний":   {"слабкий"},
		"слабкий":   {"сильний"},
		"повний":    {"порожній"},
		"порожній":  {"повний"},
		"близько":   {"далеко"},
		"далеко":    {"близько"},
		"купити":    {"продати"},
		"продати":   {"купити"},
		"дати":      {"взяти"},
		"взяти":     {"дати"},
		"любити":    {"ненавидіти"},
		"ненавидіти": {"любити"},
		"виграти":   {"програти"},
		"програти":  {"виграти"},
		"день":      {"ніч"},
		"ніч":       {"день"},
		"життя":     {"смерть"},
		"смерть":    {"життя"},
	}

	Language = Ukrainian{
		code:        LANGUAGE,
		name:        "Ukrainian",
		description: "Ukrainian is an East Slavic language and the official language of Ukraine",
		numerals: map[string][]string{
			"0":          {"нуль"},
			"1":          {"один", "перший"},
			"2":          {"два", "другий"},
			"3":          {"три", "третій"},
			"4":          {"чотири", "четвертий"},
			"5":          {"п’ять", "п’ятий"},
			"6":          {"шість", "шостий"},
			"7":          {"сім", "сьомий"},
			"8":          {"вісім", "восьмий"},
			"9":          {"дев’ять", "дев’ятий"},
			"10":         {"десять", "десятий"},
			"11":         {"одинадцять", "одинадцятий"},
			"12":         {"дванадцять", "дванадцятий"},
			"13":         {"тринадцять", "тринадцятий"},
			"14":         {"чотирнадцять", "чотирнадцятий"},
			"15":         {"п’ятнадцять", "п’ятнадцятий"},
			"16":         {"шістнадцять", "шістнадцятий"},
			"17":         {"сімнадцять", "сімнадцятий"},
			"18":         {"вісімнадцять", "вісімнадцятий"},
			"19":         {"дев’ятнадцять", "дев’ятнадцятий"},
			"20":         {"двадцять", "двадцятий"},
			"30":         {"тридцять", "тридцятий"},
			"40":         {"сорок", "сороковий"},
			"50":         {"п’ятдесят", "п’ятдесятий"},
			"60":         {"шістдесят", "шістдесятий"},
			"70":         {"сімдесят", "сімдесятий"},
			"80":         {"вісімдесят", "вісімдесятий"},
			"90":         {"дев’яносто", "дев’яностий"},
			"100":        {"сто", "сотий"},
			"1000":       {"тисяча", "тисячний"},
			"1000000":    {"мільйон", "мільйонний"},
			"1000000000": {"мільярд", "мільярдний"},
		},
		graphemes: []string{
			"а", "б", "в", "г", "ґ", "д", "е", "є", "ж", "з", "и", "і", "ї", "й",
			"к", "л", "м", "н", "о", "п", "р", "с", "т", "у", "ф", "х", "ц", "ч", "ш", "щ", "ь", "ю", "я",
		},
		vowels:       []string{"а", "е", "є", "и", "і", "ї", "о", "у", "ю", "я"},
		misspellings: ukMisspellings,
		homophones:   ukHomophones,
		antonyms:     ukAntonyms,
		homoglyphs:   ukHomoglyphs,
	}
)

func init() { languages.AddLanguage(LANGUAGE, func() internal.Language { return &Language }) }
