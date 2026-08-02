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
	"fmt"
	"time"
)

// Kind is the closed set of value kinds a prop may hold. Keeping it closed is
// what lets props encode deterministically; see docs/DESIGN.md §1.3.
type Kind uint8

const (
	KindString Kind = iota + 1
	KindInt
	KindFloat
	KindBool
	KindBytes
	KindTime
)

func (k Kind) String() string {
	switch k {
	case KindString:
		return "string"
	case KindInt:
		return "int"
	case KindFloat:
		return "float"
	case KindBool:
		return "bool"
	case KindBytes:
		return "bytes"
	case KindTime:
		return "time"
	}
	return "invalid"
}

// Value is a typed prop value. Times are held as Unix nanoseconds so the
// encoded form is an integer and carries no location or formatting.
type Value struct {
	kind Kind
	num  int64
	real float64
	str  string
	raw  []byte
}

func String(s string) Value { return Value{kind: KindString, str: s} }
func Int(i int64) Value     { return Value{kind: KindInt, num: i} }
func Float(f float64) Value { return Value{kind: KindFloat, real: f} }
func Bool(b bool) Value {
	n := int64(0)
	if b {
		n = 1
	}
	return Value{kind: KindBool, num: n}
}
func Bytes(b []byte) Value   { return Value{kind: KindBytes, raw: append([]byte(nil), b...)} }
func Time(t time.Time) Value { return Value{kind: KindTime, num: t.UTC().UnixNano()} }

func (v Value) Kind() Kind      { return v.kind }
func (v Value) Str() string     { return v.str }
func (v Value) Num() int64      { return v.num }
func (v Value) Real() float64   { return v.real }
func (v Value) Flag() bool      { return v.num != 0 }
func (v Value) Raw() []byte     { return v.raw }
func (v Value) Time() time.Time { return time.Unix(0, v.num).UTC() }
func (v Value) IsZero() bool    { return v.kind == 0 }

// equal reports value equality. Used to tell a redundant assertion (no change)
// from a genuine conflict.
func (v Value) equal(o Value) bool {
	if v.kind != o.kind {
		return false
	}
	switch v.kind {
	case KindString:
		return v.str == o.str
	case KindFloat:
		return v.real == o.real
	case KindBytes:
		return string(v.raw) == string(o.raw)
	default:
		return v.num == o.num
	}
}

// MergePolicy resolves competing assertions of the same field. Resolution must
// not depend on arrival order — under concurrent dispatch that would be decided
// by network timing. See docs/DESIGN.md §1.4.
type MergePolicy struct {
	order []string
}

// Precedence declares operator ids in priority order, highest first. Operators
// not listed rank behind every listed one, and ties break on lowest id.
func Precedence(operators ...string) MergePolicy {
	return MergePolicy{order: append([]string(nil), operators...)}
}

func (m MergePolicy) rank(op string) int {
	for i, o := range m.order {
		if o == op {
			return i
		}
	}
	return len(m.order)
}

// wins reports whether an assertion by next should replace one by cur.
func (m MergePolicy) wins(cur, next string) bool {
	rc, rn := m.rank(cur), m.rank(next)
	if rn != rc {
		return rn < rc
	}
	return next < cur
}

// FieldDef declares one field of a node or relation type. Order in the
// declaring slice is the field's stable index and part of the on-disk contract:
// fields are append-only, and removal is by tombstone (Deprecated).
type FieldDef struct {
	Name       string
	Kind       Kind
	Merge      MergePolicy
	Deprecated bool
}

// Field is a resolved handle to a declared field. Operators obtain one at
// registration, so an unknown name fails then rather than returning a
// per-access boolean at runtime.
type Field struct {
	sch   *schema
	index int
	name  string
	kind  Kind
	merge MergePolicy
}

func (f Field) Name() string  { return f.name }
func (f Field) Kind() Kind    { return f.kind }
func (f Field) Owner() string { return f.sch.owner }
func (f Field) valid() bool   { return f.sch != nil }

// schema is the ordered field list shared by every Props of one type.
type schema struct {
	owner   string
	version int
	fields  []Field
	byName  map[string]Field
}

func newSchema(owner string, version int, defs []FieldDef) (*schema, error) {
	s := &schema{owner: owner, version: version, byName: make(map[string]Field, len(defs))}
	for i, d := range defs {
		if d.Name == "" {
			return nil, fmt.Errorf("graph: %s field %d has no name", owner, i)
		}
		if d.Kind == 0 || d.Kind > KindTime {
			return nil, fmt.Errorf("graph: %s field %q has invalid kind", owner, d.Name)
		}
		if _, dup := s.byName[d.Name]; dup {
			return nil, fmt.Errorf("graph: %s field %q declared twice", owner, d.Name)
		}
		f := Field{sch: s, index: i, name: d.Name, kind: d.Kind, merge: d.Merge}
		s.fields = append(s.fields, f)
		s.byName[d.Name] = f
	}
	return s, nil
}

// Field resolves a name to a handle.
func (s *schema) Field(name string) (Field, bool) {
	f, ok := s.byName[name]
	return f, ok
}

// Props holds one value per declared field, addressed positionally. Order is a
// property of the type rather than something imposed at encode time, so
// identical values encode identically by construction.
type Props struct {
	sch   *schema
	slots []slot
}

type slot struct {
	set bool
	val Value
	op  string // operator id that won this field
}

func newProps(s *schema) Props {
	return Props{sch: s, slots: make([]slot, len(s.fields))}
}

// Get returns the materialized value of f and whether it has been set.
func (p Props) Get(f Field) (Value, bool) {
	if p.sch == nil || !f.valid() || f.sch != p.sch {
		return Value{}, false
	}
	sl := p.slots[f.index]
	return sl.val, sl.set
}

// Setter returns the operator id whose assertion is currently materialized.
func (p Props) Setter(f Field) (string, bool) {
	if p.sch == nil || !f.valid() || f.sch != p.sch {
		return "", false
	}
	sl := p.slots[f.index]
	return sl.op, sl.set
}

// Each yields every set field in declaration order.
func (p Props) Each(fn func(Field, Value)) {
	if p.sch == nil {
		return
	}
	for i, sl := range p.slots {
		if sl.set {
			fn(p.sch.fields[i], sl.val)
		}
	}
}

// assert applies one operator's assertion under the field's merge policy and
// reports whether the materialized value changed.
func (p *Props) assert(f Field, v Value, op string) (bool, error) {
	if p.sch == nil || !f.valid() || f.sch != p.sch {
		return false, fmt.Errorf("graph: field %q does not belong to %s", f.name, p.owner())
	}
	if v.Kind() != f.kind {
		return false, fmt.Errorf("graph: field %s.%s is %s, got %s", p.owner(), f.name, f.kind, v.Kind())
	}
	sl := &p.slots[f.index]
	if !sl.set {
		*sl = slot{set: true, val: v, op: op}
		return true, nil
	}
	if !f.merge.wins(sl.op, op) {
		return false, nil
	}
	changed := !sl.val.equal(v)
	sl.val, sl.op = v, op
	return changed, nil
}

func (p Props) owner() string {
	if p.sch == nil {
		return "<none>"
	}
	return p.sch.owner
}
