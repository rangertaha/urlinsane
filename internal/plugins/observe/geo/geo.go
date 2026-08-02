// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package geo locates an address.
//
// It binds On: ip, so it runs wherever an address exists rather than declaring
// that it needs dns-a — the dependency the old collector declared, which broke
// as soon as anything else produced addresses (§4.1).
//
// With no geolocation database it contributes nothing at all. An operator that
// can only ever fail should not be in the compiled plan, or --explain promises
// work that cannot happen.
package geo

import (
	"context"

	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/observe"
)

type geo struct {
	observe.Base
	db observe.GeoLocator
}

func newGeo(o observe.Options, db observe.GeoLocator) graph.Operator {
	return geo{Base: o.Base(), db: db}
}

func (geo) Id() string       { return "geo" }
func (geo) Version() int     { return 1 }
func (geo) Resource() string { return observe.ResourceGeo }

func (geo) Trigger() graph.Trigger {
	return graph.Trigger{On: graph.Selector{Types: []string{observe.TypeIP}}}
}

func (geo) Emits() graph.Effects {
	// No Nodes and no Rels: a location has no identity worth converging on, so
	// it is props on the address rather than a node of its own (DESIGN §2).
	return graph.Effects{Props: []string{
		observe.FieldCity, observe.FieldCountry, observe.FieldCountryCode, observe.FieldContinent,
		observe.FieldLatitude, observe.FieldLongitude, observe.FieldTimeZone, observe.FieldPostalCode,
	}}
}

func (o geo) Exec(ctx context.Context, v graph.View) (graph.Delta, graph.Outcome) {
	loc, err := o.db.Locate(v.Key())
	if err != nil {
		return graph.Delta{}, graph.Failed(err)
	}
	if loc == nil {
		// The database answered and has nothing for this address. Absence, not
		// failure — the same distinction NXDOMAIN gets.
		return graph.Delta{}, graph.Empty()
	}

	self := v.Ref()
	var d graph.Delta
	str := func(field, value string) {
		if value == "" {
			return
		}
		d.Props = append(d.Props, graph.PropSet{Node: &self, Field: field, Value: graph.String(value)})
	}
	str(observe.FieldCity, loc.City)
	str(observe.FieldCountry, loc.Country)
	str(observe.FieldCountryCode, loc.CountryCode)
	str(observe.FieldContinent, loc.Continent)
	str(observe.FieldTimeZone, loc.TimeZone)
	str(observe.FieldPostalCode, loc.PostalCode)

	// MaxMind uses 0,0 for "unknown", so recording it would place half the
	// internet in the Gulf of Guinea.
	if loc.Latitude != 0 || loc.Longitude != 0 {
		d.Props = append(d.Props,
			graph.PropSet{Node: &self, Field: observe.FieldLatitude, Value: graph.Float(loc.Latitude)},
			graph.PropSet{Node: &self, Field: observe.FieldLongitude, Value: graph.Float(loc.Longitude)},
		)
	}
	if len(d.Props) == 0 {
		return graph.Delta{}, graph.Empty()
	}
	return d, graph.OK()
}

// New builds the geolocation operator. Callers must not call it with a nil
// locator: an operator that can only fail must not reach the plan.
func New(o observe.Options, db observe.GeoLocator) graph.Operator { return newGeo(o, db) }
