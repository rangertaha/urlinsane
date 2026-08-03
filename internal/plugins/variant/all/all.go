// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package all composes the shipped algorithms.
//
// Separate from the variant library so each algorithm can import that library —
// for Spec, the Generate signature and the keyboard/language combinators —
// without importing its siblings through it.
//
// Every algorithm is one package, and this is the only list of them. An
// algorithm added here without a directory, or a directory without a line
// here, is the kind of drift the old registry made possible.
package all

import (
	"fmt"
	"sort"

	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/variant"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/aci"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/acs"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/afx"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/bf"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/cb"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/cm"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/cns"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/co"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/cr"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/cs"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/dhs"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/di"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/do"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/fsd"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/gi"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/gr"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/hi"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/ho"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/hr"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/hs"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/nsc"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/ons"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/rar"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/sep"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/si"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/sld"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/sp"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/tld"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/tli"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/tos"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/vs"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/xhs"
)

// Specs returns every algorithm declaration, in id order.
//
// Sorted because the plan is compiled from this set and hashed (§5): the order
// packages happen to be listed in must not reach the plan hash.
func Specs(o variant.Options) []variant.Spec {
	o = o.WithDefaults()
	specs := []variant.Spec{
		co.Spec(),
		cr.Spec(),
		cs.Spec(),
		hi.Spec(),
		ho.Spec(),
		di.Spec(),
		do.Spec(),
		dhs.Spec(),
		bf.Spec(),
		sp.Spec(),
		tos.Spec(),
		afx.Spec(),
		sep.Spec(),
		nsc.Spec(),
		si.Spec(o.Subdomains),
		tld.Spec(o.Suffixes),
		sld.Spec(o.Suffixes),
		tli.Spec(o.Suffixes),
		fsd.Spec(o.Providers),
		aci.Spec(o.Keyboards),
		acs.Spec(o.Keyboards),
		rar.Spec(o.Keyboards),
		vs.Spec(o.Languages),
		hr.Spec(o.Languages),
		hs.Spec(o.Languages),
		xhs.Spec(o.CrossHomophones),
		cm.Spec(o.Languages),
		gi.Spec(o.Languages),
		gr.Spec(o.Languages),
		cns.Spec(o.Languages),
		ons.Spec(o.Languages),
		cb.Spec(o.Combos),
	}
	specs = append(specs, o.Extra...)
	sort.Slice(specs, func(i, j int) bool { return specs[i].ID < specs[j].ID })
	return specs
}

// All returns every algorithm as an operator, in id order. Two ids that collide
// are a programming error rather than a runtime condition — the scheduler keys
// its seen-set and cache on the id — so this panics rather than silently
// dropping one.
func All(o variant.Options) []graph.Operator {
	o = o.WithDefaults()
	specs := Specs(o)
	seen := make(map[string]bool, len(specs))
	ops := make([]graph.Operator, 0, len(specs))
	for _, s := range specs {
		if seen[s.ID] {
			panic(fmt.Sprintf("variant: duplicate operator id %q", s.ID))
		}
		seen[s.ID] = true
		ops = append(ops, variant.New(s, o.Split))
	}
	return ops
}

// Select returns the named algorithms as operators, in id order (§12.10).
//
//	nil            every algorithm
//	"co","cr"      only these
//	"^cb","^afx"   everything except these
//
// An unknown id is an error rather than a silent omission: a run that quietly
// dropped half the algorithms the user asked for would report an incomplete
// graph as complete.
func Select(o variant.Options, ids ...string) ([]graph.Operator, error) {
	return graph.SelectOperators(All(o), ids)
}
