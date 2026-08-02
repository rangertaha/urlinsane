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

package graph

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/rangertaha/urlinsane/internal/graph/dag"
)

// SeedSpec is the target a plan was compiled for.
type SeedSpec struct {
	Type  string   `json:"type"`
	Key   string   `json:"key"`   // canonical
	Scope []string `json:"scope"` // nameable types to vary; empty means all in closure
}

// OpBinding is one operator as the plan sees it. This is the single definition;
// §5 and §10.6 of the design refer to it rather than restating it.
type OpBinding struct {
	Id       string            `json:"id"`
	Version  int               `json:"version"`
	Resource string            `json:"resource,omitempty"`
	OnTypes  []string          `json:"on_types,omitempty"`
	OnCaps   []string          `json:"on_caps,omitempty"`
	Where    []string          `json:"where,omitempty"`
	Reads    Reads             `json:"reads"`
	Emits    Effects           `json:"emits"`
	Config   map[string]string `json:"config,omitempty"`
	Models   []string          `json:"models,omitempty"`

	// Dead marks an operator whose Where can never be satisfied by anything
	// this plan produces. It is listed, not run, and not fatal: usually it
	// means the user narrowed scope.
	Dead    bool   `json:"dead,omitempty"`
	DeadWhy string `json:"dead_why,omitempty"`
}

// Plan is the compiled, inspectable, pinnable answer to "what will this do?".
type Plan struct {
	Hash      string      `json:"hash"`
	Seed      SeedSpec    `json:"seed"`
	Operators []OpBinding `json:"operators"`
	Model     string      `json:"model,omitempty"` // engine model CID
	Limits    Limits      `json:"limits"`

	// Pruned lists operators dropped as unreachable from the seed. Recording
	// them is what keeps a narrow plan from reading like a complete one.
	Pruned []string `json:"pruned,omitempty"`

	// flow is the type-flow graph, kept for rendering.
	flow dag.Graph
}

// PlanInput is everything compilation needs besides the operators themselves.
type PlanInput struct {
	Seed   SeedSpec
	Limits Limits
	Model  string                       // engine model CID
	Config map[string]map[string]string // operator id -> flag values
	Models map[string][]string          // operator id -> model CIDs
}

// Compile turns a registry, a seed and an operator set into a plan.
//
// Pruning is reachability over the type-flow graph, not a topological walk:
// that graph is cyclic, so reachability is transitive closure. It is an
// over-approximation because Emits is a *may*, not a *must*, and the plan says
// so rather than implying certainty.
func Compile(reg *Registry, ops []Operator, in PlanInput) (*Plan, error) {
	if _, ok := reg.Type(in.Seed.Type); !ok {
		return nil, fmt.Errorf("graph: unknown seed type %q", in.Seed.Type)
	}
	for _, s := range in.Seed.Scope {
		t, ok := reg.Type(s)
		if !ok {
			return nil, fmt.Errorf("graph: unknown scope type %q", s)
		}
		if t.cap != Nameable {
			return nil, fmt.Errorf("graph: scope type %q is %s, not nameable", s, t.cap)
		}
	}

	sorted := append([]Operator(nil), ops...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Id() < sorted[j].Id() })

	// Type flow: an operator that matches type A and emits type B is an edge
	// A -> B. Reachability from the seed type decides what stays.
	flow := dag.Graph{}
	for _, op := range sorted {
		for _, from := range matchedTypes(reg, op.Trigger().On) {
			for _, to := range op.Emits().Nodes {
				flow[from] = appendUnique(flow[from], to)
			}
			if _, seen := flow[from]; !seen {
				flow[from] = nil
			}
		}
	}
	reach := dag.Reachable(flow, in.Seed.Type)
	reach[in.Seed.Type] = true

	var kept []Operator
	var pruned []string
	for _, op := range sorted {
		if anyReachable(matchedTypes(reg, op.Trigger().On), reach) {
			kept = append(kept, op)
			continue
		}
		pruned = append(pruned, op.Id())
	}

	// What the surviving operators can produce, for dead-operator detection.
	emittedProps := map[string]bool{}
	emittedRels := map[string]bool{}
	for _, op := range kept {
		for _, p := range op.Emits().Props {
			emittedProps[p] = true
		}
		for _, r := range op.Emits().Rels {
			emittedRels[r] = true
		}
	}

	p := &Plan{Seed: in.Seed, Limits: in.Limits.withDefaults(), Model: in.Model, Pruned: pruned, flow: flow}
	for _, op := range kept {
		t := op.Trigger()
		b := OpBinding{
			Id:       op.Id(),
			Version:  op.Version(),
			Resource: op.Resource(),
			OnTypes:  append([]string(nil), t.On.Types...),
			Reads:    t.effectiveReads(),
			Emits:    op.Emits(),
			Config:   in.Config[op.Id()],
			Models:   in.Models[op.Id()],
		}
		for _, c := range t.On.Caps {
			b.OnCaps = append(b.OnCaps, c.String())
		}
		for _, c := range t.Where {
			b.Where = append(b.Where, c.describe())
		}
		if why, dead := deadReason(t, emittedProps, emittedRels); dead {
			b.Dead, b.DeadWhy = true, why
		}
		p.Operators = append(p.Operators, b)
	}
	p.Hash = p.computeHash()
	return p, nil
}

