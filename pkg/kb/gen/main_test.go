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
package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"golang.org/x/net/html"
	"os"
	"path/filepath"
	"testing"

	"github.com/rangertaha/urlinsane/pkg/kb"
)

func TestCodepoints(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want string
		ok   bool
	}{
		{"0041", "A", true},
		{"41", "A", true},   // odd length, padded
		{"0020", " ", true}, // space is graphic
		{"1F600", "\U0001F600", true},
		{"0041 0301", "Á", true}, // a base letter and a combining mark
		{"001B", "", false},       // escape: a control code, not text
		{"0000", "", false},
		{"00410301", "", false}, // one number past the end of Unicode
		{"zz", "", false},
		{"", "", false},
	} {
		got, ok := codepoints(tc.in)
		if got != tc.want || ok != tc.ok {
			t.Errorf("codepoints(%q) = %q, %v; want %q, %v", tc.in, got, ok, tc.want, tc.ok)
		}
	}
}

func TestModifiers(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want kb.Mod
		ok   bool
	}{
		{"", kb.Base, true},
		{"VK_SHIFT", kb.Shift, true},
		{"VK_CAPITAL", kb.Caps, true},
		{"VK_SHIFT VK_CAPITAL", kb.Shift | kb.Caps, true},
		{"VK_CONTROL VK_MENU", kb.AltGr, true},
		{"VK_RMENU", kb.AltGr, true},
		{"VK_SHIFT VK_CONTROL VK_MENU", kb.Shift | kb.AltGr, true},
		{"VK_CONTROL VK_MENU VK_SHIFT VK_CAPITAL", kb.Shift | kb.AltGr | kb.Caps, true},

		// Ctrl and Alt alone give control codes, not text.
		{"VK_CONTROL", 0, false},
		{"VK_MENU", 0, false},
		{"VK_SHIFT VK_CONTROL", 0, false},

		// Levels this package has no bit for. Dropping them is deliberate,
		// and the run reports what went with them.
		{"VK_KANA", 0, false},
		{"VK_OEM_8", 0, false},
		{"VK_NUMLOCK", 0, false},
	} {
		got, ok := modifiers(tc.in)
		if ok != tc.ok || (ok && got != tc.want) {
			t.Errorf("modifiers(%q) = %v, %v; want %v, %v", tc.in, got, ok, tc.want, tc.ok)
		}
	}
}

