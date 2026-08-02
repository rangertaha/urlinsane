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

// Command gen rebuilds the pkg/kb dataset from kbdlayout.info.
//
//	go generate ./pkg/kb
//
// The site lists layouts by KLID, the Windows keyboard layout identifier, but
// several KLIDs usually share one driver: KBDUS backs US English as well as
// the US variants of Chinese and Bulgarian. Requesting a KLID redirects to its
// driver, so the generator walks the KLID list, follows each redirect, and
// keeps one file per driver with the KLIDs it serves recorded as locales.
//
// Responses are cached on disk, so a re-run after a parsing change costs no
// requests. Delete the cache directory to fetch afresh.
package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode"

	"golang.org/x/net/html"

	"github.com/rangertaha/urlinsane/pkg/kb"
)

var (
	base    = flag.String("base", "https://kbdlayout.info", "site to fetch layouts from")
	out     = flag.String("out", "data", "directory to write the dataset to")
	cache   = flag.String("cache", ".cache", "directory to cache fetched responses in")
	delay   = flag.Duration("delay", 250*time.Millisecond, "pause between requests to the site")
	limit   = flag.Int("limit", 0, "stop after this many layouts, for a quick check (0 for all)")
	verbose = flag.Bool("v", false, "log every layout as it is written")
)

// ext is the extension of an encoded layout. The dataset is protocol buffers,
// so the files are opaque to a diff; the round-trip check in writeLayout and
// the tests in package kb are what stand in for reading them.
const ext = ".pb"

func main() {
	flag.Parse()
	log.SetFlags(0)
	log.SetPrefix("kb/gen: ")

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := os.MkdirAll(*out, 0o755); err != nil {
		return err
	}

	klids, err := catalogue()
	if err != nil {
		return fmt.Errorf("listing layouts: %w", err)
	}
	log.Printf("%d KLIDs listed", len(klids))

	var (
		entries []kb.Entry
		seen    = map[string]bool{}
		known   = map[string]bool{} // KLIDs already accounted for
		forms   = map[kb.Form]int{}
	)

	for _, klid := range klids {
		if known[klid] {
			continue
		}

		driver, meta, err := resolve(klid)
		if err != nil {
			return fmt.Errorf("resolving %s: %w", klid, err)
		}

		// A driver's page lists every KLID it backs, so the rest of them
		// need no request of their own.
		known[klid] = true
		for _, loc := range meta.Locales {
			known[loc.KLID] = true
		}

		if seen[driver] {
			continue
		}
		seen[driver] = true

		layout, err := build(driver, meta)
		if err != nil {
			return fmt.Errorf("building %s: %w", driver, err)
		}

		if err := writeLayout(filepath.Join(*out, driver+ext), layout); err != nil {
			return err
		}

		entries = append(entries, kb.Entry{
			ID:      layout.ID,
			Name:    layout.Name,
			Form:    layout.Form,
			Locales: layout.Locales,
		})
		forms[layout.Form]++

		if *verbose {
			log.Printf("%-12s %-34s %-4s %3d keys  %d locales",
				layout.ID, layout.Name, layout.Form, len(layout.Keys), len(layout.Locales))
		}
		if *limit > 0 && len(entries) >= *limit {
			break
		}
	}

	sort.Slice(entries, func(i, j int) bool { return entries[i].ID < entries[j].ID })

	// A limited run has only some of the layouts, so writing the index would
	// leave the dataset claiming the rest do not exist while their files sit
	// there — and the package believes the index. Leave it alone.
	if *limit > 0 {
		log.Printf("wrote %d layouts; index left untouched because -limit was set", len(entries))
		return nil
	}

	index, err := kb.MarshalCatalogue(entries)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(*out, "index"+ext), index, 0o644); err != nil {
		return err
	}

	// Drop layouts the site no longer publishes. This happens after the new
	// dataset is written rather than before, because the generator imports
	// the package it is generating for, and an empty data directory will
	// not build.
	stale, err := filepath.Glob(filepath.Join(*out, "*"+ext))
	if err != nil {
		return err
	}
	for _, path := range stale {
		id := strings.TrimSuffix(filepath.Base(path), ext)
		if id == "index" || seen[id] {
			continue
		}
		log.Printf("removing stale layout %s", id)
		if err := os.Remove(path); err != nil {
			return err
		}
	}

	log.Printf("wrote %d layouts (%d ansi, %d iso, %d jis)",
		len(entries), forms[kb.ANSI], forms[kb.ISO], forms[kb.JIS])

	return nil
}

