// Copyright (C) 2024 Rangertaha
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
package com

// Algo

// Algo is taking the target domain name and adding additional words
// to it typically separated by a hyphen. example.com may end up looking like
// support-example[.]com. Something like this might find its way into a phishing
// message targeting an organization. It gives the appearance of it being I.T. support.

// Combo-squatting is a form of typosquatting that involves creating domain names combining a legitimate brand or word with another keyword. It’s a social engineering technique that aims to mislead users into believing the domain is related to a trusted website or service.

// For example, if the legitimate domain is amazon.com, a combo-squatted domain might be something like amazon-login.com or amazon-secure.com. These additional words or phrases make the fake domain look authentic or security-focused, increasing the chances that users will click the link.

// This tactic is used to trick users into visiting the site, often for malicious purposes like phishing attacks, credential harvesting, or spreading malware. Combo-squatting is particularly effective because the domains appear more trustworthy than simple typographical errors, as the added words can look plausible and contextually relevant.

import (
	"fmt"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/pkg/domain"
	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
)

const (
	CODE        = "com"
	NAME        = "Combo Squatting"
	DESCRIPTION = "Kreating domain names combining a legitimate brand or word with another keyword"
)

type Algo struct {
	config    internal.Config
	languages []internal.Language
	keyboards []internal.Keyboard
	funcs     map[int]func(internal.Typo) []internal.Typo
}

func (n *Algo) Id() string {
	return CODE
}

func (n *Algo) Init(conf internal.Config) {
	n.funcs = make(map[int]func(internal.Typo) []internal.Typo)
	n.keyboards = conf.Keyboards()
	n.languages = conf.Languages()
	n.config = conf

	// Supported targets
	n.funcs[internal.DOMAIN] = n.domain
	n.funcs[internal.PACKAGE] = n.name
	n.funcs[internal.EMAIL] = n.email
	n.funcs[internal.NAME] = n.name
}

func (n *Algo) Name() string {
	return NAME
}
func (n *Algo) Description() string {
	return DESCRIPTION
}

func (n *Algo) Exec(typo internal.Typo) []internal.Typo {
	return n.funcs[n.config.Type()](typo)
}

func (n *Algo) domain(typo internal.Typo) (typos []internal.Typo) {
	sub, prefix, suffix := typo.Original().Domain()

	for _, variant := range n.Func(prefix) {
		if prefix != variant {
			d := domain.New(sub, variant, suffix)

			new := typo.Clone(d.String())

			typos = append(typos, new)
		}
	}
	return
}

func (n *Algo) email(typo internal.Typo) (typos []internal.Typo) {
	username, domain := typo.Original().Email()

	for _, variant := range n.Func(username) {
		if username != variant {
			new := typo.Clone(fmt.Sprintf("%s@%s", variant, domain))

			typos = append(typos, new)
		}
	}
	return
}

func (n *Algo) name(typo internal.Typo) (typos []internal.Typo) {
	original := n.config.Target().Name()
	for _, variant := range n.Func(original) {
		if original != variant {
			typos = append(typos, typo.Clone(variant))
		}
	}
	return
}

func (n *Algo) Func(original string) (results []string) {
	for i, char := range original {
		for _, board := range n.keyboards {
			for _, kchar := range board.Adjacent(string(char)) {
				variant := fmt.Sprint(original[:i], kchar, original[i+1:])
				results = append(results, variant)
			}
		}
	}
	return results
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() internal.Algorithm {
		return &Algo{}
	})
}
