package turkish

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
)

const LANGUAGE string = "tr"

type Turkish struct {
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

func (l *Turkish) Id() string                      { return l.code }
func (l *Turkish) Name() string                    { return l.name }
func (l *Turkish) Description() string             { return l.description }
func (l *Turkish) Numerals() map[string][]string   { return l.numerals }
func (l *Turkish) Cardinal() map[string]string     { return languages.NumeralMap(l.numerals, 0) }
func (l *Turkish) Ordinal() map[string]string      { return languages.NumeralMap(l.numerals, 1) }
func (l *Turkish) Graphemes() []string             { return l.graphemes }
func (l *Turkish) Vowels() []string                { return l.vowels }
func (l *Turkish) Misspellings() [][]string        { return l.misspellings }
func (l *Turkish) Homophones() [][]string          { return l.homophones }
func (l *Turkish) Antonyms() map[string][]string   { return l.antonyms }
func (l *Turkish) Homoglyphs() map[string][]string { return l.homoglyphs }
func (l *Turkish) SimilarChars(char string) []string {
	return languages.SimilarChars(l.homoglyphs, char)
}
func (l *Turkish) SimilarSpellings(word string) []string {
	return languages.SimilarSpellings(l.misspellings, word)
}
func (l *Turkish) SimilarSounds(word string) []string {
	return languages.SimilarSounds(l.homophones, word)
}
func (l *Turkish) Keyboards() (boards []internal.Keyboard) {
	for _, b := range languages.Keyboards() {
		if b.Language() == l.code {
			boards = append(boards, b)
		}
	}
	return
}

var (
	trHomoglyphs = func() map[string][]string {
		m := languages.DefaultLatinHomoglyphs()
		m["ı"] = []string{"i", "l", "1"}
		m["İ"] = []string{"I", "l", "1"}
		m["ğ"] = []string{"g"}
		m["ş"] = []string{"s", "ś"}
		return m
	}()

	trMisspellings = [][]string{
		{"istanbul", "İstanbul"},
		{"turkiye", "türkiye", "Türkiye"},
		{"sirket", "şirket"},
	}
	trHomophones = [][]string{
		{"nokta", "."},
		{"et", "@"},
		{"tire", "-"},
		{"slash", "/"},
	}
	trAntonyms = map[string][]string{
		"iyi":        {"kötü"},
		"kötü":       {"iyi"},
		"büyük":      {"küçük"},
		"küçük":      {"büyük"},
		"evet":       {"hayır"},
		"hayır":      {"evet"},
		"yüksek":     {"alçak"},
		"alçak":      {"yüksek"},
		"sıcak":      {"soğuk"},
		"soğuk":      {"sıcak"},
		"yeni":       {"eski"},
		"eski":       {"yeni"},
		"hızlı":      {"yavaş"},
		"yavaş":      {"hızlı"},
		"kolay":      {"zor"},
		"zor":        {"kolay"},
		"aydınlık":   {"karanlık"},
		"karanlık":   {"aydınlık"},
		"açık":       {"kapalı"},
		"kapalı":     {"açık"},
		"iç":         {"dış"},
		"dış":        {"iç"},
		"yukarı":     {"aşağı"},
		"aşağı":      {"yukarı"},
		"önce":       {"sonra"},
		"sonra":      {"önce"},
		"erken":      {"geç"},
		"geç":        {"erken"},
		"başlangıç":  {"bitiş"},
		"bitiş":      {"başlangıç"},
		"doğru":      {"yanlış"},
		"yanlış":     {"doğru"},
		"gerçek":     {"sahte"},
		"sahte":      {"gerçek"},
		"zengin":     {"fakir"},
		"fakir":      {"zengin"},
		"güçlü":      {"zayıf"},
		"zayıf":      {"güçlü"},
		"dolu":       {"boş"},
		"boş":        {"dolu"},
		"yakın":      {"uzak"},
		"uzak":       {"yakın"},
		"almak":      {"satmak"},
		"satmak":     {"almak"},
		"vermek":     {"almak"},
		"sevmek":     {"nefretetmek"},
		"nefretetmek": {"sevmek"},
		"kazanmak":   {"kaybetmek"},
		"kaybetmek":  {"kazanmak"},
		"gün":        {"gece"},
		"gece":       {"gün"},
		"yaşam":      {"ölüm"},
		"ölüm":       {"yaşam"},
	}

	Language = Turkish{
		code:        LANGUAGE,
		name:        "Turkish",
		description: "Turkish is the most widely spoken of the Turkic languages",
		numerals: map[string][]string{
			"0":          {"sıfır"},
			"1":          {"bir", "birinci"},
			"2":          {"iki", "ikinci"},
			"3":          {"üç", "üçüncü"},
			"4":          {"dört", "dördüncü"},
			"5":          {"beş", "beşinci"},
			"6":          {"altı", "altıncı"},
			"7":          {"yedi", "yedinci"},
			"8":          {"sekiz", "sekizinci"},
			"9":          {"dokuz", "dokuzuncu"},
			"10":         {"on", "onuncu"},
			"11":         {"onbir", "onbirinci"},
			"12":         {"oniki", "onikinci"},
			"13":         {"onüç", "onüçüncü"},
			"14":         {"ondört", "ondördüncü"},
			"15":         {"onbeş", "onbeşinci"},
			"16":         {"onaltı", "onaltıncı"},
			"17":         {"onyedi", "onyedinci"},
			"18":         {"onsekiz", "onsekizinci"},
			"19":         {"ondokuz", "ondokuzuncu"},
			"20":         {"yirmi", "yirminci"},
			"30":         {"otuz", "otuzuncu"},
			"40":         {"kırk", "kırkıncı"},
			"50":         {"elli", "ellinci"},
			"60":         {"altmış", "altmışıncı"},
			"70":         {"yetmiş", "yetmişinci"},
			"80":         {"seksen", "sekseninci"},
			"90":         {"doksan", "doksanıncı"},
			"100":        {"yüz", "yüzüncü"},
			"1000":       {"bin", "bininci"},
			"1000000":    {"milyon", "milyonuncu"},
			"1000000000": {"milyar", "milyarıncı"},
		},
		graphemes: []string{
			"a", "b", "c", "ç", "d", "e", "f", "g", "ğ",
			"h", "ı", "i", "j", "k", "l", "m", "n", "o", "ö",
			"p", "r", "s", "ş", "t", "u", "ü", "v", "y", "z",
		},
		vowels:       []string{"a", "e", "ı", "i", "o", "ö", "u", "ü"},
		misspellings: trMisspellings,
		homophones:   trHomophones,
		antonyms:     trAntonyms,
		homoglyphs:   trHomoglyphs,
	}
)

func init() { languages.AddLanguage(LANGUAGE, func() internal.Language { return &Language }) }