// catalogue returns every KLID the site's front page links to.
func catalogue() ([]string, error) {
	doc, err := fetchHTML("/")
	if err != nil {
		return nil, err
	}

	var (
		klids []string
		seen  = map[string]bool{}
	)

	walk(doc, func(n *html.Node) {
		if n.Type != html.ElementNode || n.Data != "a" {
			return
		}
		id := strings.Trim(attr(n, "href"), "/")
		if len(id) != 8 || seen[id] {
			return
		}
		for _, r := range id {
			if !strings.ContainsRune("0123456789abcdefABCDEF", r) {
				return
			}
		}
		seen[id] = true
		klids = append(klids, strings.ToLower(id))
	})

	if len(klids) == 0 {
		return nil, fmt.Errorf("no layouts linked from the front page")
	}

	sort.Strings(klids)
	return klids, nil
}

// meta is what a layout page says about itself.
type meta struct {
	File    string
	Locales []kb.Locale
}

// resolve follows a KLID to the driver that serves it and reads the metadata
// table off the driver's page, which lists every KLID that driver backs.
func resolve(klid string) (driver string, m meta, err error) {
	doc, final, err := fetchHTMLFinal("/" + klid + "/")
	if err != nil {
		return "", meta{}, err
	}

	driver = strings.Trim(final.Path, "/")
	if driver == "" {
		return "", meta{}, fmt.Errorf("%s redirected to %s", klid, final)
	}
	driver = strings.ToLower(driver)

	// Each metadata table describes one KLID: a column of headings on the
	// left and values on the right.
	walk(doc, func(n *html.Node) {
		if n.Type != html.ElementNode || n.Data != "table" {
			return
		}

		fields := map[string]string{}
		walk(n, func(row *html.Node) {
			if row.Type != html.ElementNode || row.Data != "tr" {
				return
			}
			var cells []string
			for c := row.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && (c.Data == "th" || c.Data == "td") {
					cells = append(cells, strings.TrimSpace(text(c)))
				}
			}
			if len(cells) == 2 {
				fields[cells[0]] = cells[1]
			}
		})

		// "00000409 (en-US)" — the tag is optional.
		raw, ok := fields["KLID"]
		if !ok {
			return
		}
		id, tag, _ := strings.Cut(raw, " ")
		tag = strings.Trim(tag, "()")

		m.Locales = append(m.Locales, kb.Locale{
			KLID: strings.ToLower(strings.TrimSpace(id)),
			Tag:  strings.TrimSpace(tag),
			Name: fields["Layout Display Name"],
		})
		if m.File == "" {
			m.File = fields["Layout File"]
		}
	})

	if len(m.Locales) == 0 {
		// A driver with no metadata table still deserves an entry, keyed by
		// the KLID that led us here.
		m.Locales = []kb.Locale{{KLID: klid}}
	}

	return driver, m, nil
}

// The shape of a layout's XML export. Only the parts that bear on typed text
// are modelled; the site also records ligatures and full dead-key tables.
type (
	keyTable struct {
		Keys []physicalKey `xml:"PhysicalKeys>PK"`
	}

	physicalKey struct {
		VK      string   `xml:"VK,attr"`
		SC      string   `xml:"SC,attr"`
		Name    string   `xml:"Name,attr"`
		Results []result `xml:"Result"`
	}

	result struct {
		Text       string `xml:"Text,attr"`
		Codepoints string `xml:"TextCodepoints,attr"`
		With       string `xml:"With,attr"`
		Dead       *struct {
			Accent string `xml:"Accent,attr"`
		} `xml:"DeadKeyTable"`
	}
)

