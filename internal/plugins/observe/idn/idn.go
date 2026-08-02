// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package idn records the unicode spelling of a punycode label.
//
// It makes no network call, but it is an observation: it asserts a prop that
// was not derivable from the key alone. A homoglyph variant whose xn-- form is
// unreadable is exactly the case a human reading the report needs spelled out.
package idn

import (
	"context"
	"strings"

	"golang.org/x/net/idna"

	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/observe"
)

type idn struct{ observe.Base }

func newIDN(o observe.Options) graph.Operator { return idn{o.Base()} }

func (idn) Id() string       { return "idn" }
func (idn) Version() int     { return 1 }
func (idn) Resource() string { return "" }

func (idn) Trigger() graph.Trigger {
	return graph.Trigger{On: graph.Selector{Types: []string{observe.TypeDomain}}}
}

func (idn) Emits() graph.Effects {
	return graph.Effects{Props: []string{observe.FieldUnicode}}
}

func (idn) Exec(ctx context.Context, v graph.View) (graph.Delta, graph.Outcome) {
	key := v.Key()

	// ToUnicode is lenient by design; a label it cannot decode comes back
	// unchanged, which is the right answer for an ASCII domain.
	unicode, err := idna.ToUnicode(key)
	if err != nil {
		return graph.Delta{}, graph.Failed(err)
	}
	if unicode == "" || strings.EqualFold(unicode, key) {
		// Not internationalized. Nothing was learned and nothing is missing:
		// the absence of a Unicode form is the finding.
		return graph.Delta{}, graph.Empty()
	}
	self := v.Ref()
	return graph.Delta{Props: []graph.PropSet{{
		Node: &self, Field: observe.FieldUnicode, Value: graph.String(unicode),
	}}}, graph.OK()
}

// New builds the unicode-spelling operator.
func New(o observe.Options) graph.Operator { return newIDN(o) }
