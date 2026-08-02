// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package all composes the shipped observation operators.
//
// Separate from the observe library so each operator can import that library —
// for Options, the per-call timeout and the whole schema vocabulary — without
// importing its siblings through it.
package all

import (
	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/observe"
	"github.com/rangertaha/urlinsane/internal/plugins/observe/dns"
	"github.com/rangertaha/urlinsane/internal/plugins/observe/geo"
	"github.com/rangertaha/urlinsane/internal/plugins/observe/idn"
	"github.com/rangertaha/urlinsane/internal/plugins/observe/pkg"
	"github.com/rangertaha/urlinsane/internal/plugins/observe/ptr"
	"github.com/rangertaha/urlinsane/internal/plugins/observe/repo"
	"github.com/rangertaha/urlinsane/internal/plugins/observe/usr"
	"github.com/rangertaha/urlinsane/internal/plugins/observe/whois"
)

// New builds every observation operator the supplied options can support.
//
// An operator whose dependency is missing is left out rather than included and
// left to fail: plan compilation reports what may run, and listing an operator
// that can only ever return an error would make --explain lie.
func New(o observe.Options) []graph.Operator {
	res := o.Resolver
	if res == nil {
		res = observe.SystemResolver()
	}
	who := o.Whois
	if who == nil {
		who = observe.DefaultWhois(o.Timeout)
	}
	probe := o.Prober
	if probe == nil {
		probe = observe.DefaultProber(o.Timeout)
	}

	ops := []graph.Operator{idn.New(o)}
	ops = append(ops, dns.New(o, res)...)
	ops = append(ops, ptr.New(o, res), whois.New(o, who))
	if o.Geo != nil {
		ops = append(ops, geo.New(o, o.Geo))
	}

	// Nil Sources falls back to the dataset; if that is unavailable too, the
	// three source operators are omitted rather than registered to fail.
	list := o.Sources
	if list == nil {
		list = observe.DatasetSources()
	}
	ops = append(ops, pkg.New(o, list, probe)...)
	ops = append(ops, usr.New(o, list, probe)...)
	ops = append(ops, repo.New(o, list, probe)...)
	return ops
}

// Select returns the named observation operators, in plan order (§12.10).
//
//	nil            every operator the options support
//	"dns","ptr"    only these
//	"^whois"       everything except this
//
// The set is built first and filtered second, so an operator omitted for a
// missing dependency — geo without a locator — is absent from the known ids
// too. Naming it is then an error that says so, rather than silently selecting
// nothing.
func Select(o observe.Options, ids ...string) ([]graph.Operator, error) {
	return graph.SelectOperators(New(o), ids)
}
