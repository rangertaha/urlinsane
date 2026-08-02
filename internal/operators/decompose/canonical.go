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

package decompose

import (
	"fmt"
	"net/netip"
	"strconv"
	"strings"

	"golang.org/x/net/idna"
)

// The canonicalization rules of docs/DESIGN.md §2. They run before the
// admission decision, not alongside it, so the truncation ledger's denylist
// compares canonical keys and "Example.com" cannot slip past a row recorded for
// "example.com". Returning an error refuses the candidate outright — which is
// the right answer for a key that names no entity, and the wrong one for a key
// that is merely unusual, so these lean permissive.

// hostProfile case-folds and punycodes a hostname the way a resolver lookup
// would, but deliberately does not enforce label syntax. Variant algorithms
// produce names a strict resolver would reject — a leading hyphen, a doubled
// dot, an underscore — and refusing those here would drop exactly the
// candidates a squatting scan exists to find, silently and before they are ever
// counted.
var hostProfile = idna.New(
	idna.MapForLookup(),
	idna.StrictDomainName(false),
	idna.ValidateLabels(false),
)

// canonHost is the shared hostname rule behind domain, tld and platform:
// trim, drop the root dot, case-fold, encode to punycode.
func canonHost(raw string) (string, error) {
	s := strings.TrimSuffix(strings.TrimSpace(raw), ".")
	if s == "" {
		return "", fmt.Errorf("empty host")
	}
	// A delimiter means the caller handed over a URL, an email or a path
	// without splitting it first. Silently keeping it would mint a host node
	// nothing else can ever converge on.
	if strings.ContainsAny(s, " \t\r\n/?#@:") {
		return "", fmt.Errorf("host %q contains a delimiter", raw)
	}
	s, err := hostProfile.ToASCII(s)
	if err != nil {
		return "", fmt.Errorf("host %q: %w", raw, err)
	}
	if s == "" {
		return "", fmt.Errorf("host %q canonicalized to nothing", raw)
	}
	return s, nil
}

// canonDomain: lowercase, IDNA to punycode, strip the trailing dot.
func canonDomain(raw string) (string, error) { return canonHost(raw) }

// canonPlatform: lowercase host. A platform is the site a repo or account lives
// on, keyed by host so github.com is one node however it was reached.
func canonPlatform(raw string) (string, error) { return canonHost(raw) }

// canonTLD: lowercase. The leading dot of ".com" is stripped so the dotted and
// undotted spellings converge.
func canonTLD(raw string) (string, error) {
	return canonHost(strings.TrimPrefix(strings.TrimSpace(raw), "."))
}

// canonEmail: lowercase the domain, preserve the local part's case. Local-part
// case is significant to the RFC even though nearly every provider ignores it,
// and folding it here would merge two addresses the mail system may treat as
// distinct.
func canonEmail(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	// Split on the last "@": a quoted local part may legally contain one.
	at := strings.LastIndex(s, "@")
	if at <= 0 || at == len(s)-1 {
		return "", fmt.Errorf("email %q is not local@domain", raw)
	}
	host, err := canonHost(s[at+1:])
	if err != nil {
		return "", err
	}
	return s[:at] + "@" + host, nil
}

// canonUsername: lowercase. §2 calls for "platform casing rules", but a
// username key carries no platform, so the only rule available is the one every
// major registry shares — case-insensitive matching. See the package tests for
// what this leaves unresolved.
func canonUsername(raw string) (string, error) {
	s := strings.ToLower(strings.TrimSpace(raw))
	if s == "" {
		return "", fmt.Errorf("empty username")
	}
	if strings.ContainsAny(s, " \t\r\n/@") {
		return "", fmt.Errorf("username %q contains a delimiter", raw)
	}
	return s, nil
}

