// Copyright 2026 Rangertaha. All Rights Reserved.
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

package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/rangertaha/urlinsane/internal/plugins/languages"
	_ "github.com/rangertaha/urlinsane/internal/plugins/languages/all"
	"github.com/urfave/cli/v2"
)

func defaultSynonyms(langID string) [][]string {
	// Keep this list intentionally small and word-level. It's mainly to ensure
	// every language has a non-empty synonym dataset for import/testing.
	switch langID {
	case "en":
		return [][]string{
			{"begin", "start", "commence", "initiate"},
		}
	case "fr":
		return [][]string{{"debut", "commencement", "depart"}}
	case "es":
		return [][]string{{"inicio", "comienzo", "arranque"}}
	case "pt":
		return [][]string{{"inicio", "comeco", "começo"}}
	case "it":
		return [][]string{{"inizio", "partenza", "avvio"}}
	case "de":
		return [][]string{{"anfang", "beginn", "start"}}
	case "nl":
		return [][]string{{"begin", "start"}}
	case "sv":
		return [][]string{{"start", "borja", "börja"}}
	case "no":
		return [][]string{{"start", "begynnelse", "begynne"}}
	case "da":
		return [][]string{{"start", "begyndelse", "begynde"}}
	case "fi":
		return [][]string{{"alku", "aloitus"}}
	case "ru":
		return [][]string{{"начало", "старт"}}
	case "uk":
		return [][]string{{"початок", "старт"}}
	case "pl":
		return [][]string{{"poczatek", "początek", "start"}}
	case "cs":
		return [][]string{{"zacatek", "začátek", "start"}}
	case "tr":
		return [][]string{{"baslangic", "başlangıç", "start"}}
	case "el":
		return [][]string{{"αρχη", "αρχή", "εναρξη", "έναρξη"}}
	case "ar":
		return [][]string{{"بداية", "بدء"}}
	case "fa":
		return [][]string{{"شروع", "آغاز"}}
	case "iw":
		return [][]string{{"התחלה", "תחילת"}}
	case "ps":
		return [][]string{{"پیل", "شروع"}}
	case "la":
		return [][]string{{"initium", "principium", "exordium"}}
	case "hy":
		return [][]string{{"սկիզբ", "մեկնարկ"}}
	case "ka":
		return [][]string{{"დაწყება", "დასაწყისი"}}
	case "hi":
		return [][]string{{"शुरुआत", "आरंभ"}}
	case "zh":
		return [][]string{{"开始", "开端"}}
	case "ja":
		return [][]string{{"開始", "始まり"}}
	case "ko":
		return [][]string{{"시작", "출발"}}
	case "th":
		return [][]string{{"เริ่ม", "เริ่มต้น"}}
	case "vi":
		return [][]string{{"batdau", "bắtđầu", "khoidau", "khởiđầu"}}
	default:
		// At least one safe line so the dataset isn't empty.
		return [][]string{{"start", "begin"}}
	}
}

