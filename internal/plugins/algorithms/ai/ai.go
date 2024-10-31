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
package ai

// Typo-squatting involving alphabet insertion occurs when attackers register domain names that include additional letters or characters within a legitimate brand's name. This tactic aims to exploit common typing errors by inserting one or more letters, creating a slight variation that may not be immediately noticeable to users.

// For example, if the legitimate domain is "example.com," a typo-squatter might register "exampel.com" or "exampxle.com." Users who accidentally mistype or overlook the extra characters might be redirected to these malicious sites, which could be used for phishing, distributing malware, or other harmful activities.

// This method capitalizes on the likelihood of users making small mistakes while typing, underscoring the need for brands to monitor and protect their online identities against such threats.


// func AlgoFunc(tc Result) (results []Result) {
// 	alphabet := map[string]bool{}
// 	for _, keyboard := range tc.Keyboards {
// 		for _, a := range keyboard.Language.Graphemes {
// 			alphabet[a] = true
// 		}
// 	}
// 	for i, char := range tc.Original.Domain {
// 		for alp := range alphabet {
// 			d1 := tc.Original.Domain[:i] + alp + string(char) + tc.Original.Domain[i+1:]
// 			if i == len(tc.Original.Domain)-1 {
// 				d1 = tc.Original.Domain[:i] + string(char) + alp + tc.Original.Domain[i+1:]
// 			}
// 			dm1 := Domain{tc.Original.Subdomain, d1, tc.Original.Suffix, Meta{}, false}
// 			results = append(results, Result{Original: tc.Original, Variant: dm1, Typo: tc.Typo, Data: tc.Data})
// 		}
// 	}
// 	return
// }

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
)

const (
	CODE        = "ai"
	NAME        = "Alphabet Insertion"
	DESCRIPTION = "Inserting the language specific alphabet in the target domain"
)

type Algo struct {
	types []string
}

func (n *Algo) Id() string {
	return CODE
}
func (n *Algo) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *Algo) Name() string {
	return "Alphabet Insertion"
}

func (n *Algo) Description() string {
	return "Inserting the language specific alphabet in the target domain"
}

func (n *Algo) Fields() []string {
	return []string{}
}

func (n *Algo) Headers() []string {
	return []string{}
}

func (n *Algo) Exec(in internal.Typo) (out []internal.Typo) {
	out = append(out, in)
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() internal.Algorithm {
		return &Algo{
			types: []string{algorithms.ENTITY, algorithms.DOMAIN},
		}
	})
}