// build fetches a driver's key table and turns it into a layout.
func build(driver string, m meta) (*kb.Layout, error) {
	raw, err := fetch("/" + driver + "/download/xml")
	if err != nil {
		return nil, err
	}

	var doc keyTable
	if err := xml.Unmarshal(raw, &doc); err != nil {
		return nil, fmt.Errorf("parsing key table: %w", err)
	}

	keys := make([]kb.Key, 0, len(doc.Keys))
	for _, pk := range doc.Keys {
		key := kb.NewKey(pk.SC, pk.VK, pk.Name)

		for _, r := range pk.Results {
			mod, ok := modifiers(r.With)
			if !ok {
				continue
			}

			switch {
			case r.Dead != nil && r.Dead.Accent != "":
				key.Set(mod, kb.Out{Text: r.Dead.Accent, Dead: true})
			case r.Text != "":
				key.Set(mod, kb.Out{Text: r.Text})
			case r.Codepoints != "":
				if t, ok := codepoints(r.Codepoints); ok {
					key.Set(mod, kb.Out{Text: t})
				}
			}
		}

		key.Compact()
		keys = append(keys, key)
	}

	layout := kb.New(driver, name(driver, m), form(doc.Keys), keys)
	layout.File = m.File
	layout.Locales = m.Locales

	if len(layout.Keys) == 0 {
		return nil, fmt.Errorf("no keys in the alphanumeric block")
	}

	return layout, nil
}

// form decides which physical board a layout was drawn for. Windows drivers
// always define SC 56, the key an ISO board has and an ANSI board does not,
// but layouts designed for ANSI simply repeat the backslash already on SC 2B
// there. A layout that gives the two keys different characters is expecting
// the extra switch to exist.
func form(keys []physicalKey) kb.Form {
	unshifted := map[string]string{}
	for _, pk := range keys {
		for _, r := range pk.Results {
			if r.With == "" && r.Text != "" {
				unshifted[strings.ToUpper(pk.SC)] = r.Text
				break
			}
		}
	}

	extra, ok := unshifted["56"]
	if !ok || extra == "" || extra == unshifted["2B"] {
		return kb.ANSI
	}
	return kb.ISO
}

// name picks a display name for a driver. Where several locales share it, the
// shortest name is the one that names the layout rather than the locale:
// KBDUS is "US", not "Chinese (Traditional) - US".
func name(driver string, m meta) string {
	best := ""
	for _, loc := range m.Locales {
		if loc.Name == "" {
			continue
		}
		if best == "" || len(loc.Name) < len(best) || (len(loc.Name) == len(best) && loc.Name < best) {
			best = loc.Name
		}
	}
	if best == "" {
		best = strings.ToUpper(driver)
	}
	return best
}

// modifiers maps the site's modifier spelling onto a kb.Mod, and reports
// whether the combination produces text worth keeping. Ctrl on its own gives
// control codes, and the locking modifiers of East Asian layouts select a
// different script rather than a different character.
func modifiers(with string) (kb.Mod, bool) {
	if with == "" {
		return kb.Base, true
	}

	held := map[string]bool{}
	for _, f := range strings.Fields(with) {
		held[f] = true
	}

	var mod kb.Mod

	switch {
	case held["VK_CONTROL"] && held["VK_MENU"]:
		// The portable spelling of AltGr.
		mod |= kb.AltGr
		delete(held, "VK_CONTROL")
		delete(held, "VK_MENU")
	case held["VK_RMENU"]:
		mod |= kb.AltGr
		delete(held, "VK_RMENU")
	case held["VK_CONTROL"]:
		return 0, false
	}

	if held["VK_SHIFT"] {
		mod |= kb.Shift
		delete(held, "VK_SHIFT")
	}
	if held["VK_CAPITAL"] {
		mod |= kb.Caps
		delete(held, "VK_CAPITAL")
	}

	// Anything left is a modifier this package does not model.
	return mod, len(held) == 0
}