// deadReason reports whether a trigger can never be satisfied by anything the
// plan produces. This is why Effects covers props and relations rather than
// node types alone — without it a permanently dead operator would be listed as
// active, which defeats the point of --explain.
func deadReason(t Trigger, props, rels map[string]bool) (string, bool) {
	for _, c := range t.Where {
		fs, rs := c.declares()
		for _, f := range fs {
			if !props[f] {
				return "no operator in this plan emits prop " + f, true
			}
		}
		for _, r := range rs {
			if !rels[r] {
				return "no operator in this plan emits relation " + r, true
			}
		}
	}
	return "", false
}

func matchedTypes(reg *Registry, s Selector) []string {
	var out []string
	for _, n := range s.Types {
		out = appendUnique(out, n)
	}
	for _, c := range s.Caps {
		for _, t := range reg.Types() {
			if t.cap == c {
				out = appendUnique(out, t.name)
			}
		}
	}
	sort.Strings(out)
	return out
}

func anyReachable(types []string, reach map[string]bool) bool {
	for _, t := range types {
		if reach[t] {
			return true
		}
	}
	return false
}

func appendUnique(s []string, v string) []string {
	for _, x := range s {
		if x == v {
			return s
		}
	}
	return append(s, v)
}

// computeHash identifies the plan exactly. It covers the engine model and every
// plugin model as well as the operators and their config: a plan pinning
// operator versions but letting models float would reproduce neither the
// traversal nor, where a plugin model decides what an operator emits, the
// graph.
func (p *Plan) computeHash() string {
	h := sha256.New()
	writeField(h, p.Seed.Type)
	writeField(h, p.Seed.Key)
	for _, s := range p.Seed.Scope {
		writeField(h, s)
	}
	writeField(h, p.Model)
	writeField(h, fmt.Sprintf("depth=%d rounds=%d rev=%d att=%d budget=%d frontier=%d",
		p.Limits.MaxDepth, p.Limits.MaxRounds, p.Limits.Revisions,
		p.Limits.Attempts, p.Limits.NodeBudget, p.Limits.Frontier))
	for _, t := range sortedMapKeys(p.Limits.TypeBudget) {
		writeField(h, fmt.Sprintf("%s=%d", t, p.Limits.TypeBudget[t]))
	}
	for _, b := range p.Operators {
		writeField(h, b.Id)
		writeField(h, fmt.Sprintf("v%d", b.Version))
		writeField(h, b.Resource)
		for _, s := range b.OnTypes {
			writeField(h, s)
		}
		for _, s := range b.OnCaps {
			writeField(h, s)
		}
		for _, s := range b.Where {
			writeField(h, s)
		}
		for _, s := range b.Reads.Fields {
			writeField(h, s)
		}
		for _, s := range b.Reads.Rels {
			writeField(h, s)
		}
		for _, s := range b.Emits.Nodes {
			writeField(h, s)
		}
		for _, s := range b.Emits.Rels {
			writeField(h, s)
		}
		for _, s := range b.Emits.Props {
			writeField(h, s)
		}
		for _, k := range sortedMapKeys(b.Config) {
			writeField(h, k+"="+b.Config[k])
		}
		for _, m := range b.Models {
			writeField(h, m)
		}
	}
	return hex.EncodeToString(h.Sum(nil))
}

