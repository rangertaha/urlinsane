// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package report

import (
	"fmt"
	"strings"

	"github.com/rangertaha/urlinsane/internal/graph"
)

// Filter selects nodes for display. The vocabulary is §9's: the three existence
// values and a risk comparison. It replaces --registered/--unregistered, two
// booleans that were domain-only, settable to contradict each other, and could
// not express the third case at all.
type Filter struct {
	// spec is the filter exactly as the user wrote it, for String and for
	// error messages. Nothing may match on it: it carries the user's casing,
	// so `--filter Live` and `--filter live` are the same filter with
	// different spellings, and code comparing spec to a literal sees two
	// different things. keep did, and quietly turned the OR between two
	// existence filters into an AND that no node can satisfy — `--filter
	// Live,Absent` rendered an empty report and exited 0.
	spec string

	// family names the attribute this filter tests, when that attribute can
	// hold exactly one value per node — existence, type, and risk or depth
	// compared for equality. Two filters in the same family are ALTERNATIVES,
	// because no node can satisfy both and requiring both would render an empty
	// report; filters in different families narrow, as filters normally do.
	//
	// Empty for a comparison (`risk>high`, `depth<=2`), where conjunction is
	// meaningful and is what the user asked for.
	//
	// It is set here, from the parse that already knows the answer, rather than
	// recovered afterwards from the spec text. It began as a bool covering only
	// existence, which fixed `--filter Live,Absent` and left every other family
	// broken in exactly the same way: `--filter type=domain,type=ip` matched no
	// node at all and rendered an empty report with exit 0.
	family string

	fn func(NodeRow) bool
}

// String returns the filter as written, for help and error messages.
func (f Filter) String() string { return f.spec }

// ParseFilter compiles one filter expression.
//
//	live | absent | unknown | untried   existence (§9)
//	risk>N | risk>=N | risk=N           severity by name: info..critical
//	type=NAME                           node type
//	depth<=N                            observation hops from the seed
func ParseFilter(s string) (Filter, error) {
	spec := strings.TrimSpace(s)
	low := strings.ToLower(spec)

	switch low {
	case "live":
		return existence(spec, graph.Live), nil
	case "absent":
		return existence(spec, graph.Absent), nil
	case "unknown":
		return existence(spec, graph.Unknown), nil
	case "untried":
		return existence(spec, 0), nil
	}

	field, op, arg, ok := split(low)
	if !ok {
		return Filter{}, fmt.Errorf(
			"report: unrecognised filter %q; want live, absent, unknown, risk>SEV, type=NAME or depth<=N", s)
	}

	switch field {
	case "risk":
		sev, ok := graph.ParseSeverity(arg)
		if !ok {
			return Filter{}, fmt.Errorf(
				"report: %q is not a severity; want info, low, medium, high or critical", arg)
		}
		return Filter{spec: spec, family: equalityFamily("risk", op), fn: func(n NodeRow) bool {
			return compare(op, int(n.severity), int(sev))
		}}, nil

	case "type":
		if op != "=" && op != "==" {
			return Filter{}, fmt.Errorf("report: type filters compare with '=', not %q", op)
		}
		return Filter{spec: spec, family: "type", fn: func(n NodeRow) bool {
			return strings.EqualFold(n.Type, arg)
		}}, nil

	case "depth":
		var d int
		if _, err := fmt.Sscanf(arg, "%d", &d); err != nil {
			return Filter{}, fmt.Errorf("report: depth filter wants a number, got %q", arg)
		}
		return Filter{spec: spec, family: equalityFamily("depth", op), fn: func(n NodeRow) bool {
			return compare(op, n.Depth, d)
		}}, nil
	}
	return Filter{}, fmt.Errorf("report: unknown filter field %q in %q", field, s)
}

