// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package observe

import "github.com/rangertaha/urlinsane/internal/graph"

// Fields returns the node props these operators assert that the schema does not
// already declare, grouped by node type.
//
// internal/plugins/decompose owns registration and registers field lists
// "at the minimum the design names", to be appended to as collectors land.
// These are that append. Until they are registered, every assertion below is
// refused as an unknown field — and refused quietly, because the applier
// records a rejection rather than failing the run.
//
// Append them *after* the registering package's own fields and never reorder
// them: a field's position is its stable index and part of the on-disk contract
// (DESIGN §1.3).
func Fields() map[string][]graph.FieldDef {
	return map[string][]graph.FieldDef{
		TypeDomain: {
			{Name: FieldUnicode, Kind: graph.KindString},
			{Name: FieldCNAME, Kind: graph.KindString},
			// One slot, many records: props are single-valued, so the TXT set is
			// joined rather than lost. See text.Exec.
			{Name: FieldTXT, Kind: graph.KindString},
			// rdap outranks whois wherever both answer, so the materialized
			// value does not depend on which service replied first (DESIGN §1.4).
			// "created" is already registered; these two are its siblings.
			{Name: FieldUpdated, Kind: graph.KindTime, Merge: graph.Precedence("rdap", "whois")},
			{Name: FieldExpires, Kind: graph.KindTime, Merge: graph.Precedence("rdap", "whois")},
			{Name: FieldRegistrar, Kind: graph.KindString, Merge: graph.Precedence("rdap", "whois")},
		},
		TypeIP: {
			{Name: FieldIPVersion, Kind: graph.KindInt},
			{Name: FieldCity, Kind: graph.KindString},
			{Name: FieldCountry, Kind: graph.KindString},
			{Name: FieldCountryCode, Kind: graph.KindString},
			{Name: FieldContinent, Kind: graph.KindString},
			{Name: FieldLatitude, Kind: graph.KindFloat},
			{Name: FieldLongitude, Kind: graph.KindFloat},
			{Name: FieldTimeZone, Kind: graph.KindString},
			{Name: FieldPostalCode, Kind: graph.KindString},
		},
		TypePlatform: {
			{Name: FieldCode, Kind: graph.KindString},
		},
	}
}

// RelFields returns the props these operators assert on observation edges,
// grouped by relation. The relations themselves are registered by decompose,
// with empty field lists; these are the append.
//
// Relation props exist because a relation carries data of its own: an MX
// preference describes the delivery order, not the mail host, and would be a
// lie as a prop on the host node.
func RelFields() map[string][]graph.FieldDef {
	return map[string][]graph.FieldDef{
		RelMX: {
			{Name: FieldPreference, Kind: graph.KindInt},
		},
		RelExistsOn: {
			// The page a human should be shown, which is often not the endpoint
			// the existence check used.
			{Name: FieldURL, Kind: graph.KindString},
		},
	}
}