func defaultPositive(langID string) []string {
	switch langID {
	case "en":
		return []string{"good", "great", "excellent", "safe", "trusted"}
	case "fr":
		return []string{"bon", "super", "excellent", "sûr", "fiable"}
	case "es":
		return []string{"bueno", "excelente", "seguro", "fiable", "genial"}
	case "pt":
		return []string{"bom", "excelente", "seguro", "confiável", "otimo", "ótimo"}
	case "de":
		return []string{"gut", "super", "sicher", "zuverlaessig", "zuverlässig", "toll"}
	case "it":
		return []string{"buono", "ottimo", "eccellente", "sicuro", "affidabile"}
	case "nl":
		return []string{"goed", "veilig", "betrouwbaar", "geweldig"}
	case "sv":
		return []string{"bra", "saker", "säker", "pålitlig"}
	case "no":
		return []string{"bra", "trygg", "sikker", "pålitelig"}
	case "da":
		return []string{"god", "sikker", "pålidelig"}
	case "fi":
		return []string{"hyvä", "turvallinen", "luotettava"}
	case "ru":
		return []string{"хорошо", "отлично", "безопасно", "надежно", "надёжно"}
	case "uk":
		return []string{"добре", "чудово", "безпечно", "надійно"}
	case "pl":
		return []string{"dobry", "świetny", "bezpieczny", "zaufany"}
	case "cs":
		return []string{"dobrý", "skvělý", "bezpečný", "spolehlivý"}
	case "tr":
		return []string{"iyi", "harika", "güvenli", "guvenli", "güvenilir", "guvenilir"}
	case "el":
		return []string{"καλό", "ασφαλές", "αξιόπιστο"}
	case "ar":
		return []string{"جيد", "ممتاز", "آمن", "موثوق"}
	case "fa":
		return []string{"خوب", "عالی", "امن", "مطمئن"}
	case "iw":
		return []string{"טוב", "מצוין", "בטוח", "אמין"}
	case "ps":
		return []string{"ښه", "عالي", "امن", "باوري", "ريښتینی"}
	case "la":
		return []string{"bonus", "optimus", "tutus", "fidelis", "verus"}
	case "hy":
		return []string{"լավ", "հիանալի", "անվտանգ", "վստահելի"}
	case "ka":
		return []string{"კარგი", "საუკეთესო", "უსაფრთხო"}
	case "hi":
		return []string{"अच्छा", "बेहतरीन", "सुरक्षित", "विश्वसनीय"}
	case "zh":
		return []string{"好", "优秀", "安全", "可靠"}
	case "ja":
		return []string{"良い", "優秀", "安全", "信頼できる"}
	case "ko":
		return []string{"좋다", "훌륭하다", "안전", "신뢰"}
	case "th":
		return []string{"ดี", "ยอดเยี่ยม", "ปลอดภัย", "เชื่อถือได้"}
	case "vi":
		return []string{"tốt", "tuyệtvời", "an toàn", "antoan", "đángtin"}
	default:
		return []string{"good"}
	}
}

func defaultNegative(langID string) []string {
	switch langID {
	case "en":
		return []string{"bad", "unsafe", "fake", "scam", "malicious"}
	case "fr":
		return []string{"mauvais", "dangereux", "faux", "arnaque", "malveillant"}
	case "es":
		return []string{"malo", "inseguro", "falso", "estafa", "malicioso"}
	case "pt":
		return []string{"mau", "inseguro", "falso", "golpe", "malicioso"}
	case "de":
		return []string{"schlecht", "unsicher", "falsch", "betrug", "boese", "böse"}
	case "it":
		return []string{"cattivo", "insicuro", "falso", "truffa", "malevolo"}
	case "nl":
		return []string{"slecht", "onveilig", "vals", "oplichting"}
	case "sv":
		return []string{"dålig", "osäker", "falsk", "bedrägeri"}
	case "no":
		return []string{"dårlig", "usikker", "falsk", "svindel"}
	case "da":
		return []string{"dårlig", "usikker", "falsk", "svindel"}
	case "fi":
		return []string{"huono", "turvaton", "väärennös", "huijaus"}
	case "ru":
		return []string{"плохо", "опасно", "фейк", "мошенничество", "вредоносный"}
	case "uk":
		return []string{"погано", "небезпечно", "фейк", "шахрайство", "шкідливий"}
	case "pl":
		return []string{"zły", "niebezpieczny", "fałszywy", "oszustwo"}
	case "cs":
		return []string{"špatný", "nebezpečný", "falešný", "podvod"}
	case "tr":
		return []string{"kötü", "guvensiz", "güvensiz", "sahte", "dolandırıcılık"}
	case "el":
		return []string{"κακό", "ανασφαλές", "ψεύτικο", "απάτη"}
	case "ar":
		return []string{"سيئ", "غيرآمن", "مزيف", "احتيال"}
	case "fa":
		return []string{"بد", "ناامن", "جعلی", "کلاهبرداری"}
	case "iw":
		return []string{"רע", "מסוכן", "מזויף", "הונאה"}
	case "ps":
		return []string{"بد", "ناامن", "جعلي", "درغلۍ", "خطرناک"}
	case "la":
		return []string{"malus", "periculosus", "falsus", "fraus", "nocivus"}
	case "hy":
		return []string{"վատ", "վտանգավոր", "կեղծ", "խաբեություն"}
	case "ka":
		return []string{"ცუდი", "სახიფათო", "ყალბი"}
	case "hi":
		return []string{"बुरा", "असुरक्षित", "नकली", "धोखा"}
	case "zh":
		return []string{"坏", "不安全", "假", "诈骗"}
	case "ja":
		return []string{"悪い", "危険", "偽", "詐欺"}
	case "ko":
		return []string{"나쁘다", "위험", "가짜", "사기"}
	case "th":
		return []string{"แย่", "ไม่ปลอดภัย", "ปลอม", "หลอกลวง"}
	case "vi":
		return []string{"xấu", "khôngan toàn", "khongan toan", "giả", "lừađảo"}
	default:
		return []string{"bad"}
	}
}

