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
package french

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
)

const LANGUAGE string = "fr"

type French struct {
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

func (l *French) Id() string {
	return l.code
}
func (l *French) Name() string {
	return l.name
}
func (l *French) Description() string {
	return l.description
}
func (l *French) Numerals() map[string][]string {
	return l.numerals
}
func (l *French) Cardinal() map[string]string {
	return languages.NumeralMap(l.numerals, 0)
}

func (l *French) Ordinal() map[string]string {
	return languages.NumeralMap(l.numerals, 1)
}

func (l *French) Graphemes() []string {
	return l.graphemes
}

func (l *French) Vowels() []string {
	return l.vowels
}

func (l *French) Misspellings() [][]string {
	return l.misspellings
}

func (l *French) Homophones() [][]string {
	return l.homophones
}

func (l *French) Antonyms() map[string][]string {
	return l.antonyms
}

func (l *French) Homoglyphs() map[string][]string {
	return l.homoglyphs
}

func (l *French) SimilarChars(char string) []string {
	return languages.SimilarChars(l.homoglyphs, char)
}

func (l *French) SimilarSpellings(word string) []string {
	return languages.SimilarSpellings(l.misspellings, word)
}

func (l *French) SimilarSounds(word string) []string {
	return languages.SimilarSounds(l.homophones, word)
}

func (l *French) Keyboards() (boards []internal.Keyboard) {
	for _, b := range languages.Keyboards() {
		if b.Language() == l.code {
			boards = append(boards, b)
		}
	}
	return
}

var (
	// frMisspellings are common misspellings
	frMisspellings = [][]string{
		// Accent omission / URL-friendly spellings
		{"francais", "français"},
		{"cafe", "café"},
		{"resume", "résumé"},
		{"eleve", "élève"},
		{"ecole", "école"},
		{"equipe", "équipe"},
		{"universite", "université"},
		{"societe", "société"},
		{"developpement", "développement"},
		{"tres", "très"},
		{"deja", "déjà"},
		{"noel", "noël"},
		{"hotel", "hôtel"},
		{"creme", "crème"},
		{"mere", "mère"},
		{"pere", "père"},
		{"surement", "sûrement"},

		// Common orthographic misspellings
		{"acceuil", "accueil"},
		{"acceuillir", "accueillir"},
		{"appeller", "appeler"},
		{"rappeller", "rappeler"},
		{"apparament", "apparemment"},
		{"professionel", "professionnel"},
		{"connection", "connexion"},
		{"addresse", "adresse"},
		{"occurence", "occurrence"},
		{"interressant", "intéressant", "interessant"},
		{"necessaire", "nécessaire"},
		{"definitivement", "définitivement"},
		{"independant", "indépendant"},
		{"evidement", "évidemment", "evidemment"},
		{"malgrés", "malgré", "malgre"},
		{"personel", "personnel"},
		{"organistation", "organisation"},
		{"securite", "sécurité", "securite"},
		{"competance", "compétence", "competence"},
		{"connaisance", "connaissance"},
		{"concurence", "concurrence"},
	}

	// frHomophones are words that sound alike
	frHomophones = [][]string{
		{"point", "."},

		// Classic French homophone sets (kept reasonably long to avoid exploding variants)
		{"ces", "ses"},
		{"cest", "sest"}, // "c'est" / "s'est" in URL-safe form
		{"son", "sont"},
		{"on", "ont"},
		{"peut", "peux"},
		{"ver", "vers", "vert", "verre"},
		{"mer", "mere", "mère", "maire"},
		{"compte", "conte", "comte"},
		{"seau", "saut", "sceau", "sot"},
		{"sain", "saint", "sein", "ceint"},
		{"sans", "sang", "cent", "cents"},
	}

	// frAntonyms are words opposite in meaning to another (e.g. bad and good ).
	frAntonyms = map[string][]string{
		"bien":     {"mal"},
		"mal":      {"bien"},
		"bon":      {"mauvais"},
		"mauvais":  {"bon"},
		"grand":    {"petit"},
		"petit":    {"grand"},
		"haut":     {"bas"},
		"bas":      {"haut"},
		"vrai":     {"faux"},
		"faux":     {"vrai"},
		"oui":      {"non"},
		"non":      {"oui"},
		"jour":     {"nuit"},
		"nuit":     {"jour"},
		"ouvrir":   {"fermer"},
		"fermer":   {"ouvrir"},
		"entrer":   {"sortir"},
		"sortir":   {"entrer"},
		"arriver":  {"partir"},
		"partir":   {"arriver"},
		"chaud":    {"froid"},
		"froid":    {"chaud"},
		"ancien":   {"nouveau"},
		"nouveau":  {"ancien"},
		"jeune":    {"vieux"},
		"vieux":    {"jeune"},
		"riche":    {"pauvre"},
		"pauvre":   {"riche"},
		"fort":     {"faible"},
		"faible":   {"fort"},
		"facile":   {"difficile"},
		"difficile": {"facile"},
		"possible": {"impossible"},
		"impossible": {"possible"},
		"heureux":  {"triste"},
		"triste":   {"heureux"},
		"toujours": {"jamais"},
		"jamais":   {"toujours"},
		"tot":      {"tard"},
		"tard":     {"tot"},
		"plein":    {"vide"},
		"vide":     {"plein"},
		"proche":   {"loin"},
		"loin":     {"proche"},
		"acheter":  {"vendre"},
		"vendre":   {"acheter"},
		"donner":   {"prendre"},
		"prendre":  {"donner"},
		"gagner":   {"perdre"},
		"perdre":   {"gagner"},
	}

	// French language
	Language = French{
		code:        LANGUAGE,
		name:        "French",
		description: "French is an official language in 27 countries",

		numerals: map[string][]string{
			// Number: cardinal..,  ordinal.., other...
			"0":          {"zero"},
			"1":          {"un", "premier"},
			"2":          {"deux", "deuxieme"},
			"3":          {"trois", "troisieme"},
			"4":          {"quatre", "quatrieme"},
			"5":          {"cinq", "cinquieme"},
			"6":          {"six", "sixieme"},
			"7":          {"sept", "septieme"},
			"8":          {"huit", "huitieme"},
			"9":          {"neuf", "neuvieme"},
			"10":         {"dix", "dixieme"},
			"11":         {"onze"},
			"12":         {"douze"},
			"13":         {"treize"},
			"14":         {"quatorze"},
			"15":         {"quinze"},
			"16":         {"seize"},
			"17":         {"dixsept"},
			"18":         {"dixhuit"},
			"19":         {"dixneuf"},
			"20":         {"vingt"},
			"21":         {"vingtetun"},
			"22":         {"vingtdeux"},
			"23":         {"vingttrois"},
			"24":         {"vingtquatre"},
			"25":         {"vingtcinq"},
			"26":         {"vingtsix"},
			"27":         {"vingtsept"},
			"28":         {"vingthuit"},
			"29":         {"vingtneuf"},
			"30":         {"trente"},
			"40":         {"quarante"},
			"50":         {"cinquante"},
			"60":         {"soixante"},
			"70":         {"soixantedix"},
			"80":         {"quatrevingt"},
			"90":         {"quatrevingtdix"},
			"100":        {"cent"},
			"1000":       {"mille"},
			"1000000":    {"million"},
			"1000000000": {"milliard"},
		},
		graphemes: []string{
			"a", "b", "c", "d", "e", "f", "g",
			"h", "i", "j", "k", "l", "m", "n",
			"o", "p", "q", "r", "s", "t", "u",
			"v", "w", "x", "y", "z", "ê", "û", "î", "ô", "â",
		},
		vowels:       []string{"a", "e", "i", "o", "u", "y"},
		misspellings: frMisspellings,
		homophones:   frHomophones,
		antonyms:     frAntonyms,
		homoglyphs: map[string][]string{
			// Latin + lookalikes from other scripts
			"a": {"à", "á", "â", "ã", "ä", "å", "ɑ", "а", "ạ", "ǎ", "ă", "ȧ", "ӓ", "٨"},
			"b": {"d", "lb", "ib", "ʙ", "Ь", `b̔"`, "ɓ", "Б"},
			"c": {"ϲ", "с", "ƈ", "ċ", "ć", "ç"},
			"d": {"b", "cl", "dl", "di", "ԁ", "ժ", "ɗ", "đ"},
			"e": {"é", "è", "ê", "ë", "ē", "ĕ", "ě", "ė", "е", "ẹ", "ę", "є", "ϵ", "ҽ"},
			"f": {"Ϝ", "ƒ", "Ғ"},
			"g": {"q", "ɢ", "ɡ", "Ԍ", "ġ", "ğ", "ց", "ǵ", "ģ"},
			"h": {"lh", "ih", "һ", "հ", "Ꮒ", "н"},
			"i": {"1", "l", "Ꭵ", "í", "ï", "ı", "ɩ", "ι", "ꙇ", "ǐ", "ĭ", "¡"},
			"j": {"ј", "ʝ", "ϳ", "ɉ"},
			"k": {"lk", "ik", "lc", "κ", "ⲕ"},
			"l": {"1", "i", "ɫ", "ł", "١", "ا"},
			"m": {"n", "nn", "rn", "rr", "ṃ", "ᴍ", "м", "ɱ"},
			"n": {"m", "r", "ń"},
			"o": {"0", "Ο", "ο", "О", "о", "Օ", "ȯ", "ọ", "ỏ", "ơ", "ó", "ö", "ӧ", "ه", "ة"},
			"p": {"ρ", "р", "ƿ", "Ϸ", "Þ"},
			"q": {"g", "զ", "ԛ", "գ", "ʠ"},
			"r": {"ʀ", "Г", "ᴦ", "ɼ", "ɽ"},
			"s": {"Ⴝ", "Ꮪ", "ʂ", "ś", "ѕ"},
			"t": {"τ", "т", "ţ"},
			"u": {"μ", "υ", "Ս", "ս", "ц", "ᴜ", "ǔ", "ŭ", "ù", "ú", "û", "ü"},
			"v": {"ѵ", "ν", "v̇"},
			"w": {"vv", "ѡ", "ա", "ԝ"},
			"x": {"х", "ҳ", "ẋ"},
			"y": {"ʏ", "γ", "у", "Ү", "ý", "ÿ"},
			"z": {"ʐ", "ż", "ź", "ᴢ"},

			// French-specific letters (ASCII approximations included)
			"à": {"a", "á", "â", "ä"},
			"â": {"a", "à", "á", "ä"},
			"ä": {"a", "à", "á", "â"},
			"ç": {"c", "ć", "č", "ҫ", "ϲ", "с"},
			"é": {"e", "è", "ê", "ë"},
			"è": {"e", "é", "ê", "ë"},
			"ê": {"e", "é", "è", "ë"},
			"ë": {"e", "é", "è", "ê"},
			"î": {"i", "í", "ï"},
			"ï": {"i", "í", "î"},
			"ô": {"o", "ó", "ö"},
			"ö": {"o", "ó", "ô"},
			"ù": {"u", "ú", "û", "ü"},
			"û": {"u", "ú", "ù", "ü"},
			"ü": {"u", "ú", "ù", "û"},
			"ÿ": {"y", "ý"},
			"œ": {"oe", "o", "e"},
			"æ": {"ae", "a", "e"},
		},
	}
)

func init() {
	languages.AddLanguage(LANGUAGE, func() internal.Language {
		return &Language
	})
}