func sortedMapKeys[V any](m map[string]V) []string {
	if len(m) == 0 {
		return nil
	}
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// Write serializes the plan for --plan FILE.
func (p *Plan) Write(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(p)
}

// ReadPlan loads a pinned plan.
func ReadPlan(r io.Reader) (*Plan, error) {
	var p Plan
	if err := json.NewDecoder(r).Decode(&p); err != nil {
		return nil, err
	}
	if got := p.computeHash(); got != p.Hash {
		return nil, fmt.Errorf("graph: pinned plan hash mismatch: file says %s, contents hash to %s", p.Hash, got)
	}
	return &p, nil
}

// Select returns the operators named by the plan, in plan order. Execution
// dispatches only these; runtime never reaches past the plan.
func (p *Plan) Select(ops []Operator) []Operator {
	byID := map[string]Operator{}
	for _, o := range ops {
		byID[o.Id()] = o
	}
	var out []Operator
	for _, b := range p.Operators {
		if b.Dead {
			continue
		}
		if o, ok := byID[b.Id]; ok {
			out = append(out, o)
		}
	}
	return out
}

// Explain renders the plan: the type flow in SCC-condensed layers, then the
// operators, then anything dead or pruned.
func (p *Plan) Explain(w io.Writer) error {
	var b strings.Builder

	fmt.Fprintf(&b, "plan %s\n", p.Hash[:16])
	fmt.Fprintf(&b, "seed %s %s\n", p.Seed.Type, p.Seed.Key)
	if len(p.Seed.Scope) > 0 {
		fmt.Fprintf(&b, "scope %s\n", strings.Join(p.Seed.Scope, ","))
	} else {
		fmt.Fprintf(&b, "scope every nameable node in the seed closure\n")
	}
	fmt.Fprintf(&b, "limits depth=%d rounds=%d\n", p.Limits.MaxDepth, p.Limits.MaxRounds)
	if p.Model != "" {
		fmt.Fprintf(&b, "model %s\n", p.Model)
	}

	b.WriteString("\ntype flow\n")
	for _, c := range dag.Condense(p.flow) {
		marker := ""
		if c.Cyclic {
			// Worth surfacing: a cycle here is correct data, not a defect, and
			// is why nothing in the engine topologically sorts.
			marker = "  (cycle)"
		}
		fmt.Fprintf(&b, "  %d  %s%s\n", c.Layer, strings.Join(c.Members, ", "), marker)
	}

	b.WriteString("\noperators\n")
	for _, o := range p.Operators {
		if o.Dead {
			continue
		}
		on := strings.Join(append(append([]string{}, o.OnTypes...), o.OnCaps...), "|")
		fmt.Fprintf(&b, "  %-14s on %-18s", o.Id, on)
		if len(o.Where) > 0 {
			fmt.Fprintf(&b, " where %s", strings.Join(o.Where, " "))
		}
		if o.Resource != "" {
			fmt.Fprintf(&b, " [%s]", o.Resource)
		}
		b.WriteString("\n")
	}

	var dead []OpBinding
	for _, o := range p.Operators {
		if o.Dead {
			dead = append(dead, o)
		}
	}
	if len(dead) > 0 {
		b.WriteString("\ndead — listed, never dispatched\n")
		for _, o := range dead {
			fmt.Fprintf(&b, "  %-14s %s\n", o.Id, o.DeadWhy)
		}
	}
	if len(p.Pruned) > 0 {
		fmt.Fprintf(&b, "\npruned — unreachable from a %s seed\n  %s\n",
			p.Seed.Type, strings.Join(p.Pruned, ", "))
	}

	_, err := io.WriteString(w, b.String())
	return err
}