func TestForm(t *testing.T) {
	key := func(sc, base string) physicalKey {
		return physicalKey{SC: sc, Results: []result{{Text: base}}}
	}

	for _, tc := range []struct {
		name string
		keys []physicalKey
		want kb.Form
	}{
		// Windows defines SC 56 on every layout. One drawn for ANSI just
		// repeats the backslash already on SC 2B there; one drawn for ISO
		// gives the extra switch a character of its own.
		{"repeats 2B", []physicalKey{key("56", `\`), key("2B", `\`)}, kb.ANSI},
		{"differs from 2B", []physicalKey{key("56", "<"), key("2B", "*")}, kb.ISO},
		{"no 56 at all", []physicalKey{key("2B", `\`)}, kb.ANSI},
		{"56 types nothing", []physicalKey{key("56", ""), key("2B", `\`)}, kb.ANSI},
		{"no keys", nil, kb.ANSI},
	} {
		if got := form(tc.keys); got != tc.want {
			t.Errorf("form(%s) = %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestName(t *testing.T) {
	// Where several locales share a driver, the shortest name is the one
	// naming the layout rather than the locale.
	got := name("kbdus", meta{Locales: []kb.Locale{
		{Name: "Chinese (Traditional) - US"},
		{Name: "US"},
		{Name: "Bulgarian (Latin)"},
	}})
	if got != "US" {
		t.Errorf("name = %q, want %q", got, "US")
	}

	// Ties break alphabetically, so the result does not depend on the order
	// the page happened to list them in.
	if got := name("kbdxx", meta{Locales: []kb.Locale{{Name: "BBB"}, {Name: "AAA"}}}); got != "AAA" {
		t.Errorf("name = %q, want AAA", got)
	}

	// With nothing to go on, fall back to the driver.
	if got := name("kbdxx", meta{Locales: []kb.Locale{{Name: ""}}}); got != "KBDXX" {
		t.Errorf("name = %q, want KBDXX", got)
	}
}

func TestResultText(t *testing.T) {
	dead := result{}
	dead.Dead = &struct {
		Accent string `xml:"Accent,attr"`
	}{Accent: "^"}

	for _, tc := range []struct {
		name string
		in   result
		want string
	}{
		{"plain", result{Text: "a"}, "a"},
		{"codepoints", result{Codepoints: "0041"}, "A"},
		{"control code", result{Codepoints: "001B"}, ""},
		{"dead key", dead, "^"},
		{"nothing", result{}, ""},
	} {
		if got := resultText(tc.in); got != tc.want {
			t.Errorf("resultText(%s) = %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestWriteAtomic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "entry.cache")

	if err := writeAtomic(path, []byte("first")); err != nil {
		t.Fatal(err)
	}
	if err := writeAtomic(path, []byte("second")); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "second" {
		t.Errorf("file holds %q, want the second write", got)
	}

	// The temporary file must not be left behind.
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Errorf("%d files in the cache directory, want just the entry", len(entries))
	}
}

// The parsing below runs on fixtures rather than on the site, which is what
// splitting it away from the fetching was for.

func parse(t *testing.T, s string) *html.Node {
	t.Helper()
	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		t.Fatal(err)
	}
	return doc
}

func TestParseKLIDs(t *testing.T) {
	doc := parse(t, `<html><body>
		<a href="/00000409/">US</a>
		<a href="/0000040C/">French</a>
		<a href="/00000409/">US again</a>
		<a href="/kbdus">not a KLID</a>
		<a href="/1234567/">too short</a>
		<a href="/00000zzz/">not hex</a>
		<a>no href</a>
	</body></html>`)

	got, err := parseKLIDs(doc)
	if err != nil {
		t.Fatal(err)
	}

	want := []string{"00000409", "0000040c"} // deduplicated, lowercased, sorted
	if len(got) != len(want) {
		t.Fatalf("parseKLIDs = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("parseKLIDs = %v, want %v", got, want)
			break
		}
	}
}

func TestParseKLIDsEmptyPage(t *testing.T) {
	// A front page with no layouts on it means something changed upstream,
	// and the run should stop rather than write an empty dataset.
	if _, err := parseKLIDs(parse(t, `<html><body>nothing here</body></html>`)); err == nil {
		t.Error("a page with no KLIDs should be an error")
	}
}

func TestParseMeta(t *testing.T) {
	doc := parse(t, `<html><body><div class="metaGroup">
		<table>
			<tr><th>KLID</th><td>00000409 (en-US)</td></tr>
			<tr><th>Layout Display Name</th><td>US</td></tr>
			<tr><th>Layout File</th><td>KBDUS.DLL</td></tr>
		</table>
		<table>
			<tr><th>KLID</th><td>00000804 (zh-Hans-CN)</td></tr>
			<tr><th>Layout Display Name</th><td>Chinese (Simplified) - US</td></tr>
			<tr><th>Layout File</th><td>KBDUS.DLL</td></tr>
		</table>
	</div></body></html>`)

	m := parseMeta(doc, "00000409")

	if m.File != "KBDUS.DLL" {
		t.Errorf("File = %q, want KBDUS.DLL", m.File)
	}
	if len(m.Locales) != 2 {
		t.Fatalf("got %d locales, want 2", len(m.Locales))
	}
	if m.Locales[0].KLID != "00000409" || m.Locales[0].Tag != "en-US" || m.Locales[0].Name != "US" {
		t.Errorf("first locale = %+v", m.Locales[0])
	}
	if m.Locales[1].KLID != "00000804" || m.Locales[1].Tag != "zh-Hans-CN" {
		t.Errorf("second locale = %+v", m.Locales[1])
	}
}

func TestParseMetaWithoutTables(t *testing.T) {
	// A page with no metadata still yields an entry, keyed by the KLID that
	// led there, so the layout is not lost.
	m := parseMeta(parse(t, `<html><body>nothing</body></html>`), "0000dead")

	if len(m.Locales) != 1 || m.Locales[0].KLID != "0000dead" {
		t.Errorf("locales = %+v, want the KLID we arrived with", m.Locales)
	}
}

func TestParseMetaWithoutATag(t *testing.T) {
	// The language tag is optional.
	m := parseMeta(parse(t, `<html><table>
		<tr><th>KLID</th><td>00000409</td></tr>
		<tr><th>Layout Display Name</th><td>US</td></tr>
	</table></html>`), "00000409")

	if len(m.Locales) != 1 || m.Locales[0].Tag != "" || m.Locales[0].KLID != "00000409" {
		t.Errorf("locales = %+v", m.Locales)
	}
}

const miniLayout = `<KeyboardLayout>
  <PhysicalKeys>
    <PK VK="VK_Q" SC="10"><Result Text="q"/><Result Text="Q" With="VK_SHIFT"/></PK>
    <PK VK="VK_E" SC="12">
      <Result Text="e"/>
      <Result Text="E" With="VK_SHIFT"/>
      <Result Text="E" With="VK_CAPITAL"/>
      <Result Text="e" With="VK_SHIFT VK_CAPITAL"/>
      <Result Text="€" With="VK_CONTROL VK_MENU"/>
      <Result TextCodepoints="0005" With="VK_CONTROL"/>
    </PK>
    <PK VK="VK_OEM_3" SC="29"><Result With="VK_CONTROL VK_MENU"><DeadKeyTable Accent="~"/></Result></PK>
    <PK VK="VK_SPACE" SC="39" Name="Space"><Result Text=" "/></PK>
    <PK VK="VK_NUMPAD7" SC="47"><Result Text="7"/></PK>
  </PhysicalKeys>
</KeyboardLayout>`

func TestBuildFrom(t *testing.T) {
	l, err := buildFrom("kbdtest", meta{
		File:    "KBDTEST.DLL",
		Locales: []kb.Locale{{KLID: "00000409", Tag: "en-US", Name: "Test"}},
	}, []byte(miniLayout))
	if err != nil {
		t.Fatal(err)
	}

	if l.ID != "kbdtest" || l.Name != "Test" || l.File != "KBDTEST.DLL" {
		t.Errorf("layout header = %+v", l)
	}

	// The numeric keypad key is outside the block and must be gone.
	if _, ok := l.Key("47"); ok {
		t.Error("SC 47 is on the keypad and should have been dropped")
	}

	e, ok := l.Key("12")
	if !ok {
		t.Fatal("no key at SC 12")
	}
	if e.Base() != "e" || e.Shift() != "E" || e.AltGr() != "€" {
		t.Errorf("SC 12 = base %q shift %q altgr %q", e.Base(), e.Shift(), e.AltGr())
	}

	// Caps here follows the ordinary rule, so Compact should have dropped
	// both stored caps states, leaving them to be derived.
	if got := e.Mods(); len(got) != 3 {
		t.Errorf("SC 12 stores %v; caps should have been compacted away", got)
	}
	if o, ok := e.Text(kb.Caps); !ok || o.Text != "E" {
		t.Errorf("caps derives to %q", o.Text)
	}

	// The Ctrl level produced a control code and must not have been kept.
	for _, m := range e.Mods() {
		if o, _ := e.Text(m); o.Text == "\x05" {
			t.Error("a control code was stored as text")
		}
	}

	// A dead key keeps its accent and its flag.
	d, _ := l.Key("29")
	if o, ok := d.Text(kb.AltGr); !ok || o.Text != "~" || !o.Dead {
		t.Errorf("SC 29 altgr = %q dead=%v", o.Text, o.Dead)
	}

	// The name on the space bar survives.
	if s, _ := l.Key("39"); s.Name != "Space" {
		t.Errorf("SC 39 name = %q", s.Name)
	}
}

func TestBuildFromRejectsRubbish(t *testing.T) {
	if _, err := buildFrom("kbdtest", meta{}, []byte("<not xml")); err == nil {
		t.Error("malformed XML should be an error")
	}

	// A key table with nothing in the alphanumeric block is not a layout.
	only := `<KeyboardLayout><PhysicalKeys>
		<PK VK="VK_NUMPAD7" SC="47"><Result Text="7"/></PK>
	</PhysicalKeys></KeyboardLayout>`
	if _, err := buildFrom("kbdtest", meta{}, []byte(only)); err == nil {
		t.Error("a layout with no keys in the block should be an error")
	}
}

func TestWriteLayout(t *testing.T) {
	l, err := buildFrom("kbdtest", meta{Locales: []kb.Locale{{KLID: "00000409"}}}, []byte(miniLayout))
	if err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(t.TempDir(), "kbdtest.pb")
	if err := writeLayout(path, l); err != nil {
		t.Fatal(err)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	back := new(kb.Layout)
	if err := back.Unmarshal(raw); err != nil {
		t.Fatal(err)
	}
	if back.ID != l.ID || len(back.Keys) != len(l.Keys) {
		t.Errorf("what was written back does not match what went in")
	}
}

func TestWrite(t *testing.T) {
	path := filepath.Join(t.TempDir(), "index.json")
	if err := write(path, []string{"a", "b"}); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != "[\n  \"a\",\n  \"b\"\n]\n" {
		t.Errorf("write produced %q", raw)
	}
}

func TestHTMLHelpers(t *testing.T) {
	doc := parse(t, `<html><body><div id="x" class="y">one<span>two</span></div></body></html>`)

	var div *html.Node
	walk(doc, func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			div = n
		}
	})
	if div == nil {
		t.Fatal("walk did not reach the div")
	}

	if got := attr(div, "id"); got != "x" {
		t.Errorf("attr(id) = %q", got)
	}
	if got := attr(div, "nope"); got != "" {
		t.Errorf("attr for a missing attribute = %q, want empty", got)
	}
	if got := text(div); got != "onetwo" {
		t.Errorf("text = %q, want onetwo", got)
	}
}

// The fetching below runs against a local server rather than kbdlayout.info,
// which covers the caching and redirect handling that the parsing tests skip.

func serve(t *testing.T) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			fmt.Fprint(w, `<html><a href="/00000409/">US</a><a href="/00000804/">zh</a></html>`)
		case "/00000409/", "/00000804/":
			// The site answers a KLID with a redirect to its driver.
			http.Redirect(w, r, "/kbdtest", http.StatusFound)
		case "/kbdtest":
			fmt.Fprint(w, `<html><table>
				<tr><th>KLID</th><td>00000409 (en-US)</td></tr>
				<tr><th>Layout Display Name</th><td>Test</td></tr>
				<tr><th>Layout File</th><td>KBDTEST.DLL</td></tr>
			</table><table>
				<tr><th>KLID</th><td>00000804 (zh-Hans-CN)</td></tr>
				<tr><th>Layout Display Name</th><td>Test Chinese</td></tr>
			</table></html>`)
		case "/kbdtest/download/xml":
			fmt.Fprint(w, miniLayout)
		default:
			http.NotFound(w, r)
		}
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	// Point the command at it, and take the politeness delay out.
	oldBase, oldCache, oldDelay := *base, *cache, *delay
	*base, *cache, *delay = srv.URL, t.TempDir(), 0
	t.Cleanup(func() { *base, *cache, *delay = oldBase, oldCache, oldDelay })

	return srv
}

func TestFetchAndCache(t *testing.T) {
	serve(t)

	raw, final, err := fetchFinal("/00000409/")
	if err != nil {
		t.Fatal(err)
	}
	// The redirect is followed, and where it ended up is what names the driver.
	if !strings.HasSuffix(final.Path, "/kbdtest") {
		t.Errorf("ended at %q, want /kbdtest", final.Path)
	}
	if !strings.Contains(string(raw), "KBDTEST.DLL") {
		t.Error("body is not the driver page")
	}

	// The second call is served from the cache, and agrees with the first.
	again, finalAgain, err := fetchFinal("/00000409/")
	if err != nil {
		t.Fatal(err)
	}
	if string(again) != string(raw) || finalAgain.Path != final.Path {
		t.Error("the cached answer differs from the fetched one")
	}

	if _, err := fetch("/nosuchpage"); err == nil {
		t.Error("a 404 should be an error")
	}
}

func TestCatalogueResolveBuild(t *testing.T) {
	serve(t)

	klids, err := catalogue()
	if err != nil {
		t.Fatal(err)
	}
	if len(klids) != 2 || klids[0] != "00000409" {
		t.Fatalf("catalogue = %v", klids)
	}

	driver, m, err := resolve(klids[0])
	if err != nil {
		t.Fatal(err)
	}
	if driver != "kbdtest" {
		t.Errorf("driver = %q, want kbdtest", driver)
	}
	if m.File != "KBDTEST.DLL" || len(m.Locales) != 2 {
		t.Errorf("meta = %+v", m)
	}

	l, err := build(driver, m)
	if err != nil {
		t.Fatal(err)
	}
	if l.ID != "kbdtest" || l.Name != "Test" {
		t.Errorf("layout = %s / %s", l.ID, l.Name)
	}
	if _, ok := l.Key("12"); !ok {
		t.Error("the built layout has no E key")
	}

	if _, err := build("nosuchdriver", meta{}); err == nil {
		t.Error("building an unknown driver should fail")
	}
}

func TestFetchHTMLRejectsAMissingPage(t *testing.T) {
	serve(t)

	if _, err := fetchHTML("/nosuchpage"); err == nil {
		t.Error("fetchHTML should pass the error through")
	}
}

func TestWriteAtomicIntoAMissingDirectory(t *testing.T) {
	// The temporary file is made beside the target, so a directory that is
	// not there is an error rather than a panic.
	if err := writeAtomic(filepath.Join(t.TempDir(), "nope", "x.cache"), []byte("x")); err == nil {
		t.Error("writing into a missing directory should fail")
	}
}
