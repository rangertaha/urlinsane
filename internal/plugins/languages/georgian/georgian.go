package georgian

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
)

const LANGUAGE string = "ka"

type Georgian struct {
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

func (l *Georgian) Id() string                      { return l.code }
func (l *Georgian) Name() string                    { return l.name }
func (l *Georgian) Description() string             { return l.description }
func (l *Georgian) Numerals() map[string][]string   { return l.numerals }
func (l *Georgian) Cardinal() map[string]string     { return languages.NumeralMap(l.numerals, 0) }
func (l *Georgian) Ordinal() map[string]string      { return languages.NumeralMap(l.numerals, 1) }
func (l *Georgian) Graphemes() []string             { return l.graphemes }
func (l *Georgian) Vowels() []string                { return l.vowels }
func (l *Georgian) Misspellings() [][]string        { return l.misspellings }
func (l *Georgian) Homophones() [][]string          { return l.homophones }
func (l *Georgian) Antonyms() map[string][]string   { return l.antonyms }
func (l *Georgian) Homoglyphs() map[string][]string { return l.homoglyphs }
func (l *Georgian) SimilarChars(char string) []string {
	return languages.SimilarChars(l.homoglyphs, char)
}
func (l *Georgian) SimilarSpellings(word string) []string {
	return languages.SimilarSpellings(l.misspellings, word)
}
func (l *Georgian) SimilarSounds(word string) []string {
	return languages.SimilarSounds(l.homophones, word)
}
func (l *Georgian) Keyboards() (boards []internal.Keyboard) {
	for _, b := range languages.Keyboards() {
		if b.Language() == l.code {
			boards = append(boards, b)
		}
	}
	return
}

var (
	kaHomophones = [][]string{
		{"წერტილი", "."},
	}
	kaAntonyms = map[string][]string{
		"კარგი":      {"ცუდი"},
		"ცუდი":       {"კარგი"},
		"დიდი":       {"პატარა"},
		"პატარა":     {"დიდი"},
		"მაღალი":     {"დაბალი"},
		"დაბალი":     {"მაღალი"},
		"გრძელი":     {"მოკლე"},
		"მოკლე":      {"გრძელი"},
		"ბევრი":      {"ცოტა"},
		"ცოტა":       {"ბევრი"},
		"ახალი":      {"ძველი"},
		"ძველი":      {"ახალი"},
		"სწრაფი":     {"ნელი"},
		"ნელი":       {"სწრაფი"},
		"ადვილი":     {"ძნელი"},
		"ძნელი":      {"ადვილი"},
		"ცხელი":      {"ცივი"},
		"ცივი":       {"ცხელი"},
		"ნათელი":     {"ბნელი"},
		"ბნელი":      {"ნათელი"},
		"ძლიერი":     {"სუსტი"},
		"სუსტი":      {"ძლიერი"},
		"მძიმე":      {"მსუბუქი"},
		"მსუბუქი":    {"მძიმე"},
		"სქელი":      {"თხელი"},
		"თხელი":      {"სქელი"},
		"ღია":        {"დახურული"},
		"დახურული":   {"ღია"},
		"შესვლა":     {"გამოსვლა"},
		"გამოსვლა":   {"შესვლა"},
		"ზემოთ":      {"ქვემოთ"},
		"ქვემოთ":     {"ზემოთ"},
		"წინ":        {"უკან"},
		"უკან":       {"წინ"},
		"მარცხენა":   {"მარჯვენა"},
		"მარჯვენა":   {"მარცხენა"},
		"ახლოს":      {"შორს"},
		"შორს":       {"ახლოს"},
		"სწორი":      {"არასწორი"},
		"არასწორი":   {"სწორი"},
		"ჭეშმარიტი":  {"ცრუ"},
		"ცრუ":        {"ჭეშმარიტი"},
		"გამარჯვება": {"დამარცხება"},
		"დამარცხება": {"გამარჯვება"},
		"ყიდვა":      {"გაყიდვა"},
		"გაყიდვა":    {"ყიდვა"},
		"მოსვლა":     {"წასვლა"},
		"წასვლა":     {"მოსვლა"},
		"კი":         {"არა"},
		"არა":        {"კი"},
	}
	Language = Georgian{
		code:        LANGUAGE,
		name:        "Georgian",
		description: "Georgian is the official language of Georgia, written in the Georgian script",
		numerals: map[string][]string{
			"0":          {"ნული"},
			"1":          {"ერთი", "პირველი"},
			"2":          {"ორი", "მეორე"},
			"3":          {"სამი", "მესამე"},
			"4":          {"ოთხი", "მეოთხე"},
			"5":          {"ხუთი", "მეხუთე"},
			"6":          {"ექვსი", "მეექვსე"},
			"7":          {"შვიდი", "მეშვიდე"},
			"8":          {"რვა", "მერვე"},
			"9":          {"ცხრა", "მეცხრე"},
			"10":         {"ათი", "მეათე"},
			"11":         {"თერთმეტი"},
			"12":         {"თორმეტი"},
			"13":         {"ცამეტი"},
			"14":         {"თოთხმეტი"},
			"15":         {"თხუთმეტი"},
			"16":         {"თექვსმეტი"},
			"17":         {"ჩვიდმეტი"},
			"18":         {"თვრამეტი"},
			"19":         {"ცხრამეტი"},
			"20":         {"ოცი"},
			"30":         {"ოცდაათი"},
			"40":         {"ორმოცი"},
			"50":         {"ორმოცდაათი"},
			"60":         {"სამოცი"},
			"70":         {"სამოცდაათი"},
			"80":         {"ოთხმოცი"},
			"90":         {"ოთხმოცდაათი"},
			"100":        {"ასი"},
			"1000":       {"ათასი"},
			"1000000":    {"მილიონი"},
			"1000000000": {"მილიარდი"},
		},
		graphemes: []string{
			"ა", "ბ", "გ", "დ", "ე", "ვ", "ზ", "თ", "ი", "კ", "ლ", "მ", "ნ", "ო", "პ",
			"ჟ", "რ", "ს", "ტ", "უ", "ფ", "ქ", "ღ", "ყ", "შ", "ჩ", "ც", "ძ", "წ", "ჭ", "ხ", "ჯ", "ჰ",
		},
		vowels:       []string{"ა", "ე", "ი", "ო", "უ"},
		misspellings: [][]string{},
		homophones:   kaHomophones,
		antonyms:     kaAntonyms,
		homoglyphs:   languages.DefaultLatinHomoglyphs(),
	}
)

func init() { languages.AddLanguage(LANGUAGE, func() internal.Language { return &Language }) }
