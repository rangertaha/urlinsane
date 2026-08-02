package korean

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
)

const LANGUAGE string = "ko"

type Korean struct {
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

func (l *Korean) Id() string                      { return l.code }
func (l *Korean) Name() string                    { return l.name }
func (l *Korean) Description() string             { return l.description }
func (l *Korean) Numerals() map[string][]string   { return l.numerals }
func (l *Korean) Cardinal() map[string]string     { return languages.NumeralMap(l.numerals, 0) }
func (l *Korean) Ordinal() map[string]string      { return languages.NumeralMap(l.numerals, 1) }
func (l *Korean) Graphemes() []string             { return l.graphemes }
func (l *Korean) Vowels() []string                { return l.vowels }
func (l *Korean) Misspellings() [][]string        { return l.misspellings }
func (l *Korean) Homophones() [][]string          { return l.homophones }
func (l *Korean) Antonyms() map[string][]string   { return l.antonyms }
func (l *Korean) Homoglyphs() map[string][]string { return l.homoglyphs }
func (l *Korean) SimilarChars(char string) []string {
	return languages.SimilarChars(l.homoglyphs, char)
}
func (l *Korean) SimilarSpellings(word string) []string {
	return languages.SimilarSpellings(l.misspellings, word)
}
func (l *Korean) SimilarSounds(word string) []string {
	return languages.SimilarSounds(l.homophones, word)
}
func (l *Korean) Keyboards() (boards []internal.Keyboard) {
	for _, b := range languages.Keyboards() {
		if b.Language() == l.code {
			boards = append(boards, b)
		}
	}
	return
}

var (
	koMisspellings = [][]string{
		{"되요", "돼요"},
		{"오랫만", "오랜만"},
		{"어떻해", "어떡해"},
		{"몇일", "며칠"},
		{"설레임", "설렘"},
		{"금새", "금세"},
		{"옛부터", "예부터"},
		{"무릎쓰다", "무릅쓰다"},
		{"삼가하다", "삼가다"},
		{"바램", "바람"},
		{"통채로", "통째로"},
		{"구지", "굳이"},
		{"역활", "역할"},
		{"희안하다", "희한하다"},
		{"사겨", "사귀어"},
		{"않되", "안돼"},
		{"어의없다", "어이없다"},
		{"문안하다", "무난하다"},
		{"일부로", "일부러"},
		{"할께요", "할게요"},
		{"갈께요", "갈게요"},
		{"낳았다", "나았다"},
		{"않하다", "안하다"},
		{"않되요", "안돼요"},
		{"설겆이", "설거지"},
		{"찌개", "찌게"},
		{"깨끗히", "깨끗이"},
	}

	koHomophones = [][]string{
		{"점", ".", "마침표"},
		{"골뱅이", "@"},
		{"대시", "-"},
		{"슬래시", "/"},

		{"반드시", "반듯이"},
		{"너머", "넘어"},
		{"느리다", "늘이다"},
		{"맞추다", "맞히다"},
		{"바라다", "바래다"},
		{"받치다", "받히다"},
		{"배다", "베다"},
		{"부치다", "붙이다"},
		{"비추다", "비치다"},
		{"빌다", "빌리다"},
		{"스러지다", "쓰러지다"},
		{"안치다", "앉히다"},
		{"어떻게", "어떡해"},
		{"왠지", "웬지"},
		{"이따가", "있다가"},
		{"저리다", "절이다"},
		{"조리다", "졸이다"},
		{"좇다", "쫓다"},
		{"지그시", "지긋이"},
		{"집다", "짚다"},
		{"햇빛", "햇볕", "햇살"},
		{"낫다", "낮다", "낳다"},
		{"걷히다", "거치다"},
		{"매다", "메다"},
		{"새다", "세다"},
		{"대다", "데다"},
		{"결재", "결제"},
		{"게시", "계시"},
		{"다치다", "닫히다"},
		{"맞히다", "마치다"},
		{"같이", "가치"},
		{"굳이", "구지"},
		{"시키다", "식히다"},
		{"업다", "없다"},
		{"짓다", "짖다"},
		{"바치다", "받치다"},
		{"받히다", "밭치다"},
		{"젖다", "젓다"},
		{"갔다", "같다"},
		{"섞다", "석다"},
		{"깎다", "깍다"},
		{"그러므로", "그럼으로"},
		{"햇빛", "햇볕"},
	}
	koAntonyms = map[string][]string{
		"좋다":   {"나쁘다"},
		"나쁘다":  {"좋다"},
		"크다":   {"작다"},
		"작다":   {"크다"},
		"높다":   {"낮다"},
		"낮다":   {"높다"},
		"길다":   {"짧다"},
		"짧다":   {"길다"},
		"많다":   {"적다"},
		"적다":   {"많다"},
		"새롭다":  {"낡다"},
		"낡다":   {"새롭다"},
		"빠르다":  {"느리다"},
		"느리다":  {"빠르다"},
		"쉽다":   {"어렵다"},
		"어렵다":  {"쉽다"},
		"뜨겁다":  {"차갑다"},
		"차갑다":  {"뜨겁다"},
		"밝다":   {"어둡다"},
		"어둡다":  {"밝다"},
		"강하다":  {"약하다"},
		"약하다":  {"강하다"},
		"무겁다":  {"가볍다"},
		"가볍다":  {"무겁다"},
		"두껍다":  {"얇다"},
		"얇다":   {"두껍다"},
		"굵다":   {"가늘다"},
		"가늘다":  {"굵다"},
		"열다":   {"닫다"},
		"닫다":   {"열다"},
		"들어가다": {"나가다"},
		"나가다":  {"들어가다"},
		"위":    {"아래"},
		"아래":   {"위"},
		"앞":    {"뒤"},
		"뒤":    {"앞"},
		"왼쪽":   {"오른쪽"},
		"오른쪽":  {"왼쪽"},
		"가깝다":  {"멀다"},
		"멀다":   {"가깝다"},
		"맞다":   {"틀리다"},
		"틀리다":  {"맞다"},
		"참":    {"거짓"},
		"거짓":   {"참"},
		"이기다":  {"지다"},
		"지다":   {"이기다"},
		"사다":   {"팔다"},
		"팔다":   {"사다"},
		"오다":   {"가다"},
		"가다":   {"오다"},
		"있다":   {"없다"},
		"없다":   {"있다"},
		"네":    {"아니요"},
		"아니요":  {"네"},
	}

	Language = Korean{
		code:        LANGUAGE,
		name:        "Korean",
		description: "Korean is a Koreanic language written mainly with Hangul",
		numerals: map[string][]string{
			"0":          {"영"},
			"1":          {"일", "첫째"},
			"2":          {"이", "둘째"},
			"3":          {"삼", "셋째"},
			"4":          {"사", "넷째"},
			"5":          {"오", "다섯째"},
			"6":          {"육", "여섯째"},
			"7":          {"칠", "일곱째"},
			"8":          {"팔", "여덟째"},
			"9":          {"구", "아홉째"},
			"10":         {"십", "열째"},
			"11":         {"십일"},
			"12":         {"십이"},
			"13":         {"십삼"},
			"14":         {"십사"},
			"15":         {"십오"},
			"16":         {"십육"},
			"17":         {"십칠"},
			"18":         {"십팔"},
			"19":         {"십구"},
			"20":         {"이십"},
			"30":         {"삼십"},
			"40":         {"사십"},
			"50":         {"오십"},
			"60":         {"육십"},
			"70":         {"칠십"},
			"80":         {"팔십"},
			"90":         {"구십"},
			"100":        {"백"},
			"1000":       {"천"},
			"10000":      {"만"},
			"1000000":    {"백만"},
			"1000000000": {"십억"},
		},
		graphemes:    []string{"ㄱ", "ㄴ", "ㄷ", "ㄹ", "ㅁ", "ㅂ", "ㅅ", "ㅇ", "ㅈ", "ㅊ", "ㅋ", "ㅌ", "ㅍ", "ㅎ", "ㄲ", "ㄸ", "ㅃ", "ㅆ", "ㅉ", "ㅏ", "ㅑ", "ㅓ", "ㅕ", "ㅗ", "ㅛ", "ㅜ", "ㅠ", "ㅡ", "ㅣ", "ㅐ", "ㅒ", "ㅔ", "ㅖ", "ㅘ", "ㅙ", "ㅚ", "ㅝ", "ㅞ", "ㅟ", "ㅢ"},
		vowels:       []string{"ㅏ", "ㅑ", "ㅓ", "ㅕ", "ㅗ", "ㅛ", "ㅜ", "ㅠ", "ㅡ", "ㅣ", "ㅐ", "ㅒ", "ㅔ", "ㅖ", "ㅘ", "ㅙ", "ㅚ", "ㅝ", "ㅞ", "ㅟ", "ㅢ"},
		misspellings: koMisspellings,
		homophones:   koHomophones,
		antonyms:     koAntonyms,
		homoglyphs:   languages.DefaultLatinHomoglyphs(),
	}
)

func init() { languages.AddLanguage(LANGUAGE, func() internal.Language { return &Language }) }
