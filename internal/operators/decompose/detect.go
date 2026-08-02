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
	"strconv"
	"strings"

	"golang.org/x/net/publicsuffix"
)

// Detect infers a target's node type from the string alone.
//
// This is the seeding path, and it retires entity.Classify — which read
// bob@example.com as a *user* because it contained an "@", collapsing the
// address and its local part into one thing. Here an email is an email; the
// local part becomes a username node through decomposition, at which point both
// exist and both can be varied.
//
// Detection never depends on the scope positional (§12). `urlinsane typo
// username bob@example.com` and `urlinsane typo bob@example.com` parse the
// target identically; scope narrows what gets *varied*, not what the string
// *is*. Tying the two together would make the same target mean different things
// on different runs, and the seed closure would differ with it.
//
// The rules are ordered most-specific first, and each is a positive test rather
// than a fallthrough:
//
//  1. an "@" with a host after it — email
//  2. host/owner/name, or any URL or scp-style remote of one — repo
//  3. registry:name — package
//  4. a dotted name whose suffix is in the public suffix list — domain
//  5. anything else that is a legal handle — username
//
// Rule 4 is what separates `lodash` from `lodash.com`: both are legal
// hostnames, so a hostname test alone would call every bare package name a
// domain and scan DNS for it. Requiring a real registry suffix is the only
// discriminator that does not need a list of known package names.
func Detect(target string) (string, error) {
	typ, _, err := classify(target)
	return typ, err
}

// DetectSeed returns the target's type and canonical key together, which is
// what seeding needs.
//
// The two must be derived in one pass. Detecting against the string the user
// typed and then canonicalizing the same string separately breaks the moment
// the two disagree about what to strip: "https://example.com" detects as a
// domain, but canonDomain refuses it outright because a bare host may not
// contain a delimiter.
func DetectSeed(target string) (typ, key string, err error) {
	typ, cleaned, err := classify(target)
	if err != nil {
		return "", "", err
	}
	key, err = canonicalFor(typ, cleaned)
	if err != nil {
		return "", "", err
	}
	return typ, key, nil
}

// classify returns the detected type and the form of the target that type's
// canonicalizer should be given.
func classify(target string) (typ, cleaned string, err error) {
	s := strings.TrimSpace(target)
	if s == "" {
		return "", "", fmt.Errorf("empty target")
	}

	if at := strings.LastIndex(s, "@"); at > 0 && !looksLikeRepo(s) {
		if _, err := canonEmail(s); err == nil {
			return TypeEmail, s, nil
		}
	}
	if looksLikeRepo(s) {
		if _, err := canonRepo(s); err == nil {
			return TypeRepo, s, nil
		}
	}
	if registry, _, ok := strings.Cut(s, ":"); ok && isRegistryToken(registry) {
		if _, err := canonPackage(s); err == nil {
			return TypePackage, s, nil
		}
	}

	// A pasted URL is the common case for a domain target, so the scheme, path,
	// port and query are stripped before the hostname test rather than making
	// the user do it. They carry no identity: https://example.com/about and
	// example.com are the same name, and minting separate nodes for them would
	// defeat the convergence the whole model rests on.
	host := hostOf(s)
	if key, err := canonDomain(host); err == nil && hasRegistrySuffix(key) {
		return TypeDomain, host, nil
	}
	if _, err := canonUsername(s); err == nil {
		return TypeUsername, s, nil
	}
	return "", "", fmt.Errorf(
		"cannot tell what %q is: expected an email, a repo URL, a registry:package, a domain or a username", target)
}

// hostOf reduces a URL to its hostname, and leaves anything that is not one
// alone so the later rules still see what the user typed.
func hostOf(s string) string {
	if i := strings.Index(s, "://"); i >= 0 {
		s = s[i+3:]
	}
	// Userinfo, but only when there is any: a leading "@" is a malformed email,
	// not a URL with an empty user, and stripping it would turn "@example.com"
	// into a perfectly good domain the user never asked for.
	if at := strings.Index(s, "@"); at > 0 {
		s = s[at+1:]
	}
	s = strings.TrimSuffix(s, "/")
	for _, sep := range []string{"/", "?", "#"} {
		if i := strings.Index(s, sep); i >= 0 {
			s = s[:i]
		}
	}
	// A port, but not an IPv6 literal's colons.
	if i := strings.LastIndex(s, ":"); i >= 0 && !strings.Contains(s, "]") {
		if _, err := strconv.Atoi(s[i+1:]); err == nil {
			s = s[:i]
		}
	}
	return s
}

func canonicalFor(typ, raw string) (string, error) {
	for _, d := range nodeTypes() {
		if d.Name == typ {
			return d.Canonical(raw)
		}
	}
	return "", fmt.Errorf("no canonicalizer for %q", typ)
}

// looksLikeRepo reports whether a string names a repository rather than merely
// containing a slash. Three path segments is the test — host, owner, name —
// after any scheme, userinfo or scp-style separator is removed, so a browser
// URL, a clone URL and a bare path all answer the same way.
func looksLikeRepo(s string) bool {
	s = strings.TrimSuffix(strings.TrimSpace(s), "/")
	scheme := false
	if i := strings.Index(s, "://"); i >= 0 {
		s, scheme = s[i+3:], true
	}
	if at := strings.Index(s, "@"); at >= 0 {
		s = s[at+1:]
		if !scheme {
			s = strings.Replace(s, ":", "/", 1)
		}
	}
	return strings.Count(strings.TrimSuffix(s, ".git"), "/") >= 2
}

// isRegistryToken reports whether the text before a ":" names a package
// registry rather than a URL scheme or a host with a port.
//
// The test is structural, not a list of known registries: a registry name is a
// bare word — npm, pypi, cargo, gem, maven — so a dot or a slash means the
// colon is doing some other job. Without it, "example.com:8080" reads as
// package "8080" on registry "example.com", and "https://x" as a package on
// registry "https". A closed list would have been worse: it would reject every
// registry added after this line was written.
func isRegistryToken(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" || strings.ContainsAny(s, "./\\ \t") {
		return false
	}
	switch strings.ToLower(s) {
	case "http", "https", "git", "ssh", "ftp", "file":
		return false
	}
	return true
}

// hasRegistrySuffix reports whether a host ends in a suffix someone can
// actually register under. It is what keeps `lodash` from being scanned as a
// domain, and it accepts private suffixes too — acme.blogspot.com is a real
// name a squatter can take, whatever the ICANN flag says about it.
func hasRegistrySuffix(key string) bool {
	if !strings.Contains(key, ".") {
		return false
	}
	suffix, _ := publicsuffix.PublicSuffix(key)
	return suffix != "" && suffix != key
}
