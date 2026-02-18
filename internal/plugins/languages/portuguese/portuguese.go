package portuguese

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
)

const LANGUAGE string = "pt"

type Portuguese struct {
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

func (l *Portuguese) Id() string                      { return l.code }
func (l *Portuguese) Name() string                    { return l.name }
func (l *Portuguese) Description() string             { return l.description }
func (l *Portuguese) Numerals() map[string][]string   { return l.numerals }
func (l *Portuguese) Cardinal() map[string]string     { return languages.NumeralMap(l.numerals, 0) }
func (l *Portuguese) Ordinal() map[string]string      { return languages.NumeralMap(l.numerals, 1) }
func (l *Portuguese) Graphemes() []string             { return l.graphemes }
func (l *Portuguese) Vowels() []string                { return l.vowels }
func (l *Portuguese) Misspellings() [][]string        { return l.misspellings }
func (l *Portuguese) Homophones() [][]string          { return l.homophones }
func (l *Portuguese) Antonyms() map[string][]string   { return l.antonyms }
func (l *Portuguese) Homoglyphs() map[string][]string { return l.homoglyphs }
func (l *Portuguese) SimilarChars(char string) []string {
	return languages.SimilarChars(l.homoglyphs, char)
}
func (l *Portuguese) SimilarSpellings(word string) []string {
	return languages.SimilarSpellings(l.misspellings, word)
}
func (l *Portuguese) SimilarSounds(word string) []string {
	return languages.SimilarSounds(l.homophones, word)
}
func (l *Portuguese) Keyboards() (boards []internal.Keyboard) {
	for _, b := range languages.Keyboards() {
		if b.Language() == l.code {
			boards = append(boards, b)
		}
	}
	return
}

var (
	ptMisspellings = [][]string{
		{"informacao", "informação"},
		{"telefone", "telefone"},
		{"voce", "você"},
		{"coracao", "coração"},
	}
	ptHomophones = [][]string{
		{"ponto", "."},
		{"virgula", ","},
		{"vírgula", ","},
		{"arroba", "@"},
		{"hifen", "-", "hífen"},
		{"concerto", "conserto"},
		{"sessao", "seção", "cessão"},
	}
	ptAntonyms = map[string][]string{
		"bom":       {"mau"},
		"mau":       {"bom"},
		"grande":    {"pequeno"},
		"pequeno":   {"grande"},
		"alto":      {"baixo"},
		"baixo":     {"alto"},
		"sim":       {"não", "nao"},
		"não":       {"sim"},
		"nao":       {"sim"},
		"quente":    {"frio"},
		"frio":      {"quente"},
		"novo":      {"velho"},
		"velho":     {"novo"},
		"rápido":    {"lento"},
		"rapido":    {"lento"},
		"lento":     {"rápido", "rapido"},
		"fácil":     {"difícil"},
		"facil":     {"difícil"},
		"difícil":   {"fácil", "facil"},
		"claro":     {"escuro"},
		"escuro":    {"claro"},
		"aberto":    {"fechado"},
		"fechado":   {"aberto"},
		"dentro":    {"fora"},
		"fora":      {"dentro"},
		"cima":      {"baixo"},
		"antes":     {"depois"},
		"depois":    {"antes"},
		"cedo":      {"tarde"},
		"tarde":     {"cedo"},
		"começar":   {"terminar"},
		"comecar":   {"terminar"},
		"terminar":  {"começar", "comecar"},
		"verdadeiro": {"falso"},
		"falso":     {"verdadeiro"},
		"rico":      {"pobre"},
		"pobre":     {"rico"},
		"forte":     {"fraco"},
		"fraco":     {"forte"},
		"cheio":     {"vazio"},
		"vazio":     {"cheio"},
		"perto":     {"longe"},
		"longe":     {"perto"},
		"comprar":   {"vender"},
		"vender":    {"comprar"},
		"dar":       {"tomar"},
		"tomar":     {"dar"},
		"amar":      {"odiar"},
		"odiar":     {"amar"},
		"ganhar":    {"perder"},
		"perder":    {"ganhar"},
		"dia":       {"noite"},
		"noite":     {"dia"},
	}

	Language = Portuguese{
		code:        LANGUAGE,
		name:        "Portuguese",
		description: "Portuguese is a Romance language spoken in Europe, South America, Africa, and Asia",
		numerals: map[string][]string{
			"0":          {"zero"},
			"1":          {"um", "primeiro"},
			"2":          {"dois", "segundo"},
			"3":          {"três", "terceiro"},
			"4":          {"quatro", "quarto"},
			"5":          {"cinco", "quinto"},
			"6":          {"seis", "sexto"},
			"7":          {"sete", "sétimo"},
			"8":          {"oito", "oitavo"},
			"9":          {"nove", "nono"},
			"10":         {"dez", "décimo"},
			"11":         {"onze"},
			"12":         {"doze"},
			"13":         {"treze"},
			"14":         {"catorze", "quatorze"},
			"15":         {"quinze"},
			"16":         {"dezesseis", "dezasseis"},
			"17":         {"dezessete", "dezassete"},
			"18":         {"dezoito"},
			"19":         {"dezenove"},
			"20":         {"vinte"},
			"30":         {"trinta"},
			"40":         {"quarenta"},
			"50":         {"cinquenta"},
			"60":         {"sessenta"},
			"70":         {"setenta"},
			"80":         {"oitenta"},
			"90":         {"noventa"},
			"100":        {"cem"},
			"1000":       {"mil"},
			"1000000":    {"milhão", "milhao"},
			"1000000000": {"bilhão", "bilhao", "bilião", "biliao"},
		},
		graphemes: []string{
			"a", "b", "c", "d", "e", "f", "g",
			"h", "i", "j", "k", "l", "m", "n",
			"o", "p", "q", "r", "s", "t", "u",
			"v", "w", "x", "y", "z", "á", "à", "â", "ã", "ç", "é", "ê", "í", "ó", "ô", "õ", "ú", "ü",
		},
		vowels:       []string{"a", "e", "i", "o", "u", "á", "à", "â", "ã", "é", "ê", "í", "ó", "ô", "õ", "ú", "ü"},
		misspellings: ptMisspellings,
		homophones:   ptHomophones,
		antonyms:     ptAntonyms,
		homoglyphs:   languages.DefaultLatinHomoglyphs(),
	}
)

func init() {
	languages.AddLanguage(LANGUAGE, func() internal.Language { return &Language })
}
