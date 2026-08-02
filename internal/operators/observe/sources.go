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
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rangertaha/urlinsane/internal/dataset"
	"github.com/rangertaha/urlinsane/internal/graph"
)

// Source is one registry or platform a name may exist on. Template is the page
// a human should be shown; CheckURL is the endpoint that answers existence
// cleanly, which is often an API rather than the page.
type Source struct {
	Code     string
	Template string
	CheckURL string
}

// SourceLister supplies the sources for a kind ("package", "username",
// "repository"). Behind an interface so tests need neither a database nor the
// network.
type SourceLister interface {
	Sources(kind string) ([]Source, error)
}

// Prober decides whether a URL exists. The bool is the answer; a non-nil error
// means no answer was obtained, which is a different thing entirely and is what
// keeps "not found" from swallowing "could not reach".
type Prober interface {
	Exists(ctx context.Context, rawURL string) (bool, error)
}

// datasetSources reads the source lists out of the dataset, or returns nil when
// there is no dataset — in which case New omits the source operators rather
// than planning operators that cannot answer.
func datasetSources() SourceLister {
	if dataset.DB == nil {
		return nil
	}
	return datasetLister{}
}

type datasetLister struct{}

func (datasetLister) Sources(kind string) ([]Source, error) {
	if dataset.DB == nil {
		return nil, errors.New("observe: no dataset database")
	}
	var rows []dataset.Source
	if err := dataset.DB.Where("type = ?", kind).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]Source, 0, len(rows))
	for _, r := range rows {
		out = append(out, Source{Code: r.Code, Template: r.Template, CheckURL: r.CheckURL})
	}
	return out, nil
}

// httpProber answers existence with a GET, treating any 2xx as present.
type httpProber struct{ client *http.Client }

func defaultProber(timeout time.Duration) Prober {
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	return httpProber{client: &http.Client{Timeout: timeout}}
}

func (p httpProber) Exists(ctx context.Context, rawURL string) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("User-Agent", "urlinsane")

	resp, err := p.client.Do(req)
	if err != nil {
		// A transport error proves nothing about the name. The old collector
		// returned false here, silently reporting an unreachable registry as a
		// name that is free to squat.
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return true, nil
	}
	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusGone {
		return false, nil
	}
	// 403, 429, 5xx: the registry declined to answer, which is not an answer.
	return false, errors.New("observe: " + rawURL + ": " + resp.Status)
}

// sourceOp checks whether a name exists on the platforms of one kind, and links
// it to each platform it was found on with an EXISTS_ON edge. This is the
// supply-chain surface: a package absent from every public registry is the
// dependency-confusion case, and it is only detectable because absence is
// recorded as such rather than as a failure.
type sourceOp struct {
	base
	id       string
	on       string
	kind     string
	resource string
	sources  SourceLister
	prober   Prober
}

// newSourceOps builds the three source-checking operators. They are one
// implementation bound to three types: the work is identical, only the pattern
// and the source list differ.
func newSourceOps(o Options, list SourceLister, prober Prober) []graph.Operator {
	resource := o.SourceResource
	if resource == "" {
		resource = ResourceHTTP
	}
	mk := func(id, on, kind string) graph.Operator {
		return sourceOp{
			base: o.base(), id: id, on: on, kind: kind,
			resource: resource, sources: list, prober: prober,
		}
	}
	return []graph.Operator{
		mk("pkg", TypePackage, "package"),
		mk("usr", TypeUsername, "username"),
		// The old repo collector bound to package and name entities; a repo is
		// its own type now, so it binds where it belongs.
		mk("repo", TypeRepo, "repository"),
	}
}

func (o sourceOp) Id() string       { return o.id }
func (o sourceOp) Version() int     { return 1 }
func (o sourceOp) Resource() string { return o.resource }

func (o sourceOp) Trigger() graph.Trigger {
	return graph.Trigger{On: graph.Selector{Types: []string{o.on}}}
}

func (o sourceOp) Emits() graph.Effects {
	return graph.Effects{
		Nodes: []string{TypePlatform},
		Rels:  []string{RelExistsOn},
		Props: []string{FieldURL, FieldCode},
	}
}

func (o sourceOp) Exec(ctx context.Context, v graph.View) (graph.Delta, graph.Outcome) {
	sources, err := o.sources.Sources(o.kind)
	if err != nil {
		return graph.Delta{}, graph.Failed(err)
	}
	if len(sources) == 0 {
		// Nothing was checked, so nothing was determined. Reporting this as
		// absence would turn a missing dataset into a page of free names.
		return graph.Delta{}, graph.Failed(errors.New("observe: no " + o.kind + " sources configured"))
	}

	ctx, cancel := o.call(ctx)
	defer cancel()

	var d graph.Delta
	self := v.Ref()
	var undetermined error

	for _, s := range sources {
		if ctx.Err() != nil {
			undetermined = ctx.Err()
			break
		}
		check := s.CheckURL
		if check == "" {
			check = s.Template
		}
		found, perr := o.prober.Exists(ctx, expand(check, v.Key()))
		if perr != nil {
			undetermined = perr
			continue
		}
		if !found {
			continue
		}

		display := expand(s.Template, v.Key())
		ref := graph.NodeRef{Type: TypePlatform, Key: platformKey(s)}
		edge := graph.EdgeRef{From: self, Rel: RelExistsOn, To: ref}
		d.Nodes = append(d.Nodes, ref)
		d.Edges = append(d.Edges, edge)
		d.Props = append(d.Props,
			graph.PropSet{Edge: &edge, Field: FieldURL, Value: graph.String(display)},
			graph.PropSet{Node: &ref, Field: FieldCode, Value: graph.String(s.Code)},
		)
	}

	switch {
	case len(d.Nodes) > 0:
		return d, graph.OK()
	case undetermined != nil && errors.Is(undetermined, context.DeadlineExceeded):
		return graph.Delta{}, graph.Timeout(undetermined)
	case undetermined != nil:
		// Some registry never answered, so "absent everywhere" is not proven.
		// Only a clean sweep of definite noes earns Empty.
		return graph.Delta{}, graph.Failed(undetermined)
	}
	return graph.Delta{}, graph.Empty()
}

// expand substitutes the name into a source's URL template.
func expand(template, name string) string {
	return strings.ReplaceAll(template, "%s", url.PathEscape(name))
}

// platformKey names the platform node: its host, per the registry's
// canonicalization for the type. A source whose template is not a URL falls
// back to its code, which is all the identity it has.
//
// The placeholder is dropped before parsing because "%s" is an invalid percent
// escape and url.Parse rejects the whole template over it — which silently sent
// every platform key to the fallback.
func platformKey(s Source) string {
	template := strings.ReplaceAll(s.Template, "%s", "")
	if u, err := url.Parse(template); err == nil && u.Host != "" {
		return u.Host
	}
	return s.Code
}