var SyncLanguagesCmd = cli.Command{
	Name:  "sync-languages",
	Usage: "Generate datasets/languages/<lang>/ files from registered language plugins",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "dir",
			Value: "datasets/languages",
			Usage: "output directory for language datasets",
		},
		&cli.BoolFlag{
			Name:  "overwrite",
			Value: false,
			Usage: "overwrite existing dataset files",
		},
	},
	Action: func(cCtx *cli.Context) error {
		base := cCtx.String("dir")
		overwrite := cCtx.Bool("overwrite")

		langs := languages.Languages()
		sort.Slice(langs, func(i, j int) bool { return langs[i].Id() < langs[j].Id() })

		for _, l := range langs {
			langDir := filepath.Join(base, l.Id())
			if err := os.MkdirAll(langDir, 0o755); err != nil {
				return err
			}

			// Write helpers
			write := func(filename string, content []byte) error {
				path := filepath.Join(langDir, filename)
				if !overwrite {
					if st, err := os.Stat(path); err == nil {
						// Do not overwrite curated datasets, but allow filling files that are empty (size=0).
						if st.Size() > 0 {
							return nil
						}
					}
				}
				return os.WriteFile(path, content, 0o644)
			}
			writeText := func(filename, content string) error {
				if content != "" && !strings.HasSuffix(content, "\n") {
					content += "\n"
				}
				return write(filename, []byte(content))
			}

			// Aggregate tokens that "belong" to the language plugin datasets to build a minimal word.lst.
			words := map[string]bool{}
			add := func(s string) {
				s = strings.TrimSpace(s)
				if s == "" {
					return
				}
				words[s] = true
			}

			// Numerals: "number token token ..."
			{
				var keys []string
				for k := range l.Numerals() {
					keys = append(keys, k)
				}
				sort.Slice(keys, func(i, j int) bool {
					ai, err1 := strconv.Atoi(keys[i])
					aj, err2 := strconv.Atoi(keys[j])
					if err1 == nil && err2 == nil {
						return ai < aj
					}
					return keys[i] < keys[j]
				})

				var b strings.Builder
				for _, k := range keys {
					b.WriteString(k)
					add(k)
					for _, tok := range l.Numerals()[k] {
						if strings.TrimSpace(tok) == "" {
							continue
						}
						b.WriteString(" ")
						b.WriteString(tok)
						add(tok)
					}
					b.WriteString("\n")
				}
				if err := writeText("numeral.lst", b.String()); err != nil {
					return err
				}
			}

			// Graphemes/vowels: one per line
			{
				var b strings.Builder
				for _, g := range l.Graphemes() {
					g = strings.TrimSpace(g)
					if g == "" {
						continue
					}
					b.WriteString(g)
					b.WriteString("\n")
					add(g)
				}
				if err := writeText("grapheme.lst", b.String()); err != nil {
					return err
				}
			}
			{
				var b strings.Builder
				for _, v := range l.Vowels() {
					v = strings.TrimSpace(v)
					if v == "" {
						continue
					}
					b.WriteString(v)
					b.WriteString("\n")
					add(v)
				}
				if err := writeText("vowel.lst", b.String()); err != nil {
					return err
				}
			}

			// Misspellings / homophones: one set per line
			{
				var b strings.Builder
				for _, set := range l.Misspellings() {
					var parts []string
					for _, tok := range set {
						tok = strings.TrimSpace(tok)
						if tok != "" {
							parts = append(parts, tok)
							add(tok)
						}
					}
					if len(parts) == 0 {
						continue
					}
					b.WriteString(strings.Join(parts, " "))
					b.WriteString("\n")
				}
				if err := writeText("misspelling.lst", b.String()); err != nil {
					return err
				}
			}
			{
				var b strings.Builder
				for _, set := range l.Homophones() {
					var parts []string
					for _, tok := range set {
						tok = strings.TrimSpace(tok)
						if tok != "" {
							parts = append(parts, tok)
							add(tok)
						}
					}
					if len(parts) == 0 {
						continue
					}
					b.WriteString(strings.Join(parts, " "))
					b.WriteString("\n")
				}
				if err := writeText("homophone.lst", b.String()); err != nil {
					return err
				}
			}

			// Antonyms: "word antonym1 antonym2 ..."
			{
				var keys []string
				for k := range l.Antonyms() {
					keys = append(keys, k)
				}
				sort.Strings(keys)

				var b strings.Builder
				for _, k := range keys {
					k = strings.TrimSpace(k)
					if k == "" {
						continue
					}
					add(k)
					b.WriteString(k)
					for _, a := range l.Antonyms()[k] {
						a = strings.TrimSpace(a)
						if a == "" {
							continue
						}
						add(a)
						b.WriteString(" ")
						b.WriteString(a)
					}
					b.WriteString("\n")
				}
				if err := writeText("antonym.lst", b.String()); err != nil {
					return err
				}
			}

			// Homoglyphs: "char homoglyph1 homoglyph2 ..."
			{
				var keys []string
				for k := range l.Homoglyphs() {
					keys = append(keys, k)
				}
				sort.Strings(keys)

				var buf bytes.Buffer
				for _, k := range keys {
					k = strings.TrimSpace(k)
					if k == "" {
						continue
					}
					add(k)
					buf.WriteString(k)
					for _, h := range l.Homoglyphs()[k] {
						h = strings.TrimSpace(h)
						if h == "" {
							continue
						}
						add(h)
						buf.WriteString(" ")
						buf.WriteString(h)
					}
					buf.WriteString("\n")
				}
				if err := write("homoglyph.lst", buf.Bytes()); err != nil {
					return err
				}
			}

			// Optional / currently unused by the importer, but keep directory parity with existing datasets.
			{
				var all []string
				for w := range words {
					all = append(all, w)
				}
				sort.Strings(all)
				var b strings.Builder
				for _, w := range all {
					b.WriteString(w)
					b.WriteString("\n")
				}
				if err := writeText("word.lst", b.String()); err != nil {
					return err
				}
			}
			if err := writeText("stopword.lst", ""); err != nil {
				return err
			}
			{
				// One word per line
				var b strings.Builder
				for _, w := range defaultPositive(l.Id()) {
					w = strings.TrimSpace(w)
					if w == "" {
						continue
					}
					add(w)
					b.WriteString(w)
					b.WriteString("\n")
				}
				if err := writeText("positive.lst", b.String()); err != nil {
					return err
				}
			}
			{
				// One word per line
				var b strings.Builder
				for _, w := range defaultNegative(l.Id()) {
					w = strings.TrimSpace(w)
					if w == "" {
						continue
					}
					add(w)
					b.WriteString(w)
					b.WriteString("\n")
				}
				if err := writeText("negative.lst", b.String()); err != nil {
					return err
				}
			}
			{
				var b strings.Builder
				for _, set := range defaultSynonyms(l.Id()) {
					var parts []string
					for _, tok := range set {
						tok = strings.TrimSpace(tok)
						if tok != "" {
							parts = append(parts, tok)
							add(tok)
						}
					}
					if len(parts) == 0 {
						continue
					}
					b.WriteString(strings.Join(parts, " "))
					b.WriteString("\n")
				}
				if err := writeText("synonym.lst", b.String()); err != nil {
					return err
				}
			}
			if err := writeText("token.lst", ""); err != nil {
				return err
			}
		}

		fmt.Printf("Synced %d language dataset folders into %s (overwrite=%v)\n", len(langs), base, overwrite)
		return nil
	},
}