// codepoints decodes the hex form the site uses for characters it cannot put
// in an attribute. Control codes come through this way and are dropped.
func codepoints(s string) (string, bool) {
	var b strings.Builder

	for _, f := range strings.Fields(s) {
		raw, err := hex.DecodeString(strings.Repeat("0", len(f)%2) + f)
		if err != nil {
			return "", false
		}
		var r rune
		for _, c := range raw {
			r = r<<8 | rune(c)
		}
		if !unicode.IsGraphic(r) {
			return "", false
		}
		b.WriteRune(r)
	}

	return b.String(), b.Len() > 0
}

// write saves v as indented JSON with a trailing newline, so the dataset stays
// readable and diffs cleanly.
func write(path string, v any) error {
	raw, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(raw, '\n'), 0o644)
}

// writeLayout saves an encoded layout, after checking that it decodes back to
// the same thing. A change to how kb stores a layout therefore cannot quietly
// produce a dataset the package cannot read.
func writeLayout(path string, l *kb.Layout) error {
	raw, err := l.Marshal()
	if err != nil {
		return err
	}

	reloaded := new(kb.Layout)
	if err := reloaded.Unmarshal(raw); err != nil {
		return fmt.Errorf("%s does not decode: %w", path, err)
	}

	// Compare through JSON, which spells out everything a layout carries and
	// orders it, rather than comparing the encoded bytes to themselves.
	before, err := json.Marshal(l)
	if err != nil {
		return err
	}
	after, err := json.Marshal(reloaded)
	if err != nil {
		return err
	}
	if !bytes.Equal(before, after) {
		return fmt.Errorf("%s does not round-trip", path)
	}

	return os.WriteFile(path, raw, 0o644)
}

// fetch returns the body at a path, from the cache when possible.
func fetch(path string) ([]byte, error) {
	raw, _, err := fetchFinal(path)
	return raw, err
}

func fetchFinal(path string) ([]byte, *url.URL, error) {
	key := filepath.Join(*cache, strings.NewReplacer("/", "_", "?", "_", "+", "_").Replace(path)+".cache")

	if raw, err := os.ReadFile(key); err == nil {
		// The first line records where the request ended up.
		final, body, ok := strings.Cut(string(raw), "\n")
		if ok {
			if u, err := url.Parse(final); err == nil {
				return []byte(body), u, nil
			}
		}
	}

	time.Sleep(*delay)

	resp, err := http.Get(*base + path)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("GET %s: %s", path, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	final := resp.Request.URL
	if err := os.MkdirAll(*cache, 0o755); err == nil {
		_ = os.WriteFile(key, append([]byte(final.String()+"\n"), body...), 0o644)
	}

	return body, final, nil
}

func fetchHTML(path string) (*html.Node, error) {
	doc, _, err := fetchHTMLFinal(path)
	return doc, err
}

func fetchHTMLFinal(path string) (*html.Node, *url.URL, error) {
	raw, final, err := fetchFinal(path)
	if err != nil {
		return nil, nil, err
	}
	doc, err := html.Parse(strings.NewReader(string(raw)))
	if err != nil {
		return nil, nil, err
	}
	return doc, final, nil
}

// walk calls fn on n and every node beneath it.
func walk(n *html.Node, fn func(*html.Node)) {
	fn(n)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walk(c, fn)
	}
}

// attr returns an element's attribute value, or "".
func attr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

// text returns the concatenated text beneath a node.
func text(n *html.Node) string {
	var b strings.Builder
	walk(n, func(c *html.Node) {
		if c.Type == html.TextNode {
			b.WriteString(c.Data)
		}
	})
	return b.String()
}
