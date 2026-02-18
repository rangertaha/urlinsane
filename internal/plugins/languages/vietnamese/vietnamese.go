package vietnamese

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
)

const LANGUAGE string = "vi"

type Vietnamese struct {
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

func (l *Vietnamese) Id() string                      { return l.code }
func (l *Vietnamese) Name() string                    { return l.name }
func (l *Vietnamese) Description() string             { return l.description }
func (l *Vietnamese) Numerals() map[string][]string   { return l.numerals }
func (l *Vietnamese) Cardinal() map[string]string     { return languages.NumeralMap(l.numerals, 0) }
func (l *Vietnamese) Ordinal() map[string]string      { return languages.NumeralMap(l.numerals, 1) }
func (l *Vietnamese) Graphemes() []string             { return l.graphemes }
func (l *Vietnamese) Vowels() []string                { return l.vowels }
func (l *Vietnamese) Misspellings() [][]string        { return l.misspellings }
func (l *Vietnamese) Homophones() [][]string          { return l.homophones }
func (l *Vietnamese) Antonyms() map[string][]string   { return l.antonyms }
func (l *Vietnamese) Homoglyphs() map[string][]string { return l.homoglyphs }
func (l *Vietnamese) SimilarChars(char string) []string {
	return languages.SimilarChars(l.homoglyphs, char)
}
func (l *Vietnamese) SimilarSpellings(word string) []string {
	return languages.SimilarSpellings(l.misspellings, word)
}
func (l *Vietnamese) SimilarSounds(word string) []string {
	return languages.SimilarSounds(l.homophones, word)
}
func (l *Vietnamese) Keyboards() (boards []internal.Keyboard) {
	for _, b := range languages.Keyboards() {
		if b.Language() == l.code {
			boards = append(boards, b)
		}
	}
	return
}

var (
	viHomoglyphs = func() map[string][]string {
		m := languages.DefaultLatinHomoglyphs()
		m["đ"] = []string{"d", "ð"}
		m["ă"] = []string{"a", "á", "à"}
		m["â"] = []string{"a", "á", "à", "ä"}
		m["ê"] = []string{"e", "é", "è", "ë"}
		m["ô"] = []string{"o", "ó", "ö", "0"}
		m["ơ"] = []string{"o", "ó", "ö", "0"}
		m["ư"] = []string{"u", "ú", "ü"}
		return m
	}()

	viMisspellings = [][]string{
		{"viet", "việt"},
		{"congnghe", "côngnghệ"},
	}
	viHomophones = [][]string{
		{"cham", ".", "chấm"},
	}
	viAntonyms = map[string][]string{
		"tốt":    {"xấu"},
		"xấu":    {"tốt"},
		"lớn":    {"nhỏ"},
		"nhỏ":    {"lớn"},
		"cao":    {"thấp"},
		"thấp":   {"cao"},
		"dài":    {"ngắn"},
		"ngắn":   {"dài"},
		"nhiều":  {"ít"},
		"ít":     {"nhiều"},
		"mới":    {"cũ"},
		"cũ":     {"mới"},
		"nhanh":  {"chậm"},
		"chậm":   {"nhanh"},
		"dễ":     {"khó"},
		"khó":    {"dễ"},
		"nóng":   {"lạnh"},
		"lạnh":   {"nóng"},
		"sáng":   {"tối"},
		"tối":    {"sáng"},
		"mạnh":   {"yếu"},
		"yếu":    {"mạnh"},
		"nặng":   {"nhẹ"},
		"nhẹ":    {"nặng"},
		"dày":    {"mỏng"},
		"mỏng":   {"dày"},
		"mở":     {"đóng"},
		"đóng":   {"mở"},
		"vào":    {"ra"},
		"ra":     {"vào"},
		"lên":    {"xuống"},
		"xuống":  {"lên"},
		"trước":  {"sau"},
		"sau":    {"trước"},
		"trái":   {"phải"},
		"phải":   {"trái"},
		"gần":    {"xa"},
		"xa":     {"gần"},
		"đúng":   {"sai"},
		"sai":    {"đúng"},
		"thật":   {"giả"},
		"giả":    {"thật"},
		"thắng":  {"thua"},
		"thua":   {"thắng"},
		"mua":    {"bán"},
		"bán":    {"mua"},
		"đến":    {"đi"},
		"đi":     {"đến"},
		"có":     {"không"},
		"không":  {"có"},
	}

	Language = Vietnamese{
		code:        LANGUAGE,
		name:        "Vietnamese",
		description: "Vietnamese is an Austroasiatic language written with the Latin-based Vietnamese alphabet",
		numerals: map[string][]string{
			"0":          {"không"},
			"1":          {"một", "thứnhất"},
			"2":          {"hai", "thứhai"},
			"3":          {"ba", "thứba"},
			"4":          {"bốn", "thứtư"},
			"5":          {"năm", "thứnăm"},
			"6":          {"sáu", "thứsáu"},
			"7":          {"bảy", "thứbảy"},
			"8":          {"tám", "thứtám"},
			"9":          {"chín", "thứchín"},
			"10":         {"mười", "thứmười"},
			"11":         {"mười một", "thứmườimột", "muoimot"},
			"12":         {"mười hai", "thứmườihai", "muoihai"},
			"13":         {"mười ba", "thứmườiba", "muoiba"},
			"14":         {"mười bốn", "thứmườibốn", "muoibon"},
			"15":         {"mười lăm", "thứmườilăm", "muoilam"},
			"16":         {"mười sáu", "thứmườisáu", "muoisau"},
			"17":         {"mười bảy", "thứmườibảy", "muoibay"},
			"18":         {"mười tám", "thứmườitám", "muoitam"},
			"19":         {"mười chín", "thứmườichín", "muoichin"},
			"20":         {"hai mươi", "haimuoi"},
			"30":         {"ba mươi", "bamuoi"},
			"40":         {"bốn mươi", "bonmuoi"},
			"50":         {"năm mươi", "nammuoi"},
			"60":         {"sáu mươi", "saumuoi"},
			"70":         {"bảy mươi", "baymuoi"},
			"80":         {"tám mươi", "tammuoi"},
			"90":         {"chín mươi", "chinmuoi"},
			"100":        {"trăm"},
			"1000":       {"nghìn", "nghin", "ngàn", "ngan"},
			"1000000":    {"triệu", "trieu"},
			"1000000000": {"tỷ", "ty"},
		},
		graphemes: []string{
			"a", "ă", "â", "b", "c", "d", "đ", "e", "ê", "g", "h", "i", "k", "l", "m", "n",
			"o", "ô", "ơ", "p", "q", "r", "s", "t", "u", "ư", "v", "x", "y",
		},
		vowels:       []string{"a", "ă", "â", "e", "ê", "i", "o", "ô", "ơ", "u", "ư", "y"},
		misspellings: viMisspellings,
		homophones:   viHomophones,
		antonyms:     viAntonyms,
		homoglyphs:   viHomoglyphs,
	}
)

func init() { languages.AddLanguage(LANGUAGE, func() internal.Language { return &Language }) }
