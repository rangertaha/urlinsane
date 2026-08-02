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

package observe

import (
	"context"

	"strings"

	"github.com/rangertaha/urlinsane/internal/graph"
	"golang.org/x/net/idna"
)

// idn records the Unicode form of an internationalized domain.
//
// The old collector computed the opposite direction — punycode from the name —
// but a domain key is canonicalized to punycode at admission (DESIGN §2), so
// that value is now the key itself and asserting it would be noise. What is
// genuinely missing is the human-readable form: xn--80ak6aa92e.com is not
// something a report can put in front of a user, and homograph attacks are only
// legible in Unicode.
//
// It makes no network call, so it declares no rate-limit class. It sits here
// rather than with the decomposers because it is an assertion about an existing
// node rather than a structural expansion of the seed — a defensible line, but
// a thin one.
type idn struct{ base }

func newIDN(o Options) graph.Operator { return idn{o.base()} }

func (idn) Id() string       { return "idn" }
func (idn) Version() int     { return 1 }
func (idn) Resource() string { return "" }

func (idn) Trigger() graph.Trigger {
	return graph.Trigger{On: graph.Selector{Types: []string{TypeDomain}}}
}

func (idn) Emits() graph.Effects {
	return graph.Effects{Props: []string{FieldUnicode}}
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
		Node: &self, Field: FieldUnicode, Value: graph.String(unicode),
	}}}, graph.OK()
}