// canonRepo: host, owner and name, lowercased. Every spelling of the same
// repository — browser URL, clone URL, scp-style git remote, bare path — has to
// land on one key, or a scan started from a clone URL and one started from a
// browser URL produce different graphs for the same repository.
func canonRepo(raw string) (string, error) {
	s := strings.TrimSuffix(strings.TrimSpace(raw), "/")
	scheme := false
	if i := strings.Index(s, "://"); i >= 0 {
		s, scheme = s[i+3:], true
	}
	if at := strings.Index(s, "@"); at >= 0 {
		// Either URL userinfo or an scp-style remote. Only the latter uses ":"
		// as the host/path separator; rewriting it under a scheme would mangle
		// a port instead.
		s = s[at+1:]
		if !scheme {
			s = strings.Replace(s, ":", "/", 1)
		}
	}
	s = strings.TrimSuffix(s, ".git")

	parts := strings.Split(s, "/")
	if len(parts) < 3 {
		return "", fmt.Errorf("repo %q is not host/owner/name", raw)
	}
	host, err := canonHost(parts[0])
	if err != nil {
		return "", err
	}
	owner, err := canonUsername(parts[1])
	if err != nil {
		return "", fmt.Errorf("repo %q: %w", raw, err)
	}
	name := strings.ToLower(strings.TrimSpace(parts[2]))
	if name == "" {
		return "", fmt.Errorf("repo %q has no name", raw)
	}
	// Anything past the third segment is a view of the repo — /tree/main,
	// /blob/... — not part of its identity.
	return host + "/" + owner + "/" + name, nil
}

// canonPackage: registry-qualified, then normalized per registry. The
// qualification is required rather than defaulted: npm's "lodash" and PyPI's
// "lodash" are different entities, and a bare name would converge them into one
// node and report a squat that does not exist.
func canonPackage(raw string) (string, error) {
	registry, name, ok := strings.Cut(strings.TrimSpace(raw), ":")
	if !ok {
		return "", fmt.Errorf("package %q is not registry-qualified (want e.g. npm:lodash)", raw)
	}
	registry = strings.ToLower(strings.TrimSpace(registry))
	name = strings.TrimSpace(name)
	if registry == "" || name == "" {
		return "", fmt.Errorf("package %q is not registry-qualified (want e.g. npm:lodash)", raw)
	}
	switch registry {
	case "pypi":
		name = pep503(name)
	default:
		// npm lowercases; so does every other registry whose rule we do not
		// model yet, and lowercasing is the safe default because it can only
		// merge two spellings of one name, never split one name in two.
		name = strings.ToLower(name)
	}
	if name == "" {
		return "", fmt.Errorf("package %q has no name", raw)
	}
	return registry + ":" + name, nil
}

// pep503 is PyPI's normalization: runs of "-", "_" and "." collapse to a single
// "-", then lowercase. It is what makes Foo.Bar, foo_bar and foo--bar one
// project, which is also what makes a typosquat of it detectable.
func pep503(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	sep := false
	for _, r := range strings.ToLower(s) {
		if r == '-' || r == '_' || r == '.' {
			sep = true
			continue
		}
		if sep {
			b.WriteByte('-')
			sep = false
		}
		b.WriteRune(r)
	}
	if sep {
		b.WriteByte('-')
	}
	return b.String()
}

// canonIP: the normalized v4 or v6 text form. Unmapping ::ffff:1.2.3.4 to
// 1.2.3.4 is the point — a DNS answer and a PTR answer must reach one node.
func canonIP(raw string) (string, error) {
	addr, err := netip.ParseAddr(strings.TrimSpace(raw))
	if err != nil {
		return "", fmt.Errorf("ip %q: %w", raw, err)
	}
	return addr.Unmap().String(), nil
}

// canonASN: "AS" plus decimal, no leading zeros. Whois, RDAP and routing feeds
// each spell it differently — 15169, AS15169, ASN00015169 — and all three name
// the same autonomous system.
func canonASN(raw string) (string, error) {
	s := strings.ToUpper(strings.TrimSpace(raw))
	s = strings.TrimPrefix(s, "ASN")
	s = strings.TrimPrefix(s, "AS")
	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return "", fmt.Errorf("asn %q is not AS<decimal>", raw)
	}
	return "AS" + strconv.FormatUint(n, 10), nil
}

// canonRegistrant: a normalized org or handle string. Whitespace runs collapse
// and case folds, because "Google  LLC" and "GOOGLE LLC" arriving from two
// registries are one registrant, and a campaign analyzer that cannot see that
// finds nothing.
func canonRegistrant(raw string) (string, error) {
	s := strings.ToLower(strings.Join(strings.Fields(raw), " "))
	if s == "" {
		return "", fmt.Errorf("empty registrant")
	}
	return s, nil
}