// ParseFilters compiles a filter list, reporting the first bad one. A filter
// that silently fails to parse is worse than a rejected one: the user believes
// they narrowed the report and reads a full one as if it were filtered.
func ParseFilters(specs []string) ([]Filter, error) {
	var out []Filter
	for _, s := range specs {
		for _, part := range strings.Split(s, ",") {
			if strings.TrimSpace(part) == "" {
				continue
			}
			f, err := ParseFilter(part)
			if err != nil {
				return nil, err
			}
			out = append(out, f)
		}
	}
	return out, nil
}

// TypeName returns the node type this filter selects on, if it is a type
// filter. It exists so a caller holding the registry can validate the name —
// this package deliberately does not know what types exist.
func (f Filter) TypeName() (string, bool) {
	field, op, arg, ok := split(strings.ToLower(f.spec))
	if !ok || field != "type" || (op != "=" && op != "==") {
		return "", false
	}
	return arg, true
}

// ValidateTypes rejects a type filter naming a type that is not registered.
//
// Without it `--filter type=domian` parses, matches nothing, and renders an
// empty report that looks exactly like a scan which found nothing — the same
// silent-narrowing failure the scope positional had. A filter the user believes
// is working and is not is worse than one that refuses to start.
func ValidateTypes(filters []Filter, known []string) error {
	if len(known) == 0 {
		return nil
	}
	set := make(map[string]bool, len(known))
	for _, k := range known {
		set[strings.ToLower(k)] = true
	}
	for _, f := range filters {
		name, ok := f.TypeName()
		if !ok || set[name] {
			continue
		}
		return fmt.Errorf(
			"report: filter %q names an unknown node type; want one of %s",
			f.spec, strings.Join(known, ", "))
	}
	return nil
}

func existence(spec string, want graph.Existence) Filter {
	return Filter{
		spec:   spec,
		family: "existence",
		fn:     func(n NodeRow) bool { return n.existence == want },
	}
}

// equalityFamily names the family for a comparison filter, but only when the
// comparison is equality — that is the case where two filters on one attribute
// are mutually exclusive and so must be alternatives. `risk>medium` and
// `risk>high` can both hold, so they keep narrowing.
func equalityFamily(field, op string) string {
	if op == "=" || op == "==" {
		return field
	}
	return ""
}

// split parses `field op arg`, longest operator first so ">=" is not read as
// ">" followed by a malformed argument.
func split(s string) (field, op, arg string, ok bool) {
	for _, o := range []string{">=", "<=", "==", "!=", ">", "<", "="} {
		if i := strings.Index(s, o); i > 0 {
			return s[:i], o, strings.TrimSpace(s[i+len(o):]), true
		}
	}
	return "", "", "", false
}

func compare(op string, got, want int) bool {
	switch op {
	case ">":
		return got > want
	case ">=":
		return got >= want
	case "<":
		return got < want
	case "<=":
		return got <= want
	case "!=":
		return got != want
	}
	return got == want
}

// keep reports whether a node survives the filter set.
//
// Filters within one family are alternatives — `--filter live --filter absent`,
// `--filter type=domain --filter type=ip` — because a node has exactly one
// existence and exactly one type, so requiring both would render an empty
// report. Across families they narrow: `--filter live --filter risk>medium`
// means both.
//
// The rule is per family rather than special-cased for existence. It was
// special-cased, and every other family was silently broken the same way the
// existence one had been: two type filters, or two `risk=` filters, matched no
// node and produced an empty report with exit 0 — a scan that looks like it
// found nothing rather than a query that cannot match. A filter kind declares
// its family where it is parsed and inherits this; nothing here needs to learn
// about it.
func keep(filters []Filter, n NodeRow) bool {
	if len(filters) == 0 {
		return true
	}
	// saw/matched per family, so an unmatched family fails the whole row while
	// a matched one is satisfied by any single member.
	saw := map[string]bool{}
	matched := map[string]bool{}
	for _, f := range filters {
		if f.family == "" {
			if !f.fn(n) {
				return false
			}
			continue
		}
		saw[f.family] = true
		if f.fn(n) {
			matched[f.family] = true
		}
	}
	for fam := range saw {
		if !matched[fam] {
			return false
		}
	}
	return true
}
