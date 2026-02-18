package hindi

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
)

const LANGUAGE string = "hi"

type Hindi struct {
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

func (l *Hindi) Id() string                        { return l.code }
func (l *Hindi) Name() string                      { return l.name }
func (l *Hindi) Description() string               { return l.description }
func (l *Hindi) Numerals() map[string][]string     { return l.numerals }
func (l *Hindi) Cardinal() map[string]string       { return languages.NumeralMap(l.numerals, 0) }
func (l *Hindi) Ordinal() map[string]string        { return languages.NumeralMap(l.numerals, 1) }
func (l *Hindi) Graphemes() []string               { return l.graphemes }
func (l *Hindi) Vowels() []string                  { return l.vowels }
func (l *Hindi) Misspellings() [][]string          { return l.misspellings }
func (l *Hindi) Homophones() [][]string            { return l.homophones }
func (l *Hindi) Antonyms() map[string][]string     { return l.antonyms }
func (l *Hindi) Homoglyphs() map[string][]string   { return l.homoglyphs }
func (l *Hindi) SimilarChars(char string) []string { return languages.SimilarChars(l.homoglyphs, char) }
func (l *Hindi) SimilarSpellings(word string) []string {
	return languages.SimilarSpellings(l.misspellings, word)
}
func (l *Hindi) SimilarSounds(word string) []string {
	return languages.SimilarSounds(l.homophones, word)
}
func (l *Hindi) Keyboards() (boards []internal.Keyboard) {
	for _, b := range languages.Keyboards() {
		if b.Language() == l.code {
			boards = append(boards, b)
		}
	}
	return
}

var (
	hiMisspellings = [][]string{
		{"bharat", "भारत"},
	}
	hiHomophones = [][]string{
		{"बिंदु", "."},
		{"ऐट", "@", "at"},
		{"डैश", "-"},
		{"स्लैश", "/"},
	}
	hiAntonyms = map[string][]string{
		"अच्छा":    {"बुरा"},
		"बुरा":     {"अच्छा"},
		"बड़ा":     {"छोटा"},
		"छोटा":     {"बड़ा"},
		"ऊँचा":     {"नीचा"},
		"नीचा":     {"ऊँचा"},
		"लंबा":     {"नाटा"},
		"नाटा":     {"लंबा"},
		"ज्यादा":   {"कम"},
		"कम":       {"ज्यादा"},
		"नया":      {"पुराना"},
		"पुराना":   {"नया"},
		"तेज़":      {"धीमा"},
		"धीमा":     {"तेज़"},
		"आसान":     {"कठिन"},
		"कठिन":     {"आसान"},
		"गर्म":     {"ठंडा"},
		"ठंडा":     {"गर्म"},
		"उजला":     {"अँधेरा"},
		"अँधेरा":   {"उजला"},
		"मजबूत":    {"कमजोर"},
		"कमजोर":    {"मजबूत"},
		"भारी":     {"हल्का"},
		"हल्का":    {"भारी"},
		"मोटा":     {"पतला"},
		"पतला":     {"मोटा"},
		"खुला":     {"बंद"},
		"बंद":      {"खुला"},
		"अंदर":     {"बाहर"},
		"बाहर":     {"अंदर"},
		"ऊपर":      {"नीचे"},
		"नीचे":     {"ऊपर"},
		"आगे":      {"पीछे"},
		"पीछे":     {"आगे"},
		"बायाँ":    {"दायाँ"},
		"दायाँ":    {"बायाँ"},
		"पास":      {"दूर"},
		"दूर":      {"पास"},
		"सही":      {"गलत"},
		"गलत":      {"सही"},
		"सच":       {"झूठ"},
		"झूठ":      {"सच"},
		"जीत":      {"हार"},
		"हार":      {"जीत"},
		"खरीदना":   {"बेचना"},
		"बेचना":    {"खरीदना"},
		"आना":      {"जाना"},
		"जाना":     {"आना"},
		"है":       {"नहीं"},
		"नहीं":     {"है"},
		"हाँ":      {"नहीं"},
	}
	Language = Hindi{
		code:        LANGUAGE,
		name:        "Hindi",
		description: "Hindi is an Indo-Aryan language spoken primarily in India",
		numerals: map[string][]string{
			"0":    {"शून्य"},
			"1":    {"एक", "पहला"},
			"2":    {"दो", "दूसरा"},
			"3":    {"तीन", "तीसरा"},
			"4":    {"चार", "चौथा"},
			"5":    {"पाँच", "पाँचवाँ"},
			"6":    {"छह", "छठा"},
			"7":    {"सात", "सातवाँ"},
			"8":    {"आठ", "आठवाँ"},
			"9":    {"नौ", "नौवाँ"},
			"10":   {"दस", "दसवाँ"},
			"11":   {"ग्यारह"},
			"12":   {"बारह"},
			"13":   {"तेरह"},
			"14":   {"चौदह"},
			"15":   {"पंद्रह"},
			"16":   {"सोलह"},
			"17":   {"सत्रह"},
			"18":   {"अठारह"},
			"19":   {"उन्नीस"},
			"20":   {"बीस"},
			"30":   {"तीस"},
			"40":   {"चालीस"},
			"50":   {"पचास"},
			"60":   {"साठ"},
			"70":   {"सत्तर"},
			"80":   {"अस्सी"},
			"90":   {"नब्बे"},
			"100":  {"सौ"},
			"1000": {"हज़ार", "हजार"},
			// Indian numbering system (commonly seen)
			"100000":   {"लाख"},
			"10000000": {"करोड़"},
			// International-scale words (commonly used in tech/finance)
			"1000000":    {"मिलियन"},
			"1000000000": {"अरब", "बिलियन"},
		},
		graphemes: []string{
			"अ", "आ", "इ", "ई", "उ", "ऊ", "ए", "ऐ", "ओ", "औ",
			"क", "ख", "ग", "घ", "च", "छ", "ज", "झ", "ट", "ठ", "ड", "ढ", "त", "थ", "द", "ध", "न",
			"प", "फ", "ब", "भ", "म", "य", "र", "ल", "व", "श", "ष", "स", "ह",
		},
		vowels:       []string{"अ", "आ", "इ", "ई", "उ", "ऊ", "ए", "ऐ", "ओ", "औ"},
		misspellings: hiMisspellings,
		homophones:   hiHomophones,
		antonyms:     hiAntonyms,
		homoglyphs:   languages.DefaultLatinHomoglyphs(),
	}
)

func init() { languages.AddLanguage(LANGUAGE, func() internal.Language { return &Language }) }
