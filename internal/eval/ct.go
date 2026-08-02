// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package eval

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/rangertaha/urlinsane/internal/plugins/variant"
	"github.com/rangertaha/urlinsane/pkg/fuzzy"
)

// SourceCrtSh is the Source value written on records fetched from crt.sh.
const SourceCrtSh = "crt.sh"

// crtEntry is the subset of a crt.sh JSON record this package reads.
type crtEntry struct {
	NameValue string `json:"name_value"`
	NotBefore string `json:"not_before"`
}

// Filter decides which observed names count as lookalikes of a brand.
//
// Certificate Transparency answers "what names exist", not "what names are
// squatting". A substring query for a brand returns the brand's own
// certificates, its legitimate subdomains and vanity domains, and unrelated
// names that merely share a substring. Everything here is heuristic, which is
// why Record.Reviewed exists: the output is a candidate set for human review,
// and calling it ground truth before that review would overstate what it is.
type Filter struct {
	// Brand is the seed key the observations are lookalikes of.
	Brand string
	// MaxDistance is the largest Levenshtein distance between an observed
	// registrable name and the brand's that still counts as a lookalike. Names
	// that *contain* the brand core are kept regardless, since combosquats are
	// arbitrarily longer than the brand.
	MaxDistance int
	// Exclude lists registrable domains known to belong to the brand. Without
	// it a brand's own vanity domains are scored as squats it failed to
	// generate, which drags recall down for a reason that is not the
	// algorithms' fault.
	Exclude []string
}

// Keep reports whether an observed name is a lookalike worth recording.
func (f Filter) Keep(name string) bool {
	name = Normalize(name)
	if name == "" {
		return false
	}
	_, brandCore, brandSuffix := variant.SplitDomain(Normalize(f.Brand))
	_, core, suffix := variant.SplitDomain(name)
	if core == "" || brandCore == "" {
		return false
	}
	// The brand's own registrable domain is not a lookalike of itself. This is
	// deliberately the *registrable domain* and not the core: paypal.co shares
	// the core "paypal" with paypal.com and is exactly the TLD squat the tld
	// algorithm exists to find.
	if core == brandCore && suffix == brandSuffix {
		return false
	}
	for _, ex := range f.Exclude {
		if e := Normalize(ex); e != "" {
			if _, exCore, exSuffix := variant.SplitDomain(e); exCore == core && exSuffix == suffix {
				return false
			}
		}
	}
	if strings.Contains(core, brandCore) {
		return true
	}
	max := f.MaxDistance
	if max <= 0 {
		max = 2
	}
	return fuzzy.Levenshtein(core, brandCore) <= max
}

// FetchCrtSh queries crt.sh for names containing the brand's registrable label
// and returns those the filter keeps, deduped.
//
// crt.sh is frequently overloaded and answers 502/503 under load, so callers
// should treat a failure as "try later" rather than "no squats exist" — the
// difference matters, because an empty truth set scores as perfect recall.
func FetchCrtSh(ctx context.Context, client *http.Client, f Filter) ([]Record, error) {
	brand := Normalize(f.Brand)
	_, core, _ := variant.SplitDomain(brand)
	if core == "" {
		return nil, fmt.Errorf("eval: %q has no registrable name to search for", f.Brand)
	}
	if client == nil {
		client = &http.Client{Timeout: 120 * time.Second}
	}

	// %core% is a substring match: it finds paypal-login.com and mypaypal.net,
	// which a domain-scoped CT API cannot.
	endpoint := "https://crt.sh/?" + url.Values{
		"q":      {"%" + core + "%"},
		"output": {"json"},
	}.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("eval: crt.sh request: %w", err)
	}
	req.Header.Set("User-Agent", "urlinsane-eval/1 (+https://github.com/rangertaha/urlinsane)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("eval: crt.sh: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("eval: crt.sh returned %s (it is often overloaded; retry later)", resp.Status)
	}

	var entries []crtEntry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return nil, fmt.Errorf("eval: decoding crt.sh response: %w", err)
	}
	return recordsFrom(entries, f), nil
}

// recordsFrom applies the filter to decoded crt.sh entries. Split out from the
// fetch so it is testable without a network.
func recordsFrom(entries []crtEntry, f Filter) []Record {
	brand := Normalize(f.Brand)
	firstSeen := map[string]string{}
	for _, e := range entries {
		// name_value packs a certificate's SANs one per line.
		for _, raw := range strings.Split(e.NameValue, "\n") {
			name := Normalize(raw)
			if name == "" || !f.Keep(name) {
				continue
			}
			seen := e.NotBefore
			if len(seen) >= 10 {
				seen = seen[:10]
			}
			// Keep the earliest observation, which is the useful one: it dates
			// the registration rather than the most recent renewal.
			if prev, ok := firstSeen[name]; !ok || (seen != "" && seen < prev) {
				firstSeen[name] = seen
			}
		}
	}

	names := make([]string, 0, len(firstSeen))
	for name := range firstSeen {
		names = append(names, name)
	}
	sort.Strings(names)

	out := make([]Record, 0, len(names))
	for _, name := range names {
		out = append(out, Record{
			Brand:     brand,
			Name:      name,
			Type:      variant.TypeDomain,
			Source:    SourceCrtSh,
			FirstSeen: firstSeen[name],
			Reviewed:  false,
		})
	}
	return out
}
