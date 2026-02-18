// Copyright 2024 Rangertaha. All Rights Reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
package russian

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
)

const LANGUAGE string = "ru"

type Russian struct {
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

func (l *Russian) Id() string {
	return l.code
}
func (l *Russian) Name() string {
	return l.name
}
func (l *Russian) Description() string {
	return l.description
}
func (l *Russian) Numerals() map[string][]string {
	return l.numerals
}
func (l *Russian) Cardinal() map[string]string {
	return languages.NumeralMap(l.numerals, 0)
}

func (l *Russian) Ordinal() map[string]string {
	return languages.NumeralMap(l.numerals, 1)
}

func (l *Russian) Graphemes() []string {
	return l.graphemes
}

func (l *Russian) Vowels() []string {
	return l.vowels
}

func (l *Russian) Misspellings() [][]string {
	return l.misspellings
}

func (l *Russian) Homophones() [][]string {
	return l.homophones
}

func (l *Russian) Antonyms() map[string][]string {
	return l.antonyms
}

func (l *Russian) Homoglyphs() map[string][]string {
	return l.homoglyphs
}

func (l *Russian) SimilarChars(char string) []string {
	return languages.SimilarChars(l.homoglyphs, char)
}

func (l *Russian) SimilarSpellings(word string) []string {
	return languages.SimilarSpellings(l.misspellings, word)
}

func (l *Russian) SimilarSounds(word string) []string {
	return languages.SimilarSounds(l.homophones, word)
}

func (l *Russian) Keyboards() (boards []internal.Keyboard) {
	for _, b := range languages.Keyboards() {
		if b.Language() == l.code {
			boards = append(boards, b)
		}
	}
	return
}

var (
	// ruMisspellings are common misspellings
	ruMisspellings = [][]string{
		// Domain-friendly variants (no spaces). Includes common "ё"→"е" and "й"→"и" issues.
		{"еще", "ещё"},
		{"все", "всё"},
		{"ее", "её"},
		{"semyon", "semen"}, // translit-style brand confusion (optional but useful)

		// Common misspellings / confusions
		{"сдесь", "здесь"},
		{"зделать", "сделать"},
		{"вообщем", "вообще"},
		{"вобщем", "вообще"},
		{"пожалуста", "пожалуйста"},
		{"пожалуйсто", "пожалуйста"},
		{"потомучто", "потомучто", "потомучто"},
		{"черезчур", "чересчур"},
		{"учавствовать", "участвовать"},
		{"агенство", "агентство"},
		{"инжинер", "инженер"},
		{"оканчание", "окончание"},
		{"будующее", "будущее"},
		{"прийдти", "прийти"},
		{"придти", "прийти"}, // often confused; also treated as sound-alike
		{"тся", "ться"},      // very common confusion in Russian (suffix)
		{"извените", "извините"},
		{"воин", "воинн"}, // can help catch doubled letters
		{"счастье", "щастье"},
		{"ссора", "сора"},
		{"адресс", "адрес"},
		{"проблемма", "проблема"},
		{"комфортно", "комфорто"},
		{"информация", "информацыя"},
		{"организация", "организацыя"},
		{"компания", "компаниа"},
		{"профессор", "професор"},
		{"профессия", "професия"},
		{"интерестно", "интересно"},
	}

	// ruHomophones are words that sound alike
	ruHomophones = [][]string{
		{"точка", "."},
		{"собака", "@"},
		{"дефис", "-"},
		{"слэш", "/"},

		// Final devoicing / voiced-vs-voiceless pairs (useful for domain variants)
		{"лук", "луг"},
		{"кот", "код"},
		{"прут", "пруд"},
		{"рот", "род"},
		{"лет", "лед"},

		// Commonly confused in speech / writing
		{"придти", "прийти"},
	}

	// ruAntonyms are words opposite in meaning to another (e.g. bad and good ).
	ruAntonyms = map[string][]string{
		"хорошо":  {"плохо"},
		"плохо":   {"хорошо"},
		"хороший": {"плохой"},
		"плохой":  {"хороший"},
		"большой": {"маленький"},
		"маленький": {"большой"},
		"высокий": {"низкий"},
		"низкий":  {"высокий"},
		"сильный": {"слабый"},
		"слабый":  {"сильный"},
		"быстро":  {"медленно"},
		"медленно": {"быстро"},
		"быстрый": {"медленный"},
		"медленный": {"быстрый"},
		"легко":   {"трудно"},
		"трудно":  {"легко"},
		"легкий":  {"тяжелый"},
		"тяжелый": {"легкий"},
		"горячий": {"холодный"},
		"холодный": {"горячий"},
		"старый":  {"новый"},
		"новый":   {"старый"},
		"раньше":  {"позже"},
		"позже":   {"раньше"},
		"рано":    {"поздно"},
		"поздно":  {"рано"},
		"день":    {"ночь"},
		"ночь":    {"день"},
		"да":      {"нет"},
		"нет":     {"да"},
		"правда":  {"ложь"},
		"ложь":    {"правда"},
		"внутри":  {"снаружи"},
		"снаружи": {"внутри"},
		"открыть": {"закрыть"},
		"закрыть": {"открыть"},
		"вход":    {"выход"},
		"выход":   {"вход"},
		"приход":  {"уход"},
		"уход":    {"приход"},
		"начало":  {"конец"},
		"конец":   {"начало"},
		"жизнь":   {"смерть"},
		"смерть":  {"жизнь"},
		"свет":    {"тьма"},
		"тьма":    {"свет"},

		"белый":   {"черный"},
		"черный":  {"белый"},
		"левый":   {"правый"},
		"правый":  {"левый"},
		"вверх":   {"вниз"},
		"вниз":    {"вверх"},
		"вперед":  {"назад"},
		"назад":   {"вперед"},
		"много":   {"мало"},
		"мало":    {"много"},
		"длинный": {"короткий"},
		"короткий": {"длинный"},
		"полный":  {"пустой"},
		"пустой":  {"полный"},
		"купить":  {"продать"},
		"продать": {"купить"},
		"дать":    {"взять"},
		"взять":   {"дать"},
		"любить":  {"ненавидеть"},
		"ненавидеть": {"любить"},
		"победа":  {"поражение"},
		"поражение": {"победа"},
		"плюс":    {"минус"},
		"минус":   {"плюс"},
	}

	Language = Russian{
		code:        LANGUAGE,
		name:        "Russian",
		description: "Russian is the native language of the Russian people",

		// http://www.russianlessons.net/lessons/lesson2_main.php
		numerals: map[string][]string{
			// Number: cardinal..,  ordinal.., other...
			"0":          {"ноль"},
			"1":          {"один", "первый"},
			"2":          {"два", "второй"},
			"3":          {"три", "третий"},
			"4":          {"четыре", "четвертый"},
			"5":          {"пять", "пятый"},
			"6":          {"шесть", "шестой"},
			"7":          {"семь", "седьмой"},
			"8":          {"восемь", "восьмой"},
			"9":          {"девять", "девятый"},
			"10":         {"десять", "десятый"},
			"11":         {"одиннадцать", "одиннадцатый"},
			"12":         {"двенадцать", "двенадцатый"},
			"13":         {"тринадцать", "тринадцатый"},
			"14":         {"четырнадцать", "четырнадцатый"},
			"15":         {"пятнадцать", "пятнадцатый"},
			"16":         {"шестнадцать", "шестнадцатый"},
			"17":         {"семнадцать", "семнадцатый"},
			"18":         {"восемнадцать", "восемнадцатый"},
			"19":         {"девятнадцать", "девятнадцатый"},
			"20":         {"двадцать", "двадцатый"},
			"21":         []string{"двадцатьодин"},
			"22":         []string{"двадцатьдва"},
			"23":         []string{"двадцатьтри"},
			"24":         []string{"двадцатьчетыре"},
			"30":         {"тридцать", "тридцатый"},
			"40":         {"сорок", "сороковой"},
			"50":         {"пятьдесят", "пятидесятый"},
			"60":         {"шестьдесят", "шестидесятый"},
			"70":         {"семьдесят", "семидесятый"},
			"80":         {"восемьдесят", "восьмидесятый"},
			"90":         {"девяносто", "девяностый"},
			"100":        {"сто", "сотый"},
			"200":        {"двести"},
			"300":        {"триста"},
			"400":        {"четыреста"},
			"500":        {"пятьсот"},
			"600":        {"шестьсот"},
			"700":        {"семьсот"},
			"800":        {"восемьсот"},
			"900":        {"девятьсот"},
			"1000":       {"тысяча", "тысячный"},
			"1000000":    {"миллион", "миллионный"},
			"1000000000": {"миллиард", "миллиардный"},
		},
		graphemes: []string{
			"а", "б", "в", "г", "д", "е", "ё",
			"ж", "з", "и", "й", "к", "л", "м",
			"н", "о", "п", "р", "с", "т", "у",
			"ф", "х", "ц", "ч", "ш", "щ", "ъ",
			"ы", "ь", "э", "ю", "я", "ѕ", "ѯ",
			"ѱ", "ѡ", "ѫ", "ѧ", "ѭ", "ѩ"},
		vowels:       []string{"a", "о", "у", "э", "ы", "я", "ё", "ю", "е", "и"},
		misspellings: ruMisspellings,
		homophones:   ruHomophones,
		antonyms:     ruAntonyms,
		homoglyphs: map[string][]string{
			// Latin keys (useful when the target domain label is ASCII)
			"a": {"à", "á", "â", "ã", "ä", "å", "ɑ", "а", "ạ", "ǎ", "ă", "ȧ", "ӓ", "٨"},
			"b": {"d", "lb", "ib", "ʙ", "Ь", `b̔"`, "ɓ", "Б", "в"},
			"c": {"ϲ", "с", "ƈ", "ċ", "ć", "ç"},
			"e": {"é", "ê", "ë", "ē", "ĕ", "ě", "ė", "е", "ẹ", "ę", "є", "ϵ", "ҽ"},
			"h": {"һ", "հ", "Ꮒ", "н"},
			"i": {"1", "l", "Ꭵ", "í", "ï", "ı", "ɩ", "ι", "ꙇ", "ǐ", "ĭ", "¡", "і"},
			"k": {"κ", "ⲕ", "к"},
			"m": {"ṃ", "ᴍ", "м", "ɱ"},
			"o": {"0", "Ο", "ο", "О", "о", "Օ", "ȯ", "ọ", "ỏ", "ơ", "ó", "ö", "ӧ"},
			"p": {"ρ", "р", "ƿ", "Ϸ", "Þ"},
			"s": {"Ⴝ", "Ꮪ", "ʂ", "ś", "ѕ"},
			"t": {"τ", "т", "ţ"},
			"x": {"х", "ҳ", "ẋ"},
			"y": {"у", "Ү", "ý", "ʏ", "γ"},

			// Cyrillic keys (useful when the target domain label is Cyrillic)
			"а": {"a", "à", "á", "â", "ã", "ä", "å", "ɑ", "ạ", "ǎ", "ă", "ȧ", "ӓ"},
			"е": {"e", "é", "ê", "ë", "ē", "ĕ", "ě", "ė", "ẹ", "ę", "є", "ҽ", "ё"},
			"ё": {"е", "ë", "ė", "ӗ", "ӧ"},
			"о": {"o", "0", "Ο", "ο", "О", "о", "Օ", "ȯ", "ọ", "ỏ", "ơ", "ó", "ö", "ӧ"},
			"р": {"p", "ρ"},
			"с": {"c", "ϲ"},
			"х": {"x"},
			"у": {"y", "ү", "γ"},
			"к": {"k", "κ", "ⲕ"},
			"м": {"m", "ᴍ"},
			"т": {"t", "τ"},
			"в": {"b", "B"},
			"н": {"h", "H"},
			"и": {"u", "n"},
			"й": {"и", "u", "n"},
			"л": {"n", "Λ"},
			"ж": {"zh", "x", "*"},
			"з": {"3", "z"},
			"ф": {"φ", "Φ", "o"},
			"ц": {"u"},
			"ч": {"4"},
			"ш": {"w"},
			"щ": {"w", "ш"},
			"ъ": {"b", "ʙ"},
			"э": {"3", "€", "e"},
			"ю": {"io", "u"},
			"і": {"i", "1", "l"},
			"ј": {"j"},

			// Common Cyrillic lookalikes and confusables (subset)
			"б": {"6", "b", "Ь", `b̔"`, "ɓ", "Б"},
			"г": {"r", "Г", "ᴦ"},
			"д": {"d"},
			"п": {"n", "Π"},
			"ь": {"b", "Ь"},
			"ы": {"bl"},
			"я": {"r", "ɾ"},
			"ѕ": {"s", "5", "ƽ"},
			"ѯ": {"x", "ξ"},
			"ѱ": {"ψ"},
			"ѡ": {"w", "ω"},
			"ѫ": {"o"},
			"ѧ": {"я"},
			"ѭ": {"ю"},
			"ѩ": {"я"},
		},
	}
)

func init() {
	languages.AddLanguage(LANGUAGE, func() internal.Language {
		return &Language
	})
}
